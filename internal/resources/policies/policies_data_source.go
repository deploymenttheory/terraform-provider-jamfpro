// policies_data_source.go
package policies

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProPolicies provides information about specific Jamf Pro policies by their ID or Name.
func DataSourceJamfProPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJamfProPoliciesRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"general": {
				Type:        schema.TypeList,
				Computed:    true,
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
							Computed:    true,
							Description: "The name of the policy.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Define whether the policy is enabled.",
						},
						"trigger": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event(s) triggers to use to initiate the policy. Values can be 'USER_INITIATED' for self self trigger and 'EVENT' for an event based trigger",
						},
						"trigger_checkin": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Trigger policy when device performs recurring check-in against the frequency configured in Jamf Pro",
						},
						"trigger_enrollment_complete": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Trigger policy when device enrollment is complete.",
						},
						"trigger_login": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Trigger policy when a user logs in to a computer. A login event that checks for policies must be configured in Jamf Pro for this to work",
						},
						"trigger_logout": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Trigger policy when a user logout.",
						},
						"trigger_network_state_changed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Trigger policy when it's network state changes. When a computer's network state changes (e.g., when the network connection changes, when the computer name changes, when the IP address changes)",
						},
						"trigger_startup": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Trigger policy when a computer starts up. A startup script that checks for policies must be configured in Jamf Pro for this to work",
						},
						"trigger_other": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Any other trigger for the policy.",
						},
						"frequency": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Frequency of policy execution.",
						},
						"retry_event": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event on which to retry policy execution.",
						},
						"retry_attempts": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of retry attempts for the jamf pro policy. Valid values are -1 (not configured) and 1 through 10.",
						},
						"notify_on_each_failed_retry": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Send notifications for each failed policy retry attempt. ",
						},
						"location_user_only": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Location-based policy for user only.",
						},
						"target_drive": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The drive on which to run the policy (e.g. /Volumes/Restore/ ). The policy runs on the boot drive by default",
						},
						"offline": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Make policy available offline by caching the policy to the macOS device to ensure it runs when Jamf Pro is unavailable. Only used when execution policy is set to 'ongoing'. ",
						},
						"category": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Category to add the policy to.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The category ID assigned to the jamf pro policy. Defaults to '-1' aka not used.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Category Name for assigned jamf pro policy. Value defaults to 'No category assigned' aka not used",
									},
								},
							},
						},
						"date_time_limitations": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Server-side limitations use your Jamf Pro host server's time zone and settings. The Jamf Pro host service is in UTC time.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"activation_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The activation date of the policy.",
									},
									"activation_date_epoch": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The epoch time of the activation date.",
									},
									"activation_date_utc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The UTC time of the activation date.",
									},
									"expiration_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The expiration date of the policy.",
									},
									"expiration_date_epoch": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The epoch time of the expiration date.",
									},
									"expiration_date_utc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The UTC time of the expiration date.",
									},
									"no_execute_on": {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Client-side limitations are enforced based on the settings on computers. This field sets specific days when the policy should not execute.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"no_execute_start": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Client-side limitations are enforced based on the settings on computers. This field sets the start time when the policy should not execute.",
									},
									"no_execute_end": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Client-side limitations are enforced based on the settings on computers. This field sets the end time when the policy should not execute.",
									},
								},
							},
						},
						"network_limitations": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Network limitations for the policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"minimum_network_connection": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Minimum network connection required for the policy.",
									},
									"any_ip_address": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the policy applies to any IP address.",
									},
									"network_segments": {
										Type:        schema.TypeString,
										Description: "Network segment limitations for the policy.",
										Computed:    true,
									},
								},
							},
						},
						"override_default_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Settings to override default configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_drive": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The drive on which to run the policy (e.g. '/Volumes/Restore/'). Defaults to '/' if no value is defined, which is the root of the file system.",
									},
									"distribution_point": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Distribution point for the policy.",
									},
									"force_afp_smb": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to force AFP/SMB.",
									},
									"sus": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Software Update Service for the policy.",
									},
								},
							},
						},
						"network_requirements": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network requirements for the policy.",
						},
						"site": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Jamf Pro Site-related settings of the policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Jamf Pro Site ID. Value defaults to -1 aka not used.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
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
				Computed:    true,
				Description: "Scope configuration for the profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_computers": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "scope all_computers if true, applies the profile to all computers. If false applies to specific computers.",
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
													Computed:    true,
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
													Computed:    true,
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
													Computed:    true,
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
													Computed:    true,
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
													Computed:    true,
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
										Computed:    true,
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
							Computed:    true,
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
																Computed:    true,
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
																Computed:    true,
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
																Computed:    true,
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
																Computed:    true,
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
							Computed:    true,
							Description: "Scoped exclusions to exclude from the policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"computers": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Scoped computer exclusions to exclude from the policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"computer": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The individual computer to exclude from the policy.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Computed:    true,
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
																Computed:    true,
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
																Computed:    true,
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
																Computed:    true,
																Description: "The unique identifier of the JSS user group.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
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
																Computed:    true,
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
																Computed:    true,
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
																Computed:    true,
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
																Computed:    true,
																Description: "The unique identifier of the user.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
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
																Computed:    true,
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
																Computed:    true,
																Description: "The unique identifier of the iBeacon.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
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
				Computed:    true,
				Description: "Self-service settings of the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_for_self_service": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the policy is available for self-service.",
						},
						"self_service_display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name of the policy in self-service.",
						},
						"install_button_text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Text displayed on the install button in self-service.",
						},
						"reinstall_button_text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Text displayed on the re-install button in self-service.",
						},
						"self_service_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the policy displayed in self-service.",
						},
						"force_users_to_view_description": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to force users to view the policy description in self-service.",
						},
						"self_service_icon": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Icon settings for the policy in self-service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "ID of the icon used in self-service.",
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
							Computed:    true,
							Description: "Whether to feature the policy on the main page of self-service.",
						},
						"self_service_categories": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Category settings for the policy in self-service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Category details for the policy in self-service.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Category ID for the policy in self-service.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Category name for the policy in self-service.",
												},
												"display_in": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to display the category in self-service.",
												},
												"feature_in": {
													Type:        schema.TypeBool,
													Computed:    true,
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
				Computed:    true,
				Description: "Package configuration settings of the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"packages": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of packages included in the policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"package": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Details of the package.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Unique identifier of the package.",
												},
												"name": {
													Type:        schema.TypeString,
													Description: "Name of the package.",
													Computed:    true,
												},
												"action": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Action to be performed for the package.",
												},
												"fut": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Fill User Template (FUT).",
												},
												"feu": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Fill Existing Users (FEU).",
												},
												"update_autorun": {
													Type:        schema.TypeBool,
													Computed:    true,
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
				Computed:    true,
				Description: "Scripts settings of the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"script": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Details of the scripts.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Unique identifier of the script.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the script.",
									},
									"priority": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Execution priority of the script.",
									},
									"parameter4": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom parameter 4 for the script.",
									},
									"parameter5": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom parameter 5 for the script.",
									},
									"parameter6": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom parameter 6 for the script.",
									},
									"parameter7": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom parameter 7 for the script.",
									},
									"parameter8": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom parameter 8 for the script.",
									},
									"parameter9": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom parameter 9 for the script.",
									},
									"parameter10": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Custom parameter 10 for the script.",
									},
									"parameter11": {
										Type:        schema.TypeString,
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
				Computed:    true,
				Description: "Printers settings of the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"leave_existing_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Policy for handling existing default printers.",
						},
						"printer": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Details of the printer configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Unique identifier of the printer.",
									},
									"name": {
										Type:        schema.TypeString,
										Description: "Name of the printer.",
										Computed:    true,
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action to be performed for the printer (e.g., install, uninstall).",
									},
									"make_default": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to set the printer as the default.",
									},
								},
							},
						},
					},
				},
			},
			"policy_items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Dock items settings of the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_item": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Details of the policy item configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Unique identifier of the policy item.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the policy item.",
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action to be performed for the policy item (e.g., Add To Beginning, Add To End, Remove).",
									},
								},
							},
						},
					},
				},
			},
			"account_maintenance": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Account maintenance settings of the policy. Use this section to create and delete local accounts, and to reset local account passwords. Also use this section to disable an existing local account for FileVault 2.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accounts": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of account maintenance configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Details of each account configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"action": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Action to be performed on the account (e.g., Create, Reset, Delete, DisableFileVault).",
												},
												"username": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Username/short name for the account",
												},
												"realname": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Real name associated with the account.",
												},
												"password": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Set a new account password. This does not update the account's login keychain password or FileVault 2 password.",
												},
												"archive_home_directory": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Permanently delete home directory. If set to true will archive the home directory.",
												},
												"archive_home_directory_to": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Path in which to archive the home directory to.",
												},
												"home": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Full path in which to create the home directory (e.g. /Users/username/ or /private/var/username/)",
												},
												"hint": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Hint to help the user remember the password",
												},
												"picture": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Full path to the account picture (e.g. /Library/User Pictures/Animals/Butterfly.tif )",
												},
												"admin": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the account has admin privileges.Setting this to true will set the user administrator privileges to the computer",
												},
												"filevault_enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
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
							Computed:    true,
							Description: "Directory binding settings for the policy. Use this section to bind computers to a directory service",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"binding": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Details of the directory binding.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
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
							Computed:    true,
							Description: "Management account settings for the policy. Use this section to change or reset the management account password.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action to perform on the management account.Rotates management account password at next policy execution. Valid values are 'rotate' or 'doNotChange'.",
									},
									"managed_password": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Managed password for the account. Management account passwords will be automatically randomized with 29 characters by jamf pro.",
									},
									"managed_password_length": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Length of the managed password. Only necessary when utilizing the random action",
									},
								},
							},
						},
						"open_firmware_efi_password": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Open Firmware/EFI password settings for the policy. Use this section to set or remove an Open Firmware/EFI password on computers with Intel-based processors.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"of_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mode for the open firmware/EFI password. Valid values are 'command' or 'none'.",
									},
									"of_password": {
										Type:        schema.TypeString,
										Computed:    true,
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
				Computed:    true,
				Description: "Use this section to restart computers and specify the disk to boot them to",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reboot message displayed to the user.",
						},
						"specify_startup": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reboot Method",
						},
						"startup_disk": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk to boot computers to",
						},
						"no_user_logged_in": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action to take if no user is logged in to the computer",
						},
						"user_logged_in": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action to take if a user is logged in to the computer",
						},
						"minutes_until_reboot": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Amount of time to wait before the restart begins.",
						},
						"start_reboot_timer_immediately": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Defines if the reboot timer should start immediately once the policy applies to a macOS device.",
						},
						"file_vault_2_reboot": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Perform authenticated restart on computers with FileVault 2 enabled. Restart FileVault 2-encrypted computers without requiring an unlock during the next startup",
						},
					},
				},
			},
			"maintenance": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Maintenance settings of the policy. Use this section to update inventory, reset computer names, install all cached packages, and run common maintenance tasks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recon": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to run recon (inventory update) as part of the maintenance. Forces computers to submit updated inventory information to Jamf Pro",
						},
						"reset_name": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to reset the computer name to the name stored in Jamf Pro. Changes the computer name on computers to match the computer name in Jamf Pro",
						},
						"install_all_cached_packages": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to install all cached packages. Installs packages cached by Jamf Pro",
						},
						"heal": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to heal the policy.",
						},
						"prebindings": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to update prebindings.",
						},
						"permissions": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to fix Disk Permissions (Not compatible with macOS v10.12 or later)",
						},
						"byhost": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to fix ByHost files andnpreferences.",
						},
						"system_cache": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to flush caches from /Library/Caches/ and /System/Library/Caches/, except for any com.apple.LaunchServices caches",
						},
						"user_cache": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to flush caches from ~/Library/Caches/, ~/.jpi_cache/, and ~/Library/Preferences/Microsoft/Office version #/Office Font Cache. Enabling this may cause problems with system fonts displaying unless a restart option is configured.",
						},
						"verify": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to verify system files and structure on the Startup Disk",
						},
					},
				},
			},
			"files_processes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Files and processes settings of the policy. Use this section to search for and log specific files and processes. Also use this section to execute a command.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search_by_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path of the file to search for.",
						},
						"delete_file": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to delete the file found at the specified path.",
						},
						"locate_file": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path of the file to locate. Name of the file, including the file extension. This field is case-sensitive and returns partial matches",
						},
						"update_locate_database": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to update the locate database. Update the locate database before searching for the file",
						},
						"spotlight_search": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Search For File Using Spotlight. File to search for. This field is not case-sensitive and returns partial matches",
						},
						"search_for_process": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the process to search for. This field is case-sensitive and returns partial matches",
						},
						"kill_process": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to kill the process if found. This works with exact matches only",
						},
						"run_command": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Command to execute on computers. This command is executed as the 'root' user",
						},
					},
				},
			},
			"user_interaction": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User interaction settings of the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message_start": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message to display before the policy runs",
						},
						"allow_user_to_defer": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Allow user deferral and configure deferral type. A deferral limit must be specified for this to work.",
						},
						"allow_deferral_until_utc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date/time at which deferrals are prohibited and the policy runs. Uses time zone settings of your hosting server. Standard environments hosted in Jamf Cloud use Coordinated Universal Time (UTC)",
						},
						"allow_deferral_minutes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of minutes after the user was first prompted by the policy at which the policy runs and deferrals are prohibited",
						},
						"message_finish": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message to display when the policy is complete.",
						},
					},
				},
			},
			"disk_encryption": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Disk encryption settings of the policy. Use this section to enable FileVault 2 or to issue a new recovery key.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The action to perform for disk encryption (e.g., apply, remediate).",
						},
						"disk_encryption_configuration_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the disk encryption configuration to apply.",
						},
						"auth_restart": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to allow authentication restart.",
						},
						"remediate_key_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of key to use for remediation (e.g., Individual, Institutional, Individual And Institutional).",
						},
						"remediate_disk_encryption_configuration_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk encryption ID to utilize for remediating institutional recovery key types.",
						},
					},
				},
			},
		},
	}
}

// dataSourceJamfProPoliciesRead fetches the details of specific policy items from Jamf Pro using either their unique Name or Id.
func dataSourceJamfProPoliciesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			// If fetching by ID fails, try fetching by Name from the 'general' section
			generalSettings, ok := d.GetOk("general")
			if !ok || len(generalSettings.([]interface{})) == 0 {
				return retry.NonRetryableError(fmt.Errorf("unable to find 'general' block for terraform read operation"))
			}
			generalMap := generalSettings.([]interface{})[0].(map[string]interface{})
			policyName, ok := generalMap["name"].(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' within 'general' as a string for terraform read operation"))
			}

			policy, apiErr = conn.GetPolicyByName(policyName)
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
