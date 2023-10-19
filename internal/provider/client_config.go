// config_client.go
package provider

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
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

/*
type APIClient struct {
	conn JamfProClientInterface
}
*/

// NewClientFunc is a global function variable for client creation that defaults to jamfpro.NewClient.
// It can be overridden in tests to use mock client creation functions.
var NewClientFunc = jamfpro.NewClient

// Client returns a new client for accessing Jamf Pro.
func (c *ProviderConfig) Client() (*APIClient, diag.Diagnostics) {
	var client APIClient

	jamfProConfig := jamfpro.Config{
		InstanceName: c.InstanceName,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		DebugMode:    c.DebugMode,
	}

	jamfProClient, err := NewClientFunc(jamfProConfig)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	client.conn = jamfProClient
	return &client, nil
}

func mockSDKError(cfg jamfpro.Config) (*jamfpro.Client, error) {
	return nil, fmt.Errorf("deeper error for testing propagation")
}

func TestErrorPropagation(t *testing.T) {
	NewClientFunc = mockSDKError

	defer func() {
		NewClientFunc = jamfpro.NewClient
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
