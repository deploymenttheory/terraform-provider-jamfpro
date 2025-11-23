package static_computer_group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the resource implements the ResourceWithUpgradeState interface
var _ resource.ResourceWithUpgradeState = &staticComputerGroupFrameworkResource{}

// UpgradeState returns the state upgraders for the resource
func (r *staticComputerGroupFrameworkResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
					"is_smart": schema.BoolAttribute{
						Computed: true,
					},
					"site_id": schema.Int64Attribute{
						Optional: true,
						Computed: true,
					},
					"assigned_computer_ids": schema.ListAttribute{
						Optional:    true,
						ElementType: types.Int64Type,
					},
				},
				Blocks: map[string]schema.Block{
					"timeouts": timeouts.Block(ctx, timeouts.Opts{
						Create: true,
						Read:   true,
						Update: true,
						Delete: true,
					}),
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				upgradeStateV0toV1(ctx, req, resp)
			},
		},
	}
}

// upgradeStateV0toV1 migrates state from version 0 (SDK v2 schema) to version 1 (Framework schema)
// This handles the migration from the old static_computer_group resource to the new v2 version
func upgradeStateV0toV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	type timeoutsV0 struct {
		Create types.String `tfsdk:"create"`
		Read   types.String `tfsdk:"read"`
		Update types.String `tfsdk:"update"`
		Delete types.String `tfsdk:"delete"`
	}

	type modelV0 struct {
		ID                  types.String `tfsdk:"id"`
		Name                types.String `tfsdk:"name"`
		IsSmart             types.Bool   `tfsdk:"is_smart"`
		SiteID              types.Int64  `tfsdk:"site_id"`
		AssignedComputerIDs types.List   `tfsdk:"assigned_computer_ids"`
		Timeouts            *timeoutsV0  `tfsdk:"timeouts"`
	}

	var priorStateData modelV0

	resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upgradedStateData := staticComputerGroupResourceModel{
		ID:          priorStateData.ID,
		Name:        priorStateData.Name,
		Description: types.StringNull(),
	}

	if !priorStateData.SiteID.IsNull() && !priorStateData.SiteID.IsUnknown() {
		siteIDInt := priorStateData.SiteID.ValueInt64()
		if siteIDInt == -1 {
			upgradedStateData.SiteID = types.StringValue("-1")
		} else {
			upgradedStateData.SiteID = types.StringValue(fmt.Sprintf("%d", siteIDInt))
		}
	} else {
		upgradedStateData.SiteID = types.StringNull()
	}

	if !priorStateData.AssignedComputerIDs.IsNull() && !priorStateData.AssignedComputerIDs.IsUnknown() {
		var computerIDsInt []int64
		diags := priorStateData.AssignedComputerIDs.ElementsAs(ctx, &computerIDsInt, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		computerIDsString := make([]string, len(computerIDsInt))
		for i, id := range computerIDsInt {
			computerIDsString[i] = fmt.Sprintf("%d", id)
		}

		setVal, diags := types.SetValueFrom(ctx, types.StringType, computerIDsString)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		upgradedStateData.AssignedComputerIDs = setVal
	} else {
		upgradedStateData.AssignedComputerIDs = types.SetNull(types.StringType)
	}

	timeoutsAttrTypes := map[string]attr.Type{
		"create": types.StringType,
		"read":   types.StringType,
		"update": types.StringType,
		"delete": types.StringType,
	}

	if priorStateData.Timeouts != nil {
		timeoutsAttrs := map[string]attr.Value{
			"create": priorStateData.Timeouts.Create,
			"read":   priorStateData.Timeouts.Read,
			"update": priorStateData.Timeouts.Update,
			"delete": priorStateData.Timeouts.Delete,
		}
		timeoutsObj, diags := types.ObjectValue(timeoutsAttrTypes, timeoutsAttrs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		upgradedStateData.Timeouts = timeouts.Value{Object: timeoutsObj}
	} else {
		upgradedStateData.Timeouts = timeouts.Value{Object: types.ObjectNull(timeoutsAttrTypes)}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
}
