// apiintegrations_resource.go
package apiintegrations

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProApiIntegrations defines the schema and CRUD operations for managing Jamf Pro API Integrations in Terraform.
func ResourceJamfProApiIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProApiIntegrationsCreate,
		ReadContext:   ResourceJamfProApiIntegrationsRead,
		UpdateContext: ResourceJamfProApiIntegrationsUpdate,
		DeleteContext: ResourceJamfProApiIntegrationsDelete,
		CustomizeDiff: validateResourceAPIIntegrationsDataFields,
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
				Description: "The unique identifier of the API integration.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the API integration.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if the API integration is enabled.",
			},
			"access_token_lifetime_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The access token lifetime in seconds for the API integration.",
			},
			"app_type": {
				Type:     schema.TypeString,
				Computed: true,
				//Required:     true,
				Description: "The app type of the API integration.",
				//ValidateFunc: validation.StringInSlice([]string{"CLIENT_CREDENTIALS", "NATIVE_APP_OAUTH", "NONE"}, false),
			},
			"client_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The client ID of the API integration.",
			},
			"authorization_scopes": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of authorization roles scoped to the API integration.",
			},
		},
	}
}

func constructJamfProApiIntegration(d *schema.ResourceData) (*jamfpro.ResourceApiIntegration, error) {
	integration := &jamfpro.ResourceApiIntegration{}

	// Utilize helper functions for direct field extraction
	integration.DisplayName = util.GetStringFromInterface(d.Get("display_name"))
	integration.Enabled = util.GetBoolFromInterface(d.Get("enabled"))
	integration.AccessTokenLifetimeSeconds = util.GetIntFromInterface(d.Get("access_token_lifetime_seconds"))

	// Handle 'authorization_scopes' field
	if v, ok := d.GetOk("authorization_scopes"); ok {
		integration.AuthorizationScopes = convertToStringSlice(v.(*schema.Set))
	}

	// Log the successful construction
	log.Printf("[INFO] Successfully constructed ApiIntegration with display name: %s", integration.DisplayName)

	return integration, nil
}

// Helper function to generate diagnostics based on the error type..
func generateTFDiagsFromHTTPError(err error, d *schema.ResourceData, action string) diag.Diagnostics {
	var diags diag.Diagnostics
	resourceName, exists := d.GetOk("name")
	if !exists {
		resourceName = "unknown"
	}

	// Handle the APIError in the diagnostic.
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

// convertToStringSlice is a helper function that converts a schema.Set to a string slice.
func convertToStringSlice(set *schema.Set) []string {
	interfaceSlice := set.List()
	stringSlice := make([]string, len(interfaceSlice))
	for i, v := range interfaceSlice {
		stringSlice[i] = v.(string)
	}
	return stringSlice
}

// ResourceJamfProApiIntegrationsCreate is responsible for creating a new Jamf Pro API Integration in the remote system.
func ResourceJamfProApiIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the create operation
	var createdIntegration *jamfpro.ResourceApiIntegration
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the API integration
		integration, err := constructJamfProApiIntegration(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the api integration for terraform create: %w", err))
		}

		// Log the details of the integration that is about to be created
		log.Printf("[INFO] Attempting to create ApiIntegration with display name: %s", integration.DisplayName)

		// Directly call the API to create the resource
		createdIntegration, err = conn.CreateApiIntegration(integration)
		if err != nil {
			// Log the error from the API call
			log.Printf("[ERROR] Error creating ApiIntegration with display name: %s. Error: %s", integration.DisplayName, err)

			// Check if the error is an APIError
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			// For simplicity, we're considering all other errors as retryable
			return retry.RetryableError(err)
		}

		// Log the response from the API call
		log.Printf("[INFO] Successfully created ApiIntegration with ID: %d and display name: %s", createdIntegration.ID, createdIntegration.DisplayName)

		return nil
	})

	if err != nil {
		// If there's an error while creating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "create")
	}

	// Set the ID of the created resource in the Terraform state
	d.SetId(strconv.Itoa(createdIntegration.ID))

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProApiIntegrationsRead(ctx, d, meta)
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

// ResourceJamfProApiIntegrationsRead is responsible for reading the current state of a Jamf Pro API Integration from the remote system.
func ResourceJamfProApiIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var integration *jamfpro.ResourceApiIntegration

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		integrationID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse integration ID: %v", convertErr))
		}

		// Try fetching the API integration using the ID
		var apiErr error
		integration, apiErr = conn.GetApiIntegrationByID(integrationID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			integrationName := d.Get("display_name").(string)
			integration, apiErr = conn.GetApiIntegrationByName(integrationName)
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
	if err := d.Set("display_name", integration.DisplayName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("enabled", integration.Enabled); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("access_token_lifetime_seconds", integration.AccessTokenLifetimeSeconds); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("app_type", integration.AppType); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("authorization_scopes", integration.AuthorizationScopes); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("client_id", integration.ClientID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProApiIntegrationsUpdate is responsible for updating an existing Jamf Pro API Integration on the remote system.
func ResourceJamfProApiIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		// Construct the API integration
		integration, err := constructJamfProApiIntegration(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the api integration for terraform update: %w", err))
		}

		// Convert the ID from the Terraform state into an integer to be used for the API request
		integrationID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse integration ID: %v", convertErr))
		}

		// Directly call the API to update the resource
		_, apiErr := conn.UpdateApiIntegrationByID((integrationID), integration)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by display name
			integrationName := d.Get("display_name").(string)
			_, apiErr = conn.UpdateApiIntegrationByName(integrationName, integration)
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
		readDiags := ResourceJamfProApiIntegrationsRead(ctx, d, meta)
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

// ResourceJamfProApiIntegrationsDelete is responsible for deleting a Jamf Pro API Integration.
func ResourceJamfProApiIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		integrationID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse integration ID: %v", convertErr))
		}

		// Directly call the API to delete the resource
		apiErr := conn.DeleteApiIntegrationByID(integrationID)
		if apiErr != nil {
			// If the delete by ID fails, try deleting by display name
			integrationName := d.Get("display_name").(string)
			apiErr = conn.DeleteApiIntegrationByName(integrationName)
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
