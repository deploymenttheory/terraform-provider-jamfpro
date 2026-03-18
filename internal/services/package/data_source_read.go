package packages

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resourceID := d.Get("id").(string)
	name := d.Get("package_name").(string)

	if resourceID == "" && name == "" {
		return diag.FromErr(fmt.Errorf("either 'id' or 'package_name' must be provided"))
	}

	var resource *jamfpro.ResourcePackage
	err := retry.RetryContext(ctx, 10*time.Second, func() *retry.RetryError {
		var apiErr error

		if name != "" {
			resource, apiErr = client.GetPackageByName(name)
		} else {
			resource, apiErr = client.GetPackageByID(resourceID)
		}

		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	warnIfNotfound := d.Get("warn_if_not_found").(bool)

	if err != nil {
		lookupMethod := "ID"
		lookupValue := resourceID
		if name != "" {
			lookupMethod = "name"
			lookupValue = name
		}

		if warnIfNotfound {
			return append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Package at %v not found", lookupValue),
				Detail:   fmt.Sprintf("Not erroring due to warn_if_not_found enabled\nerr: %v", err),
			})
		}

		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Package with %s '%s' after retries: %v", lookupMethod, lookupValue, err))
	}

	if resource == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf(" the Jamf Pro Package cannot be found"))
	}

	d.SetId(resource.ID)
	return updateState(d, resource)
}
