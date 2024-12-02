// computerextensionattributes_data_source.go
package computerextensionattributes

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProComputerExtensionAttributes provides information about a specific computer extension attribute by its ID or Name.
func DataSourceJamfProComputerExtensionAttributes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the computer extension attribute.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the Jamf Pro computer extension attribute.",
			},
		},
	}
}

// dataSourceRead fetches the details of a specific computer extension attribute
// from Jamf Pro using either its unique Name or its Id. The function prioritizes the 'name' attribute over the 'id'
// attribute for fetching details. If neither 'name' nor 'id' is provided, it returns an error.
// Once the details are fetched, they are set in the data source's state.
//
// Parameters:
// - ctx: The context within which the function is called. It's used for timeouts and cancellation.
// - d: The current state of the data source.
// - meta: The meta object that can be used to retrieve the API client connection.
//
// Returns:
// - diag.Diagnostics: Returns any diagnostics (errors or warnings) encountered during the function's execution.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)

	var resource *jamfpro.ResourceComputerExtensionAttribute
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = client.GetComputerExtensionAttributeByID(resourceID)
		if apiErr != nil {

			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Computer Extension Attribute with ID '%s' after retries: %v", resourceID, err))
	}

	if resource != nil {
		d.SetId(resourceID)
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Computer Extension Attribute with ID '%s': %v", resourceID, err))...)
		}
	} else {
		d.SetId("")
	}

	return diags
}

// DataSourceJamfProComputerExtensionAttributesList provides a list of all Jamf Pro computer extension attributes.
func DataSourceJamfProComputerExtensionAttributesList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadList,
		Schema: map[string]*schema.Schema{
			"attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// dataSourceReadList fetches a list of all Jamf Pro computer extension attributes
// and maps them into the Terraform state.
func dataSourceReadList(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	var diags diag.Diagnostics

	// Fetch the list of computer extension attributes
	response, err := client.GetComputerExtensionAttributes("")
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve Jamf Pro computer extension attributes: %v", err))
	}

	// Map the attributes to the Terraform state
	var attributes []map[string]interface{}
	var ids []string

	for _, attr := range response.Results {
		attributes = append(attributes, map[string]interface{}{
			"id":          attr.ID,
			"name":        attr.Name,
			"description": attr.Description,
		})
		ids = append(ids, attr.ID)
	}

	// Set the computed attributes in Terraform state
	if err := d.Set("attributes", attributes); err != nil {
		return diag.FromErr(fmt.Errorf("error setting 'attributes' attribute: %v", err))
	}
	if err := d.Set("ids", ids); err != nil {
		return diag.FromErr(fmt.Errorf("error setting 'ids' attribute: %v", err))
	}

	// Generate a unique ID for the resource
	d.SetId(fmt.Sprintf("datasource-computer-extension-attributes-list-%d", time.Now().Unix()))

	return diags
}
