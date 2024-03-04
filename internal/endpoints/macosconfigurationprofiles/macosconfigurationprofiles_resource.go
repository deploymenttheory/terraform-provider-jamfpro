package macosconfigurationprofiles

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
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

func constructJamfProMacOSConfigurationProfile(d *schema.ResourceData) (*jamfpro.ResourceMacOSConfigurationProfile, error) {
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
		out.General.Category = jamfpro.SharedResourceCategory{
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
	xmlData, _ := xml.MarshalIndent(out, "", "  ")
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
