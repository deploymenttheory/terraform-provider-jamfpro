// packages_constructor.go
package packages

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProPackageCreate constructs a ResourcePackage object from the provided schema data.
// It extracts the filename from the full path provided in the schema and uses it for the FileName field.
// If the full path is a URL, it downloads the file and uses the downloaded file path.
func constructJamfProPackageCreate(d *schema.ResourceData) (*jamfpro.ResourcePackage, error) {
	fullPath := d.Get("package_file_path").(string)
	var fileName string
	var err error

	if strings.HasPrefix(fullPath, "http") {
		log.Printf("[INFO] URL detected: %s. Attempting to download.", fullPath)
		fileName, err = downloadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to download file: %v", err)
		}
		log.Printf("[INFO] Successfully downloaded file from URL: %s", fullPath)
	} else {
		fileName = filepath.Base(fullPath)
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

// downloadFile downloads a file from the given URL and saves it to the specified path.
// It returns the local file path or an error if the download fails.
func downloadFile(url string) (string, error) {
	fileName := filepath.Base(url)

	out, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %v", fileName, err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file from %s: %v", url, err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to file %s: %v", fileName, err)
	}

	// Return the local file path
	return fileName, nil
}
