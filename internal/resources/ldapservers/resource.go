// ldapservers_resource.go
package ldapservers

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceJamfProLDAPServers defines the schema and CRUD operations for managing LDAP Servers in Terraform.
func ResourceJamfProLDAPServers() *schema.Resource {
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
			// Connection fields
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the LDAP server configuration.",
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The hostname or IP address of the LDAP server.",
			},
			"server_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The type of LDAP server. Valid values are 'Active Directory', 'Open Directory', 'eDirectory' or 'Custom'.",
				ValidateFunc: validation.StringInSlice([]string{"Active Directory", "Open Directory", "eDirectory", "Custom"}, false),
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The port number used to connect to the LDAP server (typically 389 for non-SSL or 636 for SSL).",
			},
			"use_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to use SSL for the LDAP connection.",
			},
			"authentication_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The authentication type used for LDAP binding (e.g., 'simple', 'CRAM-MD5', 'DIGEST-MD5').",
			},
			"account": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The account credentials used to bind to the LDAP server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"distinguished_username": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The distinguished username (DN) used to bind to the LDAP server.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "The password for the binding account.",
						},
					},
				},
			},
			"open_close_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout in seconds for opening and closing LDAP connections.",
			},
			"search_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The timeout in seconds for LDAP search operations.",
			},
			"referral_response": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "How to handle LDAP referrals (e.g., 'follow', 'ignore').",
			},
			"use_wildcards": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to use wildcards in LDAP searches.",
			},

			// User Mappings
			"user_mappings": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Configuration for mapping LDAP user attributes to Jamf Pro user attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"map_object_class_to_any_or_all": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether objects must match any or all object classes ('Any' or 'All').",
						},
						"object_classes": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The object classes that identify user records (e.g., 'user', 'person').",
						},
						"search_base": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The base DN where to search for users.",
						},
						"search_scope": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The scope of the search ('All Subtrees' or 'First Level Only').",
						},
						"map_user_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The LDAP attribute that maps to the user ID.",
						},
						"map_username": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The LDAP attribute that maps to the username.",
						},
						"map_realname": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The LDAP attribute that maps to the user's real name.",
						},
						"map_email_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The LDAP attribute that maps to the email address.",
						},
						"append_to_email_results": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to append to email addresses.",
						},
						"map_department": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The LDAP attribute that maps to the department.",
						},
						"map_building": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The LDAP attribute that maps to the building.",
						},
						"map_room": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The LDAP attribute that maps to the room.",
						},
						"map_telephone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The LDAP attribute that maps to the telephone number.",
						},
						"map_position": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The LDAP attribute that maps to the position/title.",
						},
						"map_user_uuid": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The LDAP attribute that maps to the user UUID.",
						},
					},
				},
			},

			// User Group Mappings
			"user_group_mappings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Configuration for mapping LDAP groups to Jamf Pro user groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"map_object_class_to_any_or_all": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether objects must match any or all object classes ('Any' or 'All').",
						},
						"object_classes": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The object classes that identify group records.",
						},
						"search_base": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The base DN where to search for groups.",
						},
						"search_scope": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The scope of the search ('All Subtrees' or 'First Level Only').",
						},
						"map_group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The LDAP attribute that maps to the group ID.",
						},
						"map_group_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The LDAP attribute that maps to the group name.",
						},
						"map_group_uuid": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The LDAP attribute that maps to the group UUID.",
						},
					},
				},
			},

			// User Group Membership Mappings
			"user_group_membership_mappings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Configuration for mapping user membership in LDAP groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_group_membership_stored_in": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Where group membership is stored ('User' or 'Group').",
						},
						"map_group_membership_to_user_field": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user attribute that contains group membership information.",
						},
						"append_to_username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text to append to usernames in group membership.",
						},
						"use_dn": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to use distinguished names for group membership.",
						},
						"recursive_lookups": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to recursively look up nested group memberships.",
						},
						"map_user_membership_to_group_field": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to map user membership to a group field.",
						},
						"map_user_membership_use_dn": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to use distinguished names when mapping user membership.",
						},
						"map_object_class_to_any_or_all": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether objects must match any or all object classes for membership.",
						},
						"object_classes": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The object classes used for membership mapping.",
						},
						"search_base": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The base DN for membership searches.",
						},
						"search_scope": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The scope for membership searches.",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The username attribute for membership mapping.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The group ID attribute for membership mapping.",
						},
						"user_group_membership_use_ldap_compare": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to use LDAP compare operations for membership checks.",
						},
						"membership_scoping_optimization": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable optimization for membership scoping.",
						},
					},
				},
			},
		},
	}
}
