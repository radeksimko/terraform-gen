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

type fieldFilterFunc func(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool)

type HelperGenerator struct {
	FieldFilterFunc fieldFilterFunc
	InputVarName    string
	OutputVarName   string

	mapVarName   string
	mapValueName string
	declarations map[string]*FunctionDeclaration
}

func (hg *HelperGenerator) FromStruct(iface interface{}) map[string]string {
	if hg.declarations == nil {
		hg.declarations = make(map[string]*FunctionDeclaration)
	}
	if hg.FieldFilterFunc == nil {
		hg.FieldFilterFunc = noFieldFilter
	}
	if hg.mapVarName == "" {
		hg.mapVarName = hg.OutputVarName
	}
	if hg.mapValueName == "" {
		hg.mapValueName = hg.InputVarName
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
	funcBody := hg.funcDeclarationBeginning(t)
	for i := 0; i < rawType.NumField(); i++ {
		sf := rawType.Field(i)
		body, err := hg.generateFlattenerFieldCode(sf.Name, sf.Type, iface, &sf, false)
		if err != nil {
			log.Printf("Skipping %s: %s", sf.Name, err)
			continue
		}
		funcBody += body
	}
	funcBody += hg.funcDeclarationEnd(t)

	hg.declarations[funcName] = &FunctionDeclaration{
		PkgPath:   t.PkgPath(),
		FuncName:  funcName,
		Arguments: argumentSignatureFromType(hg.InputVarName, t),
		Outputs:   returnInterfacesFromType(t),
		FuncBody:  funcBody,
	}

	return funcName
}

func (hg *HelperGenerator) funcDeclarationBeginning(t reflect.Type) string {
	if t.Kind() == reflect.Slice {
		body := hg.mapVarName + ` := make([]map[string]interface{}, len(in), len(in))
for i, n := range in {
m := make(map[string]interface{})
`
		hg.mapVarName = "m"
		hg.mapValueName = "n"
		return body
	}

	return hg.mapVarName + " := make(map[string]interface{})\n"
}

func (hg *HelperGenerator) funcDeclarationEnd(t reflect.Type) string {
	if t.Kind() == reflect.Slice {
		body := hg.OutputVarName + `[i] = ` + hg.mapVarName + `
}
return ` + hg.OutputVarName
		hg.mapVarName = ""
		hg.mapValueName = ""
		return body
	}

	return `return ` + hg.OutputVarName
}

func (hg *HelperGenerator) generateFlattenerFieldCode(sfName string, sfType reflect.Type, iface interface{}, sf *reflect.StructField, isNested bool) (string, error) {
	kind := u.DereferencePtrType(sfType).Kind()
	s := &schema.Schema{}

	if sf != nil {
		var ok bool
		kind, ok = hg.FieldFilterFunc(iface, sf, kind, s)
		if !ok {
			return "", fmt.Errorf("Skipping %q (filter)", sf.Name)
		}
	}

	inputVarName := hg.InputVarName
	if hg.mapValueName != "" {
		inputVarName = hg.mapValueName
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
		value = fmt.Sprintf("%s%s.%s", sfPtr, inputVarName, sf.Name)
	case reflect.Map:
		// TODO: map[string]*string
		value = fmt.Sprintf("%s.%s", inputVarName, sf.Name)
	case reflect.Slice:
		// TODO: s.Type == TypeSet
		sliceOf := sfType.Elem()
		switch sliceOf.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
			// Slice of primitive data types
			value = fmt.Sprintf("%s.%s", inputVarName, sf.Name)
		case reflect.Ptr:
			ptrTo := sliceOf.Elem()
			funcName := hg.slicePtrHelperFuncNameForType(ptrTo, sfType)
			value = fmt.Sprintf("%s(%s.%s)", funcName, inputVarName, sf.Name)
		case reflect.Struct:
			iface := reflect.New(sfType).Elem().Interface()
			funcName := hg.generateDeclarationsFromStruct(iface)
			value = fmt.Sprintf("%s(%s.%s)", funcName, inputVarName, sf.Name)
		default:
			f := fmt.Sprintf("%s %s\n", sfName, sfType.String())
			return "", fmt.Errorf("Unable to process: %s", f)
		}
	case reflect.Struct:
		iface := reflect.New(sfType).Elem().Interface()
		funcName := hg.generateDeclarationsFromStruct(iface)
		value = fmt.Sprintf("%s(%s.%s)", funcName, inputVarName, sf.Name)
	default:
		f := fmt.Sprintf("%s %s\n", sfName, sfType.String())
		return "", fmt.Errorf("Unable to process: %s", f)
	}

	leftSide := fmt.Sprintf("%s[%q]", hg.OutputVarName, u.Underscore(sf.Name))
	if hg.mapVarName != "" {
		leftSide = fmt.Sprintf("%s[%q]", hg.mapVarName, u.Underscore(sf.Name))
	}

	if s.Optional && !s.Computed {
		emptyValue, err := emptyConditionForType(inputVarName, sf)
		if err != nil {
			log.Printf("Unknown optional condition: %s", err)
		}
		body := fmt.Sprintf("if %s {\n", emptyValue)
		body += fmt.Sprintf("%s = %s\n", leftSide, value)
		body += "}\n"
		return body, nil
	}

	return fmt.Sprintf("%s = %s\n", leftSide, value), nil
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

func (hg *HelperGenerator) slicePtrHelperFuncNameForType(t reflect.Type, sfType reflect.Type) string {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "flattenIntSlice"
	case reflect.Float32, reflect.Float64:
		return "flattenFloatSlice"
	case reflect.String:
		return "flattenStringSlice"
	case reflect.Bool:
		return "flattenBoolSlice"
	case reflect.Struct:
		iface := reflect.New(sfType).Elem().Interface()
		return hg.generateDeclarationsFromStruct(iface)
	}
	return ""
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

func returnInterfacesFromType(t reflect.Type) string {
	if t.Kind() == reflect.Slice {
		return "[]map[string]interface{}"
	}
	return "map[string]interface{}"
}

func noFieldFilter(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool) {
	return k, true
}

var flattenerFuncTpl = template.Must(template.New("func-decl").Parse(`func {{.FuncName}}({{.Arguments}}) {{.Outputs}} {{"{"}}
{{.FuncBody}}
}`))
