package service_discovery_enrollment_well_known_settings

import "github.com/hashicorp/terraform-plugin-framework/types"

// serviceDiscoveryEnrollmentWellKnownSettingsDataSourceModel defines the schema for the service discovery enrollment well-known settings data source.
type serviceDiscoveryEnrollmentWellKnownSettingsDataSourceModel struct {
	ID                types.String                            `tfsdk:"id"`
	WellKnownSettings []serviceDiscoveryWellKnownSettingModel `tfsdk:"well_known_settings"`
}
