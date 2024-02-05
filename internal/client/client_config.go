// config_client.go
package client

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/stretchr/testify/mock"
)

type ProviderConfig struct {
	InstanceName string
	ClientID     string
	ClientSecret string
	LogLevel     http_client.LogLevel
	UserAgent    string
}

// APIClient is a HTTP API Client.
type APIClient struct {
	Conn *jamfpro.Client
}

// MockAPIClient is a mock version of APIClient for testing.
type MockAPIClient struct {
	mock.Mock
}

// BuildClient is a global function variable for client creation that defaults to jamfpro.NewClient.
// It can be overridden in tests to use mock client creation functions.
var BuildClient = jamfpro.NewClient

// This function maps the string values received from Terraform configuration
// to the LogLevel constants defined in the http_client package.
func convertToJamfProLogLevel(logLevelString string) (http_client.LogLevel, error) {
	switch logLevelString {
	case "debug":
		return http_client.LogLevelDebug, nil
	case "info":
		return http_client.LogLevelInfo, nil
	case "warning":
		return http_client.LogLevelWarning, nil
	case "none":
		return http_client.LogLevelNone, nil
	default:
		return http_client.LogLevelNone, fmt.Errorf("invalid log level: %s", logLevelString)
	}
}

// ConvertToLogLevel wraps the convertToJamfProLogLevel function to match the expected function signature.
func ConvertToLogLevel(logLevelString string) (http_client.LogLevel, error) {
	return convertToJamfProLogLevel(logLevelString)
}

// Client returns a new client for accessing Jamf Pro.
func (c *ProviderConfig) Client() (*APIClient, diag.Diagnostics) {
	var client APIClient

	jamfProConfig := http_client.Config{
		InstanceName: c.InstanceName,
		Auth: http_client.AuthConfig{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
		},
		LogLevel: c.LogLevel,
	}

	jamfProClient, err := BuildClient(jamfProConfig)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	client.Conn = jamfProClient
	return &client, nil
}
