package helpergen

import (
	"reflect"
	"testing"
)

func TestFlattenersFromStruct_primitives(t *testing.T) {
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

	output := hg.FlattenersFromStruct(SimpleStruct{})
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
	ptrOutput := hg.FlattenersFromStruct(&SimpleStruct{})
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

	output := hg.FlattenersFromStruct([]SimpleStruct{})
	expectedOutput := map[string]string{
		"flattenSimpleStruct": `func flattenSimpleStruct(in []helpergen.SimpleStruct) []map[string]interface{} {
att := make([]map[string]interface{}, len(in), len(in))
for i, n := range in {
m := make(map[string]interface{})
m["my_int"] = n.MyInt
m["my_int8"] = n.MyInt8
m["my_int16"] = n.MyInt16
m["my_int32"] = n.MyInt32
m["my_int64"] = n.MyInt64
m["my_float32"] = n.MyFloat32
m["my_float64"] = n.MyFloat64
m["my_string"] = n.MyString
m["my_bool"] = n.MyBool
att[i] = m
}
return att
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestFlattenersFromStruct_ptrsToPrimitives(t *testing.T) {
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

	output := hg.FlattenersFromStruct(SimpleStruct{})
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

func TestFlattenersFromStruct_nestedSingleLevel(t *testing.T) {
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

	output := hg.FlattenersFromStruct(SimpleStruct{})
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

func TestFlattenersFromStruct_ptrNestedSingleLevel(t *testing.T) {
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

	output := hg.FlattenersFromStruct(SimpleStruct{})
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

func TestFlattenersFromStruct_primitiveSlice(t *testing.T) {
	type SimpleStruct struct {
		SliceOfInt     []int
		SliceOfString  []string
		SliceOfBool    []bool
		SliceOfFloat64 []float64
	}
	hg := &HelperGenerator{
		InputVarName:  "in",
		OutputVarName: "att",
	}

	output := hg.FlattenersFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"flattenSimpleStruct": `func flattenSimpleStruct(in helpergen.SimpleStruct) map[string]interface{} {
att := make(map[string]interface{})
att["slice_of_int"] = in.SliceOfInt
att["slice_of_string"] = in.SliceOfString
att["slice_of_bool"] = in.SliceOfBool
att["slice_of_float64"] = in.SliceOfFloat64
return att
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestFlattenersFromStruct_primitivePtrSlice(t *testing.T) {
	type SimpleStruct struct {
		SliceOfInt     []*int
		SliceOfString  []*string
		SliceOfBool    []*bool
		SliceOfFloat64 []*float64
	}
	hg := &HelperGenerator{
		InputVarName:  "in",
		OutputVarName: "att",
	}

	output := hg.FlattenersFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"flattenSimpleStruct": `func flattenSimpleStruct(in helpergen.SimpleStruct) map[string]interface{} {
att := make(map[string]interface{})
att["slice_of_int"] = flattenIntSlice(in.SliceOfInt)
att["slice_of_string"] = flattenStringSlice(in.SliceOfString)
att["slice_of_bool"] = flattenBoolSlice(in.SliceOfBool)
att["slice_of_float64"] = flattenFloatSlice(in.SliceOfFloat64)
return att
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestFlattenersFromStruct_sliceOfStructs(t *testing.T) {
	type NestedStruct struct {
		SimpleString string
		SimpleBool   bool
		SimpleFloat  float64
	}
	type SimpleStruct struct {
		SimpleInt      int
		SliceOfStructs []NestedStruct
	}
	hg := &HelperGenerator{
		InputVarName:  "in",
		OutputVarName: "att",
	}

	output := hg.FlattenersFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"flattenSimpleStruct": `func flattenSimpleStruct(in helpergen.SimpleStruct) map[string]interface{} {
att := make(map[string]interface{})
att["simple_int"] = in.SimpleInt
att["slice_of_structs"] = flattenNestedStruct(in.SliceOfStructs)
return att
}`,
		"flattenNestedStruct": `func flattenNestedStruct(in []helpergen.NestedStruct) []map[string]interface{} {
att := make([]map[string]interface{}, len(in), len(in))
for i, n := range in {
m := make(map[string]interface{})
m["simple_string"] = n.SimpleString
m["simple_bool"] = n.SimpleBool
m["simple_float"] = n.SimpleFloat
att[i] = m
}
return att
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}
