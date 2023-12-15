// scripts_date_source.go
package scripts

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProScripts provides information about a specific Jamf Pro script by its ID or Name.
func DataSourceJamfProScripts() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProScriptsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the script.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the script.",
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Category of the script.",
			},
			"filename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filename of the script.",
			},
			"info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Information to display to the administrator when the script is run.",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Notes to display about the script (e.g., who created it and when it was created).",
			},
			"priority": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Execution priority of the script (Before, After, At Reboot).",
			},
			"parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Script parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter4": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter5": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter6": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter7": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter8": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter9": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter10": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter11": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"os_requirements": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OS requirements for the script.",
			},
			"script_contents": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contents of the script.",
			},
			"script_contents_encoded": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Encoded contents of the script.",
			},
		},
	}
}

// dataSourceJamfProScriptsRead fetches the details of a specific Jamf Pro script
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
func DataSourceJamfProScriptsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var attribute *jamfpro.ResourceScript

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		attributeID := d.Id()

		// Try fetching the script using the ID
		var apiErr error
		attribute, apiErr = conn.GetScriptByID(attributeID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			attributeName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string for read function"))
			}

			attribute, apiErr = conn.GetScriptByName(attributeName)
			if apiErr != nil {
				// Handle the APIError
				if apiError, ok := apiErr.(*http_client.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while reading the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	// Construct a map of script attributes
	scriptAttributes := map[string]interface{}{
		"name":            attribute.Name,
		"category_name":   attribute.CategoryName,
		"category_id":     attribute.CategoryId,
		"info":            attribute.Info,
		"notes":           attribute.Notes,
		"os_requirements": attribute.OSRequirements,
		"priority":        attribute.Priority,
		"script_contents": encodeScriptContent(attribute.ScriptContents),
		"parameter4":      attribute.Parameter4,
		"parameter5":      attribute.Parameter5,
		"parameter6":      attribute.Parameter6,
		"parameter7":      attribute.Parameter7,
		"parameter8":      attribute.Parameter8,
		"parameter9":      attribute.Parameter9,
		"parameter10":     attribute.Parameter10,
		"parameter11":     attribute.Parameter11,
	}

	// Update the Terraform state with script attributes
	for key, value := range scriptAttributes {
		if err := d.Set(key, value); err != nil {
			diags = append(diags, diag.FromErr(err)...)
			return diags
		}
	}

	return diags
}
