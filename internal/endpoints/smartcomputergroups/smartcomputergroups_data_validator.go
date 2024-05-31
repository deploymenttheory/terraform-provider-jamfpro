// smartcomputergroup_data_validator.go
package smartcomputergroups

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// customDiffComputeGroups is the top-level custom diff function.
func customDiffComputeGroups(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {

	return nil
}

// getCriteriaOperators returns a list of criteria operators for Smart Computer Groups.
func getCriteriaOperators() []string {
	var out []string
	out = []string{
		And,
		Or,
		SearchTypeIs,
		SearchTypeIsNot,
		SearchTypeHas,
		SearchTypeDoesNotHave,
		SearchTypeMemberOf,
		SearchTypeNotMemberOf,
		SearchTypeBeforeYYYYMMDD,
		SearchTypeAfterYYYYMMDD,
		SearchTypeMoreThanXDaysAgo,
		SearchTypeLessThanXDaysAgo,
		SearchTypeLike,
		SearchTypeNotLike,
		SearchTypeGreaterThan,
		SearchTypeMoreThan,
		SearchTypeLessThan,
		SearchTypeGreaterThanOrEqual,
		SearchTypeLessThanOrEqual,
		SearchTypeMatchesRegex,
		SearchTypeDoesNotMatch,
	}
	return out
}
