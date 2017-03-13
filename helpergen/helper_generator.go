package helpergen

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"text/template"

	"github.com/hashicorp/terraform/helper/schema"
)

type FunctionDeclaration struct {
	PkgPath   string
	FuncName  string
	Arguments string
	Outputs   string
	FuncBody  string
}

type fieldFilterFunc func(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool)

type HelperGenerator struct {
	FieldFilterFunc fieldFilterFunc
	InputVarName    string
	OutputVarName   string

	mapVarName   string
	mapValueName string
	declarations map[string]*FunctionDeclaration
}

func (hg *HelperGenerator) init() {
	hg.declarations = make(map[string]*FunctionDeclaration)
	if hg.FieldFilterFunc == nil {
		hg.FieldFilterFunc = noFieldFilter
	}
	if hg.mapVarName == "" {
		hg.mapVarName = hg.OutputVarName
	}
	if hg.mapValueName == "" {
		hg.mapValueName = hg.InputVarName
	}
}

func (hg *HelperGenerator) renderDeclarations() map[string]string {
	m := make(map[string]string)
	for name, decl := range hg.declarations {
		buf := bytes.NewBuffer([]byte{})
		err := funcDeclTpl.Execute(buf, decl)
		if err != nil {
			log.Fatal(err)
		}
		m[name] = buf.String()
	}
	return m
}

func emptyConditionForType(inputVarName string, sf *reflect.StructField) (string, error) {
	leftSide := inputVarName + "." + sf.Name

	switch sf.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		val := reflect.Zero(sf.Type)
		return fmt.Sprintf("%s == %v", leftSide, val.Interface()), nil
	case reflect.String:
		return fmt.Sprintf(`%s == ""`, leftSide), nil
	case reflect.Ptr:
		return fmt.Sprintf("%s == nil", leftSide), nil
	case reflect.Slice, reflect.Map:
		return fmt.Sprintf("len(%s) > 0", leftSide), nil
	}

	f := fmt.Sprintf("%s\n", sf.Type.String())
	return fmt.Sprintf(`%s == /* unknown */`, leftSide), fmt.Errorf("Unable to process: %s", f)
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

func interfaceFromType(t reflect.Type) string {
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
	return slice + ptr + t.String()
}

func mapInterfacesFromType(t reflect.Type) string {
	return "[]interface{}"
}

func noFieldFilter(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool) {
	return k, true
}

var funcDeclTpl = template.Must(template.New("func-decl").Parse(`func {{.FuncName}}({{.Arguments}}) {{.Outputs}} {{"{"}}
{{.FuncBody}}
}`))
