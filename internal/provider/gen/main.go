package main

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	//go:embed resource.go.tmpl
	tmplRaw string
	tmpl    = template.Must(
		template.New("resource").
			Funcs(template.FuncMap{
				"toPublicName": cases.Title(language.AmericanEnglish).String,
			}).
			Parse(tmplRaw),
	)
)

type ResourceDef struct {
	Name    string
	APIName string
}

func main() {
	resources := []ResourceDef{
		{Name: "session", APIName: "Session"},
		{Name: "peergroup", APIName: "PeerGroup"},
		{Name: "prefixlist", APIName: "PrefixList"},
		{Name: "prefixlistrule", APIName: "PrefixListRule"},
		{Name: "aspathlist", APIName: "AspathList"},
		{Name: "aspathlistrule", APIName: "AspathListRule"},
	}

	for _, res := range resources {
		if err := createResource(res); err != nil {
			panic(fmt.Errorf("failed to create %s: %w", res.Name, err))
		}
	}
}

func createResource(res ResourceDef) error {
	f, err := os.Create(res.Name + "_resource.gen.go")
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, res)
}
