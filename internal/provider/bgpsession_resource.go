package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SessionResource{}
var _ resource.ResourceWithImportState = &SessionResource{}

func NewSessionResource() resource.Resource {
	return &SessionResource{}
}

// SessionResource defines the resource implementation.
type SessionResource struct {
	client *client.Client
}

// SessionResourceModel describes the resource data model.
type SessionResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Comments    types.String `tfsdk:"comments"`
	Status      types.String `tfsdk:"status"`

	SiteID   types.Int64 `tfsdk:"site_id"`
	TenantID types.Int64 `tfsdk:"tenant_id"`
	DeviceID types.Int64 `tfsdk:"device_id"`

	LocalAddressID  types.Int64 `tfsdk:"local_address_id"`
	RemoteAddressID types.Int64 `tfsdk:"remote_address_id"`
	LocalASID       types.Int64 `tfsdk:"local_as_id"`
	RemoteASID      types.Int64 `tfsdk:"remote_as_id"`
	PeerGroupID     types.Int64 `tfsdk:"peer_group_id"`

	ImportPolicyIDs types.List `tfsdk:"import_policy_ids"`
	ExportPolicyIDs types.List `tfsdk:"export_policy_ids"`

	PrefixListInID  types.Int64 `tfsdk:"prefix_list_in_id"`
	PrefixListOutID types.Int64 `tfsdk:"prefix_list_out_id"`

	Tags types.List `tfsdk:"tags"`

	// todo: custom fields
}

func (m *SessionResourceModel) ToAPIModel(ctx context.Context, diags diag.Diagnostics) client.WritableBGPSessionRequest {
	p := client.WritableBGPSessionRequest{}

	p.Name = m.Name.ValueStringPointer()
	p.Description = m.Description.ValueStringPointer()
	p.Comments = m.Comments.ValueStringPointer()
	if !m.Status.IsNull() {
		status := client.WritableBGPSessionRequestStatus(m.Status.ValueString())
		p.Status = &status
	}
	p.Site = fromInt64Value(m.SiteID)
	p.Tenant = fromInt64Value(m.TenantID)
	p.Device = fromInt64Value(m.DeviceID)
	p.LocalAddress = *fromInt64Value(m.LocalAddressID)
	p.RemoteAddress = *fromInt64Value(m.RemoteAddressID)
	p.LocalAs = *fromInt64Value(m.LocalASID)
	p.RemoteAs = *fromInt64Value(m.RemoteASID)
	p.PeerGroup = fromInt64Value(m.PeerGroupID)
	if !m.ImportPolicyIDs.IsNull() {
		policies, ds := toIntListPointer(ctx, m.ImportPolicyIDs)
		for _, d := range ds {
			diags.Append(diag.WithPath(path.Root("import_policy_ids"), d))
		}
		p.ImportPolicies = &policies
	}
	if !m.ExportPolicyIDs.IsNull() {
		policies, ds := toIntListPointer(ctx, m.ExportPolicyIDs)
		for _, d := range ds {
			diags.Append(diag.WithPath(path.Root("export_policy_ids"), d))
		}
		p.ExportPolicies = &policies
	}
	p.PrefixListIn = fromInt64Value(m.PrefixListInID)
	p.PrefixListOut = fromInt64Value(m.PrefixListOutID)

	p.Tags = TagsForAPIModel(ctx, m.Tags, diags)

	// todo: custom fields

	return p
}

func (m *SessionResourceModel) FillFromAPIModel(ctx context.Context, resp *client.BGPSession, diags diag.Diagnostics) {
	m.ID = maybeInt64Value(resp.Id)
	m.Comments = maybeStringValue(resp.Comments)
	m.Description = maybeStringValue(resp.Description)
	if resp.Device != nil {
		m.DeviceID = maybeInt64Value(resp.Device.Id)
	}
	if resp.ExportPolicies != nil && len(*resp.ExportPolicies) > 0 {
		var ds diag.Diagnostics
		m.ExportPolicyIDs, ds = types.ListValueFrom(ctx, types.Int64Type, resp.ExportPolicies)
		for _, d := range ds {
			diags.Append(diag.WithPath(path.Root("export_policy_ids"), d))
		}
	}
	if resp.ImportPolicies != nil && len(*resp.ImportPolicies) > 0 {
		var ds diag.Diagnostics
		m.ImportPolicyIDs, ds = types.ListValueFrom(ctx, types.Int64Type, resp.ImportPolicies)
		for _, d := range ds {
			diags.Append(diag.WithPath(path.Root("import_policy_ids"), d))
		}
	}
	m.LocalAddressID = maybeInt64Value(resp.LocalAddress.Id)
	m.LocalASID = maybeInt64Value(resp.LocalAs.Id)
	m.Name = maybeStringValue(resp.Name)
	if resp.PeerGroup != nil {
		m.PeerGroupID = maybeInt64Value(resp.PeerGroup.Id)
	}
	if resp.PrefixListIn != nil {
		m.PrefixListInID = maybeInt64Value(resp.PrefixListIn.Id)
	}
	if resp.PrefixListOut != nil {
		m.PrefixListOutID = maybeInt64Value(resp.PrefixListOut.Id)
	}
	m.RemoteAddressID = maybeInt64Value(resp.RemoteAddress.Id)
	m.RemoteASID = maybeInt64Value(resp.RemoteAs.Id)
	if resp.Site != nil {
		m.SiteID = maybeInt64Value(resp.Site.Id)
	}
	m.Status = maybeStringValue((*string)(resp.Status.Value))
	if resp.Tenant != nil {
		m.TenantID = maybeInt64Value(resp.Tenant.Id)
	}

	m.Tags = TagsFromAPI(ctx, resp.Tags, diags)

	// todo: custom fields
}

func (r *SessionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session"
}

func (r *SessionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Session resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "ID of the resource in Netbox",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"comments": schema.StringAttribute{
				Optional: true,
			},
			"status": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(client.BGPSessionStatusValueActive),
						string(client.BGPSessionStatusValueFailed),
						string(client.BGPSessionStatusValueOffline),
						string(client.BGPSessionStatusValuePlanned),
					),
				},
				MarkdownDescription: `One of: "active", "failed", "offline", "planned"`,
			},
			"site_id": schema.Int64Attribute{
				Optional: true,
			},
			"tenant_id": schema.Int64Attribute{
				Optional: true,
			},
			"device_id": schema.Int64Attribute{
				Required: true,
			},
			"local_address_id": schema.Int64Attribute{
				Required: true,
			},
			"remote_address_id": schema.Int64Attribute{
				Required: true,
			},
			"local_as_id": schema.Int64Attribute{
				Required: true,
			},
			"remote_as_id": schema.Int64Attribute{
				Required: true,
			},
			"peer_group_id": schema.Int64Attribute{
				Optional: true,
			},
			"import_policy_ids": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
			},
			"export_policy_ids": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
			},
			"prefix_list_in_id": schema.Int64Attribute{
				Optional: true,
			},
			"prefix_list_out_id": schema.Int64Attribute{
				Optional: true,
			},
			TagFieldName: TagSchema,
		},
	}
}

func (r *SessionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*configuredProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *configuredProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = data.Client
}

func (r *SessionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SessionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := data.ToAPIModel(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	httpRes, err := r.client.PluginsBgpSessionCreate(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to create session: %s", err))
		return
	}
	res, err := client.ParsePluginsBgpSessionCreateResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse session response: %s", err))
		return
	}
	if res.JSON201 == nil {
		resp.Diagnostics.AddError("Client Error", httpError(httpRes, res.Body))
		return
	}

	data.FillFromAPIModel(ctx, res.JSON201, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "created a resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SessionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SessionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("id"), "Internal Error", "Missing ID value")
		return
	}

	httpRes, err := r.client.PluginsBgpSessionRetrieve(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to retrieve session: %s", err))
		return
	}
	res, err := client.ParsePluginsBgpSessionRetrieveResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse session: %s", err))
		return
	}
	if res.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", httpError(httpRes, res.Body))
		return
	}

	data.FillFromAPIModel(ctx, res.JSON200, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SessionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SessionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := data.ToAPIModel(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	httpRes, err := r.client.PluginsBgpSessionUpdate(ctx, int(data.ID.ValueInt64()), params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to update session: %s", err))
		return
	}
	res, err := client.ParsePluginsBgpBgpsessionUpdateResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse session response: %s", err))
		return
	}
	if res.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", httpError(httpRes, res.Body))
		return
	}

	data.FillFromAPIModel(ctx, res.JSON200, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SessionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SessionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpRes, err := r.client.PluginsBgpSessionDestroy(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to destroy session: %s", err))
		return
	}
	res, err := client.ParsePluginsBgpBgpsessionDestroyResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse response: %s", err))
		return
	}
	if res.StatusCode() != http.StatusNoContent {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to destroy session: %s", string(res.Body)))
		return
	}
}

func (r *SessionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByInt64ID(ctx, req, resp)
}
