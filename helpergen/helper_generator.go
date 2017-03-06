package helpergen

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform/helper/schema"
	u "github.com/radeksimko/terraform-gen/internal/util"
)

type FunctionDeclaration struct {
	PkgPath   string
	FuncName  string
	Arguments string
	Outputs   string
	FuncBody  string
}

type filterFunc func(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool)

type HelperGenerator struct {
	FilterFunc    filterFunc
	InputVarName  string
	OutputVarName string

	declarations map[string]*FunctionDeclaration
}

func (hg *HelperGenerator) FromStruct(iface interface{}) map[string]string {
	if hg.declarations == nil {
		hg.declarations = make(map[string]*FunctionDeclaration)
	}
	if hg.FilterFunc == nil {
		hg.FilterFunc = noFilter
	}

	hg.generateDeclarationsFromStruct(iface)

	m := make(map[string]string)
	for name, decl := range hg.declarations {
		buf := bytes.NewBuffer([]byte{})
		err := flattenerFuncTpl.Execute(buf, decl)
		if err != nil {
			log.Fatal(err)
		}
		m[name] = buf.String()
	}
	return m
}

func (hg *HelperGenerator) generateDeclarationsFromStruct(iface interface{}) string {
	t := reflect.TypeOf(iface)
	rawType := getRawType(t)

	funcName := flattenerFuncNameFromType(t)
	funcBody := hg.OutputVarName + " := make(map[string]interface{})\n"

	for i := 0; i < rawType.NumField(); i++ {
		sf := rawType.Field(i)
		body, err := hg.generateFlattenerFieldCode(sf.Name, sf.Type, iface, &sf, false)
		if err != nil {
			log.Printf("Skipping %s: %s", sf.Name, err)
			continue
		}
		funcBody += body
	}

	funcBody += "return " + returnCodeFromType(hg.OutputVarName, t)

	hg.declarations[funcName] = &FunctionDeclaration{
		PkgPath:   t.PkgPath(),
		FuncName:  funcName,
		Arguments: argumentSignatureFromType(hg.InputVarName, t),
		Outputs:   returnInterfacesFromType(t),
		FuncBody:  funcBody,
	}

	return funcName
}

func (hg *HelperGenerator) generateFlattenerFieldCode(sfName string, sfType reflect.Type, iface interface{}, sf *reflect.StructField, isNested bool) (string, error) {
	kind := u.DereferencePtrType(sfType).Kind()
	s := &schema.Schema{}

	if sf != nil {
		var ok bool
		kind, ok = hg.FilterFunc(iface, sf, kind, s)
		if !ok {
			return "", fmt.Errorf("Skipping %q (filter)", sf.Name)
		}
	}

	value := "// TODO"
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		// Primitive data types are easy
		sfPtr := ""
		if sfType.Kind() == reflect.Ptr {
			sfPtr = "*"
		}
		value = fmt.Sprintf("%s%s.%s", sfPtr, hg.InputVarName, sf.Name)
	case reflect.Map:
		// TODO: map[string]*string
		value = fmt.Sprintf("%s.%s", hg.InputVarName, sf.Name)
	// TODO: case reflect.Slice:
	case reflect.Struct:
		iface := reflect.New(sfType).Elem().Interface()
		funcName := hg.generateDeclarationsFromStruct(iface)
		value = fmt.Sprintf("%s(%s.%s)", funcName, hg.InputVarName, sf.Name)
	default:
		f := fmt.Sprintf("%s %s\n", sfName, sfType.String())
		return "", fmt.Errorf("Unable to process: %s", f)
	}

	// TODO: Optional

	return fmt.Sprintf("%s[%q] = %s\n",
		hg.OutputVarName, u.Underscore(sf.Name), value), nil
}

func flattenerFuncNameFromType(t reflect.Type) string {
	// pkg.TypeName
	parts := strings.Split(t.String(), ".")
	rawTypeName := parts[1]
	return "flatten" + rawTypeName
}

func getRawType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Slice {
		return getRawType(t.Elem())
	}
	if t.Kind() == reflect.Ptr {
		return getRawType(t.Elem())
	}
	return t
}

func argumentSignatureFromType(argName string, t reflect.Type) string {
	ptr := ""
	slice := ""
	if t.Kind() == reflect.Slice {
		slice = "[]"
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		ptr = "*"
		t = t.Elem()
	}
	return argName + " " + slice + ptr + t.String()
}

func returnCodeFromType(varName string, t reflect.Type) string {
	if t.Kind() == reflect.Slice {
		return "[]map[string]interface{}{" + varName + "}"
	}
	return varName
}

func returnInterfacesFromType(t reflect.Type) string {
	if t.Kind() == reflect.Slice {
		return "[]map[string]interface{}"
	}
	return "map[string]interface{}"
}

func noFilter(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool) {
	return k, true
}

var flattenerFuncTpl = template.Must(template.New("func-decl").Parse(`func {{.FuncName}}({{.Arguments}}) {{.Outputs}} {{"{"}}
{{.FuncBody}}
}`))
