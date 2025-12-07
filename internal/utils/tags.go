package utils

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	TagFieldName = "tags"
	TagSchema    = schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
	}
)

func TagsForAPIModel(ctx context.Context, l types.List, diags diag.Diagnostics) *[]client.NestedTagRequest {
	if l.IsUnknown() || l.IsNull() {
		return nil
	}

	var tags []string
	if diags := l.ElementsAs(ctx, &tags, true); diags.HasError() {
		diags.Append(diags...)
		return nil
	}

	reqs := make([]client.NestedTagRequest, 0, len(tags))
	for _, t := range tags {
		reqs = append(reqs, client.NestedTagRequest{
			Name: t,
			Slug: slugify(t),
		})
	}

	return &reqs
}

func TagsFromAPI(ctx context.Context, tags *[]client.NestedTag, diags diag.Diagnostics) types.List {
	if tags == nil {
		return types.ListNull(types.StringType)
	}
	if len(*tags) == 0 {
		return types.ListValueMust(types.StringType, []attr.Value{})
	}

	tagNames := make([]string, 0, len(*tags))
	for _, t := range *tags {
		tagNames = append(tagNames, t.Name)
	}

	l, d := types.ListValueFrom(ctx, types.StringType, tagNames)
	diags.Append(d...)
	return l
}
