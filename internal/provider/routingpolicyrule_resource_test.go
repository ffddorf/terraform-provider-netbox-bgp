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
					resource "netboxbgp_routingpolicy" "test" {
						name        = "%[1]s"
						description = "for some peer"
					}

					resource "netbox_prefix" "test" {
						prefix = "%[2]s/128"
						status = "active"
					}

					resource "netboxbgp_prefixlist" "test" {
						name        = "%[1]s"
						family      = "ipv6"
					}

					resource "netboxbgp_prefixlistrule" "lan" {
						action      = "permit"
						index       = 1
						prefix      = netbox_prefix.test.id
						prefix_list = netboxbgp_prefixlist.test.id
					}

					resource "netboxbgp_routingpolicyrule" "test" {
						routing_policy = netboxbgp_routingpolicy.test.id

						action = "permit"
						index  = 0

						match_ipv6_address = [netboxbgp_prefixlist.test.id]
						set_actions = jsonencode({
							local_pref = 100
						})
					}
				`, testName(t), testIP6(t, 0)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_routingpolicyrule.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_routingpolicyrule.test", "action", "permit"),
					resource.TestCheckResourceAttrPair("netboxbgp_routingpolicyrule.test", "routing_policy", "netboxbgp_routingpolicy.test", "id"),
					resource.TestCheckResourceAttrPair("netboxbgp_routingpolicyrule.test", "match_ipv6_address.0", "netboxbgp_prefixlist.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_routingpolicyrule.test", "set_actions", `{"local_pref":100}`),
				),
			},
			{
				ResourceName:      "netboxbgp_routingpolicyrule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
