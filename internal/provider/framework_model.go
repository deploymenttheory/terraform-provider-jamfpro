package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// frameworkProviderModel describes the provider data model.
type frameworkProviderModel struct {
	InstanceFQDN                      types.String `tfsdk:"jamfpro_instance_fqdn"`
	AuthMethod                        types.String `tfsdk:"auth_method"`
	ClientID                          types.String `tfsdk:"client_id"`
	ClientSecret                      types.String `tfsdk:"client_secret"`
	BasicAuthUsername                 types.String `tfsdk:"basic_auth_username"`
	BasicAuthPassword                 types.String `tfsdk:"basic_auth_password"`
	EnableClientSDKLogs               types.Bool   `tfsdk:"enable_client_sdk_logs"`
	ClientSDKLogExportPath            types.String `tfsdk:"client_sdk_log_export_path"`
	HideSensitiveData                 types.Bool   `tfsdk:"hide_sensitive_data"`
	LoadBalancerLock                  types.Bool   `tfsdk:"jamfpro_load_balancer_lock"`
	TokenRefreshBufferPeriodSeconds   types.Int64  `tfsdk:"token_refresh_buffer_period_seconds"`
	MandatoryRequestDelayMilliseconds types.Int64  `tfsdk:"mandatory_request_delay_milliseconds"`
	CustomCookies                     types.List   `tfsdk:"custom_cookies"`
}

// customCookieModel describes the custom cookie nested model.
type customCookieModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}
