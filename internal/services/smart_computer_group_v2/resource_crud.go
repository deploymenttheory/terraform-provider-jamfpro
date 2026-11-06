package smart_computer_group_v2

import (
	"context"
	"fmt"
	"time"

	frameworkCrud "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/framework_crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create creates a new smart computer group resource in Jamf Pro.
func (r *smartComputerGroupV2FrameworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object smartComputerGroupV2ResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if resp.Diagnostics.HasError() {
		return
	}

	smartGroup, constructDiags := constructResource(ctx, &object)
	resp.Diagnostics.Append(constructDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdGroup, err := r.client.CreateSmartComputerGroupV2WithContext(ctx, *smartGroup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Smart Computer Group",
			fmt.Sprintf("Could not create smart computer group: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	object.ID = types.StringValue(createdGroup.ID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actual, err := r.waitForWrite(ctx, object.ID.ValueString(), smartGroup)
	if actual != nil {
		resp.Diagnostics.Append(state(ctx, &object, object.ID.ValueString(), actual)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Smart Computer Group After Create",
			fmt.Sprintf("Could not read smart computer group after creation: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read reads the current state of a smart computer group resource from Jamf Pro.
func (r *smartComputerGroupV2FrameworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object smartComputerGroupV2ResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s", ResourceName, object.ID.ValueString()))

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	resourceID := object.ID.ValueString()
	smartGroup, err := r.client.GetSmartComputerGroupByIDV2WithContext(ctx, resourceID)
	if err != nil {
		if frameworkCrud.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Smart Computer Group V2",
			fmt.Sprintf("Could not read smart computer group ID %s: %s", resourceID, err.Error()),
		)
		return
	}

	stateDiags := state(ctx, &object, resourceID, smartGroup)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update updates an existing smart computer group resource in Jamf Pro.
func (r *smartComputerGroupV2FrameworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan smartComputerGroupV2ResourceModel
	var prior smartComputerGroupV2ResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, prior.ID.ValueString()))

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	smartGroup, constructDiags := constructResource(ctx, &plan)
	resp.Diagnostics.Append(constructDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateSmartComputerGroupByIDV2WithContext(ctx, prior.ID.ValueString(), *smartGroup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Smart Computer Group V2",
			fmt.Sprintf("Could not update smart computer group: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	plan.ID = prior.ID

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actual, err := r.waitForWrite(ctx, prior.ID.ValueString(), smartGroup)
	if actual != nil {
		resp.Diagnostics.Append(state(ctx, &plan, prior.ID.ValueString(), actual)...)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Smart Computer Group After Update",
			fmt.Sprintf("Could not read smart computer group after update: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, prior.ID.ValueString()))
}

// Delete deletes a smart computer group resource from Jamf Pro.
func (r *smartComputerGroupV2FrameworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object smartComputerGroupV2ResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.DeleteSmartComputerGroupByIDV2WithContext(ctx, object.ID.ValueString())

	if err != nil && !frameworkCrud.IsNotFound(err) {
		resp.Diagnostics.AddError(
			"Error Deleting Smart Computer Group V2",
			fmt.Sprintf("Could not delete smart computer group: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
