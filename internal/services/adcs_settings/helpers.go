package adcs_settings

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// isNotFoundError checks if the provided error indicates a "not found" (404) response from the Jamf Pro API.
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "404")
}

// stringValue safely unwraps a Terraform string value, returning "" when null or unknown.
func stringValue(value types.String) string {
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return value.ValueString()
}

// buildCertificatePayload constructs a ResourceAdcsCertificateV1 from the provided filename, data, and password.
func buildCertificatePayload(filename, data, password types.String, certType string) (*jamfpro.ResourceAdcsCertificateV1, diag.Diagnostics) {
	var certDiags diag.Diagnostics

	filenameVal := stringValue(filename)
	dataVal := stringValue(data)
	passwordVal := stringValue(password)

	if filenameVal == "" && dataVal == "" && passwordVal == "" {
		return nil, certDiags
	}

	if filenameVal == "" {
		certDiags.AddError(
			fmt.Sprintf("Missing %s certificate filename", certType),
			fmt.Sprintf("The %s certificate filename must be provided when uploading certificate data.", certType),
		)
		return nil, certDiags
	}

	if dataVal == "" {
		certDiags.AddError(
			fmt.Sprintf("Missing %s certificate data", certType),
			fmt.Sprintf("The base64 encoded %s certificate data must be provided when uploading certificate credentials.", certType),
		)
		return nil, certDiags
	}

	decoded, err := base64.StdEncoding.DecodeString(dataVal)
	if err != nil {
		certDiags.AddError(
			fmt.Sprintf("Invalid %s certificate data", certType),
			fmt.Sprintf("The provided %s certificate data is not valid base64: %v", certType, err),
		)
		return nil, certDiags
	}

	certificate := &jamfpro.ResourceAdcsCertificateV1{
		Filename: filenameVal,
		Data:     decoded,
	}

	if passwordVal != "" {
		certificate.Password = passwordVal
	}

	return certificate, certDiags
}
