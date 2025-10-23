package ldap_server

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProLDAPServers provides information about specific Jamf Pro LDAP Servers by their ID or Name.
func DataSourceJamfProLDAPServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the jamf pro ldap server.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the jamf pro ldap server.",
			},
		},
	}
}
