package macosconfigurationprofilesplist

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/plist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	hplist "howett.net/plist"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return common.Create(
		ctx,
		d,
		meta,
		construct,
		meta.(*jamfpro.Client).CreateMacOSConfigurationProfile,
		readNoCleanup,
	)
}

// read is responsible for reading the current state of a Jamf Pro Script Resource from the remote system.
func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	return common.Read(
		ctx,
		d,
		meta,
		cleanup,
		meta.(*jamfpro.Client).GetMacOSConfigurationProfileByID,
		updateState,
	)
}

// readWithCleanup reads the resource with cleanup enabled
func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

// readNoCleanup reads the resource with cleanup disabled
func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

// update is responsible for updating an existing Jamf Pro config profile on the remote system.
// The function:
// 1. Gets the existing profile to extract the payload information and UUID
// 2. Unescapes and decodes the XML-encoded plist data
// 3. Stores the current UUID
// 4. Constructs the new resource
// 5. Unescapes and decodes the new payload
// 6. Recursively replaces the UUID throughout the payload structure
// 7. Encodes and escapes the modified payload
// 8. Performs the update operation
// 9. Returns any diagnostics
// ref: https://grahamrpugh.com/2020/04/27/managing-profiles-between-multiple-jamf-instances.html
func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	var existingProfile *jamfpro.ResourceMacOSConfigurationProfile
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		existingProfile, apiErr = client.GetMacOSConfigurationProfileByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get existing Jamf Pro macOS Configuration Profile (ID: %s): %v", resourceID, err))
	}

	unescapedPayload := html.UnescapeString(existingProfile.General.Payloads)

	var configProfile plist.ConfigurationProfile
	decoder := hplist.NewDecoder(strings.NewReader(unescapedPayload))
	if err := decoder.Decode(&configProfile); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode existing profile payload: %v", err))
	}

	currentUUID := configProfile.PayloadUUID

	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro macOS Configuration Profile for update: %v", err))
	}

	unescapedNewPayload := html.UnescapeString(resource.General.Payloads)
	var newConfigProfile plist.ConfigurationProfile
	decoder = hplist.NewDecoder(strings.NewReader(unescapedNewPayload))
	if err := decoder.Decode(&newConfigProfile); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode new profile payload: %v", err))
	}

	var replaceUUIDs func(content []plist.PayloadContent)
	replaceUUIDs = func(content []plist.PayloadContent) {
		for i := range content {
			if content[i].PayloadUUID == newConfigProfile.PayloadUUID {
				content[i].PayloadUUID = currentUUID
			}
			if nestedContent, ok := content[i].ConfigurationItems["PayloadContent"].([]plist.PayloadContent); ok {
				replaceUUIDs(nestedContent)
			}
		}
	}

	newConfigProfile.PayloadUUID = currentUUID
	replaceUUIDs(newConfigProfile.PayloadContent)

	var payloadBuffer bytes.Buffer
	encoder := hplist.NewEncoder(&payloadBuffer)
	encoder.Indent("\t")
	if err := encoder.Encode(newConfigProfile); err != nil {
		return diag.FromErr(fmt.Errorf("failed to encode modified payload: %v", err))
	}

	formattedPayload := payloadBuffer.String()
	log.Printf("[DEBUG] Constructed plist payload before XML escaping:\n%s\n", formattedPayload)

	resource.General.Payloads = html.EscapeString(formattedPayload)
	log.Printf("[DEBUG] Final escaped XML payload:\n%s\n", resource.General.Payloads)

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateMacOSConfigurationProfileByID(resourceID, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro macOS Configuration Profile '%s' (ID: %s) after retries: %v", resource.General.Name, resourceID, err))
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

// resourceJamfProMacOSConfigurationProfilesPlistDelete is responsible for deleting a Jamf Pro config profile.
func resourceJamfProMacOSConfigurationProfilesPlistDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceName := d.Get("name").(string)

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteMacOSConfigurationProfileByID(resourceID)
		if apiErr != nil {
			apiErrByName := client.DeleteMacOSConfigurationProfileByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro macOS Configuration Profile '%s' (ID: %s) after retries: %v", resourceName, resourceID, err))
	}

	d.SetId("")

	return diags
}
