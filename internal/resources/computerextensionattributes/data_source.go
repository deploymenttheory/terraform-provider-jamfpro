// computerextensionattributes_data_source.go
package computerextensionattributes

import (
	"context"
	"fmt"

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
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the computer extension attribute.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique name of the Jamf Pro computer extension attribute.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the computer extension attribute.",
			},
			"data_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "String",
				Description: "Data type of the computer extension attribute. Can be String, Integer, or Date.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enabled by default, but for inputType Script we can disable it as well.Possible values are: false or true.",
			},
			"inventory_display_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Category in which to display the extension attribute in Jamf Pro.",
			},
			"input_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Extension attributes collect inventory data by using an input type.The type of the Input used to populate the extension attribute.",
			},
			"script_contents": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When we run this script it returns a data value each time a computer submits inventory to Jamf Pro. Provide scriptContents only when inputType is 'SCRIPT'.",
			},
			"popup_menu_choices": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "When added with list of choices while creating computer extension attributes these Pop-up menu can be displayed in inventory information. User can choose a value from the pop-up menu list when enrolling a computer any time using Jamf Pro. Provide popupMenuChoices only when inputType is 'POPUP'.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ldap_attribute_mapping": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Directory Service attribute use to populate the extension attribute.Required when inputType is 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING'.",
			},
			"ldap_extension_attribute_allowed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Collect multiple values for this extension attribute. ldapExtensionAttributeAllowed is disabled by default, only for inputType 'DIRECTORY_SERVICE_ATTRIBUTE_MAPPING' it can be enabled. It's value cannot be modified during edit operation.Possible values are:true or false.",
			},
		},
	}
}

// dataSourceRead fetches the details of a specific computer extension attribute
// computerextensionattributes_data_source.go
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	// Get the id and name from schema
	resourceID := d.Get("id").(string)
	name := d.Get("name").(string)

	// Validate that at least one identifier is provided
	if resourceID == "" && name == "" {
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	var resource *jamfpro.ResourceComputerExtensionAttribute
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error

		// Try to get by name first if provided
		if name != "" {
			resource, apiErr = client.GetComputerExtensionAttributeByName(name)
		} else {
			// Fall back to ID lookup
			resource, apiErr = client.GetComputerExtensionAttributeByID(resourceID)
		}

		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		lookupMethod := "ID"
		lookupValue := resourceID
		if name != "" {
			lookupMethod = "name"
			lookupValue = name
		}
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Computer Extension Attribute with %s '%s' after retries: %v", lookupMethod, lookupValue, err))
	}

	if resource == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("the Jamf Pro Computer Extension Attribute was not found"))
	}

	d.SetId(resource.ID)
	return updateState(d, resource)
}
