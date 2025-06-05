package jamfconnect

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfConnectConfigProfile defines the schema and CRUD operations for managing Jamf Connect config profiles
func ResourceJamfConnectConfigProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(70 * time.Second),
			Update: schema.DefaultTimeout(70 * time.Second),
			Delete: schema.DefaultTimeout(70 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_profile_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier for the Jamf Connect config profile",
			},
			"profile_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the configuration profile",
			},
			"profile_name": {
				Type:        schema.TypeString,
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
				Optional:    true,
				Description: "Version of Jamf Connect to deploy.Versions are listed here https://www.jamf.com/resources/product-documentation/jamf-connect-administrators-guide/",
			},
			"auto_deployment_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Determines how the server will behave regarding application" +
					"updates and installs on the devices that have the configuration profile" +
					"installed. PATCH_UPDATES - Server handles initial installation of the" +
					"application and any patch updates. MINOR_AND_PATCH_UPDATES - Server handles" +
					"initial installation of the application and any patch and minor updates." +
					"INITIAL_INSTALLATION_ONLY - Server only handles initial installation of " +
					"the application. Updates will have to be done manually. NONE - Server does" +
					"not handle any installations or updates for the application. Version is " +
					"ignored for this type.",
				ValidateFunc: validation.StringInSlice([]string{
					"NONE",
					"PATCH_UPDATES",
					"MINOR_AND_PATCH_UPDATES",
					"INITIAL_INSTALLATION_ONLY",
				}, false),
			},
		},
	}
}
