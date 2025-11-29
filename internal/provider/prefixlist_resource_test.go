package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPrefixlistResource(t *testing.T) {
	resourceName := testName(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "netboxbgp_prefixlist" "test" {
						name        = "%[1]s"
						family      = "ipv4"
						description = "Prefix belonging to that neighbor"
						comments    = "on some IX"
						tags        = []
					}
				`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_prefixlist.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_prefixlist.test", "name", resourceName),
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
