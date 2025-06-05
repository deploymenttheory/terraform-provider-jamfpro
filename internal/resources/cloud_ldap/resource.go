package cloud_ldap

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceJamfProCloudLdap() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"provider_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"GOOGLE", "AZURE"}, false),
				Description:  "The name of the cloud identity provider. Must be 'GOOGLE' or 'AZURE'.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name for the cloud LDAP configuration",
			},
			"server_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the cloud LDAP server is enabled",
			},
			"keystore_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"keystore_file_bytes": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Base64 encoded keystore file",
			},
			"keystore_file_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the keystore file",
			},
			"keystore_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"keystore_expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"keystore_subject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"use_wildcards": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to use wildcards in LDAP queries",
			},
			"connection_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"LDAPS", "START_TLS"}, false),
				Description:  "The type of LDAP connection (LDAPS or START_TLS)",
			},
			"server_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the LDAP server",
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The domain name for the LDAP server",
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "The port number for the LDAP server",
			},
			"connection_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Connection timeout in seconds",
			},
			"search_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Search timeout in seconds",
			},
			"membership_calculation_optimization_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable optimization for membership calculations",
			},
			"user_mappings_object_class_limitation": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ANY_OBJECT_CLASSES", "ALL_OBJECT_CLASSES"}, false),
				Description:  "Object class limitation for user mappings",
			},
			"user_mappings_object_classes": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Object classes for user mappings (e.g., inetOrgPerson)",
			},
			"user_mappings_search_base": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Search base for user mappings (e.g., ou=Users)",
			},
			"user_mappings_search_scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL_SUBTREES", "FIRST_LEVEL_ONLY"}, false),
				Description:  "Search scope for user mappings",
			},
			"user_mappings_additional_search_base": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional search base for user mappings",
			},
			"user_mappings_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User ID attribute mapping (e.g., mail)",
			},
			"user_mappings_username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username attribute mapping (e.g., uid)",
			},
			"user_mappings_real_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Real name attribute mapping (e.g., displayName)",
			},
			"user_mappings_email_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email address attribute mapping (e.g., mail)",
			},
			"user_mappings_department": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Department attribute mapping (e.g., departmentNumber)",
			},
			"user_mappings_building": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Building attribute mapping",
			},
			"user_mappings_room": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Room attribute mapping",
			},
			"user_mappings_phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Phone attribute mapping",
			},
			"user_mappings_position": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Position attribute mapping (e.g., title)",
			},
			"user_mappings_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User UUID attribute mapping (e.g., uid)",
			},
			"group_mappings_object_class_limitation": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ANY_OBJECT_CLASSES", "ALL_OBJECT_CLASSES"}, false),
				Description:  "Object class limitation for group mappings",
			},
			"group_mappings_object_classes": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Object classes for group mappings (e.g., groupOfNames)",
			},
			"group_mappings_search_base": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Search base for group mappings (e.g., ou=Groups)",
			},
			"group_mappings_search_scope": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL_SUBTREES", "FIRST_LEVEL_ONLY"}, false),
				Description:  "Search scope for group mappings",
			},
			"group_mappings_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group ID attribute mapping (e.g., cn)",
			},
			"group_mappings_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group name attribute mapping (e.g., cn)",
			},
			"group_mappings_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group UUID attribute mapping (e.g., gidNumber)",
			},
			"group_membership_mapping": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group membership attribute mapping (e.g., memberOf)",
			},
		},
	}
}
