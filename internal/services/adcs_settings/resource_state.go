package adcs_settings

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func state(data *adcsSettingsResourceModel, resp *jamfpro.ResponseAdcsSettingsV1) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringValue(resp.ID)
	data.DisplayName = types.StringValue(resp.DisplayName)
	data.CAName = types.StringValue(resp.CAName)
	data.FQDN = types.StringValue(resp.FQDN)
	data.AdcsURL = types.StringValue(resp.AdcsURL)
	data.APIClientID = types.StringValue(resp.APIClientID)
	data.RevocationEnabled = types.BoolValue(resp.RevocationEnabled)
	data.Outbound = types.BoolValue(resp.Outbound)
	data.ConnectorLastCheckInTimestamp = stringOrNull(resp.ConnectorLastCheckInTimestamp)

	if resp.ServerCert != nil {
		data.ServerCertFilename = types.StringValue(resp.ServerCert.Filename)
		data.ServerCertData = types.StringNull()
		data.ServerCertPassword = types.StringNull()
		data.ServerCertSerialNumber = types.StringValue(resp.ServerCert.SerialNumber)
		data.ServerCertSubject = types.StringValue(resp.ServerCert.Subject)
		data.ServerCertIssuer = types.StringValue(resp.ServerCert.Issuer)
		data.ServerCertExpiration = types.StringValue(resp.ServerCert.ExpirationDate)
	} else {
		data.ServerCertFilename = types.StringNull()
		data.ServerCertData = types.StringNull()
		data.ServerCertPassword = types.StringNull()
		data.ServerCertSerialNumber = types.StringNull()
		data.ServerCertSubject = types.StringNull()
		data.ServerCertIssuer = types.StringNull()
		data.ServerCertExpiration = types.StringNull()
	}

	if resp.ClientCert != nil {
		data.ClientCertFilename = types.StringValue(resp.ClientCert.Filename)
		data.ClientCertData = types.StringNull()
		data.ClientCertPassword = types.StringNull()
		data.ClientCertSerialNumber = types.StringValue(resp.ClientCert.SerialNumber)
		data.ClientCertSubject = types.StringValue(resp.ClientCert.Subject)
		data.ClientCertIssuer = types.StringValue(resp.ClientCert.Issuer)
		data.ClientCertExpiration = types.StringValue(resp.ClientCert.ExpirationDate)
	} else {
		data.ClientCertFilename = types.StringNull()
		data.ClientCertData = types.StringNull()
		data.ClientCertPassword = types.StringNull()
		data.ClientCertSerialNumber = types.StringNull()
		data.ClientCertSubject = types.StringNull()
		data.ClientCertIssuer = types.StringNull()
		data.ClientCertExpiration = types.StringNull()
	}

	return diags
}

func stringOrNull(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}
