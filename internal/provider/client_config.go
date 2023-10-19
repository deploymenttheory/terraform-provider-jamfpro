// config_client.go
package provider

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type ProviderConfig struct {
	InstanceName string
	ClientID     string
	ClientSecret string
	DebugMode    bool
	UserAgent    string
}

// APIClient is a HTTP API Client.
/*type APIClient struct {
	conn *jamfpro.Client
}*/

type APIClient struct {
	conn JamfProDepartmentCRUDOperations
}

// BuildClient is a global function variable for client creation that defaults to jamfpro.NewClient.
// It can be overridden in tests to use mock client creation functions.
var BuildClient = jamfpro.NewClient

// Client returns a new client for accessing Jamf Pro.
func (c *ProviderConfig) Client() (*APIClient, diag.Diagnostics) {
	var client APIClient

	jamfProConfig := jamfpro.Config{
		InstanceName: c.InstanceName,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		DebugMode:    c.DebugMode,
	}

	jamfProClient, err := BuildClient(jamfProConfig)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	client.conn = jamfProClient
	return &client, nil
}
