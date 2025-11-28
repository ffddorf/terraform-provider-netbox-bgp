package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoutingpolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        testExternalProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "netboxbgp_routingpolicy" "test" {
						name        = "Public peering"
						description = "Policy towards public peers"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_routingpolicy.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_routingpolicy.test", "name", "Public peering"),
				),
			},
			{
				ResourceName:      "netboxbgp_routingpolicy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
