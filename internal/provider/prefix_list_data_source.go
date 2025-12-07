package provider

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = (*prefixListDataSource)(nil)

func NewPrefixListDataSource() datasource.DataSource {
	return &prefixListDataSource{}
}

type prefixListDataSource struct {
	client *ProviderClient
}

func (d *prefixListDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prefix_list"
}

func (d *prefixListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasources.PrefixListDataSourceSchema(ctx)
}

func (r *prefixListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.client = configureDataSourceClient(req, resp)
}

func (d *prefixListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasources.PrefixListModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parsed, err := d.client.PluginsBgpPrefixListRetrieveWithResponse(ctx, int(data.Id.ValueInt64()))
	resp.Diagnostics.Append(MaybeAPIError("failed to fetch prefix_list", err, parsed.JSON200, parsed.HTTPResponse, parsed.Body)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.FillFromAPIModel(ctx, parsed.JSON200, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
