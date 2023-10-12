/*
 */
package jamfpro_provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDepartmentDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create a department
			{
				Config: testAccDepartmentResourceConfig("initial_name"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jamfpro_department.test", "name", "initial_name"),
				),
			},
			// Step 2: Update the department
			{
				Config: testAccDepartmentDataSourceConfig("updated_name"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jamfpro_department.test", "name", "updated_name"),
				),
			},
			// Step 3: Import the department and check if it's still available
			{
				ResourceName:      "jamfpro_department.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDepartmentDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	provider "jamfpro" {
		client_id     = "%s"
		client_secret = "%s"
		instance_name = "%s"
		debug_mode    = %s
	}

	resource "jamfpro_department" "tf_department_test" {
		name = "%s"
	}`, os.Getenv("JAMFPRO_CLIENT_ID"),
		os.Getenv("JAMFPRO_CLIENT_SECRET"),
		os.Getenv("JAMFPRO_INSTANCE_NAME"),
		os.Getenv("JAMFPRO_DEBUG_MODE"),
		name,
	)
}
