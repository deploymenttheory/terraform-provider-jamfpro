package filesharedistributionpoints

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	JamfProResourceDistributionPoint = "Distribution Point"
)

// resourceJamfProFileShareDistributionPointsCreate is responsible for creating a new file share
// distribution point object in the remote system.
// The function:
// 1. Constructs the dock item data using the provided Terraform configuration.
// 2. Calls the API to create the dock item in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created dock item.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func resourceJamfProFileShareDistributionPointsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProFileShareDistributionPoint(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Fileshare Distribution Point: %v", err))
	}

	var creationResponse *jamfpro.ResponseFileShareDistributionPointCreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateDistributionPoint(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Fileshare Distribution Point '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	return append(diags, resourceJamfProFileShareDistributionPointsReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProFileShareDistributionPointsRead is responsible for reading the current state of a
// Jamf Pro File Share Distribution Point Resource from the remote system.
// The function:
// 1. Fetches the dock item's current state using its ID. If it fails then obtain dock item's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the dock item being deleted outside of Terraform, to keep the Terraform state synchronized.
func resourceJamfProFileShareDistributionPointsRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	resourceID := d.Id()
	var diags diag.Diagnostics

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var response *jamfpro.ResourceFileShareDistributionPoint
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetDistributionPointByID(resourceIDInt)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, updateTerraformState(d, response)...)
}

// resourceJamfProFileShareDistributionPointsReadWithCleanup reads the resource with cleanup enabled
func resourceJamfProFileShareDistributionPointsReadWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProFileShareDistributionPointsRead(ctx, d, meta, true)
}

// resourceJamfProFileShareDistributionPointsReadNoCleanup reads the resource with cleanup disabled
func resourceJamfProFileShareDistributionPointsReadNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProFileShareDistributionPointsRead(ctx, d, meta, false)
}

// resourceJamfProFileShareDistributionPointsUpdate is responsible for updating an existing Jamf Pro Site on the remote system.
func resourceJamfProFileShareDistributionPointsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructJamfProFileShareDistributionPoint(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro file share distribution point for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateDistributionPointByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro file share distribution point '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	return append(diags, resourceJamfProFileShareDistributionPointsReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProFileShareDistributionPointsDeleteis responsible for deleting a Jamf Pro file share distribution point from the remote system.
func resourceJamfProFileShareDistributionPointsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteDistributionPointByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteDistributionPointByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro file share distribution point '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	d.SetId("")

	return diags
}
