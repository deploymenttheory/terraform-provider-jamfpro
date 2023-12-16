// department_data_source.go
package departments

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProDepartments provides information about a specific department in Jamf Pro.
func DataSourceJamfProDepartments() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProDepartmentsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
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

// DataSourceJamfProDepartmentsRead fetches the details of a specific department from Jamf Pro using either its unique Name or its Id.
func DataSourceJamfProDepartmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var department *jamfpro.ResponseDepartment
	var err error

	// Check if Name is provided in the data source configuration
	if v, ok := d.GetOk("name"); ok {
		departmentName, ok := v.(string)
		if !ok {
			return diag.Errorf("error asserting 'name' as string")
		}
		if departmentName != "" {
			department, err = conn.GetDepartmentByName(departmentName)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to fetch department by name: %v", err))
			}
		}
	} else if v, ok := d.GetOk("id"); ok {
		departmentIDStr, ok := v.(string)
		if !ok {
			return diag.Errorf("error asserting 'id' as string")
		}
		departmentID, convertErr := strconv.Atoi(departmentIDStr)
		if convertErr != nil {
			return diag.FromErr(fmt.Errorf("failed to convert department ID to integer: %v", convertErr))
		}
		department, err = conn.GetDepartmentByID(departmentID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch department by ID: %v", err))
		}
	} else {
		return diag.Errorf("Either 'name' or 'id' must be provided")
	}

	if department == nil {
		return diag.FromErr(fmt.Errorf("department not found"))
	}

	// Set the data source attributes using the fetched data
	if err := d.Set("id", department.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting 'id': %v", err))
	}
	if err := d.Set("name", department.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting 'name': %v", err))
	}

	return nil
}
