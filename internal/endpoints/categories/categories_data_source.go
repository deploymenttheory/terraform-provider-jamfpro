// categories_data_source.go
package categories

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProCategories provides information about a specific Category in Jamf Pro.
func DataSourceJamfProCategories() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProCategoriesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
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
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)

	Category, err := client.GetCategoryByID(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Category with ID '%s': %v", resourceID, err))
	}

	if Category != nil {
		d.SetId(resourceID)
		if err := d.Set("name", Category.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Category with ID '%s': %v", resourceID, err))...)
		}

	} else {
		d.SetId("")
	}

	return diags
}
