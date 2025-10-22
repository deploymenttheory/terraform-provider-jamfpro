package app_installer

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/sdkv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProAppInstallersCreate is responsible for creating a new Jamf Pro App Installer in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// ==== Temp removal since jamf broke this end point in v11.14.0
	// will reinstate with jamf fix it.
	//
	// Check and accept the Jamf App Catalog App Installer terms and conditions
	// err := checkJamfAppCatalogAppInstallerTermsAndConditions(ctx, client)
	// if err != nil {
	// 	return diag.FromErr(fmt.Errorf("failed to ensure Jamf Pro App Installer terms and conditions are accepted: %v", err))
	// }

	resource, err := construct(d, client)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro App Installer: %v", err))
	}

	var creationResponse *jamfpro.ResponseJamfAppCatalogDeploymentCreateAndUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateJamfAppCatalogAppInstallerDeployment(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro App Installer '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(creationResponse.ID)

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read reads and states a jamfpro building
func read(ctx context.Context, d *schema.ResourceData, meta any, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	resourceID := d.Id()
	var diags diag.Diagnostics

	var response *jamfpro.ResourceJamfAppCatalogDeployment
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetJamfAppCatalogAppInstallerByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, sdkv2.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, updateState(d, response)...)
}

// readWithCleanup reads a resources and states with cleanup
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads a resource without cleanup
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update updates a jamfpro building
func update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	return sdkv2.Update(
		ctx,
		d,
		meta,
		func(d *schema.ResourceData) (*jamfpro.ResourceJamfAppCatalogDeployment, error) {
			return construct(d, client)
		},
		client.UpdateJamfAppCatalogAppInstallerDeploymentByID,
		readNoCleanup,
	)
}

func delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return sdkv2.Delete(
		ctx,
		d,
		meta,
		meta.(*jamfpro.Client).DeleteJamfAppCatalogAppInstallerDeploymentByID,
	)
}
