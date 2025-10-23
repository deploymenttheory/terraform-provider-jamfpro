package cloud_ldap

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the latest Cloud LDAP settings information from the Jamf Pro API.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceCloudLdap) diag.Diagnostics {
	var diags diag.Diagnostics

	settings := map[string]any{
		"provider_name":      resp.CloudIdPCommon.ProviderName,
		"display_name":       resp.CloudIdPCommon.DisplayName,
		"server_enabled":     resp.Server.Enabled,
		"use_wildcards":      resp.Server.UseWildcards,
		"connection_type":    resp.Server.ConnectionType,
		"server_url":         resp.Server.ServerUrl,
		"domain_name":        resp.Server.DomainName,
		"port":               resp.Server.Port,
		"connection_timeout": resp.Server.ConnectionTimeout,
		"search_timeout":     resp.Server.SearchTimeout,
		"membership_calculation_optimization_enabled": resp.Server.MembershipCalculationOptimizationEnabled,
		"user_mappings_object_class_limitation":       resp.Mappings.UserMappings.ObjectClassLimitation,
		"user_mappings_object_classes":                resp.Mappings.UserMappings.ObjectClasses,
		"user_mappings_search_base":                   resp.Mappings.UserMappings.SearchBase,
		"user_mappings_search_scope":                  resp.Mappings.UserMappings.SearchScope,
		"user_mappings_additional_search_base":        resp.Mappings.UserMappings.AdditionalSearchBase,
		"user_mappings_id":                            resp.Mappings.UserMappings.UserID,
		"user_mappings_username":                      resp.Mappings.UserMappings.Username,
		"user_mappings_real_name":                     resp.Mappings.UserMappings.RealName,
		"user_mappings_email_address":                 resp.Mappings.UserMappings.EmailAddress,
		"user_mappings_department":                    resp.Mappings.UserMappings.Department,
		"user_mappings_building":                      resp.Mappings.UserMappings.Building,
		"user_mappings_room":                          resp.Mappings.UserMappings.Room,
		"user_mappings_phone":                         resp.Mappings.UserMappings.Phone,
		"user_mappings_position":                      resp.Mappings.UserMappings.Position,
		"user_mappings_uuid":                          resp.Mappings.UserMappings.UserUuid,
		"group_mappings_object_class_limitation":      resp.Mappings.GroupMappings.ObjectClassLimitation,
		"group_mappings_object_classes":               resp.Mappings.GroupMappings.ObjectClasses,
		"group_mappings_search_base":                  resp.Mappings.GroupMappings.SearchBase,
		"group_mappings_search_scope":                 resp.Mappings.GroupMappings.SearchScope,
		"group_mappings_id":                           resp.Mappings.GroupMappings.GroupID,
		"group_mappings_name":                         resp.Mappings.GroupMappings.GroupName,
		"group_mappings_uuid":                         resp.Mappings.GroupMappings.GroupUuid,
		"group_membership_mapping":                    resp.Mappings.MembershipMappings.GroupMembershipMapping,
	}

	if resp.Server.Keystore != nil {
		computedKeystoreFields := map[string]any{
			"keystore_file_name":       resp.Server.Keystore.FileName,
			"keystore_type":            resp.Server.Keystore.Type,
			"keystore_expiration_date": resp.Server.Keystore.ExpirationDate,
			"keystore_subject":         resp.Server.Keystore.Subject,
		}

		if password, ok := d.GetOk("keystore_password"); ok {
			settings["keystore_password"] = password
		}

		if fileBytes, ok := d.GetOk("keystore_file_bytes"); ok {
			settings["keystore_file_bytes"] = fileBytes
		}

		for k, v := range computedKeystoreFields {
			settings[k] = v
		}
	}

	for key, val := range settings {
		if err := d.Set(key, val); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
