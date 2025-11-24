package static_computer_group

import (
	"context"
	"fmt"
	"time"

	frameworkCrud "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/framework_crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create creates a new static computer group resource in Jamf Pro.
func (r *staticComputerGroupFrameworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object staticComputerGroupResourceModel

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

	staticGroup, constructDiags := constructResource(&object)
	resp.Diagnostics.Append(constructDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdGroup, err := r.client.CreateStaticComputerGroupV2(*staticGroup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Static Computer Group",
			fmt.Sprintf("Could not create static computer group: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	object.ID = types.StringValue(createdGroup.ID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &frameworkCrud.CreateResponseContainer{CreateResponse: resp}

	opts := frameworkCrud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
	opts.ResourceTypeName = ResourceName

	err = frameworkCrud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Static Computer Group After Create",
			fmt.Sprintf("Could not read static computer group after creation: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read reads the current state of a static computer group resource from Jamf Pro.
func (r *staticComputerGroupFrameworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object staticComputerGroupResourceModel

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

	staticGroup, err := r.client.GetStaticComputerGroupByIDV2(object.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Static Computer Group V2",
			fmt.Sprintf("Could not read static computer group ID %s: %s", object.ID.ValueString(), err.Error()),
		)
		return
	}

	stateDiags := state(&object, staticGroup)
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

// Update updates an existing static computer group resource in Jamf Pro.
func (r *staticComputerGroupFrameworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan staticComputerGroupResourceModel
	var state staticComputerGroupResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	staticGroup, constructDiags := constructResource(&plan)
	resp.Diagnostics.Append(constructDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateStaticComputerGroupByIDV2(state.ID.ValueString(), *staticGroup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Static Computer Group V2",
			fmt.Sprintf("Could not update static computer group: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	plan.ID = state.ID
	plan.Timeouts = state.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &frameworkCrud.UpdateResponseContainer{UpdateResponse: resp}

	opts := frameworkCrud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = ResourceName

	err = frameworkCrud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Static Computer Group After Update",
			fmt.Sprintf("Could not read static computer group after update: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete deletes a static computer group resource from Jamf Pro.
func (r *staticComputerGroupFrameworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object staticComputerGroupResourceModel

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

	err := r.client.DeleteStaticComputerGroupByIDV2(object.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Static Computer Group V2",
			fmt.Sprintf("Could not delete static computer group: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
