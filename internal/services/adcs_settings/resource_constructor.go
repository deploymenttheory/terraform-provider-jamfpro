package adcs_settings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func constructResource(data *adcsSettingsResourceModel) (*jamfpro.ResourceAdcsSettingsV1, diag.Diagnostics) {
	var diags diag.Diagnostics

	resource := &jamfpro.ResourceAdcsSettingsV1{
		DisplayName: data.DisplayName.ValueString(),
		CAName:      data.CAName.ValueString(),
		FQDN:        data.FQDN.ValueString(),
		AdcsURL:     data.AdcsURL.ValueString(),
		APIClientID: data.APIClientID.ValueString(),
	}

	revocationEnabled := data.RevocationEnabled.ValueBool()
	resource.RevocationEnabled = &revocationEnabled

	outbound := data.Outbound.ValueBool()
	resource.Outbound = &outbound

	serverCert, serverDiags := buildCertificatePayload(
		data.ServerCertFilename,
		data.ServerCertData,
		data.ServerCertPassword,
		"server",
	)
	diags.Append(serverDiags...)
	if serverCert != nil {
		resource.ServerCert = serverCert
	}

	clientCert, clientDiags := buildCertificatePayload(
		data.ClientCertFilename,
		data.ClientCertData,
		data.ClientCertPassword,
		"client",
	)
	diags.Append(clientDiags...)
	if clientCert != nil {
		resource.ClientCert = clientCert
	}

	return resource, diags
}
