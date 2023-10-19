// department_resource_test.go
package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// JamfProClientInterface captures the methods we expect on the jamfpro.Client.
type JamfProClientInterface interface {
	GetDepartmentByID(id int) (*jamfpro.Department, error)
	GetDepartmentByName(name string) (*jamfpro.Department, error)
	CreateDepartment(departmentName string) (*jamfpro.Department, error)
	UpdateDepartmentByID(id int, departmentName string) (*jamfpro.Department, error)
	UpdateDepartmentByName(oldName string, newName string) (*jamfpro.Department, error)
	DeleteDepartmentByID(id int) error
	DeleteDepartmentByName(name string) error
}

// Ensure mockJamfProClient implements JamfProClientInterface
var _ JamfProClientInterface = (*mockJamfProClient)(nil)

// Mock client for Jamf Pro
type mockJamfProClient struct {
	departments map[int]string
}

func (m *mockJamfProClient) GetDepartmentByID(id int) (*jamfpro.Department, error) {
	name, exists := m.departments[id]
	if !exists {
		return nil, fmt.Errorf("department not found")
	}
	return &jamfpro.Department{
		Id:   id,
		Name: name,
	}, nil
}

func (m *mockJamfProClient) GetDepartmentByName(name string) (*jamfpro.Department, error) {
	// Just a dummy implementation for the sake of this example
	for id, depName := range m.departments {
		if depName == name {
			return &jamfpro.Department{
				Id:   id,
				Name: name,
			}, nil
		}
	}
	return nil, fmt.Errorf("department not found")
}

func (m *mockJamfProClient) CreateDepartment(departmentName string) (*jamfpro.Department, error) {
	// For simplicity, just creating a dummy department with a random ID.
	id := len(m.departments) + 1
	m.departments[id] = departmentName

	return &jamfpro.Department{
		Id:   id,
		Name: departmentName,
	}, nil
}

func (m *mockJamfProClient) UpdateDepartmentByID(id int, departmentName string) (*jamfpro.Department, error) {
	// Check if the department with the given ID exists
	_, exists := m.departments[id]
	if !exists {
		return nil, fmt.Errorf("department with ID %d not found", id)
	}

	// Update the department name
	m.departments[id] = departmentName

	// Return the updated department
	return &jamfpro.Department{
		Id:   id,
		Name: departmentName,
	}, nil
}

func (m *mockJamfProClient) UpdateDepartmentByName(oldName string, newName string) (*jamfpro.Department, error) {
	// Find department ID by old name
	var departmentID int
	found := false
	for id, name := range m.departments {
		if name == oldName {
			departmentID = id
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("department with name %s not found", oldName)
	}

	// Update the department name
	m.departments[departmentID] = newName

	// Return the updated department
	return &jamfpro.Department{
		Id:   departmentID,
		Name: newName,
	}, nil
}

func (m *mockJamfProClient) DeleteDepartmentByID(id int) error {
	_, exists := m.departments[id]
	if !exists {
		return fmt.Errorf("department with ID %d not found", id)
	}

	// Delete the department
	delete(m.departments, id)

	return nil
}

func (m *mockJamfProClient) DeleteDepartmentByName(name string) error {
	var departmentID int
	found := false
	for id, deptName := range m.departments {
		if deptName == name {
			departmentID = id
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("department with name %s not found", name)
	}

	// Delete the department
	delete(m.departments, departmentID)

	return nil
}

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

	// Mock the APIClient and its methods to simulate real-world interactions
	mockClient := &APIClient{
		conn: &mockJamfProClient{
			departments: map[int]string{
				123: "changedDepartmentName", // Simulate that department name was changed outside of Terraform
			},
		},
	}

	// Call the Read function
	diags := resourceJamfProDepartmentsRead(context.Background(), d, mockClient)

	// Check for no errors
	assert.Len(t, diags, 0)

	// Check if Terraform detected the drift
	assert.NotEqual(t, "testDepartment", d.Get("name").(string))
	assert.Equal(t, "changedDepartmentName", d.Get("name").(string))
}
