package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSessionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        testExternalProviders,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: fmt.Sprintf(`%s
					resource "netboxbgp_session" "test" {
						name           = "My session"
						status         = "active"
						device         = netbox_device.test.id
						local_address  = netbox_ip_address.local.id
						remote_address = netbox_ip_address.remote.id
						local_as       = netbox_asn.test.id
						remote_as      = netbox_asn.test.id
					}

					data "netboxbgp_session" "test" {
						depends_on = [netboxbgp_session.test]
						id = netboxbgp_session.test.id
					}
				`, baseResources(t)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netboxbgp_session.test", "name", "My session"),
					resource.TestCheckResourceAttrPair("data.netboxbgp_session.test", "device.name", "netbox_device.test", "name"),
				),
			},
		},
	})
}
