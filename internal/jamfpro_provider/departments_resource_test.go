package jamfpro_provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDepartmentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDepartmentResourceConfig("Engineering"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jamfpro_department.test", "name", "Engineering"),
					resource.TestCheckResourceAttrSet("jamfpro_department.test", "id"),
					resource.TestCheckResourceAttrSet("jamfpro_department.test", "href"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "jamfpro_department.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDepartmentResourceConfig("HR"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("jamfpro_department.test", "name", "HR"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDepartmentResourceConfig(departmentName string) string {
	return fmt.Sprintf(`
resource "jamfpro_department" "test" {
  name = %[1]q
}
`, departmentName)
}
