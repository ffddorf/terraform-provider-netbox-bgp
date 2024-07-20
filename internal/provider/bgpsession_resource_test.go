package provider

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func testName(t *testing.T) string {
	return t.Name() + "_" + uuid.NewString()
}

func baseResources(t *testing.T) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = "%[1]s"
}

resource "netbox_site" "test" {
  name = "%[1]s"
  status = "active"
}

resource "netbox_device_role" "test" {
  name = "%[1]s"
  color_hex = "123456"
}

resource "netbox_manufacturer" "test" {
  name = "%[1]s"
}

resource "netbox_device_type" "test" {
  model = "%[1]s"
  manufacturer_id = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name = "%[1]s"
  device_type_id = netbox_device_type.test.id
  role_id = netbox_device_role.test.id
  site_id = netbox_site.test.id
}

resource "netbox_device_interface" "test" {
  name      = "eth0"
  device_id = netbox_device.test.id
  type      = "1000base-t"
}

resource "netbox_ip_address" "local" {
  ip_address   = "203.0.113.10/24"
  status       = "active"
  interface_id = netbox_device_interface.test.id
  object_type  = "dcim.interface"
}

resource "netbox_ip_address" "remote" {
  ip_address   = "203.0.113.11/24"
  status       = "active"
}

resource "netbox_rir" "test" {
  name = "%[1]s"
}

resource "netbox_asn" "test" {
  asn    = 1337
	rir_id = netbox_rir.test.id
}`, testName(t))
}

func TestAccSessionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"netbox": {
				VersionConstraint: "~> 3.8.7",
				Source:            "registry.terraform.io/e-breuninger/netbox",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`%s

					resource "netboxbgp_session" "test" {
						name              = "My session"
						status            = "active"
						device_id         = netbox_device.test.id
						local_address_id  = netbox_ip_address.local.id
						remote_address_id = netbox_ip_address.remote.id
						local_as_id       = netbox_asn.test.id
						remote_as_id      = netbox_asn.test.id
					}
				`, baseResources(t)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_session.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_session.test", "name", "My session"),
				),
			},
			{
				ResourceName:      "netboxbgp_session.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(`%s
					resource "netbox_tag" "test_a" {
						name = "Integration Test"
					}
					resource "netbox_tag" "test_b" {
						name = "temporary"
					}

					resource "netboxbgp_session" "test" {
						name              = "My session changed"
						status            = "active"
						device_id         = netbox_device.test.id
						local_address_id  = netbox_ip_address.local.id
						remote_address_id = netbox_ip_address.remote.id
						local_as_id       = netbox_asn.test.id
						remote_as_id      = netbox_asn.test.id
						site_id           = netbox_site.test.id
						tags              = ["Integration Test", "temporary"]
					}
				`, baseResources(t)),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"netboxbgp_session.test",
						tfjsonpath.New("tags"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("Integration Test"),
							knownvalue.StringExact("temporary"),
						}),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netboxbgp_session.test", "name", "My session changed"),
					resource.TestCheckResourceAttrPair("netboxbgp_session.test", "site_id", "netbox_site.test", "id"),
				),
			},
		},
	})
}
