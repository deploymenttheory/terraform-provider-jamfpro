package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

func main() {
	cwd, _ := os.Getwd()
	log.Printf("Current working directory: %s", cwd)

	client := initializeJamfProClient()

	departmentIDs, err := client.GetDepartments()
	if err != nil {
		log.Fatalf("Failed to retrieve department IDs: %v", err)
	}

	generateTerraformConfig(client, departmentIDs)
}

func initializeJamfProClient() *jamfpro.Client {
	// Load client authentication configuration from a JSON file
	authConfig, err := jamfpro.LoadClientAuthConfig("/Users/dafyddwatkins/GitHub/deploymenttheory/terraform-provider-jamfpro/.localtesting/clientauth.json")
	if err != nil {
		log.Fatalf("Failed to load client authentication configuration: %v", err)
	}

	// Construct the jamfpro.Config object
	config := jamfpro.Config{
		InstanceName: authConfig.InstanceName,
		//DebugMode:    true,
		Logger:       nil,
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
	}

	// Initialize the Jamf Pro client
	client, err := jamfpro.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to initialize Jamf Pro client: %v", err)
	}

	return client
}

func generateTerraformConfig(client *jamfpro.Client, departmentIDs *jamfpro.ResponseDepartmentsList) {
	log.Println("Starting generation of Terraform config...")

	// Check if departments.tf already exists
	filePath := "departments.tf"
	var f *os.File

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		var err error
		f, err = os.Create(filePath)
		if err != nil {
			log.Fatalf("Failed to create Terraform file: %v", err)
		}
	} else {
		// Open the file in append mode if it exists
		var err error
		f, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open Terraform file: %v", err)
		}
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("Failed to close Terraform file: %v", err)
		}
	}()

	for _, departmentID := range departmentIDs.Results {
		log.Printf("Processing department with ID: %d", departmentID.Id)

		department, err := client.GetDepartmentByID(departmentID.Id)
		if err != nil {
			log.Printf("Failed to retrieve department with ID %d: %v", departmentID.Id, err)
			continue
		}

		hcl := generateDepartmentHCL(*department)
		_, err = f.WriteString(hcl)
		if err != nil {
			log.Printf("Failed to write department (ID: %d, Name: %s) to Terraform file: %v", department.ID, department.Name, err)
			continue
		}

		log.Printf("Successfully added department (ID: %d, Name: %s) to Terraform config.", department.ID, department.Name)

		err = importIntoTerraformState(department.ID)
		if err != nil {
			log.Printf("Failed to import department (ID: %d, Name: %s) into Terraform state: %v", department.ID, department.Name, err)
		} else {
			log.Printf("Successfully imported department (ID: %d, Name: %s) into Terraform state.", department.ID, department.Name)
		}
	}

	log.Println("Finished generation of Terraform config.")
}

func generateDepartmentHCL(department jamfpro.ResponseDepartment) string {
	// Generate Terraform HCL for the given department
	hcl := fmt.Sprintf(`
resource "jamfpro_department" "department_%d" {
  name = "%s"
}
`, department.ID, department.Name)

	return hcl
}

// isResourceImported checks if a given resource is already imported in the Terraform state.
// It returns true if the resource is imported, false otherwise.
func isResourceImported(resourceAddress string) (bool, error) {
	cmd := exec.Command("terraform", "state", "list", resourceAddress)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to execute terraform state list: %v, output: %s", err, output)
	}

	// If the output contains the resource address, it means the resource is already imported.
	if strings.Contains(string(output), resourceAddress) {
		return true, nil
	}

	return false, nil
}

// importIntoTerraformState imports the department into the Terraform state if not already imported.
func importIntoTerraformState(departmentID int) error {
	resourceAddress := fmt.Sprintf("jamfpro_department.department_%d", departmentID)

	// Check if the resource is already imported
	imported, err := isResourceImported(resourceAddress)
	if err != nil {
		return fmt.Errorf("failed to check if resource is imported: %v", err)
	}

	if imported {
		log.Printf("Resource %s is already imported into Terraform state.", resourceAddress)
		return nil
	}

	// Import the resource
	cmd := exec.Command("terraform", "import", resourceAddress, fmt.Sprint(departmentID))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run terraform import: %v, output: %s", err, output)
	}

	log.Printf("Successfully imported resource %s into Terraform state.", resourceAddress)
	return nil
}
