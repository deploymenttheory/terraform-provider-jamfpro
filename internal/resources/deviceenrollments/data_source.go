// deviceenrollments_data_source.go
package deviceenrollments

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProDeviceEnrollments provides information about device enrollments in Jamf Pro.
func DataSourceJamfProDeviceEnrollments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the device enrollment.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the device enrollment.",
			},
			"supervision_identity_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The supervision identity ID for the device enrollment.",
			},
			"site_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The site ID associated with the device enrollment.",
			},
			"server_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The server name for the device enrollment.",
			},
			"server_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The server UUID for the device enrollment.",
			},
			"admin_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The administrator ID associated with the device enrollment.",
			},
			"org_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization name for the device enrollment.",
			},
			"org_email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization email for the device enrollment.",
			},
			"org_phone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization phone number for the device enrollment.",
			},
			"org_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization address for the device enrollment.",
			},
			"token_expiration_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration date of the device enrollment token.",
			},
		},
	}
}

// dataSourceRead fetches the details of a device enrollment from Jamf Pro using either its Name or Id.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	id := d.Get("id").(string)
	name := d.Get("name").(string)

	if id != "" && name != "" {

		return diag.FromErr(fmt.Errorf("please provide either 'id' or 'name', not both"))
	}

	var getFunc func() (*jamfpro.ResourceDeviceEnrollment, error)
	var identifier string

	switch {
	case id != "":
		getFunc = func() (*jamfpro.ResourceDeviceEnrollment, error) {
			return client.GetDeviceEnrollmentByID(id)
		}
		identifier = id
	case name != "":
		getFunc = func() (*jamfpro.ResourceDeviceEnrollment, error) {
			return client.GetDeviceEnrollmentByName(name)
		}
		identifier = name
	default:

		return diag.FromErr(fmt.Errorf("either 'id' or 'name' must be provided"))
	}

	var resource *jamfpro.ResourceDeviceEnrollment
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = getFunc()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {

		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Device Enrollment resource with identifier '%s' after retries: %v", identifier, err))
	}

	if resource == nil {
		d.SetId("")

		return diag.FromErr(fmt.Errorf("the Jamf Pro Device Enrollment resource was not found using identifier '%s'", identifier))
	}

	d.SetId(resource.ID)
	return updateState(d, resource)
}
