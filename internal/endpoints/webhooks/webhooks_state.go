// webhooks_state.go
package webhooks

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Webhook information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resp *jamfpro.ResourceWebhook) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update the Terraform state with the fetched data
	if err := d.Set("id", strconv.Itoa(resp.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resp.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("enabled", resp.Enabled); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("url", resp.URL); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("content_type", resp.ContentType); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("event", resp.Event); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("connection_timeout", resp.ConnectionTimeout); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("read_timeout", resp.ReadTimeout); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("authentication_type", resp.AuthenticationType); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("username", resp.Username); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("enable_display_fields_for_group", resp.EnableDisplayFieldsForGroup); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("smart_group_id", resp.SmartGroupID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle Display Fields
	displayFields := make([]interface{}, 0, len(resp.DisplayFields))
	for _, field := range resp.DisplayFields {
		df := make(map[string]interface{})

		subFields := make([]interface{}, 0, len(field.DisplayField))
		for _, subField := range field.DisplayField {
			sf := map[string]interface{}{
				"name": subField.Name,
			}
			subFields = append(subFields, sf)
		}
		df["display_field"] = subFields
		displayFields = append(displayFields, df)
	}
	if err := d.Set("display_fields", displayFields); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
