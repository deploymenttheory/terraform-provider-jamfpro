// computerextensionattributes_data_source.go
package computerextensionattributes

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProComputerExtensionAttributes provides information about a specific computer extension attribute by its ID or Name.
func DataSourceJamfProComputerExtensionAttributes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJamfProComputerExtensionAttributesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The unique identifier of the computer extension attribute.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique name of the Jamf Pro computer extension attribute.",
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Indicates if the computer extension attribute is enabled.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the computer extension attribute.",
				Computed:    true,
			},
			"data_type": {
				Type:        schema.TypeString,
				Description: "Data type of the computer extension attribute. Can be String / Integer / Date (YYYY-MM-DD hh:mm:ss)",
				Computed:    true,
			},
			"input_type": {
				Type:        schema.TypeList,
				Description: "Input type details of the computer extension attribute.",
				Computed:    true,
				//MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Description: "Type of the input for the computer extension attribute.",
							Computed:    true,
						},
						"platform": {
							Type:        schema.TypeString,
							Description: "Platform type for the computer extension attribute.",
							Computed:    true,
						},
						"script": {
							Type:        schema.TypeString,
							Description: "Script associated with the computer extension attribute.",
							Computed:    true,
						},
						"choices": {
							Type:        schema.TypeList,
							Description: "Choices associated with the computer extension attribute if it is a pop-up menu type.",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"inventory_display": {
				Type:        schema.TypeString,
				Description: "Display details for inventory for the computer extension attribute.",
				Computed:    true,
			},
			"recon_display": {
				Type:        schema.TypeString,
				Description: "Display details for recon for the computer extension attribute.",
				Computed:    true,
			},
		},
	}
}

// dataSourceJamfProComputerExtensionAttributesRead fetches the details of a specific computer extension attribute
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
func dataSourceJamfProComputerExtensionAttributesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var attribute *jamfpro.ResponseComputerExtensionAttribute
	var err error

	// Check if Name is provided in the data source configuration
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		attributeName := v.(string)
		attribute, err = conn.GetComputerExtensionAttributeByName(attributeName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch computer extension attribute by name: %v", err))
		}
	} else if v, ok := d.GetOk("id"); ok {
		attributeID := v.(int) // Correctly cast to int
		attribute, err = conn.GetComputerExtensionAttributeByID(attributeID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch computer extension attribute by ID: %v", err))
		}
	} else {
		return diag.Errorf("Either 'name' or 'id' must be provided")
	}

	// Set the data source attributes using the fetched data
	if attribute == nil {
		return diag.FromErr(fmt.Errorf("computer extension attribute not found"))
	}

	// Set the data source attributes using the fetched data
	if err := d.Set("name", attribute.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'name': %v", err))
	}
	if err := d.Set("enabled", attribute.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'enabled': %v", err))
	}
	if err := d.Set("description", attribute.Description); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'description': %v", err))
	}
	if err := d.Set("data_type", attribute.DataType); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'data_type': %v", err))
	}
	if err := d.Set("inventory_display", attribute.InventoryDisplay); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'inventory_display': %v", err))
	}
	if err := d.Set("recon_display", attribute.ReconDisplay); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'recon_display': %v", err))
	}

	// Extract the input type details and set them in the data source
	inputType := make(map[string]interface{})
	inputType["type"] = attribute.InputType.Type
	inputType["platform"] = attribute.InputType.Platform
	inputType["script"] = attribute.InputType.Script
	inputType["choices"] = attribute.InputType.Choices
	if err := d.Set("input_type", []interface{}{inputType}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set 'input_type': %v", err))
	}

	d.SetId(fmt.Sprintf("%d", attribute.ID))

	return nil
}
