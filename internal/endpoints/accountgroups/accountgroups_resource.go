// accountgroups_resource.go
package accountgroups

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/utilities"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProAccountGroup defines the schema and CRUD operations for managing buildings in Terraform.
func ResourceJamfProAccountGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProAccountGroupCreate,
		ReadContext:   ResourceJamfProAccountGroupRead,
		UpdateContext: ResourceJamfProAccountGroupUpdate,
		DeleteContext: ResourceJamfProAccountGroupDelete,
		CustomizeDiff: customDiffAccountGroups,
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
				Description: "The unique identifier of the account group.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the account group.",
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
							Description: "Jamf Pro Site ID. Value defaults to -1 aka not used.",
							Default:     -1,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Jamf Pro Site Name. Value defaults to 'None' aka not used",
							Computed:    true,
						},
					},
				},
			},
			"jss_objects_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to JSS Objects.",
				Computed:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: common.ValidateJSSObjectsPrivileges,
				},
			},
			"jss_settings_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to JSS Settings.",
				Computed:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: common.ValidateJSSSettingsPrivileges,
				},
			},
			"jss_actions_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to JSS Actions.",
				Computed:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: common.ValidateJSSActionsPrivileges,
				},
			},
			"casper_admin_privileges": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Privileges related to Casper Admin.",
				Computed:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: common.ValidateCasperAdminPrivileges,
				},
			},
			"members": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Members of the account group.",
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
		},
	}
}

const (
	JamfProResourceAccountGroup = "Account Group"
)

// constructJamfProAccountGroup constructs an AccountGroup object from the provided schema data.
func constructJamfProAccountGroup(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceAccountGroup, error) {
	// Initialize the logging subsystem for the construction operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemConstruct, hclog.Debug)

	accountGroup := &jamfpro.ResourceAccountGroup{
		Name:         util.GetStringFromInterface(d.Get("name")),
		AccessLevel:  util.GetStringFromInterface(d.Get("access_level")),
		PrivilegeSet: util.GetStringFromInterface(d.Get("privilege_set")),
	}

	// Construct Site
	if v, ok := d.GetOk("site"); ok && len(v.([]interface{})) > 0 {
		siteMap := v.([]interface{})[0].(map[string]interface{})
		accountGroup.Site = jamfpro.SharedResourceSite{
			ID:   util.GetIntFromInterface(siteMap["id"]),
			Name: util.GetStringFromInterface(siteMap["name"]),
		}
	}

	// Construct Privileges using TypeSet
	accountGroup.Privileges = jamfpro.AccountSubsetPrivileges{
		JSSObjects:  utilities.ExtractSetToStringSlice(d.Get("jss_objects_privileges").(*schema.Set)),
		JSSSettings: utilities.ExtractSetToStringSlice(d.Get("jss_settings_privileges").(*schema.Set)),
		JSSActions:  utilities.ExtractSetToStringSlice(d.Get("jss_actions_privileges").(*schema.Set)),
		CasperAdmin: utilities.ExtractSetToStringSlice(d.Get("casper_admin_privileges").(*schema.Set)),
	}

	// Construct Members
	if v, ok := d.GetOk("members"); ok {
		members := make(jamfpro.AccountGroupSubsetMembers, 0)
		for _, member := range v.([]interface{}) {
			memberMap := member.(map[string]interface{})
			memberUser := jamfpro.MemberUser{
				ID:   util.GetIntFromInterface(memberMap["id"]),
				Name: util.GetStringFromInterface(memberMap["name"]),
			}
			members = append(members, struct {
				User jamfpro.MemberUser `json:"user,omitempty" xml:"user,omitempty"`
			}{
				User: memberUser,
			})
		}
		accountGroup.Members = members
	}

	// Optional: Serialize and pretty-print the accountGroup object for logging
	resourceXML, err := xml.MarshalIndent(accountGroup, "", "  ")
	if err != nil {
		logging.LogTFConstructResourceXMLMarshalFailure(subCtx, JamfProResourceAccountGroup, err.Error())
		return nil, err
	}

	// Log the successful construction and serialization to XML
	logging.LogTFConstructedXMLResource(subCtx, JamfProResourceAccountGroup, string(resourceXML))

	return accountGroup, nil
}

// ResourceJamfProAccountGroupCreate is responsible for creating a new Jamf Pro Script in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
// ResourceJamfProAccountGroupCreate is responsible for creating a new Jamf Pro Account Group in the remote system.
func ResourceJamfProAccountGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	var creationResponse *jamfpro.ResponseAccountGroupCreated
	var apiErrorCode int
	resourceName := d.Get("name").(string)

	// Initialize the logging subsystem with the create operation context
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemCreate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Construct the accountgroup object outside the retry loop to avoid reconstructing it on each retry
	accountGroup, err := constructJamfProAccountGroup(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceAccountGroup, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceAccountGroup)

	// Retry the API call to create the accountgroup in Jamf Pro
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateAccountGroup(accountGroup)
		if apiErr != nil {
			// Extract and log the API error code if available
			if apiError, ok := apiErr.(*.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPICreateFailedAfterRetry(subCtx, JamfProResourceAccountGroup, resourceName, apiErr.Error(), apiErrorCode)
			// Return a non-retryable error to break out of the retry loop
			return retry.NonRetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final error and append it to the diagnostics
		logging.LogAPICreateFailure(subCtx, JamfProResourceAccountGroup, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Log successful creation of the accountgroup and set the resource ID in Terraform state
	logging.LogAPICreateSuccess(subCtx, JamfProResourceAccountGroup, strconv.Itoa(creationResponse.ID))

	d.SetId(strconv.Itoa(creationResponse.ID))

	// Retry reading the accountgroup to ensure the Terraform state is up to date
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProAccountGroupRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			// Log any read errors and return a retryable error to retry the read operation
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceAccountGroup, d.Id(), readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		// Successfully read the accountgroup, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final state sync failure and append it to the diagnostics
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceAccountGroup, err.Error())
		diags = append(diags, diag.FromErr(err)...)
	} else {
		// Log successful state synchronization
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceAccountGroup, d.Id())
	}

	return diags
}

// ResourceJamfProAccountGroupRead is responsible for reading the current state of a Jamf Pro Account Group Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProAccountGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var accountGroup *jamfpro.ResourceAccountGroup
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subGeneralCtx, "string", "int", JamfProResourceAccountGroup, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Read operation with retry
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		accountGroup, apiErr = conn.GetAccountGroupByID(resourceIDInt)
		if apiErr != nil {
			logging.LogFailedReadByID(subCtx, JamfProResourceAccountGroup, resourceID, apiErr.Error(), apiErrorCode)
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the account group, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		d.SetId("") // Remove from Terraform state if unable to read after retries
		logging.LogTFStateRemovalWarning(subCtx, JamfProResourceAccountGroup, resourceID)
		return diag.FromErr(err)
	}

	// Assuming successful read if no error
	logging.LogAPIReadSuccess(subCtx, JamfProResourceAccountGroup, resourceID)

	if err := d.Set("id", resourceID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("name", accountGroup.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("access_level", accountGroup.AccessLevel); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("privilege_set", accountGroup.PrivilegeSet); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Update site information
	site := make(map[string]interface{})
	site["id"] = accountGroup.Site.ID
	site["name"] = accountGroup.Site.Name
	if err := d.Set("site", []interface{}{site}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set privileges
	privilegeAttributes := map[string][]string{
		"jss_objects_privileges":  accountGroup.Privileges.JSSObjects,
		"jss_settings_privileges": accountGroup.Privileges.JSSSettings,
		"jss_actions_privileges":  accountGroup.Privileges.JSSActions,
		"casper_admin_privileges": accountGroup.Privileges.CasperAdmin,
	}

	for attrName, privileges := range privilegeAttributes {
		if err := d.Set(attrName, schema.NewSet(schema.HashString, utilities.ConvertToStringInterface(privileges))); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Update members
	members := make([]interface{}, 0)
	for _, memberStruct := range accountGroup.Members {
		member := memberStruct.User // Access the User field
		memberMap := map[string]interface{}{
			"id":   member.ID,
			"name": member.Name,
		}
		members = append(members, memberMap)
	}
	if err := d.Set("members", members); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Check if there were any errors and return the diagnostics
	if len(diags) > 0 {
		return diags
	}
	return nil
}

// ResourceJamfProAccountGroupUpdate is responsible for updating an existing Jamf Pro Account Group on the remote system.
func ResourceJamfProAccountGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceAccountGroup, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Construct the resource object
	dockItem, err := constructJamfProAccountGroup(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceAccountGroup, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceAccountGroup)

	// Update operations with retries
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateAccountGroupByID(resourceIDInt, dockItem)
		if apiErr != nil {
			if apiError, ok := apiErr.(*.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.LogAPIUpdateFailureByID(subCtx, JamfProResourceAccountGroup, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			_, apiErrByName := conn.UpdateAccountGroupByName(resourceName, dockItem)
			if apiErrByName != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErrByName.(*.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIUpdateFailureByName(subCtx, JamfProResourceAccountGroup, resourceName, apiErrByName.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErrByName)
			}
		} else {
			logging.LogAPIUpdateSuccess(subCtx, JamfProResourceAccountGroup, resourceID, resourceName)
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceAccountGroup, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Retry reading the Site to synchronize the Terraform state
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProAccountGroupRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceAccountGroup, resourceID, readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceAccountGroup, err.Error())
		return diag.FromErr(err)
	} else {
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceAccountGroup, resourceID)
	}

	return nil
}

// ResourceJamfProAccountGroupDelete is responsible for deleting a Jamf Pro account group.
func ResourceJamfProAccountGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceAccountGroup, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Delete By ID
		apiErr := conn.DeleteAccountGroupByID(resourceIDInt)
		if apiErr != nil {
			if apiError, ok := apiErr.(*.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPIDeleteFailureByID(subCtx, JamfProResourceAccountGroup, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			// If Delete by ID fails then try Delete by Name
			apiErr = conn.DeleteAccountGroupByName(resourceName)
			if apiErr != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErr.(*.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIDeleteFailureByName(subCtx, JamfProResourceAccountGroup, resourceName, apiErr.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceAccountGroup, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	logging.LogAPIDeleteSuccess(subCtx, JamfProResourceAccountGroup, resourceID, resourceName)

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return nil
}
