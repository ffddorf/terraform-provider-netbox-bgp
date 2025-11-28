package resource_prefixlistrule

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ utils.APIConvertibleModel[client.PrefixListRuleRequest, client.PrefixListRule] = (*PrefixlistruleModel)(nil)

func (p *PrefixlistruleModel) ToAPIModel(ctx context.Context, diags diag.Diagnostics) client.PrefixListRuleRequest {
	r := client.PrefixListRuleRequest{
		Action:       client.PrefixListRuleRequestAction(p.Action.ValueString()),
		Comments:     utils.FromStringValue(p.Comments),
		Description:  utils.FromStringValue(p.Description),
		Ge:           utils.FromIntValue(p.Ge),
		Index:        int(p.Index.ValueInt64()),
		Le:           utils.FromIntValue(p.Le),
		PrefixCustom: utils.FromStringValue(p.PrefixCustom),
		Tags:         utils.TagsForAPIModel(ctx, p.Tags, diags),
	}
	r.Prefix = utils.SetForeignID(r.Prefix, p.Prefix)
	r.PrefixList = *utils.SetForeignID(&r.PrefixList, p.PrefixList)
	return r
}

func (p *PrefixlistruleModel) FillFromAPIModel(ctx context.Context, resp *client.PrefixListRule, diags diag.Diagnostics) {
	*p = PrefixlistruleModel{
		Action:       types.StringValue(string(resp.Action)),
		Comments:     utils.MaybeStringValue(resp.Comments),
		Created:      utils.MaybeStringifiedValue(resp.Created, utils.TimeString),
		Description:  utils.MaybeStringValue(resp.Description),
		Display:      utils.MaybeStringValue(resp.Display),
		Ge:           utils.MaybeInt64Value(resp.Ge),
		Id:           types.Int64Value(int64(*resp.Id)),
		Index:        types.Int64Value(int64(resp.Index)),
		LastUpdated:  utils.MaybeStringifiedValue(resp.LastUpdated, utils.TimeString),
		Le:           utils.MaybeInt64Value(resp.Le),
		Prefix:       utils.MaybeInt64ValueSubfield(resp.Prefix, func(p client.BriefPrefix) *int { return p.Id }),
		PrefixCustom: utils.MaybeStringValue(resp.PrefixCustom),
		PrefixList:   types.Int64Value(int64(*resp.PrefixList.Id)),
		Tags:         utils.TagsFromAPI(ctx, resp.Tags, diags),
	}
}
