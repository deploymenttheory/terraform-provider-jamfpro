package ssocertificate

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// create is responsible for creating a new Jamf Pro SSO Certificate in the remote system.
func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, apiErr := client.CreateSSOCertificate()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}

		if err := setKeystoreData(d, response.Keystore); err != nil {
			return retry.NonRetryableError(fmt.Errorf("failed to set keystore data: %v", err))
		}

		if response.KeystoreDetails != nil {
			if err := setKeystoreDetails(d, response.KeystoreDetails); err != nil {
				return retry.NonRetryableError(fmt.Errorf("failed to set keystore details: %v", err))
			}
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create SSO certificate: %v", err))
	}

	d.SetId("jamfpro_sso_certificate_singleton")

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// read is responsible for reading the current state of a Jamf Pro SSO Certificate from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	d.SetId("jamfpro_sso_certificate_singleton")

	var response *jamfpro.ResourceSSOKeystoreResponse
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetSSOCertificate()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	if err := setKeystoreData(d, response.Keystore); err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("error setting keystore data: %v", err))...)
	}

	if response.KeystoreDetails != nil {
		if err := setKeystoreDetails(d, response.KeystoreDetails); err != nil {
			return append(diags, diag.FromErr(fmt.Errorf("error setting keystore details: %v", err))...)
		}
	}

	return append(diags, updateState(d, response)...)
}

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// delete is responsible for deleting a Jamf Pro SSO Certificate.
func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteSSOCertificate()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("failed to delete SSO certificate: %v", err))...)
	}

	d.SetId("")
	return diags
}

// setKeystoreData sets the keystore data in the resource data
func setKeystoreData(d *schema.ResourceData, keystore jamfpro.ResourceSSOCertKeystore) error {
	if err := d.Set("keystore", []interface{}{
		map[string]interface{}{
			"key":                 keystore.Key,
			"type":                keystore.Type,
			"keystore_file_name":  keystore.KeystoreFileName,
			"keystore_setup_type": keystore.KeystoreSetupType,
			"keys":                flattenCertKeys(keystore.Keys),
		},
	}); err != nil {
		return fmt.Errorf("error setting keystore: %v", err)
	}
	return nil
}

// setKeystoreDetails sets the keystore details in the resource data
func setKeystoreDetails(d *schema.ResourceData, details *jamfpro.ResourceSSOKeystoreDetails) error {
	if err := d.Set("keystore_details", []interface{}{
		map[string]interface{}{
			"keys":          details.Keys,
			"issuer":        details.Issuer,
			"subject":       details.Subject,
			"expiration":    details.Expiration,
			"serial_number": details.SerialNumber,
		},
	}); err != nil {
		return fmt.Errorf("error setting keystore details: %v", err)
	}
	return nil
}

// flattenCertKeys flattens the cert keys into a slice of interfaces
func flattenCertKeys(keys []jamfpro.ResourceCertKey) []interface{} {
	var result []interface{}
	for _, key := range keys {
		result = append(result, map[string]interface{}{
			"id":    key.ID,
			"valid": key.Valid,
		})
	}
	return result
}
