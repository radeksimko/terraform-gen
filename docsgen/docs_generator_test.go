package docsgen

import (
	"bytes"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

func TestGenerateResourceMarkdown_basic(t *testing.T) {
	resource := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Standard object's metadata. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata",
				Required:    true,
			},
			"my_int": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Sample integer.",
				Required:    true,
			},
			"computed_field": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Yada yada yada",
			},
			"my_optional_bool": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Standard boolean.",
				Optional:    true,
			},
		},
	}

	buf := bytes.NewBuffer([]byte{})
	r := &Resource{
		ProviderKey:    "cattle",
		ProviderName:   "Cattle",
		ResourceKey:    "cattle_cow",
		ResourceSlug:   "cattle-cow",
		ResourceSchema: &resource,
	}
	err := r.GenerateResourceMarkdown(buf)
	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	expectedOutput := markdown_basic_output
	if output != expectedOutput {
		t.Fatalf("Output doesn't match.\nExpected: %s\nGiven: %s\n", expectedOutput, output)
	}
}

func TestGenerateResourceMarkdown_nested(t *testing.T) {
	resource := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Standard object's metadata. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nested_string": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of a nested string",
						},
						"nested_int": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Description of a nested integer",
						},
					},
				},
			},
			"my_int": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Sample integer.",
				Required:    true,
			},
			"computed_field": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Yada yada yada",
			},
			"my_optional_bool": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Standard boolean.",
				Optional:    true,
			},
		},
	}

	buf := bytes.NewBuffer([]byte{})
	r := &Resource{
		ProviderKey:    "cattle",
		ProviderName:   "Cattle",
		ResourceKey:    "cattle_cow",
		ResourceSlug:   "cattle-cow",
		ResourceSchema: &resource,
	}
	err := r.GenerateResourceMarkdown(buf)
	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	expectedOutput := markdown_nested_output
	if output != expectedOutput {
		t.Fatalf("Output doesn't match.\nExpected: %s\nGiven: %s\n", expectedOutput, output)
	}
}

func TestGenerateResourceMarkdown_doubleNested(t *testing.T) {
	resource := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": &schema.Schema{
				Type:        schema.TypeList,
				Description: "Standard object's metadata. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nested_string": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of a nested string",
						},
						"nested_int": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Description of a nested integer",
						},
						"nested_list": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Yada yada yada",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"one": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"two": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"nested_set": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Yada yada yada",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"three": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Third nested description",
									},
									"four": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Fourth nested description",
									},
								},
							},
						},
					},
				},
			},
			"my_int": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Sample integer.",
				Required:    true,
			},
			"computed_field": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Yada yada yada",
			},
			"my_optional_bool": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Standard boolean.",
				Optional:    true,
			},
		},
	}

	buf := bytes.NewBuffer([]byte{})
	r := &Resource{
		ProviderKey:    "cattle",
		ProviderName:   "Cattle",
		ResourceKey:    "cattle_cow",
		ResourceSlug:   "cattle-cow",
		ResourceSchema: &resource,
	}
	err := r.GenerateResourceMarkdown(buf)
	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	expectedOutput := markdown_double_nested_output
	if output != expectedOutput {
		t.Fatalf("Output doesn't match.\nExpected: %s\nGiven: %s\n", expectedOutput, output)
	}
}

var markdown_basic_output = `
---
layout: "cattle"
page_title: "Cattle: cattle_cow"
sidebar_current: "docs-cattle-cow"
description: |-
  TODO
---

# cattle_cow

TODO


## Example Usage

` + "```" + `
resource "cattle_cow" "example" {
  // TODO
}
` + "```" + `

## Argument Reference

The following arguments are supported:

* ` + "`metadata`" + ` - (Required) Standard object's metadata. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
* ` + "`my_int`" + ` - (Required) Sample integer.
* ` + "`my_optional_bool`" + ` - (Optional) Standard boolean.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* ` + "`computed_field`" + ` - Yada yada yada

## Import

cattle_cow can be imported using the , e.g.

` + "```" + `
$ terraform import cattle_cow.example ...
` + "```" + `

`

var markdown_nested_output = `
---
layout: "cattle"
page_title: "Cattle: cattle_cow"
sidebar_current: "docs-cattle-cow"
description: |-
  TODO
---

# cattle_cow

TODO


## Example Usage

` + "```" + `
resource "cattle_cow" "example" {
  // TODO
}
` + "```" + `

## Argument Reference

The following arguments are supported:

* ` + "`metadata`" + ` - (Required) Standard object's metadata. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
* ` + "`my_int`" + ` - (Required) Sample integer.
* ` + "`my_optional_bool`" + ` - (Optional) Standard boolean.

## Nested Blocks

### ` + "`metadata`" + `

#### Arguments

* ` + "`nested_int`" + ` - (Optional) Description of a nested integer
* ` + "`nested_string`" + ` - (Optional) Description of a nested string

#### Attributes




## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* ` + "`computed_field`" + ` - Yada yada yada

## Import

cattle_cow can be imported using the , e.g.

` + "```" + `
$ terraform import cattle_cow.example ...
` + "```" + `

`

var markdown_double_nested_output = `
---
layout: "cattle"
page_title: "Cattle: cattle_cow"
sidebar_current: "docs-cattle-cow"
description: |-
  TODO
---

# cattle_cow

TODO


## Example Usage

` + "```" + `
resource "cattle_cow" "example" {
  // TODO
}
` + "```" + `

## Argument Reference

The following arguments are supported:

* ` + "`metadata`" + ` - (Required) Standard object's metadata. More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
* ` + "`my_int`" + ` - (Required) Sample integer.
* ` + "`my_optional_bool`" + ` - (Optional) Standard boolean.

## Nested Blocks

### ` + "`metadata`" + `

#### Arguments

* ` + "`nested_int`" + ` - (Optional) Description of a nested integer
* ` + "`nested_list`" + ` - (Required) Yada yada yada
* ` + "`nested_set`" + ` - (Required) Yada yada yada
* ` + "`nested_string`" + ` - (Optional) Description of a nested string

#### Attributes



### ` + "`nested_list`" + `

#### Arguments

* ` + "`one`" + ` - (Optional) 
* ` + "`two`" + ` - (Optional) 

#### Attributes



### ` + "`nested_set`" + `

#### Arguments

* ` + "`four`" + ` - (Optional) Fourth nested description
* ` + "`three`" + ` - (Optional) Third nested description

#### Attributes




## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* ` + "`computed_field`" + ` - Yada yada yada

## Import

cattle_cow can be imported using the , e.g.

` + "```" + `
$ terraform import cattle_cow.example ...
` + "```" + `

`
