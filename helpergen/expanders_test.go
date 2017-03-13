package helpergen

import (
	"reflect"
	"testing"
)

func TestExpanderFromStruct_primitives(t *testing.T) {
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
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.ExpandersFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"expandSimpleStruct": `func expandSimpleStruct(l []interface{}) helpergen.SimpleStruct {
if len(l) == 0 || l[0] == nil {
return helpergen.SimpleStruct{}
}
cfg := l[0].(map[string]interface{})
obj := helpergen.SimpleStruct{
MyInt: cfg["my_int"].(int),
MyInt8: cfg["my_int8"].(int8),
MyInt16: cfg["my_int16"].(int16),
MyInt32: cfg["my_int32"].(int32),
MyInt64: cfg["my_int64"].(int64),
MyFloat32: cfg["my_float32"].(float32),
MyFloat64: cfg["my_float64"].(float64),
MyString: cfg["my_string"].(string),
MyBool: cfg["my_bool"].(bool),
}
return obj
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestExpanderFromPtrToStruct_primitives(t *testing.T) {
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
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.ExpandersFromStruct(&SimpleStruct{})
	expectedOutput := map[string]string{
		"expandSimpleStruct": `func expandSimpleStruct(l []interface{}) *helpergen.SimpleStruct {
if len(l) == 0 || l[0] == nil {
return &helpergen.SimpleStruct{}
}
cfg := l[0].(map[string]interface{})
obj := helpergen.SimpleStruct{
MyInt: cfg["my_int"].(int),
MyInt8: cfg["my_int8"].(int8),
MyInt16: cfg["my_int16"].(int16),
MyInt32: cfg["my_int32"].(int32),
MyInt64: cfg["my_int64"].(int64),
MyFloat32: cfg["my_float32"].(float32),
MyFloat64: cfg["my_float64"].(float64),
MyString: cfg["my_string"].(string),
MyBool: cfg["my_bool"].(bool),
}
return &obj
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestExpanderFromPtrToStruct_ptrsToPrimitives(t *testing.T) {
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
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.ExpandersFromStruct(&SimpleStruct{})
	expectedOutput := map[string]string{
		"expandSimpleStruct": `func expandSimpleStruct(l []interface{}) *helpergen.SimpleStruct {
if len(l) == 0 || l[0] == nil {
return &helpergen.SimpleStruct{}
}
cfg := l[0].(map[string]interface{})
obj := helpergen.SimpleStruct{
MyInt: ptrToInt(cfg["my_int"].(int)),
MyInt8: ptrToInt8(cfg["my_int8"].(int8)),
MyInt16: ptrToInt16(cfg["my_int16"].(int16)),
MyInt32: ptrToInt32(cfg["my_int32"].(int32)),
MyInt64: ptrToInt64(cfg["my_int64"].(int64)),
MyFloat32: ptrToFloat32(cfg["my_float32"].(float32)),
MyFloat64: ptrToFloat64(cfg["my_float64"].(float64)),
MyString: ptrToString(cfg["my_string"].(string)),
MyBool: ptrToBool(cfg["my_bool"].(bool)),
}
return &obj
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestExpanderFromStruct_stringMap(t *testing.T) {
	type SimpleStruct struct {
		MyInt    int
		MyString string
		MyMap    map[string]string
	}
	hg := &HelperGenerator{
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.ExpandersFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"expandSimpleStruct": `func expandSimpleStruct(l []interface{}) helpergen.SimpleStruct {
if len(l) == 0 || l[0] == nil {
return helpergen.SimpleStruct{}
}
cfg := l[0].(map[string]interface{})
obj := helpergen.SimpleStruct{
MyInt: cfg["my_int"].(int),
MyString: cfg["my_string"].(string),
MyMap: expandStringMap(cfg["my_map"].(map[string]interface{})),
}
return obj
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestExpanderFromStruct_primitiveSlices(t *testing.T) {
	type SimpleStruct struct {
		SliceOfInt     []int
		SliceOfInt32   []int32
		SliceOfInt64   []int64
		SliceOfString  []string
		SliceOfFloat64 []float64
		SliceOfBool    []bool
		SimpleInt      int
		SimpleString   string
	}
	hg := &HelperGenerator{
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.ExpandersFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"expandSimpleStruct": `func expandSimpleStruct(l []interface{}) helpergen.SimpleStruct {
if len(l) == 0 || l[0] == nil {
return helpergen.SimpleStruct{}
}
cfg := l[0].(map[string]interface{})
obj := helpergen.SimpleStruct{
SliceOfInt: sliceOfInt(cfg["slice_of_int"].([]interface{})),
SliceOfInt32: sliceOfInt(cfg["slice_of_int32"].([]interface{})),
SliceOfInt64: sliceOfInt(cfg["slice_of_int64"].([]interface{})),
SliceOfString: sliceOfString(cfg["slice_of_string"].([]interface{})),
SliceOfFloat64: sliceOfFloat(cfg["slice_of_float64"].([]interface{})),
SliceOfBool: sliceOfBool(cfg["slice_of_bool"].([]interface{})),
SimpleInt: cfg["simple_int"].(int),
SimpleString: cfg["simple_string"].(string),
}
return obj
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestExpanderFromStruct_primitivePtrSlices(t *testing.T) {
	type SimpleStruct struct {
		SliceOfInt     []*int
		SliceOfInt32   []*int32
		SliceOfInt64   []*int64
		SliceOfString  []*string
		SliceOfFloat64 []*float64
		SliceOfBool    []*bool
		SimpleInt      int
		SimpleString   string
	}
	hg := &HelperGenerator{
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.ExpandersFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"expandSimpleStruct": `func expandSimpleStruct(l []interface{}) helpergen.SimpleStruct {
if len(l) == 0 || l[0] == nil {
return helpergen.SimpleStruct{}
}
cfg := l[0].(map[string]interface{})
obj := helpergen.SimpleStruct{
SliceOfInt: sliceOfPtrInt(cfg["slice_of_int"].([]interface{})),
SliceOfInt32: sliceOfPtrInt(cfg["slice_of_int32"].([]interface{})),
SliceOfInt64: sliceOfPtrInt(cfg["slice_of_int64"].([]interface{})),
SliceOfString: sliceOfPtrString(cfg["slice_of_string"].([]interface{})),
SliceOfFloat64: sliceOfPtrFloat(cfg["slice_of_float64"].([]interface{})),
SliceOfBool: sliceOfPtrBool(cfg["slice_of_bool"].([]interface{})),
SimpleInt: cfg["simple_int"].(int),
SimpleString: cfg["simple_string"].(string),
}
return obj
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestExpanderFromStruct_structSlice(t *testing.T) {
	type NestedStruct struct {
		NestedInt    int
		NestedString string
	}
	type SimpleStruct struct {
		NestedSlice  []NestedStruct
		SimpleInt    int
		SimpleString string
	}
	hg := &HelperGenerator{
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.ExpandersFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"expandSimpleStruct": `func expandSimpleStruct(l []interface{}) helpergen.SimpleStruct {
if len(l) == 0 || l[0] == nil {
return helpergen.SimpleStruct{}
}
cfg := l[0].(map[string]interface{})
obj := helpergen.SimpleStruct{
NestedSlice: expandNestedStruct(cfg["nested_slice"].([]interface{})),
SimpleInt: cfg["simple_int"].(int),
SimpleString: cfg["simple_string"].(string),
}
return obj
}`,
		"expandNestedStruct": `func expandNestedStruct(l []interface{}) []helpergen.NestedStruct {
if len(l) == 0 || l[0] == nil {
return []helpergen.NestedStruct{}
}
obj := make([]helpergen.NestedStruct, len(l), len(l))
for i, n := range l {
cfg := n.(map[string]interface{})
obj[i] = helpergen.NestedStruct{
NestedInt: cfg["nested_int"].(int),
NestedString: cfg["nested_string"].(string),
}
}
return obj
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestExpanderFromStruct_structPtrSlice(t *testing.T) {
	type NestedStruct struct {
		NestedInt    int
		NestedString string
	}
	type SimpleStruct struct {
		NestedSlice  []*NestedStruct
		SimpleInt    int
		SimpleString string
	}
	hg := &HelperGenerator{
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.ExpandersFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"expandSimpleStruct": `func expandSimpleStruct(l []interface{}) helpergen.SimpleStruct {
if len(l) == 0 || l[0] == nil {
return helpergen.SimpleStruct{}
}
cfg := l[0].(map[string]interface{})
obj := helpergen.SimpleStruct{
NestedSlice: expandNestedStruct(cfg["nested_slice"].([]interface{})),
SimpleInt: cfg["simple_int"].(int),
SimpleString: cfg["simple_string"].(string),
}
return obj
}`,
		"expandNestedStruct": `func expandNestedStruct(l []interface{}) []*helpergen.NestedStruct {
if len(l) == 0 || l[0] == nil {
return []*helpergen.NestedStruct{}
}
obj := make([]*helpergen.NestedStruct, len(l), len(l))
for i, n := range l {
cfg := n.(map[string]interface{})
obj[i] = &helpergen.NestedStruct{
NestedInt: cfg["nested_int"].(int),
NestedString: cfg["nested_string"].(string),
}
}
return obj
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}

func TestExpanderFromStruct_nestedStruct(t *testing.T) {
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
		InputVarName:  "cfg",
		OutputVarName: "obj",
	}

	output := hg.ExpandersFromStruct(SimpleStruct{})
	expectedOutput := map[string]string{
		"expandSimpleStruct": `func expandSimpleStruct(l []interface{}) helpergen.SimpleStruct {
if len(l) == 0 || l[0] == nil {
return helpergen.SimpleStruct{}
}
cfg := l[0].(map[string]interface{})
obj := helpergen.SimpleStruct{
MyInt: cfg["my_int"].(int),
MyString: cfg["my_string"].(string),
MyBool: cfg["my_bool"].(bool),
MyNested: expandNestedStruct(cfg["my_nested"].([]interface{})),
}
return obj
}`,
		"expandNestedStruct": `func expandNestedStruct(l []interface{}) helpergen.NestedStruct {
if len(l) == 0 || l[0] == nil {
return helpergen.NestedStruct{}
}
cfg := l[0].(map[string]interface{})
obj := helpergen.NestedStruct{
NestedInt: cfg["nested_int"].(int),
NestedString: cfg["nested_string"].(string),
}
return obj
}`,
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Fatalf("\nExpected: %s\n\nGiven:    %s", expectedOutput, output)
	}
}
