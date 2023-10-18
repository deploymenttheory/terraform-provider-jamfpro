package provider

import (
	"os"
	"testing"

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
	os.Setenv("JAMFPRO_INSTANCE", "testEnvInstance")
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
	os.Setenv("JAMFPRO_CLIENT_ID", "testEnvClientID")
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
	// Similar tests as for getInstanceName but for client secret
}

func TestProvider(t *testing.T) {
	// Test the Provider function, especially the ConfigureContextFunc
	// This might be a bit more complex and might require more mocking

	// Test 1: Everything set correctly

	// Test 2: Missing instance name

	// Test 3: Missing client ID

	// Test 4: Missing client secret

	// ... and so on
}
