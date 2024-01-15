// accounts_data_source.go
package accounts

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProAccounts provides information about specific Jamf Pro Dock Items by their ID or Name.
func DataSourceJamfProAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJamfProAccountRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the jamf pro account.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the jamf pro account.",
			},
			"directory_user": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the user is a directory user.",
			},
			"full_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full name of the account user.",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email of the account user.",
			},
			"enabled": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access status of the account (“enabled” or “disabled”).",
			},
			"ldap_server": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Description: "LDAP server information associated with the account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the LDAP server.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the LDAP server.",
						},
					},
				},
			},
			"force_password_change": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the user is forced to change password on next login.",
			},
			"access_level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The access level of the account. This can be either Full Access, scoped to a jamf pro site with Site Access, or scoped to a jamf pro account group with Group Access",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The password for the account.",
			},
			"privilege_set": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The privilege set assigned to the account.",
			},
			"site": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Description: "The site information associated with the account group if access_level is set to Site Access.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Jamf Pro Site ID. Value defaults to '0' aka not used.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Jamf Pro Site Name",
						},
					},
				},
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of group information associated with the account.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the group.",
						},
						"site": {
							Type:        schema.TypeList,
							Computed:    true,
							MaxItems:    1,
							Description: "The site information associated with the group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Jamf Pro Site ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Jamf Pro Site Name.",
									},
								},
							},
						},
						"privileges": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The privileges assigned to the group.",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"jss_objects": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Privileges related to JSS Objects.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"jss_settings": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Privileges related to JSS Settings.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"jss_actions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Privileges related to JSS Actions.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"jss_objects_privileges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Privileges related to JSS Objects.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"jss_settings_privileges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Privileges related to JSS Settings.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"jss_actions_privileges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Privileges related to JSS Actions.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"casper_admin_privileges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Privileges related to Casper Admin.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"casper_remote_privileges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Privileges related to Casper Remote.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"casper_imaging_privileges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Privileges related to Casper Imaging.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"recon_privileges": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Privileges related to Recon.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// dataSourceJamfProAccountRead fetches the details of specific account from Jamf Pro using either their unique Name or Id.
func dataSourceJamfProAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Asserts 'meta' as '*client.APIClient'
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var account *jamfpro.ResourceAccount

	// Use the retry function for the read operation
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		// Convert the ID from the Terraform state into an integer to be used for the API request
		accountID, err := strconv.Atoi(d.Id())
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("error converting id (%s) to integer: %s", d.Id(), err))
		}

		// Try fetching the account using the ID
		account, err = conn.GetAccountByID(accountID)
		if err != nil {
			// Handle the APIError
			if apiError, ok := err.(*http_client.APIError); ok {
				return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
			}
			// If fetching by ID fails, try fetching by Name
			accountName, ok := d.Get("name").(string)
			if !ok {
				return retry.NonRetryableError(fmt.Errorf("unable to assert 'name' as a string"))
			}

			account, err = conn.GetAccountByName(accountName)
			if err != nil {
				// Handle the APIError
				if apiError, ok := err.(*http_client.APIError); ok {
					return retry.NonRetryableError(fmt.Errorf("API Error (Code: %d): %s", apiError.StatusCode, apiError.Message))
				}
				return retry.RetryableError(err)
			}
		}
		return nil
	})

	// Handle error from the retry function
	if err != nil {
		// If there's an error while reading the resource, generate diagnostics using the helper function.
		return generateTFDiagsFromHTTPError(err, d, "read")
	}

	// Update the Terraform state with account attributes
	d.Set("name", account.Name)
	d.Set("directory_user", account.DirectoryUser)
	d.Set("full_name", account.FullName)
	d.Set("email", account.Email)
	d.Set("enabled", account.Enabled)

	// Update LDAP server information
	ldapServer := make(map[string]interface{})
	ldapServer["id"] = account.LdapServer.ID
	ldapServer["name"] = account.LdapServer.Name
	d.Set("ldap_server", []interface{}{ldapServer})

	d.Set("force_password_change", account.ForcePasswordChange)
	d.Set("access_level", account.AccessLevel)
	d.Set("password", account.Password)
	d.Set("privilege_set", account.PrivilegeSet)

	// Update site information
	site := make(map[string]interface{})
	site["id"] = account.Site.ID
	site["name"] = account.Site.Name
	d.Set("site", []interface{}{site})

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
