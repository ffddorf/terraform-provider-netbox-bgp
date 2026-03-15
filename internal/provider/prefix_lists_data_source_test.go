package provider

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccPrefixListsDataSource(t *testing.T) {
	ctx := t.Context()
	c := testClient(t)

	name := testName(t)
	r, err := c.PluginsBgpPrefixListCreateWithResponse(ctx, client.PrefixListRequest{
		Comments:    ptr("This will be great!"),
		Description: ptr("For testing"),
		Family:      client.PrefixListRequestFamilyIpv6,
		Name:        name,
	})
	require.NoError(t, err)
	require.NotNilf(t, r.JSON201, "bad API response: %d: %s", r.StatusCode(), string(r.Body))

	prefixListID := *r.JSON201.Id
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()
		_, err = c.PluginsBgpPrefixListDestroy(ctx, prefixListID)
		assert.NoError(t, err)
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: fmt.Sprintf(`
					data "netboxbgp_prefix_lists" "any" {
						name = ["%s"]
					}
				`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netboxbgp_prefix_lists.any", "results.0.id", strconv.Itoa(prefixListID)),
					resource.TestCheckResourceAttr("data.netboxbgp_prefix_lists.any", "results.0.name", name),
					resource.TestCheckResourceAttr("data.netboxbgp_prefix_lists.any", "results.0.family", "ipv6"),
				),
			},
		},
	})
}
