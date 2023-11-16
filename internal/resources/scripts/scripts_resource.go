// scripts_resource.go
package scripts

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProScripts defines the schema and CRUD operations for managing Jamf Pro Scripts in Terraform.
func ResourceJamfProScripts() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProScriptsCreate,
		ReadContext:   ResourceJamfProScriptsRead,
		UpdateContext: ResourceJamfProScriptsUpdate,
		DeleteContext: ResourceJamfProScriptsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Jamf Pro unique identifier (ID) of the script.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name for the script.",
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Category to add the script to.",
			},
			"filename": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filename of the script.",
			},
			"info": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Information to display to the administrator when the script is run.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notes to display about the script (e.g., who created it and when it was created).",
			},
			"priority": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Execution priority of the script (Before, After, At Reboot).",
				ValidateFunc: validation.StringInSlice([]string{"Before", "After", "At Reboot"}, false),
			},
			"parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Labels to use for script parameters. Parameters 1 through 3 are predefined as mount point, computer name, and username",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter4": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameter5": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameter6": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameter7": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameter8": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameter9": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameter10": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameter11": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"os_requirements": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The script can only be run on computers with these operating system versions. Each version must be separated by a comma (e.g., 10.11, 15, 16.1).",
			},
			"script_contents": {
				Type:        schema.TypeString,
				Required:    true,
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

// constructScript constructs a ResponseScript object from the provided schema data.
// It captures attributes from the schema and maps them to the fields of the structs in
// the jamf pro sdk.
func constructScript(d *schema.ResourceData) *jamfpro.ResponseScript {
	// Handling Parameters
	var params jamfpro.Parameters
	if p, ok := d.GetOk("parameters"); ok {
		paramsList := p.([]interface{})
		if len(paramsList) > 0 && paramsList[0] != nil {
			paramMap := paramsList[0].(map[string]interface{})
			params = jamfpro.Parameters{
				Parameter4:  paramMap["parameter4"].(string),
				Parameter5:  paramMap["parameter5"].(string),
				Parameter6:  paramMap["parameter6"].(string),
				Parameter7:  paramMap["parameter7"].(string),
				Parameter8:  paramMap["parameter8"].(string),
				Parameter9:  paramMap["parameter9"].(string),
				Parameter10: paramMap["parameter10"].(string),
				Parameter11: paramMap["parameter11"].(string),
			}
		}
	}

	// Constructing the ResponseScript
	return &jamfpro.ResponseScript{
		Name:           d.Get("name").(string),
		Category:       d.Get("category").(string),
		Filename:       d.Get("filename").(string),
		Info:           d.Get("info").(string),
		Notes:          d.Get("notes").(string),
		Priority:       d.Get("priority").(string),
		Parameters:     params,
		OSRequirements: d.Get("os_requirements").(string),
		ScriptContents: d.Get("script_contents").(string),
		// ScriptContentsEncoded is not set here as it's a computed field
	}
}

// Helper function to generate diagnostics based on the error type
func generateTFDiagsFromHTTPError(err error, d *schema.ResourceData, action string) diag.Diagnostics {
	var diags diag.Diagnostics
	resourceName, exists := d.GetOk("name")
	if !exists {
		resourceName = "unknown"
	}

	// Handle the APIError in the diagnostic
	if apiErr, ok := err.(*http_client.APIError); ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   fmt.Sprintf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message),
		})
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   err.Error(),
		})
	}
	return diags
}

// ResourceJamfProScriptsCreate is responsible for creating a new Jamf Pro Script in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProScriptsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the create operation
	var createdAttribute *jamfpro.ResponseScript
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the script
		attribute := constructScript(d)

		// Check if the attribute is nil (indicating an issue with input_type)
		if attribute == nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the computer extension attribute due to missing or invalid input_type"))
		}

		// Directly call the API to create the resource
		createdAttribute, err = conn.CreateScriptByID(attribute)
		if err != nil {
			// Check if the error is an APIError
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			// For simplicity, we're considering all other errors as retryable
			return retry.RetryableError(err)
		}

		return nil
	})

	if err != nil {
		// If there's an error while creating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "create")
	}

	// Set the ID of the created resource in the Terraform state
	d.SetId(strconv.Itoa(createdAttribute.ID))

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProScriptsRead(ctx, d, meta)
		if len(readDiags) > 0 {
			// If readDiags is not empty, it means there's an error, so we retry
			return retry.RetryableError(fmt.Errorf("failed to read the created resource"))
		}
		return nil
	})

	if err != nil {
		// If there's an error while updating the state for the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update state for")
	}

	return diags
}

// ResourceJamfProScriptsRead is responsible for reading the current state of a Jamf Pro Script Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProScriptsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	var attribute *jamfpro.ResponseScript

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		attributeID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse attribute ID: %v", convertErr))
		}

		// Try fetching the script using the ID
		var apiErr error
		attribute, apiErr = conn.GetScriptsByID(attributeID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			attributeName := d.Get("name").(string)
			attribute, apiErr = conn.GetScriptsByName(attributeName)
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

	// Safely set attributes in the Terraform state
	if err := d.Set("name", attribute.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProScriptsUpdate is responsible for updating an existing Jamf Pro Script on the remote system.
func ResourceJamfProScriptsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the update operation
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the updated script
		script := constructScript(d)

		// Convert the ID from the Terraform state into an integer to be used for the API request
		scriptID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse script ID: %v", convertErr))
		}

		// Directly call the API to update the resource
		_, apiErr := conn.UpdateScriptByID(scriptID, script)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by name
			scriptName := d.Get("name").(string)
			_, apiErr = conn.UpdateScriptByName(scriptName, script)
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
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	// Use the retry function for the read operation to update the Terraform state
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProScriptsRead(ctx, d, meta)
		if len(readDiags) > 0 {
			return retry.RetryableError(fmt.Errorf("failed to update the Terraform state for the updated resource"))
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	return diags
}

// ResourceJamfProScriptsDelete is responsible for deleting a Jamf Pro script.
func ResourceJamfProScriptsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the delete operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		scriptID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse script ID: %v", convertErr))
		}

		// Directly call the API to delete the resource
		apiErr := conn.DeleteScriptByID(scriptID)
		if apiErr != nil {
			// If the delete by ID fails, try deleting by name
			scriptName := d.Get("name").(string)
			apiErr = conn.DeleteScriptByName(scriptName)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while deleting the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "delete")
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
