package macosconfigurationprofiles

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const JamfProResourceMacOSConfigurationProfile = "macos_configuration_profile"

func ResourceJamfProMacOSConfigurationProfiles() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProMacOSConfigurationProfilesCreate,
		ReadContext:   ResourceJamfProMacOSConfigurationProfilesRead,
		UpdateContext: ResourceJamfProMacOSConfigurationProfilesUpdate,
		DeleteContext: ResourceJamfProMacOSConfigurationProfilesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Second),
			Read:   schema.DefaultTimeout(5 * time.Second),
			Update: schema.DefaultTimeout(5 * time.Second),
			Delete: schema.DefaultTimeout(5 * time.Second),
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
				Description: "The site to which the configuration profile is scoped.",
				Optional:    true,
				Default:     nil,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The unique identifier of the site to which the configuration profile is scoped.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the site to which the configuration profile is scoped.",
						},
					},
				},
			},
			"category": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The category to which the configuration profile is scoped.",
				Optional:    true,
				Default:     nil,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The unique identifier of the category to which the configuration profile is scoped.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the category to which the configuration profile is scoped.",
						},
					},
				},
			},
			"distribution_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Install Automatically",
				Description:  "The distribution method for the configuration profile. Available options are: 'push', 'install_enterprise', 'install_user_initiated', 'install_system', 'install_self_service'.",
				ValidateFunc: validation.StringInSlice([]string{"Make Available in Self Service", "Install Automatically"}, false),
			},
			"user_removeable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the configuration profile is user removeable.",
			},
			"level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "System",
				Description:  "The level of the configuration profile. Available options are: 'computer', 'user'.",
				ValidateFunc: validation.StringInSlice([]string{"Computer", "User", "System"}, false),
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the configuration profile.",
			},
			// "redeploy_on_update": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Default:     "true",
			// 	Description: "Whether the configuration profile is redeployed on update.",
			// },
			// "payloads": {},
			"scope": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The scope of the configuration profile.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_computers": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the configuration profile is scoped to all computers.",
						},
						// "all_jss_users": {
						// 	Type:        schema.TypeBool,
						// 	Optional:    true,
						// 	Default:     false,
						// 	Description: "Whether the configuration profile is scoped to all JSS users.",
						// },
						"computers": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									// "name": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	Description: "The name of the computer to which the configuration profile is scoped.",
									// },
									// "udid": {
									// 	Type:        schema.TypeString,
									// 	Optional:    true,
									// 	Description: "The UDID of the computer to which the configuration profile is scoped.",
									// },
								},
							},
						},
						// "computer_groups": {
						// 	Type:        schema.TypeList,
						// 	Optional:    true,
						// 	Description: "The computer groups to which the configuration profile is scoped.",
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"id": {
						// 				Type:        schema.TypeInt,
						// 				Required:    true,
						// 				Description: "The unique identifier of the computer group to which the configuration profile is scoped.",
						// 			},
						// 			"name": {
						// 				Type:        schema.TypeString,
						// 				Optional:    true,
						// 				Description: "The name of the computer group to which the configuration profile is scoped.",
						// 			},
						// 		},
						// 	},
						// },
						// "jss_users": {
						// 	Type:        schema.TypeList,
						// 	Optional:    true,
						// 	Description: "The JSS users to which the configuration profile is scoped.",
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"id": {
						// 				Type:        schema.TypeInt,
						// 				Required:    true,
						// 				Description: "The unique identifier of the JSS user to which the configuration profile is scoped.",
						// 			},
						// 			"name": {
						// 				Type:        schema.TypeString,
						// 				Optional:    true,
						// 				Description: "The name of the JSS user to which the configuration profile is scoped.",
						// 			},
						// 		},
						// 	},
						// },
						// "jss_user_groups": {
						// 	Type:        schema.TypeList,
						// 	Optional:    true,
						// 	Description: "The JSS user groups to which the configuration profile is scoped.",
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"id": {
						// 				Type:        schema.TypeInt,
						// 				Required:    true,
						// 				Description: "The unique identifier of the JSS user group to which the configuration profile is scoped.",
						// 			},
						// 			"name": {
						// 				Type:        schema.TypeString,
						// 				Optional:    true,
						// 				Description: "The name of the JSS user group to which the configuration profile is scoped.",
						// 			},
						// 		},
						// 	},
						// },
						// "buildings": {
						// 	Type:        schema.TypeList,
						// 	Optional:    true,
						// 	Description: "The buildings to which the configuration profile is scoped.",
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"id": {
						// 				Type:        schema.TypeInt,
						// 				Required:    true,
						// 				Description: "The unique identifier of the building to which the configuration profile is scoped.",
						// 			},
						// 			"name": {
						// 				Type:        schema.TypeString,
						// 				Optional:    true,
						// 				Description: "The name of the building to which the configuration profile is scoped.",
						// 			},
						// 		},
						// 	},
						// },
						// "departments": {
						// 	Type:        schema.TypeList,
						// 	Optional:    true,
						// 	Description: "The departments to which the configuration profile is scoped.",
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"id": {
						// 				Type:        schema.TypeInt,
						// 				Required:    true,
						// 				Description: "The unique identifier of the department to which the configuration profile is scoped.",
						// 			},
						// 			"name": {
						// 				Type:        schema.TypeString,
						// 				Optional:    true,
						// 				Description: "The name of the department to which the configuration profile is scoped.",
						// 			},
						// 		},
						// 	},
						// },
						// "limitations": {},
						// "exclusions":  {},
					},
				},
			},
		},
	}
}

func constructJamfProMacOSConfigurationProfile(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceMacOSConfigurationProfile, error) {
	// Main obj with fields which do not require processing
	out := jamfpro.ResourceMacOSConfigurationProfile{
		General: jamfpro.MacOSConfigurationProfileSubsetGeneral{
			Name:               d.Get("name").(string),
			Description:        d.Get("description").(string),
			DistributionMethod: d.Get("distribution_method").(string),
			UserRemovable:      d.Get("user_removeable").(bool),
			Level:              d.Get("level").(string),
			UUID:               d.Get("uuid").(string),
			// RedeployOnUpdate:   d.Get("redeploy_on_update").(string),
		},
		Scope:       jamfpro.MacOSConfigurationProfileSubsetScope{},
		SelfService: jamfpro.MacOSConfigurationProfileSubsetSelfService{},
	}

	// Fields with processing

	// Site
	if len(d.Get("site").([]interface{})) != 0 {
		out.General.Site = jamfpro.SharedResourceSite{
			ID:   d.Get("site.0.id").(int),
			Name: d.Get("site.0.name").(string),
		}
	} else {
		log.Println("NO SITE") // TODO probably put some logging here
	}

	// Category
	if len(d.Get("category").([]interface{})) != 0 {
		out.General.Category = &jamfpro.SharedResourceCategory{
			ID:   d.Get("category.0.id").(int),
			Name: d.Get("category.0.name").(string),
		}
	} else {
		log.Println("NO CATEGORY") // TODO probably put some logging here
	}

	// Scope

	out.Scope.AllComputers = d.Get("scope.0.all_computers").(bool)

	log.Println("LOGHERE-OUT")
	log.Printf("%+v\n", out)

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
	log.Println("LOGHERE-READSTART")
	// API Stuff
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

	// Stating

	// ID
	if err := d.Set("id", resourceID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Name
	if err := d.Set("name", resp.General.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Description
	if err := d.Set("description", resp.General.Description); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Site
	if resp.General.Site.ID != -1 && resp.General.Site.Name != "None" {
		out_site := []map[string]interface{}{
			{
				"id":   resp.General.Site.ID,
				"name": resp.General.Site.Name,
			},
		}

		if err := d.Set("site", out_site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		log.Println("Not stating default site response") // TODO probably put some logging here
	}

	// Category
	if resp.General.Category.ID != -1 && resp.General.Category.Name != "No category assigned" {
		out_category := []map[string]interface{}{
			{
				"id":   resp.General.Category.ID,
				"name": resp.General.Category.Name,
			},
		}
		if err := d.Set("category", out_category); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		log.Printf("Not stating default category response") // TODO probably put some logging here
	}

	// Distribution Method
	if err := d.Set("distribution_method", resp.General.DistributionMethod); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// User Removeable
	if err := d.Set("user_removeable", resp.General.UserRemovable); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Level
	if err := d.Set("level", resp.General.Level); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// UUID
	if err := d.Set("uuid", resp.General.UUID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Redeploy On Update - not in ui
	// if err := d.Set("redeploy_on_update", resp.General.RedeployOnUpdate); err != nil {
	// 	diags = append(diags, diag.FromErr(err)...)
	// }

	// Scope

	out_scope := make([]map[string]interface{}, 0)
	out_scope = append(out_scope, make(map[string]interface{}, 1))

	// All computers
	out_scope[0]["all_computers"] = resp.Scope.AllComputers
	log.Println("FLAG 1")
	log.Println(resp.Scope.AllComputers)
	log.Println(out_scope)

	//////////////////////////////////////////////////
	// Computers

	if len(resp.Scope.Computers) > 0 {

		// Define computers list
		var out_computers []map[string]interface{}

		// Add empty item to that list
		out_computers = append(out_computers, make(map[string]interface{}))

		// Get Ids from HCL
		var hclComputerIds []int
		for _, v := range resp.Scope.Computers {
			hclComputerIds = append(hclComputerIds, v.ID)
		}

		// Put IDs in empty item in computers list
		out_computers[0]["id"] = hclComputerIds

		// Add list to parent scope
		out_scope[0]["computers"] = out_computers
	}

	//////////////////////////////////////////////////

	// Write scope to state
	err = d.Set("scope", out_scope)
	if err != nil {
		log.Println("FLAG 2")
		diags = append(diags, diag.FromErr(err)...)
	}
	log.Println("LOGHERE-READEND")

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
	var resourceIDInt int
	resourceID := d.Id()
	resourceIDInt, convErr := strconv.Atoi(resourceID)
	if convErr != nil {
		return diag.FromErr(convErr)
	}

	resourceName := d.Get("name").(string)
	var apiErrorCode int

	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemDelete, hclog.Info)

	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {

		apiErr := conn.DeleteMacOSConfigurationProfileByID(resourceIDInt)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPIDeleteFailureByID(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			apiErr = conn.DeleteMacOSConfigurationProfileByName(resourceName)
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
