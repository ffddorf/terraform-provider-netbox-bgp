package provider

import (
	"context"
	"fmt"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &SessionDataSource{}

func NewPeerGroupDataSource() datasource.DataSource {
	return &PeerGroupDataSource{}
}

type PeerGroupDataSource struct {
	client *client.Client
}

type PeerGroupDataSourceModel struct {
	ID              types.Int64   `tfsdk:"id"`
	Name            types.String  `tfsdk:"name"`
	Description     types.String  `tfsdk:"description"`
	Comments        types.String  `tfsdk:"comments"`
	ExportPolicyIDs []types.Int64 `tfsdk:"export_policy_ids"`
	ImportPolicyIDs []types.Int64 `tfsdk:"import_policy_ids"`
}

func (pgm *PeerGroupDataSourceModel) FillFromAPIModel(resp *client.BGPPeerGroup) {
	pgm.ID = maybeInt64Value(resp.Id)
	pgm.Name = maybeStringValue(&resp.Name)
	pgm.Description = maybeStringValue(resp.Description)
	pgm.Comments = maybeStringValue(resp.Comments)

	if resp.ExportPolicies != nil {
		for _, policy := range *resp.ExportPolicies {
			pgm.ExportPolicyIDs = append(pgm.ExportPolicyIDs, maybeInt64Value(policy.Id))
		}
	}
	if resp.ImportPolicies != nil && len(*resp.ImportPolicies) > 0 {
		for _, policy := range *resp.ImportPolicies {
			pgm.ImportPolicyIDs = append(pgm.ImportPolicyIDs, maybeInt64Value(policy.Id))
		}
	}
}

func (pgd *PeerGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_peergroup"
}

func (pgd *PeerGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	pgd.client = configureDataSourceClient(req, resp)
}

var peerGroupDataSchema = map[string]schema.Attribute{
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
	"import_policy_ids": schema.ListAttribute{
		ElementType: types.Int64Type,
		Computed:    true,
	},
	"export_policy_ids": schema.ListAttribute{
		ElementType: types.Int64Type,
		Computed:    true,
	},
}

func (pgd *PeerGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "BGP peer group data source",
		Attributes:          peerGroupDataSchema,
	}
}

func (pgd *PeerGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PeerGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	httpRes, err := pgd.client.PluginsBgpBgppeergroupRetrieve(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to retrieve session: %s", err))
		return
	}
	res, err := client.ParsePluginsBgpBgppeergroupRetrieveResponse(httpRes)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse session: %s", err))
		return
	}
	if res.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", httpError(httpRes, res.Body))
		return
	}

	data.FillFromAPIModel(res.JSON200)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
