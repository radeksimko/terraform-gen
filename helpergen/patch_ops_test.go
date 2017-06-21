package helpergen

import (
	"reflect"
	"testing"
)

func TestPatchOpsFromStruct_primitives(t *testing.T) {
	type SimpleStruct struct {
		MyInt    int
		MyString string
		MyFloat  float64
		MyBool   bool
	}
	hg := &HelperGenerator{
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.PatchOpsFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"patchSimpleStruct": `func patchSimpleStruct(keyPrefix, pathPrefix string, d *schema.ResourceData) PatchOperations {
ops := make([]PatchOperation, 0, 0)
if d.HasChange(keyPrefix+"my_int") {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"MyInt",
Value: d.Get("my_int").(int),
})
}
if d.HasChange(keyPrefix+"my_string") {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"MyString",
Value: d.Get("my_string").(string),
})
}
if d.HasChange(keyPrefix+"my_float") {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"MyFloat",
Value: d.Get("my_float").(float64),
})
}
if d.HasChange(keyPrefix+"my_bool") {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"MyBool",
Value: d.Get("my_bool").(bool),
})
}
return ops
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestPatchOpsFromStruct_ptrsToPrimitives(t *testing.T) {
	type SimpleStruct struct {
		MyInt    *int
		MyString *string
		MyFloat  *float64
		MyBool   *bool
	}
	hg := &HelperGenerator{
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.PatchOpsFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"patchSimpleStruct": `func patchSimpleStruct(keyPrefix, pathPrefix string, d *schema.ResourceData) PatchOperations {
ops := make([]PatchOperation, 0, 0)
if d.HasChange(keyPrefix+"my_int") {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"MyInt",
Value: d.Get("my_int").(int),
})
}
if d.HasChange(keyPrefix+"my_string") {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"MyString",
Value: d.Get("my_string").(string),
})
}
if d.HasChange(keyPrefix+"my_float") {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"MyFloat",
Value: d.Get("my_float").(float64),
})
}
if d.HasChange(keyPrefix+"my_bool") {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"MyBool",
Value: d.Get("my_bool").(bool),
})
}
return ops
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestPatchOpsFromStruct_simplifiedPrimitives(t *testing.T) {
	type SimpleStruct struct {
		MyInt8    int8
		MyInt32   int32
		MyFloat32 float32
		MyUint    uint
	}
	hg := &HelperGenerator{
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.PatchOpsFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"patchSimpleStruct": `func patchSimpleStruct(keyPrefix, pathPrefix string, d *schema.ResourceData) PatchOperations {
ops := make([]PatchOperation, 0, 0)
if d.HasChange(keyPrefix+"my_int8") {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"MyInt8",
Value: d.Get("my_int8").(int),
})
}
if d.HasChange(keyPrefix+"my_int32") {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"MyInt32",
Value: d.Get("my_int32").(int),
})
}
if d.HasChange(keyPrefix+"my_float32") {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"MyFloat32",
Value: d.Get("my_float32").(float64),
})
}
if d.HasChange(keyPrefix+"my_uint") {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"MyUint",
Value: d.Get("my_uint").(int),
})
}
return ops
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestPatchOpsFromStruct_maps(t *testing.T) {
	type SimpleStruct struct {
		StringMap map[string]string
		IntMap    map[string]int
		FloatMap  map[string]float64
		BoolMap   map[string]bool
	}
	hg := &HelperGenerator{
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.PatchOpsFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"patchSimpleStruct": `func patchSimpleStruct(keyPrefix, pathPrefix string, d *schema.ResourceData) PatchOperations {
ops := make([]PatchOperation, 0, 0)
if d.HasChange(keyPrefix+"string_map") {
oldV, newV := d.GetChange(keyPrefix+"string_map")
diffOps := diffStringMap(pathPrefix+"StringMap/", oldV.(map[string]interface{}), newV.(map[string]interface{}))
ops = append(ops, diffOps...)
}
if d.HasChange(keyPrefix+"int_map") {
oldV, newV := d.GetChange(keyPrefix+"int_map")
diffOps := diffIntMap(pathPrefix+"IntMap/", oldV.(map[string]interface{}), newV.(map[string]interface{}))
ops = append(ops, diffOps...)
}
if d.HasChange(keyPrefix+"float_map") {
oldV, newV := d.GetChange(keyPrefix+"float_map")
diffOps := diffFloatMap(pathPrefix+"FloatMap/", oldV.(map[string]interface{}), newV.(map[string]interface{}))
ops = append(ops, diffOps...)
}
if d.HasChange(keyPrefix+"bool_map") {
oldV, newV := d.GetChange(keyPrefix+"bool_map")
diffOps := diffBoolMap(pathPrefix+"BoolMap/", oldV.(map[string]interface{}), newV.(map[string]interface{}))
ops = append(ops, diffOps...)
}
return ops
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestPatchOpsFromStruct_nestedStruct(t *testing.T) {
	type NestedStruct struct {
		String string
		Int    int
		Bool   bool
	}
	type SimpleStruct struct {
		String       string
		NestedStruct NestedStruct
	}
	hg := &HelperGenerator{
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.PatchOpsFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"patchSimpleStruct": `func patchSimpleStruct(keyPrefix, pathPrefix string, d *schema.ResourceData) PatchOperations {
ops := make([]PatchOperation, 0, 0)
if d.HasChange(keyPrefix+"string") {
ops = append(ops, &ReplaceOperation{
Path:  pathPrefix+"String",
Value: d.Get("string").(string),
})
}
if d.HasChange(keyPrefix+"nested_struct") {
oldV, newV := d.GetChange(keyPrefix+"int_map")
diffOps := diffIntMap(pathPrefix+"IntMap/", oldV.(map[string]interface{}), newV.(map[string]interface{}))
ops = append(ops, diffOps...)
}
return ops
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}
