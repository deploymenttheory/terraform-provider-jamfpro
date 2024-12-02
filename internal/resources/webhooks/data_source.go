// webhooks_data_source.go
package webhooks

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProWebhooks provides information about a specific Jamf Pro Webhook by its ID or Name.
func DataSourceJamfProWebhooks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the Jamf Pro Webhook.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name of the Jamf Pro Webhook.",
			},
		},
	}
}

// dataSourceRead fetches the details of a specific Jamf Pro Webhook
// from Jamf Pro using either its unique Name or its Id. The function prioritizes the 'name' attribute over the 'id'
// attribute for fetching details. If neither 'name' nor 'id' is provided, it returns an error.
// Once the details are fetched, they are set in the data source's state.
//
// Parameters:
// - ctx: The context within which the function is called. It's used for timeouts and cancellation.
// - d: The current state of the data source.
// - meta: The meta object that can be used to retrieve the API client connection.
//
// Returns:
// - diag.Diagnostics: Returns any diagnostics (errors or warnings) encountered during the function's execution.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	var diags diag.Diagnostics
	resourceID := d.Get("id").(string)
	var resource *jamfpro.ResourceWebhook

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = client.GetWebhookByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Webhook with ID '%s' after retries: %v", resourceID, err))
	}

	if resource != nil {
		d.SetId(resourceID)
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(fmt.Errorf("error setting 'name' for Jamf Pro Webhook with ID '%s': %v", resourceID, err))...)
		}
	} else {
		d.SetId("")
	}

	return diags
}

// DataSourceJamfProWebhooks List provides a list of all Jamf Pro Webhooks
func DataSourceJamfProWebhooksList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReadList,
		Schema: map[string]*schema.Schema{
			"webhooks": {
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

// dataSourceReadList retrieves a list of all Jamf Pro Webhooks
// and maps the data into the Terraform state.
//
// Parameters:
// - ctx: The context within which the function is called. Used for timeouts and cancellation.
// - d: The current state of the data source.
// - meta: The meta object that provides the API client connection.
//
// Returns:
// - diag.Diagnostics: Returns any diagnostics (errors or warnings) encountered during the function's execution.
func dataSourceReadList(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	var diags diag.Diagnostics

	response, err := client.GetWebhooks()
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve Jamf Pro Webhooks: %v", err))
	}

	var webhooks []map[string]interface{}
	var ids []string

	for _, webhook := range response.Webhooks {
		webhooks = append(webhooks, map[string]interface{}{
			"id":   strconv.Itoa(webhook.ID),
			"name": webhook.Name,
		})
		ids = append(ids, strconv.Itoa(webhook.ID))
	}

	if err := d.Set("webhooks", webhooks); err != nil {
		return diag.FromErr(fmt.Errorf("error setting 'webhooks' attribute: %v", err))
	}
	if err := d.Set("ids", ids); err != nil {
		return diag.FromErr(fmt.Errorf("error setting 'ids' attribute: %v", err))
	}

	d.SetId(fmt.Sprintf("datasource-webhooks-list-%d", time.Now().Unix()))

	return diags
}
