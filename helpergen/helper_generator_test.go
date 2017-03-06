package helpergen

import (
	"reflect"
	"testing"
)

func TestFromStruct_primitives(t *testing.T) {
	type SimpleStruct struct {
		MyInt     int
		MyInt8    int8
		MyInt16   int16
		MyInt32   int32
		MyInt64   int64
		MyFloat32 float32
		MyFloat64 float64
		MyString  string
		MyBool    bool
	}
	hg := &HelperGenerator{
		InputVarName:  "in",
		OutputVarName: "att",
	}

	output := hg.FromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"flattenSimpleStruct": `func flattenSimpleStruct(in helpergen.SimpleStruct) map[string]interface{} {
att := make(map[string]interface{})
att["my_int"] = in.MyInt
att["my_int8"] = in.MyInt8
att["my_int16"] = in.MyInt16
att["my_int32"] = in.MyInt32
att["my_int64"] = in.MyInt64
att["my_float32"] = in.MyFloat32
att["my_float64"] = in.MyFloat64
att["my_string"] = in.MyString
att["my_bool"] = in.MyBool
return att
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}

	// Pointer
	ptrOutput := hg.FromStruct(&SimpleStruct{})
	expectedPtrOutput := map[string]string{
		"flattenSimpleStruct": `func flattenSimpleStruct(in *helpergen.SimpleStruct) map[string]interface{} {
att := make(map[string]interface{})
att["my_int"] = in.MyInt
att["my_int8"] = in.MyInt8
att["my_int16"] = in.MyInt16
att["my_int32"] = in.MyInt32
att["my_int64"] = in.MyInt64
att["my_float32"] = in.MyFloat32
att["my_float64"] = in.MyFloat64
att["my_string"] = in.MyString
att["my_bool"] = in.MyBool
return att
}`,
	}
	if !reflect.DeepEqual(ptrOutput, expectedPtrOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedPtrOutput, ptrOutput)
	}
}

func TestFromSliceOfStructs_primitives(t *testing.T) {
	type SimpleStruct struct {
		MyInt     int
		MyInt8    int8
		MyInt16   int16
		MyInt32   int32
		MyInt64   int64
		MyFloat32 float32
		MyFloat64 float64
		MyString  string
		MyBool    bool
	}
	hg := &HelperGenerator{
		InputVarName:  "in",
		OutputVarName: "att",
	}

	output := hg.FromStruct([]SimpleStruct{})
	expectedOutput := map[string]string{
		"flattenSimpleStruct": `func flattenSimpleStruct(in []helpergen.SimpleStruct) []map[string]interface{} {
att := make(map[string]interface{})
att["my_int"] = in.MyInt
att["my_int8"] = in.MyInt8
att["my_int16"] = in.MyInt16
att["my_int32"] = in.MyInt32
att["my_int64"] = in.MyInt64
att["my_float32"] = in.MyFloat32
att["my_float64"] = in.MyFloat64
att["my_string"] = in.MyString
att["my_bool"] = in.MyBool
return []map[string]interface{}{att}
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestFromStruct_ptrsToPrimitives(t *testing.T) {
	type SimpleStruct struct {
		MyInt     *int
		MyInt8    *int8
		MyInt16   *int16
		MyInt32   *int32
		MyInt64   *int64
		MyFloat32 *float32
		MyFloat64 *float64
		MyString  *string
		MyBool    *bool
	}
	hg := &HelperGenerator{
		InputVarName:  "in",
		OutputVarName: "att",
	}

	output := hg.FromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"flattenSimpleStruct": `func flattenSimpleStruct(in helpergen.SimpleStruct) map[string]interface{} {
att := make(map[string]interface{})
att["my_int"] = *in.MyInt
att["my_int8"] = *in.MyInt8
att["my_int16"] = *in.MyInt16
att["my_int32"] = *in.MyInt32
att["my_int64"] = *in.MyInt64
att["my_float32"] = *in.MyFloat32
att["my_float64"] = *in.MyFloat64
att["my_string"] = *in.MyString
att["my_bool"] = *in.MyBool
return att
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestFromStruct_nestedSingleLevel(t *testing.T) {
	type NestedStruct struct {
		NestedInt    int
		NestedString string
	}
	type SimpleStruct struct {
		MyInt    int
		MyString string
		MyBool   bool
		MyNested NestedStruct
	}
	hg := &HelperGenerator{
		InputVarName:  "in",
		OutputVarName: "att",
	}

	output := hg.FromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"flattenSimpleStruct": `func flattenSimpleStruct(in helpergen.SimpleStruct) map[string]interface{} {
att := make(map[string]interface{})
att["my_int"] = in.MyInt
att["my_string"] = in.MyString
att["my_bool"] = in.MyBool
att["my_nested"] = flattenNestedStruct(in.MyNested)
return att
}`,
		"flattenNestedStruct": `func flattenNestedStruct(in helpergen.NestedStruct) map[string]interface{} {
att := make(map[string]interface{})
att["nested_int"] = in.NestedInt
att["nested_string"] = in.NestedString
return att
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestFromStruct_ptrNestedSingleLevel(t *testing.T) {
	type NestedStruct struct {
		NestedInt    int
		NestedString string
	}
	type SimpleStruct struct {
		MyInt    int
		MyString string
		MyBool   bool
		MyNested *NestedStruct
	}
	hg := &HelperGenerator{
		InputVarName:  "in",
		OutputVarName: "att",
	}

	output := hg.FromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"flattenSimpleStruct": `func flattenSimpleStruct(in helpergen.SimpleStruct) map[string]interface{} {
att := make(map[string]interface{})
att["my_int"] = in.MyInt
att["my_string"] = in.MyString
att["my_bool"] = in.MyBool
att["my_nested"] = flattenNestedStruct(in.MyNested)
return att
}`,
		"flattenNestedStruct": `func flattenNestedStruct(in *helpergen.NestedStruct) map[string]interface{} {
att := make(map[string]interface{})
att["nested_int"] = in.NestedInt
att["nested_string"] = in.NestedString
return att
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}
