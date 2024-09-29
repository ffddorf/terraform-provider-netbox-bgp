package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//go:generate sh -c "go run github.com/ffddorf/terraform-provider-netbox-bgp/cmd/gen-filters PluginsBgpBgpsessionListParams > bgpsessions_filters.gen.go && go run golang.org/x/tools/cmd/goimports -w bgpsessions_filters.gen.go"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SessionsDataSource{}

func NewSessionsDataSource() datasource.DataSource {
	return &SessionsDataSource{}
}

type SessionsDataSource struct {
	client *client.Client
}

type SessionsDataSourceModel struct {
	Filters  Filters                  `tfsdk:"filters"`
	Limit    types.Int64              `tfsdk:"limit"`
	Ordering types.String             `tfsdk:"ordering"`
	Sessions []SessionDataSourceModel `tfsdk:"sessions"`
}

func (d *SessionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sessions"
}

func (d *SessionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	sessionAttrs := map[string]attr.Type{}
	for attrName, attrSchema := range sessionDataSchema {
		sessionAttrs[attrName] = attrSchema.GetType()
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to query for multiple BGP sessions by arbitrary parameters",
		Attributes: map[string]schema.Attribute{
			"filters": schema.ListNestedAttribute{
				NestedObject: FiltersSchema(BgpsessionListParamsFields, BgpsessionListParamsOperators),
				Optional:     true,
			},
			"limit": schema.Int64Attribute{
				Optional: true,
			},
			"ordering": schema.StringAttribute{
				Optional: true,
			},
			"sessions": schema.ListAttribute{
				ElementType: types.ObjectType{
					AttrTypes: sessionAttrs,
				},
				Computed: true,
			},
		},
	}
}

func (d *SessionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = data.Client
}

func unexpectedOperator(op, name string) diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Unexpected operator",
		fmt.Sprintf(`The operator "%s" does not work with the field name "%s"`, op, name),
	)
}

func (d *SessionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SessionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// construct filters
	var params client.PluginsBgpBgpsessionListParams
	for i, filter := range data.Filters {
		if d := setBgpsessionListParamsFromFilter(filter, &params); d != nil {
			resp.Diagnostics.Append(diag.WithPath(path.Root("filters").AtListIndex(i), d))
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	params.Limit = fromInt64Value(data.Limit)
	params.Ordering = data.Ordering.ValueStringPointer()

	nextHTTPReq, err := client.NewPluginsBgpBgpsessionListRequest(d.client.Server, &params)
	for nextHTTPReq != nil {
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to create session list request: %s", err))
			return
		}

		var httpRes *http.Response
		httpRes, err = doPlainReq(ctx, nextHTTPReq, d.client)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to retrieve sessions: %s", err))
			return
		}
		var res *client.PluginsBgpBgpsessionListResponse
		res, err = client.ParsePluginsBgpBgpsessionListResponse(httpRes)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("failed to parse sessions: %s", err))
			return
		}
		if res.JSON200 == nil {
			resp.Diagnostics.AddError("Client Error", httpError(httpRes, res.Body))
			return
		}
		if res.JSON200.Results != nil {
			for _, sess := range *res.JSON200.Results {
				m := SessionDataSourceModel{}
				m.FillFromAPIModel(ctx, &sess, resp.Diagnostics)
				if resp.Diagnostics.HasError() {
					return
				}
				data.Sessions = append(data.Sessions, m)
			}
		}

		// if there was a limit configured, only return elements up to the limit
		if !data.Limit.IsNull() && len(data.Sessions) >= int(data.Limit.ValueInt64()) {
			break
		}

		if res.JSON200.Next == nil || *res.JSON200.Next == "" {
			break
		}

		// handle pagination, query next results
		nextHTTPReq, err = http.NewRequest(http.MethodGet, *res.JSON200.Next, nil)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
