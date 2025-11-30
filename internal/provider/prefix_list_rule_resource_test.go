package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPrefixListRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        testExternalProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "netboxbgp_prefix_list" "test" {
						name        = "%[1]s"
						family      = "ipv4"
						description = "RFC 1918 prefixes"
						comments    = "do not announce publicly"
						tags        = []
					}

					resource "netbox_prefix" "lan" {
						prefix = "%[2]s/32"
						status = "active"
					}

					resource "netboxbgp_prefix_list_rule" "lan" {
						action      = "permit"
						index       = 1
						prefix      = netbox_prefix.lan.id
						prefix_list = netboxbgp_prefix_list.test.id
					}

					resource "netboxbgp_prefix_list_rule" "custom" {
						action        = "permit"
						index         = 2
						prefix_custom = "10.0.0.0/8"
						prefix_list   = netboxbgp_prefix_list.test.id
					}
				`, testName(t), testIP(t, 0)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_prefix_list_rule.lan", "id"),
					resource.TestCheckResourceAttrSet("netboxbgp_prefix_list_rule.custom", "id"),
				),
			},
			{
				ResourceName:      "netboxbgp_prefix_list_rule.lan",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
