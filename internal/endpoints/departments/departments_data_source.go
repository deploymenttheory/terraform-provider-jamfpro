// department_data_source.go
package departments

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProDepartments provides information about a specific department in Jamf Pro.
func DataSourceJamfProDepartments() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProDepartmentsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the department.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the jamf pro department.",
			},
		},
	}
}

// DataSourceJamfProDepartmentsRead fetches the details of a specific department from Jamf Pro using its unique ID.
func DataSourceJamfProDepartmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	// Initialize variables
	var diags diag.Diagnostics

	// Get the department ID from the data source's arguments
	resourceID := d.Get("id").(string)

	// Attempt to fetch the department's details using its ID
	department, err := client.GetDepartmentByID(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Department with ID '%s': %v", resourceID, err))
	}

	// Check if resource data exists and set the Terraform state
	if department != nil {
		d.SetId(resourceID) // Set the id in the Terraform state
		if err := d.Set("name", department.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Department with ID '%s': %v", resourceID, err))...)
		}

	} else {
		d.SetId("") // Data not found, unset the id in the Terraform state
	}

	return diags
}
