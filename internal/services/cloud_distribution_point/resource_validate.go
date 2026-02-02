package cloud_distribution_point

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// cloudDistributionPointConfigValidator enforces cross-attribute rules at config time.
type cloudDistributionPointConfigValidator struct{}

func (cloudDistributionPointConfigValidator) Description(context.Context) string {
	return "Ensures CDN-specific attribute combinations are valid."
}

func (cloudDistributionPointConfigValidator) MarkdownDescription(context.Context) string {
	return "Ensures CDN-specific attribute combinations are valid."
}

func (cloudDistributionPointConfigValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data cloudDistributionPointResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateCloudDistributionPointPlan(&data)...)
}

// validateCloudDistributionPointPlan validates the provided cloud distribution point resource model.
func validateCloudDistributionPointPlan(data *cloudDistributionPointResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if data.CdnType.IsNull() || data.CdnType.IsUnknown() {
		diags.AddError(
			"Missing CDN Type",
			"Attribute cdn_type must be provided to manage the cloud distribution point.",
		)
		return diags
	}

	if data.Master.IsNull() || data.Master.IsUnknown() {
		diags.AddError(
			"Missing Master Flag",
			"Attribute master must be provided to manage the cloud distribution point.",
		)
		return diags
	}

	cdnType := data.CdnType.ValueString()

	usernameAllowedCdns := []string{cdnTypeRackspace, cdnTypeAmazonS3, cdnTypeAkamai}
	diags.Append(attributeAllowedOnlyForCdns("username", data.Username, cdnType, usernameAllowedCdns...)...)
	diags.Append(attributeAllowedOnlyForCdns("password", data.Password, cdnType, usernameAllowedCdns...)...)

	akamaiOnly := []string{cdnTypeAkamai}
	s3Only := []string{cdnTypeAmazonS3}

	diags.Append(attributeAllowedOnlyForCdns("directory", data.Directory, cdnType, akamaiOnly...)...)
	diags.Append(attributeAllowedOnlyForCdns("upload_url", data.UploadURL, cdnType, akamaiOnly...)...)
	diags.Append(attributeAllowedOnlyForCdns("download_url", data.DownloadURL, cdnType, akamaiOnly...)...)
	diags.Append(attributeAllowedOnlyForCdns("secondary_auth_required", data.SecondaryAuthRequired, cdnType, akamaiOnly...)...)
	diags.Append(attributeAllowedOnlyForCdns("secondary_auth_status_code", data.SecondaryAuthStatusCode, cdnType, akamaiOnly...)...)
	diags.Append(attributeAllowedOnlyForCdns("secondary_auth_time_to_live", data.SecondaryAuthTimeToLive, cdnType, akamaiOnly...)...)

	diags.Append(attributeAllowedOnlyForCdns("require_signed_urls", data.RequireSignedUrls, cdnType, s3Only...)...)
	diags.Append(attributeAllowedOnlyForCdns("key_pair_id", data.KeyPairID, cdnType, s3Only...)...)
	diags.Append(attributeAllowedOnlyForCdns("expiration_seconds", data.ExpirationSeconds, cdnType, s3Only...)...)
	diags.Append(attributeAllowedOnlyForCdns("private_key", data.PrivateKey, cdnType, s3Only...)...)

	switch cdnType {
	case cdnTypeRackspace, cdnTypeAmazonS3, cdnTypeAkamai:
		diags.Append(requireStringAttribute("username", data.Username, cdnType)...)
		diags.Append(requireStringAttribute("password", data.Password, cdnType)...)
	}

	if cdnType == cdnTypeAkamai {
		diags.Append(requireStringAttribute("directory", data.Directory, cdnType)...)
		diags.Append(requireStringAttribute("upload_url", data.UploadURL, cdnType)...)
		diags.Append(requireStringAttribute("download_url", data.DownloadURL, cdnType)...)
		diags.Append(requireBoolAttribute("secondary_auth_required", data.SecondaryAuthRequired, cdnType)...)

		if !data.SecondaryAuthRequired.IsNull() && !data.SecondaryAuthRequired.IsUnknown() && data.SecondaryAuthRequired.ValueBool() {
			diags.Append(requireIntAttribute("secondary_auth_status_code", data.SecondaryAuthStatusCode, cdnType)...)
			diags.Append(requireIntAttribute("secondary_auth_time_to_live", data.SecondaryAuthTimeToLive, cdnType)...)
		}
	}

	if cdnType == cdnTypeAmazonS3 {
		diags.Append(requireBoolAttribute("require_signed_urls", data.RequireSignedUrls, cdnType)...)
		if !data.RequireSignedUrls.IsNull() && !data.RequireSignedUrls.IsUnknown() && data.RequireSignedUrls.ValueBool() {
			diags.Append(requireStringAttribute("key_pair_id", data.KeyPairID, cdnType)...)
			diags.Append(requireIntAttribute("expiration_seconds", data.ExpirationSeconds, cdnType)...)
			diags.Append(requireStringAttribute("private_key", data.PrivateKey, cdnType)...)
		}
	}

	return diags
} // requireStringAttribute checks that a string attribute is set when required.
func requireStringAttribute(name string, value types.String, cdnType string) diag.Diagnostics {
	var diags diag.Diagnostics

	if value.IsNull() || value.IsUnknown() || value.ValueString() == "" {
		diags.AddError(
			fmt.Sprintf("Missing %s", name),
			fmt.Sprintf("Attribute %s must be set when cdn_type is %s.", name, cdnType),
		)
	}

	return diags
}

// requireBoolAttribute checks that a bool attribute is set when required.
func requireBoolAttribute(name string, value types.Bool, cdnType string) diag.Diagnostics {
	var diags diag.Diagnostics

	if value.IsNull() || value.IsUnknown() {
		diags.AddError(
			fmt.Sprintf("Missing %s", name),
			fmt.Sprintf("Attribute %s must be explicitly set when cdn_type is %s.", name, cdnType),
		)
	}

	return diags
}

// requireIntAttribute checks that an int attribute is set when required.
func requireIntAttribute(name string, value types.Int64, cdnType string) diag.Diagnostics {
	var diags diag.Diagnostics

	if value.IsNull() || value.IsUnknown() {
		diags.AddError(
			fmt.Sprintf("Missing %s", name),
			fmt.Sprintf("Attribute %s must be provided when cdn_type is %s.", name, cdnType),
		)
	}

	return diags
}

type nullableAttribute interface {
	IsNull() bool
	IsUnknown() bool
}

func attributeAllowedOnlyForCdns(name string, value nullableAttribute, cdnType string, allowed ...string) diag.Diagnostics {
	var diags diag.Diagnostics

	if value == nil || value.IsNull() || value.IsUnknown() {
		return diags
	}

	if cdnMatchesAllowed(cdnType, allowed) {
		return diags
	}

	var allowedList string
	if len(allowed) == 0 {
		allowedList = "not supported by any CDN type"
	} else {
		allowedList = fmt.Sprintf("only when cdn_type is %s", strings.Join(allowed, ", "))
	}

	diags.AddError(
		fmt.Sprintf("Attribute %s Is Not Valid For CDN Type", name),
		fmt.Sprintf("Attribute %s %s. Remove it or change cdn_type.", name, allowedList),
	)

	return diags
}

func cdnMatchesAllowed(cdnType string, allowed []string) bool {
	for _, candidate := range allowed {
		if cdnType == candidate {
			return true
		}
	}
	return false
}
