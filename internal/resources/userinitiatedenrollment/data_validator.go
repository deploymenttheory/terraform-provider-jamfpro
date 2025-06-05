package userinitiatedenrollment

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Main CustomizeDiff function that orchestrates all diff customizations
func customizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	err := customizeDiffForEnglishLanguage(ctx, d, meta)
	if err != nil {
		return err
	}

	return nil
}

// customizeDiffForEnglishLanguage ensures that the English language exists
// in the enrollment languages block
func customizeDiffForEnglishLanguage(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if v, ok := d.GetOk("messaging"); ok {
		messagingSet := v.(*schema.Set).List()

		// Check if English is included in the configuration, it should exist.
		// it can never be removed from the gui.
		var hasEnglish bool
		for _, messaging := range messagingSet {
			msg := messaging.(map[string]interface{})
			langName := msg["language_name"].(string)
			normalizedLangName := strings.ToLower(strings.TrimSpace(langName))

			// Check for English - consider both "english" and checking language code if available
			if normalizedLangName == "english" {
				hasEnglish = true
				break
			}
		}

		if !hasEnglish {
			return fmt.Errorf("english language enrollment messaging is required, please include an English enrollment messaging with 'en' and 'english'")
		}
	} else if d.Id() != "" {
		old, _ := d.GetChange("messaging")
		oldSet := old.(*schema.Set)

		if oldSet != nil && oldSet.Len() > 0 {
			return fmt.Errorf("cannot remove all messaging configurations as English language configuration is required")
		}
	}

	return nil
}
