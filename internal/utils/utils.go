package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type StringValuable interface {
	attr.Value
	ValueString() string
	ValueStringPointer() *string
}

func ToIntPointer(from *int64) *int {
	if from == nil {
		return nil
	}
	val := int(*from)
	return &val
}

func ToIntListPointer(ctx context.Context, from types.List, path path.Path, diags diag.Diagnostics) *[]int {
	if from.IsNull() || from.IsUnknown() {
		return &[]int{}
	}

	var values []int64
	errs := from.ElementsAs(ctx, &values, false)
	if errs.HasError() {
		for _, d := range errs {
			diags.Append(diag.WithPath(path, d))
		}
		return nil
	}

	out := make([]int, 0, len(values))
	for _, val := range values {
		out = append(out, int(val))
	}
	return &out
}

func MaybeStringValue(in *string) types.String {
	if in == nil {
		return types.StringNull()
	}
	if *in == "" {
		return types.StringNull()
	}
	return types.StringPointerValue(in)
}

func MaybeStringifiedValue[T any](in *T, convert func(T) string) types.String {
	if in == nil {
		return types.StringNull()
	}

	return types.StringValue(convert(*in))
}

func MaybeConvertedValue[S any, T attr.Value](ctx context.Context, in *S) T {
	var zero T
	t := zero.Type(ctx)
	v, err := t.ValueFromTerraform(ctx, tftypes.NewValue(t.TerraformType(ctx), in))
	if err != nil {
		return zero
	}
	return v.(T) // nolint:forcetypeassert
}

func SafeBytesToString[T ~[]byte](in *T) *string {
	if in == nil {
		return nil
	}
	v := string(*in)
	return &v
}

func MaybeInt64Value(in *int) types.Int64 {
	if in == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*in))
}

func MaybeInt64ValueSubfield[T any](obj *T, access func(T) *int) types.Int64 {
	if obj == nil {
		return types.Int64Null()
	}
	val := access(*obj)
	return MaybeInt64Value(val)
}

func FromInt64Value(in types.Int64) *int {
	if in.IsNull() || in.IsUnknown() {
		return nil
	}
	return ToIntPointer(in.ValueInt64Pointer())
}

func FromStringValue(in StringValuable) *string {
	if in.IsNull() || in.IsUnknown() {
		return nil
	}
	return in.ValueStringPointer()
}

func FromBoolValue(in types.Bool) *bool {
	if in.IsNull() || in.IsUnknown() {
		return nil
	}
	return in.ValueBoolPointer()
}

func FromIntValue(in types.Int64) *int {
	if in.IsNull() || in.IsUnknown() {
		return nil
	}
	v := int(in.ValueInt64())
	return &v
}

func FromListValue[T any](ctx context.Context, p path.Path, in types.List, diags diag.Diagnostics) *[]T {
	if in.IsNull() || in.IsUnknown() {
		return nil
	}

	var out []T
	convertDiags := in.ElementsAs(ctx, &out, false)
	for _, d := range convertDiags {
		diags.Append(diag.WithPath(p, d))
	}
	return &out
}

func MaybeListValueAccessor[M, T any](
	ctx context.Context,
	elementType attr.Type,
	path path.Path,
	in *[]M,
	acc func(M) T,
	diags diag.Diagnostics,
) types.List {
	if in == nil || len(*in) == 0 {
		return types.ListNull(elementType)
	}

	l := make([]T, 0, len(*in))
	for _, v := range *in {
		l = append(l, acc(v))
	}

	v, ds := types.ListValueFrom(ctx, elementType, l)
	for _, d := range ds {
		diags.Append(diag.WithPath(path, d))
	}
	return v
}

func TimeString(t time.Time) string {
	return t.Format(time.RFC3339)
}

func TimeFromStringValue(v types.String) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, v.ValueString())
	if err != nil {
		return t, fmt.Errorf("failed to parse time as RFC3339: %w", err)
	}
	return t, nil
}

func MaybeRawJSON(in StringValuable, path path.Path, diags diag.Diagnostics) *json.RawMessage {
	if in.IsUnknown() || in.IsNull() {
		return nil
	}

	val := json.RawMessage(in.ValueString())
	// validate by parsing
	var tester any
	if err := json.Unmarshal(val, &tester); err != nil {
		diags.AddAttributeError(path, "Format error", fmt.Sprintf("Unable to parse input as JSON: %v", err))
		return nil
	}

	return &val
}

func FromStringParsed[T any](ctx context.Context, p path.Path, in types.String, parser func(string) (T, error), diags diag.Diagnostics) *T {
	if in.IsUnknown() || in.IsNull() {
		return nil
	}

	parsed, err := parser(in.ValueString())
	if err != nil {
		diags.AddAttributeError(p, "Type Error", fmt.Sprintf("failed to parse string: %v", err))
		return nil
	}
	return &parsed
}

func FromListParsed[T any, Inner attr.Value](ctx context.Context, p path.Path, in types.List, parser func(Inner) (T, error), diags diag.Diagnostics) *[]T {
	var inner Inner
	if !in.ElementType(ctx).Equal(inner.Type(ctx)) {
		diags.AddError(
			"Type Error",
			fmt.Sprintf(
				"Wrong type received in list, expected %q, got %q",
				inner.Type(ctx).String(),
				in.ElementType(ctx).String(),
			),
		)
		return nil
	}

	if in.IsUnknown() || in.IsNull() {
		return nil
	}

	elems := in.Elements()
	out := make([]T, 0, len(elems))
	for i, val := range elems {
		asInner, _ := val.(Inner)
		conv, err := parser(asInner)
		if err != nil {
			diags.AddAttributeError(
				path.Empty().AtListIndex(i),
				"Type Error",
				fmt.Sprintf("failed to parse string value: %v", err),
			)
			continue
		}
		out = append(out, conv)
	}
	return &out
}
