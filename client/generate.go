package client

//go:generate go tool oapi-codegen -config config.yaml openapi.json

//go:generate go tool tfplugingen-openapi generate --config tfgen_config.yaml --output resource_spec.json openapi.json
