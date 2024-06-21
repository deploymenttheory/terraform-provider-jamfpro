// restrictedsoftware_resource.go
package restrictedsoftware

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProRestrictedSoftwares defines the schema and CRUD operations for managing Jamf Pro Restricted Software in Terraform.
func ResourceJamfProRestrictedSoftwares() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProRestrictedSoftwareCreate,
		ReadContext:   resourceJamfProRestrictedSoftwareRead,
		UpdateContext: resourceJamfProRestrictedSoftwareUpdate,
		DeleteContext: resourceJamfProRestrictedSoftwareDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the restricted software.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the restricted software.",
			},
			"process_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The process name of the restricted software.",
			},
			"match_exact_process_name": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the process name should be matched exactly.",
			},
			"send_notification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if a notification should be sent.",
			},
			"kill_process": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the process should be killed.",
			},
			"delete_executable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the executable should be deleted.",
			},
			"display_message": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The message to display when the software is restricted.",
			},
			"site": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The unique identifier of the site.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name of the site.",
						},
					},
				},
				Description: "The site associated with the restricted software.",
			},
			"scope": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The scope of the restricted software.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_computers": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if the restricted software applies to all computers.",
						},
						"computer_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "A list of computer IDs associated with the restricted software.",
						},
						"computer_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "A list of computer group IDs associated with the restricted software.",
						},
						"building_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "A list of building IDs associated with the restricted software.",
						},
						"department_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "A list of department IDs associated with the restricted software.",
						},
						"limitations": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Limitations for the restricted software.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_segment_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "A list of network segment IDs for limitations.",
									},
									"ibeacon_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "A list of iBeacon IDs for limitations.",
									},
								},
							},
						},
						"exclusions": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Exclusions for the restricted software.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"computer_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "A list of computer IDs for exclusions.",
									},
									"computer_group_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "A list of computer group IDs for exclusions.",
									},
									"building_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "A list of building IDs for exclusions.",
									},
									"department_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "A list of department IDs for exclusions.",
									},
									"directory_service_or_local_usernames": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "A list of directory service / local usernames for scoping exclusions.",
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

// scopeEntitySchema returns the schema for scope entities.
func scopeEntitySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The unique identifier of the scope entity.",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the scope entity.",
		},
	}
}
