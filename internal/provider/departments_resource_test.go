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

// Ensure mockJamfProClient implements ClientInterface
var _ ClientInterface = (*mockJamfProClient)(nil)

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
