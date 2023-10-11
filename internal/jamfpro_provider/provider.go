package jamfpro_provider

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure JamfProProvider satisfies various provider interfaces.
var _ provider.Provider = &JamfProProvider{}

// constants
const defaultMaxConcurrentRequests = 10

// JamfProProvider defines the provider implementation.
type JamfProProvider struct {
	version               string
	InstanceName          string `tfsdk:"instance_name"`
	ClientID              string `tfsdk:"client_id"`
	ClientSecret          string `tfsdk:"client_secret"`
	DebugMode             bool   `tfsdk:"debug_mode"`
	MaxConcurrentRequests *int   `tfsdk:"max_concurrent_requests"`
}

// JamfProProviderModel describes the provider data model.
type JamfProProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
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
				Description: "Client secret for authentication.",
			},
			"debug_mode": schema.BoolAttribute{
				Required:    true,
				Description: "Enable or disable debug mode for verbose logging.",
			},
			"max_concurrent_requests": schema.NumberAttribute{
				Optional:    true,
				Description: "Maximum number of simultaneous requests allowed.",
			},
			// Add other attributes you wish to add from the http client
		},
	}
}

func (p *JamfProProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Starting the configuration process
	tflog.Info(ctx, "Configuring the JamfProProvider...")

	// Extract configuration into the JamfProProvider directly
	resp.Diagnostics.Append(req.Config.Get(ctx, p)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error occurred while extracting configuration.")
		return
	}

	// Handle optional fields
	maxConcurrentRequests := defaultMaxConcurrentRequests
	if p.MaxConcurrentRequests != nil {
		maxConcurrentRequests = *p.MaxConcurrentRequests
	}

	// Configuration for the jamfpro
	config := jamfpro.Config{
		InstanceName:          p.InstanceName,
		DebugMode:             p.DebugMode,
		Logger:                jamfpro.NewDefaultLogger(),
		MaxConcurrentRequests: maxConcurrentRequests,
		TokenLifespan:         30,
		BufferPeriod:          5,
		ClientID:              p.ClientID,
		ClientSecret:          p.ClientSecret,
	}

	// Create a new jamfpro client instance
	client := jamfpro.NewClient(config)
	if client == nil {
		tflog.Error(ctx, "JamfPro client is nil after initialization. Configuration failed.")
		resp.Diagnostics.AddError("JamfPro Client Initialization Error", "JamfPro client is nil after initialization. This may be due to an internal issue in the NewClient function.")
		return
	}

	// Assign the jamfpro client to the provider
	resp.DataSourceData = client
	resp.ResourceData = client

	// Check if the client was correctly assigned
	if resp.DataSourceData == nil || resp.ResourceData == nil {
		tflog.Error(ctx, "resp.DataSourceData or resp.ResourceData is nil after assignment.")
		resp.Diagnostics.AddError("Assignment Error", "resp.DataSourceData or resp.ResourceData is nil after assignment.")
		return
	}

	tflog.Info(ctx, "JamfProProvider configuration completed successfully.")
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

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &JamfProProvider{
			version: version,
		}
	}
}
