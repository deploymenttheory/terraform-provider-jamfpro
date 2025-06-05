// jamfconnect_data_source.go
package jamf_connect

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProJamfConnect provides information about Jamf Connect config profiles
func DataSourceJamfProJamfConnect() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"config_profile_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the Jamf Connect config profile",
			},
			"profile_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the configuration profile",
			},
			"profile_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the configuration profile",
			},
			"scope_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the profile's scope",
			},
			"site_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site ID of the site this profile belongs to",
			},
			"jamf_connect_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of Jamf Connect to deploy.Versions are listed here https://www.jamf.com/resources/product-documentation/jamf-connect-administrators-guide/",
			},
			"auto_deployment_type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Determines how the server will behave regarding application" +
					"updates and installs on the devices that have the configuration profile" +
					"installed. PATCH_UPDATES - Server handles initial installation of the" +
					"application and any patch updates. MINOR_AND_PATCH_UPDATES - Server handles" +
					"initial installation of the application and any patch and minor updates." +
					"INITIAL_INSTALLATION_ONLY - Server only handles initial installation of " +
					"the application. Updates will have to be done manually. NONE - Server does" +
					"not handle any installations or updates for the application. Version is " +
					"ignored for this type.",
			},
		},
	}
}

// dataSourceRead fetches Jamf Connect profile details
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	id := d.Get("profile_id").(int)
	name := d.Get("profile_name").(string)

	if id != 0 && name != "" {
		return diag.FromErr(fmt.Errorf("please provide either 'profile_id' or 'profile_name', not both"))
	}

	var getFunc func() (*jamfpro.ResourceJamfConnectConfigProfile, error)
	var identifier string

	switch {
	case id != 0:
		getFunc = func() (*jamfpro.ResourceJamfConnectConfigProfile, error) {
			return client.GetJamfConnectConfigProfileByID(id)
		}
		identifier = strconv.Itoa(id)
	case name != "":
		getFunc = func() (*jamfpro.ResourceJamfConnectConfigProfile, error) {
			return client.GetJamfConnectConfigProfileByName(name)
		}
		identifier = name
	default:
		return diag.FromErr(fmt.Errorf("either 'profile_id' or 'profile_name' must be provided"))
	}

	var resource *jamfpro.ResourceJamfConnectConfigProfile
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = getFunc()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Connect profile with identifier '%s' after retries: %v", identifier, err))
	}

	if resource == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("the Jamf Connect profile not found using identifier '%s'", identifier))
	}

	d.SetId(strconv.Itoa(resource.ProfileID))

	return updateState(d, resource)
}
