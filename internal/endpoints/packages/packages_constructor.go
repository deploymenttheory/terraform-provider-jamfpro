// packages_constructor.go
package packages

import (
	"log"
	"path/filepath"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProPackageCreate constructs a ResourcePackage object from the provided schema data.
// It extracts the filename from the full path provided in the schema and uses it for the FileName field.
func constructJamfProPackageCreate(d *schema.ResourceData) (*jamfpro.ResourcePackage, error) {
	// Use filepath.Base to extract just the filename from the full path
	fullPath := d.Get("package_file_path").(string)
	fileName := filepath.Base(fullPath)

	// Get the category from the schema, and set it to "-1" if it's empty
	// 'Unknown' is the valid default request value for the category field when none is set
	// Jamf API returns "No category assigned" for the same field. But this is not a valid
	// request value. Why!!!! >_<
	category := d.Get("category_id").(string)
	if category == "" {
		category = "-1"
	}

	// Construct the ResourcePackage struct from the Terraform schema data
	packageResource := &jamfpro.ResourcePackage{
		PackageName:          d.Get("package_name").(string),
		FileName:             fileName,
		CategoryID:           category,
		Priority:             d.Get("priority").(int),
		FillUserTemplate:     BoolPtr(d.Get("fill_user_template").(bool)),
		SWU:                  BoolPtr(d.Get("swu").(bool)),
		RebootRequired:       BoolPtr(d.Get("reboot_required").(bool)),
		OSInstall:            BoolPtr(d.Get("os_install").(bool)),
		SuppressUpdates:      BoolPtr(d.Get("suppress_updates").(bool)),
		SuppressFromDock:     BoolPtr(d.Get("suppress_from_dock").(bool)),
		SuppressEula:         BoolPtr(d.Get("suppress_eula").(bool)),
		SuppressRegistration: BoolPtr(d.Get("suppress_registration").(bool)),
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Package: %+v\n", packageResource)

	return packageResource, nil
}

// BoolPtr is a helper function to create a pointer to a bool.
func BoolPtr(b bool) *bool {
	return &b
}
