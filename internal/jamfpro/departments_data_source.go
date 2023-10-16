package jamfpro

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DepartmentDataSource{}

func NewDepartmentDataSource() datasource.DataSource {
	return &DepartmentDataSource{}
}

func (d *DepartmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "jamfpro_department"
}

// DepartmentDataSource defines the data source implementation.
type DepartmentDataSource struct {
	client *jamfpro.Client // This points to the http_client that's part of the jamfpro package
}

// DepartmentDataSourceModel describes the data source data model.
type DepartmentDataSourceModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Href types.String `tfsdk:"href"`
}

// DepartmentDataSource provides information about a specific department in Jamf Pro.
// It can fetch department details using either the department's unique Name or its Id.
// The Name attribute is prioritized for fetching if provided. Otherwise, the Id is used.
func (d *DepartmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The unique identifier of the department.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The unique name of the jamf pro department.",
			},
			"href": schema.StringAttribute{
				Computed:    true,
				Description: "The URL link for the department.",
			},
		},
	}
}

// Configure sets up the data source for subsequent operations.
// It checks if the provider data (which should contain configuration details) is available,
// and then attempts to cast it to a specific type (in this case, a Jamf Pro client).
func (d *DepartmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	println("-----func (d *DepartmentDataSource) Configure-----")
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jamfpro.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *jamfpro.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *DepartmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DepartmentDataSourceModel

	// Extract the Terraform configuration data into the DepartmentDataSourceModel.
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var department *jamfpro.Department
	var err error
	// Check if the Name attribute is provided in the Terraform configuration.
	// If it's not, use the Id to fetch the department details.
	if data.Name.IsNull() || data.Name.ValueString() == "" {
		department, err = d.client.GetDepartmentByID(int(data.Id.ValueInt64()))
	} else {
		department, err = d.client.GetDepartmentByName(data.Name.ValueString())
	}

	// If there's an error fetching the department details, add the error to the diagnostics.
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch department", err.Error())
		return
	}

	// Populate the DepartmentDataSourceModel with the fetched department details.
	data.Id = basetypes.NewInt64Value(int64(department.Id))
	data.Name = basetypes.NewStringValue(department.Name)
	data.Href = basetypes.NewStringValue(department.Href)

	// Log the successful reading of the data source.
	tflog.Trace(ctx, "Successfully read the department data source")

	// Save the populated DepartmentDataSourceModel into the Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
