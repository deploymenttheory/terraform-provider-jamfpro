package service_discovery_enrollment_well_known_settings

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func constructPayload(data *serviceDiscoveryEnrollmentWellKnownSettingsResourceModel) (*jamfpro.ResponseServiceDiscoveryEnrollmentWellKnownSettingsV1, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(data.WellKnownSettings) == 0 {
		diags.AddError(
			"Missing Well-Known Settings",
			"Attribute well_known_settings must include at least one entry.",
		)
		return nil, diags
	}

	payload := &jamfpro.ResponseServiceDiscoveryEnrollmentWellKnownSettingsV1{
		WellKnownSettings: make([]jamfpro.ResourceServiceDiscoveryWellKnownSettingV1, 0, len(data.WellKnownSettings)),
	}

	for index, setting := range data.WellKnownSettings {
		if setting.ServerUUID.IsNull() || setting.ServerUUID.IsUnknown() || setting.ServerUUID.ValueString() == "" {
			diags.AddError(
				"Missing Server UUID",
				fmt.Sprintf("well_known_settings[%d].server_uuid must be provided", index),
			)
			continue
		}
		if setting.EnrollmentType.IsNull() || setting.EnrollmentType.IsUnknown() || setting.EnrollmentType.ValueString() == "" {
			diags.AddError(
				"Missing Enrollment Type",
				fmt.Sprintf("well_known_settings[%d].enrollment_type must be provided", index),
			)
			continue
		}

		entry := jamfpro.ResourceServiceDiscoveryWellKnownSettingV1{
			ServerUUID:     setting.ServerUUID.ValueString(),
			EnrollmentType: setting.EnrollmentType.ValueString(),
		}
		if !setting.OrgName.IsNull() && !setting.OrgName.IsUnknown() && setting.OrgName.ValueString() != "" {
			entry.OrgName = setting.OrgName.ValueString()
		}

		payload.WellKnownSettings = append(payload.WellKnownSettings, entry)
	}

	if diags.HasError() {
		return nil, diags
	}

	return payload, diags
}
