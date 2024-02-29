// accounts_resource.go
package accounts

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProAccount defines the schema and CRUD operations for managing buildings in Terraform.
func ResourceJamfProAccounts() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProAccountCreate,
		ReadContext:   ResourceJamfProAccountRead,
		UpdateContext: ResourceJamfProAccountUpdate,
		DeleteContext: ResourceJamfProAccountDelete,
		CustomizeDiff: customDiffAccounts,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the jamf pro account.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the jamf pro account.",
			},
			"directory_user": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the user is a directory user.",
			},
			"full_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The full name of the account user.",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email of the account user.",
			},
			"email_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email address of the account user.",
			},
			"enabled": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access status of the account (“enabled” or “disabled”).",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := util.GetString(val)
					if v == "Enabled" || v == "Disabled" {
						return
					}
					errs = append(errs, fmt.Errorf("%q must be either 'Enabled' or 'Disabled', got: %s", key, v))
					return warns, errs
				},
			},
			"ldap_server": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "LDAP server information associated with the account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The ID of the LDAP server.",
							Default:     "",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the LDAP server.",
							Computed:    true,
						},
					},
				},
			},
			"force_password_change": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the user is forced to change password on next login.",
			},
			"access_level": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The access level of the account. This can be either Full Access, scoped to a jamf pro site with Site Access, or scoped to a jamf pro account group with Group Access",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := util.GetString(val)
					if v == "Full Access" || v == "Site Access" || v == "Group Access" {
						return
					}
					errs = append(errs, fmt.Errorf("%q must be either 'Full Access' or 'Site Access' or 'Group Access', got: %s", key, v))
					return warns, errs
				},
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The password for the account.",
				Sensitive:   true,
			},
			"privilege_set": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The privilege set assigned to the account.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := util.GetString(val)
					validPrivileges := []string{"Administrator", "Auditor", "Enrollment Only", "Custom"}
					for _, validPriv := range validPrivileges {
						if v == validPriv {
							return // Valid value found, return without error
						}
					}
					errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validPrivileges, v))
					return warns, errs
				},
			},
			"site": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The site information associated with the account group if access_level is set to Site Access.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Jamf Pro Site ID. Value defaults to '0' aka not used.",
							Default:     "",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Jamf Pro Site Name",
							Computed:    true,
						},
					},
				},
			},
			"groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A set of group names and IDs associated with the account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"site": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"jss_objects_privileges": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Privileges related to JSS Objects.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: common.ValidateJSSObjectsPrivileges,
							},
						},
						"jss_settings_privileges": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Privileges related to JSS Settings.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: common.ValidateJSSSettingsPrivileges,
							},
						},
						"jss_actions_privileges": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Privileges related to JSS Actions.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: common.ValidateJSSActionsPrivileges,
							},
						},
						"casper_admin_privileges": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Privileges related to Casper Admin.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: common.ValidateCasperAdminPrivileges,
							},
						},
						"casper_remote_privileges": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Privileges related to Casper Remote.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"casper_imaging_privileges": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Privileges related to Casper Imaging.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"recon_privileges": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Privileges related to Recon.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"jss_objects_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to JSS Objects.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: common.ValidateJSSObjectsPrivileges,
				},
			},
			"jss_settings_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to JSS Settings.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: common.ValidateJSSSettingsPrivileges,
				},
			},
			"jss_actions_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to JSS Actions.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: common.ValidateJSSActionsPrivileges,
				},
			},
			"casper_admin_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to Casper Admin.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: common.ValidateCasperAdminPrivileges,
				},
			},
			"casper_remote_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to Casper Remote.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"casper_imaging_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to Casper Imaging.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"recon_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges related to Recon.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

const (
	JamfProResourceAccount = "Account"
)

// ResourceJamfProAccountCreate is responsible for creating a new Jamf Pro Script in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProAccount(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Account: %v", err))
	}

	// Retry the API call to create the site in Jamf Pro
	var creationResponse *jamfpro.ResponseAccountCreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateAccount(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Account '%s' after retries: %v", resource.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProAccountRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProAccountRead is responsible for reading the current state of a Jamf Pro Account Group Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	var resource *jamfpro.ResourceAccount

	// Read operation with retry
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = conn.GetAccountByID(resourceIDInt)
		if apiErr != nil {
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the resource, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		d.SetId("") // Remove from Terraform state if unable to read after retries
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Account with ID '%s' after retries: %v", resourceID, err))
	}

	// Update the Terraform state with account attributes
	d.Set("name", resource.Name)
	d.Set("directory_user", resource.DirectoryUser)
	d.Set("full_name", resource.FullName)
	d.Set("email", resource.Email)
	d.Set("enabled", resource.Enabled)

	// Update LDAP server information
	if resource.LdapServer.ID != 0 || resource.LdapServer.Name != "" {
		ldapServer := make(map[string]interface{})
		ldapServer["id"] = resource.LdapServer.ID
		ldapServer["name"] = resource.LdapServer.Name
		d.Set("ldap_server", []interface{}{ldapServer})
	} else {
		d.Set("ldap_server", []interface{}{}) // Clear the LDAP server data if not present
	}

	d.Set("force_password_change", resource.ForcePasswordChange)
	d.Set("access_level", resource.AccessLevel)
	// Set password only if it's provided in the configuration
	if _, ok := d.GetOk("password"); ok {
		d.Set("password", resource.Password)
	}
	d.Set("privilege_set", resource.PrivilegeSet)

	// Update site information
	if resource.Site.ID != 0 || resource.Site.Name != "" {
		site := make(map[string]interface{})
		site["id"] = resource.Site.ID
		site["name"] = resource.Site.Name
		d.Set("site", []interface{}{site})
	} else {
		d.Set("site", []interface{}{}) // Clear the site data if not present
	}

	// Construct and set the groups attribute
	groups := make([]interface{}, len(resource.Groups))
	for i, group := range resource.Groups {
		groupMap := make(map[string]interface{})
		groupMap["name"] = group.Name
		groupMap["id"] = group.ID

		// Construct Site subfield
		if group.Site.ID != 0 || group.Site.Name != "" {
			site := make(map[string]interface{})
			site["id"] = group.Site.ID
			site["name"] = group.Site.Name
			groupMap["site"] = []interface{}{site}
		} else {
			groupMap["site"] = []interface{}{} // Clear the Site data if not present
		}

		// Map privileges from the AccountSubsetPrivileges struct to the Terraform schema
		groupMap["jss_objects_privileges"] = group.Privileges.JSSObjects
		groupMap["jss_settings_privileges"] = group.Privileges.JSSSettings
		groupMap["jss_actions_privileges"] = group.Privileges.JSSActions
		groupMap["casper_admin_privileges"] = group.Privileges.CasperAdmin
		groupMap["casper_remote_privileges"] = group.Privileges.CasperRemote
		groupMap["casper_imaging_privileges"] = group.Privileges.CasperImaging
		groupMap["recon_privileges"] = group.Privileges.Recon

		groups[i] = groupMap
	}

	if err := d.Set("groups", groups); err != nil {
		return diag.FromErr(err)
	}

	// Update privileges
	if err := d.Set("jss_objects_privileges", resource.Privileges.JSSObjects); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("jss_settings_privileges", resource.Privileges.JSSSettings); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("jss_actions_privileges", resource.Privileges.JSSActions); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("casper_admin_privileges", resource.Privileges.CasperAdmin); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("casper_remote_privileges", resource.Privileges.CasperRemote); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("casper_imaging_privileges", resource.Privileges.CasperImaging); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("recon_privileges", resource.Privileges.Recon); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// ResourceJamfProAccountUpdate is responsible for updating an existing Jamf Pro Account Group on the remote system.
func ResourceJamfProAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	resource, err := constructJamfProAccount(ctx, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Account for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateAccountByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Account '%s' (ID: %s) after retries: %v", resource.Name, resourceID, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProAccountRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProAccountDelete is responsible for deleting a Jamf Pro account .
func ResourceJamfProAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		apiErr := conn.DeleteAccountByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeleteAccountByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Account '%s' (ID: %s) after retries: %v", d.Get("extension").(string), resourceID, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
