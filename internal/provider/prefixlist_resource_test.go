package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPrefixlistResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "netboxbgp_prefixlist" "test" {
						name        = "some-neighbor"
						family      = "ipv4"
						description = "Prefix belonging to that neighbor"
						comments    = "on some IX"
						tags        = []
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_prefixlist.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_prefixlist.test", "name", "some-neighbor"),
				),
			},
			{
				ResourceName:      "netboxbgp_prefixlist.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
