// macosconfigurationprofiles_resource.go
package macosconfigurationprofiles

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProMacOSConfigurationProfiles defines the schema and CRUD operations for managing Jamf Pro Departments in Terraform.
func ResourceJamfProMacOSConfigurationProfiles() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProMacOSConfigurationProfilesCreate,
		ReadContext:   ResourceJamfProMacOSConfigurationProfilesRead,
		UpdateContext: ResourceJamfProMacOSConfigurationProfilesUpdate,
		DeleteContext: ResourceJamfProMacOSConfigurationProfilesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// GeneralConfig fields
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
				Description: "Site information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"category": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Category information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
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
				Description: "Define whether the macOS configuration profile is removeable by the end user.",
			},
			"level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The level defines what level the MDM profile is deployed at. It can either be device wide using computer, or for an individual user.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					value := val.(string)
					allowedValues := []string{"computer", "user"}
					for _, v := range allowedValues {
						if value == v {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be either 'computer' or 'user', got: %s", key, value))
					return
				},
			},
			"uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"redeploy_on_update": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payloads": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// ScopeConfig fields
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

			// SelfServiceConfig fields
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

// constructJamfProMacOSConfigurationProfile constructs a ResponseMacOSConfigurationProfile object from the provided schema data.
// It captures each attribute from the schema and returns the constructed ResponseMacOSConfigurationProfile object.
func constructJamfProMacOSConfigurationProfile(d *schema.ResourceData) *jamfpro.ResponseMacOSConfigurationProfiles {
	// Construct General section
	general := jamfpro.MacOSConfigurationProfilesDataSubsetGeneral{
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		DistributionMethod: d.Get("distribution_method").(string),
		UserRemovable:      d.Get("user_removable").(bool),
		Level:              d.Get("level").(string),
		UUID:               d.Get("uuid").(string),
		RedeployOnUpdate:   d.Get("redeploy_on_update").(string),
		Payloads:           d.Get("payloads").(string),
	}

	if siteList, ok := d.GetOk("site"); ok {
		site := siteList.([]interface{})[0].(map[string]interface{})
		general.Site = jamfpro.MacOSConfigurationProfilesDataSubsetSite{
			ID:   site["id"].(int),
			Name: site["name"].(string),
		}
	}

	if categoryList, ok := d.GetOk("category"); ok {
		category := categoryList.([]interface{})[0].(map[string]interface{})
		general.Category = jamfpro.MacOSConfigurationProfilesDataSubsetCategory{
			ID:   category["id"].(int),
			Name: category["name"].(string),
		}
	}

	// Construct Scope section
	scopeData := d.Get("scope").([]interface{})
	var scope jamfpro.MacOSConfigurationProfilesDataSubsetScope
	if len(scopeData) > 0 {
		scopeMap := scopeData[0].(map[string]interface{})
		scope.AllComputers = scopeMap["all_computers"].(bool)
		scope.AllJSSUsers = scopeMap["all_jss_users"].(bool)

		// Construct Computers
		if computers, ok := scopeMap["computers"].([]interface{}); ok {
			for _, comp := range computers {
				compMap := comp.(map[string]interface{})
				computer := jamfpro.MacOSConfigurationProfilesDataSubsetComputer{
					Computer: jamfpro.MacOSConfigurationProfilesDataSubsetComputerItem{
						ID:   compMap["id"].(int),
						Name: compMap["name"].(string),
						UDID: compMap["udid"].(string),
					},
				}
				scope.Computers = append(scope.Computers, computer)
			}
		}

		// Construct Buildings
		if buildings, ok := scopeMap["buildings"].([]interface{}); ok {
			for _, bld := range buildings {
				bldMap := bld.(map[string]interface{})
				building := jamfpro.MacOSConfigurationProfilesDataSubsetBuilding{
					Building: jamfpro.MacOSConfigurationProfilesDataSubsetBuildingItem{
						ID:   bldMap["id"].(int),
						Name: bldMap["name"].(string),
					},
				}
				scope.Buildings = append(scope.Buildings, building)
			}
		}

		// Construct Departments
		if departments, ok := scopeMap["departments"].([]interface{}); ok {
			for _, dept := range departments {
				deptMap := dept.(map[string]interface{})
				department := jamfpro.MacOSConfigurationProfilesDataSubsetDepartment{
					Department: jamfpro.MacOSConfigurationProfilesDataSubsetDepartmentItem{
						ID:   deptMap["id"].(int),
						Name: deptMap["name"].(string),
					},
				}
				scope.Departments = append(scope.Departments, department)
			}
		}

		// Construct ComputerGroups
		if computerGroups, ok := scopeMap["computer_groups"].([]interface{}); ok {
			for _, grp := range computerGroups {
				grpMap := grp.(map[string]interface{})
				computerGroup := jamfpro.MacOSConfigurationProfilesDataSubsetComputerGroup{
					ComputerGroup: jamfpro.MacOSConfigurationProfilesDataSubsetComputerGroupItem{
						ID:   grpMap["id"].(int),
						Name: grpMap["name"].(string),
					},
				}
				scope.ComputerGroups = append(scope.ComputerGroups, computerGroup)
			}
		}

		// Construct JSSUsers
		if jssUsers, ok := scopeMap["jss_users"].([]interface{}); ok {
			for _, usr := range jssUsers {
				usrMap := usr.(map[string]interface{})
				jssUser := jamfpro.MacOSConfigurationProfilesDataSubsetJSSUser{
					JSSUser: jamfpro.MacOSConfigurationProfilesDataSubsetJSSUserItem{
						ID:   usrMap["id"].(int),
						Name: usrMap["name"].(string),
					},
				}
				scope.JSSUsers = append(scope.JSSUsers, jssUser)
			}
		}

		// Construct JSSUserGroups
		if jssUserGroups, ok := scopeMap["jss_user_groups"].([]interface{}); ok {
			for _, grp := range jssUserGroups {
				grpMap := grp.(map[string]interface{})
				jssUserGroup := jamfpro.MacOSConfigurationProfilesDataSubsetJSSUserGroup{
					JSSUserGroup: jamfpro.MacOSConfigurationProfilesDataSubsetJSSUserGroupItem{
						ID:   grpMap["id"].(int),
						Name: grpMap["name"].(string),
					},
				}
				scope.JSSUserGroups = append(scope.JSSUserGroups, jssUserGroup)
			}
		}

		// Construct Limitations
		if limitations, ok := scopeMap["limitations"].([]interface{}); ok && len(limitations) > 0 {
			limMap := limitations[0].(map[string]interface{})
			var lim jamfpro.MacOSConfigurationProfilesDataSubsetLimitations

			// Construct Users in Limitations
			if users, ok := limMap["users"].([]interface{}); ok {
				for _, usr := range users {
					usrMap := usr.(map[string]interface{})
					user := jamfpro.MacOSConfigurationProfilesDataSubsetUser{
						User: jamfpro.MacOSConfigurationProfilesDataSubsetUserItem{
							ID:   usrMap["id"].(int),
							Name: usrMap["name"].(string),
						},
					}
					lim.Users = append(lim.Users, user)
				}
			}

			// Construct UserGroups in Limitations
			if userGroups, ok := limMap["user_groups"].([]interface{}); ok {
				for _, grp := range userGroups {
					grpMap := grp.(map[string]interface{})
					userGroup := jamfpro.MacOSConfigurationProfilesDataSubsetUserGroup{
						UserGroup: jamfpro.MacOSConfigurationProfilesDataSubsetUserGroupItem{
							ID:   grpMap["id"].(int),
							Name: grpMap["name"].(string),
						},
					}
					lim.UserGroups = append(lim.UserGroups, userGroup)
				}
			}

			// Construct NetworkSegments in Limitations
			if networkSegments, ok := limMap["network_segments"].([]interface{}); ok {
				for _, seg := range networkSegments {
					segMap := seg.(map[string]interface{})
					networkSegment := jamfpro.MacOSConfigurationProfilesDataSubsetNetworkSegment{
						NetworkSegment: jamfpro.MacOSConfigurationProfilesDataSubsetNetworkSegmentItem{
							ID:   segMap["id"].(int),
							UID:  segMap["uid"].(string),
							Name: segMap["name"].(string),
						},
					}
					lim.NetworkSegments = append(lim.NetworkSegments, networkSegment)
				}
			}

			// Construct IBeacons in Limitations
			if ibeacons, ok := limMap["ibeacons"].([]interface{}); ok {
				for _, ibc := range ibeacons {
					ibcMap := ibc.(map[string]interface{})
					ibeacon := jamfpro.MacOSConfigurationProfilesDataSubsetIBeacon{
						IBeacon: jamfpro.MacOSConfigurationProfilesDataSubsetIBeaconItem{
							ID:   ibcMap["id"].(int),
							Name: ibcMap["name"].(string),
						},
					}
					lim.IBeacons = append(lim.IBeacons, ibeacon)
				}
			}

			scope.Limitations = lim
		}

		// Construct Exclusions
		if exclusions, ok := scopeMap["exclusions"].([]interface{}); ok && len(exclusions) > 0 {
			excMap := exclusions[0].(map[string]interface{})
			var exc jamfpro.MacOSConfigurationProfilesDataSubsetExclusions

			// Construct Computers in Exclusions
			if computers, ok := excMap["computers"].([]interface{}); ok {
				for _, comp := range computers {
					compMap := comp.(map[string]interface{})
					computer := jamfpro.MacOSConfigurationProfilesDataSubsetComputer{
						Computer: jamfpro.MacOSConfigurationProfilesDataSubsetComputerItem{
							ID:   compMap["id"].(int),
							Name: compMap["name"].(string),
							UDID: compMap["udid"].(string),
						},
					}
					exc.Computers = append(exc.Computers, computer)
				}
			}

			// Construct Buildings in Exclusions
			if buildings, ok := excMap["buildings"].([]interface{}); ok {
				for _, bld := range buildings {
					bldMap := bld.(map[string]interface{})
					building := jamfpro.MacOSConfigurationProfilesDataSubsetBuilding{
						Building: jamfpro.MacOSConfigurationProfilesDataSubsetBuildingItem{
							ID:   bldMap["id"].(int),
							Name: bldMap["name"].(string),
						},
					}
					exc.Buildings = append(exc.Buildings, building)
				}
			}

			// Construct Departments in Exclusions
			if departments, ok := excMap["departments"].([]interface{}); ok {
				for _, dept := range departments {
					deptMap := dept.(map[string]interface{})
					department := jamfpro.MacOSConfigurationProfilesDataSubsetDepartment{
						Department: jamfpro.MacOSConfigurationProfilesDataSubsetDepartmentItem{
							ID:   deptMap["id"].(int),
							Name: deptMap["name"].(string),
						},
					}
					exc.Departments = append(exc.Departments, department)
				}
			}

			// Construct ComputerGroups in Exclusions
			if computerGroups, ok := excMap["computer_groups"].([]interface{}); ok {
				for _, grp := range computerGroups {
					grpMap := grp.(map[string]interface{})
					computerGroup := jamfpro.MacOSConfigurationProfilesDataSubsetComputerGroup{
						ComputerGroup: jamfpro.MacOSConfigurationProfilesDataSubsetComputerGroupItem{
							ID:   grpMap["id"].(int),
							Name: grpMap["name"].(string),
						},
					}
					exc.ComputerGroups = append(exc.ComputerGroups, computerGroup)
				}
			}

			// Construct UserGroups in Exclusions
			if userGroups, ok := excMap["user_groups"].([]interface{}); ok {
				for _, grp := range userGroups {
					grpMap := grp.(map[string]interface{})
					userGroup := jamfpro.MacOSConfigurationProfilesDataSubsetUserGroup{
						UserGroup: jamfpro.MacOSConfigurationProfilesDataSubsetUserGroupItem{
							ID:   grpMap["id"].(int),
							Name: grpMap["name"].(string),
						},
					}
					exc.UserGroups = append(exc.UserGroups, userGroup)
				}
			}

			// Construct NetworkSegments in Exclusions
			if networkSegments, ok := excMap["network_segments"].([]interface{}); ok {
				for _, seg := range networkSegments {
					segMap := seg.(map[string]interface{})
					networkSegment := jamfpro.MacOSConfigurationProfilesDataSubsetNetworkSegment{
						NetworkSegment: jamfpro.MacOSConfigurationProfilesDataSubsetNetworkSegmentItem{
							ID:   segMap["id"].(int),
							UID:  segMap["uid"].(string),
							Name: segMap["name"].(string),
						},
					}
					exc.NetworkSegments = append(exc.NetworkSegments, networkSegment)
				}
			}

			// Construct IBeacons in Exclusions
			if ibeacons, ok := excMap["ibeacons"].([]interface{}); ok {
				for _, ibc := range ibeacons {
					ibcMap := ibc.(map[string]interface{})
					ibeacon := jamfpro.MacOSConfigurationProfilesDataSubsetIBeacon{
						IBeacon: jamfpro.MacOSConfigurationProfilesDataSubsetIBeaconItem{
							ID:   ibcMap["id"].(int),
							Name: ibcMap["name"].(string),
						},
					}
					exc.IBeacons = append(exc.IBeacons, ibeacon)
				}
			}

			// Construct JSSUsers in Exclusions
			if jssUsers, ok := excMap["jss_users"].([]interface{}); ok {
				for _, usr := range jssUsers {
					usrMap := usr.(map[string]interface{})
					jssUser := jamfpro.MacOSConfigurationProfilesDataSubsetJSSUser{
						JSSUser: jamfpro.MacOSConfigurationProfilesDataSubsetJSSUserItem{
							ID:   usrMap["id"].(int),
							Name: usrMap["name"].(string),
						},
					}
					exc.JSSUsers = append(exc.JSSUsers, jssUser)
				}
			}

			// Construct JSSUserGroups in Exclusions
			if jssUserGroups, ok := excMap["jss_user_groups"].([]interface{}); ok {
				for _, grp := range jssUserGroups {
					grpMap := grp.(map[string]interface{})
					jssUserGroup := jamfpro.MacOSConfigurationProfilesDataSubsetJSSUserGroup{
						JSSUserGroup: jamfpro.MacOSConfigurationProfilesDataSubsetJSSUserGroupItem{
							ID:   grpMap["id"].(int),
							Name: grpMap["name"].(string),
						},
					}
					exc.JSSUserGroups = append(exc.JSSUserGroups, jssUserGroup)
				}
			}

			scope.Exclusions = exc
		}

	}

	// Construct SelfService section
	selfServiceData := d.Get("self_service").([]interface{})
	var selfService jamfpro.MacOSConfigurationProfilesDataSubsetSelfService
	if len(selfServiceData) > 0 {
		ssMap := selfServiceData[0].(map[string]interface{})
		selfService.InstallButtonText = ssMap["install_button_text"].(string)
		selfService.SelfServiceDescription = ssMap["self_service_description"].(string)
		selfService.ForceUsersToViewDescription = ssMap["force_users_to_view_description"].(bool)
		selfService.FeatureOnMainPage = ssMap["feature_on_main_page"].(bool)

		// Constructing SelfServiceIcon
		if iconData, ok := ssMap["self_service_icon"].([]interface{}); ok && len(iconData) > 0 {
			iconMap := iconData[0].(map[string]interface{})
			selfService.SelfServiceIcon = jamfpro.MacOSConfigurationProfilesDataSubsetSelfServiceIcon{
				ID:   iconMap["id"].(int),
				URI:  iconMap["uri"].(string),
				Data: iconMap["data"].(string),
			}
		}

		// Constructing SelfServiceCategories
		if categoriesData, ok := ssMap["self_service_categories"].([]interface{}); ok && len(categoriesData) > 0 {
			catMap := categoriesData[0].(map[string]interface{})
			selfService.SelfServiceCategories = jamfpro.MacOSConfigurationProfilesDataSubsetSelfServiceCategories{
				Category: jamfpro.MacOSConfigurationProfilesDataSubsetSelfServiceCategory{
					ID:        catMap["id"].(int),
					Name:      catMap["name"].(string),
					DisplayIn: catMap["display_in"].(bool),
					FeatureIn: catMap["feature_in"].(bool),
				},
			}
		}

	}

	return &jamfpro.ResponseMacOSConfigurationProfiles{
		General:     general,
		Scope:       scope,
		SelfService: selfService,
	}

}

// Helper function to generate diagnostics based on the error type
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

// ResourceJamfProMacOSConfigurationProfilesCreate is responsible for creating a new Jamf Pro macOS Configuration Profile in the remote system.
// The function:
// 1. Constructs the macOS Configuration Profile data using the provided Terraform configuration.
// 2. Calls the API to create the profile in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created profile.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProMacOSConfigurationProfilesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the create operation
	var createdProfile *jamfpro.ResponseMacOSConfigurationProfiles
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the macOS Configuration Profile
		profile := constructJamfProMacOSConfigurationProfile(d)

		// Check if the profile is nil
		if profile == nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the macOS Configuration Profile"))
		}

		// Directly call the API to create the resource
		createdProfile, err = conn.CreateMacOSConfigurationProfile(profile)
		if err != nil {
			// Check if the error is an APIError
			if apiErr, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiErr.StatusCode, apiErr.Message))
			}
			// For simplicity, we're considering all other errors as retryable
			return retry.RetryableError(err)
		}

		return nil
	})

	if err != nil {
		// If there's an error while creating the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "create")
	}

	// Set the ID of the created resource in the Terraform state
	d.SetId(strconv.Itoa(createdProfile.General.ID))

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProMacOSConfigurationProfilesRead(ctx, d, meta)
		if len(readDiags) > 0 {
			// If readDiags is not empty, it means there's an error, so we retry
			return retry.RetryableError(fmt.Errorf("failed to read the created resource"))
		}
		return nil
	})

	if err != nil {
		// If there's an error while updating the state for the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "update state for")
	}

	return diags
}

// ResourceJamfProMacOSConfigurationProfilesRead is responsible for reading the current state of a Jamf Pro macOS Configuration Profile from the remote system.
// The function:
// 1. Fetches the profile's current state using its ID. If it fails then obtain profile's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the profile being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProMacOSConfigurationProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	var profile *jamfpro.ResponseMacOSConfigurationProfiles

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		profileID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse profile ID: %v", convertErr))
		}

		// Try fetching the profile using the ID
		var apiErr error
		profile, apiErr = conn.GetMacOSConfigurationProfileByID(profileID)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			profileName := d.Get("name").(string)
			profile, apiErr = conn.GetMacOSConfigurationProfileByName(profileName)
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

	// Safely set attributes in the Terraform state
	if err := d.Set("name", profile.General.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("description", profile.General.Description); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("distribution_method", profile.General.DistributionMethod); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("user_removable", profile.General.UserRemovable); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("level", profile.General.Level); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("uuid", profile.General.UUID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("redeploy_on_update", profile.General.RedeployOnUpdate); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("payloads", profile.General.Payloads); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("site", []interface{}{map[string]interface{}{
		"id":   profile.General.Site.ID,
		"name": profile.General.Site.Name,
	}}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("category", []interface{}{map[string]interface{}{
		"id":   profile.General.Category.ID,
		"name": profile.General.Category.Name,
	}}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Setting 'scope' attribute
	scopeAttr := map[string]interface{}{
		"all_computers":   profile.Scope.AllComputers,
		"all_jss_users":   profile.Scope.AllJSSUsers,
		"computers":       constructNestedSliceOfMaps(profile.Scope.Computers, "Computer"),
		"buildings":       constructNestedSliceOfMaps(profile.Scope.Buildings, "Building"),
		"departments":     constructNestedSliceOfMaps(profile.Scope.Departments, "Department"),
		"computer_groups": constructNestedSliceOfMaps(profile.Scope.ComputerGroups, "ComputerGroup"),
		"jss_users":       constructNestedSliceOfMaps(profile.Scope.JSSUsers, "JSSUser"),
		"jss_user_groups": constructNestedSliceOfMaps(profile.Scope.JSSUserGroups, "JSSUserGroup"),
		"limitations":     constructLimitationsExclusions(profile.Scope.Limitations),
		"exclusions":      constructLimitationsExclusions(profile.Scope.Exclusions),
	}

	// Add the 'scope' to Terraform state
	if err := d.Set("scope", []interface{}{scopeAttr}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Setting 'self_service' attribute
	selfServiceAttr := map[string]interface{}{
		"install_button_text":             profile.SelfService.InstallButtonText,
		"self_service_description":        profile.SelfService.SelfServiceDescription,
		"force_users_to_view_description": profile.SelfService.ForceUsersToViewDescription,
		"feature_on_main_page":            profile.SelfService.FeatureOnMainPage,
	}

	// Constructing 'self_service_icon'
	if profile.SelfService.SelfServiceIcon.ID != 0 || profile.SelfService.SelfServiceIcon.URI != "" || profile.SelfService.SelfServiceIcon.Data != "" {
		selfServiceIcon := map[string]interface{}{
			"id":   profile.SelfService.SelfServiceIcon.ID,
			"uri":  profile.SelfService.SelfServiceIcon.URI,
			"data": profile.SelfService.SelfServiceIcon.Data,
		}
		selfServiceAttr["self_service_icon"] = []interface{}{selfServiceIcon}
	}

	// Constructing 'self_service_categories'
	if profile.SelfService.SelfServiceCategories.Category.ID != 0 || profile.SelfService.SelfServiceCategories.Category.Name != "" {
		selfServiceCategory := map[string]interface{}{
			"id":         profile.SelfService.SelfServiceCategories.Category.ID,
			"name":       profile.SelfService.SelfServiceCategories.Category.Name,
			"display_in": profile.SelfService.SelfServiceCategories.Category.DisplayIn,
			"feature_in": profile.SelfService.SelfServiceCategories.Category.FeatureIn,
		}
		selfServiceAttr["self_service_categories"] = []interface{}{selfServiceCategory}
	}

	// Set the 'self_service' attribute in the Terraform state
	if err := d.Set("self_service", []interface{}{selfServiceAttr}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags

}

// Helper function to construct a nested slice of maps for structures like 'computers', 'buildings', etc.
func constructNestedSliceOfMaps(entities interface{}, entityName string) []interface{} {
	var result []interface{}
	// Use reflection to iterate over the slice of entities
	v := reflect.ValueOf(entities)
	for i := 0; i < v.Len(); i++ {
		entity := v.Index(i).FieldByName(entityName).Interface()
		// Convert the entity to a map
		entityMap := structToMap(entity)
		result = append(result, entityMap)
	}
	return result
}

// Helper function to convert a struct to a map
func structToMap(obj interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i)
		out[field.Name] = v.Field(i).Interface()
	}
	return out
}

// ConstructLimitationsExclusions constructs a map for limitations or exclusions
func constructLimitationsExclusions(limExc interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(limExc)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typeOfS.Field(i)
		fieldName := field.Name
		fieldSlice := v.Field(i).Interface()

		// Create a slice of maps for each field
		entitySlice := constructNestedSliceOfMaps(fieldSlice, fieldName)
		result[strings.ToLower(fieldName)] = entitySlice // Use lowercase for the key
	}
	return result
}

// ResourceJamfProMacOSConfigurationProfilesUpdate is responsible for updating an existing Jamf Pro macOS Configuration Profile on the remote system.
func ResourceJamfProMacOSConfigurationProfilesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the update operation
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		// Construct the updated macOS configuration profile
		profile := constructJamfProMacOSConfigurationProfile(d)

		// Convert the ID from the Terraform state into an integer to be used for the API request
		profileID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse profile ID: %v", convertErr))
		}

		// Directly call the API to update the profile
		_, apiErr := conn.UpdateMacOSConfigurationProfileByID(profileID, profile)
		if apiErr != nil {
			// Handle the APIError
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If the update by ID fails, try updating by name
			profileName := d.Get("name").(string)
			_, apiErr = conn.UpdateMacOSConfigurationProfileByName(profileName, profile)
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
		readDiags := ResourceJamfProMacOSConfigurationProfilesRead(ctx, d, meta)
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

// ResourceJamfProMacOSConfigurationProfilesDelete is responsible for deleting a macOS Configuration Profile.
func ResourceJamfProMacOSConfigurationProfilesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn
	var diags diag.Diagnostics

	// Use the retry function for the **DELETE** operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		macOSConfigurationProfileID, convertErr := strconv.Atoi(d.Id())
		if convertErr != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to parse department ID: %v", convertErr))
		}

		// Directly call the API to **DELETE** the resource
		apiErr := conn.DeleteMacOSConfigurationProfileByID(macOSConfigurationProfileID)
		if apiErr != nil {
			// If the **DELETE** by ID fails, try deleting by name
			siteName := d.Get("name").(string)
			apiErr = conn.DeleteDepartmentByName(siteName)
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
