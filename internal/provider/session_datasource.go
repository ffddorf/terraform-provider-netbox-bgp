package provider

import (
	"context"
	"fmt"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SessionDataSource{}

func NewSessionDataSource() datasource.DataSource {
	return &SessionDataSource{}
}

type SessionDataSource struct {
	client *ProviderClient
}

type SessionDataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Comments    types.String `tfsdk:"comments"`
	Status      types.String `tfsdk:"status"`

	Site   *NestedSite   `tfsdk:"site"`
	Tenant *NestedTenant `tfsdk:"tenant"`
	Device *NestedDevice `tfsdk:"device"`

	LocalAddress  *NestedIPAddress    `tfsdk:"local_address"`
	RemoteAddress *NestedIPAddress    `tfsdk:"remote_address"`
	LocalAS       *NestedASN          `tfsdk:"local_as"`
	RemoteAS      *NestedASN          `tfsdk:"remote_as"`
	PeerGroup     *NestedBGPPeerGroup `tfsdk:"peer_group"`

	ImportPolicyIDs []types.Int64 `tfsdk:"import_policy_ids"`
	ExportPolicyIDs []types.Int64 `tfsdk:"export_policy_ids"`

	PrefixListIn  *NestedPrefixList `tfsdk:"prefix_list_in"`
	PrefixListOut *NestedPrefixList `tfsdk:"prefix_list_out"`

	Tags types.List `tfsdk:"tags"`
}

func (m *SessionDataSourceModel) FillFromAPIModel(ctx context.Context, resp *client.BGPSession, diags diag.Diagnostics) {
	m.ID = utils.MaybeInt64Value(resp.Id)
	m.Comments = utils.MaybeStringValue(resp.Comments)
	m.Description = utils.MaybeStringValue(resp.Description)
	m.Device = NestedDeviceFromAPI(resp.Device)
	if resp.ExportPolicies != nil {
		for _, policy := range *resp.ExportPolicies {
			m.ExportPolicyIDs = append(m.ExportPolicyIDs, utils.MaybeInt64Value(policy.Id))
		}
	}
	if resp.ImportPolicies != nil && len(*resp.ImportPolicies) > 0 {
		for _, policy := range *resp.ImportPolicies {
			m.ImportPolicyIDs = append(m.ImportPolicyIDs, utils.MaybeInt64Value(policy.Id))
		}
	}
	m.LocalAddress = NestedIPAddressFromAPI(&resp.LocalAddress)
	m.LocalAS = NestedASNFromAPI(&resp.LocalAs)
	m.Name = utils.MaybeStringValue(resp.Name)
	m.PeerGroup = NestedBGPPeerGroupFromAPI(resp.PeerGroup)
	m.PrefixListIn = NestedPrefixListFromAPI(resp.PrefixListIn)
	m.PrefixListOut = NestedPrefixListFromAPI(resp.PrefixListOut)
	m.RemoteAddress = NestedIPAddressFromAPI(&resp.RemoteAddress)
	m.RemoteAS = NestedASNFromAPI(&resp.RemoteAs)
	m.Site = NestedSiteFromAPI(resp.Site)
	m.Status = utils.MaybeStringValue((*string)(resp.Status.Value))
	m.Tenant = NestedTenantFromAPI(resp.Tenant)

	m.Tags = utils.TagsFromAPI(ctx, resp.Tags, diags)

	// todo: custom fields
}

func (d *SessionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_session"
}

var sessionDataSchema = map[string]schema.Attribute{
	"id": schema.Int64Attribute{
		MarkdownDescription: "ID of the resource in Netbox to use for lookup",
		Required:            true,
	},
	"name": schema.StringAttribute{
		Computed: true,
	},
	"description": schema.StringAttribute{
		Computed: true,
	},
	"comments": schema.StringAttribute{
		Computed: true,
	},
	"status": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: `One of: "active", "failed", "offline", "planned"`,
	},
	"site": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: (*NestedSite)(nil).SchemaAttributes(),
	},
	"tenant": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: (*NestedTenant)(nil).SchemaAttributes(),
	},
	"device": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: (*NestedDevice)(nil).SchemaAttributes(),
	},
	"local_address": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: (*NestedIPAddress)(nil).SchemaAttributes(),
	},
	"remote_address": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: (*NestedIPAddress)(nil).SchemaAttributes(),
	},
	"local_as": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: (*NestedASN)(nil).SchemaAttributes(),
	},
	"remote_as": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: (*NestedASN)(nil).SchemaAttributes(),
	},
	"peer_group": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: (*NestedBGPPeerGroup)(nil).SchemaAttributes(),
	},
	"import_policy_ids": schema.ListAttribute{
		ElementType: types.Int64Type,
		Computed:    true,
	},
	"export_policy_ids": schema.ListAttribute{
		ElementType: types.Int64Type,
		Computed:    true,
	},
	"prefix_list_in": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: (*NestedPrefixList)(nil).SchemaAttributes(),
	},
	"prefix_list_out": schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: (*NestedPrefixList)(nil).SchemaAttributes(),
	},
	utils.TagFieldName: utils.TagSchema,
}

func (d *SessionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "BGP Session data source",
		Attributes:          sessionDataSchema,
	}
}

func (d *SessionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureDataSourceClient(req, resp)
}

func (d *SessionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SessionDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	httpRes, err := d.client.PluginsBgpBgpsessionRetrieve(ctx, int(data.ID.ValueInt64()))
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
