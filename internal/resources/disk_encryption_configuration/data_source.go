package disk_encryption_configuration

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProDiskEncryptionConfigurations defines the schema and CRUD operations for managing Jamf Pro Disk Encryption Configurations in Terraform.
func DataSourceJamfProDiskEncryptionConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the disk encryption configuration.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the disk encryption configuration.",
			},
			"key_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the key used in the disk encryption which can be either 'Institutional' or 'Individual and Institutional'.",
			},
			"file_vault_enabled_users": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines which user to enable for FileVault 2. Value can be either 'Management Account' or 'Current or Next User'",
			},
			"institutional_recovery_key": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details of the institutional recovery key.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_type": {
							Type:        schema.TypeString,
							Description: "The type of certificate used for the institutional recovery key. e.g 'PKCS12' for .p12 certificate types.",
							Computed:    true,
						},
						"password": {
							Type:        schema.TypeString,
							Description: "The password for the institutional recovery key certificate.",
							Computed:    true,
							Sensitive:   true,
						},
						"data": {
							Type:        schema.TypeString,
							Description: "The certificate payload.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// dataSourceRead fetches the details of a specific Jamf Pro disk encryption configuration
// from Jamf Pro and returns the details of the disk encryption configuration in the Terraform state.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	id := d.Get("id").(string)
	name := d.Get("name").(string)

	if id != "" && name != "" {
		return diag.FromErr(fmt.Errorf("please provide either 'id' or 'name', not both"))
	}

	var getFunc func() (*jamfpro.ResourceDiskEncryptionConfiguration, error)
	var identifier string

	switch {
	case id != "":
		getFunc = func() (*jamfpro.ResourceDiskEncryptionConfiguration, error) {
			return client.GetDiskEncryptionConfigurationByID(id)
		}
		identifier = id
	case name != "":
		getFunc = func() (*jamfpro.ResourceDiskEncryptionConfiguration, error) {
			return client.GetDiskEncryptionConfigurationByName(name)
		}
		identifier = name
	default:
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	var resource *jamfpro.ResourceDiskEncryptionConfiguration
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = getFunc()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Disk Encryption Configuration resource with identifier '%s' after retries: %v", identifier, err))
	}

	if resource == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("the Jamf Pro Disk Encryption Configuration resource was not found using identifier '%s'", identifier))
	}

	d.SetId(strconv.Itoa(resource.ID))
	return updateState(d, resource)
}
