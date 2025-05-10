package cloudldap

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs a ResourceCloudLdap object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourceCloudLdap, error) {
	resource := &jamfpro.ResourceCloudLdap{
		CloudIdPCommon: constructCloudIdPCommon(d),
		Server:         constructCloudLdapServer(d),
	}

	if v, ok := d.GetOk("mappings"); ok {
		resource.Mappings = constructCloudLdapMappings(v.([]interface{}))
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Cloud LDAP '%s' to JSON: %v", resource.CloudIdPCommon.DisplayName, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Cloud LDAP JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}

// constructCloudIdPCommon constructs a CloudIdPCommon object from the provided schema data.
func constructCloudIdPCommon(d *schema.ResourceData) *jamfpro.CloudIdPCommon {
	commonData := d.Get("cloud_idp_common").([]interface{})[0].(map[string]interface{})
	return &jamfpro.CloudIdPCommon{
		ID:           d.Id(),
		ProviderName: commonData["provider_name"].(string),
		DisplayName:  commonData["display_name"].(string),
	}
}

// constructCloudLdapServer constructs a CloudLdapServer object from the provided schema data.
func constructCloudLdapServer(d *schema.ResourceData) *jamfpro.CloudLdapServer {
	serverData := d.Get("server").([]interface{})[0].(map[string]interface{})
	server := &jamfpro.CloudLdapServer{
		Enabled:                                  serverData["enabled"].(bool),
		UseWildcards:                             serverData["use_wildcards"].(bool),
		ConnectionType:                           serverData["connection_type"].(string),
		ServerUrl:                                serverData["server_url"].(string),
		DomainName:                               serverData["domain_name"].(string),
		Port:                                     serverData["port"].(int),
		ConnectionTimeout:                        serverData["connection_timeout"].(int),
		SearchTimeout:                            serverData["search_timeout"].(int),
		MembershipCalculationOptimizationEnabled: serverData["membership_calculation_optimization_enabled"].(bool),
	}

	if v, ok := serverData["keystore"]; ok && len(v.([]interface{})) > 0 {
		keystoreData := v.([]interface{})[0].(map[string]interface{})
		server.Keystore = &jamfpro.CloudLdapKeystore{
			Password:  keystoreData["password"].(string),
			FileBytes: keystoreData["file_bytes"].(string),
			FileName:  keystoreData["file_name"].(string),
		}
	}

	return server
}

// constructCloudLdapMappings constructs a CloudLdapMappings object from the provided schema data.
func constructCloudLdapMappings(mappingsList []interface{}) *jamfpro.CloudLdapMappings {
	if len(mappingsList) == 0 {
		return nil
	}

	mappingsData := mappingsList[0].(map[string]interface{})
	mappings := &jamfpro.CloudLdapMappings{}

	if um, ok := mappingsData["user_mappings"].([]interface{}); ok && len(um) > 0 {
		userMap := um[0].(map[string]interface{})
		mappings.UserMappings = jamfpro.CloudIdentityProviderDefaultMappingsSubsetUserMappings{
			ObjectClassLimitation: userMap["object_class_limitation"].(string),
			ObjectClasses:         userMap["object_classes"].(string),
			SearchBase:            userMap["search_base"].(string),
			SearchScope:           userMap["search_scope"].(string),
			AdditionalSearchBase:  userMap["additional_search_base"].(string),
			UserID:                userMap["user_id"].(string),
			Username:              userMap["username"].(string),
			RealName:              userMap["real_name"].(string),
			EmailAddress:          userMap["email_address"].(string),
			Department:            userMap["department"].(string),
			Building:              userMap["building"].(string),
			Room:                  userMap["room"].(string),
			Phone:                 userMap["phone"].(string),
			Position:              userMap["position"].(string),
			UserUuid:              userMap["user_uuid"].(string),
		}
	}

	if gm, ok := mappingsData["group_mappings"].([]interface{}); ok && len(gm) > 0 {
		groupMap := gm[0].(map[string]interface{})
		mappings.GroupMappings = jamfpro.CloudIdentityProviderDefaultMappingsSubsetGroupMappings{
			ObjectClassLimitation: groupMap["object_class_limitation"].(string),
			ObjectClasses:         groupMap["object_classes"].(string),
			SearchBase:            groupMap["search_base"].(string),
			SearchScope:           groupMap["search_scope"].(string),
			GroupID:               groupMap["group_id"].(string),
			GroupName:             groupMap["group_name"].(string),
			GroupUuid:             groupMap["group_uuid"].(string),
		}
	}

	if mm, ok := mappingsData["membership_mappings"].([]interface{}); ok && len(mm) > 0 {
		membershipMap := mm[0].(map[string]interface{})
		mappings.MembershipMappings = jamfpro.CloudIdentityProviderDefaultMappingsSubsetMembershipMappings{
			GroupMembershipMapping: membershipMap["group_membership_mapping"].(string),
		}
	}

	return mappings
}
