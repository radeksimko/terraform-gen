package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/radeksimko/terraform-gen/docsgen"
)

func main() {
	providerKey := "kubernetes"
	providerName := "Kubernetes"

	buf := bytes.NewBuffer([]byte{})
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Standard namespace's metadata. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"annotations": &schema.Schema{
							Type:        schema.TypeMap,
							Description: "An unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. More info: http://kubernetes.io/docs/user-guide/annotations",
							Optional:    true,
						},
						"generate_name": &schema.Schema{
							Type:          schema.TypeString,
							Description:   "Prefix, used by the server, to generate a unique name ONLY IF the `name` field has not been provided. This value will also be combined with a unique suffix. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#idempotency",
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"metadata.name"},
						},
						"generation": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "A sequence number representing a specific generation of the desired state.",
							Computed:    true,
						},
						"labels": &schema.Schema{
							Type:        schema.TypeMap,
							Description: "Map of string keys and values that can be used to organize and categorize (scope and select) namespaces. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels",
							Optional:    true,
						},
						"name": &schema.Schema{
							Type:          schema.TypeString,
							Description:   "Name of the namespace, must be unique. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
							Optional:      true,
							ForceNew:      true,
							Computed:      true,
							ConflictsWith: []string{"metadata.generate_name"},
						},
						"resource_version": &schema.Schema{
							Type:        schema.TypeString,
							Description: "An opaque value that represents the internal version of this namespace that can be used by clients to determine when namespaces have changed. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#concurrency-control-and-consistency",
							Computed:    true,
						},
						"self_link": &schema.Schema{
							Type:        schema.TypeString,
							Description: "A URL representing this namespace.",
							Computed:    true,
						},
						"uid": &schema.Schema{
							Type:        schema.TypeString,
							Description: "The unique in time and space value for this namespace. More info: http://kubernetes.io/docs/user-guide/identifiers#uids",
							Computed:    true,
						},
					},
				},
			},
		},
	}

	err := docsgen.GenerateResourceMarkdown(providerKey, providerName, "kubernetes_namespace", "kubernetes-resource-namespace", r, buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(buf.String())
}
