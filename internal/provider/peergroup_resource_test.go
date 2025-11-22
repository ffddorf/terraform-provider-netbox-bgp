package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPeergroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        testExternalProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "netboxbgp_peergroup" "test" {
						name        = "Peering Partner"
						description = "Someone we peer with"
						comments    = "This is one some IX"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_peergroup.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_peergroup.test", "name", "Peering Partner"),
				),
			},
			{
				ResourceName:      "netboxbgp_peergroup.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
