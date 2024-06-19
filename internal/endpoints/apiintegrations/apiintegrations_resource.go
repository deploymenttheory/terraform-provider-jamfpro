// apiintegrations_resource.go
package apiintegrations

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

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
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
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
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProApiIntegration(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro API Integration: %v", err))
	}

	var creationResponse *jamfpro.ResourceApiIntegration
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateApiIntegration(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro API Integration '%s' after retries: %v", resource.DisplayName, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	intID, err := strconv.Atoi(id.(string))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
	// 	}
	// 	return client.GetApiIntegrationByID(intID)
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro API Integration", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)

	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, ResourceJamfProApiIntegrationsRead(ctx, d, meta)...)
}

// ResourceJamfProApiIntegrationsRead is responsible for reading the current state of a Jamf Pro API Integration from the remote system.
func ResourceJamfProApiIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := client.GetApiIntegrationByID(resourceIDInt)
	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProApiIntegrationsUpdate is responsible for updating an existing Jamf Pro API Integration on the remote system.
func ResourceJamfProApiIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := constructJamfProApiIntegration(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro API Integration for update: %v", err))
	}

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateApiIntegrationByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro API Integration '%s' (ID: %s) after retries: %v", resource.DisplayName, resourceID, err))
	}

	return append(diags, ResourceJamfProApiIntegrationsRead(ctx, d, meta)...)
}

// ResourceJamfProApiIntegrationsDelete is responsible for deleting a Jamf Pro API Integration.
func ResourceJamfProApiIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteApiIntegrationByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("display_name").(string)
			apiErrByName := client.DeleteApiIntegrationByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro API Integration '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
