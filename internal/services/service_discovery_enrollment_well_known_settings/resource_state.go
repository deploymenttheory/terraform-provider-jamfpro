package service_discovery_enrollment_well_known_settings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func applyResponse(data *serviceDiscoveryEnrollmentWellKnownSettingsResourceModel, response *jamfpro.ResponseServiceDiscoveryEnrollmentWellKnownSettingsV1) {
	if response == nil {
		return
	}

	data.ID = types.StringValue(serviceDiscoveryEnrollmentWellKnownSettingsSingletonID)
	data.WellKnownSettings = make([]serviceDiscoveryWellKnownSettingModel, 0, len(response.WellKnownSettings))

	for _, setting := range response.WellKnownSettings {
		model := serviceDiscoveryWellKnownSettingModel{
			OrgName:        stringValueOrNull(setting.OrgName),
			ServerUUID:     stringValueOrNull(setting.ServerUUID),
			EnrollmentType: stringValueOrNull(setting.EnrollmentType),
		}
		data.WellKnownSettings = append(data.WellKnownSettings, model)
	}
}

func stringValueOrNull(value string) types.String {
	if value == "" {
		return types.StringNull()
	}

	return types.StringValue(value)
}
