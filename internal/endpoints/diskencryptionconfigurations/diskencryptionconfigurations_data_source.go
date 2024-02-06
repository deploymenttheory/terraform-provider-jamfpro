// diskencryptionconfigurations_resource.go
package diskencryptionconfigurations

import (
	"context"
	"strconv"
	"time"

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
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the disk encryption configuration.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the disk encryption configuration.",
			},
		},
	}
}

// DataSourceJamfProDiskEncryptionConfigurationsRead fetches the details of a specific Jamf Pro disk encryption configuration
// from Jamf Pro and returns the details of the disk encryption configuration in the Terraform state.
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
	var apiErrorCode int
	var script *jamfpro.ResourceDiskEncryptionConfiguration

	// Get the distribution point ID from the data source's arguments
	resourceID := d.Get("id").(string)

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
		script, apiErr = conn.GetDiskEncryptionConfigurationByID(resourceIDInt)
		if apiErr != nil {
			logging.LogFailedReadByID(subCtx, JamfProResourceDiskEncryptionConfiguration, resourceID, apiErr.Error(), apiErrorCode)
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the data, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		return diag.FromErr(err)
	}

	// Check if resource data exists and set the Terraform state
	if script != nil {
		d.SetId(resourceID) // Set the id in the Terraform state
		if err := d.Set("name", script.Name); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		d.SetId("") // Data not found, unset the id in the Terraform state
	}

	return diags
}
