package policies

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO handle the commented attrs

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
				Description: "Category settings for the policy in self-service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Category ID for the policy in self-service.",
						},
						"display_in": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to display the category in self-service.",
						},
						"feature_in": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to feature the category in self-service.",
						},
					},
				},
			},
		},
	}

	return selfServiceSchema
}
