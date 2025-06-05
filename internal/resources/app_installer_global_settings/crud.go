package app_installer_global_settings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create initializes the App Installer Global Settings in Jamf Pro.
// This is a singleton configuration, so it always performs an update.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	settings, err := constructAppInstallerGlobalSettings(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro App Installer Global Settings: %w", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		_, apiErr := client.UpdateJamfAppCatalogAppInstallerGlobalSettings(settings)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to apply App Installer Global Settings after retries: %w", err))
	}

	d.SetId("jamfpro_app_installers_global_settings_singleton")
	return readNoCleanup(ctx, d, meta)
}

// read fetches and updates the state of the App Installer Global Settings from Jamf Pro.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	d.SetId("jamfpro_app_installers_global_settings_singleton")

	var fetched *jamfpro.JamfAppCatalogDeploymentSubsetNotificationSettings
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		fetched, apiErr = client.GetJamfAppCatalogAppInstallerGlobalSettings()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})
	if err != nil {
		return common.HandleResourceNotFoundError(err, d, cleanup)
	}

	// Update Terraform state with fetched data
	return updateState(d, &jamfpro.ResponseJamfAppCatalogGlobalSettings{
		EndUserExperienceSettings: *fetched,
	})
}

// readWithCleanup runs the read operation and allows cleanup on missing resource
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup runs the read operation without cleanup on missing resource
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update modifies the App Installer Global Settings in Jamf Pro.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	settings, err := constructAppInstallerGlobalSettings(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro App Installer Global Settings for update: %w", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateJamfAppCatalogAppInstallerGlobalSettings(settings)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update App Installer Global Settings after retries: %w", err))
	}

	d.SetId("jamfpro_app_installers_global_settings_singleton")
	return readNoCleanup(ctx, d, meta)
}

// delete removes the App Installer Global Settings from Terraform state only.
// The configuration remains in Jamf Pro (not deletable via API).
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
