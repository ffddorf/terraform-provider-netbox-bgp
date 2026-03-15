package datasources

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (p *PrefixListsModel) FillFromAPIModel(ctx context.Context, resp *client.PaginatedPrefixListList, diags diag.Diagnostics) {
	p.Next = types.StringPointerValue(resp.Next)
	p.Previous = types.StringPointerValue(resp.Previous)

	results := make([]*ResultsValue, 0, len(resp.Results))
	for _, res := range resp.Results {
		val := &ResultsValue{
			Comments:    types.StringPointerValue(res.Comments),
			Description: types.StringPointerValue(res.Description),
			Display:     types.StringPointerValue(res.Display),
			Family:      types.StringValue(string(res.Family)),
			Id:          types.Int64Value(int64(*res.Id)),
			Name:        types.StringValue(res.Name),
			Tags:        utils.TagsFromAPI(ctx, res.Tags, diags),
			Url:         types.StringPointerValue(res.Url),
			state:       attr.ValueStateKnown,
		}
		results = append(results, val)
	}
	var listDiags diag.Diagnostics
	p.Results, listDiags = types.ListValueFrom(ctx, NewResultsValueUnknown().Type(ctx), results)
	for _, d := range listDiags {
		diags.Append(diag.WithPath(path.Root("results"), d))
	}
}

func (p *PrefixListsModel) ToAPIModel(ctx context.Context, diags diag.Diagnostics) *client.PluginsBgpPrefixListListParams {
	params := &client.PluginsBgpPrefixListListParams{
		Created:      utils.FromListParsed(ctx, path.Root("created"), p.Created, utils.TimeFromStringValue, diags),
		CreatedEmpty: utils.FromListParsed(ctx, path.Root("created__empty"), p.Created_Empty, utils.TimeFromStringValue, diags),
		CreatedGt:    utils.FromListParsed(ctx, path.Root("created__gt"), p.Created_Gt, utils.TimeFromStringValue, diags),
		CreatedGte:   utils.FromListParsed(ctx, path.Root("created__gte"), p.Created_Gte, utils.TimeFromStringValue, diags),
		CreatedLt:    utils.FromListParsed(ctx, path.Root("created__lt"), p.Created_Lt, utils.TimeFromStringValue, diags),
		CreatedLte:   utils.FromListParsed(ctx, path.Root("created__lte"), p.Created_Lte, utils.TimeFromStringValue, diags),
		CreatedN:     utils.FromListParsed(ctx, path.Root("created__n"), p.Created_N, utils.TimeFromStringValue, diags),

		Description:       utils.FromListValue[string](ctx, path.Root("description"), p.Description, diags),
		DescriptionIc:     utils.FromListValue[string](ctx, path.Root("description__ic"), p.Description_Ic, diags),
		DescriptionIe:     utils.FromListValue[string](ctx, path.Root("description__ie"), p.Description_Ie, diags),
		DescriptionIew:    utils.FromListValue[string](ctx, path.Root("description__iew"), p.Description_Iew, diags),
		DescriptionIregex: utils.FromListValue[string](ctx, path.Root("description__iregex"), p.Description_Iregex, diags),
		DescriptionIsw:    utils.FromListValue[string](ctx, path.Root("description__isw"), p.Description_Isw, diags),
		DescriptionN:      utils.FromListValue[string](ctx, path.Root("description__n"), p.Description_N, diags),
		DescriptionNic:    utils.FromListValue[string](ctx, path.Root("description__nic"), p.Description_Nic, diags),
		DescriptionNie:    utils.FromListValue[string](ctx, path.Root("description__nie"), p.Description_Nie, diags),
		DescriptionNiew:   utils.FromListValue[string](ctx, path.Root("description__niew"), p.Description_Niew, diags),
		DescriptionNisw:   utils.FromListValue[string](ctx, path.Root("description__nisw"), p.Description_Nisw, diags),
		DescriptionRegex:  utils.FromListValue[string](ctx, path.Root("description__regex"), p.Description_Regex, diags),
		DescriptionEmpty:  utils.FromBoolValue(p.Description_Empty),

		Family: (*client.PluginsBgpPrefixListListParamsFamily)(utils.FromStringValue(p.Family)),

		Id:      utils.FromListValue[int32](ctx, path.Root("id"), p.Id, diags),
		IdGt:    utils.FromListValue[int32](ctx, path.Root("id__gt"), p.Id_Gt, diags),
		IdGte:   utils.FromListValue[int32](ctx, path.Root("id__gte"), p.Id_Gte, diags),
		IdLt:    utils.FromListValue[int32](ctx, path.Root("id__lt"), p.Id_Lt, diags),
		IdLte:   utils.FromListValue[int32](ctx, path.Root("id__lte"), p.Id_Lte, diags),
		IdN:     utils.FromListValue[int32](ctx, path.Root("id__n"), p.Id_N, diags),
		IdEmpty: utils.FromBoolValue(p.Id_Empty),

		LastUpdated:      utils.FromListParsed(ctx, path.Root("last_updated"), p.LastUpdated, utils.TimeFromStringValue, diags),
		LastUpdatedEmpty: utils.FromListParsed(ctx, path.Root("last_updated__empty"), p.LastUpdated_Empty, utils.TimeFromStringValue, diags),
		LastUpdatedGt:    utils.FromListParsed(ctx, path.Root("last_updated__gt"), p.LastUpdated_Gt, utils.TimeFromStringValue, diags),
		LastUpdatedGte:   utils.FromListParsed(ctx, path.Root("last_updated__gte"), p.LastUpdated_Gte, utils.TimeFromStringValue, diags),
		LastUpdatedLt:    utils.FromListParsed(ctx, path.Root("last_updated__lt"), p.LastUpdated_Lt, utils.TimeFromStringValue, diags),
		LastUpdatedLte:   utils.FromListParsed(ctx, path.Root("last_updated__lte"), p.LastUpdated_Lte, utils.TimeFromStringValue, diags),
		LastUpdatedN:     utils.FromListParsed(ctx, path.Root("last_updated__n"), p.LastUpdated_N, utils.TimeFromStringValue, diags),

		Limit: utils.FromInt64Value(p.Limit),

		Name:       utils.FromListValue[string](ctx, path.Root("name"), p.Name, diags),
		NameIc:     utils.FromListValue[string](ctx, path.Root("name__ic"), p.Name_Ic, diags),
		NameIe:     utils.FromListValue[string](ctx, path.Root("name__ie"), p.Name_Ie, diags),
		NameIew:    utils.FromListValue[string](ctx, path.Root("name__iew"), p.Name_Iew, diags),
		NameIregex: utils.FromListValue[string](ctx, path.Root("name__iregex"), p.Name_Iregex, diags),
		NameIsw:    utils.FromListValue[string](ctx, path.Root("name__isw"), p.Name_Isw, diags),
		NameN:      utils.FromListValue[string](ctx, path.Root("name__n"), p.Name_N, diags),
		NameNic:    utils.FromListValue[string](ctx, path.Root("name__nic"), p.Name_Nic, diags),
		NameNie:    utils.FromListValue[string](ctx, path.Root("name__nie"), p.Name_Nie, diags),
		NameNiew:   utils.FromListValue[string](ctx, path.Root("name__niew"), p.Name_Niew, diags),
		NameNisw:   utils.FromListValue[string](ctx, path.Root("name__nisw"), p.Name_Nisw, diags),
		NameRegex:  utils.FromListValue[string](ctx, path.Root("name__regex"), p.Name_Regex, diags),
		NameEmpty:  utils.FromBoolValue(p.Name_Empty),

		Offset:   utils.FromInt64Value(p.Offset),
		Ordering: utils.FromStringValue(p.Ordering),
		Q:        utils.FromStringValue(p.Q),

		Tag:    utils.FromListValue[string](ctx, path.Root("tag"), p.Tag, diags),
		TagN:   utils.FromListValue[string](ctx, path.Root("tag__n"), p.Tag_N, diags),
		TagId:  utils.FromListValue[string](ctx, path.Root("tag_id"), p.TagId, diags),
		TagIdN: utils.FromListValue[string](ctx, path.Root("tag_id__n"), p.TagId_N, diags),

		CreatedByRequest:  utils.FromStringParsed(ctx, path.Root("created_by_request"), p.CreatedByRequest, uuid.Parse, diags),
		ModifiedByRequest: utils.FromStringParsed(ctx, path.Root("modified_by_request"), p.ModifiedByRequest, uuid.Parse, diags),
		UpdatedByRequest:  utils.FromStringParsed(ctx, path.Root("updated_by_request"), p.UpdatedByRequest, uuid.Parse, diags),
	}

	return params
}
