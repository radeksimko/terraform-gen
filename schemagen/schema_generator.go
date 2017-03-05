package schemagen

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"sort"
	"text/template"

	"github.com/hashicorp/terraform/helper/schema"
	u "github.com/radeksimko/terraform-gen/internal/util"
)

type getDocsFunc func(iface interface{}, sf *reflect.StructField) string
type filterFunc func(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool)

type SchemaGenerator struct {
	DocsFunc   getDocsFunc
	FilterFunc filterFunc
}

func (g *SchemaGenerator) FromStruct(iface interface{}) map[string]string {
	rawType := u.DereferencePtrType(reflect.TypeOf(iface))
	fields := make(map[string]string, 0)

	for i := 0; i < rawType.NumField(); i++ {
		sf := rawType.Field(i)

		content, err := g.generateField(sf.Name, sf.Type, iface, &sf, false)
		if err != nil {
			log.Printf("ERROR: %s", err)
		} else {
			fields[u.Underscore(sf.Name)] = content
		}
	}

	return fields
}

func (g *SchemaGenerator) generateField(sfName string, sfType reflect.Type, iface interface{}, sf *reflect.StructField, isNested bool) (string, error) {
	kind := u.DereferencePtrType(sfType).Kind()
	var comment, setFunc string
	s := &schema.Schema{}

	if sf != nil {
		var ok bool
		kind, ok = g.FilterFunc(iface, sf, kind, s)
		if !ok {
			return "", fmt.Errorf("Skipping %q (filter)", sf.Name)
		}
		comment = g.DocsFunc(iface, sf)
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s.Type = schema.TypeInt
	case reflect.Float32, reflect.Float64:
		s.Type = schema.TypeFloat
	case reflect.String:
		s.Type = schema.TypeString
	case reflect.Bool:
		s.Type = schema.TypeBool
	case reflect.Slice:
		// TODO: TypeList may be more suitable for some situations
		// TODO: Proper SetFunc may be required for TypeSet
		s.Type = schema.TypeSet
		elem, err := g.generateField("", sfType.Elem(), iface, nil, true)
		if err != nil {
			return "", fmt.Errorf("Unable to generate Elem for %q: %s", sfName, err)
		}
		s.Elem = elem

		elemKind := u.DereferencePtrType(sfType.Elem()).Kind()
		if elemKind == reflect.String {
			setFunc = "schema.HashString"
		}
	case reflect.Map:
		s.Type = schema.TypeMap
	case reflect.Struct:
		structType := sfType
		if structType.Kind() == reflect.Ptr {
			structType = structType.Elem()
		}

		s.Type = schema.TypeList
		s.MaxItems = 1

		elem := "&schema.Resource{\nSchema: map[string]*schema.Schema{\n"

		iface := reflect.New(structType).Elem().Interface()

		m := g.FromStruct(iface)
		fieldNames := make([]string, len(m), len(m))
		i := 0
		for k, _ := range m {
			fieldNames[i] = k
			i++
		}
		sort.Strings(fieldNames)
		for _, k := range fieldNames {
			elem += fmt.Sprintf("%q: %s,\n", k, m[k])
		}
		elem += "},\n}"
		if isNested {
			return elem, nil
		}

		s.Elem = elem
	default:
		f := fmt.Sprintf("%s %s\n", sfName, sfType.String())
		return "", fmt.Errorf("Unable to process: %s", f)
	}

	s.Description = comment

	return schemaCode(s, setFunc, isNested)
}

func schemaCode(s *schema.Schema, setFunc string, isNested bool) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	err := schemaTemplate.Execute(buf, struct {
		Schema   *schema.Schema
		SetFunc  string
		IsNested bool
	}{
		Schema:   s,
		SetFunc:  setFunc,
		IsNested: isNested,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

var schemaTemplate = template.Must(template.New("schema").Parse(`{{if .IsNested}}&schema.Schema{{end}}{{"{"}}{{if not .IsNested}}
{{end}}Type: schema.{{.Schema.Type}},{{if ne .Schema.Description ""}}
Description: {{printf "%q" .Schema.Description}},{{end}}{{if .Schema.Required}}
Required: {{.Schema.Required}},{{end}}{{if .Schema.Optional}}
Optional: {{.Schema.Optional}},{{end}}{{if .Schema.ForceNew}}
ForceNew: {{.Schema.ForceNew}},{{end}}{{if .Schema.Computed}}
Computed: {{.Schema.Computed}},{{end}}{{if gt .Schema.MaxItems 0}}
MaxItems: {{.Schema.MaxItems}},{{end}}{{if .Schema.Elem}}
Elem: {{.Schema.Elem}},{{end}}{{if ne .SetFunc ""}}{{if not .IsNested}}
{{end}}Set: {{.SetFunc}},{{end}}{{if not .IsNested}}
{{end}}{{"}"}}`))
