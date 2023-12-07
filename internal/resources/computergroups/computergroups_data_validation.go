// computergroup_data_validation.go
package computergroups

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Define a slice containing all the valid built-in smart group criteria names.
var validBuiltInSmartGroupCriteriaNames = []string{
	"Activation Lock Enabled",
	"Activation Lock Manageable",
	"Active Directory Status",
	"Apple Silicon",
	"AppleCare ID",
	"Application Bundle ID",
	"Application Title",
	"Application Version",
	"Architecture Type",
	"Asset Tag",
	"Available RAM Slots",
	"Available SWUs",
	"Bar Code",
	"Battery Capacity",
	"BeyondCorp Enterprise Integration - Compliance Status",
	"BeyondCorp Enterprise Integration - Registration Status",
	"Bluetooth Low Energy Capability",
	"Boot Drive Available MB",
	"Boot Drive Percentage Full",
	"Boot ROM",
	"Bootstrap Token Allowed",
	"Building",
	"Bus Speed MHz",
	"Cached Packages",
	"Certificate Issuer",
	"Certificate Name",
	"Certificates Expiring",
	"Computer Azure Active Directory ID",
	"Computer Group",
	"Computer Name",
	"Conditional Access Inventory State",
	"Content Caching - Activated",
	"Content Caching - Active",
	"Content Caching - Actual Cache Used bytes",
	"Content Caching - Cache Limit bytes",
	"Content Caching - Cache Status",
	"Content Caching - Tetherator Status",
	"Core Storage Partition Scheme on Boot Partition",
	"Declarative Device Management Enabled",
	"Department",
	"Device Compliance Integration - Compliance Status",
	"Device Compliance Integration - Registration Status",
	"Disable Automatic Login",
	"Disk Encryption Configuration",
	"Drive Capacity MB",
	"Email Address",
	"Enrolled via Automated Device Enrollment",
	"Enrollment Method: PreStage enrollment",
	"External Boot Level",
	"FileVault 2 Eligibility",
	"FileVault 2 Individual Key Validation",
	"FileVault 2 Institutional Key",
	"FileVault 2 Partition Encryption State",
	"FileVault 2 Recovery Key Type",
	"FileVault 2 Status",
	"FileVault 2 User",
	"FileVault Status",
	"Firewall Enabled",
	"Font Title",
	"Font Version",
	"Full Name",
	"Gatekeeper",
	"IP Address",
	"iTunes Store Account",
	"JAMF Binary Version",
	"JSS Computer ID",
	"Last Check-in",
	"Last Enrollment",
	"Last iCloud Backup",
	"Last Inventory Update",
	"Last Reported IP Address",
	"Lease Expiration",
	"Licensed Software",
	"Life Expectancy",
	"Local User Accounts",
	"MAC Address",
	"Make",
	"Managed By",
	"Mapped Printers",
	"Master Password Set",
	"Maximum Passcode Age",
	"MDM Capability",
	"MDM Profile Expiration Date",
	"MDM Profile Renewal Needed - CA Renewed",
	"Minimum Number of Complex Characters",
	"Model",
	"Model Identifier",
	"NIC Speed",
	"Number of Available Updates",
	"Number of Processors",
	"Operating System",
	"Operating System Build",
	"Operating System Name",
	"Operating System Rapid Security Response",
	"Operating System Version",
	"Optical Drive",
	"Packages Installed By Casper",
	"Packages Installed By Installer.app/SWU",
	"Partition Name",
	"Password History",
	"Password Type",
	"Patch Reporting Software Title",
	"Phone Number",
	"Platform",
	"Plug-in Title",
	"Plug-in Version",
	"PO Date",
	"PO Number",
	"Position",
	"Processor Speed MHz",
	"Processor Type",
	"Profile Identifier",
	"Profile Name",
	"Purchase Price",
	"Purchased or Leased",
	"Purchasing Account",
	"Purchasing Contact",
	"Recovery Lock Enabled",
	"Remote Desktop Enabled",
	"Required Passcode Length",
	"Room",
	"Running Services",
	"S.M.A.R.T. Status",
	"Scheduled Tasks",
	"Secure Boot Level",
	"Serial Number",
	"Service Pack",
	"SMC Version",
	"Software Update Device ID",
	"Supervised",
	"Supports iOS and iPadOS App Installations",
	"System Integrity Protection",
	"Total Number of Cores",
	"Total RAM MB",
	"UDID",
	"User Approved MDM",
	"User Azure Active Directory ID",
	"Username",
	"Vendor",
	"Warranty Expiration",
	"XProtect Definitions Version",
}

func getComputerExtensionAttributeNames(ctx context.Context) ([]string, error) {
	// Extract the client from the context
	client, ok := ctx.Value("client").(*jamfpro.Client)
	if !ok {
		return nil, fmt.Errorf("could not retrieve Jamf Pro client from context")
	}

	// Call the SDK function to Get the list of computer extension attributes using the client from the context.
	response, err := client.GetComputerExtensionAttributes()
	if err != nil {
		return nil, fmt.Errorf("error fetching computer extension attributes: %v", err)
	}

	// Parse the response and extract the names of the extension attributes.
	var customNames []string
	for _, attr := range response.Results {
		if attr.Enabled {
			customNames = append(customNames, attr.Name)
		}
	}
	return customNames, nil
}

// validateSmartGroupCriteriaName uses the provider's context to validate the criteria name.
func validateSmartGroupCriteriaName(ctx context.Context, val interface{}, key string) (warns []string, errs []error) {
	name, ok := val.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected type of %s to be string", key))
		return
	}

	// Check if the provided name is in the list of valid built-in names.
	for _, validName := range validBuiltInSmartGroupCriteriaNames {
		if name == validName {
			return // The name is valid as a built-in criteria.
		}
	}

	// If not a built-in criteria, check if it's a valid custom extension attribute name.
	customNames, err := getComputerExtensionAttributeNames(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("error validating %s: %v", key, err))
		return
	}

	for _, customName := range customNames {
		if name == customName {
			return // The name is valid as a custom extension attribute.
		}
	}

	// If not found in either, return an error.
	errs = append(errs, fmt.Errorf("%s must be either a built-in criteria or a custom computer extension attribute", key))
	return
}

// customDiffComputeGroups is a CustomDiff function that enforces conditional logic on the 'computers' and 'criteria' fields of the JamfProComputerGroups resource based on the value of 'is_smart'.
// When is_smart is true, the criteria block is valid, and the computers block should not be set.
// When is_smart is false, the computers block is valid, and the criteria block should not be set.
func customDiffComputeGroups(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	isSmart := diff.Get("is_smart").(bool)

	// When 'is_smart' is true, 'computers' should not be set.
	if isSmart {
		if computers, exists := diff.GetOk("computers"); exists && len(computers.([]interface{})) > 0 {
			return fmt.Errorf("'computers' field is not allowed when 'is_smart' is true")
		}
	} else {
		// If 'is_smart' is false, 'criteria' should not be set.
		if criteria, exists := diff.GetOk("criteria"); exists && len(criteria.([]interface{})) > 0 {
			return fmt.Errorf("'criteria' field is not allowed when 'is_smart' is false")
		}
	}

	// Additional validations for 'criteria' when 'is_smart' is true.
	if isSmart {
		criteria, ok := diff.GetOk("criteria")
		if !ok || len(criteria.([]interface{})) == 0 {
			return fmt.Errorf("'criteria' field must be set when 'is_smart' is true")
		}

		for i, c := range criteria.([]interface{}) {
			criterion, ok := c.(map[string]interface{})
			if !ok {
				continue // Skip invalid structure.
			}

			// Validate 'name', 'and_or', and 'search_type' in each criterion.
			if criterion["name"] == nil || criterion["name"].(string) == "" {
				return fmt.Errorf("'name' field is required for 'criteria' at index %d when 'is_smart' is true", i)
			}
			if criterion["and_or"] == nil || criterion["and_or"].(string) == "" {
				return fmt.Errorf("'and_or' field is required for 'criteria' at index %d when 'is_smart' is true", i)
			}
			if criterion["search_type"] == nil || criterion["search_type"].(string) == "" {
				return fmt.Errorf("'search_type' field is required for 'criteria' at index %d when 'is_smart' is true", i)
			}
		}
	}

	return nil
}
