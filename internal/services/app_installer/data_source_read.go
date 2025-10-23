package app_installer

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Check if either ID or name is provided
	resourceID := d.Get("id").(string)
	name := d.Get("name").(string)

	if resourceID == "" && name == "" {
		return diag.FromErr(fmt.Errorf("either id or name must be provided"))
	}

	var appInstaller *jamfpro.ResourceJamfAppCatalogDeployment
	var err error

	if resourceID != "" {
		// Get deployment directly by ID
		appInstaller, err = client.GetJamfAppCatalogAppInstallerByID(resourceID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch Jamf Pro App Installer Deployment by ID %s: %v", resourceID, err))
		}
	} else {
		// Get deployment by name
		appInstaller, err = client.GetJamfAppCatalogAppInstallerByName(name)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch Jamf Pro App Installer Deployment by name %s: %v", name, err))
		}
	}

	// Set ID from the deployment
	d.SetId(appInstaller.ID)

	// Update state using the same function as the resource
	return append(diags, updateState(d, appInstaller)...)
}
