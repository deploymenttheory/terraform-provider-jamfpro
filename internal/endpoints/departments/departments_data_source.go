// department_data_source.go
package departments

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the read operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	var apiErrorCode int

	// read operation

	department, err := conn.GetDepartmentByID(resourceID)
	if err != nil {
		if apiError, ok := err.(*http_client.APIError); ok {
			apiErrorCode = apiError.StatusCode
		}
		logging.LogFailedReadByID(subCtx, JamfProResourceDepartment, resourceID, err.Error(), apiErrorCode)
		return diags
	}

	logging.LogAPIReadSuccess(subCtx, JamfProResourceDepartment, resourceID)

	if err := d.Set("id", resourceID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", department.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
