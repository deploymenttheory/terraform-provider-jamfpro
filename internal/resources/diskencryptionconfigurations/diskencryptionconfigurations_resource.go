// diskencryptionconfigurations_resource.go
package diskencryptionconfigurations

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProDiskEncryptionConfigurations defines the schema and CRUD operations for managing Jamf Pro Disk Encryption Configurations in Terraform.
func ResourceJamfProDiskEncryptionConfigurations() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the disk encryption configuration.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the disk encryption configuration.",
			},
			"key_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the key used in the disk encryption which can be either 'Institutional' or 'Individual and Institutional'.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					value := val.(string)
					validValues := []string{"Individual", "Institutional", "Individual and Institutional"}

					found := false
					for _, v := range validValues {
						if value == v {
							found = true
							break
						}
					}

					if !found {
						errs = append(errs, fmt.Errorf("%q must be one of [%s], got '%s'", key, strings.Join(validValues, ", "), value))
					}

					return
				},
			},
			"file_vault_enabled_users": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Defines which user to enable for FileVault 2. Value can be either 'Management Account' or 'Current or Next User'",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					value := val.(string)
					validValues := []string{"Management Account", "Current or Next User"}

					found := false
					for _, v := range validValues {
						if value == v {
							found = true
							break
						}
					}

					if !found {
						errs = append(errs, fmt.Errorf("%q must be one of [%s], got '%s'", key, strings.Join(validValues, ", "), value))
					}

					return
				},
			},
			"institutional_recovery_key": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Details of the institutional recovery key.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"certificate_type": {
							Type:        schema.TypeString,
							Description: "The type of certificate used for the institutional recovery key. e.g 'PKCS12' for .p12 certificate types.",
							Optional:    true,
						},
						"password": {
							Type:        schema.TypeString,
							Description: "The password for the institutional recovery key certificate.",
							Optional:    true,
							Sensitive:   true,
						},
						"data": {
							Type:        schema.TypeString,
							Description: "The certificate payload.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

const (
	JamfProResourceDiskEncryptionConfiguration = "Disk Encryption Configuration"
)
