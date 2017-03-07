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

		functions := hg.FlattenersFromStruct(s.Obj)

		err = tpl.Execute(f, struct {
			PkgName   string
			Functions map[string]string
		}{
			PkgName:   pkgName,
			Functions: functions,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

var tpl = template.Must(template.New("pod").Parse(`package {{.PkgName}}

{{range $name, $definition := .Functions}}
{{ $definition }}
{{end}}
`))
