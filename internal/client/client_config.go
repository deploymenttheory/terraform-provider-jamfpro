// config_client.go
package client

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/stretchr/testify/mock"
)

// APIClient is a HTTP API Client.
type APIClient struct {
	Conn *jamfpro.Client
}

// MockAPIClient is a mock version of APIClient for testing.
type MockAPIClient struct {
	mock.Mock
}

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
