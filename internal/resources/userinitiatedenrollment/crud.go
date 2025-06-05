package userinitiatedenrollment

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ResourceIDSingleton = "jamfpro_user_initiated_enrollment_settings_singleton"
)

// create is responsible for creating jamf pro User-initiated enrollment base settings, enrollment languages
// and ldap groups. It performs multiple api calls and therefore doesn't follow the function pattern
// of simpler resource types.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Step 1: Update main enrollment settings (API doesn't have a true "create" operation)
	enrollmentSettings, err := constructEnrollmentSettings(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct enrollment settings: %v", err))
	}

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

	// Step 2: Create language messaging configurations
	messagingList, err := constructEnrollmentMessaging(d, client)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct enrollment messaging: %v", err))
	}

	for i := range messagingList {
		message := &messagingList[i]
		languageCode := message.LanguageCode

		if languageCode == "" {
			return diag.FromErr(fmt.Errorf("cannot update language message for '%s': empty language code", message.Name))
		}

		err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			_, apiErr := client.UpdateEnrollmentMessageByLanguageID(languageCode, message)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update enrollment message for language '%s' (code: %s): %v",
				message.Name, languageCode, err))
		}
	}

	// Step 3: Create directory service group enrollment settings
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

	d.SetId(ResourceIDSingleton)

	readDiags := readNoCleanup(ctx, d, meta)
	diags = append(diags, readDiags...)

	return diags
}

// read handles reading the current state of the user initiated enrollment settings from Jamf Pro
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// --- Step 1: Read base enrollment configuration ---
	var enrollment *jamfpro.ResourceEnrollment
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		enrollment, apiErr = client.GetEnrollment()
		if apiErr != nil {
			log.Printf("[ERROR] Failed attempt to get jamf pro User-initiated enrollment config: %v", apiErr)
			return retry.RetryableError(apiErr)
		}
		return nil
	})
	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	// --- Step 2: Get all configured enrollment messages from API ---
	log.Print("[DEBUG] Fetching all configured enrollment language messages from Jamf Pro.")

	var enrollmentMessages []jamfpro.ResourceEnrollmentLanguage
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		enrollmentMessages, apiErr = client.GetEnrollmentMessages()
		if apiErr != nil {
			log.Printf("[ERROR] Failed to fetch enrollment messages: %v. Retrying...", apiErr)
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to retrieve enrollment language messages from API",
			Detail:   fmt.Sprintf("Could not fetch language messages. Error: %v", err),
		})
		enrollmentMessages = []jamfpro.ResourceEnrollmentLanguage{}
	} else {
		log.Printf("[DEBUG] Successfully fetched %d configured enrollment language messages.", len(enrollmentMessages))
	}

	// --- Step 3: Get all directory service group enrollment settings ---
	var accessGroupsList []jamfpro.ResourceAccountDrivenUserEnrollmentAccessGroup
	var currentGroups *jamfpro.ResponseAccountDrivenUserEnrollmentAccessGroupsList
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		currentGroups, apiErr = client.GetAccountDrivenUserEnrollmentAccessGroups(nil)
		if apiErr != nil {
			log.Printf("[ERROR] Failed to fetch directory service groups: %v. Retrying...", apiErr)
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to fetch directory service groups",
			Detail:   fmt.Sprintf("Error fetching all groups: %v", err),
		})
	} else if currentGroups != nil && len(currentGroups.Results) > 0 {
		// When we have access groups from the API, filter out any built-in groups (ID "1")
		for _, group := range currentGroups.Results {
			if group.ID == "1" {
				continue
			}
			accessGroupsList = append(accessGroupsList, group)
		}
	}

	// --- Step 4: Update state using the centralized function ---
	log.Print("[DEBUG] Updating Terraform state based on fetched API data.")
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

// update is responsible for updating an existing Jamf Pro UIE settings on the remote system.
// first it updates base config items.
// second it gets all existing enrollment messages, skips the built in english option and removes
// all other language settings. It then reapplies as needed. The http method is PUT.
// It then follows the same flow for LDAP Directory Service Groups.
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Step 1: Update main enrollment settings
	enrollmentSettings, err := constructEnrollmentSettings(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct updated enrollment settings: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateEnrollment(enrollmentSettings)
		if apiErr != nil {
			log.Printf("[ERROR] Failed to update base enrollment settings: %v", apiErr)
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update base enrollment settings: %v", err))
	}

	// Step 2: Process language messaging settings
	if d.HasChange("messaging") {
		log.Print("[DEBUG] Detected change in 'messaging'. Fetching current messages for cleanup.")

		// Step 2.1: Fetch currently configured enrollment messages
		var configuredMessages []jamfpro.ResourceEnrollmentLanguage
		err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
			var apiErr error
			configuredMessages, apiErr = client.GetEnrollmentMessages()
			if apiErr != nil {
				log.Printf("[ERROR] Failed to fetch existing language messages: %v. Retrying...", apiErr)
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch configured enrollment messages for cleanup: %v", err))
		}

		// Step 2.2: Identify and delete all non-English messages
		for _, msg := range configuredMessages {
			if msg.LanguageCode == "en" {
				continue // Skip built-in English
			}

			languageCode := msg.LanguageCode
			log.Printf("[DEBUG] Deleting enrollment language message for code: %s", languageCode)

			deleteErr := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				apiErr := client.DeleteEnrollmentMessageByLanguageID(languageCode)
				if apiErr != nil {
					if strings.Contains(apiErr.Error(), "404") {
						log.Printf("[WARN] Language '%s' already deleted or not found (404). Skipping.", languageCode)
						return nil
					}
					log.Printf("[ERROR] Failed to delete language '%s': %v. Retrying...", languageCode, apiErr)
					return retry.RetryableError(apiErr)
				}
				log.Printf("[DEBUG] Successfully deleted language '%s'", languageCode)
				return nil
			})

			if deleteErr != nil {
				return diag.FromErr(fmt.Errorf("failed to delete enrollment message for language '%s': %v", languageCode, deleteErr))
			}
		}

		// Step 2.3: Reconstruct and reapply updated messaging config
		messagingList, err := constructEnrollmentMessaging(d, client)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to construct enrollment messaging: %v", err))
		}

		for i := range messagingList {
			message := &messagingList[i]
			if message.LanguageCode == "" {
				return diag.FromErr(fmt.Errorf("cannot update enrollment message for language '%s': empty language code", message.Name))
			}

			log.Printf("[DEBUG] Creating/updating language configuration for '%s' (%s)", message.Name, message.LanguageCode)
			err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				_, apiErr := client.UpdateEnrollmentMessageByLanguageID(message.LanguageCode, message)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				return nil
			})

			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to update enrollment message for language '%s' (code: %s): %v", message.Name, message.LanguageCode, err))
			}
		}
	}

	// Step 3: Process directory service group enrollment settings (Delete All, Create New)
	if d.HasChange("directory_service_group_enrollment_settings") {
		log.Print("[DEBUG] Detected change in 'directory_service_group_enrollment_settings'. Applying delete-all, create-new strategy.")

		// 3.1: Delete all existing directory service group settings
		log.Print("[DEBUG] Fetching current directory service groups to delete them.")
		var currentGroups *jamfpro.ResponseAccountDrivenUserEnrollmentAccessGroupsList
		err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
			var apiErr error
			currentGroups, apiErr = client.GetAccountDrivenUserEnrollmentAccessGroups(nil)
			if apiErr != nil {
				log.Printf("[WARN] Error fetching current directory service groups for deletion: %v. Retrying...", apiErr)
				return retry.RetryableError(apiErr)
			}
			return nil
		})

		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to fetch current directory service groups before deletion",
				Detail:   fmt.Sprintf("Could not confirm existing groups to delete. Proceeding to create desired state, but old groups might remain if fetching failed. Error: %v", err),
			})
			// Continue anyway
		}

		if currentGroups != nil && len(currentGroups.Results) > 0 {
			log.Printf("[DEBUG] Found %d existing directory service groups to delete.", len(currentGroups.Results))
			for i := range currentGroups.Results {
				group := &currentGroups.Results[i]
				groupID := group.ID

				if groupID == "1" {
					log.Printf("[INFO] Skipping built-in directory service group ID '1' (Name: %s)", group.Name)
					continue
				}

				if groupID == "" {
					log.Printf("[WARN] Skipping group at index %d with empty ID.", i)
					continue
				}

				err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					apiErr := client.DeleteAccountDrivenUserEnrollmentAccessGroupByID(groupID)
					if apiErr != nil {
						if strings.Contains(apiErr.Error(), "404") {
							log.Printf("[WARN] Directory service group with ID '%s' returned 404 Not Found. Assuming already deleted.", groupID)
							return nil
						}
						log.Printf("[ERROR] Failed attempt to delete directory service group with ID '%s': %v. Retrying...", groupID, apiErr)
						return retry.RetryableError(apiErr)
					}
					return nil
				})

				if err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Failed to delete existing directory service group after retries",
						Detail:   fmt.Sprintf("Failed to delete group with API ID '%s' (Name: %s): %v", groupID, group.Name, err),
					})
					return diags
				}
			}
		} else {
			log.Print("[DEBUG] No existing directory service groups found to delete.")
		}

		// 3.2: Create new directory service group settings from Terraform config
		groupsToCreate, err := constructDirectoryServiceGroupSettings(d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to construct directory service group settings for creation: %v", err))
		}

		log.Printf("[DEBUG] Creating %d directory service groups from Terraform configuration.", len(groupsToCreate))
		for i := range groupsToCreate {
			group := groupsToCreate[i]
			log.Printf("[DEBUG] Creating directory service group: Name=%s, LDAP ID=%s, Group ID=%s", group.Name, group.LdapServerID, group.GroupID)

			err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				_, apiErr := client.CreateAccountDrivenUserEnrollmentAccessGroup(group)
				if apiErr != nil {
					log.Printf("[ERROR] Failed attempt to create directory service group for group '%s': %v", group.Name, apiErr)
					return retry.RetryableError(apiErr)
				}
				log.Printf("[DEBUG] Successfully created directory service group: Name=%s", group.Name)
				return nil
			})

			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to create directory service group",
					Detail:   fmt.Sprintf("Failed to create directory service group for '%s' (LDAP ID: %s, Group ID: %s): %v", group.Name, group.LdapServerID, group.GroupID, err),
				})
				return diags
			}
		}
	}

	log.Print("[DEBUG] Update process complete, running readNoCleanup to refresh state.")
	readDiags := readNoCleanup(ctx, d, meta)
	diags = append(diags, readDiags...)

	return diags
}

// delete handles cleanup of specific sub-configurations in Jamf Pro
// before removing the resource from Terraform state.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Starting deletion cleanup for resource %s.", ResourceIDSingleton)

	// --- Step 1: Delete language messaging (except for English) ---
	log.Print("[DEBUG] Fetching configured language messages to delete from Jamf Pro.")
	configuredMessages, err := client.GetEnrollmentMessages()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to fetch configured language messages",
			Detail:   fmt.Sprintf("Skipping language cleanup. Error: %v", err),
		})
	} else {
		for _, msg := range configuredMessages {
			if msg.LanguageCode == "en" {
				log.Print("[DEBUG] Skipping built-in English language message.")
				continue
			}

			languageCode := msg.LanguageCode
			log.Printf("[DEBUG] Deleting enrollment language message for code: %s", languageCode)

			deleteErr := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
				apiErr := client.DeleteEnrollmentMessageByLanguageID(languageCode)
				if apiErr != nil {
					if strings.Contains(apiErr.Error(), "404") {
						log.Printf("[WARN] Language code '%s' already deleted or not found (404). Skipping.", languageCode)
						return nil
					}
					log.Printf("[ERROR] Failed to delete language code '%s': %v. Retrying...", languageCode, apiErr)
					return retry.RetryableError(apiErr)
				}
				log.Printf("[DEBUG] Successfully deleted enrollment language for code: %s", languageCode)
				return nil
			})

			if deleteErr != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to delete enrollment language message",
					Detail:   fmt.Sprintf("Could not delete message for language '%s'. Error: %v", languageCode, deleteErr),
				})
			}
		}
	}

	// --- Step 2: Delete directory service group enrollment settings ---
	log.Print("[DEBUG] Attempting to fetch and delete all directory service group enrollment settings (skipping ID '1').") // Updated log message
	var accessGroups *jamfpro.ResponseAccountDrivenUserEnrollmentAccessGroupsList
	fetchErr := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		var apiErr error
		accessGroups, apiErr = client.GetAccountDrivenUserEnrollmentAccessGroups(nil)
		if apiErr != nil {
			if strings.Contains(apiErr.Error(), "404") {
				log.Printf("[WARN] Received 404 fetching directory service groups. Assuming none exist.")
				accessGroups = nil
				return nil
			}
			log.Printf("[ERROR] Failed fetch directory groups: %v. Retrying...", apiErr)
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if fetchErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to fetch directory service groups for deletion",
			Detail:   fmt.Sprintf("Could not confirm groups to delete. Error: %v", fetchErr),
		})
	} else if accessGroups != nil && len(accessGroups.Results) > 0 {
		log.Printf("[DEBUG] Found %d directory service groups to process for deletion.", len(accessGroups.Results))
		for i := range accessGroups.Results {
			group := &accessGroups.Results[i]
			groupID := group.ID

			if groupID == "1" {
				log.Printf("[INFO] Skipping deletion of directory service group with ID '1' (Name: %s) as it is considered built-in.", group.Name)
				continue
			}

			if groupID == "" {
				log.Printf("[WARN] Skipping group at index %d with empty ID.", i)
				continue
			}

			log.Printf("[DEBUG] Deleting directory service group: API ID=%s, Name=%s", groupID, group.Name)
			deleteErr := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
				apiErr := client.DeleteAccountDrivenUserEnrollmentAccessGroupByID(groupID)
				if apiErr != nil {
					if strings.Contains(apiErr.Error(), "404") {
						log.Printf("[WARN] Group ID '%s' (404). Already gone.", groupID)
						return nil
					}
					log.Printf("[ERROR] Failed delete group ID '%s': %v. Retrying...", groupID, apiErr)
					return retry.RetryableError(apiErr)
				}
				log.Printf("[DEBUG] Successfully deleted group: API ID=%s", groupID)
				return nil
			})
			if deleteErr != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to delete a directory service group",
					Detail:   fmt.Sprintf("Group API ID '%s' (Name: '%s') might still exist. Error: %v", groupID, group.Name, deleteErr),
				})
			}
		}
	} else {
		log.Print("[DEBUG] No directory service groups found in Jamf Pro to delete.")
	}

	log.Printf("[DEBUG] Deletion cleanup finished for resource %s. Removing from Terraform state.", ResourceIDSingleton)
	d.SetId("")

	return diags
}
