// packages_data_object.go
package packages

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
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

	packageResource := &jamfpro.ResourcePackage{
		Name:                       d.Get("name").(string),
		Filename:                   filename,
		Category:                   d.Get("category").(string),
		Info:                       d.Get("info").(string),
		Notes:                      d.Get("notes").(string),
		Priority:                   d.Get("priority").(int),
		RebootRequired:             d.Get("reboot_required").(bool),
		FillUserTemplate:           d.Get("fill_user_template").(bool),
		FillExistingUsers:          d.Get("fill_existing_users").(bool),
		BootVolumeRequired:         d.Get("boot_volume_required").(bool),
		AllowUninstalled:           d.Get("allow_uninstalled").(bool),
		OSRequirements:             d.Get("os_requirements").(string),
		RequiredProcessor:          d.Get("required_processor").(string),
		SwitchWithPackage:          d.Get("switch_with_package").(string),
		InstallIfReportedAvailable: d.Get("install_if_reported_available").(bool),
		ReinstallOption:            d.Get("reinstall_option").(string),
		TriggeringFiles:            d.Get("triggering_files").(string),
		SendNotification:           d.Get("send_notification").(bool),
	}

	return packageResource, nil
}

// generateFileHash accepts a file path and returns a SHA-256 hash of the file's contents.
func generateFileHash(filePath string) (string, error) {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	// Create a new SHA256 hash object
	hash := sha256.New()

	// Copy the file content into the hash object
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file contents of %s: %v", filePath, err)
	}

	// Compute the SHA256 checksum of the file
	hashBytes := hash.Sum(nil)

	// Convert the bytes to a hex string
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}
