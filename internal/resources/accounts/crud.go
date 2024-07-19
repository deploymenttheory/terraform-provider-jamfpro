package accounts

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProAccountCreate is responsible for creating a new Jamf Pro Script in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Create(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreateAccount,
		readNoCleanup,
	)
}

// resourceJamfProAccountRead is responsible for reading the current state of a Jamf Pro Account Group Resource from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	return common.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetAccountByID,
		updateTerraformState,
	)
}

// resourceJamfProAccountReadWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// resourceJamfProAccountReadNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// resourceJamfProAccountRead is responsible for reading the current state of a Jamf Pro Account Group Resource from the remote system.
// func readFromCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	client := meta.(*jamfpro.Client)
// 	var diags diag.Diagnostics
// 	resourceID := d.Id()

// 	resourceIDInt, err := strconv.Atoi(resourceID)
// 	if err != nil {
// 		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
// 	}

// 	var response *jamfpro.ResourceAccount
// 	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
// 		var apiErr error
// 		response, apiErr = client.GetAccountByID(resourceIDInt)
// 		if apiErr != nil {
// 			return retry.RetryableError(apiErr)
// 		}
// 		return nil
// 	})

// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	return append(diags, updateTerraformState(d, response)...)
// }

// resourceJamfProAccountUpdate is responsible for updating an existing Jamf Pro Account Group on the remote system.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Update(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).UpdateAccountByID,
		readNoCleanup,
	)
}

// resourceJamfProAccountDelete is responsible for deleting a Jamf Pro account .
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteAccountByID,
	)
}
