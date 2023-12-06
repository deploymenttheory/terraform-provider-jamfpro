// scripts_resource.go
package scripts

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
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
				Computed:    true,
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
				Computed:    true,
				Description: "Labels to use for script parameters. Parameters 1 through 3 are predefined as mount point, computer name, and username",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
			},
			"os_requirements": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The script can only be run on computers with these operating system versions. Each version must be separated by a comma (e.g., 10.11, 15, 16.1).",
			},
			"script_contents": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Contents of the script. The script contents must be non-compiled and in one of the following formats: Bash (.sh), Shell (.sh), Non-compiled AppleScript (.applescript), C Shell (.csh), Zsh (.zsh),Korn Shell (.ksh), Tool Command Language (.tcl), and Python (.py). Ref - https://learn.jamf.com/bundle/jamf-pro-documentation-current/page/Scripts.html . The script contents should also have trailing whitespace at the end removed, to avoid tf state false positives.",
				DiffSuppressFunc: suppressBase64EncodedScriptDiff,
			},
			"script_contents_encoded": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Jamf Pro encoded contents of the script.",
			},
		},
	}
}

// constructJamfProScript constructs a ResponseScript object from the provided schema data and returns any errors encountered.
func constructJamfProScript(d *schema.ResourceData) (*jamfpro.ResponseScript, error) {
	script := &jamfpro.ResponseScript{}

	// construct regular fields
	fields := map[string]interface{}{
		"name":     &script.Name,
		"filename": &script.Filename,
		"info":     &script.Info,
		//"script_contents": &script.ScriptContents,
		"notes":           &script.Notes,
		"priority":        &script.Priority,
		"os_requirements": &script.OSRequirements,
	}

	for key, ptr := range fields {
		if v, ok := d.GetOk(key); ok {
			*ptr.(*string) = v.(string)
		}
	}

	// construct script_contents
	scriptContents, isScriptContentPresent := d.GetOk("script_contents")
	if isScriptContentPresent {
		// If the script contents were modified, use the new value directly
		script.ScriptContents = scriptContents.(string)
	} else {
		// If the script contents were not modified, decode them from the state
		encodedScriptContents, _ := d.Get("script_contents_encoded").(string)
		decodedBytes, err := base64.StdEncoding.DecodeString(encodedScriptContents)
		if err != nil {
			return nil, fmt.Errorf("error decoding script contents: %s", err)
		}
		script.ScriptContents = string(decodedBytes)
	}

	// construct fields with default values
	if v, ok := d.GetOk("category"); ok {
		script.Category = v.(string)
	} else {
		script.Category = "No category assigned"
	}

	// construct nested fields
	if params, ok := d.GetOk("parameters"); ok {
		paramsList, ok := params.([]interface{})
		if !ok || len(paramsList) == 0 {
			return nil, fmt.Errorf("invalid data for 'parameters'")
		}
		paramMap, ok := paramsList[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid data structure for 'parameters'")
		}

		script.Parameters = jamfpro.Parameters{
			Parameter4:  getStringFromMap(paramMap, "parameter4"),
			Parameter5:  getStringFromMap(paramMap, "parameter5"),
			Parameter6:  getStringFromMap(paramMap, "parameter6"),
			Parameter7:  getStringFromMap(paramMap, "parameter7"),
			Parameter8:  getStringFromMap(paramMap, "parameter8"),
			Parameter9:  getStringFromMap(paramMap, "parameter9"),
			Parameter10: getStringFromMap(paramMap, "parameter10"),
			Parameter11: getStringFromMap(paramMap, "parameter11"),
		}
	}

	// Log the successful construction of the script
	log.Printf("[INFO] Successfully constructed Script with name: %s", script.Name)

	return script, nil
}

// getStringFromMap is a helper function to safely extract string values from a map.
// Returns an empty string if the key is not found.
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok && val != "" {
		return val.(string)
	}
	return ""
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
	var createdAttribute *jamfpro.ResponseScript
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the script
		attribute, err := constructJamfProScript(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the script for terraform create: %w", err))
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
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

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
			attributeName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

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
	if attribute.Category == "" {
		if err := d.Set("category", "No category assigned"); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		if err := d.Set("category", attribute.Category); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	if err := d.Set("filename", attribute.Filename); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("info", attribute.Info); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("notes", attribute.Notes); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("priority", attribute.Priority); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("os_requirements", attribute.OSRequirements); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	// Fetch script content from Jamf Pro
	scriptContentFromJamf := attribute.ScriptContents
	// Encode the Jamf Pro script content to base64
	encodedContentFromJamf := encodeScriptContent(scriptContentFromJamf)
	// Update the Terraform state with the base64 encoded script content
	if err := d.Set("script_contents", encodedContentFromJamf); err != nil {
		return diag.FromErr(err)
	}

	// Handling parameters
	parameters := make(map[string]interface{})
	if attribute.Parameters.Parameter4 != "" {
		parameters["parameter4"] = attribute.Parameters.Parameter4
	}
	if attribute.Parameters.Parameter5 != "" {
		parameters["parameter5"] = attribute.Parameters.Parameter5
	}
	// Apply this pattern to the rest of the parameters
	if attribute.Parameters.Parameter6 != "" {
		parameters["parameter6"] = attribute.Parameters.Parameter6
	}
	if attribute.Parameters.Parameter7 != "" {
		parameters["parameter7"] = attribute.Parameters.Parameter7
	}
	if attribute.Parameters.Parameter8 != "" {
		parameters["parameter8"] = attribute.Parameters.Parameter8
	}
	if attribute.Parameters.Parameter9 != "" {
		parameters["parameter9"] = attribute.Parameters.Parameter9
	}
	if attribute.Parameters.Parameter10 != "" {
		parameters["parameter10"] = attribute.Parameters.Parameter10
	}
	if attribute.Parameters.Parameter11 != "" {
		parameters["parameter11"] = attribute.Parameters.Parameter11
	}

	if len(parameters) > 0 {
		if err := d.Set("parameters", []interface{}{parameters}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		// Explicitly setting parameters to nil if they are absent
		if err := d.Set("parameters", nil); err != nil {
			diags = append(diags, diag.FromErr(err)...)
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
			encodedScriptContents, _ := d.Get("script_contents_encoded").(string)
			decodedBytes, err := base64.StdEncoding.DecodeString(encodedScriptContents)
			if err != nil {
				return retry.NonRetryableError(fmt.Errorf("error decoding script contents: %s", err))
			}
			script.ScriptContents = string(decodedBytes)
		}

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
			scriptName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
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
		// Convert the ID from the Terraform state into an integer to be used for the API request
		scriptID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse script ID: %v", convertErr))
		}

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
