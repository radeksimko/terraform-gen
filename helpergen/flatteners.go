package helpergen

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	u "github.com/radeksimko/terraform-gen/internal/util"
)

func (hg *HelperGenerator) FlattenersFromStruct(iface interface{}) map[string]string {
	hg.init()
	hg.generateFlattenersFromStruct(iface)
	return hg.renderDeclarations()
}

func (hg *HelperGenerator) generateFlattenersFromStruct(iface interface{}) string {
	t := reflect.TypeOf(iface)
	rawType := getRawType(t)

	funcName := flattenerFuncNameFromType(t)
	funcBody := hg.flattenerDeclarationBeginning(t)

	// Inline fields (typically those we never expect to be empty)
	for i := 0; i < rawType.NumField(); i++ {
		sf := rawType.Field(i)
		body, err := hg.inlineFlattenerField(sf.Name, sf.Type, iface, &sf, false)
		if err != nil {
			log.Printf("Skipping %s (inline): %s", sf.Name, err)
			continue
		}
		funcBody += body
	}

	// Outline fields (typically optional)
	for i := 0; i < rawType.NumField(); i++ {
		sf := rawType.Field(i)
		body, err := hg.outlineFlattenerField(sf.Name, sf.Type, iface, &sf, false)
		if err != nil {
			log.Printf("Skipping %s (outline): %s", sf.Name, err)
			continue
		}
		funcBody += body
	}

	funcBody += hg.flattenerDeclarationEnd(t)

	hg.declarations[funcName] = &FunctionDeclaration{
		PkgPath:   t.PkgPath(),
		FuncName:  funcName,
		Arguments: hg.InputVarName + " " + interfaceFromType(t),
		Outputs:   mapInterfacesFromType(t),
		FuncBody:  funcBody,
	}

	return funcName
}

func (hg *HelperGenerator) flattenerDeclarationBeginning(t reflect.Type) string {
	if t.Kind() == reflect.Slice {
		body := hg.mapVarName + ` := make([]interface{}, len(in), len(in))
for i, n := range in {
m := make(map[string]interface{})
`
		hg.mapVarName = "m"
		hg.mapValueName = "n"
		return body
	}

	return hg.mapVarName + " := make(map[string]interface{})\n"
}

func (hg *HelperGenerator) flattenerDeclarationEnd(t reflect.Type) string {
	if t.Kind() == reflect.Slice {
		body := hg.OutputVarName + `[i] = ` + hg.mapVarName + `
}
return ` + hg.OutputVarName
		hg.mapVarName = ""
		hg.mapValueName = ""
		return body
	}

	return `return []interface{}{` + hg.OutputVarName + `}`
}

func (hg *HelperGenerator) inlineFlattenerField(sfName string, sfType reflect.Type, iface interface{}, sf *reflect.StructField, isNested bool) (string, error) {
	kind := u.DereferencePtrType(sfType).Kind()
	s := &schema.Schema{}

	if sf != nil {
		var ok bool
		kind, ok = hg.InlineFieldFilterFunc(iface, sf, kind, s)
		if !ok {
			return "", fmt.Errorf("Skipping %q (inline filter)", sf.Name)
		}
	}

	value, err := hg.flattenerFieldValue(kind, sf, sfName, sfType)
	if err != nil {
		return "", err
	}

	leftSide := fmt.Sprintf("%s[%q]", hg.OutputVarName, u.Underscore(sf.Name))
	if hg.mapVarName != "" {
		leftSide = fmt.Sprintf("%s[%q]", hg.mapVarName, u.Underscore(sf.Name))
	}

	return fmt.Sprintf("%s = %s\n", leftSide, value), nil
}

func (hg *HelperGenerator) outlineFlattenerField(sfName string, sfType reflect.Type, iface interface{}, sf *reflect.StructField, isNested bool) (string, error) {
	kind := u.DereferencePtrType(sfType).Kind()
	s := &schema.Schema{}

	if sf != nil {
		var ok bool
		kind, ok = hg.OutlineFieldFilterFunc(iface, sf, kind, s)
		if !ok {
			return "", fmt.Errorf("Skipping %q (outline filter)", sf.Name)
		}
	}

	value, err := hg.flattenerFieldValue(kind, sf, sfName, sfType)
	if err != nil {
		return "", err
	}

	leftSide := fmt.Sprintf("%s[%q]", hg.OutputVarName, u.Underscore(sf.Name))
	if hg.mapVarName != "" {
		leftSide = fmt.Sprintf("%s[%q]", hg.mapVarName, u.Underscore(sf.Name))
	}

	if s.Optional || s.Computed {
		inputVarName := hg.InputVarName
		if hg.mapValueName != "" {
			inputVarName = hg.mapValueName
		}
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

func (hg *HelperGenerator) flattenerFieldValue(kind reflect.Kind, sf *reflect.StructField, sfName string, sfType reflect.Type) (string, error) {
	inputVarName := hg.InputVarName
	if hg.mapValueName != "" {
		inputVarName = hg.mapValueName
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		// Primitive data types are easy
		sfPtr := ""
		if sfType.Kind() == reflect.Ptr {
			sfPtr = "*"
		}
		return fmt.Sprintf("%s%s.%s", sfPtr, inputVarName, sf.Name), nil
	case reflect.Map:
		// TODO: map[string]*string
		// TODO: map[string]*string
		// TODO: map[string]int
		// TODO: map[string]bool
		// TODO: map[string]float
		return fmt.Sprintf("%s.%s", inputVarName, sf.Name), nil
	case reflect.Slice:
		// TODO: s.Type == TypeSet
		sliceOf := sfType.Elem()
		switch sliceOf.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
			// Slice of primitive data types
			return fmt.Sprintf("%s.%s", inputVarName, sf.Name), nil
		case reflect.Ptr:
			ptrTo := sliceOf.Elem()
			funcName := hg.primitivePtrSliceFlattenerForType(ptrTo, sfType)
			return fmt.Sprintf("%s(%s.%s)", funcName, inputVarName, sf.Name), nil
		case reflect.Struct:
			iface := reflect.New(sfType).Elem().Interface()
			funcName := hg.generateFlattenersFromStruct(iface)
			return fmt.Sprintf("%s(%s.%s)", funcName, inputVarName, sf.Name), nil
		}
	case reflect.Struct:
		iface := reflect.New(sfType).Elem().Interface()
		funcName := hg.generateFlattenersFromStruct(iface)
		return fmt.Sprintf("%s(%s.%s)", funcName, inputVarName, sf.Name), nil
	}

	f := fmt.Sprintf("%s %s\n", sfName, sfType.String())
	return "", fmt.Errorf("Unable to process: %s", f)
}

func (hg *HelperGenerator) primitivePtrSliceFlattenerForType(t reflect.Type, sfType reflect.Type) string {
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
		return hg.generateFlattenersFromStruct(iface)
	}
	return ""
}

func flattenerFuncNameFromType(t reflect.Type) string {
	// pkg.TypeName
	parts := strings.Split(t.String(), ".")
	rawTypeName := parts[1]
	return "flatten" + rawTypeName
}
