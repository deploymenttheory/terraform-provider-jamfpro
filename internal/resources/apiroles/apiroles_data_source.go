// apiroles_data_source.go
package apiroles

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProAPIRoles provides information about a specific Jamf Pro API role by its ID or Name.
func DataSourceJamfProAPIRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJamfProAPIRolesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The unique identifier of the API role.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique name of the Jamf Pro API role.",
				Computed:    true,
			},
			"privileges": {
				Type:        schema.TypeList,
				Description: "List of privileges associated with the API role.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// dataSourceJamfProAPIRolesRead fetches the details of a specific API role from Jamf Pro using either its unique Name or its Id.
// The function prioritizes the 'name' attribute over the 'id' attribute for fetching details. If neither 'name' nor 'id' is provided,
// it returns an error. Once the details are fetched, they are set in the data source's state.
func dataSourceJamfProAPIRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*client.APIClient).Conn

	var role *jamfpro.APIRole
	var err error

	// Check if Name is provided in the data source configuration
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		roleName := v.(string)
		role, err = conn.GetJamfApiRolesNameById(roleName)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch API role by name: %v", err))
		}
	} else if v, ok := d.GetOk("id"); ok {
		roleID := v.(int) // Correctly cast to int
		role, err = conn.GetJamfApiRolesByID(roleID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch API role by ID: %v", err))
		}
	} else {
		return diag.Errorf("Either 'name' or 'id' must be provided")
	}

	// Set the data source attributes using the fetched data
	d.SetId(role.ID)
	d.Set("name", role.DisplayName)
	d.Set("privileges", role.Privileges)

	return nil
}
