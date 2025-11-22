package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPrefixlistResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        testExternalProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "netboxbgp_prefixlist" "test" {
						name        = "rfc1918"
						family      = "ipv4"
						description = "RFC 1918 prefixes"
						comments    = "do not announce publicly"
						tags        = []
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_prefixlist.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_prefixlist.test", "name", "rfc1918"),
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
