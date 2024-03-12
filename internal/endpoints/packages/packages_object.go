// packages_data_object.go
package packages

import (
	"path/filepath"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProPackageCreate constructs a ResourcePackage object from the provided schema data.
// It extracts the filename from the full path provided in the schema and uses it for the Filename field.
func constructJamfProPackageCreate(d *schema.ResourceData) (*jamfpro.ResourcePackage, error) {
	// Extract the full file path from the schema
	fullPath := d.Get("package_file_path").(string)

	// Use filepath.Base to extract just the filename from the full path
	filename := filepath.Base(fullPath)

	// Get the category from the schema, and set it to "Unknown" if it's empty
	// Unknown is the valid default value for the category field not "No category assigned"
	// which is returned by the API when no category is assigned. MENTAL >_<
	category := d.Get("category").(string)
	if category == "" {
		category = "Unknown"
	}

	packageResource := &jamfpro.ResourcePackage{
		Name:               d.Get("name").(string),
		Filename:           filename,
		Category:           category,
		Info:               d.Get("info").(string),
		Notes:              d.Get("notes").(string),
		Priority:           d.Get("priority").(int),
		RebootRequired:     d.Get("reboot_required").(bool),
		FillUserTemplate:   d.Get("fill_user_template").(bool),
		FillExistingUsers:  d.Get("fill_existing_users").(bool),
		BootVolumeRequired: d.Get("boot_volume_required").(bool),
		AllowUninstalled:   d.Get("allow_uninstalled").(bool),
		OSRequirements:     d.Get("os_requirements").(string),
		// fields appear to only be relevant for jamf admin indexed packages
		// which i believe is to be deprecated.
		//RequiredProcessor:          d.Get("required_processor").(string),
		//SwitchWithPackage:          d.Get("switch_with_package").(string),
		//ReinstallOption:            d.Get("reinstall_option").(string),
		//TriggeringFiles:            d.Get("triggering_files").(string),
		InstallIfReportedAvailable: d.Get("install_if_reported_available").(bool),
		SendNotification:           d.Get("send_notification").(bool),
	}

	return packageResource, nil
}
