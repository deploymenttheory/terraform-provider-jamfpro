// computerprestages_data_source.go
package computerprestages

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProComputerPrestage provides information about a specific department in Jamf Pro.
func DataSourceJamfProComputerPrestage() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProComputerPrestageRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the computer prestage.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the computer prestage.",
			},
		},
	}
}

// DataSourceJamfProComputerPrestageRead fetches the details of a specific department from Jamf Pro using its unique ID.
func DataSourceJamfProComputerPrestageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Get the department ID from the data source's arguments
	resourceID := d.Get("id").(string)

	// Attempt to fetch the department's details using its ID
	department, err := conn.GetJamfApiRoleByID(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Department with ID '%s': %v", resourceID, err))
	}

	// Check if resource data exists and set the Terraform state
	if department != nil {
		d.SetId(resourceID) // Set the id in the Terraform state
		if err := d.Set("display_name", department.DisplayName); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'display_name' for Jamf Pro Computer Prestage with ID '%s': %v", resourceID, err))...)
		}

	} else {
		d.SetId("") // Data not found, unset the id in the Terraform state
	}

	return diags
}
