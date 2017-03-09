package main

import (
	"log"
	"os"
	"text/template"

	"github.com/radeksimko/terraform-gen/helpergen"

	api "k8s.io/kubernetes/pkg/api/v1"
)

type helperGen struct {
	Obj      interface{}
	Filename string
}

func main() {
	pkgName := "kubernetes"
	schemas := []helperGen{
		{
			Obj:      &api.PersistentVolumeSpec{},
			Filename: "structure_persistent_volume_spec.go",
		},
	}

	for _, s := range schemas {
		log.Printf("Generating %q...\n", s.Filename)
		f, err := os.Create(s.Filename)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}

		hg := &helpergen.HelperGenerator{
			InputVarName:  "in",
			OutputVarName: "att",
		}

		flatteners := hg.FlattenersFromStruct(s.Obj)
		expanders := hg.ExpandersFromStruct(s.Obj)

		err = tpl.Execute(f, struct {
			PkgName    string
			Flatteners map[string]string
			Expanders  map[string]string
		}{
			PkgName:    pkgName,
			Flatteners: flatteners,
			Expanders:  expanders,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

var tpl = template.Must(template.New("pod").Parse(`package {{.PkgName}}

// Flatteners
{{range $name, $definition := .Flatteners}}
{{ $definition }}
{{end}}

// Expanders
{{range $name, $definition := .Expanders}}
{{ $definition }}
{{end}}
`))
