package managed_software_update

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for creating a new Jamf Pro Managed Software Update in the remote system.
// The function:
// 1. Constructs the managed software update data using the provided Terraform configuration.
// 2. Calls the API to create the managed software update in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created managed software update.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	err := checkAndEnableManagedSoftwareUpdateFeatureToggle(ctx, client)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to ensure Jamf Pro Managed Software Update toggle is enabled: %v", err))
	}

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Managed Software Update: %v", err))
	}

	var creationResponse *jamfpro.ResponseManagedSoftwareUpdatePlanCreate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		if resource.Group.GroupId != "" {
			creationResponse, apiErr = client.CreateManagedSoftwareUpdatePlanByGroupID(resource)
		} else {
			creationResponse, apiErr = client.CreateManagedSoftwareUpdatePlanByDeviceID(resource)
		}
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Managed Software Update after retries: %v", err))
	}

	if len(creationResponse.Plans) > 0 {
		planUUID := creationResponse.Plans[0].PlanID
		d.SetId(planUUID)
		if err := d.Set("plan_uuid", planUUID); err != nil {
			return diag.FromErr(fmt.Errorf("error setting planID as plan_uuid: %v", err))
		}

		// Set group and object_type
		if resource.Group.GroupId != "" {
			if err := d.Set("group_id", resource.Group.GroupId); err != nil {
				return diag.FromErr(fmt.Errorf("error setting group_id: %v", err))
			}
			if err := d.Set("object_type", resource.Group.ObjectType); err != nil {
				return diag.FromErr(fmt.Errorf("error setting object_type: %v", err))
			}
		}
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceUUID := d.Id()

	var response *jamfpro.ResponseManagedSoftwareUpdatePlan

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetManagedSoftwareUpdatePlanByUUID(resourceUUID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, updateState(d, response)...)
}

// readWithCleanup reads a resources and states with cleanup
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads a resource without cleanup
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update updates a jamfpro managed software update plan
// Since there is no API endpoint to update a managed software update plan, we create a new one by calling the create function
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Call the create function to create a new plan
	return create(ctx, d, meta)
}

// delete deletes a jamfpro managed software update plan
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
