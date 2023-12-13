// policies_resource.go
package policies

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceJamfProPolicies defines the schema and CRUD operations for managing Jamf Pro Policy in Terraform.
func ResourceJamfProPolicies() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProPoliciesCreate,
		ReadContext:   ResourceJamfProPoliciesRead,
		UpdateContext: ResourceJamfProPoliciesUpdate,
		DeleteContext: ResourceJamfProPoliciesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"general": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "General settings of the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier of the Jamf Pro policy.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the policy.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the policy is enabled.",
							Default:     false,
						},
						"trigger": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Event(s) triggers to use to initiate the policy.",
							Default:     "EVENT",
						},
						"trigger_checkin": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Trigger policy when device performs recurring check-in against the frequency configured in Jamf Pro",
							Default:     false,
						},
						"trigger_enrollment_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Trigger policy when device enrollment is complete.",
							Default:     false,
						},
						"trigger_login": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Trigger policy when a user logs in to a computer. A login event that checks for policies must be configured in Jamf Pro for this to work",
							Default:     false,
						},
						"trigger_logout": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Trigger policy when a user logout.",
							Default:     false,
						},
						"trigger_network_state_changed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Trigger policy when it's network state changes. When a computer's network state changes (e.g., when the network connection changes, when the computer name changes, when the IP address changes)",
							Default:     false,
						},
						"trigger_startup": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Trigger policy when a computer starts up. A startup script that checks for policies must be configured in Jamf Pro for this to work",
							Default:     false,
						},
						"trigger_other": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Any other trigger for the policy.",
							Computed:    true,
						},
						"frequency": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Frequency of policy execution.",
							Default:     "Once per computer",
							ValidateFunc: validation.StringInSlice([]string{
								"Once per computer",
								"Once per user per computer",
								"Once per user",
								"Once every day",
								"Once every week",
								"Once every month",
								"Ongoing",
							}, false),
						},
						"retry_event": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Event on which to retry policy execution.",
							Default:     "none",
							ValidateFunc: validation.StringInSlice([]string{
								"none",
								"trigger",
								"check-in",
							}, false),
						},
						"retry_attempts": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of retry attempts for the jamf pro policy.",
							Default:     -1,
						},
						"notify_on_each_failed_retry": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Send notifications for each failed policy retry attempt. ",
							Default:     false,
						},
						"location_user_only": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Location-based policy for user only.",
							Default:     false,
						},
						"target_drive": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The drive on which to run the policy (e.g. /Volumes/Restore/ ). The policy runs on the boot drive by default",
							Default:     "/",
						},
						"offline": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the policy applies when offline.",
							Default:     false,
						},
						"category": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Category to add the policy to.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The category ID assigned to the jamf pro policy. Defaults to '-1' aka not used.",
										Computed:    true,
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Category Name for assigned jamf pro policy. Value defaults to 'No category assigned' aka not used",
										Computed:    true,
									},
								},
							},
						},
						"date_time_limitations": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Server-side limitations use your Jamf Pro host server's time zone and settings. The Jamf Pro host service is in UTC time.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"activation_date": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The activation date of the policy.",
										Computed:    true,
									},
									"activation_date_epoch": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The epoch time of the activation date.",
										Computed:    true,
									},
									"activation_date_utc": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The UTC time of the activation date.",
										Computed:    true,
									},
									"expiration_date": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The expiration date of the policy.",
										Computed:    true,
									},
									"expiration_date_epoch": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The epoch time of the expiration date.",
										Computed:    true,
									},
									"expiration_date_utc": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The UTC time of the expiration date.",
									},
									"no_execute_on": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}, false),
										},
										Description: "Client-side limitations are enforced based on the settings on computers. This field sets specific days when the policy should not execute.",
										Computed:    true,
									},
									"no_execute_start": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Client-side limitations are enforced based on the settings on computers. This field sets the start time when the policy should not execute.",
									},
									"no_execute_end": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Client-side limitations are enforced based on the settings on computers. This field sets the end time when the policy should not execute.",
									},
								},
							},
						},
						"network_limitations": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Network limitations for the policy.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"minimum_network_connection": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Minimum network connection required for the policy.",
										Default:      "No Minimum",
										ValidateFunc: validation.StringInSlice([]string{"No Minimum", "Ethernet"}, false),
									},
									"any_ip_address": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the policy applies to any IP address.",
										Default:     true,
									},
									"network_segments": {
										Type:        schema.TypeString,
										Description: "Network segment limitations for the policy.",
										Optional:    true,
									},
								},
							},
						},
						"override_default_settings": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Settings to override default configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_drive": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Target drive for the policy.",
										Default:     "/",
									},
									"distribution_point": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Distribution point for the policy.",
										Default:     "default",
									},
									"force_afp_smb": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to force AFP/SMB.",
										Default:     false,
									},
									"sus": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Software Update Service for the policy.",
										Default:     "default",
									},
								},
							},
						},
						"network_requirements": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Network requirements for the policy.",
							Default:     "Any",
						},
						"site": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Jamf Pro Site-related settings of the policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Jamf Pro Site ID. Value defaults to -1 aka not used.",
										Default:     -1,
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Jamf Pro Site Name. Value defaults to 'None' aka not used",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"scope": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Scope configuration for the profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_computers": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If true, applies the profile to all computers. If false applies to specific computers. Default is false.",
							Computed:    true,
						},
						"computers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"computer": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "The unique identifier of the scoped computer.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the scoped computer.",
												},
												"udid": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "UDID of the scoped computer.",
												},
											},
										},
									},
								},
							},
						},
						"buildings": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"building": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "The unique identifier of the scoped building.",
												},
												"name": {
													Type:        schema.TypeString,
													Description: "Name of the scoped building.",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"departments": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"department": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "The unique identifier of the scoped department.",
												},
												"name": {
													Type:        schema.TypeString,
													Description: "Name of the scoped department.",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"computer_groups": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"computer_group": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "The unique identifier of the scoped computer group.",
												},
												"name": {
													Type:        schema.TypeString,
													Description: "Name of the computer scoped group.",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"jss_users": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"jss_user": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "The unique identifier of the scoped JSS user.",
												},
												"name": {
													Type:        schema.TypeString,
													Description: "Name of the scoped JSS user.",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"jss_user_groups": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The unique identifier of the scoped JSS user group.",
									},
									"name": {
										Type:        schema.TypeString,
										Description: "Name of the scoped JSS user group.",
										Computed:    true,
									},
								},
							},
						},
						"limitations": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Scoped limitations for the policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_segments": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"network_segment": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the scoped network segment.",
															},
															"name": {
																Type:        schema.TypeString,
																Description: "Name of the scoped network segment.",
																Computed:    true,
															},
															"uid": {
																Type:        schema.TypeString,
																Description: "UID of the scoped network segment.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
									"users": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the user.",
															},
															"name": {
																Type:        schema.TypeString,
																Description: "Name of the user.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
									"user_groups": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user_group": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the scoped user group.",
															},
															"name": {
																Type:        schema.TypeString,
																Description: "Name of the scoped user group.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
									"ibeacons": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ibeacon": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the scoped iBeacon.",
															},
															"name": {
																Type:        schema.TypeString,
																Description: "Name of the scoped iBeacon.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"exclusions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Scoped exclusions to exclude from the policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"computers": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Scoped computer exclusions to exclude from the policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"computer": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "The individual computer to exclude from the policy.",
													MaxItems:    1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the computer.",
															},
															"name": {
																Type:        schema.TypeString,
																Description: "Name of the computer.",
																Computed:    true,
															},
															"udid": {
																Type:        schema.TypeString,
																Description: "UDID of the computer.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
									"computer_groups": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"computer_group": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the computer group.",
															},
															"name": {
																Type:        schema.TypeString,
																Description: "Name of the computer group.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
									"jss_users": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"jss_user": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the JSS user.",
															},
															"name": {
																Type:        schema.TypeString,
																Description: "Name of the JSS user.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
									"jss_user_groups": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"jss_user_group": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the JSS user group.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the JSS user group.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
									"buildings": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"building": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the building.",
															},
															"name": {
																Type:        schema.TypeString,
																Description: "Name of the building.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
									"departments": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"department": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the department.",
															},
															"name": {
																Type:        schema.TypeString,
																Description: "Name of the department.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
									"network_segments": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"network_segment": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the network segment.",
															},
															"uid": {
																Type:        schema.TypeString,
																Description: "UID of the network segment.",
																Computed:    true,
															},
															"name": {
																Type:        schema.TypeString,
																Description: "Name of the network segment.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
									"users": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the user.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the user.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
									"user_groups": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user_group": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the user group.",
															},
															"name": {
																Type:        schema.TypeString,
																Description: "Name of the user group.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
									"ibeacons": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ibeacon": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Required:    true,
																Description: "The unique identifier of the iBeacon.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the iBeacon.",
																Computed:    true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"self_service": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Self-service settings of the policy.",
				Elem: &schema.Resource{
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
							Default:     "Reinstall",
						},
						"self_service_description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the policy displayed in self-service.",
						},
						"force_users_to_view_description": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to force users to view the policy description in self-service.",
							Default:     false,
						},
						"self_service_icon": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Icon settings for the policy in self-service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "ID of the icon used in self-service.",
										Default:     0,
									},
									"filename": {
										Type:        schema.TypeString,
										Description: "Filename of the icon used in self-service.",
										Computed:    true,
									},
									"uri": {
										Type:        schema.TypeString,
										Description: "URI of the icon used in self-service.",
										Computed:    true,
									},
								},
							},
						},
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
													Required:    true,
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
													Default:     true,
												},
												"feature_in": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to feature the category in self-service.",
													Default:     true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"package_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Package configuration settings of the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"packages": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of packages included in the policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"package": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Details of the package.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Unique identifier of the package.",
												},
												"name": {
													Type:        schema.TypeString,
													Description: "Name of the package.",
													Computed:    true,
												},
												"action": {
													Type:         schema.TypeString,
													Optional:     true,
													Description:  "Action to be performed for the package.",
													ValidateFunc: validation.StringInSlice([]string{"Install", "Cache", "Install Cached"}, false),
													Default:      "Install",
												},
												"fut": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Fill User Template (FUT).",
												},
												"feu": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Fill Existing Users (FEU).",
												},
												"update_autorun": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Update auto-run status of the package.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"scripts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Scripts settings of the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The number of scripts included in the policy.",
						},
						"script": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Details of the scripts.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Unique identifier of the script.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name of the script.",
									},
									"priority": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Execution priority of the script.",
										ValidateFunc: validation.StringInSlice([]string{"Before", "After"}, false),
										Default:      "After",
									},
									"parameter4": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Custom parameter 4 for the script.",
									},
									"parameter5": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Custom parameter 5 for the script.",
									},
									"parameter6": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Custom parameter 6 for the script.",
									},
									"parameter7": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Custom parameter 7 for the script.",
									},
									"parameter8": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Custom parameter 8 for the script.",
									},
									"parameter9": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Custom parameter 9 for the script.",
									},
									"parameter10": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Custom parameter 10 for the script.",
									},
									"parameter11": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Custom parameter 11 for the script.",
									},
								},
							},
						},
					},
				},
			},
			"printers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Printers settings of the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"leave_existing_default": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Policy for handling existing default printers.",
						},
						"printer": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Details of the printer configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Unique identifier of the printer.",
									},
									"name": {
										Type:        schema.TypeString,
										Description: "Name of the printer.",
										Computed:    true,
									},
									"action": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Action to be performed for the printer (e.g., install, uninstall).",
										ValidateFunc: validation.StringInSlice([]string{"install", "uninstall"}, false),
									},
									"make_default": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to set the printer as the default.",
									},
								},
							},
						},
					},
				},
			},
			"dock_items": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Dock items settings of the policy.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dock_item": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Details of the dock item configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Unique identifier of the dock item.",
									},
									"name": {
										Type:        schema.TypeString,
										Description: "Name of the dock item.",
										Computed:    true,
									},
									"action": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Action to be performed for the dock item (e.g., Add To Beginning, Add To End, Remove).",
										ValidateFunc: validation.StringInSlice([]string{"Add To Beginning", "Add To End", "Remove"}, false),
									},
								},
							},
						},
					},
				},
			},
			"account_maintenance": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Account maintenance settings of the policy. Use this section to create and delete local accounts, and to reset local account passwords. Also use this section to disable an existing local account for FileVault 2.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accounts": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of account maintenance configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Details of the account configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:         schema.TypeString,
													Optional:     true,
													Description:  "Action to be performed on the account (e.g., Create, Reset, Delete, DisableFileVault).",
													ValidateFunc: validation.StringInSlice([]string{"Create", "Reset", "Delete", "DisableFileVault"}, false),
												},
												"username": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Username/short name for the account",
												},
												"realname": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Real name associated with the account.",
												},
												"password": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Set a new account password. This does not update the account's login keychain password or FileVault 2 password.",
												},
												"archive_home_directory": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Permanently delete home directory. If set to true will archive the home directory.",
												},
												"archive_home_directory_to": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Path in which to archive the home directory to.",
												},
												"home": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Full path in which to create the home directory (e.g. /Users/username/ or /private/var/username/)",
												},
												"hint": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Hint to help the user remember the password",
												},
												"picture": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Full path to the account picture (e.g. /Library/User Pictures/Animals/Butterfly.tif )",
												},
												"admin": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether the account has admin privileges.Setting this to true will set the user administrator privileges to the computer",
												},
												"filevault_enabled": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Allow the user to unlock the FileVault 2-encrypted drive",
												},
											},
										},
									},
								},
							},
						},
						"directory_bindings": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Directory binding settings for the policy. Use this section to bind computers to a directory service",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"binding": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Details of the directory binding.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "The unique identifier of the binding.",
												},
												"name": {
													Type:        schema.TypeString,
													Description: "The name of the binding.",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},
						"management_account": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Management account settings for the policy. Use this section to change or reset the management account password.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Action to perform on the management account.Rotates management account password at next policy execution. Management account passwords will be automatically randomized with 29 characters. Valid values are 'rotate' or 'doNotChange'.",
										ValidateFunc: validation.StringInSlice([]string{"rotate", "doNotChange"}, false),
										Default:      "doNotChange",
									},
									"managed_password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Managed password for the account.",
									},
									"managed_password_length": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Length of the managed password. Only necessary when utilizing the random action",
									},
								},
							},
						},
						"open_firmware_efi_password": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Open Firmware/EFI password settings for the policy. Use this section to set or remove an Open Firmware/EFI password on computers with Intel-based processors.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"of_mode": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Mode for the open firmware/EFI password. Valid values are 'command' or 'none'.",
										ValidateFunc: validation.StringInSlice([]string{"command", "none"}, false),
										Default:      "none",
									},
									"of_password": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Password for the open firmware/EFI.",
									},
								},
							},
						},
					},
				},
			},
			"reboot": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Use this section to restart computers and specify the disk to boot them to",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The reboot message displayed to the user.",
						},
						"specify_startup": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Reboot Method",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								validMethods := []string{"Standard Restart", "MDM Restart with Kernel Cache Rebuild"}
								for _, method := range validMethods {
									if v == method {
										return
									}
								}
								errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validMethods, v))
								return
							},
						},
						"startup_disk": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Disk to boot computers to",
							Default:     "Current Startup Disk",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								validDisks := []string{"Current Startup Disk", "Currently Selected Startup Disk (No Bless)", "macOS Installer", "Specify Local Startup Disk"}
								for _, disk := range validDisks {
									if v == disk {
										return
									}
								}
								errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validDisks, v))
								return
							},
						},
						"no_user_logged_in": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Action to take if no user is logged in to the computer",
							Default:     "Do not restart",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								validOptions := []string{"Restart if a package or update requires it", "Restart Immediately", "Do not restart"}
								for _, option := range validOptions {
									if v == option {
										return
									}
								}
								errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validOptions, v))
								return
							},
						},
						"user_logged_in": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "Do not restart",
							Description: "Action to take if a user is logged in to the computer",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								v := val.(string)
								validOptions := []string{"Restart if a package or update requires it", "Restart Immediately", "Restart", "Do not restart"}
								for _, option := range validOptions {
									if v == option {
										return
									}
								}
								errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validOptions, v))
								return
							},
						},
						"minutes_until_reboot": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Amount of time to wait before the restart begins.",
							Default:     5,
						},
						"start_reboot_timer_immediately": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Defines if the reboot timer should start immediately once the policy applies to a macOS device.",
							Default:     false,
						},
						"file_vault_2_reboot": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Perform authenticated restart on computers with FileVault 2 enabled. Restart FileVault 2-encrypted computers without requiring an unlock during the next startup",
							Default:     false,
						},
					},
				},
			},
			"maintenance": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Maintenance settings of the policy. Use this section to update inventory, reset computer names, install all cached packages, and run common maintenance tasks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recon": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to run recon (inventory update) as part of the maintenance. Forces computers to submit updated inventory information to Jamf Pro",
						},
						"reset_name": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to reset the computer name to the name stored in Jamf Pro. Changes the computer name on computers to match the computer name in Jamf Pro",
						},
						"install_all_cached_packages": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to install all cached packages. Installs packages cached by Jamf Pro",
						},
						"heal": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to heal the policy.",
						},
						"prebindings": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to update prebindings.",
						},
						"permissions": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to fix Disk Permissions (Not compatible with macOS v10.12 or later)",
						},
						"byhost": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to fix ByHost files andnpreferences.",
						},
						"system_cache": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to flush caches from /Library/Caches/ and /System/Library/Caches/, except for any com.apple.LaunchServices caches",
						},
						"user_cache": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to flush caches from ~/Library/Caches/, ~/.jpi_cache/, and ~/Library/Preferences/Microsoft/Office version #/Office Font Cache. Enabling this may cause problems with system fonts displaying unless a restart option is configured.",
						},
						"verify": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to verify system files and structure on the Startup Disk",
						},
					},
				},
			},
			"files_processes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Files and processes settings of the policy. Use this section to search for and log specific files and processes. Also use this section to execute a command.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search_by_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path of the file to search for.",
						},
						"delete_file": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to delete the file found at the specified path.",
							Computed:    true,
						},
						"locate_file": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Path of the file to locate. Name of the file, including the file extension. This field is case-sensitive and returns partial matches",
						},
						"update_locate_database": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to update the locate database. Update the locate database before searching for the file",
							Computed:    true,
						},
						"spotlight_search": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Search For File Using Spotlight. File to search for. This field is not case-sensitive and returns partial matches",
						},
						"search_for_process": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the process to search for. This field is case-sensitive and returns partial matches",
						},
						"kill_process": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to kill the process if found. This works with exact matches only",
						},
						"run_command": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Command to execute on computers. This command is executed as the 'root' user",
						},
					},
				},
			},
			"user_interaction": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "User interaction settings of the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message_start": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Message to display before the policy runs",
						},
						"allow_user_to_defer": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Allow user deferral and configure deferral type. A deferral limit must be specified for this to work.",
							Default:     false,
						},
						"allow_deferral_until_utc": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Date/time at which deferrals are prohibited and the policy runs. Uses time zone settings of your hosting server. Standard environments hosted in Jamf Cloud use Coordinated Universal Time (UTC)",
						},
						"allow_deferral_minutes": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of minutes after the user was first prompted by the policy at which the policy runs and deferrals are prohibited",
							Default:     "0",
						},
						"message_finish": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Message to display when the policy is complete.",
						},
					},
				},
			},
			"disk_encryption": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Disk encryption settings of the policy. Use this section to enable FileVault 2 or to issue a new recovery key.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The action to perform for disk encryption (e.g., apply, remediate).",
							ValidateFunc: validation.StringInSlice([]string{"none", "apply", "remediate"}, false),
							Default:      "none",
						},
						"disk_encryption_configuration_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "ID of the disk encryption configuration to apply.",
							Default:     "0",
						},
						"auth_restart": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to allow authentication restart.",
							Default:     false,
						},
						"remediate_key_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Type of key to use for remediation (e.g., Individual, Institutional, Individual And Institutional).",
							ValidateFunc: validation.StringInSlice([]string{"Individual", "Institutional", "Individual And Institutional"}, false),
						},
						"remediate_disk_encryption_configuration_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Disk encryption ID to utilize for remediating institutional recovery key types.",
						},
					},
				},
			},
		},
	}
}

// constructJamfProPolicy constructs a ResponsePolicy object from the provided schema data.
func constructJamfProPolicy(d *schema.ResourceData) (*jamfpro.ResponsePolicy, error) {
	// Initialize a new ResponsePolicy struct with all its sub-components.
	policy := &jamfpro.ResponsePolicy{
		General:              jamfpro.PolicyGeneral{},
		Scope:                jamfpro.PolicyScope{},
		SelfService:          jamfpro.PolicySelfService{},
		PackageConfiguration: jamfpro.PolicyPackageConfiguration{},
		Scripts:              jamfpro.PolicyScripts{},
		Printers:             jamfpro.PolicyPrinters{},
		DockItems:            jamfpro.PolicyDockItems{},
		AccountMaintenance:   jamfpro.PolicyAccountMaintenance{},
		Maintenance:          jamfpro.PolicyMaintenance{},
		FilesProcesses:       jamfpro.PolicyFilesProcesses{},
		UserInteraction:      jamfpro.PolicyUserInteraction{},
		DiskEncryption:       jamfpro.PolicyDiskEncryption{},
		Reboot:               jamfpro.PolicyReboot{},
	}

	// Construct the General section
	if v, ok := d.GetOk("general"); ok {
		generalData := v.([]interface{})[0].(map[string]interface{})
		policy.General = jamfpro.PolicyGeneral{
			Name:    getStringFromMap(generalData, "name"),
			Enabled: getBoolFromMap(generalData, "enabled"),
			Trigger: func() string {
				trigger := getStringFromMap(generalData, "trigger")
				if trigger != "" {
					return trigger
				}
				return "EVENT" // Default to "EVENT" if not specified
			}(),
			TriggerCheckin:             getBoolFromMap(generalData, "trigger_checkin"),
			TriggerEnrollmentComplete:  getBoolFromMap(generalData, "trigger_enrollment_complete"),
			TriggerLogin:               getBoolFromMap(generalData, "trigger_login"),
			TriggerLogout:              getBoolFromMap(generalData, "trigger_logout"),
			TriggerNetworkStateChanged: getBoolFromMap(generalData, "trigger_network_state_changed"),
			TriggerStartup:             getBoolFromMap(generalData, "trigger_startup"),
			TriggerOther: func() string {
				to := getStringFromMap(generalData, "trigger_other")
				if to != "" {
					return to
				}
				return "" // Default to an empty string if not specified
			}(),
			Frequency:               getStringFromMap(generalData, "frequency"),
			RetryEvent:              getStringFromMap(generalData, "retry_event"),
			RetryAttempts:           getIntFromMap(generalData, "retry_attempts"),
			NotifyOnEachFailedRetry: getBoolFromMap(generalData, "notify_on_each_failed_retry"),
			LocationUserOnly:        getBoolFromMap(generalData, "location_user_only"),
			TargetDrive: func() string {
				targetDriveValue := getStringFromMap(generalData, "target_drive")
				if targetDriveValue != "" {
					return targetDriveValue
				}
				return "/" // Set target drive default value
			}(),
			Offline: getBoolFromMap(generalData, "offline"),
			Category: func() jamfpro.PolicyCategory {
				// Initialize with default values
				defaultCategory := jamfpro.PolicyCategory{
					ID:   -1,                     // Default ID as integer
					Name: "No category assigned", // Default Name
				}

				// Check if values are provided in Terraform and override defaults if necessary
				if catData, ok := generalData["category"].([]interface{}); ok && len(catData) > 0 {
					catMap := catData[0].(map[string]interface{})

					catID := getIntFromMap(catMap, "id")
					catName := getStringFromMap(catMap, "name")

					// Only override the default if a valid value is provided
					if catID != -1 {
						defaultCategory.ID = catID
					}
					if catName != "" {
						defaultCategory.Name = catName
					}
				}

				return defaultCategory
			}(),
			// DateTimeLimitations field
			DateTimeLimitations: func() jamfpro.PolicyDateTimeLimitations {
				if dtData, ok := generalData["date_time_limitations"].([]interface{}); ok && len(dtData) > 0 {
					dateTimeMap := dtData[0].(map[string]interface{})
					dateTimeLimitations := jamfpro.PolicyDateTimeLimitations{
						ActivationDate:      getStringFromMap(dateTimeMap, "activation_date"),
						ActivationDateEpoch: getIntFromMap(dateTimeMap, "activation_date_epoch"),
						ActivationDateUTC:   getStringFromMap(dateTimeMap, "activation_date_utc"),
						ExpirationDate:      getStringFromMap(dateTimeMap, "expiration_date"),
						ExpirationDateEpoch: getIntFromMap(dateTimeMap, "expiration_date_epoch"),
						ExpirationDateUTC:   getStringFromMap(dateTimeMap, "expiration_date_utc"),
					}

					// Handling NoExecuteOn field
					if noExecOnData, ok := dateTimeMap["no_execute_on"].([]interface{}); ok {
						var noExecOnDays []jamfpro.PolicyNoExecuteOn
						for _, day := range noExecOnData {
							if dayMap, ok := day.(map[string]interface{}); ok {
								noExecOnDays = append(noExecOnDays, jamfpro.PolicyNoExecuteOn{
									Day: getStringFromMap(dayMap, "day"),
								})
							}
						}
						dateTimeLimitations.NoExecuteOn = noExecOnDays
					}

					// Handling NoExecuteStart and NoExecuteEnd fields
					dateTimeLimitations.NoExecuteStart = getStringFromMap(dateTimeMap, "no_execute_start")
					dateTimeLimitations.NoExecuteEnd = getStringFromMap(dateTimeMap, "no_execute_end")

					return dateTimeLimitations
				}
				return jamfpro.PolicyDateTimeLimitations{}
			}(),
			NetworkLimitations: func() jamfpro.PolicyNetworkLimitations {
				networkLimitations := jamfpro.PolicyNetworkLimitations{
					MinimumNetworkConnection: "No Minimum", // Default value
					AnyIPAddress:             true,         // Default value
				}

				if networkLimitationsData, ok := generalData["network_limitations"].([]interface{}); ok && len(networkLimitationsData) > 0 {
					netMap := networkLimitationsData[0].(map[string]interface{})

					if mnConn, ok := netMap["minimum_network_connection"]; ok && mnConn != "" {
						networkLimitations.MinimumNetworkConnection = getStringFromMap(netMap, "minimum_network_connection")
					}
					if _, ok := netMap["any_ip_address"]; ok {
						networkLimitations.AnyIPAddress = getBoolFromMap(netMap, "any_ip_address")
					}

					// Handling NetworkSegments field
					networkSegments := getStringFromMap(netMap, "network_segments")
					networkLimitations.NetworkSegments = networkSegments
				}

				return networkLimitations
			}(),
			// OverrideDefaultSettings field
			OverrideDefaultSettings: func() jamfpro.PolicyOverrideSettings {
				if overrideData, ok := generalData["override_default_settings"].([]interface{}); ok && len(overrideData) > 0 {
					overrideMap := overrideData[0].(map[string]interface{})
					overrideSettings := jamfpro.PolicyOverrideSettings{
						TargetDrive: getStringFromMap(overrideMap, "target_drive"),
						DistributionPoint: func() string {
							distributionPointValue := getStringFromMap(overrideMap, "distribution_point")
							if distributionPointValue != "" {
								return distributionPointValue
							}
							return "default" // Set distribution point default value
						}(),
						ForceAfpSmb: func() bool {
							forceAfpSmbValue, exists := overrideMap["force_afp_smb"].(bool)
							return exists && forceAfpSmbValue // Return the value if it exists, otherwise default (false)
						}(),
						SUS: func() string {
							susValue := getStringFromMap(overrideMap, "sus")
							if susValue != "" {
								return susValue
							}
							return "default" // Set SUS default value
						}(),
						NetbootServer: getStringFromMap(overrideMap, "netboot_server"),
					}

					return overrideSettings
				}
				return jamfpro.PolicyOverrideSettings{}
			}(),
			// NetworkRequirements field
			NetworkRequirements: func() string {
				networkRequirements := getStringFromMap(generalData, "network_requirements")
				if networkRequirements != "" {
					return networkRequirements
				}
				return "Any" // Default network requirements
			}(),
			// construct Site fields, and Initialize with default values
			Site: func() jamfpro.PolicySite {

				defaultSite := jamfpro.PolicySite{
					ID:   -1,
					Name: "None",
				}

				// Check if values are provided in Terraform and override defaults if necessary
				if siteData, ok := generalData["site"].([]interface{}); ok && len(siteData) > 0 {
					siteMap := siteData[0].(map[string]interface{})

					siteID := getIntFromMap(siteMap, "id")
					siteName := getStringFromMap(siteMap, "name")

					// Only override the default if a valid value is provided
					if siteID != -1 {
						defaultSite.ID = siteID
					}
					if siteName != "" {
						defaultSite.Name = siteName
					}
				}

				return defaultSite
			}(),
		}
	}

	// Construct the Scope section
	if v, ok := d.GetOk("scope"); ok {
		scopeData := v.([]interface{})[0].(map[string]interface{})
		var computers []jamfpro.PolicyDataSubsetComputer
		var computerGroups []jamfpro.PolicyDataSubsetComputerGroup
		var jssUsers []jamfpro.PolicyDataSubsetJSSUser
		var jssUserGroups []jamfpro.PolicyDataSubsetJSSUserGroup
		var buildings []jamfpro.PolicyDataSubsetBuilding
		var departments []jamfpro.PolicyDataSubsetDepartment

		// Construct Computers slice
		if comps, ok := scopeData["computers"].([]interface{}); ok {
			for _, comp := range comps {
				compMap := comp.(map[string]interface{})
				computers = append(computers, jamfpro.PolicyDataSubsetComputer{
					ID:   compMap["id"].(int),
					Name: compMap["name"].(string),
					UDID: compMap["udid"].(string),
				})
			}
		}

		// Construct ComputerGroups slice
		if groups, ok := scopeData["computer_groups"].([]interface{}); ok {
			for _, group := range groups {
				groupMap := group.(map[string]interface{})
				computerGroups = append(computerGroups, jamfpro.PolicyDataSubsetComputerGroup{
					ID:   groupMap["id"].(int),
					Name: groupMap["name"].(string),
				})
			}
		}

		// Construct JSSUsers slice
		if users, ok := scopeData["jss_users"].([]interface{}); ok {
			for _, user := range users {
				userMap := user.(map[string]interface{})
				jssUsers = append(jssUsers, jamfpro.PolicyDataSubsetJSSUser{
					ID:   userMap["id"].(int),
					Name: userMap["name"].(string),
				})
			}
		}

		// Construct JSSUserGroups slice
		if groups, ok := scopeData["jss_user_groups"].([]interface{}); ok {
			for _, group := range groups {
				groupMap := group.(map[string]interface{})
				jssUserGroups = append(jssUserGroups, jamfpro.PolicyDataSubsetJSSUserGroup{
					ID:   groupMap["id"].(int),
					Name: groupMap["name"].(string),
				})
			}
		}

		// Construct Buildings slice
		if bldgs, ok := scopeData["buildings"].([]interface{}); ok {
			for _, bldg := range bldgs {
				bldgMap := bldg.(map[string]interface{})
				buildings = append(buildings, jamfpro.PolicyDataSubsetBuilding{
					ID:   bldgMap["id"].(int),
					Name: bldgMap["name"].(string),
				})
			}
		}

		// Construct Departments slice
		if depts, ok := scopeData["departments"].([]interface{}); ok {
			for _, dept := range depts {
				deptMap := dept.(map[string]interface{})
				departments = append(departments, jamfpro.PolicyDataSubsetDepartment{
					ID:   deptMap["id"].(int),
					Name: deptMap["name"].(string),
				})
			}
		}

		// Construct LimitToUsers field
		var limitToUsers jamfpro.PolicyLimitToUsers
		if luData, ok := scopeData["limit_to_users"].([]interface{}); ok && len(luData) > 0 {
			luMap := luData[0].(map[string]interface{})
			var userGroups []string
			if uGroups, ok := luMap["user_groups"].([]interface{}); ok {
				for _, uGroup := range uGroups {
					userGroups = append(userGroups, uGroup.(string))
				}
			}
			limitToUsers = jamfpro.PolicyLimitToUsers{UserGroups: userGroups}
		}

		// Construct Limitations field
		var limitations jamfpro.PolicyLimitations
		if limitationsData, ok := scopeData["limitations"].([]interface{}); ok && len(limitationsData) > 0 {
			limitationsMap := limitationsData[0].(map[string]interface{})

			// Construct Directory Service/Local Users slice
			var users []jamfpro.PolicyDataSubsetUser
			if directoryServicesUsersData, ok := limitationsMap["users"].([]interface{}); ok {
				for _, user := range directoryServicesUsersData {
					userMap := user.(map[string]interface{})
					users = append(users, jamfpro.PolicyDataSubsetUser{
						ID:   userMap["id"].(int),
						Name: userMap["name"].(string),
					})
				}
			}

			// Construct Directory Service User Groups slice
			var userGroups []jamfpro.PolicyDataSubsetUserGroup
			if userGroupsData, ok := limitationsMap["user_groups"].([]interface{}); ok {
				for _, group := range userGroupsData {
					groupMap := group.(map[string]interface{})
					userGroups = append(userGroups, jamfpro.PolicyDataSubsetUserGroup{
						ID:   groupMap["id"].(int),
						Name: groupMap["name"].(string),
					})
				}
			}

			// Construct NetworkSegments slice
			var networkSegments []jamfpro.PolicyDataSubsetNetworkSegment
			if netSegsData, ok := limitationsMap["network_segments"].([]interface{}); ok {
				for _, seg := range netSegsData {
					segMap := seg.(map[string]interface{})
					networkSegments = append(networkSegments, jamfpro.PolicyDataSubsetNetworkSegment{
						ID:   segMap["id"].(int),
						Name: segMap["name"].(string),
						UID:  segMap["uid"].(string),
					})
				}
			}

			// Construct iBeacons slice
			var iBeacons []jamfpro.PolicyDataSubsetIBeacon
			if beaconsData, ok := limitationsMap["ibeacons"].([]interface{}); ok {
				for _, beacon := range beaconsData {
					beaconMap := beacon.(map[string]interface{})
					iBeacons = append(iBeacons, jamfpro.PolicyDataSubsetIBeacon{
						ID:   beaconMap["id"].(int),
						Name: beaconMap["name"].(string),
					})
				}
			}

			// Assign constructed slices to limitations struct
			limitations = jamfpro.PolicyLimitations{
				Users:           users,
				UserGroups:      userGroups,
				NetworkSegments: networkSegments,
				IBeacons:        iBeacons,
			}
		}

		// Assign Limitations to policy's Scope
		policy.Scope.Limitations = limitations

		// Construct Exclusions field
		var exclusions jamfpro.PolicyExclusions
		if exclusionsData, ok := scopeData["exclusions"].([]interface{}); ok && len(exclusionsData) > 0 {
			exclusionsMap := exclusionsData[0].(map[string]interface{})

			// Construct exclusion Computers slice
			var exclusionComputers []jamfpro.PolicyDataSubsetComputer
			if comps, ok := exclusionsMap["computers"].([]interface{}); ok {
				for _, comp := range comps {
					compMap := comp.(map[string]interface{})
					exclusionComputers = append(exclusionComputers, jamfpro.PolicyDataSubsetComputer{
						ID:   compMap["id"].(int),
						Name: compMap["name"].(string),
						UDID: compMap["udid"].(string),
					})
				}
			}

			// Construct exclusion ComputerGroups slice
			var exclusionComputerGroups []jamfpro.PolicyDataSubsetComputerGroup
			if groups, ok := exclusionsMap["computer_groups"].([]interface{}); ok {
				for _, group := range groups {
					groupMap := group.(map[string]interface{})
					exclusionComputerGroups = append(exclusionComputerGroups, jamfpro.PolicyDataSubsetComputerGroup{
						ID:   groupMap["id"].(int),
						Name: groupMap["name"].(string),
					})
				}
			}

			// Construct exclusion Users slice
			var exclusionUsers []jamfpro.PolicyDataSubsetUser
			if users, ok := exclusionsMap["users"].([]interface{}); ok {
				for _, user := range users {
					userMap := user.(map[string]interface{})
					exclusionUsers = append(exclusionUsers, jamfpro.PolicyDataSubsetUser{
						ID:   userMap["id"].(int),
						Name: userMap["name"].(string),
					})
				}
			}

			// Construct exclusion UserGroups slice
			var exclusionUserGroups []jamfpro.PolicyDataSubsetUserGroup
			if userGroups, ok := exclusionsMap["user_groups"].([]interface{}); ok {
				for _, group := range userGroups {
					groupMap := group.(map[string]interface{})
					exclusionUserGroups = append(exclusionUserGroups, jamfpro.PolicyDataSubsetUserGroup{
						ID:   groupMap["id"].(int),
						Name: groupMap["name"].(string),
					})
				}
			}

			// Construct exclusion Buildings slice
			var exclusionBuildings []jamfpro.PolicyDataSubsetBuilding
			if bldgs, ok := exclusionsMap["buildings"].([]interface{}); ok {
				for _, bldg := range bldgs {
					bldgMap := bldg.(map[string]interface{})
					exclusionBuildings = append(exclusionBuildings, jamfpro.PolicyDataSubsetBuilding{
						ID:   bldgMap["id"].(int),
						Name: bldgMap["name"].(string),
					})
				}
			}

			// Construct exclusion Departments slice
			var exclusionDepartments []jamfpro.PolicyDataSubsetDepartment
			if depts, ok := exclusionsMap["departments"].([]interface{}); ok {
				for _, dept := range depts {
					deptMap := dept.(map[string]interface{})
					exclusionDepartments = append(exclusionDepartments, jamfpro.PolicyDataSubsetDepartment{
						ID:   deptMap["id"].(int),
						Name: deptMap["name"].(string),
					})
				}
			}

			// Construct exclusion NetworkSegments slice
			var exclusionNetworkSegments []jamfpro.PolicyDataSubsetNetworkSegment
			if netSegments, ok := exclusionsMap["network_segments"].([]interface{}); ok {
				for _, segment := range netSegments {
					segmentMap := segment.(map[string]interface{})
					exclusionNetworkSegments = append(exclusionNetworkSegments, jamfpro.PolicyDataSubsetNetworkSegment{
						ID:   segmentMap["id"].(int),
						Name: segmentMap["name"].(string),
						UID:  segmentMap["uid"].(string),
					})
				}
			}

			// Construct exclusion JSSUsers slice
			var exclusionJSSUsers []jamfpro.PolicyDataSubsetJSSUser
			if jssUsers, ok := exclusionsMap["jss_users"].([]interface{}); ok {
				for _, jssUser := range jssUsers {
					jssUserMap := jssUser.(map[string]interface{})
					exclusionJSSUsers = append(exclusionJSSUsers, jamfpro.PolicyDataSubsetJSSUser{
						ID:   jssUserMap["id"].(int),
						Name: jssUserMap["name"].(string),
					})
				}
			}

			// Construct exclusion JSSUserGroups slice
			var exclusionJSSUserGroups []jamfpro.PolicyDataSubsetJSSUserGroup
			if jssUserGroups, ok := exclusionsMap["jss_user_groups"].([]interface{}); ok {
				for _, jssUserGroup := range jssUserGroups {
					jssUserGroupMap := jssUserGroup.(map[string]interface{})
					exclusionJSSUserGroups = append(exclusionJSSUserGroups, jamfpro.PolicyDataSubsetJSSUserGroup{
						ID:   jssUserGroupMap["id"].(int),
						Name: jssUserGroupMap["name"].(string),
					})
				}
			}

			// Construct exclusion iBeacons slice
			var exclusionIBeacons []jamfpro.PolicyDataSubsetIBeacon
			if beacons, ok := exclusionsMap["ibeacons"].([]interface{}); ok {
				for _, beacon := range beacons {
					beaconMap := beacon.(map[string]interface{})
					exclusionIBeacons = append(exclusionIBeacons, jamfpro.PolicyDataSubsetIBeacon{
						ID:   beaconMap["id"].(int),
						Name: beaconMap["name"].(string),
					})
				}
			}

			// Assign constructed slices to exclusions struct
			exclusions = jamfpro.PolicyExclusions{
				Computers:       exclusionComputers,
				ComputerGroups:  exclusionComputerGroups,
				Users:           exclusionUsers,
				UserGroups:      exclusionUserGroups,
				Buildings:       exclusionBuildings,
				Departments:     exclusionDepartments,
				NetworkSegments: exclusionNetworkSegments,
				JSSUsers:        exclusionJSSUsers,
				JSSUserGroups:   exclusionJSSUserGroups,
				IBeacons:        exclusionIBeacons,
			}

		}

		// Assign Exclusions to policy's Scope
		policy.Scope.Exclusions = exclusions

		// Assign constructed fields to the policy's Scope
		policy.Scope = jamfpro.PolicyScope{
			AllComputers: func() bool {
				if val, ok := scopeData["all_computers"].(bool); ok {
					return val
				}
				return false // Default value if not specified
			}(),
			Computers:      computers,
			ComputerGroups: computerGroups,
			JSSUsers:       jssUsers,
			JSSUserGroups:  jssUserGroups,
			Buildings:      buildings,
			Departments:    departments,
			LimitToUsers:   limitToUsers,
			Limitations:    limitations,
			Exclusions:     exclusions,
		}
	}

	// Construct the SelfService section
	if v, ok := d.GetOk("self_service"); ok {
		selfServiceData := v.([]interface{})[0].(map[string]interface{})
		policy.SelfService = jamfpro.PolicySelfService{
			UseForSelfService:           getBoolFromMap(selfServiceData, "use_for_self_service"),
			SelfServiceDisplayName:      getStringFromMap(selfServiceData, "self_service_display_name"),
			InstallButtonText:           getStringFromMap(selfServiceData, "install_button_text"),
			ReinstallButtonText:         getStringFromMap(selfServiceData, "reinstall_button_text"),
			SelfServiceDescription:      getStringFromMap(selfServiceData, "self_service_description"),
			ForceUsersToViewDescription: getBoolFromMap(selfServiceData, "force_users_to_view_description"),
			FeatureOnMainPage:           getBoolFromMap(selfServiceData, "feature_on_main_page"),
			SelfServiceIcon: func() jamfpro.PolicySelfServiceIcon {
				if iconData, ok := selfServiceData["self_service_icon"].([]interface{}); ok && len(iconData) > 0 {
					iconMap := iconData[0].(map[string]interface{})
					return jamfpro.PolicySelfServiceIcon{
						ID:       getIntFromMap(iconMap, "id"),
						Filename: getStringFromMap(iconMap, "filename"),
						URI:      getStringFromMap(iconMap, "uri"),
					}
				}
				return jamfpro.PolicySelfServiceIcon{}
			}(),
			SelfServiceCategories: func() []jamfpro.PolicySelfServiceCategory {
				var categories []jamfpro.PolicySelfServiceCategory
				if catData, ok := selfServiceData["self_service_categories"].([]interface{}); ok {
					for _, cat := range catData {
						catMap := cat.(map[string]interface{})
						if catDetails, ok := catMap["category"].([]interface{}); ok && len(catDetails) > 0 {
							details := catDetails[0].(map[string]interface{})
							category := jamfpro.PolicySelfServiceCategory{
								Category: jamfpro.PolicyCategory{
									ID:        getIntFromMap(details, "id"),
									Name:      getStringFromMap(details, "name"),
									DisplayIn: getBoolFromMap(details, "display_in"),
									FeatureIn: getBoolFromMap(details, "feature_in"),
								},
							}
							categories = append(categories, category)
						}
					}
				}
				return categories
			}(),
		}
	}
	// Construct the PackageConfiguration section
	if v, ok := d.GetOk("package_configuration"); ok {
		packageConfigData := v.([]interface{})[0].(map[string]interface{})

		packageItems := func() []jamfpro.PolicyPackage {
			var items []jamfpro.PolicyPackage
			if pkgs, ok := packageConfigData["packages"].([]interface{}); ok {
				for _, pkg := range pkgs {
					pkgMap := pkg.(map[string]interface{})

					items = append(items, jamfpro.PolicyPackage{
						ID:                getIntFromMap(pkgMap, "id"),
						Name:              getStringFromMap(pkgMap, "name"),
						Action:            getStringFromMap(pkgMap, "action"),
						FillUserTemplate:  getBoolFromMap(pkgMap, "fut"),
						FillExistingUsers: getBoolFromMap(pkgMap, "feu"),
						UpdateAutorun:     getBoolFromMap(pkgMap, "update_autorun"),
					})
				}
			}
			return items
		}()

		policy.PackageConfiguration = jamfpro.PolicyPackageConfiguration{
			Packages: packageItems,
		}
	}
	// Construct the Scripts section
	if v, ok := d.GetOk("scripts"); ok {
		scriptsData := v.([]interface{})[0].(map[string]interface{})

		scriptItems := func() []jamfpro.PolicyScriptItem {
			var items []jamfpro.PolicyScriptItem
			if scripts, ok := scriptsData["script"].([]interface{}); ok {
				for _, script := range scripts {
					scriptMap := script.(map[string]interface{})

					items = append(items, jamfpro.PolicyScriptItem{
						ID:          getStringFromMap(scriptMap, "id"),
						Name:        getStringFromMap(scriptMap, "name"),
						Priority:    getStringFromMap(scriptMap, "priority"),
						Parameter4:  getStringFromMap(scriptMap, "parameter4"),
						Parameter5:  getStringFromMap(scriptMap, "parameter5"),
						Parameter6:  getStringFromMap(scriptMap, "parameter6"),
						Parameter7:  getStringFromMap(scriptMap, "parameter7"),
						Parameter8:  getStringFromMap(scriptMap, "parameter8"),
						Parameter9:  getStringFromMap(scriptMap, "parameter9"),
						Parameter10: getStringFromMap(scriptMap, "parameter10"),
						Parameter11: getStringFromMap(scriptMap, "parameter11"),
					})
				}
			}
			return items
		}()

		size := getIntFromMap(scriptsData, "size")

		policy.Scripts = jamfpro.PolicyScripts{
			Size:   size,
			Script: scriptItems,
		}
	}

	// Construct the Printers section
	if v, ok := d.GetOk("printers"); ok {
		printersData := v.([]interface{})[0].(map[string]interface{})

		var printerItems []jamfpro.PolicyPrinterItem
		if printers, ok := printersData["printer"].([]interface{}); ok {
			for _, printer := range printers {
				printerMap := printer.(map[string]interface{})
				printerItems = append(printerItems, jamfpro.PolicyPrinterItem{
					ID:          getIntFromMap(printerMap, "id"),
					Name:        getStringFromMap(printerMap, "name"),
					Action:      getStringFromMap(printerMap, "action"),
					MakeDefault: getBoolFromMap(printerMap, "make_default"),
				})
			}
		}

		leaveExistingDefault := false
		if val, ok := printersData["leave_existing_default"].(bool); ok {
			leaveExistingDefault = val
		}

		policy.Printers = jamfpro.PolicyPrinters{
			LeaveExistingDefault: leaveExistingDefault,
			Printer:              printerItems,
		}
	}

	// Construct the DockItems section
	if v, ok := d.GetOk("dock_items"); ok {
		dockItemsData := v.([]interface{})[0].(map[string]interface{})

		dockItems := func() []jamfpro.PolicyDockItem {
			var items []jamfpro.PolicyDockItem
			if docks, ok := dockItemsData["dock_item"].([]interface{}); ok {
				for _, dock := range docks {
					dockMap := dock.(map[string]interface{})

					items = append(items, jamfpro.PolicyDockItem{
						ID:     getIntFromMap(dockMap, "id"),
						Name:   getStringFromMap(dockMap, "name"),
						Action: getStringFromMap(dockMap, "action"),
					})
				}
			}
			return items
		}()

		// Safely get 'size', default to 0 if not found or if type is different
		size := getIntFromMap(dockItemsData, "size")

		policy.DockItems = jamfpro.PolicyDockItems{
			Size:     size,
			DockItem: dockItems,
		}
	}

	// Construct the AccountMaintenance section
	if v, ok := d.GetOk("account_maintenance"); ok {
		accountMaintenanceData := v.([]interface{})[0].(map[string]interface{})

		accounts := func() []jamfpro.PolicyAccount {
			var items []jamfpro.PolicyAccount
			if accs, ok := accountMaintenanceData["accounts"].([]interface{}); ok {
				for _, acc := range accs {
					accMap := acc.(map[string]interface{})
					items = append(items, jamfpro.PolicyAccount{
						Action:                 getStringFromMap(accMap, "action"),
						Username:               getStringFromMap(accMap, "username"),
						Realname:               getStringFromMap(accMap, "realname"),
						Password:               getStringFromMap(accMap, "password"),
						ArchiveHomeDirectory:   getBoolFromMap(accMap, "archive_home_directory"),
						ArchiveHomeDirectoryTo: getStringFromMap(accMap, "archive_home_directory_to"),
						Home:                   getStringFromMap(accMap, "home"),
						Hint:                   getStringFromMap(accMap, "hint"),
						Picture:                getStringFromMap(accMap, "picture"),
						Admin:                  getBoolFromMap(accMap, "admin"),
						FilevaultEnabled:       getBoolFromMap(accMap, "filevault_enabled"),
					})
				}
			}
			return items
		}()

		directoryBindings := func() []jamfpro.PolicyDirectoryBinding {
			var items []jamfpro.PolicyDirectoryBinding
			if bindingsList, ok := accountMaintenanceData["directory_bindings"].([]interface{}); ok {
				for _, bindingEntry := range bindingsList {
					bindingData := bindingEntry.(map[string]interface{})
					if bindings, ok := bindingData["binding"].([]interface{}); ok {
						for _, binding := range bindings {
							bindingMap := binding.(map[string]interface{})
							items = append(items, jamfpro.PolicyDirectoryBinding{
								ID:   getIntFromMap(bindingMap, "id"),
								Name: getStringFromMap(bindingMap, "name"),
							})
						}
					}
				}
			}
			return items
		}()

		managementAccount := func() jamfpro.PolicyManagementAccount {
			// Initialize with default values
			defaultManagementAccount := jamfpro.PolicyManagementAccount{
				Action:                "doNotChange",
				ManagedPassword:       "",
				ManagedPasswordLength: 0,
			}

			// Check if values are provided in Terraform and override defaults if necessary
			if maData, ok := accountMaintenanceData["management_account"].(map[string]interface{}); ok {
				defaultManagementAccount.Action = getStringFromMap(maData, "action")
				defaultManagementAccount.ManagedPassword = getStringFromMap(maData, "managed_password")
				defaultManagementAccount.ManagedPasswordLength = getIntFromMap(maData, "managed_password_length")
			}

			return defaultManagementAccount
		}()

		openFirmwareEfiPassword := func() jamfpro.PolicyOpenFirmwareEfiPassword {
			// Initialize with default values
			defaultOpenFirmwareEfiPassword := jamfpro.PolicyOpenFirmwareEfiPassword{
				OfMode:           "none",
				OfPassword:       "",
				OfPasswordSHA256: "",
			}

			// Check if values are provided in Terraform and override defaults if necessary
			if ofData, ok := accountMaintenanceData["open_firmware_efi_password"].(map[string]interface{}); ok {
				defaultOpenFirmwareEfiPassword.OfMode = getStringFromMap(ofData, "of_mode")
				defaultOpenFirmwareEfiPassword.OfPassword = getStringFromMap(ofData, "of_password")
				defaultOpenFirmwareEfiPassword.OfPasswordSHA256 = getStringFromMap(ofData, "of_password_sha256")
			}

			return defaultOpenFirmwareEfiPassword
		}()

		// Assign all constructed components to AccountMaintenance
		policy.AccountMaintenance = jamfpro.PolicyAccountMaintenance{
			Accounts:                accounts,
			DirectoryBindings:       directoryBindings,
			ManagementAccount:       managementAccount,
			OpenFirmwareEfiPassword: openFirmwareEfiPassword,
		}
	}

	// Construct the Reboot section and set default values if none are defined in terraform
	policy.Reboot = func() jamfpro.PolicyReboot {
		// Set default values
		defaultReboot := jamfpro.PolicyReboot{
			Message:                     "This computer will restart in 5 minutes. Please save anything you are working on and log out by choosing Log Out from the bottom of the Apple menu.",
			StartupDisk:                 "Current Startup Disk",
			SpecifyStartup:              "",
			NoUserLoggedIn:              "Do not restart",
			UserLoggedIn:                "Do not restart",
			MinutesUntilReboot:          5,
			StartRebootTimerImmediately: false,
			FileVault2Reboot:            false,
		}

		// Check if values are provided in Terraform and override defaults if necessary
		if v, ok := d.GetOk("reboot"); ok {
			rebootData := v.(*schema.Set).List()[0].(map[string]interface{})

			// Override with provided values if available
			defaultReboot.Message = getStringFromMap(rebootData, "message")
			defaultReboot.SpecifyStartup = getStringFromMap(rebootData, "specify_startup")
			defaultReboot.StartupDisk = getStringFromMap(rebootData, "startup_disk")
			defaultReboot.NoUserLoggedIn = getStringFromMap(rebootData, "no_user_logged_in")
			defaultReboot.UserLoggedIn = getStringFromMap(rebootData, "user_logged_in")
			defaultReboot.MinutesUntilReboot = getIntFromMap(rebootData, "minutes_until_reboot")
			defaultReboot.StartRebootTimerImmediately = getBoolFromMap(rebootData, "start_reboot_timer_immediately")
			defaultReboot.FileVault2Reboot = getBoolFromMap(rebootData, "file_vault_2_reboot")
		}

		return defaultReboot
	}()

	// Construct the Maintenance section
	if v, ok := d.GetOk("maintenance"); ok {
		maintenanceData := v.([]interface{})[0].(map[string]interface{})
		policy.Maintenance = jamfpro.PolicyMaintenance{
			Recon:                    getBoolFromMap(maintenanceData, "recon"),
			ResetName:                getBoolFromMap(maintenanceData, "reset_name"),
			InstallAllCachedPackages: getBoolFromMap(maintenanceData, "install_all_cached_packages"),
			Heal:                     getBoolFromMap(maintenanceData, "heal"),
			Prebindings:              getBoolFromMap(maintenanceData, "prebindings"),
			Permissions:              getBoolFromMap(maintenanceData, "permissions"),
			Byhost:                   getBoolFromMap(maintenanceData, "byhost"),
			SystemCache:              getBoolFromMap(maintenanceData, "system_cache"),
			UserCache:                getBoolFromMap(maintenanceData, "user_cache"),
			Verify:                   getBoolFromMap(maintenanceData, "verify"),
		}
	}

	// Construct the FilesProcesses section
	if v, ok := d.GetOk("files_processes"); ok {
		filesProcessesData := v.([]interface{})[0].(map[string]interface{})
		policy.FilesProcesses = jamfpro.PolicyFilesProcesses{
			SearchByPath:         getStringFromMap(filesProcessesData, "search_by_path"),
			DeleteFile:           getBoolFromMap(filesProcessesData, "delete_file"),
			LocateFile:           getStringFromMap(filesProcessesData, "locate_file"),
			UpdateLocateDatabase: getBoolFromMap(filesProcessesData, "update_locate_database"),
			SpotlightSearch:      getStringFromMap(filesProcessesData, "spotlight_search"),
			SearchForProcess:     getStringFromMap(filesProcessesData, "search_for_process"),
			KillProcess:          getBoolFromMap(filesProcessesData, "kill_process"),
			RunCommand:           getStringFromMap(filesProcessesData, "run_command"),
		}
	}

	// Construct the UserInteraction section
	if v, ok := d.GetOk("user_interaction"); ok {
		userInteractionData := v.([]interface{})[0].(map[string]interface{})
		policy.UserInteraction = jamfpro.PolicyUserInteraction{
			MessageStart:          getStringFromMap(userInteractionData, "message_start"),
			AllowUserToDefer:      getBoolFromMap(userInteractionData, "allow_user_to_defer"),
			AllowDeferralUntilUtc: getStringFromMap(userInteractionData, "allow_deferral_until_utc"),
			AllowDeferralMinutes:  getIntFromMap(userInteractionData, "allow_deferral_minutes"),
			MessageFinish:         getStringFromMap(userInteractionData, "message_finish"),
		}
	}

	// Construct the DiskEncryption section
	policy.DiskEncryption = func() jamfpro.PolicyDiskEncryption {
		// Initialize with default values
		defaultDiskEncryption := jamfpro.PolicyDiskEncryption{
			Action:                                 "none",
			DiskEncryptionConfigurationID:          0,
			AuthRestart:                            false,
			RemediateKeyType:                       "",
			RemediateDiskEncryptionConfigurationID: 0,
		}

		// Check if values are provided in Terraform and override defaults if necessary
		if v, ok := d.GetOk("disk_encryption"); ok && len(v.([]interface{})) > 0 {
			diskEncryptionData := v.([]interface{})[0].(map[string]interface{})

			// Override with provided values if available
			defaultDiskEncryption.Action = getStringFromMap(diskEncryptionData, "action")
			defaultDiskEncryption.DiskEncryptionConfigurationID = getIntFromMap(diskEncryptionData, "disk_encryption_configuration_id")
			defaultDiskEncryption.AuthRestart = getBoolFromMap(diskEncryptionData, "auth_restart")
			defaultDiskEncryption.RemediateKeyType = getStringFromMap(diskEncryptionData, "remediate_key_type")
			defaultDiskEncryption.RemediateDiskEncryptionConfigurationID = getIntFromMap(diskEncryptionData, "remediate_disk_encryption_configuration_id")
		}

		return defaultDiskEncryption
	}()

	log.Printf("[INFO] Successfully constructed Jamf Pro Policy with name: %s", policy.General.Name)

	return policy, nil
}

// Helper function to safely get an int value from a map. Returns 0 if key is absent or nil.
func getIntFromMap(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok && val != nil {
		if intVal, ok := val.(int); ok {
			return intVal
		}
	}
	return 0 // Return default zero value if key is not found or nil
}

// Helper function to safely get a bool value from a map. Returns false if key is absent or nil.
func getBoolFromMap(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok && val != nil {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return false // Return default false value if key is not found or nil
}

// Helper function to safely get a string value from a map. Returns an empty string if key is absent or nil.
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok && val != nil {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return "" // Return default empty string if key is not found or nil
}

// Helper function to safely get an int64 value from a map. Returns 0 if key is absent, nil, or not an int64.
func getInt64FromMap(m map[string]interface{}, key string) int64 {
	if val, ok := m[key]; ok && val != nil {
		if int64Val, ok := val.(int64); ok {
			return int64Val
		}
	}
	return 0 // Return default zero value if key is not found, nil, or not an int64
}

// Helper function to safely get a string value from an array at a given index. Returns an empty string if index is out of range or value is not a string.
func getStringFromArray(arr []interface{}, index int) string {
	if len(arr) > index {
		if strVal, ok := arr[index].(string); ok {
			return strVal
		}
	}
	return "" // Return default empty string if index is out of range or value is not a string
}

// Helper function to generate diagnostics based on the error type.
func generateTFDiagsFromHTTPError(err error, d *schema.ResourceData, action string) diag.Diagnostics {
	var diags diag.Diagnostics
	resourceName, exists := d.GetOk("name")
	if !exists {
		resourceName = "unknown"
	}

	// Handle the APIError in the diagnostic
	if apiErr, ok := err.(*http_client.APIError); ok {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   fmt.Sprintf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message),
		})
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed to %s the resource with name: %s", action, resourceName),
			Detail:   err.Error(),
		})
	}
	return diags
}

// ResourceJamfProPoliciesCreate is responsible for reading the current state of a Jamf Pro policy from the remote system.
// The function:
// 1. Fetches the policies current state using its ID. If it fails then obtain profile's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the profile being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProPoliciesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the Policy
		policy, err := constructJamfProPolicy(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the policy for terraform create: %w", err))
		}

		// Log the details of the policy that is about to be created
		log.Printf("[INFO] Attempting to create Jamf Pro Policy with name: %s", policy.General.Name)

		// Call the API to create the policy and get its ID
		createdPolicy, err := conn.CreatePolicyByID(policy)
		if err != nil {
			log.Printf("[ERROR] Error creating Jamf Pro Policy with name: %s. Error: %s", policy.General.Name, err)
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			return retry.RetryableError(err)
		}

		// Log the successfully resource creation
		log.Printf("[INFO] Successfully created Jamf Pro Policy with ID: %d", createdPolicy.ID)

		// Set the ID in the Terraform state
		d.SetId(strconv.Itoa(createdPolicy.ID))

		return nil
	})

	if err != nil {
		return generateTFDiagsFromHTTPError(err, d, "create")
	}

	// Log the ID that was set in the Terraform state
	log.Printf("[INFO] Terraform state was successfully updated with new Jamf Pro Policy with ID: %s", d.Id())

	// Perform a read operation to update the Terraform state
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProPoliciesRead(ctx, d, meta)
		if len(readDiags) > 0 {
			return retry.RetryableError(fmt.Errorf("failed to read the created resource"))
		}
		return nil
	})

	if err != nil {
		return generateTFDiagsFromHTTPError(err, d, "update state for")
	}

	return diags
}

// ResourceJamfProPoliciesRead is responsible for reading the current state of a Jamf Pro policy Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProPoliciesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var policy *jamfpro.ResponsePolicy

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		policyID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse policy ID: %v", convertErr))
		}

		// Try fetching the policy using the ID
		var apiErr error
		policy, apiErr = conn.GetPolicyByID(policyID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				if apiError.StatusCode == 404 {
					// If fetching by ID fails, try fetching by Name
					policyName, ok := d.Get("name").(string)
					if !ok {
						return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
					}
					// If fetching by ID fails, try fetching by Name
					policy, apiErr = conn.GetPolicyByName(policyName)
					if apiErr != nil {
						return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
					}
				} else {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
			}
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while reading the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	// Update the Terraform state with the fetched data
	// Set 'general' attributes
	generalAttributes := map[string]interface{}{
		"name":                          policy.General.Name,
		"enabled":                       policy.General.Enabled,
		"trigger":                       policy.General.Trigger,
		"trigger_checkin":               policy.General.TriggerCheckin,
		"trigger_enrollment_complete":   policy.General.TriggerEnrollmentComplete,
		"trigger_login":                 policy.General.TriggerLogin,
		"trigger_logout":                policy.General.TriggerLogout,
		"trigger_network_state_changed": policy.General.TriggerNetworkStateChanged,
		"trigger_startup":               policy.General.TriggerStartup,
		"trigger_other":                 policy.General.TriggerOther,
		"frequency":                     policy.General.Frequency,
		"retry_event":                   policy.General.RetryEvent,
		"retry_attempts":                policy.General.RetryAttempts,
		"notify_on_each_failed_retry":   policy.General.NotifyOnEachFailedRetry,
		"location_user_only":            policy.General.LocationUserOnly,
		"target_drive":                  policy.General.TargetDrive,
		"offline":                       policy.General.Offline,
		"category": []interface{}{map[string]interface{}{
			"id":   policy.General.Category.ID,
			"name": policy.General.Category.Name,
		}},
		"date_time_limitations": []interface{}{
			map[string]interface{}{
				"activation_date":       policy.General.DateTimeLimitations.ActivationDate,
				"activation_date_epoch": policy.General.DateTimeLimitations.ActivationDateEpoch,
				"activation_date_utc":   policy.General.DateTimeLimitations.ActivationDateUTC,
				"expiration_date":       policy.General.DateTimeLimitations.ExpirationDate,
				"expiration_date_epoch": policy.General.DateTimeLimitations.ExpirationDateEpoch,
				"expiration_date_utc":   policy.General.DateTimeLimitations.ExpirationDateUTC,
				"no_execute_on": func() []interface{} {
					noExecOnDays := make([]interface{}, len(policy.General.DateTimeLimitations.NoExecuteOn))
					for i, noExecOn := range policy.General.DateTimeLimitations.NoExecuteOn {
						noExecOnDays[i] = map[string]interface{}{"day": noExecOn.Day}
					}
					return noExecOnDays
				}(),
				"no_execute_start": policy.General.DateTimeLimitations.NoExecuteStart,
				"no_execute_end":   policy.General.DateTimeLimitations.NoExecuteEnd,
			},
		},
		"network_limitations": []interface{}{map[string]interface{}{
			"minimum_network_connection": policy.General.NetworkLimitations.MinimumNetworkConnection,
			"any_ip_address":             policy.General.NetworkLimitations.AnyIPAddress,
			"network_segments":           policy.General.NetworkLimitations.NetworkSegments,
		}},
		"override_default_settings": []interface{}{map[string]interface{}{
			"target_drive":       policy.General.OverrideDefaultSettings.TargetDrive,
			"distribution_point": policy.General.OverrideDefaultSettings.DistributionPoint,
			"force_afp_smb":      policy.General.OverrideDefaultSettings.ForceAfpSmb,
			"sus":                policy.General.OverrideDefaultSettings.SUS,
		}},
		"network_requirements": policy.General.NetworkRequirements,
		"site": []interface{}{map[string]interface{}{
			"id":   policy.General.Site.ID,
			"name": policy.General.Site.Name,
		}},
	}

	if err := d.Set("general", []interface{}{generalAttributes}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set 'scope' attributes
	scopeAttributes := map[string]interface{}{
		"all_computers": policy.Scope.AllComputers,
		"computers": func() []interface{} {
			computersInterfaces := make([]interface{}, len(policy.Scope.Computers))
			for i, computer := range policy.Scope.Computers {
				computersInterfaces[i] = map[string]interface{}{
					"id":   computer.ID,
					"name": computer.Name,
					"udid": computer.UDID,
				}
			}
			return computersInterfaces
		}(),
		"computer_groups": func() []interface{} {
			groupInterfaces := make([]interface{}, len(policy.Scope.ComputerGroups))
			for i, group := range policy.Scope.ComputerGroups {
				groupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			return groupInterfaces
		}(),
		"jss_users": func() []interface{} {
			userInterfaces := make([]interface{}, len(policy.Scope.JSSUsers))
			for i, user := range policy.Scope.JSSUsers {
				userInterfaces[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			return userInterfaces
		}(),
		"jss_user_groups": func() []interface{} {
			userGroupInterfaces := make([]interface{}, len(policy.Scope.JSSUserGroups))
			for i, group := range policy.Scope.JSSUserGroups {
				userGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			return userGroupInterfaces
		}(),
		"buildings": func() []interface{} {
			buildingInterfaces := make([]interface{}, len(policy.Scope.Buildings))
			for i, building := range policy.Scope.Buildings {
				buildingInterfaces[i] = map[string]interface{}{
					"id":   building.ID,
					"name": building.Name,
				}
			}
			return buildingInterfaces
		}(),
		"departments": func() []interface{} {
			departmentInterfaces := make([]interface{}, len(policy.Scope.Departments))
			for i, department := range policy.Scope.Departments {
				departmentInterfaces[i] = map[string]interface{}{
					"id":   department.ID,
					"name": department.Name,
				}
			}
			return departmentInterfaces
		}(),
		"limitations": func() []interface{} {
			limitationInterfaces := make([]interface{}, 0)
			limitationData := map[string]interface{}{}

			// Network Segments
			networkSegmentInterfaces := make([]interface{}, len(policy.Scope.Limitations.NetworkSegments))
			for i, segment := range policy.Scope.Limitations.NetworkSegments {
				networkSegmentInterfaces[i] = map[string]interface{}{
					"id":   segment.ID,
					"name": segment.Name,
					"uid":  segment.UID,
				}
			}
			limitationData["network_segments"] = networkSegmentInterfaces

			// Users
			userInterfaces := make([]interface{}, len(policy.Scope.Limitations.Users))
			for i, user := range policy.Scope.Limitations.Users {
				userInterfaces[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			limitationData["users"] = userInterfaces

			// User Groups
			userGroupInterfaces := make([]interface{}, len(policy.Scope.Limitations.UserGroups))
			for i, group := range policy.Scope.Limitations.UserGroups {
				userGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			limitationData["user_groups"] = userGroupInterfaces

			// iBeacons
			iBeaconInterfaces := make([]interface{}, len(policy.Scope.Limitations.IBeacons))
			for i, beacon := range policy.Scope.Limitations.IBeacons {
				iBeaconInterfaces[i] = map[string]interface{}{
					"id":   beacon.ID,
					"name": beacon.Name,
				}
			}
			limitationData["ibeacons"] = iBeaconInterfaces

			limitationInterfaces = append(limitationInterfaces, limitationData)
			return limitationInterfaces
		}(),
		"exclusions": func() []interface{} {
			exclusionsInterfaces := make([]interface{}, 0)
			exclusionsData := map[string]interface{}{}

			// Computers
			computerInterfaces := make([]interface{}, len(policy.Scope.Exclusions.Computers))
			for i, computer := range policy.Scope.Exclusions.Computers {
				computerInterfaces[i] = map[string]interface{}{
					"id":   computer.ID,
					"name": computer.Name,
					"udid": computer.UDID,
				}
			}
			exclusionsData["computers"] = computerInterfaces

			// Computer Groups
			computerGroupInterfaces := make([]interface{}, len(policy.Scope.Exclusions.ComputerGroups))
			for i, group := range policy.Scope.Exclusions.ComputerGroups {
				computerGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			exclusionsData["computer_groups"] = computerGroupInterfaces

			// Users
			userInterfaces := make([]interface{}, len(policy.Scope.Exclusions.Users))
			for i, user := range policy.Scope.Exclusions.Users {
				userInterfaces[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			exclusionsData["users"] = userInterfaces

			// User Groups
			userGroupInterfaces := make([]interface{}, len(policy.Scope.Exclusions.UserGroups))
			for i, group := range policy.Scope.Exclusions.UserGroups {
				userGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			exclusionsData["user_groups"] = userGroupInterfaces

			// Buildings
			buildingInterfaces := make([]interface{}, len(policy.Scope.Exclusions.Buildings))
			for i, building := range policy.Scope.Exclusions.Buildings {
				buildingInterfaces[i] = map[string]interface{}{
					"id":   building.ID,
					"name": building.Name,
				}
			}
			exclusionsData["buildings"] = buildingInterfaces

			// Departments
			departmentInterfaces := make([]interface{}, len(policy.Scope.Exclusions.Departments))
			for i, department := range policy.Scope.Exclusions.Departments {
				departmentInterfaces[i] = map[string]interface{}{
					"id":   department.ID,
					"name": department.Name,
				}
			}
			exclusionsData["departments"] = departmentInterfaces

			// Network Segments
			networkSegmentInterfaces := make([]interface{}, len(policy.Scope.Exclusions.NetworkSegments))
			for i, segment := range policy.Scope.Exclusions.NetworkSegments {
				networkSegmentInterfaces[i] = map[string]interface{}{
					"id":   segment.ID,
					"name": segment.Name,
					"uid":  segment.UID,
				}
			}
			exclusionsData["network_segments"] = networkSegmentInterfaces

			// JSS Users
			jssUserInterfaces := make([]interface{}, len(policy.Scope.Exclusions.JSSUsers))
			for i, user := range policy.Scope.Exclusions.JSSUsers {
				jssUserInterfaces[i] = map[string]interface{}{
					"id":   user.ID,
					"name": user.Name,
				}
			}
			exclusionsData["jss_users"] = jssUserInterfaces

			// JSS User Groups
			jssUserGroupInterfaces := make([]interface{}, len(policy.Scope.Exclusions.JSSUserGroups))
			for i, group := range policy.Scope.Exclusions.JSSUserGroups {
				jssUserGroupInterfaces[i] = map[string]interface{}{
					"id":   group.ID,
					"name": group.Name,
				}
			}
			exclusionsData["jss_user_groups"] = jssUserGroupInterfaces

			// IBeacons
			iBeaconInterfaces := make([]interface{}, len(policy.Scope.Exclusions.IBeacons))
			for i, beacon := range policy.Scope.Exclusions.IBeacons {
				iBeaconInterfaces[i] = map[string]interface{}{
					"id":   beacon.ID,
					"name": beacon.Name,
				}
			}
			exclusionsData["ibeacons"] = iBeaconInterfaces

			exclusionsInterfaces = append(exclusionsInterfaces, exclusionsData)
			return exclusionsInterfaces
		}(),
	}

	if err := d.Set("scope", []interface{}{scopeAttributes}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set 'self_service' attributes
	selfServiceAttributes := map[string]interface{}{
		"use_for_self_service":            policy.SelfService.UseForSelfService,
		"self_service_display_name":       policy.SelfService.SelfServiceDisplayName,
		"install_button_text":             policy.SelfService.InstallButtonText,
		"reinstall_button_text":           policy.SelfService.ReinstallButtonText,
		"self_service_description":        policy.SelfService.SelfServiceDescription,
		"force_users_to_view_description": policy.SelfService.ForceUsersToViewDescription,
		"self_service_icon": []interface{}{map[string]interface{}{
			"id":       policy.SelfService.SelfServiceIcon.ID,
			"filename": policy.SelfService.SelfServiceIcon.Filename,
			"uri":      policy.SelfService.SelfServiceIcon.URI,
		}},
		"feature_on_main_page": policy.SelfService.FeatureOnMainPage,
		"self_service_categories": func() []interface{} {
			categories := make([]interface{}, len(policy.SelfService.SelfServiceCategories))
			for i, cat := range policy.SelfService.SelfServiceCategories {
				categoryMap := map[string]interface{}{
					"id":         cat.Category.ID,
					"name":       cat.Category.Name,
					"display_in": cat.Category.DisplayIn,
					"feature_in": cat.Category.FeatureIn,
				}
				categories[i] = map[string]interface{}{
					"category": []interface{}{categoryMap},
				}
			}
			return categories
		}(),
	}

	if err := d.Set("self_service", []interface{}{selfServiceAttributes}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Fetch the package configuration from the policy and set it in Terraform state
	packageConfigurations := make([]interface{}, 0)
	for _, packageItem := range policy.PackageConfiguration.Packages {
		pkg := make(map[string]interface{})
		pkg["id"] = packageItem.ID
		pkg["name"] = packageItem.Name
		pkg["action"] = packageItem.Action
		pkg["fut"] = packageItem.FillUserTemplate
		pkg["feu"] = packageItem.FillExistingUsers
		pkg["update_autorun"] = packageItem.UpdateAutorun
		packageConfigurations = append(packageConfigurations, pkg)
	}

	// Wrap packageConfigurations in a map under the key 'packages'
	packageConfiguration := map[string]interface{}{
		"packages": packageConfigurations,
	}

	// Wrap this map in a slice to set in the Terraform state
	if err := d.Set("package_configuration", []interface{}{packageConfiguration}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the scripts from the policy and set them in Terraform state
	scriptConfigurations := make([]interface{}, 0)
	for _, scriptItem := range policy.Scripts.Script {
		script := make(map[string]interface{})
		script["id"] = scriptItem.ID
		script["name"] = scriptItem.Name
		script["priority"] = scriptItem.Priority
		script["parameter4"] = scriptItem.Parameter4
		script["parameter5"] = scriptItem.Parameter5
		script["parameter6"] = scriptItem.Parameter6
		script["parameter7"] = scriptItem.Parameter7
		script["parameter8"] = scriptItem.Parameter8
		script["parameter9"] = scriptItem.Parameter9
		script["parameter10"] = scriptItem.Parameter10
		script["parameter11"] = scriptItem.Parameter11
		scriptConfigurations = append(scriptConfigurations, script)
	}

	if err := d.Set("scripts", scriptConfigurations); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the printers from the policy and set them in Terraform state
	printerConfigurations := make([]interface{}, 0)
	for _, printerItem := range policy.Printers.Printer {
		printer := make(map[string]interface{})
		printer["id"] = printerItem.ID
		printer["name"] = printerItem.Name
		printer["action"] = printerItem.Action
		printer["make_default"] = printerItem.MakeDefault
		printerConfigurations = append(printerConfigurations, printer)
	}

	if err := d.Set("printers", []interface{}{
		map[string]interface{}{
			"leave_existing_default": policy.Printers.LeaveExistingDefault,
			"printer":                printerConfigurations,
		},
	}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the dock items from the policy and set them in Terraform state
	dockItemConfigurations := make([]interface{}, 0)
	for _, dockItem := range policy.DockItems.DockItem {
		dock := make(map[string]interface{})
		dock["id"] = dockItem.ID
		dock["name"] = dockItem.Name
		dock["action"] = dockItem.Action
		dockItemConfigurations = append(dockItemConfigurations, dock)
	}

	if err := d.Set("dock_items", []interface{}{
		map[string]interface{}{
			"dock_item": dockItemConfigurations,
		},
	}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the account maintenance data from the policy and set it in Terraform state
	accountMaintenanceState := make(map[string]interface{})
	accountMaintenanceState["accounts"] = []interface{}{}

	// Add account data if present
	if len(policy.AccountMaintenance.Accounts) > 0 {
		accountsState := make([]interface{}, len(policy.AccountMaintenance.Accounts))
		for i, account := range policy.AccountMaintenance.Accounts {
			accountMap := map[string]interface{}{
				"action":                    account.Action,
				"username":                  account.Username,
				"realname":                  account.Realname,
				"password":                  account.Password,
				"archive_home_directory":    account.ArchiveHomeDirectory,
				"archive_home_directory_to": account.ArchiveHomeDirectoryTo,
				"home":                      account.Home,
				"hint":                      account.Hint,
				"picture":                   account.Picture,
				"admin":                     account.Admin,
				"filevault_enabled":         account.FilevaultEnabled,
			}
			accountsState[i] = map[string]interface{}{"account": accountMap}
		}
		accountMaintenanceState["accounts"] = accountsState
	}

	// Add directory bindings data if present
	if len(policy.AccountMaintenance.DirectoryBindings) > 0 {
		bindingsState := make([]interface{}, len(policy.AccountMaintenance.DirectoryBindings))
		for i, binding := range policy.AccountMaintenance.DirectoryBindings {
			bindingMap := map[string]interface{}{
				"id":   binding.ID,
				"name": binding.Name,
			}
			bindingsState[i] = map[string]interface{}{"binding": bindingMap}
		}
		accountMaintenanceState["directory_bindings"] = bindingsState
	}

	// Add management account data
	accountMaintenanceState["management_account"] = []interface{}{
		map[string]interface{}{
			"action":                  policy.AccountMaintenance.ManagementAccount.Action,
			"managed_password":        policy.AccountMaintenance.ManagementAccount.ManagedPassword,
			"managed_password_length": policy.AccountMaintenance.ManagementAccount.ManagedPasswordLength,
		},
	}

	// Add open firmware/EFI password data
	accountMaintenanceState["open_firmware_efi_password"] = []interface{}{
		map[string]interface{}{
			"of_mode":     policy.AccountMaintenance.OpenFirmwareEfiPassword.OfMode,
			"of_password": policy.AccountMaintenance.OpenFirmwareEfiPassword.OfPassword,
		},
	}

	// Set the account_maintenance in state
	if err := d.Set("account_maintenance", []interface{}{accountMaintenanceState}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the reboot data from the policy and set it in Terraform state
	rebootConfig := make(map[string]interface{})
	rebootConfig["message"] = policy.Reboot.Message
	rebootConfig["specify_startup"] = policy.Reboot.SpecifyStartup
	rebootConfig["startup_disk"] = policy.Reboot.StartupDisk
	rebootConfig["no_user_logged_in"] = policy.Reboot.NoUserLoggedIn
	rebootConfig["user_logged_in"] = policy.Reboot.UserLoggedIn
	rebootConfig["minutes_until_reboot"] = policy.Reboot.MinutesUntilReboot
	rebootConfig["start_reboot_timer_immediately"] = policy.Reboot.StartRebootTimerImmediately
	rebootConfig["file_vault_2_reboot"] = policy.Reboot.FileVault2Reboot

	if err := d.Set("reboot", []interface{}{rebootConfig}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the maintenance data from the policy and set it in Terraform state
	maintenanceConfig := make(map[string]interface{})
	maintenanceConfig["recon"] = policy.Maintenance.Recon
	maintenanceConfig["reset_name"] = policy.Maintenance.ResetName
	maintenanceConfig["install_all_cached_packages"] = policy.Maintenance.InstallAllCachedPackages
	maintenanceConfig["heal"] = policy.Maintenance.Heal
	maintenanceConfig["prebindings"] = policy.Maintenance.Prebindings
	maintenanceConfig["permissions"] = policy.Maintenance.Permissions
	maintenanceConfig["byhost"] = policy.Maintenance.Byhost
	maintenanceConfig["system_cache"] = policy.Maintenance.SystemCache
	maintenanceConfig["user_cache"] = policy.Maintenance.UserCache
	maintenanceConfig["verify"] = policy.Maintenance.Verify

	if err := d.Set("maintenance", []interface{}{maintenanceConfig}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the files and processes data from the policy and set it in Terraform state
	filesProcessesConfig := make(map[string]interface{})
	filesProcessesConfig["search_by_path"] = policy.FilesProcesses.SearchByPath
	filesProcessesConfig["delete_file"] = policy.FilesProcesses.DeleteFile
	filesProcessesConfig["locate_file"] = policy.FilesProcesses.LocateFile
	filesProcessesConfig["update_locate_database"] = policy.FilesProcesses.UpdateLocateDatabase
	filesProcessesConfig["spotlight_search"] = policy.FilesProcesses.SpotlightSearch
	filesProcessesConfig["search_for_process"] = policy.FilesProcesses.SearchForProcess
	filesProcessesConfig["kill_process"] = policy.FilesProcesses.KillProcess
	filesProcessesConfig["run_command"] = policy.FilesProcesses.RunCommand

	if err := d.Set("files_processes", []interface{}{filesProcessesConfig}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the user interaction data from the policy and set it in Terraform state
	userInteractionConfig := make(map[string]interface{})
	userInteractionConfig["message_start"] = policy.UserInteraction.MessageStart
	userInteractionConfig["allow_user_to_defer"] = policy.UserInteraction.AllowUserToDefer
	userInteractionConfig["allow_deferral_until_utc"] = policy.UserInteraction.AllowDeferralUntilUtc
	userInteractionConfig["allow_deferral_minutes"] = policy.UserInteraction.AllowDeferralMinutes
	userInteractionConfig["message_finish"] = policy.UserInteraction.MessageFinish

	if err := d.Set("user_interaction", []interface{}{userInteractionConfig}); err != nil {
		return diag.FromErr(err)
	}

	// Fetch the disk encryption data from the policy and set it in Terraform state
	diskEncryptionConfig := make(map[string]interface{})
	diskEncryptionConfig["action"] = policy.DiskEncryption.Action
	diskEncryptionConfig["disk_encryption_configuration_id"] = policy.DiskEncryption.DiskEncryptionConfigurationID
	diskEncryptionConfig["auth_restart"] = policy.DiskEncryption.AuthRestart
	diskEncryptionConfig["remediate_key_type"] = policy.DiskEncryption.RemediateKeyType
	diskEncryptionConfig["remediate_disk_encryption_configuration_id"] = policy.DiskEncryption.RemediateDiskEncryptionConfigurationID

	if err := d.Set("disk_encryption", []interface{}{diskEncryptionConfig}); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// ResourceJamfProPoliciesUpdate is responsible for updating an existing Jamf Pro policy on the remote system.
func ResourceJamfProPoliciesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the update operation
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the updated jamf pro policy
		policy, err := constructJamfProPolicy(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the configuration policy for terraform update: %w", err))
		}

		// Convert the ID from the Terraform state into an integer to be used for the API request
		policyID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse policy ID: %v", convertErr))
		}

		// Directly call the API to update the policy
		_, apiErr := conn.UpdatePolicyByID(policyID, policy)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by name
			policyName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			_, apiErr = conn.UpdatePolicyByName(policyName, policy)
			if apiErr != nil {
				// Handle the APIError
				if apiError, ok := apiErr.(*http_client.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	// Use the retry function for the read operation to update the Terraform state
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProPoliciesRead(ctx, d, meta)
		if len(readDiags) > 0 {
			return retry.RetryableError(fmt.Errorf("failed to update the Terraform state for the updated resource"))
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while updating the Terraform state, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update")
	}

	return diags
}

// ResourceJamfProPoliciesDelete is responsible for deleting a Jamf Pro policy.
func ResourceJamfProPoliciesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the DELETE operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		policyID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse dock item ID: %v", convertErr))
		}

		// Directly call the API to DELETE the resource
		apiErr := conn.DeletePolicyByID(policyID)
		if apiErr != nil {
			// If the DELETE by ID fails, try deleting by name
			policyName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			apiErr = conn.DeletePolicyByName(policyName)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while deleting the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "delete")
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
