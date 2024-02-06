// accounts_resource.go
package accounts

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
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

// constructJamfProAccount constructs an Account object from the provided schema data.
func constructJamfProAccount(ctx context.Context, d *schema.ResourceData, client *jamfpro.Client) (*jamfpro.ResourceAccount, error) {
	// Initialize the logging subsystem for the construction operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemConstruct, hclog.Debug)

	//func constructJamfProAccount(d *schema.ResourceData) (*jamfpro.ResourceAccount, error) {
	account := &jamfpro.ResourceAccount{}

	// Utilize type assertion helper functions for direct field extraction
	account.Name = util.GetStringFromInterface(d.Get("name"))
	account.DirectoryUser = d.Get("directory_user").(bool)
	account.FullName = util.GetStringFromInterface(d.Get("full_name"))
	account.Email = util.GetStringFromInterface(d.Get("email"))
	account.Enabled = util.GetStringFromInterface(d.Get("enabled"))
	account.ForcePasswordChange = d.Get("force_password_change").(bool)
	account.AccessLevel = util.GetStringFromInterface(d.Get("access_level"))
	account.Password = util.GetStringFromInterface(d.Get("password"))
	account.PrivilegeSet = util.GetStringFromInterface(d.Get("privilege_set"))

	// Construct LDAP Server
	if v, ok := d.GetOk("ldap_server"); ok {
		ldapServerList := v.([]interface{})
		if len(ldapServerList) > 0 && ldapServerList[0] != nil {
			ldapServerMap := ldapServerList[0].(map[string]interface{})
			account.LdapServer = jamfpro.AccountSubsetLdapServer{
				ID:   util.GetIntFromInterface(ldapServerMap["id"]),
				Name: util.GetStringFromInterface(ldapServerMap["name"]),
			}
		}
	}

	// Construct Site
	if v, ok := d.GetOk("site"); ok {
		siteList := v.([]interface{})
		if len(siteList) > 0 && siteList[0] != nil {
			siteMap := siteList[0].(map[string]interface{})
			account.Site = jamfpro.SharedResourceSite{
				ID:   util.GetIntFromInterface(siteMap["id"]),
				Name: util.GetStringFromInterface(siteMap["name"]),
			}
		}
	}

	// Get all accounts to map group names to IDs
	allAccounts, err := client.GetAccounts()
	if err != nil {
		return nil, err
	}

	groupNameToID := make(map[string]int)
	for _, group := range allAccounts.Groups {
		groupNameToID[group.Name] = group.ID
	}

	// Construct Groups with inferred IDs
	if v, ok := d.GetOk("groups"); ok {
		groupsSet := v.(*schema.Set)
		for _, groupItem := range groupsSet.List() {
			groupMap := groupItem.(map[string]interface{})
			groupName := groupMap["name"].(string)
			groupID, exists := groupNameToID[groupName]
			if !exists {
				return nil, fmt.Errorf("group name %s does not exist", groupName)
			}
			group := jamfpro.AccountsListSubsetGroups{
				ID:   groupID,
				Name: groupName,
			}
			// ... Process other fields like 'site' and 'privileges' similarly
			account.Groups = append(account.Groups, group)
		}
	}

	// Construct Privileges
	account.Privileges = jamfpro.AccountSubsetPrivileges{
		JSSObjects:    util.GetStringSliceFromInterface(d.Get("jss_objects_privileges")),
		JSSSettings:   util.GetStringSliceFromInterface(d.Get("jss_settings_privileges")),
		JSSActions:    util.GetStringSliceFromInterface(d.Get("jss_actions_privileges")),
		CasperAdmin:   util.GetStringSliceFromInterface(d.Get("casper_admin_privileges")),
		CasperRemote:  util.GetStringSliceFromInterface(d.Get("casper_remote_privileges")),
		CasperImaging: util.GetStringSliceFromInterface(d.Get("casper_imaging_privileges")),
		Recon:         util.GetStringSliceFromInterface(d.Get("recon_privileges")),
	}

	// Optional: Serialize and pretty-print the accountGroup object for logging
	resourceXML, err := xml.MarshalIndent(account, "", "  ")
	if err != nil {
		logging.LogTFConstructResourceXMLMarshalFailure(subCtx, JamfProResourceAccount, err.Error())
		return nil, err
	}

	// Log the successful construction and serialization to XML
	logging.LogTFConstructedXMLResource(subCtx, JamfProResourceAccount, string(resourceXML))

	return account, nil
}

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
	var creationResponse *jamfpro.ResponseAccountCreatedAndUpdated
	var apiErrorCode int
	resourceName := d.Get("name").(string)

	// Initialize the logging subsystem with the create operation context
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemCreate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Construct the account object outside the retry loop to avoid reconstructing it on each retry
	account, err := constructJamfProAccount(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceAccount, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceAccount)

	// Retry the API call to create the account in Jamf Pro
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateAccount(account)
		if apiErr != nil {
			// Extract and log the API error code if available
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPICreateFailedAfterRetry(subCtx, JamfProResourceAccount, resourceName, apiErr.Error(), apiErrorCode)
			// Return a non-retryable error to break out of the retry loop
			return retry.NonRetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final error and append it to the diagnostics
		logging.LogAPICreateFailure(subCtx, JamfProResourceAccount, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Log successful creation of the accountgroup and set the resource ID in Terraform state
	logging.LogAPICreateSuccess(subCtx, JamfProResourceAccount, strconv.Itoa(creationResponse.ID))

	d.SetId(strconv.Itoa(creationResponse.ID))

	// Retry reading the accountgroup to ensure the Terraform state is up to date
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProAccountRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			// Log any read errors and return a retryable error to retry the read operation
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceAccount, d.Id(), readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		// Successfully read the accountgroup, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final state sync failure and append it to the diagnostics
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceAccount, err.Error())
		diags = append(diags, diag.FromErr(err)...)
	} else {
		// Log successful state synchronization
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceAccount, d.Id())
	}

	return diags
}

// ResourceJamfProAccountRead is responsible for reading the current state of a Jamf Pro Account Group Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the read operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)
	subGeneralCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemGeneral, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	var apiErrorCode int
	var account *jamfpro.ResourceAccount
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subGeneralCtx, "string", "int", JamfProResourceAccount, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Read operation with retry
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		account, apiErr = conn.GetAccountByID(resourceIDInt)
		if apiErr != nil {
			logging.LogFailedReadByID(subCtx, JamfProResourceAccount, resourceID, apiErr.Error(), apiErrorCode)
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the account group, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		d.SetId("") // Remove from Terraform state if unable to read after retries
		logging.LogTFStateRemovalWarning(subCtx, JamfProResourceAccount, resourceID)
		return diag.FromErr(err)
	}

	// Update the Terraform state with account attributes
	d.Set("name", account.Name)
	d.Set("directory_user", account.DirectoryUser)
	d.Set("full_name", account.FullName)
	d.Set("email", account.Email)
	d.Set("enabled", account.Enabled)

	// Update LDAP server information
	if account.LdapServer.ID != 0 || account.LdapServer.Name != "" {
		ldapServer := make(map[string]interface{})
		ldapServer["id"] = account.LdapServer.ID
		ldapServer["name"] = account.LdapServer.Name
		d.Set("ldap_server", []interface{}{ldapServer})
	} else {
		d.Set("ldap_server", []interface{}{}) // Clear the LDAP server data if not present
	}

	d.Set("force_password_change", account.ForcePasswordChange)
	d.Set("access_level", account.AccessLevel)
	// Set password only if it's provided in the configuration
	if _, ok := d.GetOk("password"); ok {
		d.Set("password", account.Password)
	}
	d.Set("privilege_set", account.PrivilegeSet)

	// Update site information
	if account.Site.ID != 0 || account.Site.Name != "" {
		site := make(map[string]interface{})
		site["id"] = account.Site.ID
		site["name"] = account.Site.Name
		d.Set("site", []interface{}{site})
	} else {
		d.Set("site", []interface{}{}) // Clear the site data if not present
	}

	// Construct and set the groups attribute
	groups := make([]interface{}, len(account.Groups))
	for i, group := range account.Groups {
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
	if err := d.Set("jss_objects_privileges", account.Privileges.JSSObjects); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("jss_settings_privileges", account.Privileges.JSSSettings); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("jss_actions_privileges", account.Privileges.JSSActions); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("casper_admin_privileges", account.Privileges.CasperAdmin); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("casper_remote_privileges", account.Privileges.CasperRemote); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("casper_imaging_privileges", account.Privileges.CasperImaging); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("recon_privileges", account.Privileges.Recon); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// ResourceJamfProAccountUpdate is responsible for updating an existing Jamf Pro Account Group on the remote system.
func ResourceJamfProAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the update operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemUpdate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceAccount, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Construct the resource object
	account, err := constructJamfProAccount(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceAccount, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceAccount)

	// Update operations with retries
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateAccountByID(resourceIDInt, account)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.LogAPIUpdateFailureByID(subCtx, JamfProResourceAccount, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			_, apiErrByName := conn.UpdateAccountByName(resourceName, account)
			if apiErrByName != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErrByName.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIUpdateFailureByName(subCtx, JamfProResourceAccount, resourceName, apiErrByName.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErrByName)
			}
		} else {
			logging.LogAPIUpdateSuccess(subCtx, JamfProResourceAccount, resourceID, resourceName)
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceAccount, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Retry reading the Site to synchronize the Terraform state
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProAccountRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceAccount, resourceID, readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceAccount, err.Error())
		return diag.FromErr(err)
	} else {
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceAccount, resourceID)
	}

	return nil
}

// ResourceJamfProAccountDelete is responsible for deleting a Jamf Pro account .
func ResourceJamfProAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the delete operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemDelete, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceAccount, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Delete By ID
		apiErr := conn.DeleteAccountByID(resourceIDInt)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPIDeleteFailureByID(subCtx, JamfProResourceAccount, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			// If Delete by ID fails then try Delete by Name
			apiErr = conn.DeleteAccountByName(resourceName)
			if apiErr != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErr.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIDeleteFailureByName(subCtx, JamfProResourceAccount, resourceName, apiErr.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceAccount, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	logging.LogAPIDeleteSuccess(subCtx, JamfProResourceAccount, resourceID, resourceName)

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return nil
}
