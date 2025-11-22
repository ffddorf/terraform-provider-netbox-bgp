package provider

import (
	"context"
	"net/http"

	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/resource_session"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SessionResource{}
var _ resource.ResourceWithImportState = &SessionResource{}

func NewSessionResource() resource.Resource {
	return &SessionResource{}
}

// SessionResource defines the resource implementation.
type SessionResource struct {
	client *ProviderClient
}

func (r *SessionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session"
}

func (r *SessionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_session.SessionResourceSchema(ctx)
}

func (r *SessionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = configureResourceClient(req, resp)
}

func (r *SessionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_session.SessionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := data.ToAPIModel(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	parsed, err := r.client.PluginsBgpSessionCreateWithResponse(ctx, params)
	res := MaybeAPIError("failed to create session", err, parsed.JSON201, parsed.HTTPResponse, parsed.Body, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	data.FillFromAPIModel(ctx, res, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SessionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_session.SessionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Internal Error", "Missing ID value")
		return
	}

	parsed, err := r.client.PluginsBgpSessionRetrieveWithResponse(ctx, int(data.Id.ValueInt64()))
	res := MaybeAPIError("failed to fetch session", err, parsed.JSON200, parsed.HTTPResponse, parsed.Body, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	data.FillFromAPIModel(ctx, res, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SessionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_session.SessionModel
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
	parsed, err := r.client.PluginsBgpSessionUpdateWithResponse(ctx, int(id), params)
	res := MaybeAPIError("failed to update session", err, parsed.JSON200, parsed.HTTPResponse, parsed.Body, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	data.FillFromAPIModel(ctx, res, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SessionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_session.SessionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parsed, err := r.client.PluginsBgpSessionDestroyWithResponse(ctx, int(data.Id.ValueInt64()))
	toCheck := parsed
	if parsed.StatusCode() != http.StatusNoContent {
		toCheck = nil // response not usable
	}
	MaybeAPIError("failed to delete session", err, toCheck, parsed.HTTPResponse, parsed.Body, resp.Diagnostics)
}

func (r *SessionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByInt64ID(ctx, req, resp)
}
