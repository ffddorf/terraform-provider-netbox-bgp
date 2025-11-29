package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAspathlistResource(t *testing.T) {
	resourceName := testName(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "netboxbgp_aspathlist" "test" {
						name        = "%[1]s"
						description = "Allow AS path for internal routes"
						comments    = "some foo"
					}
				`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netboxbgp_aspathlist.test", "id"),
					resource.TestCheckResourceAttr("netboxbgp_aspathlist.test", "name", resourceName),
				),
			},
			{
				ResourceName:      "netboxbgp_aspathlist.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
