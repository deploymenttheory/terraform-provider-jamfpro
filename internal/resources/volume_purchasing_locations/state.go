// volumepurchasinglocations_state.go
package volumepurchasinglocations

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Volume Purchasing Location information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceVolumePurchasingLocation) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceData := map[string]interface{}{
		"id":                      resp.ID,
		"name":                    resp.Name,
		"apple_id":                resp.AppleID,
		"organization_name":       resp.OrganizationName,
		"token_expiration":        resp.TokenExpiration,
		"country_code":            resp.CountryCode,
		"location_name":           resp.LocationName,
		"client_context_mismatch": resp.ClientContextMismatch,
		"automatically_populate_purchased_content":  resp.AutomaticallyPopulatePurchasedContent,
		"send_notification_when_no_longer_assigned": resp.SendNotificationWhenNoLongerAssigned,
		"auto_register_managed_users":               resp.AutoRegisterManagedUsers,
		"site_id":                                   resp.SiteID,
		"last_sync_time":                            resp.LastSyncTime,
		"total_purchased_licenses":                  resp.TotalPurchasedLicenses,
		"total_used_licenses":                       resp.TotalUsedLicenses,
	}

	if serviceToken, ok := d.GetOk("service_token"); ok {
		resourceData["service_token"] = serviceToken
	}

	if resp.Content != nil {
		content := flattenContent(resp.Content)
		resourceData["content"] = content
	}

	for key, val := range resourceData {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

// flattenContent converts the Content slice to a format suitable for the Terraform state
func flattenContent(content []jamfpro.VolumePurchasingSubsetContent) []interface{} {
	var out []interface{}
	for _, c := range content {
		m := map[string]interface{}{
			"name":                   c.Name,
			"license_count_total":    c.LicenseCountTotal,
			"license_count_in_use":   c.LicenseCountInUse,
			"license_count_reported": c.LicenseCountReported,
			"icon_url":               c.IconURL,
			"device_types":           c.DeviceTypes,
			"content_type":           c.ContentType,
			"pricing_param":          c.PricingParam,
			"adam_id":                c.AdamId,
		}
		out = append(out, m)
	}
	return out
}
