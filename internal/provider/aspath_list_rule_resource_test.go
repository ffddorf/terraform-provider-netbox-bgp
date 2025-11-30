package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAspathListRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "netboxbgp_aspath_list" "test" {
						name        = "%[1]s"
						description = "Rules for ASPaths on interchanges"
						comments    = "some foo"
					}

					resource "netboxbgp_aspath_list_rule" "test" {
						aspath_list = netboxbgp_aspath_list.test.id

						action = "permit"
						index  = 0
						pattern = "_65002_"

						description = "some text here"
					}
				`, testName(t)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_aspath_list_rule.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_aspath_list_rule.test", "action", "permit"),
					resource.TestCheckResourceAttr("netboxbgp_aspath_list_rule.test", "index", "0"),
					resource.TestCheckResourceAttr("netboxbgp_aspath_list_rule.test", "description", "some text here"),
				),
			},
			{
				ResourceName:      "netboxbgp_aspath_list_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
