package departments_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/acceptance/check"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type JamfProDepartmentsResource struct{}

func TestAccJamfProDepartments_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "provider_jamfpro_department", "test")
	r := JamfProDepartmentsResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInJamfPro(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("test-dept-%d", data.RandomInteger)),
			),
		},
	})
}

func TestAccJamfProDepartments_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "provider_jamfpro_department", "test")
	r := JamfProDepartmentsResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInJamfPro(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("complete-dept-%d", data.RandomInteger)),
			),
		},
	})
}

func TestAccJamfProDepartments_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "provider_jamfpro_department", "test")
	r := JamfProDepartmentsResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInJamfPro(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("test-dept-%d", data.RandomInteger)),
			),
		},
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInJamfPro(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("complete-dept-%d", data.RandomInteger)),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInJamfPro(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("test-dept-%d", data.RandomInteger)),
			),
		},
	})
}

func TestAccJamfProDepartments_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "provider_jamfpro_department", "test")
	r := JamfProDepartmentsResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInJamfPro(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("provider_jamfpro_department"),
		},
	})
}

func (r JamfProDepartmentsResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	// Use the clients to check if the department exists.
	exists := false

	// Try fetching the department by its ID.
	departmentID, err := strconv.Atoi(state.ID)
	if err == nil {
		department, err := clients.conn.GetDepartmentByID(departmentID)
		if err != nil {
			return nil, fmt.Errorf("Error checking existence of department with ID %d: %v", departmentID, err)
		}
		if department != nil && department.Id == departmentID {
			exists = true
			return &exists, nil
		}
	}

	// If fetching by ID failed, try fetching by its name.
	// Note: This assumes that the name of the department is stored in the state.
	// If not, you might need to adjust this part.
	departmentName, ok := state.Attributes["name"]
	if !ok {
		return nil, fmt.Errorf("Error retrieving department name from state")
	}

	department, err := clients.conn.GetDepartmentByName(departmentName)
	if err != nil {
		return nil, fmt.Errorf("Error checking existence of department with name %s: %v", departmentName, err)
	}

	if department != nil && department.Name == departmentName {
		exists = true
	}

	return &exists, nil
}

func (r JamfProDepartmentsResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "provider_jamfpro_department" "test" {
  name = "test-dept-%d"
}
`, data.RandomInteger)
}

func (r JamfProDepartmentsResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "provider_jamfpro_department" "test" {
  name = "complete-dept-%d"
}
`, data.RandomInteger)
}

func (r JamfProDepartmentsResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "provider_jamfpro_department" "test" {
  name = "updated-test-dept-%d"
}
`, data.RandomInteger)
}

func (r JamfProDepartmentsResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "provider_jamfpro_department" "import" {
  name = provider_jamfpro_department.test.name
}
`, r.basic(data))
}
