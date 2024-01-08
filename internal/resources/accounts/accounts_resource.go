// accounts_resource.go
package accounts

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProAccount defines the schema and CRUD operations for managing buildings in Terraform.
func ResourceJamfProAccountGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProAccountCreate,
		ReadContext:   ResourceJamfProAccountRead,
		UpdateContext: ResourceJamfProAccountUpdate,
		DeleteContext: ResourceJamfProAccountDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the account.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the account.",
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
			"enabled": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The enabled status of the account.",
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
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the LDAP server.",
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
				Optional:    true,
				Description: "The access level of the account.",
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
			},
			"site": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The site information associated with the account group.",
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
			"privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The privileges associated with the account.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// constructJamfProAccount constructs an Account object from the provided schema data.
func constructJamfProAccount(d *schema.ResourceData) (*jamfpro.ResourceAccount, error) {
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

	// Construct Privileges
	// This needs to be adjusted based on how privileges are represented in your schema
	if v, ok := d.GetOk("privileges"); ok {
		privilegesList := v.([]interface{})
		var privileges []string
		for _, priv := range privilegesList {
			privileges = append(privileges, util.GetStringFromInterface(priv))
		}
		account.Privileges = jamfpro.AccountSubsetPrivileges{
			JSSObjects:    privileges,
			JSSSettings:   privileges,
			JSSActions:    privileges,
			Recon:         privileges,
			CasperAdmin:   privileges,
			CasperRemote:  privileges,
			CasperImaging: privileges,
		}
	}

	// Log the successful construction of the account
	log.Printf("[INFO] Successfully constructed Account with name: %s", account.Name)

	return account, nil
}

// Helper function to generate diagnostics based on the error type.
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

// ResourceJamfProAccountCreate is responsible for creating a new Jamf Pro Script in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Use the retry function for the create operation
	var createdAccount *jamfpro.ResponseAccountCreatedAndUpdated
	var err error
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		// Construct the account
		account, err := constructJamfProAccount(d)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to construct the account for terraform create: %w", err))
		}

		// Directly call the API to create the resource
		createdAccount, err = conn.CreateAccount(account)
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
	d.SetId(strconv.Itoa(createdAccount.ID))

	// Use the retry function for the read operation to update the Terraform state with the resource attributes
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProAccountRead(ctx, d, meta)
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

// ResourceJamfProAccountRead is responsible for reading the current state of a Jamf Pro Account Group Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var accountGroup *jamfpro.ResourceAccount

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		accountGroupID, err := strconv.Atoi(d.Id())
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("error converting id (%s) to integer: %s", d.Id(), err))
		}

		// Try fetching the account group using the ID
		accountGroup, err = conn.GetAccountByID(accountGroupID)
		if err != nil {
			// Handle the APIError
			if apiError, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			return retry.RetryableError(err)
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while reading the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	// Update the Terraform state with account group attributes
	d.Set("name", accountGroup.Name)
	d.Set("access_level", accountGroup.AccessLevel)
	d.Set("privilege_set", accountGroup.PrivilegeSet)

	// Update site information
	site := make(map[string]interface{})
	site["id"] = accountGroup.Site.ID
	site["name"] = accountGroup.Site.Name
	d.Set("site", []interface{}{site})

	// Update privileges
	privileges := make(map[string]interface{})
	privileges["jss_objects"] = accountGroup.Privileges.JSSObjects
	privileges["jss_settings"] = accountGroup.Privileges.JSSSettings
	privileges["jss_actions"] = accountGroup.Privileges.JSSActions
	privileges["recon"] = accountGroup.Privileges.Recon
	privileges["casper_admin"] = accountGroup.Privileges.CasperAdmin
	privileges["casper_remote"] = accountGroup.Privileges.CasperRemote
	privileges["casper_imaging"] = accountGroup.Privileges.CasperImaging
	d.Set("privileges", []interface{}{privileges})

	// Update members
	members := make([]interface{}, 0)
	for _, member := range accountGroup.Members {
		memberMap := map[string]interface{}{
			"id":   member.ID,
			"name": member.Name,
		}
		members = append(members, memberMap)
	}
	d.Set("members", members)

	return diags
}
