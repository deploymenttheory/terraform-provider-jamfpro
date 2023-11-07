// computerextensionattributes_resource.go
package computerextensionattributes

import (
	"context"
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

// ResourceJamfProComputerExtensionAttributes defines the schema and CRUD operations (Create, Read, Update, Delete)
// for managing Jamf Pro Computer Extension Attributes in Terraform.
func ResourceJamfProComputerExtensionAttributes() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProComputerExtensionAttributesCreate,
		ReadContext:   ResourceJamfProComputerExtensionAttributesRead,
		UpdateContext: ResourceJamfProComputerExtensionAttributesUpdate,
		DeleteContext: ResourceJamfProComputerExtensionAttributesDelete,
		CustomizeDiff: validateResourceComputerExtensionAttributesDataFields,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute), // default timeout for create operation
			Read:   schema.DefaultTimeout(1 * time.Minute),  // default timeout for read operation
			Update: schema.DefaultTimeout(30 * time.Minute), // default timeout for update operation
			Delete: schema.DefaultTimeout(15 * time.Minute), // default timeout for delete operation
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the computer extension attribute.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro computer extension attribute.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if the computer extension attribute is enabled.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the computer extension attribute.",
			},
			"data_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Data type of the computer extension attribute. Can be String / Integer / Date (YYYY-MM-DD hh:mm:ss)",
				ValidateFunc: validateDataType,
			},
			"input_type": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"script", "Text Field", "LDAP Mapping", "Pop-up Menu"}, false),
						},
						"platform": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "Platform type for the computer extension attribute.",
							ValidateFunc: validation.StringInSlice([]string{"Mac", "Windows"}, false),
						},
						"script": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"choices": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				Description: "Input type details of the computer extension attribute.",
			},
			"inventory_display": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Display details for inventory for the computer extension attribute.",
				ValidateFunc: validation.StringInSlice([]string{"General", "Hardware", "Operating System", "User and Location", "Purchasing", "Extension Attributes"}, false),
			},
			"recon_display": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Display details for recon for the computer extension attribute.",
			},
		},
	}
}

// constructComputerExtensionAttribute constructs a ResponseComputerExtensionAttribute object from the provided schema data.
// It captures attributes from the schema, including nested attributes under "input_type", and returns the constructed object.
func constructComputerExtensionAttribute(d *schema.ResourceData) *jamfpro.ResponseComputerExtensionAttribute {
	// Extract the first item from the input_type list, which should be a map
	inputTypes := d.Get("input_type").([]interface{})
	if len(inputTypes) == 0 {
		// Handle this case as you see fit; here I'll just return nil to indicate an error
		return nil
	}

	inputTypeMap := inputTypes[0].(map[string]interface{})

	// Construct the InputType struct, capturing attributes such as type, platform, and script
	inputType := jamfpro.ComputerExtensionAttributeInputType{
		Type:     inputTypeMap["type"].(string),
		Platform: inputTypeMap["platform"].(string),
		Script:   inputTypeMap["script"].(string),
	}

	// If choices are provided under "input_type", extract and append them to the InputType struct
	if choices, exists := inputTypeMap["choices"]; exists {
		for _, choice := range choices.([]interface{}) {
			inputType.Choices = append(inputType.Choices, choice.(string))
		}
	}

	// Construct the ResponseComputerExtensionAttribute object
	attribute := &jamfpro.ResponseComputerExtensionAttribute{
		Name:             d.Get("name").(string),
		Enabled:          d.Get("enabled").(bool),
		Description:      d.Get("description").(string),
		DataType:         d.Get("data_type").(string),
		InventoryDisplay: d.Get("inventory_display").(string),
		ReconDisplay:     d.Get("recon_display").(string),
		InputType:        inputType,
	}

	// Log the successful construction of the attribute
	log.Printf("[INFO] Successfully constructed ComputerExtensionAttribute with name: %s", attribute.Name)

	return attribute
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

// ResourceJamfProComputerExtensionAttributesCreate is responsible for creating a new Jamf Pro Computer Extension Attribute in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProComputerExtensionAttributesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the create operation
	var createdAttribute *jamfpro.ResponseComputerExtensionAttribute
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the computer extension attribute
		attribute := constructComputerExtensionAttribute(d)

		// Check if the attribute is nil (indicating an issue with input_type)
		if attribute == nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the computer extension attribute due to missing or invalid input_type"))
		}

		// Log the details of the attribute that is about to be created
		log.Printf("[INFO] Attempting to create ComputerExtensionAttribute with name: %s", attribute.Name)

		// Directly call the API to create the resource
		createdAttribute, err = conn.CreateComputerExtensionAttribute(attribute)
		if err != nil {
			// Log the error from the API call
			log.Printf("[ERROR] Error creating ComputerExtensionAttribute with name: %s. Error: %s", attribute.Name, err)

			// Check if the error is an APIError
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			// For simplicity, we're considering all other errors as retryable
			return retry.RetryableError(err)
		}

		// Log the response from the API call
		log.Printf("[INFO] Successfully created ComputerExtensionAttribute with ID: %d and name: %s", createdAttribute.ID, createdAttribute.Name)

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
		readDiags := ResourceJamfProComputerExtensionAttributesRead(ctx, d, meta)
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

// ResourceJamfProComputerExtensionAttributesRead is responsible for reading the current state of a Jamf Pro Computer Extension Attribute from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProComputerExtensionAttributesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	var attribute *jamfpro.ResponseComputerExtensionAttribute

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		attributeID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse attribute ID: %v", convertErr))
		}

		// Try fetching the computer extension attribute using the ID
		var apiErr error
		attribute, apiErr = conn.GetComputerExtensionAttributeByID(attributeID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			attributeName := d.Get("name").(string)
			attribute, apiErr = conn.GetComputerExtensionAttributeByName(attributeName)
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
	if err := d.Set("enabled", attribute.Enabled); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("description", attribute.Description); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("data_type", attribute.DataType); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("inventory_display", attribute.InventoryDisplay); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("recon_display", attribute.ReconDisplay); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Extract the input type details and set them in the state
	inputType := map[string]interface{}{
		"type":     attribute.InputType.Type,
		"platform": attribute.InputType.Platform,
		"script":   attribute.InputType.Script,
		"choices":  attribute.InputType.Choices,
	}
	if attribute.InputType.Choices == nil || len(attribute.InputType.Choices) == 0 {
		inputType["choices"] = []string{}
	}

	// Wrap the map in a slice and set it to the Terraform state
	if err := d.Set("input_type", []interface{}{inputType}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProComputerExtensionAttributesUpdate is responsible for updating an existing Jamf Pro Computer Extension Attribute on the remote system.
func ResourceJamfProComputerExtensionAttributesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the update operation
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the updated computer extension attribute
		attribute := constructComputerExtensionAttribute(d)

		// Convert the ID from the Terraform state into an integer to be used for the API request
		attributeID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse attribute ID: %v", convertErr))
		}

		// Directly call the API to update the resource
		_, apiErr := conn.UpdateComputerExtensionAttributeByID(attributeID, attribute)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by name
			attributeName := d.Get("name").(string)
			_, apiErr = conn.UpdateComputerExtensionAttributeByName(attributeName, attribute)
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
		// If there's an error while update the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	// Use the retry function for the read operation to update the Terraform state
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProComputerExtensionAttributesRead(ctx, d, meta)
		if len(readDiags) > 0 {
			return retry.RetryableError(fmt.Errorf("failed to update the Terraform state for the updated resource"))
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while update the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	return diags
}

// ResourceJamfProComputerExtensionAttributesDelete is responsible for deleting a Jamf Pro Computer Extension Attribute.
func ResourceJamfProComputerExtensionAttributesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the delete operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		attributeID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse attribute ID: %v", convertErr))
		}

		// Directly call the API to delete the resource
		apiErr := conn.DeleteComputerExtensionAttributeByID(attributeID)
		if apiErr != nil {
			// If the delete by ID fails, try deleting by name
			attributeName := d.Get("name").(string)
			apiErr = conn.DeleteComputerExtensionAttributeByNameByID(attributeName)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while update the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "delete")
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
