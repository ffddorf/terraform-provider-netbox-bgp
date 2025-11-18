package provider

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	tagPath = path.Root("tags")
)

func AddTags(orig schema.Schema) {
	orig.Attributes["tags"] = schema.ListAttribute{
		ElementType: types.StringType,
		Optional:    true,
	}
}

func (pc *ProviderClient) TagsForAPI(
	ctx context.Context,
	plan tfsdk.Plan,
	diags diag.Diagnostics,
) *[]client.NestedTagRequest {

	var tagValues []string
	d := plan.GetAttribute(ctx, tagPath, &tagValues)
	diags.Append(d...)
	if d.HasError() {
		return nil
	}

	var unknownTags []string
	var result []client.NestedTagRequest
	pc.tagCacheMu.Lock()
	for _, t := range tagValues {
		if val, ok := pc.tagCache[t]; ok {
			result = append(result, val)
		} else {
			unknownTags = append(unknownTags, t)
		}
	}
	pc.tagCacheMu.Unlock()

	if len(unknownTags) > 0 {
		resp, err := pc.ExtrasTagsListWithResponse(ctx, &client.ExtrasTagsListParams{
			Slug: &unknownTags,
		})
		if err != nil || resp.JSON200 == nil {
			diags.Append(diag.WithPath(tagPath, APIErrorDiagnostic(
				"Tag Retrieval Failed",
				"Could not retrieve tag info from the API",
				err, resp.StatusCode(), resp.Body,
			)))
			return nil
		}

		pc.tagCacheMu.Lock()
		for _, apiTag := range resp.JSON200.Results {
			tag := client.NestedTagRequest{
				Name: apiTag.Name,
				Slug: apiTag.Slug,
			}
			pc.tagCache[apiTag.Slug] = tag
			result = append(result, tag)
		}
		pc.tagCacheMu.Unlock()
	}

	return &result
}

func (pc *ProviderClient) TagsFromAPI(ctx context.Context, state tfsdk.State, tags *[]client.NestedTag) diag.Diagnostics {
	if tags == nil || len(*tags) == 0 {
		return state.SetAttribute(ctx, tagPath, types.ListNull(types.StringType))
	}

	var tagValues []string
	pc.tagCacheMu.Lock()
	for _, apiTag := range *tags {
		pc.tagCache[apiTag.Slug] = client.NestedTagRequest{
			Name: apiTag.Name,
			Slug: apiTag.Slug,
		}
		tagValues = append(tagValues, apiTag.Slug)
	}
	pc.tagCacheMu.Unlock()

	return state.SetAttribute(ctx, tagPath, tagValues)
}
