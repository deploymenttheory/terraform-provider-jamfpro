// department_resource_test.go
/*
Unit Tests: These test individual functions or methods without making real-world side effects.
Mocks are frequently used in this context to simulate interactions with external systems (like the SDK or an API).
These tests are faster and can be run without setting up any real infrastructure.

Acceptance Tests: These tests interact with the real-world resources, making actual API calls to create, read, update, or delete resources.
They are used to ensure that the Terraform provider is correctly interfacing with the real service.
*/
package provider

/*
import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"jamfpro": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

// MockAPIClient Jamf Pro Department methods to simulate the behavior of the actual jamfpro.Client
// without making real API calls
func (m *MockAPIClient) GetDepartments() (*jamfpro.ResponseDepartments, error) {
	args := m.Called()
	return args.Get(0).(*jamfpro.ResponseDepartments), args.Error(1)
}

func (m *MockAPIClient) GetDepartmentByID(id int) (*jamfpro.Department, error) {
	args := m.Called(id)
	return args.Get(0).(*jamfpro.Department), args.Error(1)
}

func (m *MockAPIClient) GetDepartmentByName(name string) (*jamfpro.Department, error) {
	args := m.Called(name)
	return args.Get(0).(*jamfpro.Department), args.Error(1)
}

func (m *MockAPIClient) GetDepartmentIdByName(name string) (int, error) {
	args := m.Called(name)
	return args.Int(0), args.Error(1)
}

func (m *MockAPIClient) CreateDepartment(departmentName string) (*jamfpro.Department, error) {
	args := m.Called(departmentName)
	return args.Get(0).(*jamfpro.Department), args.Error(1)
}

func (m *MockAPIClient) UpdateDepartmentByID(id int, departmentName string) (*jamfpro.Department, error) {
	args := m.Called(id, departmentName)
	return args.Get(0).(*jamfpro.Department), args.Error(1)
}

func (m *MockAPIClient) UpdateDepartmentByName(oldName string, newName string) (*jamfpro.Department, error) {
	args := m.Called(oldName, newName)
	return args.Get(0).(*jamfpro.Department), args.Error(1)
}

func (m *MockAPIClient) DeleteDepartmentByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAPIClient) DeleteDepartmentByName(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

// testAccPreCheck ensures necessary environment variables are set for acceptance tests.
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("JAMFPRO_INSTANCE"); v == "" {
		t.Fatal("JAMFPRO_INSTANCE must be set for acceptance tests")
	}

	if v := os.Getenv("JAMFPRO_CLIENT_ID"); v == "" {
		t.Fatal("JAMFPRO_CLIENT_ID must be set for acceptance tests")
	}

	if v := os.Getenv("JAMFPRO_CLIENT_SECRET"); v == "" {
		t.Fatal("JAMFPRO_CLIENT_SECRET must be set for acceptance tests")
	}
}


// TestAccResourceJamfProDepartment_minimum tests the creation of a department with minimal configuration.
func TestAccResourceJamfProDepartment_minimum(t *testing.T) {
	resourceName := "jamfpro_departments.department"
	departmentName := "testDepartment"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDepartmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDepartmentConfig(departmentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDepartmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", departmentName),
				),
			},
		},
	})
}

// Maximum configuration for the department
func testAccResourceDepartmentMaximumConfig(departmentName string) string {
	return fmt.Sprintf(`
resource "jamfpro_departments" "department" {
	name = "%s"
	description = "%s department description"
	metadata {
		tag = "%s"
	}
	related_id = jamfpro_related.resource.id
}

resource "jamfpro_related" "resource" {
	value = "relatedValue%s"
}
`, departmentName, departmentName, departmentName, departmentName)
}

// testAccResourceDepartmentMaximumConfig returns a string representation of a Terraform configuration
// for a department resource with maximum attributes set.
func TestAccResourceJamfProDepartment_maximum(t *testing.T) {
	resourceName := "jamfpro_departments.department"
	departmentNameInitial := "initialTestDepartment"
	departmentNameUpdated := "updatedTestDepartment"
	relatedResourceName := "jamfpro_related.resource"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDepartmentDestroy,
		Steps: []resource.TestStep{
			// Step 1: Initial configuration
			{
				Config: testAccResourceDepartmentMaximumConfig(departmentNameInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDepartmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", departmentNameInitial),
					resource.TestCheckResourceAttr(resourceName, "description", "Initial department description"),
					resource.TestCheckResourceAttr(resourceName, "metadata.tag", "initial"),
					resource.TestCheckResourceAttrSet(resourceName, "related_id"),
					resource.TestCheckResourceAttr(relatedResourceName, "value", "relatedValue1"),
				),
			},
			// Step 2: Update configuration
			{
				Config: testAccResourceDepartmentMaximumConfig(departmentNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDepartmentExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", departmentNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated department description"),
					resource.TestCheckResourceAttr(resourceName, "metadata.tag", "updated"),
					resource.TestCheckResourceAttrSet(resourceName, "related_id"),
					resource.TestCheckResourceAttr(relatedResourceName, "value", "relatedValue2"),
				),
			},
		},
	})
}

// TestAccResourceJamfProDepartment_maximum tests the creation, update, and deletion of a department with maximum configuration.
func testAccCheckDepartmentExists(n string, expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Find the resource in the Terraform state
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found in state: %s", n)
		}

		// The ID is typically used to query the real system
		departmentID := rs.Primary.ID

		// Create a client to interact with the Jamf Pro API
		client := BuildClient() // This should be your actual client creation function

		// Query the Jamf Pro API
		department, err := client.GetDepartmentByID(departmentID)
		if err != nil {
			return fmt.Errorf("Error fetching department with resource ID %s. Error: %s", departmentID, err)
		}

		if department == nil {
			return fmt.Errorf("Department with ID %s not found", departmentID)
		}

		// Check the name of the department from the service
		if department.Name != expectedName {
			return fmt.Errorf("Incorrect name from service: expected %s but got %s", expectedName, department.Name)
		}

		return nil
	}
}


// testAccResourceDepartmentConfig returns a string representation of a Terraform configuration
// for a department resource with the provided department name.
func testAccResourceDepartmentConfig(departmentName string) string {
	return fmt.Sprintf(`
resource "jamfpro_departments" "department" {
	name = "%s"
	// other attributes here...
}
`, departmentName)
}

// testAccCheckDepartmentStateVerification checks the Terraform state to ensure that the department resource
// has the expected attributes.
func testAccCheckDepartmentStateVerification(n string, expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Find the resource in the Terraform state
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found in state: %s", n)
		}

		// Check if the ID exists in the state
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		// Check the name attribute in the state
		departmentName, ok := rs.Primary.Attributes["name"]
		if !ok || departmentName != expectedName {
			return fmt.Errorf("Incorrect name: expected %s but got %s", expectedName, departmentName)
		}

		return nil
	}
}

// TestStateDriftForJamfProDepartments tests if Terraform correctly detects drift for a department resource.
// This test simulates a scenario where the department name is changed outside of Terraform.
func TestStateDriftForJamfProDepartments(t *testing.T) {
	// Setup a mock schema.ResourceData
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"id": {
			Type: schema.TypeString,
		},
		"name": {
			Type: schema.TypeString,
		},
	}, map[string]interface{}{
		"id":   "123",
		"name": "testDepartment",
	})

	// Create the mock client and set expectations
	mockClient := new(MockAPIClient)
	department := &jamfpro.Department{
		Id:   123,
		Name: "changedDepartmentName", // Simulate that department name was changed outside of Terraform
	}

	// Mock the expected response when department with ID 123 is fetched
	mockClient.On("GetDepartmentByID", 123).Return(department, nil)

	// Mock the GetDepartmentByName to return an error (since we're not expecting it to be called in this test scenario)
	mockClient.On("GetDepartmentByName", "testDepartment").Return(nil, fmt.Errorf("Not expected to be called in this scenario"))

	// Create a meta object with the mock client set as the connection
	meta := &APIClient{mockConn: mockClient}

	// Call the Read function
	diags := resourceJamfProDepartmentsRead(context.Background(), d, meta)

	// Check for no errors
	assert.Len(t, diags, 0)

	// Check if Terraform detected the drift
	assert.NotEqual(t, "testDepartment", d.Get("name").(string))
	assert.Equal(t, "changedDepartmentName", d.Get("name").(string))

	// Ensure the mock was called
	mockClient.AssertExpectations(t)
}


func TestResourceJamfProDepartmentsCreate_Success(t *testing.T) {
	mockClient := new(MockAPIClient)
	meta := &APIClient{mockConn: mockClient}

	// Setting up the mock to return a department when CreateDepartment is called.
	department := &jamfpro.Department{
		Id:   1,
		Name: "testDepartment",
	}
	mockClient.On("CreateDepartment", "testDepartment").Return(department, nil)

	d := &schema.ResourceData{}
	d.Set("name", "testDepartment")

	// Calling the Create function
	diags := resourceJamfProDepartmentsCreate(context.Background(), d, meta)

	// Assertions
	assert.Empty(t, diags, "Diags should be empty for successful creation")
	assert.Equal(t, "1", d.Id(), "ID should be set to 1 for successful creation")

	// Ensure the mock was called
	mockClient.AssertExpectations(t)
}

func TestResourceJamfProDepartmentsCreate_AlreadyExists(t *testing.T) {
	mockClient := new(MockAPIClient)
	meta := &APIClient{mockConn: mockClient}

	// Setting up the mock to return an error when CreateDepartment is called with a name that already exists.
	err := fmt.Errorf("Department with name 'testDepartment' already exists")
	mockClient.On("CreateDepartment", "testDepartment").Return(nil, err)

	d := &schema.ResourceData{}
	d.Set("name", "testDepartment")

	// Calling the Create function
	diags := resourceJamfProDepartmentsCreate(context.Background(), d, meta)

	// Assertions
	assert.NotEmpty(t, diags, "Diags should not be empty for duplicate department creation")
	assert.Contains(t, diags[0].Summary, "already exists", "Error message should mention that the department already exists")

	// Ensure the mock was called
	mockClient.AssertExpectations(t)
}

func TestResourceJamfProDepartmentsRead_Success(t *testing.T) {
	mockClient := new(MockAPIClient)
	meta := &APIClient{mockConn: mockClient}

	// Setting up the mock to return a department when GetDepartmentByID or GetDepartmentByName is called.
	department := &jamfpro.Department{
		Id:   1,
		Name: "testDepartment",
	}
	mockClient.On("GetDepartmentByID", 1).Return(department, nil)
	mockClient.On("GetDepartmentByName", "testDepartment").Return(department, nil)

	d := &schema.ResourceData{}
	d.SetId("1")
	d.Set("name", "testDepartment")

	// Calling the Read function
	diags := resourceJamfProDepartmentsRead(context.Background(), d, meta)

	// Assertions
	assert.Empty(t, diags, "Diags should be empty for successful read")
	assert.Equal(t, "testDepartment", d.Get("name").(string), "Name should be set to 'testDepartment' for successful read")

	// Ensure the mock was called
	mockClient.AssertExpectations(t)
}

func TestResourceJamfProDepartmentsRead_NotFound(t *testing.T) {
	mockClient := new(MockAPIClient)
	meta := &APIClient{mockConn: mockClient}

	// Setting up the mock to return an error when GetDepartmentByID or GetDepartmentByName is called with a non-existent department.
	err := fmt.Errorf("Department not found")
	mockClient.On("GetDepartmentByID", 1).Return(nil, err)
	mockClient.On("GetDepartmentByName", "testDepartment").Return(nil, err)

	d := &schema.ResourceData{}
	d.SetId("1")
	d.Set("name", "testDepartment")

	// Calling the Read function
	diags := resourceJamfProDepartmentsRead(context.Background(), d, meta)

	// Assertions
	assert.NotEmpty(t, diags, "Diags should not be empty for non-existent department read")
	assert.Contains(t, diags[0].Summary, "Failed to fetch department", "Error message should mention that the department was not found")

	// Ensure the mock was called
	mockClient.AssertExpectations(t)
}

func TestResourceJamfProDepartmentsUpdate_Success(t *testing.T) {
	mockClient := new(MockAPIClient)
	meta := &APIClient{mockConn: mockClient}

	// Setting up the mock to return an updated department when UpdateDepartmentByID or UpdateDepartmentByName is called.
	updatedDepartment := &jamfpro.Department{
		Id:   1,
		Name: "updatedDepartment",
	}
	mockClient.On("UpdateDepartmentByID", 1, "updatedDepartment").Return(updatedDepartment, nil)
	mockClient.On("UpdateDepartmentByName", "testDepartment", "updatedDepartment").Return(updatedDepartment, nil)

	d := &schema.ResourceData{}
	d.SetId("1")
	d.Set("name", "testDepartment")

	// Make a change to the name attribute to trigger the update
	d.Set("name", "updatedDepartment")

	// Calling the Update function
	diags := resourceJamfProDepartmentsUpdate(context.Background(), d, meta)

	// Assertions
	assert.Empty(t, diags, "Diags should be empty for successful update")
	assert.Equal(t, "updatedDepartment", d.Get("name").(string), "Name should be updated to 'updatedDepartment'")

	// Ensure the mock was called
	mockClient.AssertExpectations(t)
}

func TestResourceJamfProDepartmentsUpdate_NameInUse(t *testing.T) {
	mockClient := new(MockAPIClient)
	meta := &APIClient{mockConn: mockClient}

	// Setting up the mock to return an error indicating the department name is already in use.
	err := fmt.Errorf("Department with name 'updatedDepartment' already exists")
	mockClient.On("UpdateDepartmentByID", 1, "updatedDepartment").Return(nil, err)
	mockClient.On("UpdateDepartmentByName", "testDepartment", "updatedDepartment").Return(nil, err)

	d := &schema.ResourceData{}
	d.SetId("1")
	d.Set("name", "testDepartment")

	// Make a change to the name attribute to trigger the update
	d.Set("name", "updatedDepartment")

	// Calling the Update function
	diags := resourceJamfProDepartmentsUpdate(context.Background(), d, meta)

	// Assertions
	assert.NotEmpty(t, diags, "Diags should not be empty for name already in use")
	assert.Contains(t, diags[0].Summary, "already exists", "Error message should mention that the department name is already in use")

	// Ensure the mock was called
	mockClient.AssertExpectations(t)
}

func TestResourceJamfProDepartmentsUpdate_NotFound(t *testing.T) {
	mockClient := new(MockAPIClient)
	meta := &APIClient{mockConn: mockClient}

	// Setting up the mock to return an error indicating the department was not found.
	err := fmt.Errorf("Department not found")
	mockClient.On("UpdateDepartmentByID", 1, "updatedDepartment").Return(nil, err)
	mockClient.On("UpdateDepartmentByName", "testDepartment", "updatedDepartment").Return(nil, err)

	d := &schema.ResourceData{}
	d.SetId("1")
	d.Set("name", "testDepartment")

	// Make a change to the name attribute to trigger the update
	d.Set("name", "updatedDepartment")

	// Calling the Update function
	diags := resourceJamfProDepartmentsUpdate(context.Background(), d, meta)

	// Assertions
	assert.NotEmpty(t, diags, "Diags should not be empty for non-existent department update")
	assert.Contains(t, diags[0].Summary, "Failed to update department", "Error message should mention that the department was not found")

	// Ensure the mock was called
	mockClient.AssertExpectations(t)
}

func TestResourceJamfProDepartmentsDelete_Success(t *testing.T) {
	mockClient := new(MockAPIClient)
	meta := &APIClient{mockConn: mockClient}

	// Setting up the mock to handle successful deletion when DeleteDepartmentByID or DeleteDepartmentByName is called.
	mockClient.On("DeleteDepartmentByID", 1).Return(nil)
	mockClient.On("DeleteDepartmentByName", "testDepartment").Return(nil)

	d := &schema.ResourceData{}
	d.SetId("1")
	d.Set("name", "testDepartment")

	// Calling the Delete function
	diags := resourceJamfProDepartmentsDelete(context.Background(), d, meta)

	// Assertions
	assert.Empty(t, diags, "Diags should be empty for successful deletion")

	// Ensure the mock was called
	mockClient.AssertExpectations(t)
}

func TestResourceJamfProDepartmentsDelete_NotFound(t *testing.T) {
	mockClient := new(MockAPIClient)
	meta := &APIClient{mockConn: mockClient}

	// Setting up the mock to return an error indicating the department was not found.
	err := fmt.Errorf("Department not found")
	mockClient.On("DeleteDepartmentByID", 1).Return(err)
	mockClient.On("DeleteDepartmentByName", "testDepartment").Return(err)

	d := &schema.ResourceData{}
	d.SetId("1")
	d.Set("name", "testDepartment")

	// Calling the Delete function
	diags := resourceJamfProDepartmentsDelete(context.Background(), d, meta)

	// Assertions
	assert.NotEmpty(t, diags, "Diags should not be empty for non-existent department deletion")
	assert.Contains(t, diags[0].Summary, "Failed to delete department", "Error message should mention that the department was not found")

	// Ensure the mock was called
	mockClient.AssertExpectations(t)
}
*/
