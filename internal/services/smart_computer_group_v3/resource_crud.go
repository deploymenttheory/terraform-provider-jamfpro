package smart_computer_group_v3

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	frameworkcrud "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/framework_crud"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *smartComputerGroupV3FrameworkResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var data smartComputerGroupV3ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() || !r.configured(&resp.Diagnostics) {
		return
	}
	opCtx, cancel := frameworkcrud.HandleTimeout(ctx, data.Timeouts.Create, defaultTimeout, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	input, diags := constructResource(opCtx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, fallback, err := r.createGroup(opCtx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Smart Computer Group V3", err.Error())
		return
	}
	data.ID = types.StringValue(id)
	r.recordWrite(ctx, opCtx, &data, input, fallback, fallback && input.Description != "", &resp.State, &resp.Diagnostics)
}

func (r *smartComputerGroupV3FrameworkResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var data smartComputerGroupV3ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() || !r.configured(&resp.Diagnostics) {
		return
	}
	opCtx, cancel := frameworkcrud.HandleTimeout(ctx, data.Timeouts.Read, defaultTimeout, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	group, err := r.client.GetSmartComputerGroupByIDV3WithContext(opCtx, data.ID.ValueString())
	if isNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Smart Computer Group V3", err.Error())
		return
	}
	if group == nil {
		resp.Diagnostics.AddError("Error Reading Smart Computer Group V3", errMissingResponse.Error())
		return
	}
	resp.Diagnostics.Append(state(ctx, &data, group)...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *smartComputerGroupV3FrameworkResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan, prior smartComputerGroupV3ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() || !r.configured(&resp.Diagnostics) {
		return
	}
	opCtx, cancel := frameworkcrud.HandleTimeout(ctx, plan.Timeouts.Update, defaultTimeout, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	input, diags := constructResource(opCtx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = prior.ID
	fallback, err := r.updateGroup(opCtx, prior.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Smart Computer Group V3", err.Error())
		return
	}
	patchDescription := fallback && plan.Description.ValueString() != prior.Description.ValueString()
	r.recordWrite(ctx, opCtx, &plan, input, fallback, patchDescription, &resp.State, &resp.Diagnostics)
}

func (r *smartComputerGroupV3FrameworkResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var data smartComputerGroupV3ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() || !r.configured(&resp.Diagnostics) {
		return
	}
	opCtx, cancel := frameworkcrud.HandleTimeout(ctx, data.Timeouts.Delete, defaultTimeout, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	err := r.client.DeleteSmartComputerGroupByIDV3WithContext(opCtx, data.ID.ValueString())
	if err != nil && !isNotFound(err) {
		resp.Diagnostics.AddError("Error Deleting Smart Computer Group V3", err.Error())
		return
	}
	if err == nil {
		r.waitUntilGone(opCtx, data.ID.ValueString())
	}
	resp.State.RemoveResource(ctx)
}

func (r *smartComputerGroupV3FrameworkResource) configured(diags *diag.Diagnostics) bool {
	if r.client == nil {
		diags.AddError("Unconfigured Smart Computer Group V3", "The provider has not supplied a Jamf Pro SDK client.")
		return false
	}
	return true
}

func (r *smartComputerGroupV3FrameworkResource) recordWrite(
	ctx, operationCtx context.Context,
	data *smartComputerGroupV3ResourceModel,
	input *jamfpro.ResourceSmartComputerGroupV3,
	fallback, patchDescription bool,
	result *tfsdk.State,
	diags *diag.Diagnostics,
) {
	// Record the identifier before any additional calls. A failed description
	// patch must leave the created group manageable, not orphan it in Jamf Pro.
	diags.Append(result.Set(ctx, data)...)
	if diags.HasError() {
		return
	}
	if fallback {
		diags.AddWarning("Jamf Pro Criterion Validation Workaround",
			"The v3 endpoint refused a criterion or operator. The complete group definition was saved through "+
				"the Classic API; description changes are saved separately through the unified groups API.")
	}
	got, err := r.completeWrite(operationCtx, data.ID.ValueString(), input, patchDescription)
	if got != nil {
		diags.Append(state(ctx, data, got)...)
		if !diags.HasError() {
			diags.Append(result.Set(ctx, data)...)
		}
	}
	if err != nil {
		diags.AddError("Smart Computer Group Write Incomplete",
			fmt.Sprintf("Group %s is recorded in Terraform state, but verifying the write or saving its description failed: %v. "+
				"Refresh and apply again to finish the operation.", data.ID.ValueString(), err))
	}
}
