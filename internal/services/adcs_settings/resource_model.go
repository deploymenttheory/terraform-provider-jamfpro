package adcs_settings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type adcsSettingsResourceModel struct {
	ID                            types.String   `tfsdk:"id"`
	DisplayName                   types.String   `tfsdk:"display_name"`
	CAName                        types.String   `tfsdk:"ca_name"`
	FQDN                          types.String   `tfsdk:"fqdn"`
	AdcsURL                       types.String   `tfsdk:"adcs_url"`
	APIClientID                   types.String   `tfsdk:"api_client_id"`
	RevocationEnabled             types.Bool     `tfsdk:"revocation_enabled"`
	Outbound                      types.Bool     `tfsdk:"outbound"`
	ConnectorLastCheckInTimestamp types.String   `tfsdk:"connector_last_check_in_timestamp"`
	ServerCertFilename            types.String   `tfsdk:"server_certificate_filename"`
	ServerCertData                types.String   `tfsdk:"server_certificate_data"`
	ServerCertPassword            types.String   `tfsdk:"server_certificate_password"`
	ServerCertSerialNumber        types.String   `tfsdk:"server_certificate_serial_number"`
	ServerCertSubject             types.String   `tfsdk:"server_certificate_subject"`
	ServerCertIssuer              types.String   `tfsdk:"server_certificate_issuer"`
	ServerCertExpiration          types.String   `tfsdk:"server_certificate_expiration_date"`
	ClientCertFilename            types.String   `tfsdk:"client_certificate_filename"`
	ClientCertData                types.String   `tfsdk:"client_certificate_data"`
	ClientCertPassword            types.String   `tfsdk:"client_certificate_password"`
	ClientCertSerialNumber        types.String   `tfsdk:"client_certificate_serial_number"`
	ClientCertSubject             types.String   `tfsdk:"client_certificate_subject"`
	ClientCertIssuer              types.String   `tfsdk:"client_certificate_issuer"`
	ClientCertExpiration          types.String   `tfsdk:"client_certificate_expiration_date"`
	Timeouts                      timeouts.Value `tfsdk:"timeouts"`
}
