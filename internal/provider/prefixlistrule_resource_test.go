package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPrefixlistruleResource(t *testing.T) {
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

					resource "netbox_prefix" "lan" {
						prefix = "192.168.0.0/16"
						status = "active"
					}

					resource "netboxbgp_prefixlistrule" "lan" {
						action      = "permit"
						index       = 1
						prefix      = netbox_prefix.lan.id
						prefix_list = netboxbgp_prefixlist.test.id
					}

					resource "netboxbgp_prefixlistrule" "custom" {
						action        = "permit"
						index         = 2
						prefix_custom = "10.0.0.0/8"
						prefix_list   = netboxbgp_prefixlist.test.id
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_prefixlistrule.lan", "id"),
					resource.TestCheckResourceAttrSet("netboxbgp_prefixlistrule.custom", "id"),
				),
			},
			{
				ResourceName:      "netboxbgp_prefixlistrule.lan",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
