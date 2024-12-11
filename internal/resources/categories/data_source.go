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
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The unique identifier of the Category.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The unique name of the jamf pro Category.",
			},
		},
	}
}

// dataSourceRead fetches the details of a specific category from Jamf Pro using its unique ID.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)
	resourceName := d.Get("name").(string)

	var resource *jamfpro.ResourceCategory
	var err error

	if resourceID != "" {
		resource, err = client.GetCategoryByID(resourceID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Category with ID '%s': %v", resourceID, err))
		}
	} else if resourceName != "" {
		resource, err = client.GetCategoryByName(resourceName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Category with Name '%s': %v", resourceName, err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	if resource != nil {
		// Set the ID from the Category object
		d.SetId(resource.Id)
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Category with ID '%s': %v", resource.Id, err))...)
		}
	} else {
		d.SetId("")
	}

	return diags
}
