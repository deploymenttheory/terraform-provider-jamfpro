// mobiledeviceconfigurationprofiles_data_handling.go
package mobiledeviceconfigurationprofiles

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// customDiffMobileDeviceConfigurationProfiles is the top-level custom diff function.
func customDiffMobileDeviceConfigurationProfiles(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if err := comparePayloadHashes(d); err != nil {
		return err
	}

	return nil
}

// comparePayloadHashes compares the payload hashes and ignores changes if they are the same.
func comparePayloadHashes(d *schema.ResourceDiff) error {

	return nil
}
