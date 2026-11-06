package smart_computer_group_v3

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (r *smartComputerGroupV3FrameworkResource) waitForWrite(
	ctx context.Context, id string, desired *jamfpro.ResourceSmartComputerGroupV3, description bool,
) (*jamfpro.ResourceSmartComputerGroupV3, error) {
	ctx, cancel := context.WithTimeout(ctx, r.readBackBound)
	defer cancel()
	var last *jamfpro.ResourceSmartComputerGroupV3
	for {
		if err := ctx.Err(); err != nil {
			return last, fmt.Errorf("%w: %w", errReadBackNotConverged, err)
		}
		got, err := r.client.GetSmartComputerGroupByIDV3WithContext(ctx, id)
		if err != nil && !isNotFound(err) {
			return last, fmt.Errorf("read smart computer group after write: %w", err)
		}
		if err == nil && got != nil {
			last = got
			if groupReflectsWrite(got, desired, description) {
				return got, nil
			}
		}
		select {
		case <-ctx.Done():
			return last, fmt.Errorf("%w: %w", errReadBackNotConverged, ctx.Err())
		case <-time.After(r.pollInterval):
		}
	}
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
	ctx, cancel := context.WithTimeout(ctx, r.readBackBound)
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
