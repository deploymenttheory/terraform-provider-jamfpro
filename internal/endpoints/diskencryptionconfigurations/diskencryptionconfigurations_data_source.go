// diskencryptionconfigurations_resource.go
package diskencryptionconfigurations

import (
	"context"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProDiskEncryptionConfigurations defines the schema and CRUD operations for managing Jamf Pro Disk Encryption Configurations in Terraform.
func DataSourceJamfProDiskEncryptionConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProDiskEncryptionConfigurationsRead,
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

func DataSourceJamfProDiskEncryptionConfigurationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the read operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	var apiErrorCode int
	var diskEncryptionConfig *jamfpro.ResourceDiskEncryptionConfiguration

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceDiskEncryptionConfiguration, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Read operation with retry
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		diskEncryptionConfig, apiErr = conn.GetDiskEncryptionConfigurationByID(resourceIDInt)
		if apiErr != nil {
			logging.LogFailedReadByID(subCtx, JamfProResourceDiskEncryptionConfiguration, resourceID, apiErr.Error(), apiErrorCode)
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the account group, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	// Assuming successful retrieval, proceed to set the resource attributes in Terraform state
	d.SetId(strconv.Itoa(resourceIDInt)) // Update the ID in the state
	d.Set("name", diskEncryptionConfig.Name)
	d.Set("key_type", diskEncryptionConfig.KeyType)
	d.Set("file_vault_enabled_users", diskEncryptionConfig.FileVaultEnabledUsers)

	// Institutional Recovery Key
	if diskEncryptionConfig.InstitutionalRecoveryKey == nil ||
		(diskEncryptionConfig.InstitutionalRecoveryKey.Key == "" &&
			diskEncryptionConfig.InstitutionalRecoveryKey.CertificateType == "" &&
			diskEncryptionConfig.InstitutionalRecoveryKey.Password == "" &&
			diskEncryptionConfig.InstitutionalRecoveryKey.Data == "") {

		// If InstitutionalRecoveryKey is nil or empty, ensure it is not set in the Terraform state
		d.Set("institutional_recovery_key", []interface{}{})
	} else {
		// If InstitutionalRecoveryKey has data, set it in the Terraform state
		irk := make(map[string]interface{})
		irk["certificate_type"] = diskEncryptionConfig.InstitutionalRecoveryKey.CertificateType
		irk["password"] = diskEncryptionConfig.InstitutionalRecoveryKey.Password
		irk["data"] = diskEncryptionConfig.InstitutionalRecoveryKey.Data

		d.Set("institutional_recovery_key", []interface{}{irk})
	}

	return diags
}
