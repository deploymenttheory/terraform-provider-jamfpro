package cloud_distribution_point

import (
	"context"
	"fmt"
	"strings"
	"time"

	frameworkCrud "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/framework_crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create ensures the singleton Cloud Distribution Point converges to the desired
// CDN type without requiring manual imports. The logic is:
//  1. Read the current Jamf configuration.
//  2. If the CDN type already matches, issue an UpdateCloudDistributionPointV1 to update mutable fields.
//  3. If Jamf is effectively unconfigured (NONE) or the resource has not been
//     touched before, POST the desired configuration.
//  4. Otherwise delete the existing CDN (Jamf only allows one) and recreate it.
//
// After any CreateCloudDistributionPointV1 we immediately issue a follow-up UpdateCloudDistributionPointV1 when `master = false`,
// because Jamf forces `master=true` during creation even if false was requested.
func (r *cloudDistributionPointFrameworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan cloudDistributionPointResourceModel

	tflog.Debug(ctx, "Starting creation of resource", map[string]any{"resource": ResourceName})

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	payload, diags := constructCloudDistributionPointPayload(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	existingConfig, err := r.client.GetCloudDistributionPointV1()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Inspecting Cloud Distribution Point",
			fmt.Sprintf("Could not check existing cloud distribution point configuration: %s", err),
		)
		return
	}

	desiredType := strings.ToUpper(strings.TrimSpace(plan.CdnType.ValueString()))
	currentType := ""
	if existingConfig != nil {
		currentType = strings.ToUpper(strings.TrimSpace(existingConfig.CdnType))
	}

	plan.ID = types.StringValue(cloudDistributionPointSingletonID)

	switch {
	case existingConfig != nil && currentType == desiredType:
		if _, err := r.client.UpdateCloudDistributionPointV1(payload); err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Cloud Distribution Point",
				fmt.Sprintf("Could not update existing cloud distribution point: %s", err),
			)
			return
		}

		plan.InventoryID = stringValueOrNull(existingConfig.InventoryId)
	case existingConfig == nil || currentType == cdnTypeNone:
		created, err := r.client.CreateCloudDistributionPointV1(payload)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating Cloud Distribution Point",
				fmt.Sprintf("Could not create cloud distribution point: %s", err),
			)
			return
		}

		plan.InventoryID = stringValueOrNull(created.InventoryId)

		if !plan.Master.ValueBool() {
			if _, err := r.client.UpdateCloudDistributionPointV1(payload); err != nil {
				resp.Diagnostics.AddError(
					"Error Setting Master Flag",
					fmt.Sprintf("Could not update master flag after creation: %s", err),
				)
				return
			}
		}
	default:
		if err := r.client.DeleteCloudDistributionPointV1(); err != nil {
			resp.Diagnostics.AddError(
				"Error Resetting Cloud Distribution Point",
				fmt.Sprintf("Could not delete existing cloud distribution point configuration: %s", err),
			)
			return
		}

		created, err := r.client.CreateCloudDistributionPointV1(payload)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Creating Cloud Distribution Point",
				fmt.Sprintf("Could not create cloud distribution point: %s", err),
			)
			return
		}

		plan.InventoryID = stringValueOrNull(created.InventoryId)

		if !plan.Master.ValueBool() {
			if _, err := r.client.UpdateCloudDistributionPointV1(payload); err != nil {
				resp.Diagnostics.AddError(
					"Error Setting Master Flag",
					fmt.Sprintf("Could not update master flag after creation: %s", err),
				)
				return
			}
		}
	}

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
			"Error Reading Cloud Distribution Point After Create",
			fmt.Sprintf("Could not refresh state after create: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Finished Create Method", map[string]any{"resource": ResourceName})
}

func (r *cloudDistributionPointFrameworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state cloudDistributionPointResourceModel

	tflog.Debug(ctx, "Starting Read method", map[string]any{"resource": ResourceName})

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	savedPassword := state.Password
	savedPrivateKey := state.PrivateKey

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	config, err := r.client.GetCloudDistributionPointV1()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Cloud Distribution Point",
			fmt.Sprintf("Could not read cloud distribution point: %s", err),
		)
		return
	}

	uploadCap, err := r.client.GetCloudDistributionPointUploadCapabilityV1()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Cloud Distribution Point Upload Capability",
			fmt.Sprintf("Could not read upload capability: %s", err),
		)
		return
	}

	state.applyResponse(config, uploadCap)
	state.Password = savedPassword
	state.PrivateKey = savedPrivateKey

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Finished Read Method", map[string]any{"resource": ResourceName})
}

func (r *cloudDistributionPointFrameworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan cloudDistributionPointResourceModel
	var state cloudDistributionPointResourceModel

	tflog.Debug(ctx, "Starting Update method", map[string]any{"resource": ResourceName})

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	payload, diags := constructCloudDistributionPointPayload(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.UpdateCloudDistributionPointV1(payload); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Cloud Distribution Point",
			fmt.Sprintf("Could not update cloud distribution point: %s", err),
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
			"Error Reading Cloud Distribution Point After Update",
			fmt.Sprintf("Could not refresh state after update: %s", err),
		)
		return
	}

	tflog.Debug(ctx, "Finished Update Method", map[string]any{"resource": ResourceName})
}

// Delete removes the Cloud Distribution Point configuration from Jamf Pro.
// If no configuration exists, the delete is a no-op.
func (r *cloudDistributionPointFrameworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state cloudDistributionPointResourceModel

	tflog.Debug(ctx, "Starting deletion of resource", map[string]any{"resource": ResourceName})

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := frameworkCrud.HandleTimeout(ctx, state.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	config, err := r.client.GetCloudDistributionPointV1()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Inspecting Cloud Distribution Point Before Delete",
			fmt.Sprintf("Could not read current cloud distribution point configuration: %s", err),
		)
		return
	}

	if config == nil || strings.EqualFold(strings.TrimSpace(config.CdnType), cdnTypeNone) {
		resp.State.RemoveResource(ctx)
		tflog.Debug(ctx, "Skipped delete; cloud distribution point already absent", map[string]any{"resource": ResourceName})
		return
	}

	if err := r.client.DeleteCloudDistributionPointV1(); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Cloud Distribution Point",
			fmt.Sprintf("Could not delete cloud distribution point: %s", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Debug(ctx, "Finished Delete Method", map[string]any{"resource": ResourceName})
}
