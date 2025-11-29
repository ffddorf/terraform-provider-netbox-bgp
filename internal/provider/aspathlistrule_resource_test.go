package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAspathlistruleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "netboxbgp_aspathlist" "test" {
						name        = "IXP"
						description = "Rules for ASPaths on interchanges"
						comments    = "some foo"
					}

					resource "netboxbgp_aspathlistrule" "test" {
						aspath_list = netboxbgp_aspathlist.test.id

						action = "permit"
						index  = 0
						pattern = "_65002_"

						description = "some text here"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_aspathlistrule.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_aspathlistrule.test", "action", "permit"),
					resource.TestCheckResourceAttr("netboxbgp_aspathlistrule.test", "index", "0"),
					resource.TestCheckResourceAttr("netboxbgp_aspathlistrule.test", "description", "some text here"),
				),
			},
			{
				ResourceName:      "netboxbgp_aspathlistrule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
