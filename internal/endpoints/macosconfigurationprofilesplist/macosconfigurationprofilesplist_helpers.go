// macosconfigurationprofilesplist_helpers.go
package macosconfigurationprofilesplist

import (
	"fmt"
	"log"
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

// FixDuplicateNotificationKey handles the double key issue in the notification field of the self_service block.
/*
<self_service>
        <self_service_display_name>WiFi Test</self_service_display_name>
        <install_button_text>Install</install_button_text>
        <self_service_description>null</self_service_description>
        <force_users_to_view_description>false</force_users_to_view_description>
        <security>
            <removal_disallowed>Never</removal_disallowed>
        </security>
        <self_service_icon/>
        <feature_on_main_page>false</feature_on_main_page>
        <self_service_categories/>
        <notification>false</notification>				<-- This is the issue
        <notification>Self Service</notification>  <-- This is the issue
        <notification_subject/>
        <notification_message/>
    </self_service>
*/
func FixDuplicateNotificationKey(resp *jamfpro.ResourceMacOSConfigurationProfile) (bool, error) {
	for _, k := range resp.SelfService.Notification {
		strValue := fmt.Sprintf("%v", k)
		if strValue == "true" || strValue == "false" {
			correctNotifValue, err := strconv.ParseBool(strValue)
			if err != nil {
				return false, err
			}
			return correctNotifValue, nil
		} else {
			log.Printf("Ignoring non-boolean notification value: %s", strValue)
		}
	}
	// Return default value if no valid boolean value is found
	return false, nil
}
