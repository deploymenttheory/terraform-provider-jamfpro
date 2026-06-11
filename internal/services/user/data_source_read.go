package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceRead fetches the details of a specific user from Jamf Pro using ID, name, or email.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	id := d.Get("id").(string)
	name := d.Get("name").(string)
	email := d.Get("email").(string)

	// Validate mutual exclusion
	providedCount := 0
	if id != "" {
		providedCount++
	}
	if name != "" {
		providedCount++
	}
	if email != "" {
		providedCount++
	}

	if providedCount == 0 {
		return diag.FromErr(fmt.Errorf("one of 'id', 'name', or 'email' must be provided"))
	}

	if providedCount > 1 {
		return diag.FromErr(fmt.Errorf("please provide only one of 'id', 'name', or 'email', not multiple"))
	}

	var getFunc func() (*jamfpro.ResourceUser, error)
	var identifier string
	var lookupMethod string

	switch {
	case name != "":
		getFunc = func() (*jamfpro.ResourceUser, error) {
			return client.GetUserByName(name)
		}
		identifier = name
		lookupMethod = "name"
	case email != "":
		getFunc = func() (*jamfpro.ResourceUser, error) {
			return client.GetUserByEmail(email)
		}
		identifier = email
		lookupMethod = "email"
	case id != "":
		getFunc = func() (*jamfpro.ResourceUser, error) {
			return client.GetUserByID(id)
		}
		identifier = id
		lookupMethod = "id"
	}

	var resource *jamfpro.ResourceUser
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = getFunc()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro User resource with %s '%s' after retries: %v", lookupMethod, identifier, err))
	}

	if resource == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("the Jamf Pro User resource was not found using %s '%s'", lookupMethod, identifier))
	}

	// Convert ID from int to string
	d.SetId(strconv.Itoa(resource.ID))
	return updateState(d, resource)
}
