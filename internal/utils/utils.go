package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ToIntPointer(from *int64) *int {
	if from == nil {
		return nil
	}
	val := int(*from)
	return &val
}

func ToIntListPointer(ctx context.Context, from types.List) ([]int, diag.Diagnostics) {
	var values []int64
	diags := from.ElementsAs(ctx, &values, false)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]int, 0, len(values))
	for _, val := range values {
		out = append(out, int(val))
	}
	return out, diags
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

func MaybeInt64Value(in *int) types.Int64 {
	if in == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*in))
}

func FromInt64Value(in types.Int64) *int {
	if in.IsNull() || in.IsUnknown() {
		return nil
	}
	return ToIntPointer(in.ValueInt64Pointer())
}
