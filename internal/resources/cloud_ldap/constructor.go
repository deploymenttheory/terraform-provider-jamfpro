package cloud_ldap

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructCloudLdap constructs a ResourceCloudLdap object from the provided schema data
func construct(d *schema.ResourceData) (*jamfpro.ResourceCloudLdap, error) {
	cloudIdpCommon := jamfpro.CloudIdPCommon{
		ID:           d.Id(),
		ProviderName: d.Get("provider_name").(string),
		DisplayName:  d.Get("display_name").(string),
	}

	keystore := jamfpro.CloudLdapKeystore{
		Password:  d.Get("keystore_password").(string),
		FileBytes: d.Get("keystore_file_bytes").(string),
		FileName:  d.Get("keystore_file_name").(string),
	}

	server := jamfpro.CloudLdapServer{
		Enabled:                                  d.Get("server_enabled").(bool),
		Keystore:                                 &keystore,
		UseWildcards:                             d.Get("use_wildcards").(bool),
		ConnectionType:                           d.Get("connection_type").(string),
		ServerUrl:                                d.Get("server_url").(string),
		DomainName:                               d.Get("domain_name").(string),
		Port:                                     d.Get("port").(int),
		ConnectionTimeout:                        d.Get("connection_timeout").(int),
		SearchTimeout:                            d.Get("search_timeout").(int),
		MembershipCalculationOptimizationEnabled: d.Get("membership_calculation_optimization_enabled").(bool),
	}

	userMappings := jamfpro.CloudIdentityProviderDefaultMappingsSubsetUserMappings{
		ObjectClassLimitation: d.Get("user_mappings_object_class_limitation").(string),
		ObjectClasses:         d.Get("user_mappings_object_classes").(string),
		SearchBase:            d.Get("user_mappings_search_base").(string),
		SearchScope:           d.Get("user_mappings_search_scope").(string),
		AdditionalSearchBase:  d.Get("user_mappings_additional_search_base").(string),
		UserID:                d.Get("user_mappings_id").(string),
		Username:              d.Get("user_mappings_username").(string),
		RealName:              d.Get("user_mappings_real_name").(string),
		EmailAddress:          d.Get("user_mappings_email_address").(string),
		Department:            d.Get("user_mappings_department").(string),
		Building:              d.Get("user_mappings_building").(string),
		Room:                  d.Get("user_mappings_room").(string),
		Phone:                 d.Get("user_mappings_phone").(string),
		Position:              d.Get("user_mappings_position").(string),
		UserUuid:              d.Get("user_mappings_uuid").(string),
	}

	groupMappings := jamfpro.CloudIdentityProviderDefaultMappingsSubsetGroupMappings{
		ObjectClassLimitation: d.Get("group_mappings_object_class_limitation").(string),
		ObjectClasses:         d.Get("group_mappings_object_classes").(string),
		SearchBase:            d.Get("group_mappings_search_base").(string),
		SearchScope:           d.Get("group_mappings_search_scope").(string),
		GroupID:               d.Get("group_mappings_id").(string),
		GroupName:             d.Get("group_mappings_name").(string),
		GroupUuid:             d.Get("group_mappings_uuid").(string),
	}

	membershipMappings := jamfpro.CloudIdentityProviderDefaultMappingsSubsetMembershipMappings{
		GroupMembershipMapping: d.Get("group_membership_mapping").(string),
	}

	mappings := jamfpro.CloudLdapMappings{
		UserMappings:       userMappings,
		GroupMappings:      groupMappings,
		MembershipMappings: membershipMappings,
	}

	resource := &jamfpro.ResourceCloudLdap{
		CloudIdPCommon: &cloudIdpCommon,
		Server:         &server,
		Mappings:       &mappings,
	}

	resourceJSON, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Cloud LDAP '%s' to JSON: %v", resource.CloudIdPCommon.DisplayName, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Cloud LDAP JSON:\n%s\n", string(resourceJSON))

	return resource, nil
}
