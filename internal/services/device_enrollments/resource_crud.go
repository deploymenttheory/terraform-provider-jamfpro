package device_enrollments

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for creating a new Jamf Pro Device Enrollment in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tokenUpload := &jamfpro.ResourceDeviceEnrollmentTokenUpload{
		EncodedToken: d.Get("encoded_token").(string),
	}

	client := meta.(*jamfpro.Client)
	enrollment, err := client.CreateDeviceEnrollmentWithMDMServerToken(tokenUpload)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(enrollment.ID)

	wrappedUpdate := func(payload *jamfpro.ResourceDeviceEnrollmentUpdate) (*jamfpro.ResourceDeviceEnrollment, error) {
		return client.UpdateDeviceEnrollmentMetadataByID(enrollment.ID, payload)
	}

	return common.Create(
		ctx,
		d,
		meta,
		construct,
		wrappedUpdate,
		readNoCleanup,
	)
}

// read is responsible for reading the current state of a Jamf Pro Device Enrollment from Jamf Pro
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	return common.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetDeviceEnrollmentByID,
		updateState,
	)
}

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update is responsible for updating an existing Jamf Pro Device Enrollment on the remote system
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	if d.HasChange("encoded_token") {
		tokenUpload := &jamfpro.ResourceDeviceEnrollmentTokenUpload{
			EncodedToken: d.Get("encoded_token").(string),
		}
		_, err := client.UpdateDeviceEnrollmentMDMServerToken(d.Id(), tokenUpload)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return common.Update(
		ctx,
		d,
		meta,
		construct,
		client.UpdateDeviceEnrollmentMetadataByID,
		readNoCleanup,
	)
}

// delete is responsible for deleting a Jamf Pro Device Enrollment
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteDeviceEnrollmentByID,
	)
}
