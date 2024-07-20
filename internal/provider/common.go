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

func (tfo NestedSite) ToAPIModel() client.NestedSite {
	return client.NestedSite{
		Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:     tfo.URL.ValueStringPointer(),
		Display: tfo.Display.ValueStringPointer(),
		Name:    tfo.Name.ValueString(),
		Slug:    tfo.Slug.ValueString(),
	}
}

func (tfo NestedASN) ToAPIModel() client.NestedASN {
	return client.NestedASN{
		Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:     tfo.URL.ValueStringPointer(),
		Display: tfo.Display.ValueStringPointer(),
		Asn:     tfo.ASN.ValueInt64(),
	}
}

func (tfo NestedIPAddress) ToAPIModel() client.NestedIPAddress {
	return client.NestedIPAddress{
		Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:     tfo.URL.ValueStringPointer(),
		Display: tfo.Display.ValueStringPointer(),
		Family:  toIntPointer(tfo.Family.ValueInt64Pointer()),
		Address: tfo.Address.ValueString(),
	}
}

func (tfo NestedDevice) ToAPIModel() client.NestedDevice {
	return client.NestedDevice{
		Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:     tfo.URL.ValueStringPointer(),
		Display: tfo.Display.ValueStringPointer(),
		Name:    tfo.Name.ValueStringPointer(),
	}
}

func (tfo NestedBGPPeerGroup) ToAPIModel() client.NestedBGPPeerGroup {
	return client.NestedBGPPeerGroup{
		Id:          toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:         tfo.URL.ValueStringPointer(),
		Display:     tfo.Display.ValueStringPointer(),
		Name:        tfo.Name.ValueString(),
		Description: tfo.Description.ValueStringPointer(),
	}
}

func (tfo NestedPrefixList) ToAPIModel() client.NestedPrefixList {
	return client.NestedPrefixList{
		Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:     tfo.URL.ValueStringPointer(),
		Display: tfo.Display.ValueStringPointer(),
		Name:    tfo.Name.ValueString(),
	}
}

func (tfo NestedTenant) ToAPIModel() client.NestedTenant {
	return client.NestedTenant{
		Id:      toIntPointer(tfo.ID.ValueInt64Pointer()),
		Url:     tfo.URL.ValueStringPointer(),
		Display: tfo.Display.ValueStringPointer(),
		Name:    tfo.Name.ValueString(),
		Slug:    tfo.Slug.ValueString(),
	}
}

func (tfo *NestedSite) FillFromAPI(resp *client.NestedSite) {
	if resp == nil {
		return
	}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.Name = types.StringValue(resp.Name)
	tfo.Slug = types.StringValue(resp.Slug)
}

func (tfo *NestedASN) FillFromAPI(resp *client.NestedASN) {
	if resp == nil {
		return
	}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.ASN = types.Int64Value(resp.Asn)
}

func (tfo *NestedIPAddress) FillFromAPI(resp *client.NestedIPAddress) {
	if resp == nil {
		return
	}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.Family = maybeInt64Value(resp.Family)
	tfo.Address = types.StringValue(resp.Address)
}

func (tfo *NestedDevice) FillFromAPI(resp *client.NestedDevice) {
	if resp == nil {
		return
	}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.Name = maybeStringValue(resp.Name)
}

func (tfo *NestedBGPPeerGroup) FillFromAPI(resp *client.NestedBGPPeerGroup) {
	if resp == nil {
		return
	}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.Name = types.StringValue(resp.Name)
	tfo.Description = maybeStringValue(resp.Description)
}

func (tfo *NestedPrefixList) FillFromAPI(resp *client.NestedPrefixList) {
	if resp == nil {
		return
	}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.Name = types.StringValue(resp.Name)
}

func (tfo *NestedTenant) FillFromAPI(resp *client.NestedTenant) {
	if resp == nil {
		return
	}
	tfo.ID = types.Int64Value(int64(*resp.Id))
	tfo.URL = maybeStringValue(resp.Url)
	tfo.Display = maybeStringValue(resp.Display)
	tfo.Name = types.StringValue(resp.Name)
	tfo.Slug = types.StringValue(resp.Slug)
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
		"name": schema.StringAttribute{
			Required: true,
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
