package user

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest User information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceUser) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("id", strconv.Itoa(resp.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resp.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("full_name", resp.FullName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle email field - prefer Email over EmailAddress
	emailValue := resp.Email
	if emailValue == "" {
		emailValue = resp.EmailAddress
	}
	if err := d.Set("email", emailValue); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("email_address", emailValue); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("phone_number", resp.PhoneNumber); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("position", resp.Position); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Handle LDAP server fields
	if err := d.Set("ldap_server_id", resp.LDAPServer.ID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("ldap_server_name", resp.LDAPServer.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
