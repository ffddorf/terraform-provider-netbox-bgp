package resource_session

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (m *SessionModel) ToAPIModel(ctx context.Context, diags diag.Diagnostics) client.WritableBGPSessionRequest {
	p := client.WritableBGPSessionRequest{}

	p.Name = m.Name.ValueStringPointer()
	p.Description = m.Description.ValueStringPointer()
	p.Comments = m.Comments.ValueStringPointer()
	if !m.Status.IsNull() {
		status := client.WritableBGPSessionRequestStatus(m.Status.ValueString())
		p.Status = &status
	}
	p.Site = utils.FromInt64Value(m.Site)
	p.Tenant = utils.FromInt64Value(m.Tenant)
	p.Device = utils.FromInt64Value(m.Device)
	p.LocalAddress = *utils.FromInt64Value(m.LocalAddress)
	p.RemoteAddress = *utils.FromInt64Value(m.RemoteAddress)
	p.LocalAs = *utils.FromInt64Value(m.LocalAs)
	p.RemoteAs = *utils.FromInt64Value(m.RemoteAs)
	p.PeerGroup = utils.FromInt64Value(m.PeerGroup)
	{
		policies, ds := utils.ToIntListPointer(ctx, m.ImportPolicies)
		for _, d := range ds {
			diags.Append(diag.WithPath(path.Root("import_policies"), d))
		}
		p.ImportPolicies = &policies
	}
	{
		policies, ds := utils.ToIntListPointer(ctx, m.ExportPolicies)
		for _, d := range ds {
			diags.Append(diag.WithPath(path.Root("export_policies"), d))
		}
		p.ExportPolicies = &policies
	}
	p.PrefixListIn = utils.FromInt64Value(m.PrefixListIn)
	p.PrefixListOut = utils.FromInt64Value(m.PrefixListOut)

	p.Tags = utils.TagsForAPIModel(ctx, m.Tags, diags)

	// todo: custom fields

	return p
}

func (m *SessionModel) FillFromAPIModel(ctx context.Context, resp *client.BGPSession, diags diag.Diagnostics) {
	m.Id = utils.MaybeInt64Value(resp.Id)
	m.Comments = utils.MaybeStringValue(resp.Comments)
	m.Description = utils.MaybeStringValue(resp.Description)
	if resp.Device != nil {
		m.Device = utils.MaybeInt64Value(resp.Device.Id)
	}
	if resp.ExportPolicies != nil && len(*resp.ExportPolicies) > 0 {
		var ds diag.Diagnostics
		m.ExportPolicies, ds = types.ListValueFrom(ctx, types.Int64Type, resp.ExportPolicies)
		for _, d := range ds {
			diags.Append(diag.WithPath(path.Root("export_policy_ids"), d))
		}
	}
	if resp.ImportPolicies != nil && len(*resp.ImportPolicies) > 0 {
		var ds diag.Diagnostics
		m.ImportPolicies, ds = types.ListValueFrom(ctx, types.Int64Type, resp.ImportPolicies)
		for _, d := range ds {
			diags.Append(diag.WithPath(path.Root("import_policy_ids"), d))
		}
	}
	m.LocalAddress = utils.MaybeInt64Value(resp.LocalAddress.Id)
	m.LocalAs = utils.MaybeInt64Value(resp.LocalAs.Id)
	m.Name = utils.MaybeStringValue(resp.Name)
	if resp.PeerGroup != nil {
		m.PeerGroup = utils.MaybeInt64Value(resp.PeerGroup.Id)
	}
	if resp.PrefixListIn != nil {
		m.PrefixListIn = utils.MaybeInt64Value(resp.PrefixListIn.Id)
	}
	if resp.PrefixListOut != nil {
		m.PrefixListOut = utils.MaybeInt64Value(resp.PrefixListOut.Id)
	}
	m.RemoteAddress = utils.MaybeInt64Value(resp.RemoteAddress.Id)
	m.RemoteAs = utils.MaybeInt64Value(resp.RemoteAs.Id)
	if resp.Site != nil {
		m.Site = utils.MaybeInt64Value(resp.Site.Id)
	}
	m.Status = utils.MaybeStringValue((*string)(resp.Status.Value))
	if resp.Tenant != nil {
		m.Tenant = utils.MaybeInt64Value(resp.Tenant.Id)
	}

	m.Tags = utils.TagsFromAPI(ctx, resp.Tags, diags)

	// todo: custom fields
}
