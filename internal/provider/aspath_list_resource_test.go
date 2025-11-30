package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAspathListResource(t *testing.T) {
	resourceName := testName(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "netboxbgp_aspath_list" "test" {
						name        = "%[1]s"
						description = "Allow AS path for internal routes"
						comments    = "some foo"
					}
				`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_aspath_list.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_aspath_list.test", "name", resourceName),
				),
			},
			{
				ResourceName:      "netboxbgp_aspath_list.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
