package webhook

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Webhook information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceWebhook) diag.Diagnostics {
	var diags diag.Diagnostics

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
	if err := d.Set("header", resp.Header); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("smart_group_id", resp.SmartGroupID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if len(resp.DisplayFields) > 0 {
		var displayFieldList []string
		for _, v := range resp.DisplayFields {
			displayFieldList = append(displayFieldList, v.Name)
		}

		d.Set("display_fields", displayFieldList)
	}

	return diags
}
