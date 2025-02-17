package smtpserver

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceJamfProSMTPServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: customDiff,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the SMTP server is enabled to allow Jamf Pro to send emails and invitations",
			},
			"authentication_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Type of SMTP authentication type to use",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "BASIC", "GRAPH_API", "GOOGLE_MAIL"}, false),
			},
			"connection_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "SMTP server hostname",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "SMTP server port",
						},
						"encryption_type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Type of encryption to use",
							ValidateFunc: validation.StringInSlice([]string{"NONE", "SSL", "TLS_1_2", "TLS_1_1", "TLS_1", "TLS_1_3"}, false),
						},
						"connection_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Connection timeout in seconds",
						},
					},
				},
			},
			"sender_settings": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "Jamf Pro Server",
							Description: "Display name for the sender",
						},
						"email_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Email address of the sender",
						},
					},
				},
			},
			"basic_auth_credentials": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Username for basic authentication",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Password for basic authentication",
						},
					},
				},
			},
			"graph_api_credentials": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Microsoft tenant ID for Graph API. Must be a valid GUID/UUID.",
							ValidateDiagFunc: validateGUID(),
						},
						"client_id": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Microsoft client ID for Graph API. Must be a valid GUID/UUID.",
							ValidateDiagFunc: validateGUID(),
						},
						"client_secret": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Microsoft client secret for Graph API",
						},
					},
				},
			},
			"google_mail_credentials": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Google client ID for Gmail",
							ValidateDiagFunc: validateGUID(),
						},
						"client_secret": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Google client secret for Gmail",
						},
					},
				},
			},
			"authentications": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Email address for authentication",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Authentication status",
						},
					},
				},
			},
		},
	}
}
