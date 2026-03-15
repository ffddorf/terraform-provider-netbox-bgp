package provider

import (
	"context"

	"github.com/ffddorf/terraform-provider-netbox-bgp/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = (*prefixListsDataSource)(nil)

func NewPrefixListsDataSource() datasource.DataSource {
	return &prefixListsDataSource{}
}

type prefixListsDataSource struct {
	client *ProviderClient
}

func (d *prefixListsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prefix_lists"
}

func (d *prefixListsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasources.PrefixListsDataSourceSchema(ctx)
}

func (d *prefixListsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = configureDataSourceClient(req, resp)
}

func (d *prefixListsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasources.PrefixListsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := data.ToAPIModel(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := d.client.PluginsBgpPrefixListListWithResponse(ctx, params)
	resp.Diagnostics.Append(MaybeAPIError("failed to get prefix lists", err, apiResp.JSON200, apiResp.HTTPResponse, apiResp.Body)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.FillFromAPIModel(ctx, apiResp.JSON200, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
