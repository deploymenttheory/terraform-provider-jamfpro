package computer_prestage_enrollment

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

// TestSkipSetupItems_NilPointers tests that skipSetupItems handles nil pointers gracefully
// This addresses the bug where upgrading from v0.23.0 would crash when SoftwareUpdate and
// AdditionalPrivacySettings fields were nil in existing state.
func TestSkipSetupItems_NilPointers(t *testing.T) {
	// Simulate state from v0.23.0 where new fields don't exist (nil pointers)
	skipSetupItemsOldState := jamfpro.ComputerPrestageSubsetSkipSetupItems{
		Biometric:         boolPtr(true),
		TermsOfAddress:    boolPtr(false),
		FileVault:         boolPtr(false),
		ICloudDiagnostics: boolPtr(true),
		Diagnostics:       boolPtr(true),
		Accessibility:     boolPtr(false),
		AppleID:           boolPtr(true),
		ScreenTime:        boolPtr(true),
		Siri:              boolPtr(true),
		DisplayTone:       boolPtr(false),
		Restore:           boolPtr(true),
		Appearance:        boolPtr(true),
		Privacy:           boolPtr(true),
		Payment:           boolPtr(true),
		Registration:      boolPtr(true),
		TOS:               boolPtr(true),
		ICloudStorage:     boolPtr(true),
		Location:          boolPtr(true),
		Intelligence:      boolPtr(true),
		EnableLockdownMode: boolPtr(true),
		Welcome:           boolPtr(true),
		Wallpaper:         boolPtr(true),
		// SoftwareUpdate and AdditionalPrivacySettings are nil (as they would be in v0.23.0 state)
		SoftwareUpdate:            nil,
		AdditionalPrivacySettings: nil,
	}

	// This should not panic
	result := skipSetupItems(skipSetupItemsOldState)

	// Verify the new fields default to false when nil
	if softwareUpdate, ok := result["software_update"].(bool); !ok || softwareUpdate != false {
		t.Errorf("Expected software_update to be false when nil, got %v", result["software_update"])
	}

	if additionalPrivacy, ok := result["additional_privacy_settings"].(bool); !ok || additionalPrivacy != false {
		t.Errorf("Expected additional_privacy_settings to be false when nil, got %v", result["additional_privacy_settings"])
	}

	// Verify old fields still work correctly
	if biometric, ok := result["biometric"].(bool); !ok || biometric != true {
		t.Errorf("Expected biometric to be true, got %v", result["biometric"])
	}
}

// TestSkipSetupItems_WithNewFields tests that skipSetupItems works correctly with all fields present
func TestSkipSetupItems_WithNewFields(t *testing.T) {
	// Simulate state from v0.25.0+ where all fields exist
	skipSetupItemsNewState := jamfpro.ComputerPrestageSubsetSkipSetupItems{
		Biometric:                 boolPtr(true),
		TermsOfAddress:            boolPtr(false),
		FileVault:                 boolPtr(false),
		ICloudDiagnostics:         boolPtr(true),
		Diagnostics:               boolPtr(true),
		Accessibility:             boolPtr(false),
		AppleID:                   boolPtr(true),
		ScreenTime:                boolPtr(true),
		Siri:                      boolPtr(true),
		DisplayTone:               boolPtr(false),
		Restore:                   boolPtr(true),
		Appearance:                boolPtr(true),
		Privacy:                   boolPtr(true),
		Payment:                   boolPtr(true),
		Registration:              boolPtr(true),
		TOS:                       boolPtr(true),
		ICloudStorage:             boolPtr(true),
		Location:                  boolPtr(true),
		Intelligence:              boolPtr(true),
		EnableLockdownMode:        boolPtr(true),
		Welcome:                   boolPtr(true),
		Wallpaper:                 boolPtr(true),
		SoftwareUpdate:            boolPtr(true),
		AdditionalPrivacySettings: boolPtr(false),
	}

	result := skipSetupItems(skipSetupItemsNewState)

	// Verify the new fields work correctly when not nil
	if softwareUpdate, ok := result["software_update"].(bool); !ok || softwareUpdate != true {
		t.Errorf("Expected software_update to be true, got %v", result["software_update"])
	}

	if additionalPrivacy, ok := result["additional_privacy_settings"].(bool); !ok || additionalPrivacy != false {
		t.Errorf("Expected additional_privacy_settings to be false, got %v", result["additional_privacy_settings"])
	}
}

// boolPtr is a helper function to get a pointer to a bool value
func boolPtr(b bool) *bool {
	return &b
}

