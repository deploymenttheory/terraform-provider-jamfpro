// config_client.go
package client

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/stretchr/testify/mock"
)

// MockAPIClient is a mock version of APIClient for testing.
type MockAPIClient struct {
	mock.Mock
}

// APIClient wraps the Jamf Pro SDK Client.
type APIClient struct {
	Conn            *jamfpro.Client // Use the SDK client
	EnableCookieJar bool            // Use the cookie jar value
}

// This function maps the string log level from the Terraform configuration
// to the appropriate log level in the httpclient package.
func convertToJamfProLogLevel(logLevelString string) (jamfpro.LogLevel, error) {
	switch logLevelString {
	case "debug":
		return jamfpro.LogLevelDebug, nil
	case "info":
		return jamfpro.LogLevelInfo, nil
	case "warning":
		return jamfpro.LogLevelWarning, nil
	case "none":
		return jamfpro.LogLevelNone, nil
	default:
		return jamfpro.LogLevelNone, fmt.Errorf("invalid log level: %s", logLevelString)
	}
}

// ConvertToLogLevel is a helper function to convert log level strings
// from the Terraform configuration to the SDK's log level type.
func ConvertToLogLevel(logLevelString string) (jamfpro.LogLevel, error) {
	return convertToJamfProLogLevel(logLevelString)
}
