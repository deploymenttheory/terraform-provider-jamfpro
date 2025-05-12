package cloudldap

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProCloudLdap defines the schema and CRUD operations for managing Jamf Pro Cloud LDAP in Terraform.
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
			"cloud_idp_common": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
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
					},
				},
			},
			"server": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the cloud LDAP server is enabled",
						},
						"keystore": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"file_bytes": {
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
										Description: "Base64 encoded keystore file",
									},
									"file_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Name of the keystore file",
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expiration_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"subject": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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
					},
				},
			},
			"mappings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_mappings": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_class_limitation": {
										Type:     schema.TypeString,
										Required: true,
									},
									"object_classes": {
										Type:     schema.TypeString,
										Required: true,
									},
									"search_base": {
										Type:     schema.TypeString,
										Required: true,
									},
									"search_scope": {
										Type:     schema.TypeString,
										Required: true,
									},
									"additional_search_base": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"user_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"username": {
										Type:     schema.TypeString,
										Required: true,
									},
									"real_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"email_address": {
										Type:     schema.TypeString,
										Required: true,
									},
									"department": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"building": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"room": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"phone": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"position": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"user_uuid": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"group_mappings": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_class_limitation": {
										Type:     schema.TypeString,
										Required: true,
									},
									"object_classes": {
										Type:     schema.TypeString,
										Required: true,
									},
									"search_base": {
										Type:     schema.TypeString,
										Required: true,
									},
									"search_scope": {
										Type:     schema.TypeString,
										Required: true,
									},
									"group_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"group_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"group_uuid": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"membership_mappings": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_membership_mapping": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
