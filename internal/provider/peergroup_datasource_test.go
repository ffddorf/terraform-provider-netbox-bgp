package provider

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAccPeerGroupDataSource(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	c := testClient(t)
	comments := "This will be great!"
	description := "For testing"
	r, err := c.PluginsBgpBgppeergroupCreate(ctx, client.BGPPeerGroupRequest{
		Name:        "Example Peer Group",
		Description: &description,
		Comments:    &comments,
	})
	require.NoError(t, err)
	pg, err := client.ParsePluginsBgpBgppeergroupCreateResponse(r)
	require.NoError(t, err)
	require.NotNil(t, pg.JSON201, "bad API response: %d", r.StatusCode)

	peerGroupID := *pg.JSON201.Id
	t.Cleanup(func() {
		_, _ = c.PluginsBgpBgppeergroupDestroy(ctx, peerGroupID)
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        testExternalProviders,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: fmt.Sprintf(`
					data "netboxbgp_peergroup" "test" {
						id = %d
					}
				`, peerGroupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netboxbgp_peergroup.test", "id", strconv.Itoa(peerGroupID)),
					resource.TestCheckResourceAttr("data.netboxbgp_peergroup.test", "name", "Example Peer Group"),
					resource.TestCheckResourceAttr("data.netboxbgp_peergroup.test", "description", "For testing"),
				),
			},
		},
	})
}
