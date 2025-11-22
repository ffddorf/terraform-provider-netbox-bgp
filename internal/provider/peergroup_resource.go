package provider

import (
	"context"
	"net/http"

	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/resource_peergroup"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PeergroupResource{}
var _ resource.ResourceWithImportState = &PeergroupResource{}

func NewPeergroupResource() resource.Resource {
	return &PeergroupResource{}
}

// SessionResource defines the resource implementation.
type PeergroupResource struct {
	client *ProviderClient
}

// Metadata implements resource.Resource.
func (p *PeergroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_peergroup"
}

// Schema implements resource.Resource.
func (p *PeergroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_peergroup.PeergroupResourceSchema(ctx)
}

func (p *PeergroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	p.client = configureResourceClient(req, resp)
}

// Create implements resource.Resource.
func (p *PeergroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_peergroup.PeergroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := data.ToAPIModel(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	parsed, err := p.client.PluginsBgpPeerGroupCreateWithResponse(ctx, params)
	res := MaybeAPIError("failed to create peer group", err, parsed.JSON201, parsed.HTTPResponse, parsed.Body, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	data.FillFromAPIModel(ctx, res, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete implements resource.Resource.
func (p *PeergroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_peergroup.PeergroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parsed, err := p.client.PluginsBgpPeerGroupDestroyWithResponse(ctx, int(data.Id.ValueInt64()))
	toCheck := parsed
	if parsed.StatusCode() != http.StatusNoContent {
		toCheck = nil
	}
	MaybeAPIError("failed to delete peer group", err, toCheck, parsed.HTTPResponse, parsed.Body, resp.Diagnostics)
}

// Update implements resource.Resource.
func (p *PeergroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_peergroup.PeergroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	var id int64
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := data.ToAPIModel(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	parsed, err := p.client.PluginsBgpPeerGroupUpdateWithResponse(ctx, int(id), params)
	res := MaybeAPIError("failed to update peer group", err, parsed.JSON200, parsed.HTTPResponse, parsed.Body, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	data.FillFromAPIModel(ctx, res, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read implements resource.Resource.
func (p *PeergroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_peergroup.PeergroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Internal Error", "Missing ID value")
		return
	}

	parsed, err := p.client.PluginsBgpPeerGroupRetrieveWithResponse(ctx, int(data.Id.ValueInt64()))
	res := MaybeAPIError("failed to fetch peer group", err, parsed.JSON200, parsed.HTTPResponse, parsed.Body, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	data.FillFromAPIModel(ctx, res, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ImportState implements resource.ResourceWithImportState.
func (p *PeergroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByInt64ID(ctx, req, resp)
}
