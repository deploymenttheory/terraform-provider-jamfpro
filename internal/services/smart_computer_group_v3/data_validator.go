package smart_computer_group_v3

import (
	"context"
	"fmt"
	"regexp"

	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var siteIDPattern = regexp.MustCompile(`^(?:-1|[1-9][0-9]*)$`)

func (r *smartComputerGroupV3FrameworkResource) ValidateConfig(
	ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse,
) {
	var data smartComputerGroupV3ResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !data.SiteID.IsNull() && !data.SiteID.IsUnknown() && !siteIDPattern.MatchString(data.SiteID.ValueString()) {
		resp.Diagnostics.AddAttributeError(path.Root("site_id"), "Invalid Site ID", "Use -1 for None or a positive numeric site ID.")
	}
	if data.Criteria.IsNull() || data.Criteria.IsUnknown() {
		return
	}
	criteria, diags := schemahelpers.Expand[criterionModel](ctx, data.Criteria)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	balance := 0
	unknownParentheses := false
	for i, criterion := range criteria {
		if !criterion.Priority.IsUnknown() && criterion.Priority.ValueInt64() != int64(i) {
			resp.Diagnostics.AddAttributeError(path.Root("criteria").AtListIndex(i).AtName("priority"),
				"Invalid Criterion Priority", fmt.Sprintf("Priority must equal its zero-based position, %d.", i))
		}
		if criterion.OpeningParen.IsUnknown() || criterion.ClosingParen.IsUnknown() {
			unknownParentheses = true
			continue
		}
		if criterion.OpeningParen.ValueBool() {
			balance++
		}
		if criterion.ClosingParen.ValueBool() {
			balance--
		}
		if balance < 0 && !unknownParentheses {
			resp.Diagnostics.AddAttributeError(path.Root("criteria").AtListIndex(i),
				"Unbalanced Criteria Parentheses", "A closing parenthesis must have a preceding opening parenthesis.")
			return
		}
	}
	if balance != 0 && !unknownParentheses {
		resp.Diagnostics.AddAttributeError(path.Root("criteria"), "Unbalanced Criteria Parentheses",
			"Opening and closing parentheses must be balanced.")
	}
}
