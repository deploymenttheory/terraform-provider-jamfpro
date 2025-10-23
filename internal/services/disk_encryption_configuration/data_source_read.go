package disk_encryption_configuration

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
