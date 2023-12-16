// scripts_resource.go
package scripts

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"

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
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
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
			"category_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the category to add the script to.",
			},
			"category_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Jamf Pro unique identifier (ID) of the category.",
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
			"os_requirements": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The script can only be run on computers with these operating system versions. Each version must be separated by a comma (e.g., 10.11, 15, 16.1).",
			},
			"priority": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Execution priority of the script (BEFORE, AFTER, AT_REBOOT).",
				ValidateFunc: validation.StringInSlice([]string{"BEFORE", "AFTER", "AT_REBOOT"}, false),
			},
			"script_contents": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Contents of the script. Must be non-compiled and in an accepted format.",
				DiffSuppressFunc: suppressBase64EncodedScriptDiff,
			},
			"parameter4": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Script parameter label 4",
			},
			"parameter5": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Script parameter label 5",
			},
			"parameter6": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Script parameter label 6",
			},
			"parameter7": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Script parameter label 7",
			},
			"parameter8": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Script parameter label 8",
			},
			"parameter9": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Script parameter label 9",
			},
			"parameter10": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Script parameter label 10",
			},
			"parameter11": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Script parameter label 11",
			},
		},
	}
}

// constructJamfProScript constructs a ResponsePolicy object from the provided schema data.
func constructJamfProScript(d *schema.ResourceData) (*jamfpro.ResourceScript, error) {
	script := &jamfpro.ResourceScript{}

	// Utilize type assertion helper functions for direct field extraction
	script.Name = util.GetStringFromInterface(d.Get("name"))
	script.CategoryName = util.GetStringFromInterface(d.Get("category_name"))
	script.CategoryId = util.GetStringFromInterface(d.Get("category_id"))
	script.Info = util.GetStringFromInterface(d.Get("info"))
	script.Notes = util.GetStringFromInterface(d.Get("notes"))
	script.OSRequirements = util.GetStringFromInterface(d.Get("os_requirements"))
	script.Priority = util.GetStringFromInterface(d.Get("priority"))

	// Extracting script parameters
	script.Parameter4 = util.GetStringFromInterface(d.Get("parameter4"))
	script.Parameter5 = util.GetStringFromInterface(d.Get("parameter5"))
	script.Parameter6 = util.GetStringFromInterface(d.Get("parameter6"))
	script.Parameter7 = util.GetStringFromInterface(d.Get("parameter7"))
	script.Parameter8 = util.GetStringFromInterface(d.Get("parameter8"))
	script.Parameter9 = util.GetStringFromInterface(d.Get("parameter9"))
	script.Parameter10 = util.GetStringFromInterface(d.Get("parameter10"))
	script.Parameter11 = util.GetStringFromInterface(d.Get("parameter11"))

	// Handle script_contents
	if scriptContent, ok := d.GetOk("script_contents"); ok {
		script.ScriptContents = util.GetStringFromInterface(scriptContent)
	} else {
		// Decode script contents from the state if not directly modified
		encodedScriptContents := util.GetStringFromInterface(d.Get("script_contents_encoded"))
		decodedBytes, err := base64.StdEncoding.DecodeString(encodedScriptContents)
		if err != nil {
			return nil, fmt.Errorf("error decoding script contents: %s", err)
		}
		script.ScriptContents = string(decodedBytes)
	}

	log.Printf("[INFO] Successfully constructed Script with name: %s", script.Name)

	return script, nil
}

// Helper function to generate diagnostics based on the error type.
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
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the create operation
	var createdAttribute *jamfpro.ResponseScriptCreate
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the script
		attribute, err := constructJamfProScript(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the script for terraform create: %w", err))
		}

		// Directly call the API to create the resource
		createdAttribute, err = conn.CreateScript(attribute)
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
	d.SetId(createdAttribute.ID)

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

// ResourceJamfProScriptsUpdate is responsible for updating an existing Jamf Pro Script on the remote system.
func ResourceJamfProScriptsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the update operation
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the updated script
		script, err := constructJamfProScript(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the script for terraform update: %w", err))
		}

		// If script_contents has not been modified, decode it from the state to ensure that base64
		// decoded version of script payload is sent to jamf.
		if !d.HasChange("script_contents") {
			encodedScriptContents := d.Get("script_contents_encoded").(string)
			decodedBytes, err := base64.StdEncoding.DecodeString(encodedScriptContents)
			if err != nil {
				return retry.NonRetryableError(fmt.Errorf("error decoding script contents: %s", err))
			}
			script.ScriptContents = string(decodedBytes)
		}

		// Obtain the ID from the Terraform state to be used for the API request
		scriptID := d.Id()

		// Directly call the API to update the resource
		_, apiErr := conn.UpdateScriptByID(scriptID, script)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by name
			scriptName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string in update"))
			}

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
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the delete operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Obtain the ID from the Terraform state to be used for the API request
		scriptID := d.Id()

		// Directly call the API to delete the resource
		apiErr := conn.DeleteScriptByID(scriptID)
		if apiErr != nil {
			// If the delete by ID fails, try deleting by name
			scriptName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

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
