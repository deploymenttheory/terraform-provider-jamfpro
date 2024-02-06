package macosconfigurationprofiles

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const JamfProResourceMacOSConfigurationProfile = "macos_configuration_profile"

func ResourceJamfProMacOSConfigurationProfiles() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProMacOSConfigurationProfilesCreate,
		ReadContext:   ResourceJamfProMacOSConfigurationProfilesRead,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the macOS configuration profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the configuration profile.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the configuration profile.",
			},
			"site": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Jamf Pro Site information for the assigned configuration profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,

							Description: "Jamf Pro Site ID. Value defaults to -1 aka not used.",
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,

							Description: "Jamf Pro Site Name. Value defaults to 'None' aka not used",
						},
					},
				},
			},
			"category": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Configuration Profile Category information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,

							Description: "Category ID. Value defaults to -1 aka not used.",
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,

							Description: "Category Name for assigned configuration profile. Value defaults to 'No category assigned' aka not used",
						},
					},
				},
			},
			"distribution_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The distribution methods for the macOS configuration profile, can be either Install Automatically or Make Available in Self Service.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					value := val.(string)
					allowedValues := []string{"Install Automatically", "Make Available in Self Service"}
					for _, v := range allowedValues {
						if value == v {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of [%s], got: %s", key, strings.Join(allowedValues, ", "), value))
					return
				},
			},
			"user_removable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Define whether the macOS configuration profile is removeable by the end user using jamf self service.",
			},
			"level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The level defines what level the MDM profile is deployed at. It can either be device wide using 'System'', or for an individual user using 'User'.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					value := val.(string)
					allowedValues := []string{"System", "User"}
					for _, v := range allowedValues {
						if value == v {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be either 'System' or 'User', got: %s", key, value))
					return
				},
			},
			"uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The uuid of the macos configuration profile",
			},
			"redeploy_on_update": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Newly Assigned",
				Description: "",
			},
			"payloads": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The configuration profile payload in xml and delivered as a plist to the macOS device by Jamf Pro.",
				DiffSuppressFunc: suppressPayloadDiff,
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
								},
							},
						},
					},
				},
			},

			"self_service": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Self Service configuration for the profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"install_button_text": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"self_service_description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"force_users_to_view_description": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"self_service_icon": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"uri": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"data": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"feature_on_main_page": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"self_service_categories": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of categories under self service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The category information for the self service.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "The unique identifier of the category.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the category.",
												},
												"display_in": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to display in self-service.",
												},
												"feature_in": {
													Type:        schema.TypeBool,
													Optional:    true,
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
							Optional: true,
						},
						"notification_subject": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"notification_message": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func constructJamfProMacOSConfigurationProfile(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceMacOSConfigurationProfile, error) {

	out := jamfpro.ResourceMacOSConfigurationProfile{
		General: jamfpro.MacOSConfigurationProfileSubsetGeneral{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Site: jamfpro.SharedResourceSite{
				ID:   d.Get("site.0.id").(int),
				Name: d.Get("site.0.name").(string),
			},
			Category: jamfpro.SharedResourceCategory{
				ID:   d.Get("category.0.id").(int),
				Name: d.Get("category.0.name").(string),
			},
			DistributionMethod: d.Get("distribution_method").(string),
			UserRemovable:      d.Get("user_removable").(bool),
			Level:              d.Get("level").(string),
			UUID:               d.Get("uuid").(string),
			RedeployOnUpdate:   d.Get("redeploy_on_update").(string),
		},
		Scope: jamfpro.MacOSConfigurationProfileSubsetScope{
			AllComputers: d.Get("scope.0.all_computers").(bool),
			AllJSSUsers:  d.Get("scope.0.all_jss_users").(bool),
		},
		SelfService: jamfpro.MacOSConfigurationProfileSubsetSelfService{
			InstallButtonText:           d.Get("self_service.0.install_button_text").(string),
			SelfServiceDescription:      d.Get("self_service.0.self_service_description").(string),
			ForceUsersToViewDescription: d.Get("self_service.0.force_users_to_view_description").(bool),
			SelfServiceIcon: jamfpro.SharedResourceSelfServiceIcon{
				ID:  d.Get("self_service.0.self_service_icon.0.id").(int),
				URI: d.Get("self_service.0.self_service_icon.0.uri").(string),
			},
			FeatureOnMainPage: d.Get("self_service.0.feature_on_main_page").(bool),
			SelfServiceCategories: jamfpro.MacOSConfigurationProfileSubsetSelfServiceCategories{
				Category: jamfpro.MacOSConfigurationProfileSubsetSelfServiceCategory{
					ID:        d.Get("self_service.0.self_service_categories.0.category.0.id").(int),
					Name:      d.Get("self_service.0.self_service_categories.0.category.0.name").(string),
					DisplayIn: d.Get("self_service.0.self_service_categories.0.category.0.display_in").(bool),
					FeatureIn: d.Get("self_service.0.self_service_categories.0.category.0.feature_in").(bool),
				},
			},
			Notification:        d.Get("self_service.0.notification").(string),
			NotificationSubject: d.Get("self_service.0.notification_subject").(string),
			NotificationMessage: d.Get("self_service.0.notification_message").(string),
		},
	}

	return &out, nil
}

func ResourceJamfProMacOSConfigurationProfilesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics
	var creationResponse *jamfpro.ResponseMacOSConfigurationProfileCreationUpdate
	var apiErrorCode int
	resourceName := d.Get("name").(string)

	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemCreate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	out, err := constructJamfProMacOSConfigurationProfile(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceMacOSConfigurationProfile, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceMacOSConfigurationProfile)

	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateMacOSConfigurationProfile(out)
		if apiErr != nil {

			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPICreateFailedAfterRetry(subCtx, JamfProResourceMacOSConfigurationProfile, resourceName, apiErr.Error(), apiErrorCode)

			return retry.NonRetryableError(apiErr)
		}

		return nil
	})

	if err != nil {

		logging.LogAPICreateFailure(subCtx, JamfProResourceMacOSConfigurationProfile, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	idString := strconv.Itoa(creationResponse.ID)
	logging.LogAPICreateSuccess(subCtx, JamfProResourceMacOSConfigurationProfile, idString)
	d.SetId(idString)

	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProMacOSConfigurationProfilesRead(subCtx, d, meta)
		if len(readDiags) > 0 {

			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceMacOSConfigurationProfile, d.Id(), readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}

		return nil
	})

	if err != nil {

		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceMacOSConfigurationProfile, err.Error())
		diags = append(diags, diag.FromErr(err)...)
	} else {

		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceMacOSConfigurationProfile, d.Id())
	}

	return diags
}

func ResourceJamfProMacOSConfigurationProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)

	var diags diag.Diagnostics
	var apiErrorCode int
	var resp *jamfpro.ResourceMacOSConfigurationProfile
	resourceID := d.Id()
	resourceIDString, convErr := strconv.Atoi(resourceID)
	if convErr != nil {
		return diag.FromErr(convErr)

	}

	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resp, apiErr = conn.GetMacOSConfigurationProfileByID(resourceIDString)
		if apiErr != nil {
			logging.LogFailedReadByID(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, apiErr.Error(), apiErrorCode)

			return retry.RetryableError(apiErr)
		}

		return nil
	})

	if err != nil {

		logging.LogTFStateRemovalWarning(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID)
		return diag.FromErr(err)
	}

	logging.LogAPIReadSuccess(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID)

	if err := d.Set("id", resourceID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("name", resp.General.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

func ResourceJamfProMacOSConfigurationProfilesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemUpdate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceIDString, convErr := strconv.Atoi(resourceID)
	if convErr != nil {
		return diag.FromErr(convErr)

	}
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	constructedPayload, err := constructJamfProMacOSConfigurationProfile(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceMacOSConfigurationProfile, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceMacOSConfigurationProfile)

	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateMacOSConfigurationProfileByID(resourceIDString, constructedPayload)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.LogAPIUpdateFailureByID(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			_, apiErrByName := conn.UpdateMacOSConfigurationProfileByName(resourceName, constructedPayload)
			if apiErrByName != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErrByName.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIUpdateFailureByName(subCtx, JamfProResourceMacOSConfigurationProfile, resourceName, apiErrByName.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErrByName)
			}
		} else {
			logging.LogAPIUpdateSuccess(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, resourceName)
		}
		return nil
	})

	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProMacOSConfigurationProfilesRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceMacOSConfigurationProfile, resourceID, readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceMacOSConfigurationProfile, err.Error())
		return diag.FromErr(err)
	} else {
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceMacOSConfigurationProfile, resourceID)
	}

	return nil
}

func ResourceJamfProMacOSConfigurationProfilesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemDelete, hclog.Info)

	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {

		apiErr := conn.DeleteDepartmentByID(resourceID)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPIDeleteFailureByID(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			apiErr = conn.DeleteDepartmentByName(resourceName)
			if apiErr != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErr.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIDeleteFailureByName(subCtx, JamfProResourceMacOSConfigurationProfile, resourceName, apiErr.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	logging.LogAPIDeleteSuccess(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, resourceName)

	d.SetId("")

	return nil
}
