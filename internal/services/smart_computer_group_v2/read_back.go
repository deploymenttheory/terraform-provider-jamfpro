package smart_computer_group_v2

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	frameworkcrud "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/framework_crud"
)

func (r *smartComputerGroupV2FrameworkResource) waitForWrite(ctx context.Context, id string, desired *jamfpro.ResourceSmartComputerGroupV2) (*jamfpro.ResourceSmartComputerGroupV2, error) {
	var last *jamfpro.ResourceSmartComputerGroupV2
	interval := r.pollInterval
	for {
		if err := ctx.Err(); err != nil {
			return last, fmt.Errorf("smart group read-back did not converge: %w", err)
		}
		got, err := r.client.GetSmartComputerGroupByIDV2WithContext(ctx, id)
		if err == nil && got != nil {
			last = got
			if groupReflectsWrite(got, desired) {
				return got, nil
			}
		} else if err != nil {
			code := frameworkcrud.APIStatusCode(err)
			if !frameworkcrud.IsNotFound(err) && code != http.StatusTooManyRequests && code < http.StatusInternalServerError {
				return last, fmt.Errorf("read smart computer group after write: %w", err)
			}
		}
		select {
		case <-ctx.Done():
			return last, fmt.Errorf("smart group read-back did not converge: %w", ctx.Err())
		case <-time.After(interval):
		}
		interval = min(interval*2, 5*time.Second)
	}
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
