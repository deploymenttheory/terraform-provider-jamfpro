// macosconfigurationprofilesplistgenerator_data_source.go
package macos_configuration_profile_plist_generator

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProMacOSConfigurationProfilesPlistGenerator provides information about a specific macOS configuration profile in Jamf Pro.
func DataSourceJamfProMacOSConfigurationProfilesPlistGenerator() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProMacOSConfigurationProfilePlistRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier for the macOS configuration profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the macOS configuration profile.",
			},
		},
	}
}

// DataSourceJamfProMacOSConfigurationProfilePlistRead fetches the details of a macOS configuration profile.
func DataSourceJamfProMacOSConfigurationProfilePlistRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)

	var resource *jamfpro.ResourceMacOSConfigurationProfile

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = client.GetMacOSConfigurationProfileByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro macOS Configuration Profile with ID '%s' after retries: %v", resourceID, err))
	}

	if resource != nil {
		d.SetId(resourceID)
		if err := d.Set("name", resource.General.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro macOS Configuration Profile with ID '%s': %v", resourceID, err))...)
		}
	} else {
		d.SetId("")
	}

	return diags
}
