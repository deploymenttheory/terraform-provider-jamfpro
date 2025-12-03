package policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getPolicySchemaSelfService() *schema.Resource {
	selfServiceSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"use_for_self_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the policy is available for self-service.",
				Default:     false,
			},
			"self_service_display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Display name of the policy in self-service.",
				Default:     "",
			},
			"install_button_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Text displayed on the install button in self-service.",
				Default:     "Install",
			},
			"reinstall_button_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Text displayed on the re-install button in self-service.",
				Default:     "REINSTALL",
			},
			"self_service_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the policy displayed in self-service.",
				Default:     "",
			},
			"force_users_to_view_description": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to force users to view the policy description in self-service.",
				Default:     false,
			},
			"self_service_icon_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Icon for policy to use in self-service",
			},
			"feature_on_main_page": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to feature the policy on the main page of self-service.",
				Default:     false,
			},
			"self_service_category": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Category settings for the policy in self-service. Multiple categories can be specified.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Category ID for the policy in self-service.",
						},
						"display_in": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Whether to display the policy in this category in self-service.",
						},
						"feature_in": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to feature the policy in this category in self-service.",
						},
					},
				},
			},
			"notification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable notifications for this self-service policy.",
			},
			"notification_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Self Service",
				Description: "The type of notification. Valid values are 'Self Service' and 'Self Service and Notification Center'.",
				ValidateFunc: validation.StringInSlice([]string{
					"Self Service",
					"Self Service and Notification Center",
				}, false),
			},
			"notification_subject": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The subject of the notification message.",
			},
			"notification_message": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The body of the notification message.",
			},
		},
	}

	return selfServiceSchema
}
