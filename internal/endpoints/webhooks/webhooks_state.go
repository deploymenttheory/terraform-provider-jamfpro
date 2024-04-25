// webhooks_state.go
package webhooks

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateTerraformState updates the Terraform state with the latest Webhook information from the Jamf Pro API.
func updateTerraformState(d *schema.ResourceData, resource *jamfpro.ResourceWebhook) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update the Terraform state with the fetched data
	if err := d.Set("id", strconv.Itoa(resource.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resource.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("enabled", resource.Enabled); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("url", resource.URL); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("content_type", resource.ContentType); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("event", resource.Event); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("connection_timeout", resource.ConnectionTimeout); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("read_timeout", resource.ReadTimeout); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("authentication_type", resource.AuthenticationType); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("username", resource.Username); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("enable_display_fields_for_group", resource.EnableDisplayFieldsForGroup); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("smart_group_id", resource.SmartGroupID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle Display Fields
	displayFields := make([]interface{}, 0, len(resource.DisplayFields))
	for _, field := range resource.DisplayFields {
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
