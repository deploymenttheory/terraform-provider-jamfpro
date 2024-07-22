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
		ReadContext: dataSourceRead,
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

// dataSourceRead fetches the details of a specific department from Jamf Pro using its unique ID.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resourceID := d.Get("id").(string)

	department, err := client.GetDepartmentByID(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Department with ID '%s': %v", resourceID, err))
	}

	if department != nil {
		d.SetId(resourceID)
		if err := d.Set("name", department.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Department with ID '%s': %v", resourceID, err))...)
		}

	} else {
		d.SetId("")
	}

	return diags
}
