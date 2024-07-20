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

	if !m.Name.IsNull() {
		p.Name = m.Name.ValueStringPointer()
	}
	if !m.Description.IsNull() {
		p.Description = m.Description.ValueStringPointer()
	}
	if !m.Comments.IsNull() {
		p.Comments = m.Comments.ValueStringPointer()
	}
	if !m.Status.IsNull() {
		status := client.WritableBGPSessionRequestStatus(m.Status.ValueString())
		p.Status = &status
	}
	if !m.SiteID.IsNull() {
		p.Site = toIntPointer(m.SiteID.ValueInt64())
	}
	if !m.TenantID.IsNull() {
		p.Tenant = toIntPointer(m.TenantID.ValueInt64())
	}
	if !m.DeviceID.IsNull() {
		p.Device = toIntPointer(m.DeviceID.ValueInt64())
	}
	if !m.LocalAddressID.IsNull() {
		p.LocalAddress = int(m.LocalAddressID.ValueInt64())
	}
	if !m.RemoteAddressID.IsNull() {
		p.RemoteAddress = int(m.RemoteAddressID.ValueInt64())
	}
	if !m.LocalASID.IsNull() {
		p.LocalAs = int(m.LocalASID.ValueInt64())
	}
	if !m.RemoteASID.IsNull() {
		p.RemoteAs = int(m.RemoteASID.ValueInt64())
	}
	if !m.PeerGroupID.IsNull() {
		p.PeerGroup = toIntPointer(m.PeerGroupID.ValueInt64())
	}
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
	if !m.PrefixListInID.IsNull() {
		p.PrefixListIn = toIntPointer(m.PrefixListInID.ValueInt64())
	}
	if !m.PrefixListOutID.IsNull() {
		p.PrefixListOut = toIntPointer(m.PrefixListOutID.ValueInt64())
	}

	p.Tags = TagsForAPIModel(ctx, m.Tags, diags)

	// todo: custom fields

	return p
}

func (m *SessionResourceModel) FillFromAPIModel(ctx context.Context, resp *client.BGPSession, diags diag.Diagnostics) {
	if resp.Id != nil {
		m.ID = types.Int64Value(int64(*resp.Id))
	}
	if resp.Comments != nil && *resp.Comments != "" {
		m.Comments = types.StringPointerValue(resp.Comments)
	}
	if resp.Description != nil && *resp.Description != "" {
		m.Description = types.StringPointerValue(resp.Description)
	}
	if resp.Device != nil {
		m.DeviceID = types.Int64Value(int64(*resp.Device.Id))
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
	m.LocalAddressID = types.Int64Value(int64(*resp.LocalAddress.Id))
	m.LocalASID = types.Int64Value(int64(*resp.LocalAs.Id))
	if resp.Name != nil {
		m.Name = types.StringPointerValue(resp.Name)
	}
	if resp.PeerGroup != nil {
		m.PeerGroupID = types.Int64Value(int64(*resp.PeerGroup.Id))
	}
	if resp.PrefixListIn != nil {
		m.PrefixListInID = types.Int64Value(int64(*resp.PrefixListIn.Id))
	}
	if resp.PrefixListOut != nil {
		m.PrefixListOutID = types.Int64Value(int64(*resp.PrefixListOut.Id))
	}
	m.RemoteAddressID = types.Int64Value(int64(*resp.RemoteAddress.Id))
	m.RemoteASID = types.Int64Value(int64(*resp.RemoteAs.Id))
	if resp.Site != nil {
		m.SiteID = types.Int64Value(int64(*resp.Site.Id))
	}
	if resp.Status != nil {
		m.Status = types.StringPointerValue((*string)(resp.Status.Value))
	}
	if resp.Tenant != nil {
		m.TenantID = types.Int64Value(int64(*resp.Tenant.Id))
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
