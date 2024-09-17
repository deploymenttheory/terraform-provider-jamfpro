package managedsoftwareupdates

import (
	"context"
	"fmt"
	"sync"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Create requires a mutex need to lock Create requests during parallel runs
var mu sync.Mutex

// create is responsible for creating a new Jamf Pro Managed Software Update in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Lock the mutex to ensure only one create operation can run this function at a time
	mu.Lock()
	defer mu.Unlock()

	// Check and accept the Jamf Managed Software Update terms and conditions
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
		creationResponse, apiErr = client.CreateManagedSoftwareUpdatePlanByGroupID(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Managed Software Update '%s' after retries: %v", resource.Group.GroupId, err))
	}

	d.SetId(creationResponse.Plans[0].PlanID)

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read reads and states a jamfpro managed software update plan
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	resourceUUID := d.Id()
	var diags diag.Diagnostics

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
