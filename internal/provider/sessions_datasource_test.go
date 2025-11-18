package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func testSessions(t *testing.T) string {
	return fmt.Sprintf(`%s
	resource "netboxbgp_session" "test1" {
		name              = "Session 1"
		status            = "active"
		device         = netbox_device.test.id
		local_address  = netbox_ip_address.local.id
		remote_address = netbox_ip_address.remote.id
		local_as       = netbox_asn.test.id
		remote_as      = netbox_asn.test.id
	}

	resource "netbox_ip_address" "remote2" {
		ip_address   = "%s/24"
		status       = "active"
	}

	resource "netboxbgp_session" "test2" {
		name              = "Session 2"
		status            = "planned"
		device         = netbox_device.test.id
		local_address  = netbox_ip_address.local.id
		remote_address = netbox_ip_address.remote2.id
		local_as       = netbox_asn.test.id
		remote_as      = netbox_asn.test.id
	}

	resource "netbox_ip_address" "remote3" {
		ip_address   = "%s/24"
		status       = "active"
	}

	resource "netboxbgp_session" "test3" {
		name              = "Session 3"
		status            = "active"
		device         = netbox_device.test.id
		local_address  = netbox_ip_address.local.id
		remote_address = netbox_ip_address.remote3.id
		local_as       = netbox_asn.test.id
		remote_as      = netbox_asn.test.id
	}`, baseResources(t), testIP(t, 2), testIP(t, 3))
}

func TestAccSessionsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        testExternalProviders,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: fmt.Sprintf(`%s
					data "netboxbgp_sessions" "test_active" {
						depends_on = [
							netboxbgp_session.test1,
							netboxbgp_session.test2,
							netboxbgp_session.test3,
						]

						filters = [
							{ name: "status", value: "active" }
						]

						ordering = "name"
					}

					data "netboxbgp_sessions" "test_limit" {
						depends_on = [
							netboxbgp_session.test1,
							netboxbgp_session.test2,
							netboxbgp_session.test3,
						]

						limit    = 2
						ordering = "name"
					}
				`, testSessions(t)),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.netboxbgp_sessions.test_active",
						tfjsonpath.New("sessions"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("Session 1"),
							}),
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("Session 3"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.netboxbgp_sessions.test_limit",
						tfjsonpath.New("sessions"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("Session 1"),
							}),
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("Session 2"),
							}),
						}),
					),
				},
			},
		},
	})
}
