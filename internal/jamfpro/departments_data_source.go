package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DepartmentDataSource provides information about a specific department in Jamf Pro.
// It can fetch department details using either the department's unique Name or its Id.
// The Name attribute is prioritized for fetching if provided. Otherwise, the Id is used.
func dataSourceJamfProDepartments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJamfProDepartmentsRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The unique identifier of the department.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The unique name of the jamf pro department.",
			},
		},
	}
}

func dataSourceJamfProDepartmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*APIClient) // Ensure your meta object is of type *APIClient, which contains methods to interact with the Jamf Pro API.

	var department *jamfpro.Department // Ensure to use the correct namespace
	var err error

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		departmentName := v.(string)
		department, err = conn.conn.GetDepartmentByName(departmentName) // Using the method from the jamfpro package
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch department by name: %v", err))
		}
	} else if v, ok := d.GetOk("id"); ok && v.(string) != "" {
		departmentID, convertErr := strconv.Atoi(v.(string))
		if convertErr != nil {
			return diag.FromErr(fmt.Errorf("failed to convert department ID to integer: %v", convertErr))
		}
		department, err = conn.conn.GetDepartmentByID(departmentID) // Using the method from the jamfpro package
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch department by ID: %v", err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("either 'name' or 'id' must be specified"))
	}

	if department == nil {
		return diag.FromErr(fmt.Errorf("department not found"))
	}

	d.SetId(fmt.Sprintf("%d", department.Id))
	d.Set("name", department.Name)

	return nil
}
