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
type APIClient struct {
	conn *jamfpro.Client
}

// Client returns a new client for accessing Jamf Pro.
func (c *ProviderConfig) Client() (*APIClient, diag.Diagnostics) {
	var client APIClient

	jamfProConfig := jamfpro.Config{ // Ensure it's using jamfpro.Config
		InstanceName: c.InstanceName,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		DebugMode:    c.DebugMode,
	}

	jamfProClient, err := jamfpro.NewClient(jamfProConfig) // Using jamfpro.NewClient
	if err != nil {
		return nil, diag.FromErr(err)
	}

	client.conn = jamfProClient
	return &client, nil
}
