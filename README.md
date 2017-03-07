# terraform-gen

Generator(s) of various things for [Terraform](https://github.com/hashicorp/terraform)

## Why

After you spend time implementing many **resources** and **data sources** for Terraform
you might realise there's a few parts of the work that are repetitive and :scream: scream for automation.

These include:

 - schema (more or less enriched copy of SDK struct & its fields)
 - `flatten*` & `expand*` helper functions (SDK struct <-> schema)
 - documentation for all fields

## Unresolved challenges

 - `TypeSet` vs `TypeList`
 - `ValidateFunc`
 - `MinItems`
 - `Set` for complex (non-primitive) fields
 - ... many others

## Caveats

The current version generates code that is unlikely to be accepted/production-ready without manual tweaks.
As such it's **not recommended** to call this as part of `go generate` nor submit PRs to Terraform
with raw generated code.

Also `gofmt` is your friend. :shower:

## Examples

See [`/_examples`](https://github.com/radeksimko/terraform-gen/tree/master/_examples).

## License

See [`LICENSE`](https://github.com/radeksimko/terraform-gen/tree/master/LICENSE).
