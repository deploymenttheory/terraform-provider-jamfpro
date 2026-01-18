package computer_extension_attribute

import (
	"context"
	"errors"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	errIDOrNameRequired                 = errors.New("either 'id' or 'name' must be provided")
	errReadComputerExtensionAttribute   = errors.New("failed to read Jamf Pro Computer Extension Attribute after retries")
	errComputerExtensionAttributeAbsent = errors.New("the Jamf Pro Computer Extension Attribute was not found")
)

// dataSourceRead fetches the details of a specific computer extension attribute
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	resourceID := d.Get("id").(string)
	name := d.Get("name").(string)

	if resourceID == "" && name == "" {
		return diag.FromErr(errIDOrNameRequired)
	}

	var resource *jamfpro.ResourceComputerExtensionAttribute
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error

		if name != "" {
			resource, apiErr = client.GetComputerExtensionAttributeByName(name)
		} else {
			resource, apiErr = client.GetComputerExtensionAttributeByID(resourceID)
		}

		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		lookupMethod := "ID"
		lookupValue := resourceID
		if name != "" {
			lookupMethod = "name"
			lookupValue = name
		}
		joinedErr := errors.Join(errReadComputerExtensionAttribute, err)
		return diag.FromErr(fmt.Errorf("%w with %s %q", joinedErr, lookupMethod, lookupValue))
	}

	if resource == nil {
		d.SetId("")
		return diag.FromErr(errComputerExtensionAttributeAbsent)
	}

	d.SetId(resource.ID)
	return updateState(d, resource)
}
