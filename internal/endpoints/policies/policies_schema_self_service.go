package policies

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			// "self_service_icon": {
			// 	Type:        schema.TypeList,
			// 	Optional:    true,
			// 	Description: "Icon settings for the policy in self-service.",
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"id": {
			// 				Type:        schema.TypeInt,
			// 				Optional:    true,
			// 				Description: "ID of the icon used in self-service.",
			// 				Default:     0,
			// 			},
			// 			"filename": {
			// 				Type:        schema.TypeString,
			// 				Description: "Filename of the icon used in self-service.",
			// 				Computed:    true,
			// 			},
			// 			"uri": {
			// 				Type:        schema.TypeString,
			// 				Description: "URI of the icon used in self-service.",
			// 				Computed:    true,
			// 			},
			// 		},
			// 	},
			// },
			"feature_on_main_page": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to feature the policy on the main page of self-service.",
				Default:     false,
			},
			"self_service_categories": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Category settings for the policy in self-service.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Category details for the policy in self-service.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Category ID for the policy in self-service.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Category name for the policy in self-service.",
									},
									"display_in": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to display the category in self-service.",
										Default:     false,
									},
									"feature_in": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to feature the category in self-service.",
										Default:     false,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return selfServiceSchema
}
