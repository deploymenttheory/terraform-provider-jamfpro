package computergroups

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/waitfor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProComputerGroupsCreate is responsible for creating a new Jamf Pro Computer Group in the remote system.
func ResourceJamfProComputerGroupsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics

	resource, err := constructJamfProComputerGroup(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Computer Group: %v", err))
	}

	var creationResponse *jamfpro.ResourceComputerGroup
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateComputerGroup(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Computer Group '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	checkResourceExists := func(id interface{}) (interface{}, error) {
		intID, err := strconv.Atoi(id.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
		}
		return apiclient.Conn.GetComputerGroupByID(intID)
	}

	_, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Computer Group", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second, apiclient.EnableCookieJar)

	if waitDiags.HasError() {
		return waitDiags
	}

	readDiags := ResourceJamfProComputerGroupsRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProComputerGroupsRead is responsible for reading the current state of a Jamf Pro Computer Group from the remote system.
func ResourceJamfProComputerGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := conn.GetComputerGroupByID(resourceIDInt)

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	diags = updateTerraformState(d, resource)

	if len(diags) > 0 {
		return diags
	}
	return nil
}

// ResourceJamfProComputerGroupsUpdate is responsible for updating an existing Jamf Pro Computer Group on the remote system.
func ResourceJamfProComputerGroupsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructJamfProComputerGroup(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Computer Group for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateComputerGroupByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Computer Group '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	readDiags := ResourceJamfProComputerGroupsRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProComputerGroupsDelete is responsible for deleting a Jamf Pro Computer Group.
func ResourceJamfProComputerGroupsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {

		apiErr := conn.DeleteComputerGroupByID(resourceIDInt)
		if apiErr != nil {

			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeleteComputerGroupByName(resourceName)
			if apiErrByName != nil {

				return retry.RetryableError(apiErrByName)
			}
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Computer Group '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	d.SetId("")

	return diags
}
