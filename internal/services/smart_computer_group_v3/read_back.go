package smart_computer_group_v3

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	frameworkcrud "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/framework_crud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ReadWithRetry is intentionally not used for post-write smart-group reads:
// Jamf can return HTTP 200 with stale criteria, which ReadWithRetry treats as success.
// Both v2 and v3 use WaitForWrite to compare the definition until the operation
// deadline, retaining the ID and last observed state without replaying writes.
// SDK reads also avoid Read removing state for a transient 404/410 after a write.
func (r *smartComputerGroupV3FrameworkResource) waitForWrite(
	ctx context.Context, id string, desired *jamfpro.ResourceSmartComputerGroupV3, description bool,
) (*jamfpro.ResourceSmartComputerGroupV3, error) {
	got, err := frameworkcrud.WaitForWrite(ctx,
		func(ctx context.Context) (*jamfpro.ResourceSmartComputerGroupV3, error) {
			return r.client.GetSmartComputerGroupByIDV3WithContext(ctx, id)
		},
		func(got *jamfpro.ResourceSmartComputerGroupV3) bool {
			return groupReflectsWrite(got, desired, description)
		},
		r.pollInterval,
	)
	if err != nil {
		return got, fmt.Errorf("verify smart computer group v3 write: %w", err)
	}
	return got, nil
}

func (r *smartComputerGroupV3FrameworkResource) completeWrite(
	ctx context.Context, id string, input *jamfpro.ResourceSmartComputerGroupV3, patchDescription bool,
) (*jamfpro.ResourceSmartComputerGroupV3, error) {
	if !patchDescription {
		return r.waitForWrite(ctx, id, input, true)
	}
	// A description PATCH can read-modify-write stale data on Jamf Pro 11.32.
	// Wait for the complete Classic write before making that separate request.
	got, err := r.waitForWrite(ctx, id, input, false)
	if err != nil {
		return got, err
	}
	platformID, err := r.client.GetComputerSmartGroupPlatformIDByJamfProIDWithContext(ctx, id)
	if err != nil {
		return got, fmt.Errorf("resolve Platform identifier for description update: %w", err)
	}
	if err := r.client.UpdateGroupDescriptionByIDV2WithContext(ctx, platformID, input.Description); err != nil {
		return got, fmt.Errorf("group saved but description update failed: %w", err)
	}
	return r.waitForWrite(ctx, id, input, true)
}

func (r *smartComputerGroupV3FrameworkResource) waitUntilGone(ctx context.Context, id string) {
	ctx, cancel := context.WithTimeout(ctx, r.deleteConfirmationTimeout)
	defer cancel()
	for ctx.Err() == nil {
		_, err := r.client.GetSmartComputerGroupByIDV3WithContext(ctx, id)
		if isNotFound(err) {
			return
		}
		if err != nil {
			tflog.Warn(ctx, "Delete succeeded but confirmation read failed", map[string]any{"id": id, "error": err.Error()})
			return
		}
		select {
		case <-ctx.Done():
		case <-time.After(r.pollInterval):
		}
	}
	tflog.Warn(ctx, "Delete succeeded but the group was still readable before the confirmation timeout", map[string]any{"id": id})
}
