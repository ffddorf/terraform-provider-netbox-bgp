package resource_prefix_list

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ utils.APIConvertibleModel[client.PrefixListRequest, client.PrefixList] = (*PrefixListModel)(nil)

func (p *PrefixListModel) ToAPIModel(ctx context.Context, diags diag.Diagnostics) client.PrefixListRequest {
	return client.PrefixListRequest{
		Comments:    utils.FromStringValue(p.Comments),
		Description: utils.FromStringValue(p.Description),
		Family:      client.PrefixListRequestFamily(p.Family.ValueString()),
		Name:        p.Name.ValueString(),
		Tags:        utils.TagsForAPIModel(ctx, p.Tags, diags),
	}
}

func (p *PrefixListModel) FillFromAPIModel(ctx context.Context, resp *client.PrefixList, diags diag.Diagnostics) {
	*p = PrefixListModel{
		Comments:    utils.MaybeStringValue(resp.Comments),
		Description: utils.MaybeStringValue(resp.Description),
		Display:     utils.MaybeStringValue(resp.Display),
		Family:      types.StringValue(string(resp.Family)),
		Id:          types.Int64Value(int64(*resp.Id)),
		Name:        types.StringValue(resp.Name),
		Tags:        utils.TagsFromAPI(ctx, resp.Tags, diags),
		Url:         utils.MaybeStringValue(resp.Url),
	}
}
