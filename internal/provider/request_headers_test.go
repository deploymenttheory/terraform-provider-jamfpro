package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func headerConfigValue(ty tftypes.Type, input any) tftypes.Value {
	if input == nil {
		return tftypes.NewValue(ty, nil)
	}
	switch ty := ty.(type) {
	case tftypes.Object:
		data := input.(map[string]any)
		values := make(map[string]tftypes.Value, len(ty.AttributeTypes))
		for name, child := range ty.AttributeTypes {
			values[name] = headerConfigValue(child, data[name])
		}
		return tftypes.NewValue(ty, values)
	case tftypes.List:
		data := input.([]any)
		values := make([]tftypes.Value, len(data))
		for i, child := range data {
			values[i] = headerConfigValue(ty.ElementType, child)
		}
		return tftypes.NewValue(ty, values)
	default:
		return tftypes.NewValue(ty, input)
	}
}

func configureHeaderTestProvider(t *testing.T, framework bool, input map[string]any) (*jamfpro.Client, bool) {
	t.Helper()
	ctx := context.Background()
	if !framework {
		p := Provider()
		data := schema.TestResourceDataRaw(t, p.Schema, input)
		client, diags := p.ConfigureContextFunc(ctx, data)
		if diags.HasError() {
			return nil, true
		}
		return client.(*jamfpro.Client), false
	}
	p := FrameworkProvider("test")()
	var definition frameworkprovider.SchemaResponse
	p.Schema(ctx, frameworkprovider.SchemaRequest{}, &definition)
	protocol := providerserver.NewProtocol6(p)()
	protocolSchema, err := protocol.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	config := tfsdk.Config{Schema: definition.Schema, Raw: headerConfigValue(protocolSchema.Provider.ValueType(), input)}
	var response frameworkprovider.ConfigureResponse
	p.Configure(ctx, frameworkprovider.ConfigureRequest{Config: config}, &response)
	if response.Diagnostics.HasError() {
		return nil, true
	}
	return response.ResourceData.(*jamfpro.Client), false
}

func headerTestConfig(baseURL string) map[string]any {
	return map[string]any{ // #nosec G101 -- Dummy credentials accepted only by the local httptest OAuth server.
		"jamfpro_instance_fqdn": baseURL + "/proxypath/", "auth_method": "oauth2", "auth_provider": "direct",
		"client_id": "test-client", "client_secret": "test-client-secret", "hide_sensitive_data": true,
		"enable_client_sdk_logs": false, "jamfpro_load_balancer_lock": false,
		"token_refresh_buffer_period_seconds": 0, "mandatory_request_delay_milliseconds": 0,
		"bearer_token_header": "X-Jamf-Token",
		"custom_headers": []any{
			map[string]any{"name": "Authorization", "value": "Basic proxy-credentials"},
			map[string]any{"name": "X-Company", "value": "company"},
		},
	}
}

func TestProviderRequestHeadersAndTokenRefresh(t *testing.T) {
	for _, framework := range []bool{false, true} {
		t.Run(map[bool]string{false: "SDKv2", true: "Framework"}[framework], func(t *testing.T) {
			var tokens atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Header.Get("Authorization") != "Basic proxy-credentials" || req.Header.Get("X-Company") != "company" {
					t.Error("proxy credentials or routing header missing")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if req.URL.Path == "/proxypath/api/v1/oauth/token" {
					if err := req.ParseForm(); err != nil {
						t.Error(err)
					}
					if req.Method != http.MethodPost || req.PostForm.Get("client_id") != "test-client" || req.PostForm.Get("client_secret") != "test-client-secret" || req.PostForm.Get("grant_type") != "client_credentials" {
						t.Error("OAuth form changed")
					}
					token := fmt.Sprintf("token-%d", tokens.Add(1))
					if err := json.NewEncoder(w).Encode(map[string]any{"access_token": token, "expires_in": 1, "token_type": "Bearer"}); err != nil {
						t.Error(err)
					}
					return
				}
				if req.URL.Path != "/proxypath/api/v1/jamf-pro-version" || req.Header.Get("X-Jamf-Token") != fmt.Sprintf("Bearer token-%d", tokens.Load()) {
					t.Error("API path or refreshed bearer token incorrect")
				}
				if err := json.NewEncoder(w).Encode(map[string]string{"version": "11.33.3"}); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			client, failed := configureHeaderTestProvider(t, framework, headerTestConfig(server.URL))
			require.False(t, failed, "provider configuration failed")
			require.Equal(t, int32(1), tokens.Load())
			time.Sleep(1100 * time.Millisecond)
			version, err := client.GetJamfProVersion()
			require.NoError(t, err)
			require.Equal(t, "11.33.3", *version.Version)
			require.Equal(t, int32(2), tokens.Load(), "the header must use the automatically refreshed token")
		})
	}
}

func TestProviderHeaderSchemasMatch(t *testing.T) {
	ctx := context.Background()
	legacy, err := tf5to6server.UpgradeServer(ctx, Provider().GRPCProvider)
	require.NoError(t, err)
	mux, err := tf6muxserver.NewMuxServer(ctx,
		func() tfprotov6.ProviderServer { return legacy }, providerserver.NewProtocol6(FrameworkProvider("test")()))
	require.NoError(t, err)
	response, err := mux.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	for _, d := range response.Diagnostics {
		require.NotEqual(t, tfprotov6.DiagnosticSeverityError, d.Severity, "%s: %s", d.Summary, d.Detail)
	}
	foundTokenHeader, foundCustomHeaders := false, false
	for _, attribute := range response.Provider.Block.Attributes {
		if attribute.Name == "bearer_token_header" {
			foundTokenHeader = true
			require.NotEmpty(t, attribute.Description)
		}
	}
	for _, block := range response.Provider.Block.BlockTypes {
		if block.TypeName == "custom_headers" {
			foundCustomHeaders = true
			require.NotEmpty(t, block.Block.Description)
			for _, attribute := range block.Block.Attributes {
				require.NotEmpty(t, attribute.Description)
				if attribute.Name == "value" {
					require.True(t, attribute.Sensitive)
				}
			}
		}
	}
	require.True(t, foundTokenHeader)
	require.True(t, foundCustomHeaders)
}

func TestProviderDefaultHeadersRemainUnchanged(t *testing.T) {
	t.Setenv("JAMFPRO_BEARER_TOKEN_HEADER", "")
	for _, framework := range []bool{false, true} {
		t.Run(map[bool]string{false: "SDKv2", true: "Framework"}[framework], func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				result := map[string]any{"version": "11.33.3"}
				if req.URL.Path == "/proxypath/api/v1/oauth/token" {
					result = map[string]any{"access_token": "direct-token", "expires_in": 3600}
				} else if req.Header.Get("Authorization") != "Bearer direct-token" {
					t.Error("default bearer header changed")
				}
				if err := json.NewEncoder(w).Encode(result); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			input := headerTestConfig(server.URL)
			delete(input, "bearer_token_header")
			delete(input, "custom_headers")
			_, failed := configureHeaderTestProvider(t, framework, input)
			require.False(t, failed)
		})
	}
}

func TestProviderRejectsInvalidHeadersBeforeAuthentication(t *testing.T) {
	for _, framework := range []bool{false, true} {
		t.Run(map[bool]string{false: "SDKv2", true: "Framework"}[framework], func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				calls.Add(1)
				w.WriteHeader(http.StatusUnauthorized)
			}))
			defer server.Close()
			for name, headers := range map[string][]any{
				"duplicate":       {map[string]any{"name": "X-Key", "value": "first"}, map[string]any{"name": "x-key", "value": "second"}},
				"invalid":         {map[string]any{"name": "X-Key", "value": "value\r\nInjected: value"}},
				"token collision": {map[string]any{"name": "X-Jamf-Token", "value": "stale-token"}},
			} {
				t.Run(name, func(t *testing.T) {
					input := headerTestConfig(server.URL)
					input["custom_headers"] = headers
					_, failed := configureHeaderTestProvider(t, framework, input)
					require.True(t, failed)
					require.Zero(t, calls.Load())
				})
			}
		})
	}
}
