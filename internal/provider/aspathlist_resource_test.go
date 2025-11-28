package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAspathlistResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        testExternalProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "netboxbgp_aspathlist" "test" {
						name        = "Internal"
						description = "Allow AS path for internal routes"
						comments    = "some foo"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_aspathlist.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_aspathlist.test", "name", "Internal"),
				),
			},
			{
				ResourceName:      "netboxbgp_aspathlist.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
