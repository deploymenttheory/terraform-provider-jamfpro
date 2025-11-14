package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-http-client-integrations/jamf/jamfprointegration"
	"github.com/deploymenttheory/go-api-http-client/httpclient"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"go.uber.org/zap"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &frameworkProvider{}
)

const minimumJamfProVersion = "11.20"

// frameworkProvider defines the provider implementation for Framework-based resources.
type frameworkProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// FrameworkProvider is a helper function to simplify provider setup.
func FrameworkProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &frameworkProvider{
			version: version,
		}
	}
}

func (p *frameworkProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jamfpro"
	resp.Version = p.version
}

func (p *frameworkProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"jamfpro_instance_fqdn": schema.StringAttribute{
				Optional:    true,
				Description: "The Jamf Pro FQDN (fully qualified domain name). example: https://mycompany.jamfcloud.com",
			},
			"auth_method": schema.StringAttribute{
				Optional:    true,
				Description: "The auth method chosen for interacting with Jamf Pro. Options are 'basic' for username/password or 'oauth2' for client id/secret.",
				Validators: []validator.String{
					stringvalidator.OneOf("basic", "oauth2"),
				},
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Jamf Pro Client ID for authentication when auth_method is 'oauth2'.",
			},
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The Jamf Pro Client secret for authentication when auth_method is 'oauth2'.",
			},
			"basic_auth_username": schema.StringAttribute{
				Optional:    true,
				Description: "The Jamf Pro username used for authentication when auth_method is 'basic'.",
			},
			"basic_auth_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The Jamf Pro password used for authentication when auth_method is 'basic'.",
			},
			"enable_client_sdk_logs": schema.BoolAttribute{
				Optional:    true,
				Description: "Debug option to propagate logs from the SDK and HttpClient",
			},
			"client_sdk_log_export_path": schema.StringAttribute{
				Optional:    true,
				Description: "Specify the path to export http client logs to.",
			},
			"hide_sensitive_data": schema.BoolAttribute{
				Optional:    true,
				Description: "Define whether sensitive fields should be hidden in logs. Default to hiding sensitive data in logs",
			},
			"jamfpro_load_balancer_lock": schema.BoolAttribute{
				Optional:    true,
				Description: "Programatically determines all available web app members in the load balancer and locks all instances of httpclient to the app for faster executions. \nTEMP SOLUTION UNTIL JAMF PROVIDES SOLUTION",
			},
			"token_refresh_buffer_period_seconds": schema.Int64Attribute{
				Optional:    true,
				Description: "The buffer period in seconds for token refresh.",
			},
			"mandatory_request_delay_milliseconds": schema.Int64Attribute{
				Optional:    true,
				Description: "A mandatory delay after each request before returning to reduce high volume of requests in a short time",
			},
		},
		Blocks: map[string]schema.Block{
			"custom_cookies": schema.ListNestedBlock{
				Description: "Persistent custom cookies used by HTTP Client in all requests.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "cookie key",
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "cookie value",
						},
					},
				},
			},
		},
	}
}

func (p *frameworkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Read configuration from the request
	var config frameworkProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get values with environment variable fallbacks - matching SDKv2 behavior exactly
	instanceFQDN := getStringValueWithEnvFallback(config.InstanceFQDN, "JAMFPRO_INSTANCE_FQDN")
	authMethod := getStringValueWithEnvFallback(config.AuthMethod, "JAMFPRO_AUTH_METHOD")
	clientID := getStringValueWithEnvFallback(config.ClientID, "JAMFPRO_CLIENT_ID")
	clientSecret := getStringValueWithEnvFallback(config.ClientSecret, "JAMFPRO_CLIENT_SECRET")
	basicUsername := getStringValueWithEnvFallback(config.BasicAuthUsername, "JAMFPRO_BASIC_USERNAME")
	basicPassword := getStringValueWithEnvFallback(config.BasicAuthPassword, "JAMFPRO_BASIC_PASSWORD")

	// Validation - matching SDKv2 provider logic
	if instanceFQDN == "" {
		resp.Diagnostics.AddError(
			"Error getting instance name",
			"instance_name must be provided either as an environment variable (JAMFPRO_INSTANCE_FQDN) or in the Terraform configuration",
		)
		return
	}

	if authMethod == "" {
		resp.Diagnostics.AddError(
			"Error getting auth method",
			"auth_method must be provided either as an environment variable (JAMFPRO_AUTH_METHOD) or in the Terraform configuration",
		)
		return
	}

	// Validate auth method
	if authMethod != "basic" && authMethod != "oauth2" {
		resp.Diagnostics.AddError(
			"invalid auth method supplied",
			"Auth method must be 'basic' or 'oauth2'",
		)
		return
	}

	// Auth method specific validation
	switch authMethod {
	case "oauth2":
		if clientID == "" {
			resp.Diagnostics.AddError(
				"Error getting client ID",
				"client_id must be provided either as an environment variable (JAMFPRO_CLIENT_ID) or in the Terraform configuration when using oauth2 auth method",
			)
			return
		}
		if clientSecret == "" {
			resp.Diagnostics.AddError(
				"Error getting client secret",
				"client_secret must be provided either as an environment variable (JAMFPRO_CLIENT_SECRET) or in the Terraform configuration when using oauth2 auth method",
			)
			return
		}
	case "basic":
		if basicUsername == "" {
			resp.Diagnostics.AddError(
				"Error getting basic auth username",
				"basic_auth_username must be provided either as an environment variable (JAMFPRO_BASIC_USERNAME) or in the Terraform configuration when using basic auth method",
			)
			return
		}
		if basicPassword == "" {
			resp.Diagnostics.AddError(
				"Error getting basic auth password",
				"basic_auth_password must be provided either as an environment variable (JAMFPRO_BASIC_PASSWORD) or in the Terraform configuration when using basic auth method",
			)
			return
		}
	}

	// Get configuration values with proper defaults
	enableClientSDKLogs := getBoolWithDefault(config.EnableClientSDKLogs, false)
	clientSDKLogExportPath := getStringWithDefault(config.ClientSDKLogExportPath, "")
	hideSensitiveData := getBoolWithDefault(config.HideSensitiveData, true)
	loadBalancerLock := getBoolWithDefault(config.LoadBalancerLock, false)
	tokenRefreshBuffer := time.Duration(getInt64WithDefault(config.TokenRefreshBufferPeriodSeconds, 300)) * time.Second
	mandatoryRequestDelay := time.Duration(getInt64WithDefault(config.MandatoryRequestDelayMilliseconds, 100)) * time.Millisecond

	// Create logger configuration - matching SDKv2 provider
	var sugaredLogger *zap.SugaredLogger
	if enableClientSDKLogs {
		var logger *zap.Logger
		var err error

		if clientSDKLogExportPath != "" {
			// Create file logger if path specified
			config := zap.NewProductionConfig()
			config.OutputPaths = []string{clientSDKLogExportPath}
			config.ErrorOutputPaths = []string{clientSDKLogExportPath}
			logger, err = config.Build()
		} else {
			// Use development logger for console output
			logger, err = zap.NewDevelopment()
		}

		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating logger",
				fmt.Sprintf("Unable to create logger: %s", err),
			)
			return
		}

		sugaredLogger = logger.Sugar()
	} else {
		// Create no-op logger
		logger := zap.NewNop()
		sugaredLogger = logger.Sugar()
	}

	// Create bootstrap HTTP client - matching SDKv2 provider configuration
	bootstrapClient := http.Client{
		Timeout: 60 * time.Second,
	}

	// Create Jamf Pro integration based on auth method
	var jamfIntegration *jamfprointegration.Integration
	var err error

	switch authMethod {
	case "oauth2":
		jamfIntegration, err = jamfprointegration.BuildWithOAuth(
			instanceFQDN,
			sugaredLogger,
			tokenRefreshBuffer,
			clientID,
			clientSecret,
			hideSensitiveData,
			bootstrapClient,
		)
	case "basic":
		jamfIntegration, err = jamfprointegration.BuildWithBasicAuth(
			instanceFQDN,
			sugaredLogger,
			tokenRefreshBuffer,
			basicUsername,
			basicPassword,
			hideSensitiveData,
			bootstrapClient,
		)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error building jamf integration",
			fmt.Sprintf("Error: %v", err),
		)
		return
	}

	// Handle cookies - matching SDKv2 provider logic
	var cookiesList []*http.Cookie

	if loadBalancerLock {
		cookies, err := jamfIntegration.GetSessionCookies()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting session cookies",
				fmt.Sprintf("error: %v", err),
			)
			return
		}
		cookiesList = append(cookiesList, cookies...)
	}

	// Handle custom cookies
	if !config.CustomCookies.IsNull() && !config.CustomCookies.IsUnknown() {
		var customCookies []customCookieModel
		resp.Diagnostics.Append(config.CustomCookies.ElementsAs(ctx, &customCookies, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		jamfLoadBalancerCookieName := "jpro-ingress"
		for _, cookieConfig := range customCookies {
			cookieName := cookieConfig.Name.ValueString()
			cookieValue := cookieConfig.Value.ValueString()

			if cookieName == jamfLoadBalancerCookieName && loadBalancerLock {
				resp.Diagnostics.AddError(
					"Error: Conflicts in Load balancer configuration",
					"Both 'jamfpro_load_balancer_lock' and 'custom_cookies' with 'jpro-ingress' are set. Please use only one method.",
				)
				return
			}

			cookie := &http.Cookie{
				Name:  cookieName,
				Value: cookieValue,
			}
			cookiesList = append(cookiesList, cookie)
		}
	}

	// Build HTTP client configuration - exactly matching SDKv2 provider
	clientConfig := httpclient.ClientConfig{
		Integration:              jamfIntegration,
		Sugar:                    sugaredLogger,
		HideSensitiveData:        hideSensitiveData,
		TokenRefreshBufferPeriod: tokenRefreshBuffer,
		CustomCookies:            cookiesList,
		MandatoryRequestDelay:    mandatoryRequestDelay,
		RetryEligiableRequests:   false, // Forced because terraform handles concurrency
		HTTP:                     http.Client{},
	}

	httpClient, err := clientConfig.Build()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error building HTTP client",
			fmt.Sprintf("Error: %v", err),
		)
		return
	}

	// Create Jamf Pro SDK client - exactly matching SDKv2 provider
	jamfProSdk := jamfpro.Client{
		HTTP: httpClient,
	}

	versionResp, err := jamfProSdk.GetJamfProVersion()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error checking Jamf Pro version",
			fmt.Sprintf("Unable to determine Jamf Pro version: %v", err),
		)
		return
	}

	if versionResp == nil || versionResp.Version == nil || strings.TrimSpace(*versionResp.Version) == "" {
		resp.Diagnostics.AddError(
			"Error checking Jamf Pro version",
			"Jamf Pro version response was empty",
		)
		return
	}

	meetsRequirement, err := versionSupportsRequirement(*versionResp.Version, minimumJamfProVersion)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing Jamf Pro version",
			fmt.Sprintf("Unable to parse Jamf Pro version %q: %v", *versionResp.Version, err),
		)
		return
	}

	if !meetsRequirement {
		resp.Diagnostics.AddWarning(
			"Unsupported Jamf Pro Version",
			fmt.Sprintf("Your Jamf Pro instance (%s) does not meet the minimum supported version (%s). Operation against unsupported versions is not guaranteed and may result in data inconsistencies or provider errors. Continue only if you understand the risks.", *versionResp.Version, minimumJamfProVersion),
		)
		return
	}

	// Store client for use by resources and data sources
	resp.ResourceData = &jamfProSdk
	resp.DataSourceData = &jamfProSdk
}
