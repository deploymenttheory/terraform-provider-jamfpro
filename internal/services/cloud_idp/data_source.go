// cloudidp_data_source.go
package cloudidp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProCloudIdp provides information about a specific cloud identity provider in Jamf Pro.
func DataSourceJamfProCloudIdp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The jamf pro unique identifier of the cloud identity provider.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The display name of the cloud identity provider.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the cloud identity provider is enabled.",
			},
			"provider_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the cloud identity provider. e.g AZURE",
			},
		},
	}
}

// dataSourceRead fetches the details of a specific cloud identity provider from Jamf Pro using either its display name or its ID.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	resourceID := d.Get("id").(string)
	displayName := d.Get("display_name").(string)

	if resourceID == "" && displayName == "" {
		return diag.FromErr(fmt.Errorf("either 'id' or 'display_name' must be provided"))
	}

	var resource *jamfpro.ResourceCloudIdentityProvider
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error

		if displayName != "" {
			resource, apiErr = client.GetCloudIdentityProviderConfigurationByName(displayName)
		} else {
			details, detailsErr := client.GetCloudIdentityProviderConfigurationByID(resourceID)
			if detailsErr == nil && details != nil {
				resource = &jamfpro.ResourceCloudIdentityProvider{
					ID:           details.ID,
					DisplayName:  details.DisplayName,
					ProviderName: details.ProviderName,
					// Enabled is not available in the details response
				}
			} else {
				return retry.RetryableError(detailsErr)
			}
		}

		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		lookupMethod := "ID"
		lookupValue := resourceID
		if displayName != "" {
			lookupMethod = "display_name"
			lookupValue = displayName
		}
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Cloud Identity Provider with %s '%s' after retries: %v", lookupMethod, lookupValue, err))
	}

	if resource == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("the Jamf Pro Cloud Identity Provider was not found"))
	}

	d.SetId(resource.ID)

	return updateState(d, resource)
}
