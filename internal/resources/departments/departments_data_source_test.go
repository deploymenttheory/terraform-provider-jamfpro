package departments

import (
	"fmt"
	"testing"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/acceptance/check"
)

type JamfProDepartmentsDataSource struct{}

func TestAccJamfProDepartmentsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.provider_jamfpro_department", "test")
	r := JamfProDepartmentsDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("test-dept-%d", data.RandomInteger)),
			),
		},
	})
}

func (JamfProDepartmentsDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "provider_jamfpro_department" "test" {
  name = "test-dept-%d"
}

data "provider_jamfpro_department" "test" {
  name = provider_jamfpro_department.test.name
}
`, data.RandomInteger)
}
