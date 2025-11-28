package resource_routingpolicy

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ utils.APIConvertibleModel[client.RoutingPolicyRequest, client.RoutingPolicy] = (*RoutingpolicyModel)(nil)

// ToAPIModel implements utils.APIConvertibleModel.
func (r *RoutingpolicyModel) ToAPIModel(ctx context.Context, diags diag.Diagnostics) client.RoutingPolicyRequest {
	return client.RoutingPolicyRequest{
		Comments:    utils.FromStringValue(r.Comments),
		Description: utils.FromStringValue(r.Description),
		Name:        r.Name.ValueString(),
		Tags:        utils.TagsForAPIModel(ctx, r.Tags, diags),
	}
}

// FillFromAPIModel implements utils.APIConvertibleModel.
func (r *RoutingpolicyModel) FillFromAPIModel(ctx context.Context, p *client.RoutingPolicy, diags diag.Diagnostics) {
	*r = RoutingpolicyModel{
		Comments:    utils.MaybeStringValue(p.Comments),
		Description: utils.MaybeStringValue(p.Description),
		Display:     utils.MaybeStringValue(p.Display),
		Id:          types.Int64Value(int64(*p.Id)),
		Name:        types.StringValue(p.Name),
		Tags:        utils.TagsFromAPI(ctx, p.Tags, diags),
		Url:         utils.MaybeStringValue(p.Url),
	}
}
