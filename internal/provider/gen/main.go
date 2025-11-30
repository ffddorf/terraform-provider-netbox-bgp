package main

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"github.com/iancoleman/strcase"
)

var (
	//go:embed resource.go.tmpl
	tmplRaw string
	tmpl    = template.Must(
		template.New("resource").
			Funcs(template.FuncMap{
				"toCamelCase": strcase.ToCamel,
			}).
			Parse(tmplRaw),
	)
)

func main() {
	resources := []string{
		"session",
		"peer_group",
		"prefix_list",
		"prefix_list_rule",
		"aspath_list",
		"aspath_list_rule",
		"routing_policy",
		"routing_policy_rule",
	}

	for _, res := range resources {
		if err := createResource(res); err != nil {
			panic(fmt.Errorf("failed to create %s: %w", res, err))
		}
	}
}

func createResource(res string) error {
	f, err := os.Create(res + "_resource.gen.go")
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, struct{ Name string }{Name: res})
}
