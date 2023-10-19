// interfaces.go
package provider

import "github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

var _ JamfProDepartmentCRUDOperations = (*jamfpro.Client)(nil)

type APIClient struct {
	mockConn JamfProDepartmentCRUDOperations `used_for:"mocking"`
	conn     *jamfpro.Client                 `used_for:"api_calls"`
}

type JamfProDepartmentCRUDOperations interface {
	GetDepartments() (*jamfpro.ResponseDepartments, error)
	GetDepartmentByID(id int) (*jamfpro.Department, error)
	GetDepartmentByName(name string) (*jamfpro.Department, error)
	GetDepartmentIdByName(name string) (int, error)
	CreateDepartment(departmentName string) (*jamfpro.Department, error)
	UpdateDepartmentByID(id int, departmentName string) (*jamfpro.Department, error)
	UpdateDepartmentByName(oldName string, newName string) (*jamfpro.Department, error)
	DeleteDepartmentByID(id int) error
	DeleteDepartmentByName(name string) error
}
