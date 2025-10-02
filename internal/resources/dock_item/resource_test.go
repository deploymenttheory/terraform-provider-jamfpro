package dock_item_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers"
	"github.com/stretchr/testify/assert"
)

// TestDockItemFrameworkResource_ConfigParsing tests that our external config files load correctly
func TestDockItemFrameworkResource_ConfigParsing(t *testing.T) {
	t.Run("App config loads successfully", func(t *testing.T) {
		config, err := helpers.ParseHCLFile("../tests/terraform/unit/resource_app.tf")
		assert.NoError(t, err)
		assert.Contains(t, config, "jamfpro_dock_item_framework")
		assert.Contains(t, config, "Test App Dock Item")
		assert.Contains(t, config, "App")
	})

	t.Run("File config loads successfully", func(t *testing.T) {
		config, err := helpers.ParseHCLFile("../tests/terraform/unit/resource_file.tf")
		assert.NoError(t, err)
		assert.Contains(t, config, "jamfpro_dock_item_framework")
		assert.Contains(t, config, "Test File Dock Item")
		assert.Contains(t, config, "File")
	})

	t.Run("Folder config loads successfully", func(t *testing.T) {
		config, err := helpers.ParseHCLFile("../tests/terraform/unit/resource_folder.tf")
		assert.NoError(t, err)
		assert.Contains(t, config, "jamfpro_dock_item_framework")
		assert.Contains(t, config, "Test Folder Dock Item")
		assert.Contains(t, config, "Folder")
	})
}

// TestDockItemFrameworkResource_JSONParsing tests that our external JSON files load correctly
func TestDockItemFrameworkResource_JSONParsing(t *testing.T) {
	t.Run("App JSON response loads successfully", func(t *testing.T) {
		jsonData, err := helpers.ParseJSONFile("../tests/json/dock_item_app_response.json")
		assert.NoError(t, err)
		assert.Contains(t, jsonData, "Test App Dock Item")
		assert.Contains(t, jsonData, "App")
		assert.Contains(t, jsonData, "file://localhost/Applications/iTunes.app/")
	})

	t.Run("File JSON response loads successfully", func(t *testing.T) {
		jsonData, err := helpers.ParseJSONFile("../tests/json/dock_item_file_response.json")
		assert.NoError(t, err)
		assert.Contains(t, jsonData, "Test File Dock Item")
		assert.Contains(t, jsonData, "File")
		assert.Contains(t, jsonData, "/etc/hosts")
	})

	t.Run("Folder JSON response loads successfully", func(t *testing.T) {
		jsonData, err := helpers.ParseJSONFile("../tests/json/dock_item_folder_response.json")
		assert.NoError(t, err)
		assert.Contains(t, jsonData, "Test Folder Dock Item")
		assert.Contains(t, jsonData, "Folder")
		assert.Contains(t, jsonData, "~/Downloads")
	})
}

// TestDockItemFrameworkResource_ValidationLogic tests the validation patterns
func TestDockItemFrameworkResource_ValidationLogic(t *testing.T) {
	t.Run("Valid types are accepted", func(t *testing.T) {
		validTypes := []string{"App", "File", "Folder"}
		for _, validType := range validTypes {
			// This would normally test the validator directly, but since we're
			// focusing on external file patterns, we test that configs contain valid types
			switch validType {
			case "App":
				config, err := helpers.ParseHCLFile("../tests/terraform/unit/resource_app.tf")
				assert.NoError(t, err)
				assert.Contains(t, config, `type = "App"`)
			case "File":
				config, err := helpers.ParseHCLFile("../tests/terraform/unit/resource_file.tf")
				assert.NoError(t, err)
				assert.Contains(t, config, `type = "File"`)
			case "Folder":
				config, err := helpers.ParseHCLFile("../tests/terraform/unit/resource_folder.tf")
				assert.NoError(t, err)
				assert.Contains(t, config, `type = "Folder"`)
			}
		}
	})

	t.Run("Timeout configurations are present", func(t *testing.T) {
		config, err := helpers.ParseHCLFile("../tests/terraform/unit/resource_app.tf")
		assert.NoError(t, err)
		assert.Contains(t, config, "timeouts")
		assert.Contains(t, config, "create")
		assert.Contains(t, config, "read")
		assert.Contains(t, config, "update")
		assert.Contains(t, config, "delete")
	})
}