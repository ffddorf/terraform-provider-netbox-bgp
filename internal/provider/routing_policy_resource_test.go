package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoutingPolicyResource(t *testing.T) {
	resourceName := testName(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "netboxbgp_routing_policy" "test" {
						name        = "%[1]s"
						description = "Policy towards public peers"
					}
				`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_routing_policy.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_routing_policy.test", "name", resourceName),
				),
			},
			{
				ResourceName:      "netboxbgp_routing_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
