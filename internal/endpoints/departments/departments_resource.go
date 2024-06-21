// department_resource.go
package departments

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProDepartments defines the schema and CRUD operations for managing Jamf Pro Departments in Terraform.
func ResourceJamfProDepartments() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProDepartmentsCreate,
		ReadContext:   resourceJamfProDepartmentsRead,
		UpdateContext: resourceJamfProDepartmentsUpdate,
		DeleteContext: resourceJamfProDepartmentsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the department.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro department.",
			},
		},
	}
}
