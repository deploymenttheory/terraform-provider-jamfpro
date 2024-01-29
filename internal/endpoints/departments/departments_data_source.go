// department_data_source.go
package departments

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
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

	// Initialize the logging subsystem for the read operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)

	// Initialize variables
	attributeID := d.Id()
	attributeName := d.Get("name").(string)

	// Get resource with timeout context
	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		attribute, apiErr := conn.GetDepartmentByID(attributeID)
		if apiErr != nil {
			attribute, apiErr = conn.GetDepartmentByName(attributeName)
			if apiErr != nil {
				logging.Error(subCtx, logging.SubsystemRead, "Error fetching department", map[string]interface{}{
					"id":    attributeID,
					"name":  attributeName,
					"error": apiErr.Error(),
				})

				return retry.RetryableError(apiErr)
			}
		}

		if attribute != nil {
			logging.Info(subCtx, logging.SubsystemRead, "Successfully fetched department", map[string]interface{}{
				"id":   attributeID,
				"name": attribute.Name,
			})

			// Set resource values into terraform state
			if err := d.Set("id", attribute.ID); err != nil {
				return retry.RetryableError(err)
			}
			if err := d.Set("name", attribute.Name); err != nil {
				return retry.RetryableError(err)
			}
		}

		return nil
	})

	if err != nil {
		logging.Error(subCtx, logging.SubsystemRead, "Failed to read department", map[string]interface{}{
			"id":    attributeID,
			"error": err.Error(),
		})

		return diag.FromErr(err)
	}

	return nil
}
