package jamfproprovider

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DepartmentResource{}
var _ resource.ResourceWithImportState = &DepartmentResource{}

func NewDepartmentResource() resource.Resource {
	return &DepartmentResource{}
}

// DepartmentResource defines the resource implementation.
type DepartmentResource struct {
	client *jamfpro.Client
}

// DepartmentResourceModel describes the resource data model.
type DepartmentResourceModel struct {
	Id   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Href types.String `tfsdk:"href"`
}

func (r *DepartmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "jamfpro_department"
}

func (r *DepartmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Jamf Pro Department",
		Attributes: map[string]schema.Attribute{
			"id": schema.NumberAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the department.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique name of the Jamf Pro department.",
			},
			"href": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The URL link for the department.",
			},
		},
	}
}

func (r *DepartmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := req.ProviderData.(*jamfpro.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *jamfpro.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create handles the creation of a department resource in Jamf Pro.
func (r *DepartmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DepartmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	department, err := r.client.CreateDepartment(data.Name.ValueString())

	if err != nil {
		tflog.Error(ctx, "Failed to create department", map[string]interface{}{"error": err.Error()})
		resp.Diagnostics.AddError("Failed to create department", err.Error())
		return
	}

	// Convert the data responses to strings for consistent data handling in Terraform.
	data.Id = types.Int64Value(int64(department.Id))
	data.Name = types.StringValue(department.Name)
	data.Href = types.StringValue(department.Href)

	tflog.Trace(ctx, "Successfully created a department resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read retrieves details of a department resource from Jamf Pro based on its ID or name.
func (r *DepartmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DepartmentResourceModel
	// Appending diagnostics data (logs, errors, etc.) from the current state of the resource
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	department, err := r.client.GetDepartmentByID(int(data.Id.ValueInt64()))
	if err != nil {
		tflog.Warn(ctx, "Failed to fetch department by ID. Attempting to fetch by name.", map[string]interface{}{"reason": err.Error()})
		department, err = r.client.GetDepartmentByName(data.Name.ValueString())
		if err != nil {
			tflog.Error(ctx, "Failed to read department", map[string]interface{}{"error": err.Error()})
			resp.Diagnostics.AddError("Failed to read department", err.Error())
			return
		}
	}
	// Update the data model with the retrieved department details
	data.Name = types.StringValue(department.Name)
	data.Href = types.StringValue(department.Href)
	data.Id = types.Int64Value(int64(department.Id))

	// Append the updated data to the terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update modifies an existing department resource in Jamf Pro.
func (r *DepartmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Initializing a model to store the resource data
	var data DepartmentResourceModel

	// Appending diagnostics data (logs, errors, etc.) from the planned state of the resource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// If there's any error in the diagnostics, exit early
	if resp.Diagnostics.HasError() {
		return
	}

	// Attempt to update the department details using the ID and name from the planned state
	updatedDepartment, err := r.client.UpdateDepartmentByID(int(data.Id.ValueInt64()), data.Name.ValueString())

	// If there's an error in updating the department by ID...
	if err != nil {
		// Log a warning about the failure to update by ID and try updating by name
		tflog.Warn(ctx, "Failed to update department by ID. Attempting to update by name.", map[string]interface{}{"reason": err.Error()})
		updatedDepartment, err = r.client.UpdateDepartmentByName(data.Name.ValueString(), data.Name.ValueString())

		// If there's still an error in updating the department by name...
		if err != nil {
			// Log an error and append the error to diagnostics
			tflog.Error(ctx, "Failed to update department", map[string]interface{}{"error": err.Error()})
			resp.Diagnostics.AddError("Failed to update department", err.Error())
			return
		}
	}

	// Update the data model with the updated department details
	data.Name = types.StringValue(updatedDepartment.Name)   // Set the name from the updated department
	data.Href = types.StringValue(updatedDepartment.Href)   // Set the href from the updated department
	data.Id = types.Int64Value(int64(updatedDepartment.Id)) // Set the ID from the updated department

	// Append the updated data to the diagnostics
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete removes a department resource from Jamf Pro.
func (r *DepartmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Initializing a model to store the resource data
	var data DepartmentResourceModel

	// Appending diagnostics data (logs, errors, etc.) from the current state of the resource
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// If there's any error in the diagnostics, exit early
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the Int64Value ID to an integer
	departmentID := int(data.Id.ValueInt64())

	// Attempt to delete the department using the ID from the current state
	err := r.client.DeleteDepartmentByID(departmentID)

	// If there's an error in deleting the department by ID...
	if err != nil {
		// Log a warning about the failure to delete by ID and try deleting by name
		tflog.Warn(ctx, "Failed to delete department by ID. Attempting to delete by name.", map[string]interface{}{"reason": err.Error()})
		err = r.client.DeleteDepartmentByName(data.Name.ValueString())

		// If there's still an error in deleting the department by name...
		if err != nil {
			// Log an error and append the error to diagnostics
			tflog.Error(ctx, "Failed to delete department", map[string]interface{}{"error": err.Error()})
			resp.Diagnostics.AddError("Failed to delete department", err.Error())
			return
		}
	}
}

func (r *DepartmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
