// ldapserver_state.go
package ldapservers

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest LDAP Server information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceLDAPServers) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update Connection fields

	if err := d.Set("id", strconv.Itoa(resp.Connection.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resp.Connection.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("hostname", resp.Connection.Hostname); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("server_type", resp.Connection.ServerType); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("port", resp.Connection.Port); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("use_ssl", resp.Connection.UseSSL); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("authentication_type", resp.Connection.AuthenticationType); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("open_close_timeout", resp.Connection.OpenCloseTimeout); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("search_timeout", resp.Connection.SearchTimeout); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("referral_response", resp.Connection.ReferralResponse); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("use_wildcards", resp.Connection.UseWildcards); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Update Account information (excluding password)
	account := []interface{}{
		map[string]interface{}{
			"distinguished_username": resp.Connection.Account.DistinguishedUsername,
			// Password is not included as it's sensitive and not returned by the API
		},
	}
	if err := d.Set("account", account); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Update User Mappings
	userMappings := []interface{}{
		map[string]interface{}{
			"map_object_class_to_any_or_all": resp.MappingsForUsers.UserMappings.MapObjectClassToAnyOrAll,
			"object_classes":                 resp.MappingsForUsers.UserMappings.ObjectClasses,
			"search_base":                    resp.MappingsForUsers.UserMappings.SearchBase,
			"search_scope":                   resp.MappingsForUsers.UserMappings.SearchScope,
			"map_user_id":                    resp.MappingsForUsers.UserMappings.MapUserID,
			"map_username":                   resp.MappingsForUsers.UserMappings.MapUsername,
			"map_realname":                   resp.MappingsForUsers.UserMappings.MapRealName,
			"map_email_address":              resp.MappingsForUsers.UserMappings.MapEmailAddress,
			"append_to_email_results":        resp.MappingsForUsers.UserMappings.AppendToEmailResults,
			"map_department":                 resp.MappingsForUsers.UserMappings.MapDepartment,
			"map_building":                   resp.MappingsForUsers.UserMappings.MapBuilding,
			"map_room":                       resp.MappingsForUsers.UserMappings.MapRoom,
			"map_phone":                      resp.MappingsForUsers.UserMappings.MapPhone,
			"map_position":                   resp.MappingsForUsers.UserMappings.MapPosition,
			"map_user_uuid":                  resp.MappingsForUsers.UserMappings.MapUserUUID,
		},
	}
	if err := d.Set("user_mappings", userMappings); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Update User Group Mappings
	userGroupMappings := []interface{}{
		map[string]interface{}{
			"map_object_class_to_any_or_all": resp.MappingsForUsers.UserGroupMappings.MapObjectClassToAnyOrAll,
			"object_classes":                 resp.MappingsForUsers.UserGroupMappings.ObjectClasses,
			"search_base":                    resp.MappingsForUsers.UserGroupMappings.SearchBase,
			"search_scope":                   resp.MappingsForUsers.UserGroupMappings.SearchScope,
			"map_group_id":                   resp.MappingsForUsers.UserGroupMappings.MapGroupID,
			"map_group_name":                 resp.MappingsForUsers.UserGroupMappings.MapGroupName,
			"map_group_uuid":                 resp.MappingsForUsers.UserGroupMappings.MapGroupUUID,
		},
	}
	if err := d.Set("user_group_mappings", userGroupMappings); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Update User Group Membership Mappings
	membershipMappings := []interface{}{
		map[string]interface{}{
			"user_group_membership_stored_in":    resp.MappingsForUsers.UserGroupMembershipMappings.UserGroupMembershipStoredIn,
			"map_group_membership_to_user_field": resp.MappingsForUsers.UserGroupMembershipMappings.MapGroupMembershipToUserField,
			"append_to_username":                 resp.MappingsForUsers.UserGroupMembershipMappings.AppendToUsername,
			"use_dn":                             resp.MappingsForUsers.UserGroupMembershipMappings.UseDN,
			"group_membership_enabled_when_user_membership_selected": resp.MappingsForUsers.UserGroupMembershipMappings.GroupMembershipEnabledWhenUserMembershipSelected,
			"recursive_lookups":                      resp.MappingsForUsers.UserGroupMembershipMappings.RecursiveLookups,
			"map_user_membership_to_group_field":     resp.MappingsForUsers.UserGroupMembershipMappings.MapUserMembershipToGroupField,
			"map_user_membership_use_dn":             resp.MappingsForUsers.UserGroupMembershipMappings.MapUserMembershipUseDN,
			"map_object_class_to_any_or_all":         resp.MappingsForUsers.UserGroupMembershipMappings.MapObjectClassToAnyOrAll,
			"object_classes":                         resp.MappingsForUsers.UserGroupMembershipMappings.ObjectClasses,
			"search_base":                            resp.MappingsForUsers.UserGroupMembershipMappings.SearchBase,
			"search_scope":                           resp.MappingsForUsers.UserGroupMembershipMappings.SearchScope,
			"username":                               resp.MappingsForUsers.UserGroupMembershipMappings.Username,
			"group_id":                               resp.MappingsForUsers.UserGroupMembershipMappings.GroupID,
			"user_group_membership_use_ldap_compare": resp.MappingsForUsers.UserGroupMembershipMappings.UserGroupMembershipUseLDAPCompare,
			"membership_scoping_optimization":        resp.MappingsForUsers.UserGroupMembershipMappings.MembershipScopingOptimization,
		},
	}
	if err := d.Set("user_group_membership_mappings", membershipMappings); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
