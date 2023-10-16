package jamfpro

import (
	"context"
	"strconv"

	"os"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure JamfProProvider satisfies various provider interfaces.
var _ provider.Provider = (*JamfProProvider)(nil)

// JamfProProvider defines the provider implementation.
type JamfProProvider struct {
	version string
}

// JamfPro ProviderModel maps provider schema data to a Go type.
type JamfProProviderModel struct {
	InstanceName types.String `tfsdk:"instance_name"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	DebugMode    types.Bool   `tfsdk:"debug_mode"`
}

func New(version string) func() provider.Provider {
	tflog.Info(context.Background(), "Initializing JamfPro provider")
	return func() provider.Provider {
		return &JamfProProvider{
			version: version,
		}
	}
}

func (p *JamfProProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jamfpro"
	resp.Version = p.version
}

// Schema defines the configuration attributes for the http_client within the JamfPro provider.
func (p *JamfProProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"instance_name": schema.StringAttribute{
				Required:    true,
				Description: "The Jamf Pro instance name. For mycompany.jamfcloud.com, define mycompany in this field.",
			},
			"client_id": schema.StringAttribute{
				Required:    true,
				Description: "Client ID for authentication.",
			},
			"client_secret": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Client secret for authentication.",
			},
			"debug_mode": schema.BoolAttribute{
				Required:    true,
				Description: "Enable or disable debug mode for verbose logging.",
			},
		},
	}
}

func (p *JamfProProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	println("hello")
	// Retrieve provider data from configuration
	var config JamfProProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide a instance name to the provider
	var instanceName string
	if config.InstanceName.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the Jamf Pro Server",
			"Cannot use unknown value as instance_name",
		)
	} else {
		if config.InstanceName.IsNull() {
			instanceName = os.Getenv("JAMFPRO_PROVIDER_INSTANCE_NAME")
		} else {
			instanceName = config.InstanceName.ValueString()
		}
		if instanceName == "" {
			resp.Diagnostics.AddError(
				"Unable to find instance_name",
				"instance_name cannot be an empty string. Either set it in the configuration or use the JAMFPRO_PROVIDER_INSTANCE_NAME environment variable.",
			)
		}
	}

	// User must provide a client id to the provider
	var clientID string
	if config.ClientID.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the Jamf Pro Server",
			"Cannot use unknown value as client_id",
		)
	} else {
		if config.ClientID.IsNull() {
			clientID = os.Getenv("JAMFPRO_PROVIDER_CLIENTID")
		} else {
			clientID = config.ClientID.ValueString()
		}
		if clientID == "" {
			resp.Diagnostics.AddError(
				"Unable to find client_id",
				"client_id cannot be an empty string. Either set it in the configuration or use the JAMFPRO_PROVIDER_CLIENTID environment variable.",
			)
		}
	}

	// User must provide a client secret to the provider
	var clientSecret string
	if config.ClientSecret.IsUnknown() {
		// Cannot connect to PingFederate with an unknown value
		resp.Diagnostics.AddError(
			"Unable to connect to the Jamf Pro Server",
			"Cannot use unknown value as client_secret",
		)
	} else {
		if config.ClientSecret.IsNull() {
			clientSecret = os.Getenv("JAMFPRO_PROVIDER_CLIENT_SECRET")
		} else {
			clientSecret = config.ClientSecret.ValueString()
		}
		if clientSecret == "" {
			resp.Diagnostics.AddError(
				"Unable to find client_secret",
				"client_secret cannot be an empty string. Either set it in the configuration or use the JAMFPRO_PROVIDER_CLIENT_SECRET environment variable.",
			)
		}
	}

	// Optional attributes
	var debugMode bool
	var err error
	if !config.DebugMode.IsUnknown() && !config.DebugMode.IsNull() {
		debugMode = config.DebugMode.ValueBool()
	} else {
		debugMode, err = strconv.ParseBool(os.Getenv("JAMFPRO_PROVIDER_DEBUG_MODE"))
		if err != nil {
			debugMode = false
			tflog.Info(ctx, "Failed to parse boolean from 'JAMFPRO_PROVIDER_DEBUG_MODE' environment variable, defaulting 'debug_mode' to false")
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a configuration for the JamfPro client based on the provider's configuration
	jamfProConfig := jamfpro.Config{
		InstanceName: instanceName,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		DebugMode:    debugMode,
		// Add other necessary fields as required
	}

	// Initialize the JamfPro client
	jamfProClient, err := jamfpro.NewClient(jamfProConfig)
	if err != nil {
		tflog.Error(ctx, "Error while initializing the JamfPro client", map[string]interface{}{"error": err.Error()})
		resp.Diagnostics.AddError("Failed to create JamfPro client", err.Error())
		return
	} else {
		tflog.Info(ctx, "Successfully initialized the JamfPro client", map[string]interface{}{"success": true})
	}

	// Validate the client's OAuth credentials
	if jamfProClient.HTTP.GetOAuthCredentials().ClientID == "" || jamfProClient.HTTP.GetOAuthCredentials().ClientSecret == "" {
		resp.Diagnostics.AddError("OAuth credentials error", "OAuth credentials (ClientID and ClientSecret) must be provided")
		return
	}

	// Store the JamfPro client in the response so it's available to resources and data sources
	resp.ResourceData = jamfProClient
	resp.DataSourceData = jamfProClient

	tflog.Info(ctx, "Configured JamfPro client", map[string]interface{}{"success": true})
}

func (p *JamfProProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDepartmentResource,
	}
}

func (p *JamfProProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDepartmentDataSource, // Use the NewDepartmentDataSource function here
	}
}
