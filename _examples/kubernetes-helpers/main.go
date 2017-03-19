package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/radeksimko/terraform-gen/helpergen"

	api "k8s.io/kubernetes/pkg/api/v1"
)

type helperGen struct {
	Obj      interface{}
	Filename string
}

func main() {
	pkgName := "kubernetes"
	schemas := []helperGen{
		{
			Obj:      api.PersistentVolumeSpec{},
			Filename: "structure_persistent_volume_spec.go",
		},
	}

	for _, s := range schemas {
		log.Printf("Generating %q...\n", s.Filename)
		f, err := os.Create(s.Filename)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}

		hg := &helpergen.HelperGenerator{
			InputVarName:           "in",
			OutputVarName:          "att",
			InlineFieldFilterFunc:  inlineFilterFunc,
			OutlineFieldFilterFunc: outlineFilterFunc,
		}

		flatteners := hg.FlattenersFromStruct(s.Obj)
		expanders := hg.ExpandersFromStruct(s.Obj)

		err = tpl.Execute(f, struct {
			PkgName    string
			Flatteners map[string]string
			Expanders  map[string]string
		}{
			PkgName:    pkgName,
			Flatteners: flatteners,
			Expanders:  expanders,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func inlineFilterFunc(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool) {
	tag := sf.Tag
	jsonTag := strings.Split(tag.Get("json"), ",")
	jsonName := jsonTag[0]

	if jsonName == "-" {
		return k, false
	}

	docs, err := getSwaggerDocs(iface, sf, jsonName)
	if err != nil {
		return k, true
	}

	if strings.Contains(docs, "Deprecated:") {
		log.Printf("Ignoring %q (deprecated)\n", sf.Name)
		return k, false
	}
	if strings.Contains(docs, "NOT YET IMPLEMENTED.") {
		log.Printf("Ignoring %q (not implemented)\n", sf.Name)
		return k, false
	}

	if strings.Contains(docs, "Read-only.") {
		s.Computed = true
	} else {
		if strings.Contains(docs, "Required.") || strings.Contains(docs, "Required:") {
			s.Required = true
		} else {
			s.Optional = true
		}
	}
	return k, s.Required
}

func outlineFilterFunc(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool) {
	tag := sf.Tag
	jsonTag := strings.Split(tag.Get("json"), ",")
	jsonName := jsonTag[0]

	if jsonName == "-" {
		return k, false
	}

	docs, err := getSwaggerDocs(iface, sf, jsonName)
	if err != nil {
		return k, true
	}

	if strings.Contains(docs, "Deprecated:") {
		log.Printf("Ignoring %q (deprecated)\n", sf.Name)
		return k, false
	}
	if strings.Contains(docs, "NOT YET IMPLEMENTED.") {
		log.Printf("Ignoring %q (not implemented)\n", sf.Name)
		return k, false
	}

	if strings.Contains(docs, "Read-only.") {
		s.Computed = true
	} else {
		if strings.Contains(docs, "Required.") || strings.Contains(docs, "Required:") {
			s.Required = true
		} else {
			s.Optional = true
		}
	}
	return k, s.Optional
}

func getSwaggerDocs(iface interface{}, sf *reflect.StructField, jsonName string) (string, error) {
	structType := reflect.TypeOf(iface)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	val := reflect.New(structType).Elem()
	method := val.MethodByName("SwaggerDoc")
	if method.IsValid() {
		out := method.Call([]reflect.Value{})
		m := out[0]
		docs := m.MapIndex(reflect.ValueOf(jsonName))
		if !docs.IsValid() {
			docs = m.MapIndex(reflect.ValueOf(""))
		}
		return docs.String(), nil
	}
	return "", fmt.Errorf("Docs not found for %s -> %s (%s)", structType.Name(), sf.Name, sf.Type.String())
}

var tpl = template.Must(template.New("pod").Parse(`package {{.PkgName}}

// Flatteners
{{range $name, $definition := .Flatteners}}
{{ $definition }}
{{end}}

// Expanders
{{range $name, $definition := .Expanders}}
{{ $definition }}
{{end}}
`))
