package macosconfigurationprofiles

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

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
				Default:     false,
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
			// "payload": {},
			// "redeploy_on_update": { // TODO Review this, missing from the gui
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
						"all_jss_users": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether the configuration profile is scoped to all JSS users.",
						},
						"computer_ids": {
							Type:        schema.TypeList,
							Description: "The computers to which the configuration profile is scoped",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"computer_group_ids": {
							Type:        schema.TypeList,
							Description: "The computer groups to which the configuration profile is scoped",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"jss_user_ids": {
							Type:        schema.TypeList,
							Description: "The jss users to which the configuration profile is scoped",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"jss_user_group_ids": {
							Type:        schema.TypeList,
							Description: "The jss user groups to which the configuration profile is scoped",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"building_ids": {
							Type:        schema.TypeList,
							Description: "The buildings to which the configuration profile is scoped",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"department_ids": {
							Type:        schema.TypeList,
							Description: "The departments to which the configuration profile is scoped",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"limitations": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "The limitations within the scope",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_names": {
										Type:        schema.TypeList,
										Description: "The limited users",
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"network_segment_ids": {
										Type:        schema.TypeList,
										Description: "The limited network segments",
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"ibeacon_ids": {
										Type:        schema.TypeList,
										Description: "The limited ibeacons",
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"user_group_ids": {
										Type:        schema.TypeList,
										Description: "The limited user groups",
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"exclusions": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "The limitations within the scope",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"computer_ids": {
										Type:        schema.TypeList,
										Description: "excluded computers",
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"computer_group_ids": {
										Type:        schema.TypeList,
										Description: "excluded computer groups",
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									// "user_ids": {},
									// "user_group_ids": {},
									"building_ids": {
										Type:        schema.TypeList,
										Description: "excluded buildings",
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"department_ids": {
										Type:        schema.TypeList,
										Description: "excluded departments",
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"network_segment_ids": {
										Type:        schema.TypeList,
										Description: "excluded network segments",
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"jss_user_ids": {
										Type:        schema.TypeList,
										Description: "excluded jss users",
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"jss_user_group_ids": {
										Type:        schema.TypeList,
										Description: "excluded jss user groups",
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"ibeacon_ids": {
										Type:        schema.TypeList,
										Description: "excluded ibeacons",
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
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

func constructJamfProMacOSConfigurationProfile(d *schema.ResourceData) (*jamfpro.ResourceMacOSConfigurationProfile, error) {
	// Main obj with fields which do not require processing
	out := jamfpro.ResourceMacOSConfigurationProfile{
		General: jamfpro.MacOSConfigurationProfileSubsetGeneral{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			// Site processed
			// Category processed
			DistributionMethod: d.Get("distribution_method").(string),
			UserRemovable:      d.Get("user_removeable").(bool),
			Level:              d.Get("level").(string),
			UUID:               d.Get("uuid").(string), // TODO not sure if this is needed as it's computed
			// RedeployOnUpdate:   d.Get("redeploy_on_update").(string), // TODO Review this, I don't think it's in the UI
			// Paylods processed
		},
		Scope:       jamfpro.MacOSConfigurationProfileSubsetScope{},
		SelfService: jamfpro.MacOSConfigurationProfileSubsetSelfService{},
	}

	// Fields which require processing

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
		out.General.Category = jamfpro.SharedResourceCategory{
			ID:   d.Get("category.0.id").(int),
			Name: d.Get("category.0.name").(string),
		}
	} else {
		log.Println("NO CATEGORY") // TODO probably put some logging here
	}

	// Scope
	var err error

	// Scope - Targets

	// Bools
	out.Scope.AllComputers = d.Get("scope.0.all_computers").(bool)
	out.Scope.AllJSSUsers = d.Get("scope.0.all_jss_users").(bool)

	// Computers
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetComputer, int]("scope.0.computer_ids", "ID", d, &out.Scope.Computers)
	if err != nil {
		return nil, err
	}

	// Computer Groups
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetComputerGroup, int]("scope.0.computer_group_ids", "ID", d, &out.Scope.ComputerGroups)
	if err != nil {
		return nil, err
	}

	// JSS Users
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetJSSUser, int]("scope.0.jss_user_ids", "ID", d, &out.Scope.JSSUsers)
	if err != nil {
		return nil, err
	}

	// JSS User Groups
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetJSSUserGroup, int]("scope.0.jss_user_group_ids", "ID", d, &out.Scope.JSSUserGroups)
	if err != nil {
		return nil, err
	}

	// Buildings
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetBuilding, int]("scope.0.building_ids", "ID", d, &out.Scope.Buildings)
	if err != nil {
		return nil, err
	}

	// Departments
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetDepartment, int]("scope.0.department_ids", "ID", d, &out.Scope.Departments)
	if err != nil {
		return nil, err
	}

	// Scope - Limitations

	// Users
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetUser, string]("scope.0.limitations.0.user_names", "Name", d, &out.Scope.Limitations.Users)
	if err != nil {
		return nil, err
	}

	// Network Segment
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetNetworkSegment, int]("scope.0.limitations.0.network_segment_ids", "ID", d, &out.Scope.Limitations.NetworkSegments)
	if err != nil {
		return nil, err
	}

	// IBeacons
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetIBeacon, int]("scope.0.limitations.0.ibeacon_ids", "ID", d, &out.Scope.Limitations.IBeacons)
	if err != nil {
		return nil, err
	}

	// User Groups
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetUserGroup, int]("scope.0.limitations.0.user_group_ids", "ID", d, &out.Scope.Limitations.UserGroups)
	if err != nil {
		return nil, err
	}

	// Scope - Limitations

	// Computers
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetComputer, int]("scope.0.exclusions.0.computer_ids", "ID", d, &out.Scope.Exclusions.Computers)
	if err != nil {
		return nil, err
	}

	// Computer Groups
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetComputerGroup, int]("scope.0.exclusions.0.computer_group_ids", "ID", d, &out.Scope.Exclusions.ComputerGroups)
	if err != nil {
		return nil, err
	}

	// Buildings
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetBuilding, int]("scope.0.exclusions.0.building_ids", "ID", d, &out.Scope.Exclusions.Buildings)
	if err != nil {
		return nil, err
	}

	// Departments
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetDepartment, int]("scope.0.exclusions.0.department_ids", "ID", d, &out.Scope.Exclusions.Departments)
	if err != nil {
		return nil, err
	}

	// Network Segments
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetNetworkSegment, int]("scope.0.exclusions.0.network_segment_ids", "ID", d, &out.Scope.Exclusions.NetworkSegments)
	if err != nil {
		return nil, err
	}

	// JSS Users
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetJSSUser, int]("scope.0.exclusions.0.jss_user_ids", "ID", d, &out.Scope.Exclusions.JSSUsers)
	if err != nil {
		return nil, err
	}

	// JSS User Groups
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetJSSUserGroup, int]("scope.0.exclusions.0.jss_user_group_ids", "ID", d, &out.Scope.Exclusions.JSSUserGroups)
	if err != nil {
		return nil, err
	}

	// IBeacons
	err = GetAttrsListFromHCL[jamfpro.MacOSConfigurationProfileSubsetIBeacon, int]("scope.0.exclusions.0.ibeacon_ids", "ID", d, &out.Scope.Exclusions.IBeacons)
	if err != nil {
		return nil, err
	}

	// Debug
	xmlData, _ := xml.MarshalIndent(out, "", "	")
	log.Println(string(xmlData))

	return &out, nil
}

func ResourceJamfProMacOSConfigurationProfilesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProMacOSConfigurationProfile(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro MacOs Configuration Profile: %v", err))
	}

	// Retry the API call to create the MacOs Configuration Profile in Jamf Pro
	var creationResponse *jamfpro.ResponseMacOSConfigurationProfileCreationUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateMacOSConfigurationProfile(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro MacOs Configuration Profile '%s' after retries: %v", resource.General.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Read the MacOs Configuration Profile to ensure the Terraform state is up to date
	readDiags := ResourceJamfProMacOSConfigurationProfilesRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

func ResourceJamfProMacOSConfigurationProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var resp *jamfpro.ResourceMacOSConfigurationProfile

	// Read operation with retry
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resp, apiErr = conn.GetMacOSConfigurationProfileByID(resourceIDInt)
		if apiErr != nil {
			if strings.Contains(apiErr.Error(), "404") || strings.Contains(apiErr.Error(), "410") {
				// Return non-retryable error with a message to avoid SDK issues
				return retry.NonRetryableError(fmt.Errorf("resource not found, marked for deletion"))
			}
			// Retry for other types of errors
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	// If err is not nil, check if it's due to the resource being not found
	if err != nil {
		if err.Error() == "resource not found, marked for deletion" {
			// Resource not found, remove from Terraform state
			d.SetId("")
			// Append a warning diagnostic and return
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("Jamf Pro Site with ID '%s' was not found on the server and is marked for deletion from terraform state.", resourceID),
			})
			return diags
		}

		// For other errors, return an error diagnostic
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Site with ID '%s' after retries: %v", resourceID, err))
	}

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
		log.Println("Not stating default category response") // TODO probably put some logging here
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

	// Setup scope obj
	out_scope := make([]map[string]interface{}, 0)
	out_scope = append(out_scope, make(map[string]interface{}, 1))

	// Bools
	out_scope[0]["all_computers"] = resp.Scope.AllComputers
	out_scope[0]["all_jss_users"] = resp.Scope.AllJSSUsers

	// Computers
	if len(resp.Scope.Computers) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Computers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["computer_ids"] = listOfIds
	}

	// Computer Groups
	if len(resp.Scope.ComputerGroups) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.ComputerGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["computer_group_ids"] = listOfIds
	}

	// JSS Users
	if len(resp.Scope.JSSUsers) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.JSSUsers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["jss_user_ids"] = listOfIds
	}

	// JSS User Groups
	if len(resp.Scope.JSSUserGroups) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.JSSUserGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["jss_user_group_ids"] = listOfIds
	}

	// Buildings
	if len(resp.Scope.Buildings) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Buildings {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["building_ids"] = listOfIds
	}

	// Departments
	if len(resp.Scope.Departments) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Departments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope[0]["department_ids"] = listOfIds
	}

	// Scope Limitations ////////////////////////////////

	out_scope_limitations := make([]map[string]interface{}, 0)
	out_scope_limitations = append(out_scope_limitations, make(map[string]interface{}))

	// Users
	if len(resp.Scope.Limitations.Users) > 0 {
		var listOfNames []string
		for _, v := range resp.Scope.Limitations.Users {
			listOfNames = append(listOfNames, v.Name)
		}
		out_scope_limitations[0]["user_names"] = listOfNames
	}

	// Network Segments
	if len(resp.Scope.Limitations.NetworkSegments) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Limitations.NetworkSegments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_limitations[0]["network_segment_ids"] = listOfIds
	}

	// IBeacons
	if len(resp.Scope.Limitations.IBeacons) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Limitations.IBeacons {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_limitations[0]["ibeacon_ids"] = listOfIds
	}

	// User Groups
	if len(resp.Scope.Limitations.UserGroups) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Limitations.UserGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_limitations[0]["user_group_ids"] = listOfIds
	}

	// Append Limitations
	out_scope[0]["limitations"] = out_scope_limitations

	// Scope Exclusions ////////////////////////////

	out_scope_exclusions := make([]map[string]interface{}, 0)
	out_scope_exclusions = append(out_scope_exclusions, make(map[string]interface{}))

	// Computers
	if len(resp.Scope.Exclusions.Computers) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Exclusions.Computers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["computer_ids"] = listOfIds
	}

	// Computer Groups
	if len(resp.Scope.Exclusions.ComputerGroups) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Exclusions.ComputerGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["computer_group_ids"] = listOfIds
	}

	// Buildings
	if len(resp.Scope.Exclusions.Buildings) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Exclusions.Buildings {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["building_ids"] = listOfIds
	}

	// Departments
	if len(resp.Scope.Exclusions.Departments) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Exclusions.Departments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["department_ids"] = listOfIds
	}

	// Network Segments
	if len(resp.Scope.Exclusions.NetworkSegments) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Exclusions.NetworkSegments {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["network_segment_ids"] = listOfIds
	}

	// JSS Users
	if len(resp.Scope.Exclusions.JSSUsers) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Exclusions.JSSUsers {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["jss_user_ids"] = listOfIds
	}

	// JSS User Groups
	if len(resp.Scope.Exclusions.JSSUserGroups) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Exclusions.JSSUserGroups {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["jss_user_group_ids"] = listOfIds
	}

	// IBeacons
	if len(resp.Scope.Exclusions.IBeacons) > 0 {
		var listOfIds []int
		for _, v := range resp.Scope.Exclusions.IBeacons {
			listOfIds = append(listOfIds, v.ID)
		}
		out_scope_exclusions[0]["ibeacon_ids"] = listOfIds
	}

	// Append Exclusions
	out_scope[0]["exclusions"] = out_scope_exclusions

	// Set Scope to state
	err = d.Set("scope", out_scope)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

func ResourceJamfProMacOSConfigurationProfilesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Construct the resource object
	resource, err := constructJamfProMacOSConfigurationProfile(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro MacOs Configuration Profile for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateMacOSConfigurationProfileByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro MacOs Configuration Profile '%s' (ID: %d) after retries: %v", resource.General.Name, resourceIDInt, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProMacOSConfigurationProfilesRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

func ResourceJamfProMacOSConfigurationProfilesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Attempt to delete by ID
		apiErr := conn.DeleteMacOSConfigurationProfileByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeleteMacOSConfigurationProfileByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro MacOs Configuration Profile '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}

// TODO rename this func and put it somewhere else
func GetAttrsListFromHCL[NestedObjectType any, ListItemPrimitiveType any](path string, target_field string, d *schema.ResourceData, home *[]NestedObjectType) (err error) {
	getAttr, ok := d.GetOk(path)

	if len(getAttr.([]interface{})) == 0 {
		return nil
	}

	if ok {
		outList := make([]NestedObjectType, 0)
		for _, v := range getAttr.([]interface{}) {
			var newObj NestedObjectType
			newObjReflect := reflect.ValueOf(&newObj).Elem()
			idField := newObjReflect.FieldByName(target_field)

			if idField.IsValid() && idField.CanSet() {
				idField.Set(reflect.ValueOf(v.(ListItemPrimitiveType)))
			} else {
				return fmt.Errorf("error cannot set field line 695") // TODO write this error
			}

			outList = append(outList, newObj)

		}

		if len(outList) > 0 {
			*home = outList
		} else {
			log.Println("list is empty")
		}

		return nil
	}
	return fmt.Errorf("no path found/no scoped items at %v", path)
}
