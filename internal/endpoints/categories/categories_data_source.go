// categories_data_source.go
package categories

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProCategories provides information about a specific Category in Jamf Pro.
func DataSourceJamfProCategories() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProCategoriesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The unique identifier of the Category.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the jamf pro Category.",
			},
		},
	}
}

// DataSourceJamfProCategoriesRead fetches the details of a specific category from Jamf Pro using its unique ID.
func DataSourceJamfProCategoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Get the Category ID from the data source's arguments
	resourceID := d.Get("id").(string)

	// Attempt to fetch the Category's details using its ID
	Category, err := conn.GetCategoryByID(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Category with ID '%s': %v", resourceID, err))
	}

	// Check if resource data exists and set the Terraform state
	if Category != nil {
		d.SetId(resourceID) // Set the id in the Terraform state
		if err := d.Set("name", Category.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Category with ID '%s': %v", resourceID, err))...)
		}

	} else {
		d.SetId("") // Data not found, unset the id in the Terraform state
	}

	return diags
}
