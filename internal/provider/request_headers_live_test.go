//go:build integration

package provider

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestLiveProxyHeaders uses a separately started reverse proxy forwarding to a real
// Jamf Pro instance. The proxy shortens token expiry metadata to exercise refresh.
// Credentials live in an external private file, never in the repository.
func TestLiveProxyHeaders(t *testing.T) {
	path := os.Getenv("JAMFPRO_PROXY_LIVE_CONFIG")
	if path == "" {
		t.Skip("set JAMFPRO_PROXY_LIVE_CONFIG to the private proxy test configuration")
	}
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	var config struct {
		ProxyURL           string `json:"proxy_url"`
		ClientID           string `json:"client_id"`
		ClientSecret       string `json:"client_secret"`
		ProxyAuthorization string `json:"proxy_authorization"`
		ExpectedVersion    string `json:"expected_version"`
	}
	require.NoError(t, json.Unmarshal(data, &config))
	for _, framework := range []bool{false, true} {
		t.Run(map[bool]string{false: "SDKv2", true: "Framework"}[framework], func(t *testing.T) {
			input := headerTestConfig(config.ProxyURL)
			input["client_id"], input["client_secret"] = config.ClientID, config.ClientSecret
			input["custom_headers"] = []any{map[string]any{"name": "Authorization", "value": config.ProxyAuthorization}}
			client, failed := configureHeaderTestProvider(t, framework, input)
			require.False(t, failed, "provider configuration through the live proxy failed")
			time.Sleep(1100 * time.Millisecond)
			version, err := client.GetJamfProVersion()
			require.NoError(t, err)
			require.Equal(t, config.ExpectedVersion, *version.Version)
		})
	}
}

func TestLivePlatformEnvironmentScope(t *testing.T) {
	environmentID := os.Getenv("JAMFPLATFORM_ENVIRONMENT_ID")
	clientID := os.Getenv("JAMFPLATFORM_CLIENT_ID")
	clientSecret := os.Getenv("JAMFPLATFORM_CLIENT_SECRET")
	baseURL := os.Getenv("JAMFPLATFORM_BASE_URL")
	if environmentID == "" || clientID == "" || clientSecret == "" || baseURL == "" {
		t.Skip("set JAMFPLATFORM_ENVIRONMENT_ID, JAMFPLATFORM_CLIENT_ID, JAMFPLATFORM_CLIENT_SECRET, and JAMFPLATFORM_BASE_URL")
	}

	for _, framework := range []bool{false, true} {
		t.Run(map[bool]string{false: "SDKv2", true: "Framework"}[framework], func(t *testing.T) {
			input := platformEnvironmentTestConfig(baseURL)
			input["platform_environment_id"] = environmentID
			input["client_id"] = clientID
			input["client_secret"] = clientSecret

			client, failed := configureHeaderTestProvider(t, framework, input)
			require.False(t, failed, "provider configuration against the live Jamf environment failed")

			version, err := client.GetJamfProVersion()
			require.NoError(t, err)
			require.NotNil(t, version)
			require.NotNil(t, version.Version)
			require.NotEmpty(t, *version.Version)
		})
	}
}
