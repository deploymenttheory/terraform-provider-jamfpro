package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

const (
	// Define ANSI color codes
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

// TerraformPlan represents the top-level structure of a Terraform plan in JSON format.
type TerraformPlan struct {
	FormatVersion    string              `json:"format_version"`
	TerraformVersion string              `json:"terraform_version"`
	Variables        map[string]Variable `json:"variables"`
	PlannedValues    PlannedValues       `json:"planned_values"`
	ResourceChanges  []ResourceChange    `json:"resource_changes"`
	Configuration    Configuration       `json:"configuration"`
	Timestamp        string              `json:"timestamp"`
	Errored          bool                `json:"errored"`
}

// Variable defines the structure for Terraform variables in the plan.
type Variable struct {
	Value       interface{} `json:"value"`
	Description string      `json:"description,omitempty"`
	Sensitive   bool        `json:"sensitive,omitempty"`
}

// PlannedValues encapsulates all planned values including outputs and root module resources.
type PlannedValues struct {
	Outputs    map[string]Output `json:"outputs,omitempty"` // Added Outputs field
	RootModule RootModule        `json:"root_module"`
}

// Output defines an output as defined in the Terraform plan.
type Output struct {
	Sensitive bool `json:"sensitive"`
}

// RootModule represents the root module of the Terraform configuration and contains resources.
type RootModule struct {
	Resources []Resource `json:"resources"`
}

// Resource describes a single Terraform resource with its properties.
type Resource struct {
	Address       string                 `json:"address"`
	Mode          string                 `json:"mode"`
	Type          string                 `json:"type"`
	Name          string                 `json:"name"`
	ProviderName  string                 `json:"provider_name"`
	SchemaVersion int                    `json:"schema_version"`
	Values        map[string]interface{} `json:"values"`
}

// Resource Name Value from tf plan
// Duplicate resource names are checked based on the name field
type ResourceValues struct {
	Name string `json:"name"`
}

// ResourceChange details changes to a specific resource including before and after states.
type ResourceChange struct {
	Address      string `json:"address"`
	Mode         string `json:"mode"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	ProviderName string `json:"provider_name"`
	Change       Change `json:"change"`
}

// Change captures the differences for a resource between its current state and the planned state.
type Change struct {
	Actions         []string               `json:"actions"`
	Before          map[string]interface{} `json:"before"`
	After           map[string]interface{} `json:"after"`
	AfterUnknown    map[string]interface{} `json:"after_unknown"`
	BeforeSensitive SensitiveType          `json:"before_sensitive,omitempty"`
	AfterSensitive  SensitiveType          `json:"after_sensitive"`
}

// SensitiveType is used to handle Terraform's sensitive values which can be either a boolean flag or a complex structure.
// This struct provides flexibility by allowing sensitive values to be parsed correctly regardless of their underlying type.
type SensitiveType struct {
	BoolValue   bool
	MapValue    map[string]interface{}
	IsBool      bool
	IsPopulated bool
}

// Configuration

// Configuration represents the Terraform configuration including provider configs and root module configurations.
type Configuration struct {
	ProviderConfig map[string]ProviderConfig `json:"provider_config"`
	RootModule     RootModuleConfig          `json:"root_module"`
}

// ProviderConfig defines a Terraform provider configuration including expressions used within the provider block.
type ProviderConfig struct {
	Name              string              `json:"name"`
	FullName          string              `json:"full_name"`
	VersionConstraint string              `json:"version_constraint"`
	Expressions       ProviderExpressions `json:"expressions"`
}

// ProviderExpressions captures expressions related to a Terraform provider configuration.
type ProviderExpressions struct {
	ClientID     Expression `json:"client_id"`
	ClientSecret Expression `json:"client_secret"`
	InstanceName Expression `json:"instance_name"`
	LogLevel     Expression `json:"log_level"`
}

// Expression represents a Terraform expression, which could be a constant value or a set of references.
type Expression struct {
	ConstantValue string   `json:"constant_value,omitempty"`
	References    []string `json:"references,omitempty"`
}

// RootModuleConfig holds the configurations for resources and variables within the root module of the Terraform configuration.
type RootModuleConfig struct {
	Resources []ResourceConfig          `json:"resources"`
	Variables map[string]VariableConfig `json:"variables"`
}

// ResourceConfig describes the configuration for a single Terraform resource within a module.
type ResourceConfig struct {
	Address           string      `json:"address"`
	Mode              string      `json:"mode"`
	Type              string      `json:"type"`
	Name              string      `json:"name"`
	ProviderConfigKey string      `json:"provider_config_key"`
	Expressions       Expressions `json:"expressions"`
	SchemaVersion     int         `json:"schema_version"`
}

type Expressions struct {
	Name Expression `json:"name"`
}

// VariableConfig defines the configuration for a Terraform variable, including its default value and other attributes.
type VariableConfig struct {
	Default     interface{} `json:"default"`
	Description string      `json:"description,omitempty"`
	Sensitive   bool        `json:"sensitive,omitempty"`
}

// main function parses command line arguments to locate the Terraform plan file, unmarshals it, and checks for duplicate resource names.
func main() {
	// Define a string flag for the Terraform plan file path
	tfPlanPath := flag.String("tfplan", "", "Path to the Terraform plan file in JSON format")

	// Parse the command-line flags
	flag.Parse()

	// Check if the tfplan flag has been set
	if *tfPlanPath == "" {
		fmt.Println("Usage: -tfplan <path to terraform plan json>")
		return
	}

	// Read the Terraform plan from the file using os.ReadFile
	planFile, err := os.ReadFile(*tfPlanPath)
	if err != nil {
		fmt.Printf("Error reading plan file: %v\n", err)
		return
	}

	var plan TerraformPlan
	err = json.Unmarshal(planFile, &plan)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return
	}

	// Specified the resource types to validate duplicates for
	interestedResourceTypes := map[string]bool{
		"jamfpro_account":                       true,
		"jamfpro_account_group":                 true,
		"jamfpro_advanced_computer_search":      true,
		"jamfpro_advanced_mobile_device_search": true,
		"jamfpro_advanced_user_search":          true,
		"jamfpro_allowed_file_extension":        true,
		"jamfpro_api_integration":               true,
		"jamfpro_api_role":                      true,
		"jamfpro_building":                      true,
		"jamfpro_category":                      true,
		"jamfpro_computer_checkin":              true,
		"jamfpro_computer_extension_attribute":  true,
		"jamfpro_computer_group":                true,
		"jamfpro_computer_prestage":             true,
		"jamfpro_department":                    true,
		"jamfpro_disk_encryption_configuration": true,
		"jamfpro_dock_item":                     true,
		"jamfpro_file_share_distribution_point": true,
		"jamfpro_site":                          true,
		"jamfpro_script":                        true,
		"jamfpro_network_segment":               true,
		"jamfpro_package":                       true,
		"jamfpro_policy":                        true,
		"jamfpro_printer":                       true,
	}

	// Store resource names and their occurrences
	resourceNames := make(map[string]int)

	// Iterate over resources in the plan
	for _, resource := range plan.PlannedValues.RootModule.Resources {
		if _, ok := interestedResourceTypes[resource.Type]; ok {
			// Attempt to extract the name from the Values map
			if name, exists := resource.Values["name"]; exists && name != nil {
				if nameStr, ok := name.(string); ok {
					resourceNames[nameStr]++
				}
			}
		}
	}

	// Check for duplicates
	foundDuplicates := false
	for name, count := range resourceNames {
		if count > 1 {
			errorMessage := fmt.Sprintf("Error: Duplicate Jamf Pro resource name found: %s, Count: %d", name, count)
			printColor(errorMessage, colorRed)
			foundDuplicates = true
		}
	}

	if !foundDuplicates {
		printColor("Check completed: No duplicate Jamf Pro resource names found within the specified Terraform plan.", colorGreen)
	}
}

// printColor prints a message with the specified color to the console.
func printColor(message string, colorCode string) {
	fmt.Println(colorCode, message, colorReset)
}

// UnmarshalJSON is a custom unmarshaler for SensitiveType that handles both boolean and structured sensitive values.
// It first attempts to unmarshal the data into a boolean; if that fails, it then tries to unmarshal into a map.
// This allows for correctly parsing the "before_sensitive" and "after_sensitive" fields which can vary in type.
func (s *SensitiveType) UnmarshalJSON(data []byte) error {
	s.IsPopulated = true // Mark as populated for further logic if needed

	// First, try to unmarshal as a bool
	if err := json.Unmarshal(data, &s.BoolValue); err == nil {
		s.IsBool = true
		return nil
	}

	// If unmarshalling as a bool fails, try as a map
	if err := json.Unmarshal(data, &s.MapValue); err != nil {
		return err // Return error if it fails to unmarshal as both bool and map
	}

	// Successfully unmarshalled as a map
	s.IsBool = false
	return nil
}
