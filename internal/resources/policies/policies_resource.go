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
			Create: schema.DefaultTimeout(3 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
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
						},
						"trigger": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Event(s) triggers to use to initiate the policy.",
						},
						"trigger_checkin": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Trigger policy when device performs recurring check-in against the frequency configured in Jamf Pro",
						},
						"trigger_enrollment_complete": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Trigger policy when device enrollment is complete.",
						},
						"trigger_login": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Trigger policy when a user logs in to a computer. A login event that checks for policies must be configured in Jamf Pro for this to work",
						},
						"trigger_logout": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Trigger policy when a user logout.",
						},
						"trigger_network_state_changed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Trigger policy when it's network state changes. When a computer's network state changes (e.g., when the network connection changes, when the computer name changes, when the IP address changes)",
						},
						"trigger_startup": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Trigger policy when a computer starts up. A startup script that checks for policies must be configured in Jamf Pro for this to work",
						},
						"trigger_other": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Any other trigger for the policy.",
						},
						"frequency": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Frequency of policy execution.",
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
						},
						"notify_on_each_failed_retry": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Send notifications for each failed policy retry attempt. ",
						},
						"location_user_only": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Location-based policy for user only.",
						},
						"target_drive": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The drive on which to run the policy (e.g. /Volumes/Restore/ ). The policy runs on the boot drive by default",
						},
						"offline": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the policy applies when offline.",
						},
						"category": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Category to add the policy to.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "-1",
										Description: "The category ID assigned to the jamf pro policy. Defaults to '-1' aka not used.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "No category assigned",
										Description: "Category Name for assigned jamf pro policy. Value defaults to 'No category assigned' aka not used",
									},
								},
							},
						},
						"date_time_limitations": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Server-side limitations use your Jamf Pro host server's time zone and settings. The Jamf Pro host service is in UTC time.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"activation_date": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The activation date of the policy.",
									},
									"activation_date_epoch": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The epoch time of the activation date.",
									},
									"activation_date_utc": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The UTC time of the activation date.",
									},
									"expiration_date": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The expiration date of the policy.",
									},
									"expiration_date_epoch": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The epoch time of the expiration date.",
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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"minimum_network_connection": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Minimum network connection required for the policy.",
										ValidateFunc: validation.StringInSlice([]string{"Any", "Ethernet"}, false),
									},
									"any_ip_address": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the policy applies to any IP address.",
									},
								},
							},
						},
						"override_default_settings": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Settings to override default configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_drive": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Target drive for the policy.",
									},
									"distribution_point": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Distribution point for the policy.",
									},
									"force_afp_smb": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to force AFP/SMB.",
									},
									"sus": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Software Update Service for the policy.",
									},
								},
							},
						},
						"network_requirements": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Network requirements for the policy.",
						},
						"site": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Jamf Pro Site-related settings of the policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "-1", // Set default value as string "-1"
										Description: "Jamf Pro Site ID. Value defaults to -1 aka not used.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "None", // Set default value as "None"
										Description: "Jamf Pro Site Name. Value defaults to 'None' aka not used",
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
				Optional:    true,
				Description: "Scope configuration for the profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_computers": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If true, applies the profile to all computers.",
						},
						"all_jss_users": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If true, applies the profile to all JSS users.",
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
													Optional:    true,
													Description: "The unique identifier of the scoped computer.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the scoped computer.",
												},
												"udid": {
													Type:        schema.TypeString,
													Optional:    true,
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
													Optional:    true,
													Description: "The unique identifier of the scoped building.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the scoped building.",
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
													Optional:    true,
													Description: "The unique identifier of the scoped department.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the scoped department.",
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
													Optional:    true,
													Description: "The unique identifier of the scoped computer group.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the computer scoped group.",
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
													Optional:    true,
													Description: "The unique identifier of the scoped JSS user.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the scoped JSS user.",
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
										Optional:    true,
										Description: "The unique identifier of the scoped JSS user group.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name of the scoped JSS user group.",
									},
								},
							},
						},
						"limitations": {
							Type:     schema.TypeList,
							Optional: true,
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
																Optional:    true,
																Description: "The unique identifier of the scoped network segment.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the scoped network segment.",
															},
															"uid": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "UID of the scoped network segment.",
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
																Optional:    true,
																Description: "The unique identifier of the user.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the user.",
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
																Optional:    true,
																Description: "The unique identifier of the scoped user group.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the scoped user group.",
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
																Optional:    true,
																Description: "The unique identifier of the scoped iBeacon.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the scoped iBeacon.",
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
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
																Optional:    true,
																Description: "The unique identifier of the computer.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the computer.",
															},
															"udid": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "UDID of the computer.",
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
																Optional:    true,
																Description: "The unique identifier of the computer group.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the computer group.",
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
																Optional:    true,
																Description: "The unique identifier of the JSS user.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the JSS user.",
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
																Optional:    true,
																Description: "The unique identifier of the JSS user group.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the JSS user group.",
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
																Optional:    true,
																Description: "The unique identifier of the building.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the building.",
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
																Optional:    true,
																Description: "The unique identifier of the department.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the department.",
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
																Optional:    true,
																Description: "The unique identifier of the network segment.",
															},
															"uid": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "UID of the network segment.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the network segment.",
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
																Optional:    true,
																Description: "The unique identifier of the user.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the user.",
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
																Optional:    true,
																Description: "The unique identifier of the user group.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the user group.",
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
																Optional:    true,
																Description: "The unique identifier of the iBeacon.",
															},
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the iBeacon.",
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
				Optional:    true,
				Description: "Self-service settings of the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_for_self_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the policy is available for self-service.",
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
						},
						"re_install_button_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Text displayed on the re-install button in self-service.",
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
						},
						"self_service_icon": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Icon settings for the policy in self-service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "ID of the icon used in self-service.",
									},
									"filename": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filename of the icon used in self-service.",
									},
									"uri": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "URI of the icon used in self-service.",
									},
								},
							},
						},
						"feature_on_main_page": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to feature the policy on the main page of self-service.",
						},
						"self_service_categories": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Category settings for the policy in self-service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Category details for the policy in self-service.",
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
												},
												"feature_in": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to feature the category in self-service.",
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
									"size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The number of packages included in the policy.",
									},
									"package": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Details of the package.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Unique identifier of the package.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the package.",
												},
												"action": {
													Type:         schema.TypeString,
													Optional:     true,
													Description:  "Action to be performed for the package.",
													ValidateFunc: validation.StringInSlice([]string{"Install", "Cache", "Install Cached"}, false),
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
									},
									"parameter4": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Custom parameter 4 for the script.",
									},
									"parameter5": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Custom parameter 5 for the script.",
									},
									"parameter6": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Custom parameter 6 for the script.",
									},
									"parameter7": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Custom parameter 7 for the script.",
									},
									"parameter8": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Custom parameter 8 for the script.",
									},
									"parameter9": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Custom parameter 9 for the script.",
									},
									"parameter10": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Custom parameter 10 for the script.",
									},
									"parameter11": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
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
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of printer configurations in the policy.",
						},
						"leave_existing_default": {
							Type:        schema.TypeString,
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
										Optional:    true,
										Description: "Unique identifier of the printer.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name of the printer.",
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of dock items in the policy.",
						},
						"dock_item": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Details of the dock item configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Unique identifier of the dock item.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name of the dock item.",
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
									"size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of accounts in the policy.",
									},
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
									"size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of directory bindings.",
									},
									"binding": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Details of the directory binding.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The unique identifier of the binding.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The name of the binding.",
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The action to perform for disk encryption (e.g., apply, remediate).",
							ValidateFunc: validation.StringInSlice([]string{"apply", "remediate"}, false),
						},
						"disk_encryption_configuration_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "ID of the disk encryption configuration to apply.",
						},
						"auth_restart": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to allow authentication restart.",
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
			ID:                         generalData["id"].(int),
			Name:                       generalData["name"].(string),
			Enabled:                    generalData["enabled"].(bool),
			Trigger:                    generalData["trigger"].(string),
			TriggerCheckin:             generalData["trigger_checkin"].(bool),
			TriggerEnrollmentComplete:  generalData["trigger_enrollment_complete"].(bool),
			TriggerLogin:               generalData["trigger_login"].(bool),
			TriggerLogout:              generalData["trigger_logout"].(bool),
			TriggerNetworkStateChanged: generalData["trigger_network_state_changed"].(bool),
			TriggerStartup:             generalData["trigger_startup"].(bool),
			TriggerOther:               generalData["trigger_other"].(string),
			Frequency:                  generalData["frequency"].(string),
			RetryEvent:                 generalData["retry_event"].(string),
			RetryAttempts:              generalData["retry_attempts"].(int),
			NotifyOnEachFailedRetry:    generalData["notify_on_each_failed_retry"].(bool),
			LocationUserOnly:           generalData["location_user_only"].(bool),
			TargetDrive:                generalData["target_drive"].(string),
			Offline:                    generalData["offline"].(bool),
			Category: func() jamfpro.PolicyCategory {
				if catData, ok := generalData["category"].([]interface{}); ok && len(catData) > 0 {
					catMap := catData[0].(map[string]interface{})
					return jamfpro.PolicyCategory{
						ID:   catMap["id"].(string),
						Name: catMap["name"].(string),
					}
				}
				return jamfpro.PolicyCategory{}
			}(),
			// DateTimeLimitations field
			DateTimeLimitations: func() jamfpro.PolicyDateTimeLimitations {
				if dtData, ok := generalData["date_time_limitations"].([]interface{}); ok && len(dtData) > 0 {
					dateTimeMap := dtData[0].(map[string]interface{})
					dateTimeLimitations := jamfpro.PolicyDateTimeLimitations{
						ActivationDate:      dateTimeMap["activation_date"].(string),
						ActivationDateEpoch: dateTimeMap["activation_date_epoch"].(int64),
						ActivationDateUTC:   dateTimeMap["activation_date_utc"].(string),
						ExpirationDate:      dateTimeMap["expiration_date"].(string),
						ExpirationDateEpoch: dateTimeMap["expiration_date_epoch"].(int64),
						ExpirationDateUTC:   dateTimeMap["expiration_date_utc"].(string),
					}

					// Handling NoExecuteOn field
					if noExecOn, ok := dateTimeMap["no_execute_on"].([]interface{}); ok && len(noExecOn) > 0 {
						// Assuming no_execute_on is a single string value
						policy.General.DateTimeLimitations.NoExecuteOn = jamfpro.PolicyNoExecuteOn{
							Day: noExecOn[0].(string),
						}
					}

					// Handling NoExecuteStart and NoExecuteEnd fields
					if noExecStart, ok := dateTimeMap["no_execute_start"].(string); ok {
						dateTimeLimitations.NoExecuteStart = noExecStart
					}
					if noExecEnd, ok := dateTimeMap["no_execute_end"].(string); ok {
						dateTimeLimitations.NoExecuteEnd = noExecEnd
					}

					return dateTimeLimitations
				}
				return jamfpro.PolicyDateTimeLimitations{}
			}(),
			// NetworkLimitations field
			NetworkLimitations: func() jamfpro.PolicyNetworkLimitations {
				if networkLimitationsData, ok := generalData["network_limitations"].([]interface{}); ok && len(networkLimitationsData) > 0 {
					netMap := networkLimitationsData[0].(map[string]interface{})
					networkLimitations := jamfpro.PolicyNetworkLimitations{
						MinimumNetworkConnection: netMap["minimum_network_connection"].(string),
						AnyIPAddress:             netMap["any_ip_address"].(bool),
					}

					// Handling NetworkSegments field
					if networkSegments, ok := netMap["network_segments"].(string); ok {
						networkLimitations.NetworkSegments = networkSegments
					}

					return networkLimitations
				}
				return jamfpro.PolicyNetworkLimitations{}
			}(),
			// OverrideDefaultSettings field
			OverrideDefaultSettings: func() jamfpro.PolicyOverrideSettings {
				if overrideData, ok := generalData["override_default_settings"].([]interface{}); ok && len(overrideData) > 0 {
					overrideMap := overrideData[0].(map[string]interface{})
					overrideSettings := jamfpro.PolicyOverrideSettings{
						TargetDrive:       overrideMap["target_drive"].(string),
						DistributionPoint: overrideMap["distribution_point"].(string),
						ForceAfpSmb:       overrideMap["force_afp_smb"].(bool),
						SUS:               overrideMap["sus"].(string),
						NetbootServer:     overrideMap["netboot_server"].(string),
					}

					return overrideSettings
				}
				return jamfpro.PolicyOverrideSettings{}
			}(),
			// NetworkRequirements field
			NetworkRequirements: generalData["network_requirements"].(string),
			// Site field
			Site: func() jamfpro.PolicySite {
				if siteData, ok := generalData["site"].([]interface{}); ok && len(siteData) > 0 {
					siteMap := siteData[0].(map[string]interface{})
					return jamfpro.PolicySite{
						ID:   siteMap["id"].(int),
						Name: siteMap["name"].(string),
					}
				}
				return jamfpro.PolicySite{}
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
			AllComputers:   scopeData["all_computers"].(bool),
			AllJSSUsers:    scopeData["all_jss_users"].(bool),
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
			UseForSelfService:           selfServiceData["use_for_self_service"].(bool),
			SelfServiceDisplayName:      selfServiceData["self_service_display_name"].(string),
			InstallButtonText:           selfServiceData["install_button_text"].(string),
			ReinstallButtonText:         selfServiceData["re_install_button_text"].(string),
			SelfServiceDescription:      selfServiceData["self_service_description"].(string),
			ForceUsersToViewDescription: selfServiceData["force_users_to_view_description"].(bool),
			FeatureOnMainPage:           selfServiceData["feature_on_main_page"].(bool),
			SelfServiceIcon: func() jamfpro.PolicySelfServiceIcon {
				if iconData, ok := selfServiceData["self_service_icon"].([]interface{}); ok && len(iconData) > 0 {
					iconMap := iconData[0].(map[string]interface{})
					return jamfpro.PolicySelfServiceIcon{
						ID:       iconMap["id"].(int),
						Filename: iconMap["filename"].(string),
						URI:      iconMap["uri"].(string),
					}
				}
				return jamfpro.PolicySelfServiceIcon{}
			}(),
			SelfServiceCategories: func() []jamfpro.PolicySelfServiceCategory {
				categories := []jamfpro.PolicySelfServiceCategory{}
				if catData, ok := selfServiceData["self_service_categories"].([]interface{}); ok {
					for _, cat := range catData {
						catMap := cat.(map[string]interface{})
						category := jamfpro.PolicySelfServiceCategory{
							Category: jamfpro.PolicyCategory{
								ID:        catMap["id"].(string),
								Name:      catMap["name"].(string),
								DisplayIn: catMap["display_in"].(bool),
								FeatureIn: catMap["feature_in"].(bool),
							},
						}
						categories = append(categories, category)
					}
				}
				return categories
			}(),
		}
	}

	// Construct the PackageConfiguration section
	if v, ok := d.GetOk("package_configuration"); ok {
		packageConfigData := v.([]interface{})[0].(map[string]interface{})
		var packageItems []jamfpro.PolicyPackage

		if pkgs, ok := packageConfigData["packages"].([]interface{}); ok {
			for _, pkg := range pkgs {
				pkgMap := pkg.(map[string]interface{})
				packageItem := jamfpro.PolicyPackage{
					ID:                pkgMap["id"].(int),
					Name:              pkgMap["name"].(string),
					Action:            pkgMap["action"].(string),
					FillUserTemplate:  pkgMap["fut"].(bool),
					FillExistingUsers: pkgMap["feu"].(bool),
					UpdateAutorun:     pkgMap["update_autorun"].(bool),
				}
				packageItems = append(packageItems, packageItem)
			}
		}

		policy.PackageConfiguration = jamfpro.PolicyPackageConfiguration{
			Packages: packageItems,
		}
	}

	// Construct the Scripts section
	if v, ok := d.GetOk("scripts"); ok {
		scriptsData := v.([]interface{})[0].(map[string]interface{})
		var scriptItems []jamfpro.PolicyScriptItem

		if scripts, ok := scriptsData["script"].([]interface{}); ok {
			for _, script := range scripts {
				scriptMap := script.(map[string]interface{})
				scriptItem := jamfpro.PolicyScriptItem{
					ID:          scriptMap["id"].(string),
					Name:        scriptMap["name"].(string),
					Priority:    scriptMap["priority"].(string),
					Parameter4:  scriptMap["parameter4"].(string),
					Parameter5:  scriptMap["parameter5"].(string),
					Parameter6:  scriptMap["parameter6"].(string),
					Parameter7:  scriptMap["parameter7"].(string),
					Parameter8:  scriptMap["parameter8"].(string),
					Parameter9:  scriptMap["parameter9"].(string),
					Parameter10: scriptMap["parameter10"].(string),
					Parameter11: scriptMap["parameter11"].(string),
				}
				scriptItems = append(scriptItems, scriptItem)
			}
		}

		policy.Scripts = jamfpro.PolicyScripts{
			Size:   scriptsData["size"].(int),
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
				printerItem := jamfpro.PolicyPrinterItem{
					ID:          printerMap["id"].(int),
					Name:        printerMap["name"].(string),
					Action:      printerMap["action"].(string),
					MakeDefault: printerMap["make_default"].(bool),
				}
				printerItems = append(printerItems, printerItem)
			}
		}

		policy.Printers = jamfpro.PolicyPrinters{
			Size:                 printersData["size"].(int),
			LeaveExistingDefault: printersData["leave_existing_default"].(bool),
			Printer:              printerItems,
		}
	}

	// Construct the DockItems section
	if v, ok := d.GetOk("dock_items"); ok {
		dockItemsData := v.([]interface{})[0].(map[string]interface{})
		var dockItems []jamfpro.PolicyDockItem

		if docks, ok := dockItemsData["dock_item"].([]interface{}); ok {
			for _, dock := range docks {
				dockMap := dock.(map[string]interface{})
				dockItem := jamfpro.PolicyDockItem{
					ID:     dockMap["id"].(int),
					Name:   dockMap["name"].(string),
					Action: dockMap["action"].(string),
				}
				dockItems = append(dockItems, dockItem)
			}
		}

		policy.DockItems = jamfpro.PolicyDockItems{
			Size:     dockItemsData["size"].(int),
			DockItem: dockItems,
		}
	}

	// Construct the AccountMaintenance section
	if v, ok := d.GetOk("account_maintenance"); ok {
		accountMaintenanceData := v.([]interface{})[0].(map[string]interface{})
		var accounts []jamfpro.PolicyAccount
		var directoryBindings []jamfpro.PolicyDirectoryBinding

		// Construct Accounts slice
		if accs, ok := accountMaintenanceData["accounts"].([]interface{}); ok {
			for _, acc := range accs {
				accMap := acc.(map[string]interface{})
				account := jamfpro.PolicyAccount{
					Action:                 accMap["action"].(string),
					Username:               accMap["username"].(string),
					Realname:               accMap["realname"].(string),
					Password:               accMap["password"].(string),
					ArchiveHomeDirectory:   accMap["archive_home_directory"].(bool),
					ArchiveHomeDirectoryTo: accMap["archive_home_directory_to"].(string),
					Home:                   accMap["home"].(string),
					Hint:                   accMap["hint"].(string),
					Picture:                accMap["picture"].(string),
					Admin:                  accMap["admin"].(bool),
					FilevaultEnabled:       accMap["filevault_enabled"].(bool),
				}
				accounts = append(accounts, account)
			}
		}

		// Construct DirectoryBindings slice
		if bindings, ok := accountMaintenanceData["directory_bindings"].([]interface{}); ok {
			for _, binding := range bindings {
				bindingMap := binding.(map[string]interface{})
				directoryBinding := jamfpro.PolicyDirectoryBinding{
					ID:   bindingMap["id"].(int),
					Name: bindingMap["name"].(string),
				}
				directoryBindings = append(directoryBindings, directoryBinding)
			}
		}

		// Construct ManagementAccount field
		var managementAccount jamfpro.PolicyManagementAccount
		if maData, ok := accountMaintenanceData["management_account"].(map[string]interface{}); ok {
			managementAccount = jamfpro.PolicyManagementAccount{
				Action:                maData["action"].(string),
				ManagedPassword:       maData["managed_password"].(string),
				ManagedPasswordLength: maData["managed_password_length"].(int),
			}
		}

		// Construct OpenFirmwareEfiPassword field
		var openFirmwareEfiPassword jamfpro.PolicyOpenFirmwareEfiPassword
		if ofData, ok := accountMaintenanceData["open_firmware_efi_password"].(map[string]interface{}); ok {
			openFirmwareEfiPassword = jamfpro.PolicyOpenFirmwareEfiPassword{
				OfMode:     ofData["of_mode"].(string),
				OfPassword: ofData["of_password"].(string),
			}
		}

		// Assign all constructed components to AccountMaintenance
		policy.AccountMaintenance = jamfpro.PolicyAccountMaintenance{
			Accounts:                accounts,
			DirectoryBindings:       directoryBindings,
			ManagementAccount:       managementAccount,
			OpenFirmwareEfiPassword: openFirmwareEfiPassword,
		}
	}

	// Construct the Maintenance section
	if v, ok := d.GetOk("maintenance"); ok {
		maintenanceData := v.([]interface{})[0].(map[string]interface{})
		policy.Maintenance = jamfpro.PolicyMaintenance{
			Recon:                    maintenanceData["recon"].(bool),
			ResetName:                maintenanceData["reset_name"].(bool),
			InstallAllCachedPackages: maintenanceData["install_all_cached_packages"].(bool),
			Heal:                     maintenanceData["heal"].(bool),
			Prebindings:              maintenanceData["prebindings"].(bool),
			Permissions:              maintenanceData["permissions"].(bool),
			Byhost:                   maintenanceData["byhost"].(bool),
			SystemCache:              maintenanceData["system_cache"].(bool),
			UserCache:                maintenanceData["user_cache"].(bool),
			Verify:                   maintenanceData["verify"].(bool),
		}
	}

	// Construct the FilesProcesses section
	if v, ok := d.GetOk("files_processes"); ok {
		filesProcessesData := v.([]interface{})[0].(map[string]interface{})
		policy.FilesProcesses = jamfpro.PolicyFilesProcesses{
			SearchByPath:         filesProcessesData["search_by_path"].(string),
			DeleteFile:           filesProcessesData["delete_file"].(bool),
			LocateFile:           filesProcessesData["locate_file"].(string),
			UpdateLocateDatabase: filesProcessesData["update_locate_database"].(bool),
			SpotlightSearch:      filesProcessesData["spotlight_search"].(string),
			SearchForProcess:     filesProcessesData["search_for_process"].(string),
			KillProcess:          filesProcessesData["kill_process"].(bool),
			RunCommand:           filesProcessesData["run_command"].(string),
		}
	}

	// Construct the UserInteraction section
	if v, ok := d.GetOk("user_interaction"); ok {
		userInteractionData := v.([]interface{})[0].(map[string]interface{})
		policy.UserInteraction = jamfpro.PolicyUserInteraction{
			MessageStart:          userInteractionData["message_start"].(string),
			AllowUserToDefer:      userInteractionData["allow_user_to_defer"].(bool),
			AllowDeferralUntilUtc: userInteractionData["allow_deferral_until_utc"].(string),
			AllowDeferralMinutes:  userInteractionData["allow_deferral_minutes"].(int),
			MessageFinish:         userInteractionData["message_finish"].(string),
		}
	}

	// Construct the DiskEncryption section
	if v, ok := d.GetOk("disk_encryption"); ok {
		diskEncryptionData := v.([]interface{})[0].(map[string]interface{})
		policy.DiskEncryption = jamfpro.PolicyDiskEncryption{
			Action:                                 diskEncryptionData["action"].(string),
			DiskEncryptionConfigurationID:          diskEncryptionData["disk_encryption_configuration_id"].(int),
			AuthRestart:                            diskEncryptionData["auth_restart"].(bool),
			RemediateKeyType:                       diskEncryptionData["remediate_key_type"].(string),
			RemediateDiskEncryptionConfigurationID: diskEncryptionData["remediate_disk_encryption_configuration_id"].(int),
		}
	}

	log.Printf("[INFO] Successfully constructed Jamf Pro Policy with name: %s", policy.General.Name)

	return policy, nil
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
		policyID, err := conn.CreatePolicyByID(policy)
		if err != nil {
			log.Printf("[ERROR] Error creating Jamf Pro Policy with name: %s. Error: %s", policy.General.Name, err)
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			return retry.RetryableError(err)
		}

		// Log the successfully resource creation
		log.Printf("[INFO] Successfully created Jamf Pro Policy with ID: %d", policyID)

		// Set the ID in the Terraform state
		d.SetId(strconv.Itoa(policy.General.ID))

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
		"id":                            policy.General.ID,
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
		"date_time_limitations": []interface{}{map[string]interface{}{
			"activation_date":       policy.General.DateTimeLimitations.ActivationDate,
			"activation_date_epoch": policy.General.DateTimeLimitations.ActivationDateEpoch,
			"activation_date_utc":   policy.General.DateTimeLimitations.ActivationDateUTC,
			"expiration_date":       policy.General.DateTimeLimitations.ExpirationDate,
			"expiration_date_epoch": policy.General.DateTimeLimitations.ExpirationDateEpoch,
			"expiration_date_utc":   policy.General.DateTimeLimitations.ExpirationDateUTC,
			"no_execute_on":         policy.General.DateTimeLimitations.NoExecuteOn,
			"no_execute_start":      policy.General.DateTimeLimitations.NoExecuteStart,
			"no_execute_end":        policy.General.DateTimeLimitations.NoExecuteEnd,
		}},
		"network_limitations": []interface{}{map[string]interface{}{
			"minimum_network_connection": policy.General.NetworkLimitations.MinimumNetworkConnection,
			"any_ip_address":             policy.General.NetworkLimitations.AnyIPAddress,
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
		"all_jss_users": policy.Scope.AllJSSUsers,
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
		"re_install_button_text":          policy.SelfService.ReinstallButtonText,
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
				categories[i] = map[string]interface{}{
					"id":         cat.Category.ID,
					"name":       cat.Category.Name,
					"display_in": cat.Category.DisplayIn,
					"feature_in": cat.Category.FeatureIn,
				}
			}
			return categories
		}(),
	}

	if err := d.Set("self_service", []interface{}{selfServiceAttributes}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// PackageConfiguration section of the policy
	if v, ok := d.GetOk("package_configuration"); ok {
		packageConfigData := v.([]interface{})[0].(map[string]interface{})
		var packageItems []jamfpro.PolicyPackage

		// Iterate over the packages
		if pkgs, ok := packageConfigData["packages"].([]interface{}); ok {
			for _, pkg := range pkgs {
				pkgMap := pkg.(map[string]interface{})
				packageItem := jamfpro.PolicyPackage{
					ID:                pkgMap["id"].(int),
					Name:              pkgMap["name"].(string),
					Action:            pkgMap["action"].(string),
					FillUserTemplate:  pkgMap["fut"].(bool),
					FillExistingUsers: pkgMap["feu"].(bool),
					UpdateAutorun:     pkgMap["update_autorun"].(bool),
				}
				packageItems = append(packageItems, packageItem)
			}
		}

		// Set the PackageConfiguration in the policy
		policy.PackageConfiguration = jamfpro.PolicyPackageConfiguration{
			Packages: packageItems,
		}
	}

	// set scripts attributes
	if v, ok := d.GetOk("scripts"); ok {
		scriptsData := v.([]interface{})[0].(map[string]interface{})
		var scriptItems []jamfpro.PolicyScriptItem

		// Iterate over the scripts
		if scripts, ok := scriptsData["script"].([]interface{}); ok {
			for _, script := range scripts {
				scriptMap := script.(map[string]interface{})
				scriptItem := jamfpro.PolicyScriptItem{
					ID:          scriptMap["id"].(string),
					Name:        scriptMap["name"].(string),
					Priority:    scriptMap["priority"].(string),
					Parameter4:  scriptMap["parameter4"].(string),
					Parameter5:  scriptMap["parameter5"].(string),
					Parameter6:  scriptMap["parameter6"].(string),
					Parameter7:  scriptMap["parameter7"].(string),
					Parameter8:  scriptMap["parameter8"].(string),
					Parameter9:  scriptMap["parameter9"].(string),
					Parameter10: scriptMap["parameter10"].(string),
					Parameter11: scriptMap["parameter11"].(string),
				}
				scriptItems = append(scriptItems, scriptItem)
			}
		}

		// Set the Scripts in the policy
		policy.Scripts = jamfpro.PolicyScripts{
			Size:   scriptsData["size"].(int),
			Script: scriptItems,
		}
	}

	// Set printer attributes
	if v, ok := d.GetOk("printers"); ok {
		printersData := v.([]interface{})[0].(map[string]interface{})
		var printerItems []jamfpro.PolicyPrinterItem

		if printers, ok := printersData["printer"].([]interface{}); ok {
			for _, printer := range printers {
				printerMap := printer.(map[string]interface{})
				printerItem := jamfpro.PolicyPrinterItem{
					ID:          printerMap["id"].(int),
					Name:        printerMap["name"].(string),
					Action:      printerMap["action"].(string),
					MakeDefault: printerMap["make_default"].(bool),
				}
				printerItems = append(printerItems, printerItem)
			}
		}

		policy.Printers = jamfpro.PolicyPrinters{
			Size:                 printersData["size"].(int),
			LeaveExistingDefault: printersData["leave_existing_default"].(bool),
			Printer:              printerItems,
		}
	}

	// Set Dock Items attributes
	if v, ok := d.GetOk("dock_items"); ok {
		dockItemsData := v.([]interface{})[0].(map[string]interface{})
		var dockItems []jamfpro.PolicyDockItem

		if docks, ok := dockItemsData["dock_item"].([]interface{}); ok {
			for _, dock := range docks {
				dockMap := dock.(map[string]interface{})
				dockItem := jamfpro.PolicyDockItem{
					ID:     dockMap["id"].(int),
					Name:   dockMap["name"].(string),
					Action: dockMap["action"].(string),
				}
				dockItems = append(dockItems, dockItem)
			}
		}

		policy.DockItems = jamfpro.PolicyDockItems{
			Size:     dockItemsData["size"].(int),
			DockItem: dockItems,
		}
	}

	// Set Account Maintenance attributes
	if v, ok := d.GetOk("account_maintenance"); ok {
		accountMaintenanceData := v.([]interface{})[0].(map[string]interface{})
		var accounts []jamfpro.PolicyAccount

		if accs, ok := accountMaintenanceData["accounts"].([]interface{}); ok {
			for _, acc := range accs {
				accMap := acc.(map[string]interface{})
				account := jamfpro.PolicyAccount{
					Action:                 accMap["action"].(string),
					Username:               accMap["username"].(string),
					Realname:               accMap["realname"].(string),
					Password:               accMap["password"].(string),
					ArchiveHomeDirectory:   accMap["archive_home_directory"].(bool),
					ArchiveHomeDirectoryTo: accMap["archive_home_directory_to"].(string),
					Home:                   accMap["home"].(string),
					Picture:                accMap["picture"].(string),
					Admin:                  accMap["admin"].(bool),
					FilevaultEnabled:       accMap["filevault_enabled"].(bool),
				}
				accounts = append(accounts, account)
			}
		}

		policy.AccountMaintenance = jamfpro.PolicyAccountMaintenance{
			Accounts: accounts,
		}
	}

	// Set Directory Bindings attributes
	if v, ok := d.GetOk("directory_bindings"); ok {
		directoryBindingsData := v.([]interface{})[0].(map[string]interface{})
		var directoryBindings []jamfpro.PolicyDirectoryBinding

		if bindings, ok := directoryBindingsData["binding"].([]interface{}); ok {
			for _, binding := range bindings {
				bindingMap := binding.(map[string]interface{})
				directoryBinding := jamfpro.PolicyDirectoryBinding{
					ID:   bindingMap["id"].(int),
					Name: bindingMap["name"].(string),
				}
				directoryBindings = append(directoryBindings, directoryBinding)
			}
		}

		policy.AccountMaintenance.DirectoryBindings = directoryBindings
	}

	// Set Management Account attributes
	if v, ok := d.GetOk("management_account"); ok {
		managementAccountData := v.([]interface{})[0].(map[string]interface{})
		var managementAccount jamfpro.PolicyManagementAccount

		managementAccount = jamfpro.PolicyManagementAccount{
			Action:                managementAccountData["action"].(string),
			ManagedPassword:       managementAccountData["managed_password"].(string),
			ManagedPasswordLength: managementAccountData["managed_password_length"].(int),
		}

		policy.AccountMaintenance.ManagementAccount = managementAccount
	}

	// Set Open Firmware/EFI Password attributes
	if v, ok := d.GetOk("open_firmware_efi_password"); ok {
		openFirmwareEfiPasswordData := v.([]interface{})[0].(map[string]interface{})
		var openFirmwareEfiPassword jamfpro.PolicyOpenFirmwareEfiPassword

		openFirmwareEfiPassword = jamfpro.PolicyOpenFirmwareEfiPassword{
			OfMode:     openFirmwareEfiPasswordData["of_mode"].(string),
			OfPassword: openFirmwareEfiPasswordData["of_password"].(string),
		}

		policy.AccountMaintenance.OpenFirmwareEfiPassword = openFirmwareEfiPassword
	}

	// Set Maintenance attributes
	if v, ok := d.GetOk("maintenance"); ok {
		maintenanceData := v.([]interface{})[0].(map[string]interface{})
		policy.Maintenance = jamfpro.PolicyMaintenance{
			Recon:                    maintenanceData["recon"].(bool),
			ResetName:                maintenanceData["reset_name"].(bool),
			InstallAllCachedPackages: maintenanceData["install_all_cached_packages"].(bool),
			Heal:                     maintenanceData["heal"].(bool),
			Prebindings:              maintenanceData["prebindings"].(bool),
			Permissions:              maintenanceData["permissions"].(bool),
			Byhost:                   maintenanceData["byhost"].(bool),
			SystemCache:              maintenanceData["system_cache"].(bool),
			UserCache:                maintenanceData["user_cache"].(bool),
			Verify:                   maintenanceData["verify"].(bool),
		}
	}

	// set Files and Processes attributes
	if v, ok := d.GetOk("files_processes"); ok {
		filesProcessesData := v.([]interface{})[0].(map[string]interface{})
		policy.FilesProcesses = jamfpro.PolicyFilesProcesses{
			SearchByPath:         filesProcessesData["search_by_path"].(string),
			DeleteFile:           filesProcessesData["delete_file"].(bool),
			LocateFile:           filesProcessesData["locate_file"].(string),
			UpdateLocateDatabase: filesProcessesData["update_locate_database"].(bool),
			SpotlightSearch:      filesProcessesData["spotlight_search"].(string),
			SearchForProcess:     filesProcessesData["search_for_process"].(string),
			KillProcess:          filesProcessesData["kill_process"].(bool),
			RunCommand:           filesProcessesData["run_command"].(string),
		}
	}

	// Set User Interaction attributes
	if v, ok := d.GetOk("user_interaction"); ok {
		userInteractionData := v.([]interface{})[0].(map[string]interface{})
		policy.UserInteraction = jamfpro.PolicyUserInteraction{
			MessageStart:          userInteractionData["message_start"].(string),
			AllowUserToDefer:      userInteractionData["allow_user_to_defer"].(bool),
			AllowDeferralUntilUtc: userInteractionData["allow_deferral_until_utc"].(string),
			AllowDeferralMinutes:  userInteractionData["allow_deferral_minutes"].(int),
			MessageFinish:         userInteractionData["message_finish"].(string),
		}
	}

	// Set Disk Encryption attributes
	if v, ok := d.GetOk("disk_encryption"); ok {
		diskEncryptionData := v.([]interface{})[0].(map[string]interface{})
		policy.DiskEncryption = jamfpro.PolicyDiskEncryption{
			Action:                                 diskEncryptionData["action"].(string),
			DiskEncryptionConfigurationID:          diskEncryptionData["disk_encryption_configuration_id"].(int),
			AuthRestart:                            diskEncryptionData["auth_restart"].(bool),
			RemediateKeyType:                       diskEncryptionData["remediate_key_type"].(string),
			RemediateDiskEncryptionConfigurationID: diskEncryptionData["remediate_disk_encryption_configuration_id"].(int),
		}
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
