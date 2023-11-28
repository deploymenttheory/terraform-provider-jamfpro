// macosconfigurationprofiles_data_source.go
package macosconfigurationprofiles

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProMacOSConfigurationProfiles provides information about a specific macOS Configuration Profile by its ID or Name.
func DataSourceJamfProMacOSConfigurationProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJamfProMacOSConfigurationProfilesRead,
		Schema: map[string]*schema.Schema{
			// GeneralConfig fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the macOS configuration profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the configuration profile.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the configuration profile.",
			},
			"site": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true, // Allow users to optionally filter by 'site'
				Description: "The jamf pro site the configuration profile is assigned to.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true, // User can filter by 'id'
							Computed: true, // 'id' will also be populated based on fetched data
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true, // User can filter by 'name'
							Computed: true, // 'name' will also be populated based on fetched data
						},
					},
				},
			},
			"category": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true, // Allow users to optionally filter by 'site'
				Description: "The jamf pro category the configuration profile is assigned to.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true, // User can filter by 'id'
							Computed: true, // 'id' will also be populated based on fetched data
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true, // User can filter by 'name'
							Computed: true, // 'name' will also be populated based on fetched data
						},
					},
				},
			},
			"distribution_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The distribution methods for the macOS configuration profile, can be either Install Automatically or Make Available in Self Service.",
			},
			"user_removable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Define whether the macOS configuration profile is removeable by the end user.",
			},
			"level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The level defines what level the MDM profile is deployed at. It can either be device wide using computer, or for an individual user.",
			},
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"redeploy_on_update": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"payloads": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// ScopeConfig fields
			"scope": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Description: "Scope configuration for the profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_computers": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If true, applies the profile to all computers.",
						},
						"all_jss_users": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If true, applies the profile to all JSS users.",
						},
						"computers": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"computer": {
										Type: schema.TypeList,
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
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"building": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The unique identifier of the scoped building.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the scoped building.",
												},
											},
										},
									},
								},
							},
						},
						"departments": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"department": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The unique identifier of the scoped department.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the scoped department.",
												},
											},
										},
									},
								},
							},
						},
						"computer_groups": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"computer_group": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The unique identifier of the scoped computer group.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the computer scoped group.",
												},
											},
										},
									},
								},
							},
						},
						"jss_users": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"jss_user": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The unique identifier of the scoped JSS user.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the scoped JSS user.",
												},
											},
										},
									},
								},
							},
						},
						"jss_user_groups": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"jss_user_group": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The unique identifier of the scoped JSS user group.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the scoped JSS user group.",
												},
											},
										},
									},
								},
							},
						},
						"limitations": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"users": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user": {
													Type: schema.TypeList,
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
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user_group": {
													Type: schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The unique identifier of the scoped user group.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the scoped user group.",
															},
														},
													},
												},
											},
										},
									},
									"network_segments": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"network_segment": {
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The unique identifier of the scoped network segment.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the scoped network segment.",
															},
															"uid": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "UID of the scoped network segment.",
															},
														},
													},
												},
											},
										},
									},
									"ibeacons": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ibeacon": {
													Type: schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The unique identifier of the scoped iBeacon.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
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
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"computers": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"computer": {
													Type: schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The unique identifier of the computer.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the computer.",
															},
															"udid": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "UDID of the computer.",
															},
														},
													},
												},
											},
										},
									},
									"buildings": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"building": {
													Type: schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The unique identifier of the building.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the building.",
															},
														},
													},
												},
											},
										},
									},
									"departments": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"department": {
													Type: schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The unique identifier of the department.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the department.",
															},
														},
													},
												},
											},
										},
									},
									"computer_groups": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"computer_group": {
													Type: schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The unique identifier of the computer group.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the computer group.",
															},
														},
													},
												},
											},
										},
									},
									"users": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user": {
													Type: schema.TypeList,
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
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user_group": {
													Type: schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The unique identifier of the user group.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the user group.",
															},
														},
													},
												},
											},
										},
									},
									"network_segments": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"network_segment": {
													Type: schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The unique identifier of the network segment.",
															},
															"uid": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "UID of the network segment.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the network segment.",
															},
														},
													},
												},
											},
										},
									},
									"ibeacons": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ibeacon": {
													Type: schema.TypeList,
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
									"jss_users": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"jss_user": {
													Type: schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The unique identifier of the JSS user.",
															},
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the JSS user.",
															},
														},
													},
												},
											},
										},
									},
									"jss_user_groups": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"jss_user_group": {
													Type: schema.TypeList,
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
								},
							},
						},
					},
				},
			},

			// SelfServiceConfig fields
			"self_service": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Description: "Self Service configuration for the profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"install_button_text": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"self_service_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"force_users_to_view_description": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"self_service_icon": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"uri": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"data": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"feature_on_main_page": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"self_service_categories": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of categories under self service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Computed:    true,
										Description: "The category information for the self service.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The unique identifier of the category.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the category.",
												},
												"display_in": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to display in self-service.",
												},
												"feature_in": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to feature in self-service.",
												},
											},
										},
									},
								},
							},
						},
						"notification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_subject": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// dataSourceJamfProMacOSConfigurationProfilesRead fetches the details of a specific macOS Configuration Profile
// from Jamf Pro using either its unique Name or its ID. The function prioritizes the 'name' attribute over the 'id'
// attribute for fetching details. If neither 'name' nor 'id' is provided, it returns an error.
// Once the details are fetched, they are set in the data source's state.
func dataSourceJamfProMacOSConfigurationProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var profile *jamfpro.ResponseMacOSConfigurationProfiles
	var err error

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		profileName := v.(string)
		profile, err = conn.GetMacOSConfigurationProfileByName(profileName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch macOS Configuration Profile by name: %v", err))
		}
	} else if v, ok := d.GetOk("id"); ok {
		profileID := v.(int)
		profile, err = conn.GetMacOSConfigurationProfileByID(profileID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch macOS Configuration Profile by ID: %v", err))
		}
	} else {
		return diag.Errorf("Either 'name' or 'id' must be provided")
	}

	// Set the data source attributes using the fetched data
	d.SetId(fmt.Sprintf("%d", profile.General.ID))
	d.Set("name", profile.General.Name)
	// Set other necessary attributes

	return nil
}
