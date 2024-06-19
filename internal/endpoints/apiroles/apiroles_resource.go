// apiroles_resource.go
package apiroles

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProAPIRoles defines the schema for managing Jamf Pro API Roles in Terraform.
func ResourceJamfProAPIRoles() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProAPIRolesCreate,
		ReadContext:   ResourceJamfProAPIRolesRead,
		UpdateContext: ResourceJamfProAPIRolesUpdate,
		DeleteContext: ResourceJamfProAPIRolesDelete,
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
				Description: "The unique identifier of the Jamf API Role.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the Jamf API Role.",
			},
			"privileges": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "List of privileges associated with the Jamf API Role.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: func(val interface{}, key string) ([]string, []error) {
						return validateResourceApiRolesDataFields(val, key)
					},
				},
			},
		},
	}
}

// ResourceJamfProAPIRolesCreate handles the creation of a Jamf Pro API Role.
// The function:
// 1. Constructs the API role data using the provided Terraform configuration.
// 2. Calls the API to create the role in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created role.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProAPIRolesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProApiRole(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Account Group: %v", err))
	}

	var creationResponse *jamfpro.ResourceAPIRole
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateJamfApiRole(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro API Role '%s' after retries: %v", resource.DisplayName, err))
	}

	d.SetId(creationResponse.ID)

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	intID, err := strconv.Atoi(id.(string))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
	// 	}
	// 	return client.GetJamfApiRoleByID(strconv.Itoa(intID))
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro API Role", creationResponse.ID, checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)

	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, ResourceJamfProAPIRolesRead(ctx, d, meta)...)
}

// ResourceJamfProAPIRolesRead handles reading a Jamf Pro API Role from the remote system.
// The function:
// 1. Tries to fetch the API role based on the ID from the Terraform state.
// 2. If fetching by ID fails, attempts to fetch it by the display name.
// 3. Updates the Terraform state with the fetched data.
func ResourceJamfProAPIRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := client.GetJamfApiRoleByID(resourceID)

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProAPIRolesUpdate handles updating a Jamf Pro API Role.
// The function:
// 1. Constructs the updated API role data using the provided Terraform configuration.
// 2. Calls the API to update the role in Jamf Pro.
// 3. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProAPIRolesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resource, err := constructJamfProApiRole(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro API Role for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateJamfApiRoleByID(resourceID, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro API Role '%s' (ID: %s) after retries: %v", resource.DisplayName, resourceID, err))
	}

	return append(diags, ResourceJamfProAPIRolesRead(ctx, d, meta)...)
}

// ResourceJamfProAPIRolesDelete handles the deletion of a Jamf Pro API Role.
func ResourceJamfProAPIRolesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteJamfApiRoleByID(resourceID)
		if apiErr != nil {
			resourceName := d.Get("display_name").(string)
			apiErrByName := client.DeleteJamfApiRoleByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro API role '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
