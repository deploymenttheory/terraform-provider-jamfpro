package macos_configuration_profile_plist

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceRead fetches the details of a macOS configuration profile.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	id := d.Get("id").(string)
	name := d.Get("name").(string)

	if id != "" && name != "" {
		return diag.FromErr(fmt.Errorf("please provide either 'id' or 'name', not both"))
	}

	var getFunc func() (*jamfpro.ResourceMacOSConfigurationProfile, error)
	var identifier string

	switch {
	case id != "":
		getFunc = func() (*jamfpro.ResourceMacOSConfigurationProfile, error) {
			return client.GetMacOSConfigurationProfileByID(id)
		}
		identifier = id
	case name != "":
		getFunc = func() (*jamfpro.ResourceMacOSConfigurationProfile, error) {
			return client.GetMacOSConfigurationProfileByName(name)
		}
		identifier = name
	default:
		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	var resource *jamfpro.ResourceMacOSConfigurationProfile
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = getFunc()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		//nolint:err113
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro macOS Configuration Profile resource with identifier '%s' after retries: %w", identifier, err))
	}

	if resource == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("the Jamf Pro Dock Item Configuration resource was not found using identifier '%s'", identifier))
	}

	d.SetId(strconv.Itoa(resource.General.ID))
	return updateState(d, resource)
}
