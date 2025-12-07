package datasources

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
