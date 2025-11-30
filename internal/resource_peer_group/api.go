package resource_peer_group

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ utils.APIConvertibleModel[client.BGPPeerGroupRequest, client.BGPPeerGroup] = (*PeerGroupModel)(nil)

func (m *PeerGroupModel) ToAPIModel(ctx context.Context, diags diag.Diagnostics) client.BGPPeerGroupRequest {
	r := client.BGPPeerGroupRequest{
		Comments:       utils.FromStringValue(m.Comments),
		Description:    utils.FromStringValue(m.Description),
		Name:           m.Name.ValueString(),
		ExportPolicies: utils.ToIntListPointer(ctx, m.ExportPolicies, path.Root("export_policies"), diags),
		ImportPolicies: utils.ToIntListPointer(ctx, m.ImportPolicies, path.Root("import_policies"), diags),
	}
	return r
}

func (m *PeerGroupModel) FillFromAPIModel(ctx context.Context, resp *client.BGPPeerGroup, diags diag.Diagnostics) {
	*m = PeerGroupModel{
		Id:          utils.MaybeInt64Value(resp.Id),
		Comments:    utils.MaybeStringValue(resp.Comments),
		Description: utils.MaybeStringValue(resp.Description),
		Display:     utils.MaybeStringValue(resp.Display),
		ExportPolicies: utils.MaybeListValueAccessor(ctx,
			types.Int64Type,
			path.Root("export_policies"),
			resp.ExportPolicies,
			func(in client.RoutingPolicy) int64 { return int64(*in.Id) },
			diags,
		),
		ImportPolicies: utils.MaybeListValueAccessor(ctx,
			types.Int64Type,
			path.Root("import_policies"),
			resp.ImportPolicies,
			func(in client.RoutingPolicy) int64 { return int64(*in.Id) },
			diags,
		),
		Name: utils.MaybeStringValue(&resp.Name),
		Url:  utils.MaybeStringValue(resp.Url),
	}
}
