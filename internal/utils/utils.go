package utils

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

func FromStringValue(in types.String) *string {
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

func MaybeListValue[T any](ctx context.Context, elementType attr.Type, path path.Path, in *[]T, diags diag.Diagnostics) types.List {
	if in == nil || len(*in) == 0 {
		return types.ListNull(elementType)
	}

	result, ds := types.ListValueFrom(ctx, elementType, in)
	for _, d := range ds {
		diags.Append(diag.WithPath(path, d))
	}
	return result
}

func TimeString(t time.Time) string {
	return t.Format(time.RFC3339)
}
