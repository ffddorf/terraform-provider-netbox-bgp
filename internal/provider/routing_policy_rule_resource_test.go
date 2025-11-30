package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoutingpolicyruleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        testExternalProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "netboxbgp_routing_policy" "test" {
						name        = "%[1]s"
						description = "for some peer"
					}

					resource "netbox_prefix" "test" {
						prefix = "%[2]s/128"
						status = "active"
					}

					resource "netboxbgp_prefix_list" "test" {
						name        = "%[1]s"
						family      = "ipv6"
					}

					resource "netboxbgp_prefix_list_rule" "lan" {
						action      = "permit"
						index       = 1
						prefix      = netbox_prefix.test.id
						prefix_list = netboxbgp_prefix_list.test.id
					}

					resource "netboxbgp_routing_policy_rule" "test" {
						routing_policy = netboxbgp_routing_policy.test.id

						action = "permit"
						index  = 0

						match_ipv6_address = [netboxbgp_prefix_list.test.id]
						set_actions = jsonencode({
							local_pref = 100
						})
					}
				`, testName(t), testIP6(t, 0)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_routing_policy_rule.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_routing_policy_rule.test", "action", "permit"),
					resource.TestCheckResourceAttrPair("netboxbgp_routing_policy_rule.test", "routing_policy", "netboxbgp_routing_policy.test", "id"),
					resource.TestCheckResourceAttrPair("netboxbgp_routing_policy_rule.test", "match_ipv6_address.0", "netboxbgp_prefix_list.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_routing_policy_rule.test", "set_actions", `{"local_pref":100}`),
				),
			},
			{
				ResourceName:      "netboxbgp_routing_policy_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
