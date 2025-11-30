package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPrefixListResource(t *testing.T) {
	resourceName := testName(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "netboxbgp_prefix_list" "test" {
						name        = "%[1]s"
						family      = "ipv4"
						description = "Prefix belonging to that neighbor"
						comments    = "on some IX"
						tags        = []
					}
				`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_prefix_list.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_prefix_list.test", "name", resourceName),
				),
			},
			{
				ResourceName:      "netboxbgp_prefix_list.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
