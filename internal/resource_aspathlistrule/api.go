package resource_aspathlistrule

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ utils.APIConvertibleModel[client.ASPathListRuleRequest, client.ASPathListRule] = (*AspathlistruleModel)(nil)

// ToAPIModel implements utils.APIConvertibleModel.
func (a *AspathlistruleModel) ToAPIModel(ctx context.Context, diags diag.Diagnostics) client.ASPathListRuleRequest {
	r := client.ASPathListRuleRequest{
		Action:      client.ASPathListRuleRequestAction(a.Action.ValueString()),
		Comments:    utils.FromStringValue(a.Comments),
		Description: utils.FromStringValue(a.Description),
		Index:       int(a.Index.ValueInt64()),
		Pattern:     a.Pattern.ValueString(),
		Tags:        utils.TagsForAPIModel(ctx, a.Tags, diags),
	}
	r.AspathList = *utils.SetForeignID(&r.AspathList, a.AspathList)
	return r
}

// FillFromAPIModel implements utils.APIConvertibleModel.
func (a *AspathlistruleModel) FillFromAPIModel(ctx context.Context, r *client.ASPathListRule, diags diag.Diagnostics) {
	*a = AspathlistruleModel{
		Action:      types.StringValue(string(r.Action)),
		AspathList:  types.Int64Value(int64(*r.AspathList.Id)),
		Comments:    utils.MaybeStringValue(r.Comments),
		Created:     utils.MaybeStringifiedValue(r.Created, utils.TimeString),
		Description: utils.MaybeStringValue(r.Description),
		Display:     utils.MaybeStringValue(r.Display),
		Id:          types.Int64Value(int64(*r.Id)),
		Index:       types.Int64Value(int64(r.Index)),
		LastUpdated: utils.MaybeStringifiedValue(r.LastUpdated, utils.TimeString),
		Pattern:     types.StringValue(r.Pattern),
		Tags:        utils.TagsFromAPI(ctx, r.Tags, diags),
	}
}
