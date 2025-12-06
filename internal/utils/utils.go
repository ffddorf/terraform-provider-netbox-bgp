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

func FromIntValue(in types.Int64) *int {
	if in.IsNull() || in.IsUnknown() {
		return nil
	}
	v := int(in.ValueInt64())
	return &v
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
