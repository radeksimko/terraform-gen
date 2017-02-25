package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/radeksimko/terraform-gen/schemagen"

	api "k8s.io/kubernetes/pkg/api/v1"
)

type schemaGen struct {
	Obj          interface{}
	Filename     string
	VariableName string
}

func main() {
	pkgName := "kubernetes"
	schemas := []schemaGen{
		{
			Obj:          &api.PodSpec{},
			Filename:     "pod_spec_schema.go",
			VariableName: "podSpecSchema",
		},
	}

	for _, s := range schemas {
		log.Printf("Generating %q...\n", s.Filename)
		f, err := os.Create(s.Filename)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}

		sg := &schemagen.SchemaGenerator{DocsFunc: docsFunc, FilterFunc: filterFunc}
		fields := sg.FromStruct(s.Obj)

		err = podTemplate.Execute(f, struct {
			PkgName      string
			VariableName string
			Fields       map[string]string
		}{
			PkgName:      pkgName,
			VariableName: s.VariableName,
			Fields:       fields,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func docsFunc(iface interface{}, sf *reflect.StructField) string {
	tag := sf.Tag
	jsonTag := strings.Split(tag.Get("json"), ",")
	jsonName := jsonTag[0]

	docs, err := getSwaggerDocs(iface, sf, jsonName)
	if err != nil {
		log.Printf("Docs not found: %s", err)
	}
	return docs
}

func filterFunc(iface interface{}, sf *reflect.StructField, kind reflect.Kind, s *schema.Schema) (reflect.Kind, bool) {
	tag := sf.Tag
	jsonTag := strings.Split(tag.Get("json"), ",")
	jsonName := jsonTag[0]

	if jsonName == "-" {
		return kind, false
	}

	t := reflect.TypeOf(iface)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if sf.Type.String() == "v1.Time" {
		log.Printf("Converting %s field to TypeString - %q\n", sf.Type.String(), sf.Name)
		return reflect.String, true
	}

	if sf.Type.String() == "resource.Quantity" {
		log.Printf("Converting %s field to TypeInt - %q\n", sf.Type.String(), sf.Name)
		return reflect.Int, true
	}

	if sf.Type.String() == "intstr.IntOrString" {
		log.Printf("Converting %s field to TypeInt - %q\n", sf.Type.String(), sf.Name)
		return reflect.Int, true
	}

	// Pod
	if (t.String() == "v1.Pod" && sf.Name == "Status") ||
		(t.String() == "v1.Pod" && sf.Name == "PodSpec") ||
		(t.String() == "v1.Volume" && sf.Name == "VolumeSource") {
		log.Printf("Ignoring %q -> %q (will be implemented as data source)", t.String(), sf.Name)
		return kind, false
	}
	// Service
	if t.String() == "v1.Service" && sf.Name == "Status" {
		log.Printf("Ignoring %q -> %q (will be implemented as data source)", t.String(), sf.Name)
		return kind, false
	}
	// ReplicationController
	if (t.String() == "v1.PodTemplateSpec" && sf.Name == "Spec") ||
		(t.String() == "v1.ReplicationController" && sf.Name == "Status") {
		log.Printf("Ignoring %q -> %q (will be implemented as data source)", t.String(), sf.Name)
		return kind, false
	}

	docs, err := getSwaggerDocs(iface, sf, jsonName)
	if err != nil {
		return kind, true
	}

	if strings.Contains(docs, "Deprecated:") {
		log.Printf("Ignoring %q (deprecated)\n", sf.Name)
		return kind, false
	}
	if strings.Contains(docs, "NOT YET IMPLEMENTED.") {
		log.Printf("Ignoring %q (not implemented)\n", sf.Name)
		return kind, false
	}

	if strings.Contains(docs, "Read-only.") {
		s.Computed = true
	} else {
		if strings.Contains(docs, "Required.") || strings.Contains(docs, "Required:") {
			s.Required = true
		} else {
			s.Optional = true
		}
	}

	if strings.Contains(docs, "Cannot be updated.") {
		s.ForceNew = true
	}

	return kind, true
}

func getSwaggerDocs(iface interface{}, sf *reflect.StructField, jsonName string) (string, error) {
	structType := reflect.TypeOf(iface)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	val := reflect.New(structType).Elem()
	method := val.MethodByName("SwaggerDoc")
	if method.IsValid() {
		out := method.Call([]reflect.Value{})
		m := out[0]
		docs := m.MapIndex(reflect.ValueOf(jsonName))
		if !docs.IsValid() {
			docs = m.MapIndex(reflect.ValueOf(""))
		}
		return docs.String(), nil
	}
	return "", fmt.Errorf("Docs not found for %s -> %s (%s)", structType.Name(), sf.Name, sf.Type.String())
}

var podTemplate = template.Must(template.New("pod").Parse(`package {{.PkgName}}

import (
	"github.com/hashicorp/terraform/helper/schema"
)

var {{.VariableName}} = map[string]*schema.Schema{
{{range $name, $schema := .Fields}}
	"{{ $name }}": {{ $schema }},{{end}}
}
`))
