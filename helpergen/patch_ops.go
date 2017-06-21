package helpergen

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	//"github.com/hashicorp/terraform/helper/schema"
	u "github.com/radeksimko/terraform-gen/internal/util"
)

func (hg *HelperGenerator) PatchOpsFromStruct(iface interface{}) map[string]string {
	hg.init()
	hg.generatePatchOpsFromStruct(iface)
	return hg.renderDeclarations()
}

func (hg *HelperGenerator) generatePatchOpsFromStruct(iface interface{}) string {
	t := reflect.TypeOf(iface)
	rawType := getRawType(t)

	funcName := patchOpsFuncNameFromType(t)
	funcBody := hg.patchOpsBodyBeginning(t)

	log.Printf("Found %d fields for %q", rawType.NumField(), rawType.String())
	for i := 0; i < rawType.NumField(); i++ {
		sf := rawType.Field(i)
		log.Printf("Processing %q (%q)", sf.Name, sf.Type.String())
		body, err := hg.patchOpsField(sf.Name, sf.Type, iface, &sf)
		if err != nil {
			log.Printf("Skipping %s: %s", sf.Name, err)
			continue
		}
		funcBody += body
	}

	funcBody += hg.patchOpsBodyEnd(t)
	args := "keyPrefix, pathPrefix string, d *schema.ResourceData"

	hg.declarations[funcName] = &FunctionDeclaration{
		PkgPath:   t.PkgPath(),
		FuncName:  funcName,
		Arguments: args,
		Outputs:   "PatchOperations",
		FuncBody:  funcBody,
	}

	return funcName
}

func (hg *HelperGenerator) patchOpsField(sfName string, sfType reflect.Type, iface interface{}, sf *reflect.StructField) (string, error) {
	rawType := u.DereferencePtrType(sfType)
	kind := rawType.Kind()
	//s := &schema.Schema{}

	// if sf != nil {
	// 	var ok bool
	// 	kind, ok = hg.PatchFieldFilterFunc(iface, sf, kind, s)
	// 	if !ok {
	// 		return "", fmt.Errorf("Skipping %q (patch filter)", sf.Name)
	// 	}
	// }

	return hg.patchOpsFieldValue(kind, sf, sfName, sfType)
}

func (hg *HelperGenerator) patchOpsFieldValue(kind reflect.Kind, sf *reflect.StructField, sfName string, sfType reflect.Type) (string, error) {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		castType := "int"
		return fmt.Sprintf(`if d.HasChange(keyPrefix+%q) {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"%s",
Value: d.Get(%q).(%v),
})
}
`, u.Underscore(sf.Name), sfName, u.Underscore(sf.Name), castType), nil
	case reflect.Float32, reflect.Float64:
		castType := "float64"
		return fmt.Sprintf(`if d.HasChange(keyPrefix+%q) {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"%s",
Value: d.Get(%q).(%v),
})
}
`, u.Underscore(sf.Name), sfName, u.Underscore(sf.Name), castType), nil
	case reflect.String, reflect.Bool:
		castType := sfType.String()
		if sfType.Kind() == reflect.Ptr {
			castType = sfType.Elem().String()
		}
		return fmt.Sprintf(`if d.HasChange(keyPrefix+%q) {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"%s",
Value: d.Get(%q).(%v),
})
}
`, u.Underscore(sf.Name), sfName, u.Underscore(sf.Name), castType), nil
	case reflect.Map:
		funcName := hg.patchMapFuncForType(sfType)
		code := fmt.Sprintf(`if d.HasChange(keyPrefix+%q) {
oldV, newV := d.GetChange(keyPrefix+%q)
diffOps := %s(pathPrefix+"%s/", oldV.(map[string]interface{}), newV.(map[string]interface{}))
ops = append(ops, diffOps...)
}
`, u.Underscore(sf.Name), u.Underscore(sf.Name), funcName, sfName)
		return code, nil
		// case reflect.Slice:
		// 	// TODO: s.Type == TypeSet
		// 	sliceOf := sfType.Elem()
		// 	switch sliceOf.Kind() {
		// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		// 		reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		// 		// Slice of primitive data types
		// 		funcName := hg.primitiveSlicePatchOpsForType(sliceOf, sfType)
		// 		return fmt.Sprintf("%s[%q].([]interface{})", hg.InputVarName, u.Underscore(sf.Name)), nil
		// 	case reflect.Ptr:
		// 		ptrTo := sliceOf.Elem()
		// 		funcName := hg.primitiveSlicePatchOpsForType(ptrTo, sfType)
		// 		return fmt.Sprintf("%s[%q].([]interface{})", hg.InputVarName, u.Underscore(sf.Name)), nil
		// 	case reflect.Struct:
		// 		iface := reflect.New(sfType).Elem().Interface()
		// 		funcName := hg.generatePatchOpsFromStruct(iface)
		// 		return fmt.Sprintf("%s[%q].([]interface{})", hg.InputVarName, u.Underscore(sf.Name)), nil
		// 	}
		// case reflect.Struct:
		// 	iface := reflect.New(sfType).Elem().Interface()
		// 	funcName := hg.generatePatchOpsFromStruct(iface)
		// 	return fmt.Sprintf("%s[%q].([]interface{})", hg.InputVarName, u.Underscore(sf.Name)), nil
	}

	f := fmt.Sprintf("%s %s\n", sfName, sfType.String())
	return "", fmt.Errorf("Unable to process: %s", f)
}

func (hg *HelperGenerator) patchOpsBodyBeginning(t reflect.Type) string {
	return "ops := make([]PatchOperation, 0, 0)\n"
}

func (hg *HelperGenerator) patchOpsBodyEnd(t reflect.Type) string {
	return "return ops"
}

func (hg *HelperGenerator) patchMapFuncForType(sfType reflect.Type) string {
	t := sfType.Elem()
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "diffIntMap"
	case reflect.Float32, reflect.Float64:
		return "diffFloatMap"
	case reflect.String:
		return "diffStringMap"
	case reflect.Bool:
		return "diffBoolMap"
		// case reflect.Struct:
		// 	iface := reflect.New(sfType).Elem().Interface()
		// 	return hg.generatePatchOpsFromStruct(iface)
	}
	return "unknown"
}

func patchOpsFuncNameFromType(t reflect.Type) string {
	// pkg.TypeName
	parts := strings.Split(t.String(), ".")
	rawTypeName := parts[1]
	return "patch" + rawTypeName
}
