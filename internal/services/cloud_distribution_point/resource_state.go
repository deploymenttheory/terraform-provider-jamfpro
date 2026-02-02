package cloud_distribution_point

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (data *cloudDistributionPointResourceModel) applyResponse(config *jamfpro.ResponseCloudDistributionPointV1, uploadCap *jamfpro.ResourceCloudDistributionPointUploadCapabilityV1) {
	if config == nil {
		return
	}

	data.ID = types.StringValue(cloudDistributionPointSingletonID)
	if config.InventoryId != "" {
		data.InventoryID = types.StringValue(config.InventoryId)
	} else {
		data.InventoryID = types.StringNull()
	}

	data.CdnType = types.StringValue(config.CdnType)
	data.Master = types.BoolValue(config.Master)
	data.Username = stringValueOrNull(config.Username)
	data.Directory = stringValueOrNull(config.Directory)
	data.CdnURL = stringValueOrNull(config.CdnUrl)
	data.UploadURL = stringValueOrNull(config.UploadUrl)
	data.DownloadURL = stringValueOrNull(config.DownloadUrl)
	data.SecondaryAuthRequired = types.BoolValue(config.SecondaryAuthRequired)
	data.SecondaryAuthStatusCode = int64ValueOrNull(config.SecondaryAuthStatusCode)
	data.SecondaryAuthTimeToLive = int64ValueOrNull(config.SecondaryAuthTimeToLive)
	data.RequireSignedUrls = types.BoolValue(config.RequireSignedUrls)
	data.KeyPairID = stringValueOrNull(config.KeyPairId)
	data.ExpirationSeconds = int64ValueOrNull(config.ExpirationSeconds)
	data.HasConnectionSucceeded = types.BoolValue(config.HasConnectionSucceeded)
	data.Message = stringValueOrNull(config.Message)

	if uploadCap != nil {
		data.PrincipalDistributionTechnology = types.BoolValue(uploadCap.ID)
		data.DirectUploadCapable = types.BoolValue(uploadCap.Name)
	} else {
		data.PrincipalDistributionTechnology = types.BoolNull()
		data.DirectUploadCapable = types.BoolNull()
	}
}

func stringValueOrNull(value string) types.String {
	if value == "" {
		return types.StringNull()
	}

	return types.StringValue(value)
}

func int64ValueOrNull(value int) types.Int64 {
	if value == 0 {
		return types.Int64Null()
	}

	return types.Int64Value(int64(value))
}
