package service_discovery_enrollment_well_known_settings

import (
	"context"
	"fmt"
	"time"

	frameworkCrud "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/framework_crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (r *serviceDiscoveryEnrollmentWellKnownSettingsFrameworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serviceDiscoveryEnrollmentWellKnownSettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	payload, diags := constructPayload(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateServiceDiscoveryEnrollmentWellKnownSettingsV1(*payload); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Service Discovery Enrollment Well-Known Settings",
			fmt.Sprintf("Could not update service discovery enrollment well-known settings: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(serviceDiscoveryEnrollmentWellKnownSettingsSingletonID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	container := &frameworkCrud.CreateResponseContainer{CreateResponse: resp}
	opts := frameworkCrud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
	opts.ResourceTypeName = ResourceName

	if err := frameworkCrud.ReadWithRetry(ctx, r.Read, readReq, container, opts); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Service Discovery Enrollment Well-Known Settings After Create",
			fmt.Sprintf("Could not refresh state after create: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

func (r *serviceDiscoveryEnrollmentWellKnownSettingsFrameworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serviceDiscoveryEnrollmentWellKnownSettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	response, err := r.client.GetServiceDiscoveryEnrollmentWellKnownSettingsV1()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Service Discovery Enrollment Well-Known Settings",
			fmt.Sprintf("Could not read service discovery enrollment well-known settings: %s", err),
		)
		return
	}

	applyResponse(&state, response)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

func (r *serviceDiscoveryEnrollmentWellKnownSettingsFrameworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan serviceDiscoveryEnrollmentWellKnownSettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	payload, diags := constructPayload(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateServiceDiscoveryEnrollmentWellKnownSettingsV1(*payload); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Service Discovery Enrollment Well-Known Settings",
			fmt.Sprintf("Could not update service discovery enrollment well-known settings: %s", err),
		)
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	container := &frameworkCrud.UpdateResponseContainer{UpdateResponse: resp}
	opts := frameworkCrud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = ResourceName

	if err := frameworkCrud.ReadWithRetry(ctx, r.Read, readReq, container, opts); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Service Discovery Enrollment Well-Known Settings After Update",
			fmt.Sprintf("Could not refresh state after update: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

func (r *serviceDiscoveryEnrollmentWellKnownSettingsFrameworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
