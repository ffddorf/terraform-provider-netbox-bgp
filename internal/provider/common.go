package provider

import (
	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NestedSite struct {
	Display types.String `tfsdk:"display"`
	ID      types.Int64  `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Slug    types.String `tfsdk:"slug"`
	URL     types.String `tfsdk:"url"`
}

type NestedASN struct {
	ASN     types.Int64  `tfsdk:"asn"`
	Display types.String `tfsdk:"display"`
	ID      types.Int64  `tfsdk:"id"`
	URL     types.String `tfsdk:"url"`
}

type NestedIPAddress struct {
	Address types.String `tfsdk:"address"`
	Display types.String `tfsdk:"display"`
	Family  types.Int64  `tfsdk:"family"`
	ID      types.Int64  `tfsdk:"id"`
	URL     types.String `tfsdk:"url"`
}

type NestedDevice struct {
	Display types.String `tfsdk:"display"`
	ID      types.Int64  `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	URL     types.String `tfsdk:"url"`
}

type NestedBGPPeerGroup struct {
	Description types.String `tfsdk:"description"`
	Display     types.String `tfsdk:"display"`
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	URL         types.String `tfsdk:"url"`
}

type NestedPrefixList struct {
	Display types.String `tfsdk:"display"`
	ID      types.Int64  `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	URL     types.String `tfsdk:"url"`
}

type NestedTenant struct {
	Display types.String `tfsdk:"display"`
	ID      types.Int64  `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Slug    types.String `tfsdk:"slug"`
	URL     types.String `tfsdk:"url"`
}

func (tfo NestedSite) ToAPIModel() client.BriefSite {
	return client.BriefSite{
		Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:     tfo.URL.ValueStringPointer(),
		Display: tfo.Display.ValueStringPointer(),
		Name:    tfo.Name.ValueString(),
		Slug:    tfo.Slug.ValueString(),
	}
}

func (tfo NestedASN) ToAPIModel() client.BriefASN {
	return client.BriefASN{
		Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:     tfo.URL.ValueStringPointer(),
		Display: tfo.Display.ValueStringPointer(),
		Asn:     tfo.ASN.ValueInt64(),
	}
}

func (tfo NestedIPAddress) ToAPIModel() client.BriefIPAddress {
	var fam client.BriefIPAddressFamilyValue
	switch tfo.Family.ValueInt64() {
	case 4:
		fam = client.BriefIPAddressFamilyValueN4
	case 6:
		fam = client.BriefIPAddressFamilyValueN6
	}
	ipa := client.BriefIPAddress{
		Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:     tfo.URL.ValueStringPointer(),
		Display: tfo.Display.ValueStringPointer(),
		Family: &struct {
			Label *client.BriefIPAddressFamilyLabel "json:\"label,omitempty\""
			Value *client.BriefIPAddressFamilyValue "json:\"value,omitempty\""
		}{
			Value: &fam,
		},
		Address: tfo.Address.ValueString(),
	}
	return ipa
}

func (tfo NestedDevice) ToAPIModel() client.BriefDevice {
	return client.BriefDevice{
		Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:     tfo.URL.ValueStringPointer(),
		Display: tfo.Display.ValueStringPointer(),
		Name:    tfo.Name.ValueStringPointer(),
	}
}

func (tfo NestedBGPPeerGroup) ToAPIModel() client.BriefBGPPeerGroup {
	return client.BriefBGPPeerGroup{
		Id:          toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:         tfo.URL.ValueStringPointer(),
		Display:     tfo.Display.ValueStringPointer(),
		Name:        tfo.Name.ValueString(),
		Description: tfo.Description.ValueStringPointer(),
	}
}

func (tfo NestedPrefixList) ToAPIModel() client.BriefPrefixList {
	return client.BriefPrefixList{
		Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:     tfo.URL.ValueStringPointer(),
		Display: tfo.Display.ValueStringPointer(),
		Name:    tfo.Name.ValueString(),
	}
}

func (tfo NestedTenant) ToAPIModel() client.BriefTenant {
	return client.BriefTenant{
		Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:     tfo.URL.ValueStringPointer(),
		Display: tfo.Display.ValueStringPointer(),
		Name:    tfo.Name.ValueString(),
		Slug:    tfo.Slug.ValueString(),
	}
}

func NestedSiteFromAPI(resp *client.BriefSite) *NestedSite {
	if resp == nil {
		return nil
	}
	tfo := &NestedSite{}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.Name = types.StringValue(resp.Name)
	tfo.Slug = types.StringValue(resp.Slug)
	return tfo
}

func NestedASNFromAPI(resp *client.BriefASN) *NestedASN {
	if resp == nil {
		return nil
	}
	tfo := &NestedASN{}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.ASN = types.Int64Value(resp.Asn)
	return tfo
}

func NestedIPAddressFromAPI(resp *client.BriefIPAddress) *NestedIPAddress {
	if resp == nil {
		return nil
	}
	tfo := &NestedIPAddress{}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.Family = maybeInt64Value((*int)(resp.Family.Value))
	tfo.Address = types.StringValue(resp.Address)
	return tfo
}

func NestedDeviceFromAPI(resp *client.BriefDevice) *NestedDevice {
	if resp == nil {
		return nil
	}
	tfo := &NestedDevice{}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.Name = maybeStringValue(resp.Name)
	return tfo
}

func NestedBGPPeerGroupFromAPI(resp *client.BriefBGPPeerGroup) *NestedBGPPeerGroup {
	if resp == nil {
		return nil
	}
	tfo := &NestedBGPPeerGroup{}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.Name = types.StringValue(resp.Name)
	tfo.Description = maybeStringValue(resp.Description)
	return tfo
}

func NestedPrefixListFromAPI(resp *client.BriefPrefixList) *NestedPrefixList {
	if resp == nil {
		return nil
	}
	tfo := &NestedPrefixList{}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.Name = types.StringValue(resp.Name)
	return tfo
}

func NestedTenantFromAPI(resp *client.BriefTenant) *NestedTenant {
	if resp == nil {
		return nil
	}
	tfo := &NestedTenant{}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.Name = types.StringValue(resp.Name)
	tfo.Slug = types.StringValue(resp.Slug)
	return tfo
}

func (*NestedSite) SchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed: true,
		},
		"display": schema.StringAttribute{
			Computed: true,
		},
		"url": schema.StringAttribute{
			Optional: true,
		},
		"name": schema.StringAttribute{
			Required: true,
		},
		"slug": schema.StringAttribute{
			Required: true,
		},
	}
}

func (*NestedASN) SchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed: true,
		},
		"display": schema.StringAttribute{
			Computed: true,
		},
		"url": schema.StringAttribute{
			Optional: true,
		},
		"asn": schema.Int64Attribute{
			Required: true,
		},
	}
}

func (*NestedIPAddress) SchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed: true,
		},
		"display": schema.StringAttribute{
			Computed: true,
		},
		"url": schema.StringAttribute{
			Optional: true,
		},
		"family": schema.Int64Attribute{
			Optional: true,
		},
		"address": schema.StringAttribute{
			Required: true,
		},
	}
}

func (*NestedDevice) SchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed: true,
		},
		"display": schema.StringAttribute{
			Computed: true,
		},
		"url": schema.StringAttribute{
			Optional: true,
		},
		"name": schema.StringAttribute{
			Optional: true,
		},
	}
}

func (*NestedBGPPeerGroup) SchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed: true,
		},
		"display": schema.StringAttribute{
			Computed: true,
		},
		"url": schema.StringAttribute{
			Optional: true,
		},
		"name": schema.StringAttribute{
			Required: true,
		},
		"description": schema.StringAttribute{
			Optional: true,
		},
	}
}

func (*NestedPrefixList) SchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed: true,
		},
		"display": schema.StringAttribute{
			Computed: true,
		},
		"url": schema.StringAttribute{
			Optional: true,
		},
		"name": schema.StringAttribute{
			Required: true,
		},
	}
}

func (*NestedTenant) SchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed: true,
		},
		"display": schema.StringAttribute{
			Computed: true,
		},
		"url": schema.StringAttribute{
			Optional: true,
		},
		"name": schema.StringAttribute{
			Required: true,
		},
		"slug": schema.StringAttribute{
			Required: true,
		},
	}
}
