// sso_settings_resource.go
package ssosettings

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProSsoSettings defines the schema and CRUD operations for the SSO settings resource in Jamf Pro.
func ResourceJamfProSsoSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"sso_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable or disable SSO",
			},
			"configuration_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "SSO configuration type. Supported values are 'SAML', 'OIDC', 'OIDC_WITH_SAML'",
				ValidateFunc: validation.StringInSlice(getConfigurationTypes(), false),
			},
			"sso_bypass_allowed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow SSO bypass",
			},
			"sso_for_enrollment_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable SSO for enrollment",
			},
			"sso_for_macos_self_service_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable SSO for macOS Self Service",
			},
			"enrollment_sso_for_account_driven_enrollment_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable enrollment SSO for account-driven enrollment",
			},
			"group_enrollment_access_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable group enrollment access",
			},
			"group_enrollment_access_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group enrollment access name",
			},
			"oidc_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_mapping": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "OIDC user mapping type. Supported values are 'USERNAME', 'EMAIL'",
							ValidateFunc: validation.StringInSlice(getUserMappingTypes(), false),
						},
					},
				},
			},
			"saml_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"idp_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Identity Provider URL",
						},
						"entity_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Entity ID",
						},
						"metadata_source": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Metadata source type. Supported values are 'URL', 'FILE', 'UNKNOWN'",
							ValidateFunc: validation.StringInSlice(getMetadataSourceTypes(), false),
						},
						"user_mapping": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "SAML user mapping type. Supported values are 'USERNAME', 'EMAIL'",
							ValidateFunc: validation.StringInSlice(getUserMappingTypes(), false),
						},
						"idp_provider_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Identity Provider type. Supported values are 'ADFS', 'OKTA', 'GOOGLE', 'SHIBBOLETH', 'ONELOGIN', 'PING', 'CENTRIFY', 'AZURE', 'OTHER'",
							ValidateFunc: validation.StringInSlice(getIdpProviderTypes(), false),
						},
						"group_rdn_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Group RDN key",
						},
						"user_attribute_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User attribute name",
						},
						"group_attribute_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Group attribute name",
						},
						"user_attribute_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Enable user attributes",
						},
						"metadata_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Metadata file name. If metadata_source is set to URL, remove this field.",
						},
						"other_provider_type_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name for other provider type",
						},
						"federation_metadata_file": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Federation metadata file in base64 format. If metadata_source is set to URL, remove this field.",
						},
						"token_expiration_disabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Disable token expiration",
						},
						"session_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Session timeout in minutes",
						},
					},
				},
			},
			"enrollment_sso_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hosts": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "List of enrollment SSO hosts",
						},
						"management_hint": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Management hint for enrollment SSO",
						},
					},
				},
			},
		},
	}
}
