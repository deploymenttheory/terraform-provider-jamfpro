// smartcomputergroup_crud.go
package smartcomputergroups

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProSmartComputerGroupsCreate is responsible for creating a new Jamf Pro Smart Computer Group in the remote system.
func ResourceJamfProSmartComputerGroupsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProSmartComputerGroup(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Smart Computer Group: %v", err))
	}

	var creationResponse *jamfpro.ResponseComputerGroupreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateComputerGroup(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Smart Computer Group '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	// Wait for the resource to be fully available before reading it
	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	intID, err := strconv.Atoi(id.(string))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
	// 	}
	// 	return client.GetComputerGroupByID(intID)
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Smart Computer Group", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)

	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, ResourceJamfProSmartComputerGroupsRead(ctx, d, meta)...)
}

// ResourceJamfProSmartComputerGroupsRead is responsible for reading the current state of a Jamf Pro Smart Computer Group from the remote system.
func ResourceJamfProSmartComputerGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := client.GetComputerGroupByID(resourceIDInt)

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// ResourceJamfProSmartComputerGroupsUpdate is responsible for updating an existing Jamf Pro Smart Computer Group on the remote system.
func ResourceJamfProSmartComputerGroupsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructJamfProSmartComputerGroup(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Smart Computer Group for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateComputerGroupByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Smart Computer Group '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	return append(diags, ResourceJamfProSmartComputerGroupsRead(ctx, d, meta)...)
}

// ResourceJamfProSmartComputerGroupsDelete is responsible for deleting a Jamf Pro Smart Computer Group.
func ResourceJamfProSmartComputerGroupsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteComputerGroupByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteComputerGroupByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Smart Computer Group '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	d.SetId("")

	return diags
}
