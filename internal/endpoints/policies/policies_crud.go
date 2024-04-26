package policies

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Constructs, creates states
func ResourceJamfProPoliciesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics

	resource, err := constructPolicy(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Policy: %v", err))
	}

	// Retry the API call to create the policy in Jamf Pro
	var creationResponse *jamfpro.ResponsePolicyCreateAndUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreatePolicy(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Policy '%s' after retries: %v", resource.General.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	readDiags := ResourceJamfProPoliciesRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// Reads and states
func ResourceJamfProPoliciesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var resp *jamfpro.ResourcePolicy

	// Extract policy name from schema

	// Use the retry function for the read operation
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resp, apiErr = conn.GetPolicyByID(resourceIDInt)
		if apiErr != nil {
			if strings.Contains(apiErr.Error(), "404") || strings.Contains(apiErr.Error(), "410") {
				return retry.NonRetryableError(fmt.Errorf("resource not found, marked for deletion"))
			}
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		d.SetId("") // Remove from Terraform state if unable to read after retries
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Policy '%s' (ID: %d) after retries: %v", "no", resourceIDInt, err))
	}

	err = d.Set("name", resp.General.Name)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("enabled", resp.General.Enabled)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_checkin", resp.General.TriggerCheckin)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_enrollment_complete", resp.General.TriggerEnrollmentComplete)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_login", resp.General.TriggerLogin)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_network_state_changed", resp.General.TriggerNetworkStateChanged)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_startup", resp.General.TriggerStartup)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("trigger_other", resp.General.TriggerOther)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("frequency", resp.General.Frequency)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("retry_event", resp.General.RetryEvent)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("retry_attempts", resp.General.RetryAttempts)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("notify_on_each_failed_retry", resp.General.NotifyOnEachFailedRetry)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	err = d.Set("offline", resp.General.Offline)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Site
	if resp.General.Site.ID != -1 && resp.General.Site.Name != "None" {
		out_site := []map[string]interface{}{
			{
				"id": resp.General.Site.ID,
				// "name": resp.General.Site.Name,
			},
		}

		if err := d.Set("site", out_site); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		log.Println("Not stating default site response") // TODO logging
	}

	// Category
	if resp.General.Category.ID != -1 && resp.General.Category.Name != "No category assigned" {
		out_category := []map[string]interface{}{
			{
				"id": resp.General.Category.ID,
				// "name": resp.General.Category.Name,
			},
		}
		if err := d.Set("category", out_category); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		log.Println("Not stating default category response") // TODO logging
	}

	log.Println(resp)

	return diags
}

// Constructs, updates and reads
func ResourceJamfProPoliciesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceIDInt, err := strconv.Atoi(resourceID)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructPolicy(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Policy for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdatePolicyByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Policy '%s' (ID: %d) after retries: %v", resource.General.Name, resourceIDInt, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProPoliciesRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// Deletes and removes from state
func ResourceJamfProPoliciesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resourceName := d.Get("name").(string)

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := conn.DeletePolicyByID(resourceIDInt)
		if apiErr != nil {

			apiErrByName := conn.DeletePolicyByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the site, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro policy '%s' (ID: %s) after retries: %v", resourceName, d.Id(), err))
	}

	d.SetId("")

	return diags
}
