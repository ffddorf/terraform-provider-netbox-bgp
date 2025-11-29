package resource_session

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ utils.APIConvertibleModel[client.WritableBGPSessionRequest, client.BGPSession] = (*SessionModel)(nil)

func (m *SessionModel) ToAPIModel(ctx context.Context, diags diag.Diagnostics) client.WritableBGPSessionRequest {
	p := client.WritableBGPSessionRequest{}

	p.Name = utils.FromStringValue(m.Name)
	p.Description = utils.FromStringValue(m.Description)
	p.Comments = utils.FromStringValue(m.Comments)
	if !m.Status.IsNull() {
		status := client.WritableBGPSessionRequestStatus(m.Status.ValueString())
		p.Status = &status
	}
	p.Site = utils.SetForeignID(p.Site, m.Site)
	p.Tenant = utils.SetForeignID(p.Tenant, m.Tenant)
	p.Device = utils.SetForeignID(p.Device, m.Device)
	p.LocalAddress = *utils.SetForeignID(&p.LocalAddress, m.LocalAddress)
	p.RemoteAddress = *utils.SetForeignID(&p.RemoteAddress, m.RemoteAddress)
	p.LocalAs = *utils.SetForeignID(&p.LocalAs, m.LocalAs)
	p.RemoteAs = *utils.SetForeignID(&p.RemoteAs, m.RemoteAs)
	p.PeerGroup = utils.SetForeignID(p.PeerGroup, m.PeerGroup)
	p.ImportPolicies = utils.ToIntListPointer(ctx, m.ImportPolicies, path.Root("import_policies"), diags)
	p.ExportPolicies = utils.ToIntListPointer(ctx, m.ExportPolicies, path.Root("export_policies"), diags)
	utils.SetForeignID(p.PrefixListIn, m.PrefixListIn)
	utils.SetForeignID(p.PrefixListOut, m.PrefixListOut)
	utils.SetForeignID(p.Virtualmachine, m.Virtualmachine)

	p.Tags = utils.TagsForAPIModel(ctx, m.Tags, diags)

	// todo: custom fields

	return p
}

func (m *SessionModel) FillFromAPIModel(ctx context.Context, resp *client.BGPSession, diags diag.Diagnostics) {
	m.Id = utils.MaybeInt64Value(resp.Id)
	m.Comments = utils.MaybeStringValue(resp.Comments)
	m.Created = utils.MaybeStringifiedValue(resp.Created, utils.TimeString)
	m.Description = utils.MaybeStringValue(resp.Description)
	m.Device = utils.MaybeInt64ValueSubfield(resp.Device, func(d client.BriefDevice) *int { return d.Id })
	m.Display = utils.MaybeStringValue(resp.Display)
	m.ExportPolicies = utils.MaybeListValueAccessor(ctx,
		types.Int64Type,
		path.Root("export_policies"),
		resp.ExportPolicies,
		func(in client.RoutingPolicy) int64 { return int64(*in.Id) },
		diags,
	)
	m.ImportPolicies = utils.MaybeListValueAccessor(ctx,
		types.Int64Type,
		path.Root("import_policies"),
		resp.ImportPolicies,
		func(in client.RoutingPolicy) int64 { return int64(*in.Id) },
		diags,
	)
	m.LastUpdated = utils.MaybeStringifiedValue(resp.LastUpdated, utils.TimeString)
	m.LocalAddress = utils.MaybeInt64Value(resp.LocalAddress.Id)
	m.LocalAs = utils.MaybeInt64Value(resp.LocalAs.Id)
	m.Name = utils.MaybeStringValue(resp.Name)
	m.PeerGroup = utils.MaybeInt64ValueSubfield(resp.PeerGroup, func(pg client.BriefBGPPeerGroup) *int { return pg.Id })
	m.PrefixListIn = utils.MaybeInt64ValueSubfield(resp.PrefixListIn, func(pfxl client.BriefPrefixList) *int { return pfxl.Id })
	m.PrefixListOut = utils.MaybeInt64ValueSubfield(resp.PrefixListOut, func(pfxl client.BriefPrefixList) *int { return pfxl.Id })
	m.RemoteAddress = utils.MaybeInt64Value(resp.RemoteAddress.Id)
	m.RemoteAs = utils.MaybeInt64Value(resp.RemoteAs.Id)
	m.Site = utils.MaybeInt64ValueSubfield(resp.Site, func(s client.BriefSite) *int { return s.Id })
	m.Status = utils.MaybeStringValue((*string)(resp.Status.Value))
	m.Tenant = utils.MaybeInt64ValueSubfield(resp.Tenant, func(t client.BriefTenant) *int { return t.Id })
	m.Virtualmachine = utils.MaybeInt64ValueSubfield(resp.Virtualmachine, func(vm client.BriefVirtualMachine) *int { return vm.Id })
	m.Url = utils.MaybeStringValue(resp.Url)

	m.Tags = utils.TagsFromAPI(ctx, resp.Tags, diags)

	// todo: custom fields
}
