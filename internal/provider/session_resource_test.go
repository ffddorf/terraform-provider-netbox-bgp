package provider

import (
	"fmt"
	"hash/fnv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var testExternalProviders = map[string]resource.ExternalProvider{
	"netbox": {
		VersionConstraint: "~> 3.8.7",
		Source:            "registry.terraform.io/e-breuninger/netbox",
	},
}

func testName(t *testing.T) string {
	return t.Name() + "_" + uuid.NewString()
}

func testNum(t *testing.T) uint64 {
	h := fnv.New64()
	fmt.Fprint(h, testName(t))
	return h.Sum64()
}

func testIP(t *testing.T, offset uint64) string {
	num := testNum(t)
	shortNum := num % 250
	return fmt.Sprintf("203.0.113.%d", shortNum+offset)
}

func baseResources(t *testing.T) string {
	num := testNum(t)
	shortNum := num % 250
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
  ip_address   = "%[2]s/24"
  status       = "active"
  interface_id = netbox_device_interface.test.id
  object_type  = "dcim.interface"
}

resource "netbox_ip_address" "remote" {
  ip_address   = "%[3]s/24"
  status       = "active"
}

resource "netbox_rir" "test" {
  name = "%[1]s"
}

resource "netbox_asn" "test" {
  asn    = %[4]d
	rir_id = netbox_rir.test.id
}`, testName(t), testIP(t, 1), testIP(t, 2), shortNum+1337)
}

func TestAccSessionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        testExternalProviders,
		Steps: []resource.TestStep{
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
						device         = netbox_device.test.id
						local_address  = netbox_ip_address.local.id
						remote_address = netbox_ip_address.remote.id
						local_as       = netbox_asn.test.id
						remote_as      = netbox_asn.test.id
						site           = netbox_site.test.id
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
					resource.TestCheckResourceAttrPair("netboxbgp_session.test", "site", "netbox_site.test", "id"),
				),
			},
		},
	})
}
