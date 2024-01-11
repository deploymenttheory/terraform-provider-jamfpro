// department_data_source.go
package departments

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProDepartments provides information about a specific department in Jamf Pro.
func DataSourceJamfProDepartments() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceJamfProDepartmentsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the department.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the jamf pro department.",
			},
		},
	}
}

// DataSourceJamfProDepartmentsRead fetches the details of a specific department from Jamf Pro using either its unique Name or its Id.
func DataSourceJamfProDepartmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var attribute *jamfpro.ResourceDepartment

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Get the ID from the Terraform state into an integer to be used for the API request
		attributeID := d.Id()

		// Try fetching the site using the ID
		var apiErr error
		attribute, apiErr = conn.GetDepartmentByID(attributeID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			attributeName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			attribute, apiErr = conn.GetDepartmentByName(attributeName)
			if apiErr != nil {
				// Handle the APIError
				if apiError, ok := apiErr.(*http_client.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while reading the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	// Safely set attributes in the Terraform state
	if err := d.Set("name", attribute.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
