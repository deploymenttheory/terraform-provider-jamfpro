package sdk

import (
	"github.com/deploymenttheory/go-api-http-client-integrations/jamf/jamfprointegration"
	"github.com/deploymenttheory/go-api-http-client/httpclient"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

// The types and builders the provider uses to authenticate and to configure
// the HTTP client beneath the SDK.
type (
	Integration         = jamfprointegration.Integration
	HTTPClientConfig    = httpclient.ClientConfig
	RequestHeaderConfig = httpclient.RequestHeaderConfig
)

// Assigned, not wrapped, so errors reach the provider's diagnostics unchanged.
var (
	BuildWithOAuth                = jamfprointegration.BuildWithOAuth
	BuildWithBasicAuth            = jamfprointegration.BuildWithBasicAuth
	BuildWithPlatformGatewayOAuth = jamfprointegration.BuildWithPlatformGatewayOAuth
)

// NewClient returns a Client that sends its requests through httpClient.
//
// The concrete type is *jamfpro.Client, which services still assert on.
func NewClient(httpClient *httpclient.Client) Client {
	return &jamfpro.Client{HTTP: httpClient}
}
