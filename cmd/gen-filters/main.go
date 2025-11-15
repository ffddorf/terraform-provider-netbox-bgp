package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"io"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

var tagRegex = regexp.MustCompile(`([a-z]+):"([^"]+)"\s*`)

type opDef struct {
	jsonName    string
	goFieldName string
	goType      types.Type
}

type fieldDef struct {
	goName    string
	goType    types.Type
	jsonName  string
	operators []opDef
}

func resolveInnerType(t types.Type) (resolved types.Type, isPointer bool, isSlice bool) {
	resolved = t
	pointer, isPointer := resolved.(*types.Pointer)
	if isPointer {
		resolved = pointer.Elem()
	}
	slice, isSlice := resolved.(*types.Slice)
	if isSlice {
		resolved = slice.Elem()
	}
	return
}

func addStringParser(w io.Writer, t types.Type) {
	defer fmt.Fprintln(w)

	typeName := t.String()
	switch typeName {
	case "string":
		fmt.Fprint(w, `v := value.ValueString()`)
		return
	case "time.Time":
		fmt.Fprint(w, `v, err := time.Parse(time.RFC3339, value.ValueString())`)
	case "bool":
		fmt.Fprint(w, `v, err := strconv.ParseBool(value.ValueString())`)
	case "int64":
		fmt.Fprint(w, `v, err := strconv.ParseInt(value.ValueString(), 10, 64)`)
	case "int32":
		fmt.Fprint(w, `v64, err := strconv.ParseInt(value.ValueString(), 10, 32)`)
		fmt.Fprintf(w, `
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as %s value: %%s", value.ValueString()),
			)
		}`, typeName)
		fmt.Fprintln(w)
		fmt.Fprint(w, `v := int32(v64)`)
		return
	case "int":
		fmt.Fprint(w, `v, err := strconv.Atoi(value.ValueString())`)
	case "github.com/google/uuid.UUID", "github.com/oapi-codegen/runtime/types.UUID":
		fmt.Fprint(w, `v, err := uuid.Parse(value.ValueString())`)
	default:
		if t.Underlying() != t {
			addStringParser(w, t.Underlying())
			return
		}
		panic(fmt.Errorf("unhandled type: %s", typeName))
	}
	fmt.Fprintf(w, `
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as %s value: %%s", value.ValueString()),
			)
		}`, typeName)
}

type wrapperOpts struct {
	Logic     string
	ShortName string
	FullName  string

	Fields    []string
	Operators []string
}

var wrapper = template.Must(template.New("wrapper").
	Funcs(template.FuncMap{"StringsJoin": strings.Join}).
	Parse(`// Generated code. DO NOT EDIT!
package provider

var (
	{{ .ShortName }}Fields = []string{"{{ StringsJoin .Fields "\",\"" }}"}
	{{ .ShortName }}Operators = []string{"{{ StringsJoin .Operators "\",\"" }}"}
)

func set{{ .ShortName }}FromFilter(filter Filter, params *client.{{ .FullName }}) diag.Diagnostic {
	name := filter.Name.ValueString()
	op := filter.Operator.ValueString()
	if op == "" {
		op = "eq"
	}
	value := filter.Value
	{{ .Logic }}
	return nil
}
`))

// these have a special meaning and are set manually.
var specialFields = []string{"limit", "offset", "ordering"}

func main() {
	packagePath := os.Args[1]
	paramsName := os.Args[2]

	cfg := &packages.Config{
		Mode: packages.NeedTypes,
	}
	pkgs, err := packages.Load(cfg, packagePath)
	if err != nil {
		panic(err)
	}

	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		fmt.Fprintf(os.Stderr, "Errors encountered while loading package:\n")
		for _, err := range pkg.Errors {
			fmt.Fprintf(os.Stderr, "%+v\n", err)
		}
		os.Exit(1)
	}

	scope := pkg.Types.Scope()
	paramsObj := scope.Lookup(paramsName)
	if paramsObj == nil {
		panic(fmt.Errorf("type not found: %s", paramsName))
	}

	paramsType := paramsObj.(*types.TypeName)          //nolint:forcetypeassert
	paramsNamed := paramsType.Type().(*types.Named)    //nolint:forcetypeassert
	params := paramsNamed.Underlying().(*types.Struct) //nolint:forcetypeassert
	fields := make(map[string]*fieldDef)
	for i := range params.NumFields() {
		f := params.Field(i)
		tag := params.Tag(i)
		tagKVs := tagRegex.FindAllStringSubmatch(tag, -1)
		for _, m := range tagKVs {
			if m[1] == "json" {
				jsonBaseName, _, _ := strings.Cut(m[2], ",")
				jsonName, op, hasSuffix := strings.Cut(jsonBaseName, "__")
				if !hasSuffix {
					op = "eq"
					jsonName = jsonBaseName
				}
				def, ok := fields[jsonName]
				if !ok {
					def = &fieldDef{
						jsonName: jsonName,
					}
					fields[jsonName] = def
				}
				if !hasSuffix {
					// update field def
					def.goName = f.Name()
					def.goType = f.Type()
				}
				def.operators = append(def.operators, opDef{
					jsonName:    op,
					goFieldName: f.Name(),
					goType:      f.Type(),
				})
				break
			}
		}
	}

	var output strings.Builder
	output.WriteString(`switch name {`)
	output.WriteByte('\n')

	fieldList := make([]*fieldDef, 0, len(fields))
	for _, field := range fields {
		fieldList = append(fieldList, field)
	}
	slices.SortFunc(fieldList, func(a, b *fieldDef) int {
		return strings.Compare(a.jsonName, b.jsonName)
	})

	allFields := make([]string, 0, len(fieldList))
	operatorSet := make(map[string]struct{})
	for _, field := range fieldList {
		if slices.Contains(specialFields, field.jsonName) {
			continue
		}

		allFields = append(allFields, field.jsonName)

		fmt.Fprintf(&output, `case "%s":`, field.jsonName)
		output.WriteByte('\n')

		fieldType, _, _ := resolveInnerType(field.goType)
		if fieldType == nil {
			fmt.Printf("type on %s resolved to no type: %#v\n", field.jsonName, field.goType)
			continue
		}
		addStringParser(&output, fieldType)

		output.WriteString(`switch op {`)
		output.WriteByte('\n')

		for _, op := range field.operators {
			operatorSet[op.jsonName] = struct{}{}

			fmt.Fprintf(&output, `case "%s":`, op.jsonName)
			output.WriteByte('\n')

			opType, isPointer, isSlice := resolveInnerType(op.goType)
			if fieldType != opType {
				addStringParser(&output, opType)
			}

			fmt.Fprintf(&output, `params.%s = `, op.goFieldName)
			valName := "v"
			if opType.Underlying() != opType {
				if opTypeNamed, ok := opType.(*types.Named); ok {
					valName = opTypeNamed.Obj().Pkg().Name() + "." + opTypeNamed.Obj().Name() + "(v)"
				}
			}
			if isSlice {
				if isPointer {
					output.WriteString("appendPointerSlice(")
				} else {
					output.WriteString("append(")
				}
				fmt.Fprintf(&output, "params.%s, %s)", op.goFieldName, valName)
			} else {
				if isPointer {
					output.WriteString("ptr(")
					output.WriteString(valName)
					output.WriteByte(')')
				} else {
					output.WriteString(valName)
				}
			}

			output.WriteByte('\n')
		}

		output.WriteString(`		default:
			return unexpectedOperator(op, name)
		}`)
		output.WriteByte('\n')
	}

	output.WriteString(`	default:
		return diag.NewErrorDiagnostic(
			"Unexpected filter name",
			fmt.Sprintf("Did not recognize field name: %s", name),
		)
	}`)

	allOperators := make([]string, 0, len(operatorSet))
	for op := range operatorSet {
		allOperators = append(allOperators, op)
	}
	sort.Strings(allOperators)

	shortName := strings.TrimPrefix(paramsName, "PluginsBgp")
	var wrapped bytes.Buffer
	err = wrapper.Execute(&wrapped, wrapperOpts{
		Logic:     output.String(),
		ShortName: shortName,
		FullName:  paramsName,
		Fields:    allFields,
		Operators: allOperators,
	})
	if err != nil {
		panic(err)
	}

	formatted, err := format.Source(wrapped.Bytes())
	if err != nil {
		panic(err)
	}

	_, _ = io.Copy(os.Stdout, bytes.NewReader(formatted))
}
