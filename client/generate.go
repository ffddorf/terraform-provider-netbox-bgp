package client

//go:generate go tool oapi-codegen -config config.yaml openapi.json

//go:generate go run github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi generate --config tfgen_config.yaml --output resource_spec.json openapi.json
