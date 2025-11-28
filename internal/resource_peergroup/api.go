package resource_peergroup

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ utils.APIConvertibleModel[client.BGPPeerGroupRequest, client.BGPPeerGroup] = (*PeergroupModel)(nil)

func (m *PeergroupModel) ToAPIModel(ctx context.Context, diags diag.Diagnostics) client.BGPPeerGroupRequest {
	r := client.BGPPeerGroupRequest{
		Comments:    utils.FromStringValue(m.Comments),
		Description: utils.FromStringValue(m.Description),
		Name:        m.Name.ValueString(),
	}
	{
		policies, ds := utils.ToIntListPointer(ctx, m.ImportPolicies)
		for _, d := range ds {
			diags.Append(diag.WithPath(path.Root("import_policies"), d))
		}
		r.ImportPolicies = &policies
	}
	{
		policies, ds := utils.ToIntListPointer(ctx, m.ExportPolicies)
		for _, d := range ds {
			diags.Append(diag.WithPath(path.Root("export_policies"), d))
		}
		r.ExportPolicies = &policies
	}
	return r
}

func (m *PeergroupModel) FillFromAPIModel(ctx context.Context, resp *client.BGPPeerGroup, diags diag.Diagnostics) {
	*m = PeergroupModel{
		Id:             utils.MaybeInt64Value(resp.Id),
		Comments:       utils.MaybeStringValue(resp.Comments),
		Description:    utils.MaybeStringValue(resp.Description),
		Display:        utils.MaybeStringValue(resp.Display),
		ExportPolicies: utils.MaybeListValue(ctx, types.Int64Type, path.Root("export_policies"), resp.ExportPolicies, diags),
		ImportPolicies: utils.MaybeListValue(ctx, types.Int64Type, path.Root("import_policies"), resp.ImportPolicies, diags),
		Name:           utils.MaybeStringValue(&resp.Name),
		Url:            utils.MaybeStringValue(resp.Url),
	}
}
