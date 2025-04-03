package userinitiatedenrollment

import (
	"context"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resource ID constants
const (
	ResourceIDSingleton = "jamfpro_user_initiated_enrollment_settings_singleton"
)

// create handles creation of the user initiated enrollment settings in Jamf Pro
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Step 1: Construct main enrollment settings
	enrollmentSettings, err := constructEnrollmentSettings(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct enrollment settings: %v", err))
	}

	// Step 2: Update main enrollment settings (API doesn't have a true "create" operation)
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		_, apiErr := client.UpdateEnrollment(enrollmentSettings)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update enrollment settings: %v", err))
	}

	// Step 3: Create language messaging configurations
	messagingList, err := constructEnrollmentMessaging(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct enrollment messaging: %v", err))
	}

	for i := range messagingList {
		// Get a pointer to the message in the slice
		message := &messagingList[i]
		languageCode := message.LanguageCode
		err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			_, apiErr := client.UpdateEnrollmentMessageByLanguageID(languageCode, message)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update enrollment message for language '%s': %v", message.Name, err))
		}
	}

	// Step 4: Create directory service group enrollment settings
	accessGroups, err := constructDirectoryServiceGroupSettings(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct directory service group settings: %v", err))
	}

	for i := range accessGroups {
		group := accessGroups[i]
		err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			_, apiErr := client.CreateAccountDrivenUserEnrollmentAccessGroup(group)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to create directory service group enrollment setting for group '%s': %v", group.Name, err))
		}
	}

	// Set ID to indicate resource was created successfully
	d.SetId(ResourceIDSingleton)

	// Read the resource to update the state with any computed values
	readDiags := readNoCleanup(ctx, d, meta)
	diags = append(diags, readDiags...)

	return diags
}

// read handles reading the current state of the user initiated enrollment settings from Jamf Pro
// read handles reading the current state of the user initiated enrollment settings from Jamf Pro
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Ensure consistent ID for singleton resource
	d.SetId(ResourceIDSingleton)

	// Step 1: Read main enrollment configuration
	var enrollment *jamfpro.ResourceEnrollment
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		enrollment, apiErr = client.GetEnrollment()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	// Step 2: Get all enrollment messages
	var enrollmentMessages []jamfpro.ResourceEnrollmentLanguage

	// Get current messaging configurations from state to know which language codes to fetch
	if v, ok := d.GetOk("messaging"); ok {
		messagingSet := v.(*schema.Set).List()

		for _, messaging := range messagingSet {
			if msg, ok := messaging.(map[string]interface{}); ok {
				languageCode := msg["language_code"].(string)

				var message *jamfpro.ResourceEnrollmentLanguage
				// For language messages
				err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
					var apiErr error
					message, apiErr = client.GetEnrollmentMessageByLanguageID(languageCode)
					if apiErr != nil {
						return retry.RetryableError(apiErr)
					}
					return nil
				})

				if err != nil {
					log.Printf("[WARN] Failed to get enrollment message for language code '%s': %v", languageCode, err)
					continue
				}

				if message != nil {
					enrollmentMessages = append(enrollmentMessages, *message)
				}
			}
		}
	}

	// Step 3: Get directory service group enrollment settings
	var accessGroupsList []jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		accessGroups, apiErr := client.GetAccountDrivenUserEnrollmentAccessGroups("")
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		if accessGroups != nil && len(accessGroups.Results) > 0 {
			accessGroupsList = accessGroups.Results
		}

		return nil
	})

	if err != nil {
		// Add a warning but continue with the state update
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to retrieve directory service group enrollment settings",
			Detail:   fmt.Sprintf("Error: %v", err),
		})
	}

	// Step 4: Update state using the centralized function
	stateDiags := updateState(d, enrollment, enrollmentMessages, accessGroupsList)
	diags = append(diags, stateDiags...)

	return diags
}

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update handles updating the user initiated enrollment settings in Jamf Pro
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Step 1: Update main enrollment configuration if needed
	if hasEnrollmentSettingsChange(d) {
		enrollmentSettings, err := constructEnrollmentSettings(d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to construct enrollment settings: %v", err))
		}

		err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			_, apiErr := client.UpdateEnrollment(enrollmentSettings)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update enrollment settings: %v", err))
		}
	}

	// Step 2: Process language messaging settings
	if d.HasChange("messaging") {
		old, new := d.GetChange("messaging")
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)

		// Find languages to delete (in old but not in new)
		languagesToDelete, err := findLanguagesToDelete(oldSet.List(), newSet.List())
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to determine languages to delete: %v", err))
		}

		// Delete languages no longer in the configuration
		if len(languagesToDelete) > 0 {
			err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				apiErr := client.DeleteMultipleEnrollmentMessagesByLanguageIDs(languagesToDelete)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to delete enrollment message languages: %v", err))
			}
		}

		// Update or create languages in the new configuration
		messagingList, err := constructEnrollmentMessagingFromSet(newSet.List())
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to construct enrollment messaging: %v", err))
		}

		for i := range messagingList {
			// Get a pointer to the message in the slice
			message := &messagingList[i]
			languageCode := message.LanguageCode
			err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				_, apiErr := client.UpdateEnrollmentMessageByLanguageID(languageCode, message)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to update enrollment message for language '%s': %v", message.Name, err))
			}
		}
	}

	// Step 3: Process directory service group enrollment settings
	if d.HasChange("directory_service_group_enrollment_settings") {
		old, new := d.GetChange("directory_service_group_enrollment_settings")

		// Get current groups from API to have their IDs
		currentGroups, err := client.GetAccountDrivenUserEnrollmentAccessGroups("")
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to get directory service group enrollment settings: %v", err))
		}

		// Process updates based on old and new sets
		toDelete, toUpdate, toCreate, err := processDirectoryServiceGroupChanges(
			old.(*schema.Set).List(),
			new.(*schema.Set).List(),
			currentGroups.Results,
		)

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to process directory service group changes: %v", err))
		}

		// Delete groups
		for i := range toDelete {
			group := &toDelete[i]
			err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				apiErr := client.DeleteAccountDrivenUserEnrollmentAccessGroupByID(group.ID)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to delete directory service group with ID '%s': %v", group.ID, err))
			}
		}

		// Update groups
		for i := range toUpdate {
			group := &toUpdate[i]
			err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				_, apiErr := client.UpdateAccountDrivenUserEnrollmentAccessGroupByID(group.ID, group)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to update directory service group with ID '%s': %v", group.ID, err))
			}
		}

		// Create new groups
		for i := range toCreate {
			group := &toCreate[i]
			err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				_, apiErr := client.CreateAccountDrivenUserEnrollmentAccessGroup(group)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to create directory service group for group '%s': %v", group.Name, err))
			}
		}
	}

	// Read the resource to update the state with any computed values
	readDiags := readNoCleanup(ctx, d, meta)
	diags = append(diags, readDiags...)

	return diags
}

// delete handles cleanup of the user initiated enrollment settings in Jamf Pro
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Step 1: Reset main enrollment configuration to defaults
	defaultSettings := constructDefaultEnrollmentSettings()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		_, apiErr := client.UpdateEnrollment(defaultSettings)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		// Add warning to diagnostics but continue with other cleanups
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to reset enrollment settings to defaults",
			Detail:   fmt.Sprintf("Error: %v", err),
		})
	}

	// Step 2: Delete language messaging (except for English which is required)
	if v, ok := d.GetOk("messaging"); ok {
		languageCodesToDelete := extractLanguageCodesToDelete(v.(*schema.Set).List())

		if len(languageCodesToDelete) > 0 {
			err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
				apiErr := client.DeleteMultipleEnrollmentMessagesByLanguageIDs(languageCodesToDelete)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				// Add warning to diagnostics but continue with other cleanups
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to delete language messaging configurations",
					Detail:   fmt.Sprintf("Error: %v", err),
				})
			}
		}
	}

	// Step 3: Delete directory service group enrollment settings
	accessGroups, err := client.GetAccountDrivenUserEnrollmentAccessGroups("")
	if err != nil {
		// Add warning to diagnostics but continue with cleanup
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to get directory service group enrollment settings",
			Detail:   fmt.Sprintf("Error: %v", err),
		})
	} else {
		for i := range accessGroups.Results {
			group := accessGroups.Results[i]
			groupID := group.ID
			err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
				apiErr := client.DeleteAccountDrivenUserEnrollmentAccessGroupByID(groupID)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				// Add warning to diagnostics but continue with other deletions
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to delete directory service group",
					Detail:   fmt.Sprintf("Group ID '%s': %v", groupID, err),
				})
			}
		}
	}

	// Remove from state
	d.SetId("")
	return diags
}
