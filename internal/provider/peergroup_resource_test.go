package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPeergroupResource(t *testing.T) {
	resourceName := testName(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "netboxbgp_peergroup" "test" {
						name        = "%[1]s"
						description = "Someone we peer with"
						comments    = "This is one some IX"
					}
				`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_peergroup.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_peergroup.test", "name", resourceName),
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
