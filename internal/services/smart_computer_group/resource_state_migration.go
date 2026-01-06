package smart_computer_group

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
var _ resource.ResourceWithUpgradeState = &smartComputerGroupFrameworkResource{}

// UpgradeState returns the state upgraders for the resource
func (r *smartComputerGroupFrameworkResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
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
					"site_id": schema.Int32Attribute{
						Optional: true,
						Computed: true,
					},
				},
				Blocks: map[string]schema.Block{
					"criteria": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Optional: true,
								},
								"priority": schema.Int32Attribute{
									Optional: true,
									Computed: true,
								},
								"and_or": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"search_type": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"value": schema.StringAttribute{
									Optional: true,
								},
								"opening_paren": schema.BoolAttribute{
									Optional: true,
									Computed: true,
								},
								"closing_paren": schema.BoolAttribute{
									Optional: true,
									Computed: true,
								},
							},
						},
					},
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
func upgradeStateV0toV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	type timeoutsV0 struct {
		Create types.String `tfsdk:"create"`
		Read   types.String `tfsdk:"read"`
		Update types.String `tfsdk:"update"`
		Delete types.String `tfsdk:"delete"`
	}

	type criteriaV0 struct {
		Name         types.String `tfsdk:"name"`
		Priority     types.Int32  `tfsdk:"priority"`
		AndOr        types.String `tfsdk:"and_or"`
		SearchType   types.String `tfsdk:"search_type"`
		Value        types.String `tfsdk:"value"`
		OpeningParen types.Bool   `tfsdk:"opening_paren"`
		ClosingParen types.Bool   `tfsdk:"closing_paren"`
	}

	type modelV0 struct {
		ID       types.String `tfsdk:"id"`
		Name     types.String `tfsdk:"name"`
		IsSmart  types.Bool   `tfsdk:"is_smart"`
		SiteID   types.Int32  `tfsdk:"site_id"`
		Criteria types.List   `tfsdk:"criteria"`
		Timeouts *timeoutsV0  `tfsdk:"timeouts"`
	}

	var priorStateData modelV0

	resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upgradedStateData := smartComputerGroupResourceModel{
		ID:          priorStateData.ID,
		Name:        priorStateData.Name,
		Description: types.StringNull(),
	}

	if !priorStateData.SiteID.IsNull() && !priorStateData.SiteID.IsUnknown() {
		siteIDInt := priorStateData.SiteID.ValueInt32()
		if siteIDInt == -1 {
			upgradedStateData.SiteID = types.StringValue("-1")
		} else {
			upgradedStateData.SiteID = types.StringValue(fmt.Sprintf("%d", siteIDInt))
		}
	} else {
		upgradedStateData.SiteID = types.StringNull()
	}

	if !priorStateData.Criteria.IsNull() && !priorStateData.Criteria.IsUnknown() {
		var criteriaV0List []criteriaV0
		diags := priorStateData.Criteria.ElementsAs(ctx, &criteriaV0List, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		upgradedStateData.Criteria = make([]smartComputerGroupCriteriaDataModel, 0, len(criteriaV0List))
		for _, oldCriteria := range criteriaV0List {
			newCriteria := smartComputerGroupCriteriaDataModel(oldCriteria)
			upgradedStateData.Criteria = append(upgradedStateData.Criteria, newCriteria)
		}
	} else {
		upgradedStateData.Criteria = []smartComputerGroupCriteriaDataModel{}
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
