package adcs_settings

import (
	"context"
	"fmt"
	"time"

	frameworkCrud "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/framework_crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (r *adcsSettingsFrameworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object adcsSettingsResourceModel

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

	payload, constructDiags := constructResource(&object)
	resp.Diagnostics.Append(constructDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateAdcsSettingsV1(payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating AD CS Settings",
			fmt.Sprintf("Could not create AD CS settings: %s", err),
		)
		return
	}

	object.ID = types.StringValue(created.ID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &frameworkCrud.CreateResponseContainer{CreateResponse: resp}

	opts := frameworkCrud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
	opts.ResourceTypeName = ResourceName

	if err := frameworkCrud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading AD CS Settings After Create",
			fmt.Sprintf("Could not refresh AD CS settings after create: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

func (r *adcsSettingsFrameworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object adcsSettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if object.ID.IsNull() || object.ID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing ID",
			"Cannot read AD CS settings without an ID in state.",
		)
		return
	}

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	settings, err := r.client.GetAdcsSettingsByIDV1(object.ID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			tflog.Warn(ctx, fmt.Sprintf("AD CS settings %s not found, removing from state", object.ID.ValueString()))
			resp.Diagnostics.AddWarning(
				"AD CS Settings Not Found",
				fmt.Sprintf("AD CS settings with ID %s no longer exists in Jamf Pro and will be recreated on the next apply.", object.ID.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading AD CS Settings",
			fmt.Sprintf("Could not read AD CS settings ID %s: %s", object.ID.ValueString(), err),
		)
		return
	}

	stateDiags := state(&object, settings)
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

func (r *adcsSettingsFrameworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan adcsSettingsResourceModel
	var stateModel adcsSettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	payload, constructDiags := constructResource(&plan)
	resp.Diagnostics.Append(constructDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateAdcsSettingsByIDV1(stateModel.ID.ValueString(), payload); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating AD CS Settings",
			fmt.Sprintf("Could not update AD CS settings ID %s: %s", stateModel.ID.ValueString(), err),
		)
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &frameworkCrud.UpdateResponseContainer{UpdateResponse: resp}

	opts := frameworkCrud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = ResourceName

	if err := frameworkCrud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading AD CS Settings After Update",
			fmt.Sprintf("Could not refresh AD CS settings after update: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, stateModel.ID.ValueString()))
}

func (r *adcsSettingsFrameworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object adcsSettingsResourceModel

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

	err := r.client.DeleteAdcsSettingsByIDV1(object.ID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			tflog.Warn(ctx, fmt.Sprintf("AD CS settings %s already removed", object.ID.ValueString()))
		} else {
			resp.Diagnostics.AddError(
				"Error Deleting AD CS Settings",
				fmt.Sprintf("Could not delete AD CS settings ID %s: %s", object.ID.ValueString(), err),
			)
			return
		}
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
