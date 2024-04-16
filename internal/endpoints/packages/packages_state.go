// packages_state.go
package packages

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Package information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourcePackage) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update Terraform state with the resource information
	if err := d.Set("id", strconv.Itoa(resource.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resource.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	// Check if the category is "No category assigned" and set it to "Unknown"
	// This is necessary because the API returns "No category assigned" when no category is assigned
	// but the request expects "Unknown" when no category is assigned.
	if resource.Category == "No category assigned" {
		// Set the category to "Unknown"
		if err := d.Set("category", "Unknown"); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		// Set the category normally if it's not "No category assigned"
		if err := d.Set("category", resource.Category); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	if err := d.Set("filename", resource.Filename); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("info", resource.Info); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("notes", resource.Notes); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("priority", resource.Priority); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("reboot_required", resource.RebootRequired); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("fill_user_template", resource.FillUserTemplate); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("fill_existing_users", resource.FillExistingUsers); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("boot_volume_required", resource.BootVolumeRequired); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("allow_uninstalled", resource.AllowUninstalled); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("os_requirements", resource.OSRequirements); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	/* Fields are in the data model but don't appear to serve a purpose in jamf 11.3 onwards
	// these fields may only be relevant if a file is indexed by JAMF Admin. which i *think*
	// is to be deprecated in favor of JCDS 2.0
	if err := d.Set("required_processor", resource.RequiredProcessor); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("switch_with_package", resource.SwitchWithPackage); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("reinstall_option", resource.ReinstallOption); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("triggering_files", resource.TriggeringFiles); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	*/
	if err := d.Set("install_if_reported_available", resource.InstallIfReportedAvailable); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("send_notification", resource.SendNotification); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
