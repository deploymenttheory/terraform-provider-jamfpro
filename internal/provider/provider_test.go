package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestGetInstanceName(t *testing.T) {
	// Mock schema.ResourceData
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"instance_name": {
			Type: schema.TypeString,
		},
	}, map[string]interface{}{})

	// Test 1: Get instance name from d
	d.Set("instance_name", "testInstance")
	name, err := GetInstanceName(d)
	assert.NoError(t, err)
	assert.Equal(t, "testInstance", name)

	// Test 2: Get instance name from environment variable
	t.Setenv("JAMFPRO_INSTANCE", "testEnvInstance")
	d.Set("instance_name", "") // Clear the previous set value
	name, err = GetInstanceName(d)
	assert.NoError(t, err)
	assert.Equal(t, "testEnvInstance", name)

	// Test 3: No instance name set anywhere
	os.Unsetenv("JAMFPRO_INSTANCE")
	_, err = GetInstanceName(d)
	assert.Error(t, err)
}

func TestGetClientID(t *testing.T) {
	// Mock schema.ResourceData
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"client_id": {
			Type: schema.TypeString,
		},
	}, map[string]interface{}{})

	// Test 1: Get client ID from d
	d.Set("client_id", "testClientID")
	clientID, err := GetClientID(d)
	assert.NoError(t, err)
	assert.Equal(t, "testClientID", clientID)

	// Test 2: Get client ID from environment variable
	t.Setenv("JAMFPRO_CLIENT_ID", "testEnvClientID")
	d.Set("client_id", "") // Clear the previous set value
	clientID, err = GetClientID(d)
	assert.NoError(t, err)
	assert.Equal(t, "testEnvClientID", clientID)

	// Test 3: No client ID set anywhere
	os.Unsetenv("JAMFPRO_CLIENT_ID")
	_, err = GetClientID(d)
	assert.Error(t, err)
}

func TestGetClientSecret(t *testing.T) {
	// Mock schema.ResourceData
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"client_secret": {
			Type: schema.TypeString,
		},
	}, map[string]interface{}{})

	// Test 1: Get client secret from d
	d.Set("client_secret", "testClientSecret")
	clientSecret, err := GetClientSecret(d)
	assert.NoError(t, err)
	assert.Equal(t, "testClientSecret", clientSecret)

	// Test 2: Get client secret from environment variable
	t.Setenv("JAMFPRO_CLIENT_SECRET", "testEnvClientSecret")
	d.Set("client_secret", "") // Clear the previous set value
	clientSecret, err = GetClientSecret(d)
	assert.NoError(t, err)
	assert.Equal(t, "testEnvClientSecret", clientSecret)

	// Test 3: No client secret set anywhere
	os.Unsetenv("JAMFPRO_CLIENT_SECRET")
	_, err = GetClientSecret(d)
	assert.Error(t, err)
}

func TestProvider(t *testing.T) {
	// Test 1: Everything set correctly
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"instance_name": {
			Type: schema.TypeString,
		},
		"client_id": {
			Type: schema.TypeString,
		},
		"client_secret": {
			Type: schema.TypeString,
		},
		"debug_mode": {
			Type: schema.TypeBool,
		},
	}, map[string]interface{}{
		"instance_name": "testInstance",
		"client_id":     "testClientID",
		"client_secret": "testClientSecret",
		"debug_mode":    true,
	})

	_, diags := Provider().ConfigureContextFunc(context.Background(), d)
	assert.Len(t, diags, 0)

	// Test 2: Missing instance name
	d.Set("instance_name", "")
	_, diags = Provider().ConfigureContextFunc(context.Background(), d)
	assert.Len(t, diags, 1)

	// Test 3: Missing client ID
	d.Set("instance_name", "testInstance") // reset instance_name
	d.Set("client_id", "")
	_, diags = Provider().ConfigureContextFunc(context.Background(), d)
	assert.Len(t, diags, 1)

	// Test 4: Missing client secret
	d.Set("client_id", "testClientID") // reset client_id
	d.Set("client_secret", "")
	_, diags = Provider().ConfigureContextFunc(context.Background(), d)
	assert.Len(t, diags, 1)
}

// Return a mock client.
func mockNewClientSuccess(cfg jamfpro.Config) (*jamfpro.Client, error) {
	return &jamfpro.Client{}, nil
}

func mockNewClientFail(cfg jamfpro.Config) (*jamfpro.Client, error) {
	return nil, fmt.Errorf("mocked client initialization failure")
}

func TestProviderWithSuccessfulClientInitialization(t *testing.T) {
	// Setup your schema.ResourceData
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"instance_name": {Type: schema.TypeString},
		"client_id":     {Type: schema.TypeString},
		"client_secret": {Type: schema.TypeString},
		"debug_mode":    {Type: schema.TypeBool},
	}, map[string]interface{}{
		"instance_name": "testInstance",
		"client_id":     "testClientID",
		"client_secret": "testClientSecret",
		"debug_mode":    true,
	})

	// Override the function for this test
	client.BuildClient = mockNewClientSuccess

	// Ensure BuildClient is reset after the test
	defer func() {
		client.BuildClient = jamfpro.NewClient
	}()

	// Now invoke the provider with the mock setup
	_, diags := Provider().ConfigureContextFunc(context.Background(), d)
	assert.Len(t, diags, 0) // Expect no diagnostics (errors)
}

func TestProviderWithFailedClientInitialization(t *testing.T) {
	// Setup your schema.ResourceData
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"instance_name": {Type: schema.TypeString},
		"client_id":     {Type: schema.TypeString},
		"client_secret": {Type: schema.TypeString},
		"debug_mode":    {Type: schema.TypeBool},
	}, map[string]interface{}{
		"instance_name": "testInstance",
		"client_id":     "testClientID",
		"client_secret": "testClientSecret",
		"debug_mode":    true,
	})
	// Override the function for this test
	client.BuildClient = mockNewClientFail

	// Ensure BuildClient is reset after the test
	defer func() {
		client.BuildClient = jamfpro.NewClient
	}()

	// Now invoke the provider with the mock setup
	_, diags := Provider().ConfigureContextFunc(context.Background(), d)
	assert.Len(t, diags, 1) // Expect one diagnostic (error)
}

/*
func TestUserAgentInitialization(t *testing.T) {
	// This is the expected UserAgent format
	expectedUserAgent := fmt.Sprintf("%s/%s", TerraformProviderProductUserAgent, version.ProviderVersion)

	// Prepare a map of schema and values for the test
	resourceDataMap := map[string]*schema.Schema{
		"instance_name": {
			Type: schema.TypeString,
		},
		"client_id": {
			Type: schema.TypeString,
		},
		"client_secret": {
			Type: schema.TypeString,
		},
		"log_level": {
			Type: schema.TypeString,
		},
	}
	resourceDataValues := map[string]interface{}{
		"instance_name": "testInstance",
		"client_id":     "testClientID",
		"client_secret": "testClientSecret",
		"log_level":     "info",
	}

	// Create a new ResourceData object for the test
	d := schema.TestResourceDataRaw(t, resourceDataMap, resourceDataValues)

	// Call the provider's ConfigureContextFunc to Get a configured client
	clientInterface, diags := Provider().ConfigureContextFunc(context.Background(), d)
	assert.Len(t, diags, 0) // Ensure no errors are returned

	// Cast the clientInterface to the expected *client.APIClient type
	apiClient, ok := clientInterface.(*client.APIClient)
	assert.True(t, ok) // Assert that the interface{} is indeed of *client.APIClient type

	// Assert that the UserAgent in the returned client matches the expected format
	assert.Equal(t, expectedUserAgent, apiClient.Conn.UserAgent)
}

func mockSDKError(cfg jamfpro.Config) (*jamfpro.Client, error) {
	return nil, fmt.Errorf("deeper error for testing propagation")
}

func TestErrorPropagation(t *testing.T) {
	client.BuildClient = mockSDKError

	defer func() {
		client.BuildClient = jamfpro.NewClient
	}()

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"instance_name": {Type: schema.TypeString},
		"client_id":     {Type: schema.TypeString},
		"client_secret": {Type: schema.TypeString},
		"debug_mode":    {Type: schema.TypeBool},
	}, map[string]interface{}{
		"instance_name": "testInstance",
		"client_id":     "testClientID",
		"client_secret": "testClientSecret",
		"debug_mode":    true,
	})

	_, diags := Provider().ConfigureContextFunc(context.Background(), d)
	assert.Len(t, diags, 1)
	assert.Contains(t, diags[0].Summary, "deeper error for testing propagation")
}

func TestSensitiveInformationLogging(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"instance_name": {Type: schema.TypeString},
		"client_id":     {Type: schema.TypeString},
		"client_secret": {Type: schema.TypeString},
		"debug_mode":    {Type: schema.TypeBool},
	}, map[string]interface{}{
		"instance_name": "testInstance",
		"client_id":     "testClientID",
		"client_secret": "testClientSecret",
		"debug_mode":    true,
	})

	Provider().ConfigureContextFunc(context.Background(), d)

	logs := buf.String()
	assert.NotContains(t, logs, "testClientSecret")
}
*/
