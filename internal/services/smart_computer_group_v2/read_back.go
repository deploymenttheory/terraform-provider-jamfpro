package smart_computer_group_v2

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	frameworkcrud "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/framework_crud"
)

// ReadWithRetry is intentionally not used for post-write smart-group reads:
// Jamf can return HTTP 200 with stale criteria, which ReadWithRetry treats as success.
// Both v2 and v3 use WaitForWrite to compare the definition until the operation
// deadline, retaining the ID and last observed state without replaying writes.
// SDK reads also avoid Read removing state for a transient 404/410 after a write.
func (r *smartComputerGroupV2FrameworkResource) waitForWrite(ctx context.Context, id string, desired *jamfpro.ResourceSmartComputerGroupV2) (*jamfpro.ResourceSmartComputerGroupV2, error) {
	got, err := frameworkcrud.WaitForWrite(ctx,
		func(ctx context.Context) (*jamfpro.ResourceSmartComputerGroupV2, error) {
			return r.client.GetSmartComputerGroupByIDV2WithContext(ctx, id)
		},
		func(got *jamfpro.ResourceSmartComputerGroupV2) bool {
			return groupReflectsWrite(got, desired)
		},
		r.pollInterval,
	)
	if err != nil {
		return got, fmt.Errorf("verify smart computer group v2 write: %w", err)
	}
	return got, nil
}

func groupReflectsWrite(got, desired *jamfpro.ResourceSmartComputerGroupV2) bool {
	if got.Name != desired.Name || got.Description != desired.Description || siteID(got.SiteId) != siteID(desired.SiteId) {
		return false
	}
	a, b := orderedCriteria(got.Criteria), orderedCriteria(desired.Criteria)
	return slices.EqualFunc(a, b, func(a, b jamfpro.SharedSubsetCriteriaJamfProAPI) bool {
		return a.Name == b.Name && a.Priority == b.Priority && strings.EqualFold(a.AndOr, b.AndOr) &&
			a.SearchType == b.SearchType && a.Value == b.Value &&
			boolValue(a.OpeningParen) == boolValue(b.OpeningParen) && boolValue(a.ClosingParen) == boolValue(b.ClosingParen)
	})
}

func siteID(id *string) string {
	if id == nil || *id == "" {
		return "-1"
	}
	return *id
}

func boolValue(value *bool) bool {
	return value != nil && *value
}

func orderedCriteria(criteria []jamfpro.SharedSubsetCriteriaJamfProAPI) []jamfpro.SharedSubsetCriteriaJamfProAPI {
	result := slices.Clone(criteria)
	slices.SortFunc(result, func(a, b jamfpro.SharedSubsetCriteriaJamfProAPI) int {
		if a.Priority < b.Priority {
			return -1
		}
		if a.Priority > b.Priority {
			return 1
		}
		return 0
	})
	return result
}
