package user_initiated_enrollment_settings

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Main CustomizeDiff function that orchestrates all diff customizations
func customizeDiff(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	err := customizeDiffForEnglishLanguage(ctx, d, meta)
	if err != nil {
		return err
	}

	if err := customizeDiffForAccessGroupOne(ctx, d, meta); err != nil {
		return err
	}

	return nil
}

// customizeDiffForEnglishLanguage ensures that the English language exists
// in the enrollment languages block
func customizeDiffForEnglishLanguage(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	if v, ok := d.GetOk("messaging"); ok {
		messagingSet := v.(*schema.Set).List()

		// Check if English is included in the configuration, it should exist.
		// it can never be removed from the gui.
		var hasEnglish bool
		for _, messaging := range messagingSet {
			msg := messaging.(map[string]any)
			langName := msg["language_name"].(string)
			normalizedLangName := strings.ToLower(strings.TrimSpace(langName))

			// Check for English - consider both "english" and checking language code if available
			if normalizedLangName == "english" {
				hasEnglish = true
				break
			}
		}

		if !hasEnglish {
			return fmt.Errorf("english language enrollment messaging is required, please include an English enrollment messaging with 'en' and 'english'") //nolint:err113
		}
	} else if d.Id() != "" {
		old, _ := d.GetChange("messaging")
		oldSet := old.(*schema.Set)

		if oldSet != nil && oldSet.Len() > 0 {
			return fmt.Errorf("cannot remove all messaging configurations as English language configuration is required") //nolint:err113
		}
	}

	return nil
}

// customizeDiffForAccessGroupOne ensures that Access Group ID 1 is handled correctly
func customizeDiffForAccessGroupOne(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	if v, ok := d.GetOk("directory_service_group_enrollment_settings"); ok {
		groupSet := v.(*schema.Set).List()

		for _, group := range groupSet {
			groupMap := group.(map[string]any)

			if groupID, ok := groupMap["directory_service_group_id"].(string); ok && groupID == "1" {
				if d.Id() == "" {
					return fmt.Errorf("access Group ID 1 is built-in and cannot be created") //nolint:err113
				}

				if d.HasChange("directory_service_group_enrollment_settings") {
					oldGroups, _ := d.GetChange("directory_service_group_enrollment_settings")
					oldGroupSet := oldGroups.(*schema.Set).List()

					var oldGroup1 map[string]any
					for _, og := range oldGroupSet {
						ogMap := og.(map[string]any)
						if ogID, ok := ogMap["directory_service_group_id"].(string); ok && ogID == "1" {
							oldGroup1 = ogMap
							break
						}
					}

					if oldGroup1 != nil {
						immutableFields := map[string]string{
							"directory_service_group_id":   "Directory Service Group ID",
							"ldap_server_id":               "LDAP Server ID",
							"directory_service_group_name": "Directory Service Group Name",
						}

						for field, displayName := range immutableFields {
							if oldGroup1[field] != groupMap[field] {
								return fmt.Errorf("%s cannot be modified for Access Group ID 1", displayName) //nolint:err113
							}
						}
					}
				}
			}
		}
	}

	return nil
}
