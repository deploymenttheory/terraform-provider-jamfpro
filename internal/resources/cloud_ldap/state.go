package cloud_ldap

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// updateState updates the Terraform state with the provided ResourceCloudLdap object.
func updateState(d *schema.ResourceData, resp *jamfpro.ResourceCloudLdap) diag.Diagnostics {
	var diags diag.Diagnostics

	cloudIdpCommon := []interface{}{
		map[string]interface{}{
			"provider_name": resp.CloudIdPCommon.ProviderName,
			"display_name":  resp.CloudIdPCommon.DisplayName,
		},
	}
	if err := d.Set("cloud_idp_common", cloudIdpCommon); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	server := setCloudLdapServer(d, resp.Server)
	if err := d.Set("server", []interface{}{server}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if resp.Mappings != nil {
		mappings := setCloudLdapMappings(resp.Mappings)
		if err := d.Set("mappings", []interface{}{mappings}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	} else {
		if err := d.Set("mappings", []interface{}{}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}

// setCloudLdapServer flattens a CloudLdapServer object into a format suitable for Terraform state.
func setCloudLdapServer(d *schema.ResourceData, server *jamfpro.CloudLdapServer) map[string]interface{} {
	serverMap := map[string]interface{}{
		"enabled":            server.Enabled,
		"use_wildcards":      server.UseWildcards,
		"connection_type":    server.ConnectionType,
		"server_url":         server.ServerUrl,
		"domain_name":        server.DomainName,
		"port":               server.Port,
		"connection_timeout": server.ConnectionTimeout,
		"search_timeout":     server.SearchTimeout,
		"membership_calculation_optimization_enabled": server.MembershipCalculationOptimizationEnabled,
	}

	if server.Keystore != nil {
		keystoreMap := map[string]interface{}{
			"file_name":       server.Keystore.FileName,
			"type":            server.Keystore.Type,
			"expiration_date": server.Keystore.ExpirationDate,
			"subject":         server.Keystore.Subject,
		}

		if old, ok := d.GetOk("server.0.keystore.0"); ok {
			oldKeystore := old.(map[string]interface{})
			if password, ok := oldKeystore["password"]; ok {
				keystoreMap["password"] = password
			}
			if fileBytes, ok := oldKeystore["file_bytes"]; ok {
				keystoreMap["file_bytes"] = fileBytes
			}
		}

		serverMap["keystore"] = []interface{}{keystoreMap}
	}

	return serverMap
}

// setCloudLdapMappings flattens a CloudLdapMappings object into a format suitable for Terraform state.
func setCloudLdapMappings(mappings *jamfpro.CloudLdapMappings) map[string]interface{} {
	userMappings := []interface{}{
		map[string]interface{}{
			"object_class_limitation": mappings.UserMappings.ObjectClassLimitation,
			"object_classes":          mappings.UserMappings.ObjectClasses,
			"search_base":             mappings.UserMappings.SearchBase,
			"search_scope":            mappings.UserMappings.SearchScope,
			"additional_search_base":  mappings.UserMappings.AdditionalSearchBase,
			"user_id":                 mappings.UserMappings.UserID,
			"username":                mappings.UserMappings.Username,
			"real_name":               mappings.UserMappings.RealName,
			"email_address":           mappings.UserMappings.EmailAddress,
			"department":              mappings.UserMappings.Department,
			"building":                mappings.UserMappings.Building,
			"room":                    mappings.UserMappings.Room,
			"phone":                   mappings.UserMappings.Phone,
			"position":                mappings.UserMappings.Position,
			"user_uuid":               mappings.UserMappings.UserUuid,
		},
	}

	groupMappings := []interface{}{
		map[string]interface{}{
			"object_class_limitation": mappings.GroupMappings.ObjectClassLimitation,
			"object_classes":          mappings.GroupMappings.ObjectClasses,
			"search_base":             mappings.GroupMappings.SearchBase,
			"search_scope":            mappings.GroupMappings.SearchScope,
			"group_id":                mappings.GroupMappings.GroupID,
			"group_name":              mappings.GroupMappings.GroupName,
			"group_uuid":              mappings.GroupMappings.GroupUuid,
		},
	}

	membershipMappings := []interface{}{
		map[string]interface{}{
			"group_membership_mapping": mappings.MembershipMappings.GroupMembershipMapping,
		},
	}

	return map[string]interface{}{
		"user_mappings":       userMappings,
		"group_mappings":      groupMappings,
		"membership_mappings": membershipMappings,
	}
}
