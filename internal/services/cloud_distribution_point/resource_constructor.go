package cloud_distribution_point

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func constructCloudDistributionPointPayload(data *cloudDistributionPointResourceModel) (*jamfpro.ResourceCloudDistributionPointV1, diag.Diagnostics) {
	var diags diag.Diagnostics

	if data.CdnType.IsNull() || data.CdnType.IsUnknown() {
		diags.AddError(
			"Missing CDN Type",
			"Attribute cdn_type must be provided to manage the cloud distribution point.",
		)
		return nil, diags
	}

	resource := &jamfpro.ResourceCloudDistributionPointV1{
		CdnType: data.CdnType.ValueString(),
		Master:  data.Master.ValueBool(),
	}

	if !data.Username.IsNull() && !data.Username.IsUnknown() {
		resource.Username = data.Username.ValueString()
	}

	if !data.Password.IsNull() && !data.Password.IsUnknown() {
		resource.Password = data.Password.ValueString()
	}

	if !data.Directory.IsNull() && !data.Directory.IsUnknown() {
		resource.Directory = data.Directory.ValueString()
	}

	if !data.UploadURL.IsNull() && !data.UploadURL.IsUnknown() {
		resource.UploadUrl = data.UploadURL.ValueString()
	}

	if !data.DownloadURL.IsNull() && !data.DownloadURL.IsUnknown() {
		resource.DownloadUrl = data.DownloadURL.ValueString()
	}

	resource.SecondaryAuthRequired = boolPointerFromValue(data.SecondaryAuthRequired)
	resource.SecondaryAuthStatusCode = intPointerFromValue(data.SecondaryAuthStatusCode)
	resource.SecondaryAuthTimeToLive = intPointerFromValue(data.SecondaryAuthTimeToLive)
	resource.RequireSignedUrls = boolPointerFromValue(data.RequireSignedUrls)

	if !data.KeyPairID.IsNull() && !data.KeyPairID.IsUnknown() {
		resource.KeyPairId = data.KeyPairID.ValueString()
	}

	resource.ExpirationSeconds = intPointerFromValue(data.ExpirationSeconds)

	if !data.PrivateKey.IsNull() && !data.PrivateKey.IsUnknown() {
		resource.PrivateKey = data.PrivateKey.ValueString()
	}

	return resource, diags
}
