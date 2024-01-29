// department_data_source.go
package departments

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/sdkv2"

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
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the centralized logger
	logger := sdkv2.ConsoleLogger{}

	var attribute *jamfpro.ResourceDepartment

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {

		attributeID := d.Id()

		var apiErr error

		attribute, apiErr = conn.GetDepartmentByID(attributeID)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				logger.Errorf("API Error (Code: %d): %s while fetching by ID", apiError.StatusCode, apiError.Message, "id", attributeID)
				return retry.NonRetryableError(apiErr)
			}

			attributeName, ok := d.Get("name").(string)
			if !ok {
				logger.Error("Unable to assert 'name' as a string for fetching department", "", "id", attributeID)
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			attribute, apiErr = conn.GetDepartmentByName(attributeName)
			if apiErr != nil {
				logger.Errorf("Error fetching department by name: %s", apiErr.Error(), "name", attributeName)
				return retry.RetryableError(apiErr)
			}
		}

		logger.Infof("Successfully fetched department: %s", attribute.Name)
		return nil
	})

	if err != nil {
		logger.Errorf("Failed to read department: %s", err.Error(), "id", d.Id())
		return logger.Diagnostics
	}

	if err := d.Set("name", attribute.Name); err != nil {
		logger.Errorf("Failed to set 'name' attribute in Terraform state: %s", err.Error(), "name", attribute.Name)
		return logger.Diagnostics
	}

	return nil
}
