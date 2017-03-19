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
	funcBody := hg.expanderBodyBeginning(t)

	// Inline fields (typically those we never expect to be empty)
	funcBody += hg.inlineExpanderDeclarationBeginning(t)
	for i := 0; i < rawType.NumField(); i++ {
		sf := rawType.Field(i)
		body, err := hg.inlineExpanderField(sf.Name, sf.Type, iface, &sf)
		if err != nil {
			log.Printf("Skipping %s (inline): %s", sf.Name, err)
			continue
		}
		funcBody += body
	}
	funcBody += hg.inlineExpanderDeclarationEnd(t)

	// Outline fields (typically optional)
	for i := 0; i < rawType.NumField(); i++ {
		sf := rawType.Field(i)
		body, err := hg.outlineExpanderField(sf.Name, sf.Type, iface, &sf)
		if err != nil {
			log.Printf("Skipping %s (outline): %s", sf.Name, err)
			continue
		}
		funcBody += body
	}

	funcBody += hg.expanderBodyEnd(t)
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

func (hg *HelperGenerator) inlineExpanderDeclarationBeginning(t reflect.Type) string {
	if t.Kind() == reflect.Slice {
		ptr := ""
		t := t.Elem()
		if t.Kind() == reflect.Ptr {
			ptr = "&"
			t = t.Elem()
		}
		return `obj[i] = ` + ptr + t.String() + "{\n"
	}

	ptr := ""
	if t.Kind() == reflect.Ptr {
		ptr = "&"
		t = t.Elem()
	}

	return "obj := " + ptr + t.String() + "{\n"
}

func (hg *HelperGenerator) inlineExpanderDeclarationEnd(t reflect.Type) string {
	return "}\n"
}

func (hg *HelperGenerator) inlineExpanderField(sfName string, sfType reflect.Type, iface interface{}, sf *reflect.StructField) (string, error) {
	rawType := u.DereferencePtrType(sfType)
	kind := rawType.Kind()
	s := &schema.Schema{}

	if sf != nil {
		var ok bool
		kind, ok = hg.InlineFieldFilterFunc(iface, sf, kind, s)
		if !ok {
			return "", fmt.Errorf("Skipping %q (inline filter)", sf.Name)
		}
	}

	wrapperFunc, value, err := hg.expanderFieldValue(kind, sf, sfName, sfType)
	if err != nil {
		return "", err
	}
	leftSide := sf.Name

	if wrapperFunc != "" {
		value = fmt.Sprintf("%s(%s)", wrapperFunc, value)
	}

	return fmt.Sprintf("%s: %s,\n", leftSide, value), nil
}

func (hg *HelperGenerator) outlineExpanderField(sfName string, sfType reflect.Type, iface interface{}, sf *reflect.StructField) (string, error) {
	rawType := u.DereferencePtrType(sfType)
	kind := rawType.Kind()
	s := &schema.Schema{}

	if sf != nil {
		var ok bool
		kind, ok = hg.OutlineFieldFilterFunc(iface, sf, kind, s)
		if !ok {
			return "", fmt.Errorf("Skipping %q (outline filter)", sf.Name)
		}
	}

	wrapperFunc, value, err := hg.expanderFieldValue(kind, sf, sfName, sfType)
	if err != nil {
		return "", err
	}
	leftSide := sf.Name
	assignedValue := "v"
	if wrapperFunc != "" {
		assignedValue = fmt.Sprintf("%s(v)", wrapperFunc)
	}

	lengthCondition := ""
	switch kind {
	case reflect.Struct, reflect.Slice, reflect.Map:
		lengthCondition = " && len(v) > 0"
	}

	return fmt.Sprintf(`if v, ok := %s; ok%s {
%s.%s = %s
}
`, value, lengthCondition, "obj", leftSide, assignedValue), nil
}

func (hg *HelperGenerator) expanderFieldValue(kind reflect.Kind, sf *reflect.StructField, sfName string, sfType reflect.Type) (string, string, error) {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		castType := sfType.String()

		if sfType.Kind() == reflect.Ptr {
			castType = sfType.Elem().String()
			firstLetter := strings.ToUpper(string(castType[0]))
			ptrHelperFunc := "ptrTo" + firstLetter + castType[1:]
			return ptrHelperFunc, fmt.Sprintf("%s[%q].(%v)", hg.InputVarName, u.Underscore(sf.Name), castType), nil
		}

		return "", fmt.Sprintf("%s[%q].(%v)", hg.InputVarName, u.Underscore(sf.Name), castType), nil
	case reflect.Map:
		// TODO: map[string]*string
		// TODO: map[string]int
		// TODO: map[string]bool
		// TODO: map[string]float
		return "expandStringMap", fmt.Sprintf("%s[%q].(map[string]interface{})", hg.InputVarName, u.Underscore(sf.Name)), nil
	case reflect.Slice:
		// TODO: s.Type == TypeSet
		sliceOf := sfType.Elem()
		switch sliceOf.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
			// Slice of primitive data types
			funcName := hg.primitiveSliceExpanderForType(sliceOf, sfType)
			return funcName, fmt.Sprintf("%s[%q].([]interface{})", hg.InputVarName, u.Underscore(sf.Name)), nil
		case reflect.Ptr:
			ptrTo := sliceOf.Elem()
			funcName := hg.primitiveSliceExpanderForType(ptrTo, sfType)
			return funcName, fmt.Sprintf("%s[%q].([]interface{})", hg.InputVarName, u.Underscore(sf.Name)), nil
		case reflect.Struct:
			iface := reflect.New(sfType).Elem().Interface()
			funcName := hg.generateExpandersFromStruct(iface)
			return funcName, fmt.Sprintf("%s[%q].([]interface{})", hg.InputVarName, u.Underscore(sf.Name)), nil
		}
	case reflect.Struct:
		iface := reflect.New(sfType).Elem().Interface()
		funcName := hg.generateExpandersFromStruct(iface)
		return funcName, fmt.Sprintf("%s[%q].([]interface{})", hg.InputVarName, u.Underscore(sf.Name)), nil
	}

	f := fmt.Sprintf("%s %s\n", sfName, sfType.String())
	return "", "", fmt.Errorf("Unable to process: %s", f)
}

func (hg *HelperGenerator) expanderBodyBeginning(t reflect.Type) string {
	code := ""
	if t.Kind() == reflect.Slice {
		code += `if len(l) == 0 || l[0] == nil {
return ` + t.String() + `{}
}
obj := make(` + t.String() + `, len(l), len(l))
for i, n := range l {
cfg := n.(map[string]interface{})
`
		hg.mapVarName = "cfg"
		return code
	}

	ptr := ""
	if t.Kind() == reflect.Ptr {
		ptr = "&"
		t = t.Elem()
	}

	code += `if len(l) == 0 || l[0] == nil {
return ` + ptr + t.String() + `{}
}
` + hg.InputVarName + " := l[0].(map[string]interface{})\n"

	return code
}

func (hg *HelperGenerator) expanderBodyEnd(t reflect.Type) string {
	code := ""
	if t.Kind() == reflect.Slice {
		code += "}\n"
	}
	return code + "return obj"
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
