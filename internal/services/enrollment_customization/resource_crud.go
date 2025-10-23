package enrollment_customization

import (
	"context"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for creating a new enrollment customization
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Handle image upload first if image source is provided
	imagePath, err := constructImageUpload(d)
	if err == nil && imagePath != "" {
		err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			uploadResponse, apiErr := client.UploadEnrollmentCustomizationsImage(imagePath)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}

			brandingSettings := d.Get("branding_settings").([]interface{})
			if len(brandingSettings) > 0 {
				settings := brandingSettings[0].(map[string]interface{})
				settings["icon_url"] = uploadResponse.Url
				brandingSettingsList := []interface{}{settings}
				if err := d.Set("branding_settings", brandingSettingsList); err != nil {
					return retry.NonRetryableError(fmt.Errorf("failed to set icon_url in schema: %v", err))
				}
			}
			return nil
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to upload enrollment customization image: %v", err))
		}
	}

	resource, err := constructBaseResource(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct enrollment customization: %v", err))
	}

	var response *jamfpro.ResponseEnrollmentCustomizationCreate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.CreateEnrollmentCustomization(*resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create enrollment customization: %v", err))
	}

	d.SetId(response.Id)

	textPanes := d.Get("text_pane").([]interface{})
	ldapPanes := d.Get("ldap_pane").([]interface{})
	ssoPanes := d.Get("sso_pane").([]interface{})

	// Check for valid combinations
	hasSSO := len(ssoPanes) > 0
	hasLDAP := len(ldapPanes) > 0
	hasText := len(textPanes) > 0

	if hasSSO && hasLDAP {
		return diag.FromErr(fmt.Errorf("invalid combination: SSO and LDAP panes cannot be used together"))
	}

	// Create text panes if present
	if hasText {
		for _, paneData := range textPanes {
			textPane, err := constructTextPane(paneData.(map[string]interface{}))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to construct text pane: %v", err))
			}

			err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
				_, apiErr := client.CreateTextPrestagePane(response.Id, *textPane)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to create text pane: %v", err))
			}
		}
	}

	// Create LDAP panes if present
	if hasLDAP {
		for _, paneData := range ldapPanes {
			ldapPane, err := constructLDAPPane(paneData.(map[string]interface{}))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to construct LDAP pane: %v", err))
			}

			err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
				_, apiErr := client.CreateLDAPPrestagePane(response.Id, *ldapPane)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to create LDAP pane: %v", err))
			}
		}
	}

	// Create SSO panes if present
	if hasSSO {
		for _, paneData := range ssoPanes {
			ssoPane, err := constructSSOPane(paneData.(map[string]interface{}))
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to construct SSO pane: %v", err))
			}

			err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
				_, apiErr := client.CreateSSOPrestagePane(response.Id, *ssoPane)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to create SSO pane: %v", err))
			}
		}
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the enrollment customization
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	var response *jamfpro.ResourceEnrollmentCustomization
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetEnrollmentCustomizationByID(d.Id())
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	// Update state with the base resource information
	if err := updateState(d, response); err != nil {
		return append(diags, err...)
	}

	// Get all prestage panes for this customization
	var panesList *jamfpro.ResponsePrestagePanesList
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		panesList, apiErr = client.GetPrestagePanes(d.Id())
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get prestage panes for enrollment customization: %v", err))
	}

	// Process and update state for each pane type
	readTextPanes(ctx, d, client, panesList)
	readLDAPPanes(ctx, d, client, panesList)
	readSSOPanes(ctx, d, client, panesList)

	return diags
}

// readTextPanes fetches text panes and updates state
func readTextPanes(ctx context.Context, d *schema.ResourceData, client *jamfpro.Client, panesList *jamfpro.ResponsePrestagePanesList) error {
	var textPanes []map[string]interface{}

	for _, panel := range panesList.Panels {
		if panel.Type != "text" {
			continue
		}

		paneID := fmt.Sprintf("%d", panel.ID)
		var textPane *jamfpro.ResourceEnrollmentCustomizationTextPane

		err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
			var apiErr error
			textPane, apiErr = client.GetTextPrestagePaneByID(d.Id(), paneID)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			log.Printf("[WARNING] failed to get text pane %s: %v", paneID, err)
			continue
		}

		textPaneMap := stateTextPrestagePane(textPane)
		textPanes = append(textPanes, textPaneMap)
	}

	if len(textPanes) > 0 {
		if err := d.Set("text_pane", textPanes); err != nil {
			return fmt.Errorf("error setting text_pane: %v", err)
		}
	}

	return nil
}

// readLDAPPanes fetches LDAP panes and updates state
func readLDAPPanes(ctx context.Context, d *schema.ResourceData, client *jamfpro.Client, panesList *jamfpro.ResponsePrestagePanesList) error {
	var ldapPanes []map[string]interface{}

	for _, panel := range panesList.Panels {
		if panel.Type != "ldap" {
			continue
		}

		paneID := fmt.Sprintf("%d", panel.ID)
		var ldapPane *jamfpro.ResourceEnrollmentCustomizationLDAPPane

		err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
			var apiErr error
			ldapPane, apiErr = client.GetLDAPPrestagePaneByID(d.Id(), paneID)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			log.Printf("[WARNING] failed to get LDAP pane %s: %v", paneID, err)
			continue
		}

		ldapPaneMap := stateLDAPPrestagePane(ldapPane)
		ldapPanes = append(ldapPanes, ldapPaneMap)
	}

	if len(ldapPanes) > 0 {
		if err := d.Set("ldap_pane", ldapPanes); err != nil {
			return fmt.Errorf("error setting ldap_pane: %v", err)
		}
	}

	return nil
}

// readSSOPanes fetches SSO panes and updates state
func readSSOPanes(ctx context.Context, d *schema.ResourceData, client *jamfpro.Client, panesList *jamfpro.ResponsePrestagePanesList) error {
	var ssoPanes []map[string]interface{}

	for _, panel := range panesList.Panels {
		if panel.Type != "sso" {
			continue
		}

		paneID := fmt.Sprintf("%d", panel.ID)
		var ssoPane *jamfpro.ResourceEnrollmentCustomizationSSOPane

		err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
			var apiErr error
			ssoPane, apiErr = client.GetSSOPrestagePaneByID(d.Id(), paneID)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			log.Printf("[WARNING] failed to get SSO pane %s: %v", paneID, err)
			continue
		}

		ssoPaneMap := stateSSOPrestagePane(ssoPane)
		ssoPanes = append(ssoPanes, ssoPaneMap)
	}

	if len(ssoPanes) > 0 {
		if err := d.Set("sso_pane", ssoPanes); err != nil {
			return fmt.Errorf("error setting sso_pane: %v", err)
		}
	}

	return nil
}

func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update is responsible for updating the enrollment customization
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Handle image upload if image source has changed
	if d.HasChange("enrollment_customization_image_source") {
		imagePath, err := constructImageUpload(d)
		if err == nil && imagePath != "" {
			err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				uploadResponse, apiErr := client.UploadEnrollmentCustomizationsImage(imagePath)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				// Store the URL in the schema for the main resource construction
				brandingSettings := d.Get("branding_settings").([]interface{})
				if len(brandingSettings) > 0 {
					settings := brandingSettings[0].(map[string]interface{})
					settings["icon_url"] = uploadResponse.Url
					brandingSettingsList := []interface{}{settings}
					if err := d.Set("branding_settings", brandingSettingsList); err != nil {
						return retry.NonRetryableError(fmt.Errorf("failed to set icon_url in schema: %v", err))
					}
				}
				return nil
			})
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to upload enrollment customization image: %v", err))
			}
		}
	}

	resource, err := constructBaseResource(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct enrollment customization: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateEnrollmentCustomizationByID(d.Id(), *resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update enrollment customization: %v", err))
	}

	// Handle changes to panes
	if d.HasChange("text_pane") || d.HasChange("ldap_pane") || d.HasChange("sso_pane") {
		// Validate pane combinations
		if err := validatePaneCombinations(d); err != nil {
			return diag.FromErr(err)
		}

		// Get current panes to determine what needs to be added, updated, or removed
		var panesList *jamfpro.ResponsePrestagePanesList
		err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			var apiErr error
			panesList, apiErr = client.GetPrestagePanes(d.Id())
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to get prestage panes for enrollment customization: %v", err))
		}

		// Create maps of existing panes by type
		existingPanes := mapExistingPanesByType(panesList)

		// Handle each pane type
		if d.HasChange("text_pane") {
			if err := handleTextPaneChanges(ctx, d, client, existingPanes["text"]); err != nil {
				return diag.FromErr(err)
			}
		}

		if d.HasChange("ldap_pane") {
			if err := handleLDAPPaneChanges(ctx, d, client, existingPanes["ldap"]); err != nil {
				return diag.FromErr(err)
			}
		}

		if d.HasChange("sso_pane") {
			if err := handleSSOPaneChanges(ctx, d, client, existingPanes["sso"]); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// mapExistingPanesByType organizes existing panes by their type
func mapExistingPanesByType(panesList *jamfpro.ResponsePrestagePanesList) map[string]map[int]jamfpro.PrestagePaneSummary {
	result := map[string]map[int]jamfpro.PrestagePaneSummary{
		"text": make(map[int]jamfpro.PrestagePaneSummary),
		"ldap": make(map[int]jamfpro.PrestagePaneSummary),
		"sso":  make(map[int]jamfpro.PrestagePaneSummary),
	}

	for _, panel := range panesList.Panels {
		if panel.Type == "text" || panel.Type == "ldap" || panel.Type == "sso" {
			result[panel.Type][panel.ID] = panel
		}
	}

	return result
}

// handleTextPaneChanges processes changes to text panes
func handleTextPaneChanges(ctx context.Context, d *schema.ResourceData, client *jamfpro.Client, existingPanes map[int]jamfpro.PrestagePaneSummary) error {
	old, new := d.GetChange("text_pane")
	oldPanes := old.([]interface{})
	newPanes := new.([]interface{})

	// Process deletions - panes that exist in old but not in new
	for _, oldPane := range oldPanes {
		oldPaneMap := oldPane.(map[string]interface{})
		if id, ok := oldPaneMap["id"].(int); ok && id > 0 {
			// Check if this pane exists in the current API state
			if _, exists := existingPanes[id]; !exists {
				continue // Skip if not found in API
			}

			// Check if this pane still exists in new config
			found := false
			for _, newPane := range newPanes {
				newPaneMap := newPane.(map[string]interface{})
				if newID, ok := newPaneMap["id"].(int); ok && newID == id {
					found = true
					break
				}
			}

			// If not found in new config, delete it
			if !found {
				paneIDStr := fmt.Sprintf("%d", id)
				err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					apiErr := client.DeletePrestagePane(d.Id(), paneIDStr)
					if apiErr != nil {
						return retry.RetryableError(apiErr)
					}
					return nil
				})

				if err != nil {
					return fmt.Errorf("failed to delete text pane %d: %v", id, err)
				}
			}
		}
	}

	// Process additions and updates
	for _, newPane := range newPanes {
		newPaneMap := newPane.(map[string]interface{})

		textPane, err := constructTextPane(newPaneMap)
		if err != nil {
			return fmt.Errorf("failed to construct text pane: %v", err)
		}

		isUpdate := false
		var paneID int

		if id, ok := newPaneMap["id"].(int); ok && id > 0 {
			paneID = id
			if _, exists := existingPanes[id]; exists {
				isUpdate = true
			}
		}

		if isUpdate {
			paneIDStr := fmt.Sprintf("%d", paneID)
			err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				_, apiErr := client.UpdateTextPrestagePaneByID(d.Id(), paneIDStr, *textPane)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("failed to update text pane %d: %v", paneID, err)
			}
		} else {
			err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				_, apiErr := client.CreateTextPrestagePane(d.Id(), *textPane)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("failed to create text pane: %v", err)
			}
		}
	}

	return nil
}

// handleLDAPPaneChanges processes changes to LDAP panes
func handleLDAPPaneChanges(ctx context.Context, d *schema.ResourceData, client *jamfpro.Client, existingPanes map[int]jamfpro.PrestagePaneSummary) error {
	old, new := d.GetChange("ldap_pane")
	oldPanes := old.([]interface{})
	newPanes := new.([]interface{})

	// Process deletions - panes that exist in old but not in new
	for _, oldPane := range oldPanes {
		oldPaneMap := oldPane.(map[string]interface{})
		if id, ok := oldPaneMap["id"].(int); ok && id > 0 {
			// Check if this pane exists in the current API state
			if _, exists := existingPanes[id]; !exists {
				continue // Skip if not found in API
			}

			// Check if this pane still exists in new config
			found := false
			for _, newPane := range newPanes {
				newPaneMap := newPane.(map[string]interface{})
				if newID, ok := newPaneMap["id"].(int); ok && newID == id {
					found = true
					break
				}
			}

			// If not found in new config, delete it
			if !found {
				paneIDStr := fmt.Sprintf("%d", id)
				err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					apiErr := client.DeletePrestagePane(d.Id(), paneIDStr)
					if apiErr != nil {
						return retry.RetryableError(apiErr)
					}
					return nil
				})

				if err != nil {
					return fmt.Errorf("failed to delete LDAP pane %d: %v", id, err)
				}
			}
		}
	}

	// Process additions and updates
	for _, newPane := range newPanes {
		newPaneMap := newPane.(map[string]interface{})

		ldapPane, err := constructLDAPPane(newPaneMap)
		if err != nil {
			return fmt.Errorf("failed to construct LDAP pane: %v", err)
		}

		isUpdate := false
		var paneID int

		if id, ok := newPaneMap["id"].(int); ok && id > 0 {
			paneID = id
			if _, exists := existingPanes[id]; exists {
				isUpdate = true
			}
		}

		if isUpdate {
			paneIDStr := fmt.Sprintf("%d", paneID)
			err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				_, apiErr := client.UpdateLDAPPrestagePaneByID(d.Id(), paneIDStr, *ldapPane)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("failed to update LDAP pane %d: %v", paneID, err)
			}
		} else {
			err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				_, apiErr := client.CreateLDAPPrestagePane(d.Id(), *ldapPane)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("failed to create LDAP pane: %v", err)
			}
		}
	}

	return nil
}

// handleSSOPaneChanges processes changes to SSO panes
func handleSSOPaneChanges(ctx context.Context, d *schema.ResourceData, client *jamfpro.Client, existingPanes map[int]jamfpro.PrestagePaneSummary) error {
	old, new := d.GetChange("sso_pane")
	oldPanes := old.([]interface{})
	newPanes := new.([]interface{})

	// Process deletions - panes that exist in old but not in new
	for _, oldPane := range oldPanes {
		oldPaneMap := oldPane.(map[string]interface{})
		if id, ok := oldPaneMap["id"].(int); ok && id > 0 {
			// Check if this pane exists in the current API state
			if _, exists := existingPanes[id]; !exists {
				continue // Skip if not found in API
			}

			// Check if this pane still exists in new config
			found := false
			for _, newPane := range newPanes {
				newPaneMap := newPane.(map[string]interface{})
				if newID, ok := newPaneMap["id"].(int); ok && newID == id {
					found = true
					break
				}
			}

			// If not found in new config, delete it
			if !found {
				paneIDStr := fmt.Sprintf("%d", id)
				err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					apiErr := client.DeleteSSOPrestagePane(d.Id(), paneIDStr)
					if apiErr != nil {
						return retry.RetryableError(apiErr)
					}
					return nil
				})

				if err != nil {
					return fmt.Errorf("failed to delete SSO pane %d: %v", id, err)
				}
			}
		}
	}

	// Process additions and updates
	for _, newPane := range newPanes {
		newPaneMap := newPane.(map[string]interface{})

		// Build the pane object
		ssoPane, err := constructSSOPane(newPaneMap)
		if err != nil {
			return fmt.Errorf("failed to construct SSO pane: %v", err)
		}

		// Check if this is an update or create
		isUpdate := false
		var paneID int

		if id, ok := newPaneMap["id"].(int); ok && id > 0 {
			paneID = id
			if _, exists := existingPanes[id]; exists {
				isUpdate = true
			}
		}

		if isUpdate {
			// Update existing pane
			paneIDStr := fmt.Sprintf("%d", paneID)
			err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				_, apiErr := client.UpdateSSOPrestagePaneByID(d.Id(), paneIDStr, *ssoPane)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("failed to update SSO pane %d: %v", paneID, err)
			}
		} else {
			// Create new pane
			err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				_, apiErr := client.CreateSSOPrestagePane(d.Id(), *ssoPane)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return fmt.Errorf("failed to create SSO pane: %v", err)
			}
		}
	}

	return nil
}

// delete is responsible for removing the enrollment customization
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteEnrollmentCustomizationByID(d.Id())
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete enrollment customization: %v", err))
	}

	d.SetId("")
	return diags
}
