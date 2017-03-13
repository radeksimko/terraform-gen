package helpergen

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	u "github.com/radeksimko/terraform-gen/internal/util"
)

func (hg *HelperGenerator) ExpandersFromStruct(iface interface{}) map[string]string {
	hg.init()
	hg.generateExpandersFromStruct(iface)
	return hg.renderDeclarations()
}

func (hg *HelperGenerator) generateExpandersFromStruct(iface interface{}) string {
	t := reflect.TypeOf(iface)
	rawType := getRawType(t)

	funcName := expanderFuncNameFromType(t)
	funcBody := hg.expanderDeclarationBeginning(t)
	for i := 0; i < rawType.NumField(); i++ {
		sf := rawType.Field(i)
		body, err := hg.generateExpanderFieldCode(sf.Name, sf.Type, iface, &sf)
		if err != nil {
			log.Printf("Skipping %s: %s", sf.Name, err)
			continue
		}
		funcBody += body
	}
	funcBody += hg.expanderDeclarationEnd(t)

	args := "l" + " []interface{}"

	hg.declarations[funcName] = &FunctionDeclaration{
		PkgPath:   t.PkgPath(),
		FuncName:  funcName,
		Arguments: args,
		Outputs:   interfaceFromType(t),
		FuncBody:  funcBody,
	}

	return funcName
}

func (hg *HelperGenerator) generateExpanderFieldCode(sfName string, sfType reflect.Type, iface interface{}, sf *reflect.StructField) (string, error) {
	kind := u.DereferencePtrType(sfType).Kind()
	s := &schema.Schema{}

	if sf != nil {
		var ok bool
		kind, ok = hg.FieldFilterFunc(iface, sf, kind, s)
		if !ok {
			return "", fmt.Errorf("Skipping %q (filter)", sf.Name)
		}
	}

	rightSide := "// TODO"

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		// Primitive data types are easy
		castType := sfType.String()
		rightSide = fmt.Sprintf("%s[%q].(%v)", hg.InputVarName, u.Underscore(sf.Name), castType)

		if sfType.Kind() == reflect.Ptr {
			castType = sfType.Elem().String()
			firstLetter := strings.ToUpper(string(castType[0]))
			ptrHelperFunc := "ptrTo" + firstLetter + castType[1:]
			rightSide = fmt.Sprintf("%s(%s[%q].(%v))", ptrHelperFunc, hg.InputVarName, u.Underscore(sf.Name), castType)
		}
	case reflect.Map:
		// TODO: map[string]*string
		// TODO: map[string]int
		// TODO: map[string]bool
		// TODO: map[string]float
		rightSide = fmt.Sprintf("expandStringMap(%s[%q].(map[string]interface{}))", hg.InputVarName, u.Underscore(sf.Name))
	case reflect.Slice:
		// TODO: s.Type == TypeSet
		sliceOf := sfType.Elem()
		switch sliceOf.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
			// Slice of primitive data types
			funcName := hg.primitiveSliceExpanderForType(sliceOf, sfType)
			rightSide = fmt.Sprintf("%s(%s[%q].([]interface{}))", funcName, hg.InputVarName, u.Underscore(sf.Name))
		case reflect.Ptr:
			ptrTo := sliceOf.Elem()
			funcName := hg.primitiveSliceExpanderForType(ptrTo, sfType)
			rightSide = fmt.Sprintf("%s(%s[%q].([]interface{}))", funcName, hg.InputVarName, u.Underscore(sf.Name))
		case reflect.Struct:
			iface := reflect.New(sfType).Elem().Interface()
			funcName := hg.generateExpandersFromStruct(iface)
			rightSide = fmt.Sprintf("%s(%s[%q].([]interface{}))", funcName, hg.InputVarName, u.Underscore(sf.Name))
		default:
			f := fmt.Sprintf("%s %s\n", sfName, sfType.String())
			return "", fmt.Errorf("Unable to process: %s", f)
		}
	case reflect.Struct:
		iface := reflect.New(sfType).Elem().Interface()
		funcName := hg.generateExpandersFromStruct(iface)
		rightSide = fmt.Sprintf("%s(%s[%q].([]interface{}))", funcName, hg.InputVarName, u.Underscore(sf.Name))
	default:
		f := fmt.Sprintf("%s %s\n", sfName, sfType.String())
		return "", fmt.Errorf("Unable to process: %s", f)
	}

	leftSide := sf.Name

	return fmt.Sprintf("%s: %s,\n", leftSide, rightSide), nil
}

func (hg *HelperGenerator) expanderDeclarationBeginning(t reflect.Type) string {
	code := ""
	if t.Kind() == reflect.Slice {
		ptr := ""
		if t.Elem().Kind() == reflect.Ptr {
			ptr = "&"
		}
		code += `obj := make([]helpergen.NestedStruct, len(l), len(l))
for i, n := range l {
cfg := n.(map[string]interface{})
obj[i] = ` + ptr + `helpergen.NestedStruct{
`
		hg.mapVarName = "cfg"
		return code
	}

	// TODO: if len(l) > 0
	code += hg.InputVarName + " := l[0].(map[string]interface{})\n"
	if t.Kind() == reflect.Ptr {
		// Pointer will be created at return stage
		t = t.Elem()
	}
	return code + "obj := " + t.String() + "{\n"
}

func (hg *HelperGenerator) expanderDeclarationEnd(t reflect.Type) string {
	code := ""
	ptr := ""
	if t.Kind() == reflect.Ptr {
		ptr = "&"
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		code += "}\n"
	}
	return code + "}\nreturn " + ptr + "obj"
}

func (hg *HelperGenerator) primitiveSliceExpanderForType(t reflect.Type, sfType reflect.Type) string {
	sliceOf := sfType.Elem()
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if sliceOf.Kind() == reflect.Ptr {
			return "sliceOfPtrInt"
		}
		return "sliceOfInt"
	case reflect.Float32, reflect.Float64:
		if sliceOf.Kind() == reflect.Ptr {
			return "sliceOfPtrFloat"
		}
		return "sliceOfFloat"
	case reflect.String:
		if sliceOf.Kind() == reflect.Ptr {
			return "sliceOfPtrString"
		}
		return "sliceOfString"
	case reflect.Bool:
		if sliceOf.Kind() == reflect.Ptr {
			return "sliceOfPtrBool"
		}
		return "sliceOfBool"
	case reflect.Struct:
		iface := reflect.New(sfType).Elem().Interface()
		return hg.generateExpandersFromStruct(iface)
	}
	return ""
}

func expanderFuncNameFromType(t reflect.Type) string {
	// pkg.TypeName
	parts := strings.Split(t.String(), ".")
	rawTypeName := parts[1]
	return "expand" + rawTypeName
}
