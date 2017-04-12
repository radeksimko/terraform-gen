package main

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform/builtin/providers/kubernetes"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/radeksimko/terraform-gen/docsgen"
)

func main() {
	buf := bytes.NewBuffer([]byte{})
	p := kubernetes.Provider().(*schema.Provider)
	r := &docsgen.Resource{
		ProviderKey:    "kubernetes",
		ProviderName:   "Kubernetes",
		ResourceKey:    "kubernetes_config_map",
		ResourceSlug:   "kubernetes-config-map",
		ResourceSchema: p.ResourcesMap["kubernetes_config_map"],
	}
	r.GenerateResourceMarkdown(buf)
	fmt.Print(buf.String())
}
