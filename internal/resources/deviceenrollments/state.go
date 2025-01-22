// deviceenrollments_state.go
package deviceenrollments

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Device Enrollment information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceDeviceEnrollment) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]interface{}{
		"id":                      resp.ID,
		"name":                    resp.Name,
		"supervision_identity_id": resp.SupervisionIdentityId,
		"site_id":                 resp.SiteId,
		"server_name":             resp.ServerName,
		"server_uuid":             resp.ServerUuid,
		"admin_id":                resp.AdminId,
		"org_name":                resp.OrgName,
		"org_email":               resp.OrgEmail,
		"org_phone":               resp.OrgPhone,
		"org_address":             resp.OrgAddress,
		"token_expiration_date":   resp.TokenExpirationDate,
	}

	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
