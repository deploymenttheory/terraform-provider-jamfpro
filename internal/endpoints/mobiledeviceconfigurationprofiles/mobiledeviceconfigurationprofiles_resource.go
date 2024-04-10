// mobiledeviceconfigurationprofiles_resource.go
package mobiledeviceconfigurationprofiles

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/waitfor"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/utilities"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProMobileDeviceConfigurationProfile defines the schema for mobile device configuration profiles in Terraform.
func ResourceJamfProMobileDeviceConfigurationProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProMobileDeviceConfigurationProfileCreate,
		ReadContext:   ResourceJamfProMobileDeviceConfigurationProfileRead,
		UpdateContext: ResourceJamfProMobileDeviceConfigurationProfileUpdate,
		DeleteContext: ResourceJamfProMobileDeviceConfigurationProfileDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier for the mobile device configuration profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the mobile device configuration profile.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the mobile device configuration profile.",
			},
			// Additional attributes as needed by your provider...
		},
	}
}

// ResourceJamfProMobileDeviceConfigurationProfileCreate is responsible for creating a new Jamf Pro Mobile Device Configuration Profile in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
// ResourceJamfProMobileDeviceConfigurationProfileCreate is responsible for creating a new Jamf Pro Mobile Device Configuration Profile in the remote system.
func ResourceJamfProMobileDeviceConfigurationProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProMobileDeviceConfigurationProfile(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Mobile Device Configuration Profile: %v", err))
	}

	// Retry the API call to create the resource in Jamf Pro
	var creationResponse *jamfpro.ResponseMobileDeviceConfigurationProfileCreateAndUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateMobileDeviceConfigurationProfile(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Mobile Device Configuration Profile '%s' after retries: %v", resource.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Wait for the resource to be fully available before reading it
	checkResourceExists := func(id interface{}) (interface{}, error) {
		intID, err := strconv.Atoi(id.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
		}
		return apiclient.Conn.GetMobileDeviceConfigurationProfileByID(intID)
	}

	_, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Mobile Device Configuration Profile", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second, apiclient.EnableCookieJar)

	if waitDiags.HasError() {
		return waitDiags
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProMobileDeviceConfigurationProfileRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProMobileDeviceConfigurationProfileRead is responsible for reading the current state of a Jamf Pro Mobile Device Configuration Profile Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProMobileDeviceConfigurationProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Attempt to fetch the resource by ID
	resource, err := conn.GetMobileDeviceConfigurationProfileByID(resourceIDInt)

	if err != nil {
		// Skip resource state removal if this is a create operation
		if !d.IsNewResource() {
			// If the error is a "not found" error, remove the resource from the state
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "410") {
				d.SetId("") // Remove the resource from Terraform state
				return diag.Diagnostics{
					{
						Severity: diag.Warning,
						Summary:  "Resource not found",
						Detail:   fmt.Sprintf("Jamf Pro Mobile Device Configuration Profile resource with ID '%s' was not found and has been removed from the Terraform state.", resourceID),
					},
				}
			}
		}
		// For other errors, or if this is a create operation, return a diagnostic error
		return diag.FromErr(err)
	}

	// Update the Terraform state with the fetched data
	if err := d.Set("id", resourceID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("name", resource.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("access_level", resource.AccessLevel); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("privilege_set", resource.PrivilegeSet); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Update LDAP server information
	if resource.LDAPServer.ID != 0 {
		ldapServer := make(map[string]interface{})
		ldapServer["id"] = resource.LDAPServer.ID
		d.Set("identity_server", []interface{}{ldapServer})
	} else {
		d.Set("identity_server", []interface{}{}) // Clear the LDAP server data if not present
	}

	// Update site information
	site := make(map[string]interface{})
	site["id"] = resource.Site.ID
	site["name"] = resource.Site.Name
	if err := d.Set("site", []interface{}{site}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Set privileges
	privilegeAttributes := map[string][]string{
		"jss_objects_privileges":  resource.Privileges.JSSObjects,
		"jss_settings_privileges": resource.Privileges.JSSSettings,
		"jss_actions_privileges":  resource.Privileges.JSSActions,
		"casper_admin_privileges": resource.Privileges.CasperAdmin,
	}

	for attrName, privileges := range privilegeAttributes {
		if err := d.Set(attrName, schema.NewSet(schema.HashString, utilities.ConvertToStringInterface(privileges))); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	// Update members
	members := make([]interface{}, 0)
	for _, memberStruct := range resource.Members {
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

// ResourceJamfProMobileDeviceConfigurationProfileUpdate is responsible for updating an existing Jamf Pro Mobile Device Configuration Profile on the remote system.
func ResourceJamfProMobileDeviceConfigurationProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	resource, err := constructJamfProMobileDeviceConfigurationProfile(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Mobile Device Configuration Profile for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateMobileDeviceConfigurationProfileByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Mobile Device Configuration Profile '%s' (ID: %s) after retries: %v", resource.Name, resourceID, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProMobileDeviceConfigurationProfileRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProMobileDeviceConfigurationProfileDelete is responsible for deleting a Jamf Pro Mobile Device Configuration Profile.
func ResourceJamfProMobileDeviceConfigurationProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		apiErr := conn.DeleteMobileDeviceConfigurationProfileByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeleteMobileDeviceConfigurationProfileByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Mobile Device Configuration Profile '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
