// accounts_data_source.go
package account

import (
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	errIDOrNameRequired = errors.New("either 'id' or 'name' must be provided")
)

// DataSourceJamfProAccounts provides information about specific Jamf Pro Accounts by their ID or Name.
func DataSourceJamfProAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the jamf pro account.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the jamf pro account.",
			},
		},
	}
}
