package sharedschemas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func GetSharedMobileDeviceSchemaScope() *schema.Resource {
	scope := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"all_mobile_devices": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, the profile is applied to all mobile devices.",
			},
			"all_jss_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, the profile is applied to all JSS users.",
			},
			"mobile_device_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of mobile device IDs associated with the profile.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"mobile_device_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of mobile device group IDs associated with the profile.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"jss_user_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of JSS user IDs associated with the profile.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"jss_user_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of JSS user group IDs associated with the profile.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"building_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of building IDs associated with the profile.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"department_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of department IDs associated with the profile.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"limitations": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The scope limitations from the mobile device configuration profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_segment_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of network segment IDs for limitations.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"ibeacon_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of iBeacon IDs for limitations.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"directory_service_or_local_usernames": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of directory service / local usernames for scoping limitations.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"directory_service_usergroup_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of directory service user group IDs for limitations.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"exclusions": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The scope exclusions from the mobile device configuration profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mobile_device_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of mobile device IDs for exclusions.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"mobile_device_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of mobile device group IDs for exclusions.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"building_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of building IDs for exclusions.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"jss_user_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "A list of user names for exclusions.",
						},
						"jss_user_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of JSS user group IDs for exclusions.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"department_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of department IDs for exclusions.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"network_segment_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of network segment IDs for exclusions.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"directory_service_or_local_usernames": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of directory service / local usernames for scoping limitations.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"directory_service_or_local_usergroup_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of directory service / local user group IDs for limitations.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"ibeacon_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of iBeacon IDs for exclusions.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
		},
	}

	return scope
}
