package client

//go:generate go tool oapi-codegen -config config.yaml openapi.json

//go:generate sh -c "go tool openapi-overlay apply tfgen-overlay.yaml openapi.json > tfgen_openapi.yaml"

//go:generate go tool tfplugingen-openapi generate --config tfgen_config.yaml --output resource_spec.json tfgen_openapi.yaml

//go:generate sh -c "go tool json-patch --patch-file resource_spec_patches.json < resource_spec.json > resource_spec.patched.json"
