// sharedschemas/shared_schemas.go
package sharedschemas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// GetSharedmacOSComputerSchemaScope defines the reusable scope schema for macOS computer resources.
func GetSharedmacOSComputerSchemaScope() *schema.Resource {
	scope := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"all_computers": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the configuration profile is scoped to all computers.",
			},
			"all_jss_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the configuration profile is scoped to all JSS users.",
			},
			"computer_ids": {
				Type:        schema.TypeSet,
				Description: "The computers to which the configuration profile is scoped by Jamf ID.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"computer_group_ids": {
				Type:        schema.TypeSet,
				Description: "The computer groups to which the configuration profile is scoped by Jamf ID.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"jss_user_ids": {
				Type:        schema.TypeSet,
				Description: "The JSS users to which the configuration profile is scoped by Jamf ID.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"jss_user_group_ids": {
				Type:        schema.TypeSet,
				Description: "The JSS user groups to which the configuration profile is scoped by Jamf ID.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"building_ids": {
				Type:        schema.TypeSet,
				Description: "The buildings to which the configuration profile is scoped by Jamf ID.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"department_ids": {
				Type:        schema.TypeSet,
				Description: "The departments to which the configuration profile is scoped by Jamf ID.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"limitations": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The scope limitations from the macOS configuration profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_segment_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "A set of network segment IDs for limitations.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"directory_service_or_local_usernames": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "A set of directory service / local usernames for scoping limitations.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"directory_service_usergroup_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "A set of directory service user group IDs for limitations.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"ibeacon_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "A set of iBeacon IDs for limitations.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"exclusions": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The scope exclusions from the macOS configuration profile.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"computer_ids": {
							Type:        schema.TypeSet,
							Description: "Computers excluded from scope by Jamf ID.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"computer_group_ids": {
							Type:        schema.TypeSet,
							Description: "Computer Groups excluded from scope by Jamf ID.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"jss_user_ids": {
							Type:        schema.TypeSet,
							Description: "JSS Users excluded from scope by Jamf ID.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"jss_user_group_ids": {
							Type:        schema.TypeSet,
							Description: "JSS User Groups excluded from scope by Jamf ID.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"building_ids": {
							Type:        schema.TypeSet,
							Description: "Buildings excluded from scope by Jamf ID.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"department_ids": {
							Type:        schema.TypeSet,
							Description: "Departments excluded from scope by Jamf ID.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"network_segment_ids": {
							Type:        schema.TypeSet,
							Description: "Network segments excluded from scope by Jamf ID.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"directory_service_or_local_usernames": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "A set of directory service / local usernames for scoping exclusions.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"directory_service_usergroup_ids": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "A set of directory service / local user group IDs for exclusions.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"ibeacon_ids": {
							Type:        schema.TypeSet,
							Description: "Ibeacons excluded from scope by Jamf ID.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
		},
	}
	return scope
}
