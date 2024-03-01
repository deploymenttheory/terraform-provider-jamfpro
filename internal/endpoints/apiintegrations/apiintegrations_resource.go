// apiintegrations_resource.go
package apiintegrations

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

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
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
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

// ResourceJamfProApiIntegrationsCreate is responsible for creating a new Jamf Pro API Integration in the remote system.
func ResourceJamfProApiIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProApiIntegration(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Site: %v", err))
	}

	// Retry the API call to create the resource in Jamf Pro
	var creationResponse *jamfpro.ResourceApiIntegration
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateApiIntegration(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro API Integration '%s' after retries: %v", resource.DisplayName, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProApiIntegrationsRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProApiIntegrationsRead is responsible for reading the current state of a Jamf Pro API Integration from the remote system.
func ResourceJamfProApiIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var resource *jamfpro.ResourceApiIntegration

	// Read operation with retry
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = conn.GetApiIntegrationByID(resourceIDInt)
		if apiErr != nil {
			if strings.Contains(apiErr.Error(), "404") || strings.Contains(apiErr.Error(), "410") {
				// Return non-retryable error with a message to avoid SDK issues
				return retry.NonRetryableError(fmt.Errorf("resource not found, marked for deletion"))
			}
			// Retry for other types of errors
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	// If err is not nil, check if it's due to the resource being not found
	if err != nil {
		if err.Error() == "resource not found, marked for deletion" {
			// Resource not found, remove from Terraform state
			d.SetId("")
			// Append a warning diagnostic and return
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("Jamf Pro API Integration with ID '%s' was not found on the server and is marked for deletion from terraform state.", resourceID),
			})
			return diags
		}

		// For other errors, return an error diagnostic
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro API Integration with ID '%s' after retries: %v", resourceID, err))
	}

	// Map the configuration fields from the API response to a structured map
	apiIntegrationData := map[string]interface{}{
		"display_name":                  resource.DisplayName,
		"enabled":                       resource.Enabled,
		"access_token_lifetime_seconds": resource.AccessTokenLifetimeSeconds,
		"app_type":                      resource.AppType,
		"authorization_scopes":          resource.AuthorizationScopes,
		"client_id":                     resource.ClientID,
	}

	// Set the structured map in the Terraform state
	for key, val := range apiIntegrationData {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("failed to set '%s': %v", key, err))...)
		}
	}

	return diags
}

// ResourceJamfProApiIntegrationsUpdate is responsible for updating an existing Jamf Pro API Integration on the remote system.
func ResourceJamfProApiIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Construct the resource object
	resource, err := constructJamfProApiIntegration(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro API Integration for update: %v", err))
	}

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateApiIntegrationByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro API Integration '%s' (ID: %s) after retries: %v", resource.DisplayName, resourceID, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProApiIntegrationsRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProApiIntegrationsDelete is responsible for deleting a Jamf Pro API Integration.
func ResourceJamfProApiIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Attempt to delete by ID
		apiErr := conn.DeleteApiIntegrationByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("display_name").(string)
			apiErrByName := conn.DeleteApiIntegrationByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro API Integration '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
