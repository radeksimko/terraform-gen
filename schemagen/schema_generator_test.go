package schemagen

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

func TestGenerateField_primitive(t *testing.T) {
	type SimpleStruct struct {
		MyInt     int `json:"myInt"`
		MyInt8    int8
		MyInt16   int16
		MyInt32   int32
		MyInt64   int64
		MyFloat32 float32
		MyFloat64 float64
		MyString  string `json:"myString"`
		MyBool    bool
	}
	docsF := func(_struct interface{}, sf *reflect.StructField) string {
		docs := map[string]string{
			"MyInt":    "Description for my integer",
			"MyString": "Description for my string",
		}
		return docs[sf.Name]
	}
	filterF := func(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool) {
		return k, true
	}

	g := &SchemaGenerator{DocsFunc: docsF, FilterFunc: filterF}
	schema := g.FromStruct(&SimpleStruct{})
	expectedSchema := map[string]string{
		"my_int":     "&schema.Schema{\nType: schema.TypeInt,\nDescription: \"Description for my integer\",\n}",
		"my_int8":    "&schema.Schema{\nType: schema.TypeInt,\n}",
		"my_int16":   "&schema.Schema{\nType: schema.TypeInt,\n}",
		"my_int32":   "&schema.Schema{\nType: schema.TypeInt,\n}",
		"my_int64":   "&schema.Schema{\nType: schema.TypeInt,\n}",
		"my_float32": "&schema.Schema{\nType: schema.TypeFloat,\n}",
		"my_float64": "&schema.Schema{\nType: schema.TypeFloat,\n}",
		"my_string":  "&schema.Schema{\nType: schema.TypeString,\nDescription: \"Description for my string\",\n}",
		"my_bool":    "&schema.Schema{\nType: schema.TypeBool,\n}",
	}
	if !reflect.DeepEqual(schema, expectedSchema) {
		t.Fatalf("Expected: %s\n\nGiven: %s\n", expectedSchema, schema)
	}
}

func TestGenerateField_primitivePointers(t *testing.T) {
	type SimpleStruct struct {
		MyInt     *int `json:"myInt"`
		MyInt8    *int8
		MyInt16   *int16
		MyInt32   *int32
		MyInt64   *int64 `json:"myInt64"`
		MyFloat32 *float32
		MyFloat64 *float64
		MyString  *string
		MyBool    *bool `json:"myBool"`
	}
	docsF := func(_struct interface{}, sf *reflect.StructField) string {
		docs := map[string]string{
			"MyInt":   "Description for my integer",
			"MyInt64": "Description for my integer64",
			"MyBool":  "Description for my boolean",
		}
		return docs[sf.Name]
	}
	filterF := func(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool) {
		return k, true
	}

	g := &SchemaGenerator{DocsFunc: docsF, FilterFunc: filterF}
	schema := g.FromStruct(&SimpleStruct{})
	expectedSchema := map[string]string{
		"my_int":     "&schema.Schema{\nType: schema.TypeInt,\nDescription: \"Description for my integer\",\n}",
		"my_int8":    "&schema.Schema{\nType: schema.TypeInt,\n}",
		"my_int16":   "&schema.Schema{\nType: schema.TypeInt,\n}",
		"my_int32":   "&schema.Schema{\nType: schema.TypeInt,\n}",
		"my_int64":   "&schema.Schema{\nType: schema.TypeInt,\nDescription: \"Description for my integer64\",\n}",
		"my_float32": "&schema.Schema{\nType: schema.TypeFloat,\n}",
		"my_float64": "&schema.Schema{\nType: schema.TypeFloat,\n}",
		"my_string":  "&schema.Schema{\nType: schema.TypeString,\n}",
		"my_bool":    "&schema.Schema{\nType: schema.TypeBool,\nDescription: \"Description for my boolean\",\n}",
	}
	if !reflect.DeepEqual(schema, expectedSchema) {
		t.Fatalf("Expected: %s\n\nGiven: %s\n", expectedSchema, schema)
	}
}

func TestGenerateField_primitiveMixed(t *testing.T) {
	type SimpleStruct struct {
		MyInt    *int
		MyInt8   int8
		MyInt16  *int16
		MyInt32  int32
		MyInt64  *int64
		MyString *string
		MyBool   bool
	}
	docsF := func(_struct interface{}, sf *reflect.StructField) string {
		return ""
	}
	filterF := func(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool) {
		return k, true
	}

	g := &SchemaGenerator{DocsFunc: docsF, FilterFunc: filterF}
	schema := g.FromStruct(&SimpleStruct{})
	expectedSchema := map[string]string{
		"my_int":    "&schema.Schema{\nType: schema.TypeInt,\n}",
		"my_int8":   "&schema.Schema{\nType: schema.TypeInt,\n}",
		"my_int16":  "&schema.Schema{\nType: schema.TypeInt,\n}",
		"my_int32":  "&schema.Schema{\nType: schema.TypeInt,\n}",
		"my_int64":  "&schema.Schema{\nType: schema.TypeInt,\n}",
		"my_string": "&schema.Schema{\nType: schema.TypeString,\n}",
		"my_bool":   "&schema.Schema{\nType: schema.TypeBool,\n}",
	}
	if !reflect.DeepEqual(schema, expectedSchema) {
		t.Fatalf("Expected: %s\n\nGiven: %s\n", expectedSchema, schema)
	}
}

func TestGenerateField_sliceOfPrimitives(t *testing.T) {
	type SimpleStruct struct {
		MyInt     []int
		MyInt8    []int8
		MyInt16   []int16
		MyInt32   []int32
		MyInt64   []int64
		MyFloat32 []float32
		MyFloat64 []float64
		MyString  []string
		MyBool    []bool
	}
	docsF := func(_struct interface{}, sf *reflect.StructField) string {
		return ""
	}
	filterF := func(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool) {
		return k, true
	}

	g := &SchemaGenerator{DocsFunc: docsF, FilterFunc: filterF}
	schema := g.FromStruct(&SimpleStruct{})
	expectedSchema := map[string]string{
		"my_int":     "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
		"my_int8":    "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
		"my_int16":   "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
		"my_int32":   "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
		"my_int64":   "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
		"my_float32": "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeFloat,},\n}",
		"my_float64": "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeFloat,},\n}",
		"my_string":  "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeString,},\nSet: schema.HashString,\n}",
		"my_bool":    "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeBool,},\n}",
	}
	if !reflect.DeepEqual(schema, expectedSchema) {
		t.Fatalf("Expected: %#v\n\nGiven: %#v\n", expectedSchema, schema)
	}
}

func TestGenerateField_sliceOfPtrsToPrimitives(t *testing.T) {
	type SimpleStruct struct {
		MyInt     []*int
		MyInt8    []*int8
		MyInt16   []*int16
		MyInt32   []*int32
		MyInt64   []*int64
		MyFloat32 []*float32
		MyFloat64 []*float64
		MyString  []*string
		MyBool    []*bool
	}
	docsF := func(_struct interface{}, sf *reflect.StructField) string {
		return ""
	}
	filterF := func(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool) {
		return k, true
	}

	g := &SchemaGenerator{DocsFunc: docsF, FilterFunc: filterF}
	schema := g.FromStruct(&SimpleStruct{})
	expectedSchema := map[string]string{
		"my_int":     "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
		"my_int8":    "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
		"my_int16":   "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
		"my_int32":   "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
		"my_int64":   "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
		"my_float32": "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeFloat,},\n}",
		"my_float64": "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeFloat,},\n}",
		"my_string":  "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeString,},\nSet: schema.HashString,\n}",
		"my_bool":    "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeBool,},\n}",
	}
	if !reflect.DeepEqual(schema, expectedSchema) {
		t.Fatalf("Expected: %#v\n\nGiven: %#v\n", expectedSchema, schema)
	}
}

func TestGenerateField_struct(t *testing.T) {
	type NestedStruct struct {
		MyInt    int
		MyString string
	}
	type SimpleStruct struct {
		Nested *NestedStruct
		MyInt  []int
	}

	docsF := func(_struct interface{}, sf *reflect.StructField) string {
		return ""
	}
	filterF := func(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool) {
		return k, true
	}

	g := &SchemaGenerator{DocsFunc: docsF, FilterFunc: filterF}
	schema := g.FromStruct(&SimpleStruct{})
	expectedSchema := map[string]string{
		"nested": "&schema.Schema{\nType: schema.TypeList,\nMaxItems: 1,\nElem: &schema.Resource{\nSchema: map[string]*schema.Schema{\n\"my_int\": &schema.Schema{\nType: schema.TypeInt,\n},\n\"my_string\": &schema.Schema{\nType: schema.TypeString,\n},\n},\n},\n}",
		"my_int": "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
	}
	if !reflect.DeepEqual(schema, expectedSchema) {
		t.Fatalf("Expected: %s\n\nGiven: %s\n", expectedSchema, schema)
	}
}

func TestGenerateField_sliceOfStructs(t *testing.T) {
	type NestedStruct struct {
		MyInt    int
		MyString string
	}
	type SimpleStruct struct {
		Nested []*NestedStruct
		MyInt  []int
	}

	docsF := func(_struct interface{}, sf *reflect.StructField) string {
		return ""
	}
	filterF := func(iface interface{}, sf *reflect.StructField, k reflect.Kind, s *schema.Schema) (reflect.Kind, bool) {
		return k, true
	}

	g := &SchemaGenerator{DocsFunc: docsF, FilterFunc: filterF}
	schema := g.FromStruct(&SimpleStruct{})
	expectedSchema := map[string]string{
		"nested": "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Resource{\nSchema: map[string]*schema.Schema{\n\"my_int\": &schema.Schema{\nType: schema.TypeInt,\n},\n\"my_string\": &schema.Schema{\nType: schema.TypeString,\n},\n},\n},\n}",
		"my_int": "&schema.Schema{\nType: schema.TypeSet,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
	}
	if !reflect.DeepEqual(schema, expectedSchema) {
		t.Fatalf("Expected: %s\n\nGiven: %s\n", expectedSchema, schema)
	}
}
