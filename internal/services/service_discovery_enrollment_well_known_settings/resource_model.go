package service_discovery_enrollment_well_known_settings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	ResourceName  = "jamfpro_service_discovery_enrollment_well_known_settings_framework"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	serviceDiscoveryEnrollmentWellKnownSettingsSingletonID = "jamfpro_service_discovery_enrollment_well_known_settings_singleton"
)

type serviceDiscoveryEnrollmentWellKnownSettingsResourceModel struct {
	ID                types.String                            `tfsdk:"id"`
	WellKnownSettings []serviceDiscoveryWellKnownSettingModel `tfsdk:"well_known_settings"`
	Timeouts          timeouts.Value                          `tfsdk:"timeouts"`
}

type serviceDiscoveryWellKnownSettingModel struct {
	OrgName        types.String `tfsdk:"org_name"`
	ServerUUID     types.String `tfsdk:"server_uuid"`
	EnrollmentType types.String `tfsdk:"enrollment_type"`
}
