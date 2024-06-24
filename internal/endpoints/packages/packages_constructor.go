// packages_constructor.go
package packages

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProPackageCreate constructs a ResourcePackage object from the provided schema data.
// It extracts the filename from the full path provided in the schema and uses it for the FileName field.
// constructJamfProPackageCreate constructs a ResourcePackage object from the provided schema data.
// It extracts the filename from the full path provided in the schema and uses it for the FileName field.
func constructJamfProPackageCreate(d *schema.ResourceData) (*jamfpro.ResourcePackage, error) {
	// Use filepath.Base to extract just the filename from the full path
	fullPath := d.Get("package_file_path").(string)
	fileName := filepath.Base(fullPath)

	// Construct the ResourcePackage struct from the Terraform schema data
	resource := &jamfpro.ResourcePackage{
		PackageName:          d.Get("package_name").(string),
		FileName:             fileName,
		CategoryID:           d.Get("category_id").(string),
		Info:                 d.Get("info").(string),
		Notes:                d.Get("notes").(string),
		Priority:             d.Get("priority").(int),
		OSRequirements:       d.Get("os_requirements").(string),
		FillUserTemplate:     BoolPtr(d.Get("fill_user_template").(bool)),
		Indexed:              BoolPtr(d.Get("indexed").(bool)),
		FillExistingUsers:    BoolPtr(d.Get("fill_existing_users").(bool)),
		SWU:                  BoolPtr(d.Get("swu").(bool)),
		RebootRequired:       BoolPtr(d.Get("reboot_required").(bool)),
		SelfHealNotify:       BoolPtr(d.Get("self_heal_notify").(bool)),
		SelfHealingAction:    d.Get("self_healing_action").(string),
		OSInstall:            BoolPtr(d.Get("os_install").(bool)),
		SerialNumber:         d.Get("serial_number").(string),
		ParentPackageID:      d.Get("parent_package_id").(string),
		BasePath:             d.Get("base_path").(string),
		SuppressUpdates:      BoolPtr(d.Get("suppress_updates").(bool)),
		IgnoreConflicts:      BoolPtr(d.Get("ignore_conflicts").(bool)),
		SuppressFromDock:     BoolPtr(d.Get("suppress_from_dock").(bool)),
		SuppressEula:         BoolPtr(d.Get("suppress_eula").(bool)),
		SuppressRegistration: BoolPtr(d.Get("suppress_registration").(bool)),
		InstallLanguage:      d.Get("install_language").(string),
		OSInstallerVersion:   d.Get("os_installer_version").(string),
		Manifest:             d.Get("manifest").(string),
		ManifestFileName:     d.Get("manifest_file_name").(string),
	}

	// Serialize and pretty-print the Network Segment object as JSON for logging
	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Package '%s' to JSON: %v", resource.FileName, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Package JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}

// BoolPtr is a helper function to create a pointer to a bool.
func BoolPtr(b bool) *bool {
	return &b
}
