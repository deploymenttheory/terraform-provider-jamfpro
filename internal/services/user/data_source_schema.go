package user

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProUsers provides information about a specific user in Jamf Pro.
func DataSourceJamfProUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the user.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The username/login name of the user.",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The email address of the user for lookup purposes.",
			},
			"full_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full display name of the user.",
			},
			"email_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of the user.",
			},
			"phone_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The phone number of the user.",
			},
			"position": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The job position/title of the user.",
			},
			"ldap_server_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the associated LDAP server.",
			},
			"ldap_server_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the associated LDAP server.",
			},
		},
	}
}
