// scripts_data_source.go
package scripts

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProScripts provides information about a specific Jamf Pro script by its ID or Name.
func DataSourceJamfProScripts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Jamf Pro unique identifier (ID) of the script.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display name for the script.",
			},
		},
	}
}

// dataSourceRead fetches the details of a specific Jamf Pro script
// from Jamf Pro using either its unique Name or its Id.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)

	var resource *jamfpro.ResourceScript

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = client.GetScriptByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Script with ID '%s' after retries: %v", resourceID, err))
	}

	if resource != nil {
		d.SetId(resourceID)
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Script with ID '%s': %v", resourceID, err))...)
		}
	} else {
		d.SetId("")
	}

	return diags
}

// DataSourceJamfProScriptsList provides a list of all Jamf Pro scripts.
func DataSourceJamfProScriptsList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadList,
		Schema: map[string]*schema.Schema{
			"scripts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// dataSourceReadList fetches a list of all Jamf Pro scripts and maps them to the Terraform state.
func dataSourceReadList(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *jamfpro.Client")
	}

	var diags diag.Diagnostics

	// Fetch the list of scripts
	response, err := client.GetScripts("")
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve Jamf Pro scripts: %v", err))
	}

	// Map the scripts to the Terraform state
	var scripts []map[string]interface{}
	var ids []string

	for _, script := range response.Results {
		scripts = append(scripts, map[string]interface{}{
			"id":            script.ID,
			"name":          script.Name,
			"category_name": script.CategoryName,
		})
		ids = append(ids, script.ID)
	}

	// Set the computed attributes in Terraform state
	if err := d.Set("scripts", scripts); err != nil {
		return diag.FromErr(fmt.Errorf("error setting 'scripts' attribute: %v", err))
	}
	if err := d.Set("ids", ids); err != nil {
		return diag.FromErr(fmt.Errorf("error setting 'ids' attribute: %v", err))
	}

	// Generate a unique ID for the resource
	d.SetId(fmt.Sprintf("datasource-scripts-list-%d", time.Now().Unix()))

	return diags
}
