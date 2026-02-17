package service_discovery_enrollment_well_known_settings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Read method for the service discovery enrollment well-known settings data source.
func (d *serviceDiscoveryEnrollmentWellKnownSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serviceDiscoveryEnrollmentWellKnownSettingsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.client.GetServiceDiscoveryEnrollmentWellKnownSettingsV1()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Service Discovery Enrollment Well-Known Settings",
			fmt.Sprintf("Could not read service discovery enrollment well-known settings: %s", err),
		)
		return
	}

	state.ID = types.StringValue(serviceDiscoveryEnrollmentWellKnownSettingsSingletonID)
	state.WellKnownSettings = make([]serviceDiscoveryWellKnownSettingModel, 0, len(response.WellKnownSettings))

	for _, setting := range response.WellKnownSettings {
		state.WellKnownSettings = append(state.WellKnownSettings, serviceDiscoveryWellKnownSettingModel{
			OrgName:        stringValueOrNull(setting.OrgName),
			ServerUUID:     stringValueOrNull(setting.ServerUUID),
			EnrollmentType: stringValueOrNull(setting.EnrollmentType),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
