package resource_aspath_list

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ utils.APIConvertibleModel[client.ASPathListRequest, client.ASPathList] = (*AspathListModel)(nil)

// ToAPIModel implements utils.APIConvertibleModel.
func (a *AspathListModel) ToAPIModel(ctx context.Context, diags diag.Diagnostics) client.ASPathListRequest {
	return client.ASPathListRequest{
		Comments:    utils.FromStringValue(a.Comments),
		Description: utils.FromStringValue(a.Description),
		Name:        a.Name.ValueString(),
		Tags:        utils.TagsForAPIModel(ctx, a.Tags, diags),
	}
}

// FillFromAPIModel implements utils.APIConvertibleModel.
func (a *AspathListModel) FillFromAPIModel(ctx context.Context, r *client.ASPathList, diags diag.Diagnostics) {
	*a = AspathListModel{
		Comments:    utils.MaybeStringValue(r.Comments),
		Description: utils.MaybeStringValue(r.Description),
		Display:     utils.MaybeStringValue(r.Display),
		Id:          types.Int64Value(int64(*r.Id)),
		Name:        types.StringValue(r.Name),
		Tags:        utils.TagsFromAPI(ctx, r.Tags, diags),
		Url:         utils.MaybeStringValue(r.Url),
	}
}
