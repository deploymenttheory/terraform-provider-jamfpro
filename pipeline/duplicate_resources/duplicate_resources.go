package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// Define the top-level TerraformPlan struct
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

// Variables

type Variable struct {
	Value       interface{} `json:"value"`
	Description string      `json:"description,omitempty"`
	Sensitive   bool        `json:"sensitive,omitempty"`
}

// Planned Values

type PlannedValues struct {
	RootModule RootModule `json:"root_module"`
}

// Define a struct for the RootModule part
type RootModule struct {
	Resources []Resource `json:"resources"`
}

// Define a struct for each Resource
type Resource struct {
	Address string         `json:"address"`
	Type    string         `json:"type"`
	Values  ResourceValues `json:"values"`
}

// Define a struct for the Values part
type ResourceValues struct {
	Name string `json:"name"`
}

// Resource Change

type ResourceChange struct {
	Address  string `json:"address"`
	Mode     string `json:"mode"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Provider string `json:"provider_name"`
	Change   Change `json:"change"`
}

type Change struct {
	Actions         []string               `json:"actions"`
	Before          map[string]interface{} `json:"before"`
	After           map[string]interface{} `json:"after"`
	AfterUnknown    map[string]interface{} `json:"after_unknown"`
	BeforeSensitive bool                   `json:"before_sensitive"`
	AfterSensitive  map[string]interface{} `json:"after_sensitive"`
}

// Configuration

type Configuration struct {
	ProviderConfig map[string]ProviderConfig `json:"provider_config"`
	RootModule     RootModuleConfig          `json:"root_module"`
}

type ProviderConfig struct {
	Name              string              `json:"name"`
	FullName          string              `json:"full_name"`
	VersionConstraint string              `json:"version_constraint"`
	Expressions       ProviderExpressions `json:"expressions"`
}

type ProviderExpressions struct {
	ClientID     Expression `json:"client_id"`
	ClientSecret Expression `json:"client_secret"`
	InstanceName Expression `json:"instance_name"`
	LogLevel     Expression `json:"log_level"`
}

type Expression struct {
	ConstantValue string   `json:"constant_value,omitempty"`
	References    []string `json:"references,omitempty"`
}

type RootModuleConfig struct {
	Resources []ResourceConfig          `json:"resources"`
	Variables map[string]VariableConfig `json:"variables"`
}

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

type VariableConfig struct {
	Default     interface{} `json:"default"`
	Description string      `json:"description,omitempty"`
	Sensitive   bool        `json:"sensitive,omitempty"`
}

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
			resourceNames[resource.Values.Name]++
		}
	}

	// Check for duplicates
	foundDuplicates := false
	for name, count := range resourceNames {
		if count > 1 {
			fmt.Printf("Error: Duplicate resource name found: %s, Count: %d\n", name, count)
			foundDuplicates = true
		}
	}

	if !foundDuplicates {
		fmt.Println("No duplicates found.")
	}
}
