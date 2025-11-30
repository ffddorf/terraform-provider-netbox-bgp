package resource_routing_policy_rule

import (
	"context"
	"encoding/json"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ utils.APIConvertibleModel[client.RoutingPolicyRuleRequest, client.RoutingPolicyRule] = (*RoutingPolicyRuleModel)(nil)

// ToAPIModel implements utils.APIConvertibleModel.
func (r *RoutingPolicyRuleModel) ToAPIModel(ctx context.Context, diags diag.Diagnostics) client.RoutingPolicyRuleRequest {
	req := client.RoutingPolicyRuleRequest{
		Action:             client.RoutingPolicyRuleRequestAction(r.Action.ValueString()),
		Comments:           utils.FromStringValue(r.Comments),
		ContinueEntry:      utils.FromInt64Value(r.ContinueEntry),
		Description:        utils.FromStringValue(r.Description),
		Index:              int(r.Index.ValueInt64()),
		MatchAspathList:    utils.ToIntListPointer(ctx, r.MatchAspathList, path.Root("match_aspath_list"), diags),
		MatchCommunity:     utils.ToIntListPointer(ctx, r.MatchCommunity, path.Root("match_community"), diags),
		MatchCommunityList: utils.ToIntListPointer(ctx, r.MatchCommunityList, path.Root("match_community_list"), diags),
		MatchCustom:        utils.MaybeRawJSON(r.MatchCustom, path.Root("match_custom"), diags),
		MatchIpAddress:     utils.ToIntListPointer(ctx, r.MatchIpAddress, path.Root("match_ip_address"), diags),
		MatchIpv6Address:   utils.ToIntListPointer(ctx, r.MatchIpv6Address, path.Root("match_ipv6_address"), diags),
		SetActions:         utils.MaybeRawJSON(r.SetActions, path.Root("set_actions"), diags),
		Tags:               utils.TagsForAPIModel(ctx, r.Tags, diags),
	}
	req.RoutingPolicy = *utils.SetForeignID(&req.RoutingPolicy, r.RoutingPolicy)
	return req
}

// FillFromAPIModel implements utils.APIConvertibleModel.
func (r *RoutingPolicyRuleModel) FillFromAPIModel(ctx context.Context, p *client.RoutingPolicyRule, diags diag.Diagnostics) {
	*r = RoutingPolicyRuleModel{
		Action:        types.StringValue(string(p.Action)),
		Comments:      utils.MaybeStringValue(p.Comments),
		ContinueEntry: utils.MaybeInt64Value(p.ContinueEntry),
		Description:   utils.MaybeStringValue(p.Description),
		Display:       utils.MaybeStringValue(p.Display),
		Id:            types.Int64Value(int64(*p.Id)),
		Index:         types.Int64Value(int64(p.Index)),
		MatchAspathList: utils.MaybeListValueAccessor(ctx,
			attr.Type(types.Int64Type),
			path.Root("match_aspath_list"),
			p.MatchAspathList,
			func(in client.ASPathList) int64 { return int64(*in.Id) },
			diags,
		),
		MatchCommunity: utils.MaybeListValueAccessor(ctx,
			attr.Type(types.Int64Type),
			path.Root("match_community"),
			p.MatchCommunity,
			func(in client.Community) int64 { return int64(*in.Id) },
			diags,
		),
		MatchCommunityList: utils.MaybeListValueAccessor(ctx,
			attr.Type(types.Int64Type),
			path.Root("match_community_list"),
			p.MatchCommunityList,
			func(in client.CommunityList) int64 { return int64(*in.Id) },
			diags,
		),
		MatchCustom: utils.MaybeStringifiedValue(p.MatchCustom, func(msg json.RawMessage) string { return string(msg) }),
		MatchIpAddress: utils.MaybeListValueAccessor(ctx,
			attr.Type(types.Int64Type),
			path.Root("match_ip_address"),
			p.MatchIpAddress,
			func(in client.PrefixList) int64 { return int64(*in.Id) },
			diags,
		),
		MatchIpv6Address: utils.MaybeListValueAccessor(ctx,
			attr.Type(types.Int64Type),
			path.Root("match_ipv6_address"),
			p.MatchIpv6Address,
			func(in client.PrefixList) int64 { return int64(*in.Id) },
			diags,
		),
		RoutingPolicy: types.Int64Value(int64(*p.RoutingPolicy.Id)),
		SetActions:    utils.MaybeStringifiedValue(p.SetActions, func(msg json.RawMessage) string { return string(msg) }),
		Tags:          utils.TagsFromAPI(ctx, p.Tags, diags),
	}
}
