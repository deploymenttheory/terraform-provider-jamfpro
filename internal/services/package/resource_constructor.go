// packages_constructor.go
package packages

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/files"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProPackageCreate constructs a ResourcePackage object from the provided schema data.
// It extracts the filename from the full path provided in the schema and uses it for the FileName field.
// If the full path is a URL, it downloads the file and uses the downloaded file path.
// The function returns the constructed ResourcePackage, the local file path, and an error if any.
func construct(d *schema.ResourceData) (*jamfpro.ResourcePackage, string, error) {
	fullPath := d.Get("package_file_source").(string)
	var fileName string
	var localFilePath string
	var err error

	if fullPath != "" {
		if strings.HasPrefix(fullPath, "http") {
			log.Printf("[INFO] URL detected: %s. Attempting to download.", fullPath)
			localFilePath, err = files.DownloadFile(fullPath)
			if err != nil {
				return nil, "", fmt.Errorf("failed to download file: %v", err)
			}
			fileName = filepath.Base(localFilePath)
			log.Printf("[INFO] Successfully downloaded file from URL: %s", fullPath)
		} else {
			fileName = filepath.Base(fullPath)
			localFilePath = fullPath
		}
	} else {
		fileName = d.Get("filename").(string)
		log.Printf("[INFO] No package_file_source specified, creating metadata-only package with filename: %s", fileName)
	}

	// Construct the ResourcePackage struct from the Terraform schema data
	resource := &jamfpro.ResourcePackage{
		PackageName:          d.Get("package_name").(string),
		FileName:             fileName,
		CategoryID:           d.Get("category_id").(string),
		Info:                 d.Get("info").(string),
		Notes:                d.Get("notes").(string),
		Priority:             d.Get("priority").(int),
		OSRequirements:       d.Get("os_requirements").(string),
		FillUserTemplate:     jamfpro.BoolPtr(d.Get("fill_user_template").(bool)),
		Indexed:              jamfpro.BoolPtr(d.Get("indexed").(bool)),
		FillExistingUsers:    jamfpro.BoolPtr(d.Get("fill_existing_users").(bool)),
		SWU:                  jamfpro.BoolPtr(d.Get("swu").(bool)),
		RebootRequired:       jamfpro.BoolPtr(d.Get("reboot_required").(bool)),
		SelfHealNotify:       jamfpro.BoolPtr(d.Get("self_heal_notify").(bool)),
		SelfHealingAction:    d.Get("self_healing_action").(string),
		OSInstall:            jamfpro.BoolPtr(d.Get("os_install").(bool)),
		SerialNumber:         d.Get("serial_number").(string),
		ParentPackageID:      d.Get("parent_package_id").(string),
		BasePath:             d.Get("base_path").(string),
		SuppressUpdates:      jamfpro.BoolPtr(d.Get("suppress_updates").(bool)),
		IgnoreConflicts:      jamfpro.BoolPtr(d.Get("ignore_conflicts").(bool)),
		SuppressFromDock:     jamfpro.BoolPtr(d.Get("suppress_from_dock").(bool)),
		SuppressEula:         jamfpro.BoolPtr(d.Get("suppress_eula").(bool)),
		SuppressRegistration: jamfpro.BoolPtr(d.Get("suppress_registration").(bool)),
		InstallLanguage:      d.Get("install_language").(string),
		OSInstallerVersion:   d.Get("os_installer_version").(string),
		Manifest:             d.Get("manifest").(string),
		ManifestFileName:     d.Get("manifest_file_name").(string),
	}

	// When no file source is provided, populate hash fields from user-supplied schema values
	if fullPath == "" {
		resource.MD5 = d.Get("md5").(string)
		resource.SHA256 = d.Get("sha256").(string)
		resource.SHA3512 = d.Get("sha3512").(string)
		resource.HashType = d.Get("hash_type").(string)
		resource.HashValue = d.Get("hash_value").(string)
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal Jamf Pro Package '%s' to JSON: %v", resource.FileName, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Package JSON:\n%s\n", string(resourceJSON))

	return resource, localFilePath, nil
}
