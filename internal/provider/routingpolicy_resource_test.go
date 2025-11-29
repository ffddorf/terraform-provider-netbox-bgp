package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoutingpolicyResource(t *testing.T) {
	resourceName := testName(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "netboxbgp_routingpolicy" "test" {
						name        = "%[1]s"
						description = "Policy towards public peers"
					}
				`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_routingpolicy.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_routingpolicy.test", "name", resourceName),
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
