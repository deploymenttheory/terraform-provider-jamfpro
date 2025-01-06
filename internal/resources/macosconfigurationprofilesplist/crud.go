package macosconfigurationprofilesplist

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	pliststruct "github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common/configurationprofiles/plist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"howett.net/plist"
)

// Create requires a mutex need to lock Create requests during parallel runs
// var mu sync.Mutex

// resourceJamfProMacOSConfigurationProfilesPlistCreate is responsible for creating a new Jamf Pro macOS Configuration Profile in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func resourceJamfProMacOSConfigurationProfilesPlistCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// mu.Lock()
	// defer mu.Unlock()

	resource, err := constructJamfProMacOSConfigurationProfilePlist(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro macOS Configuration Profile: %v", err))
	}

	var creationResponse *jamfpro.ResponseMacOSConfigurationProfileCreationUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateMacOSConfigurationProfile(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro macOS Configuration Profile '%s' after retries: %v", resource.General.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	return append(diags, resourceJamfProMacOSConfigurationProfilesPlistReadNoCleanup(ctx, d, meta)...)
}

// resourceJamfProMacOSConfigurationProfilesPlistRead is responsible for reading the current state of a Jamf Pro config profile Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func resourceJamfProMacOSConfigurationProfilesPlistRead(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	resourceID := d.Id()
	var diags diag.Diagnostics

	var response *jamfpro.ResourceMacOSConfigurationProfile
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetMacOSConfigurationProfileByID(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, common.HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, updateState(d, response)...)
}

// resourceJamfProMacOSConfigurationProfilesPlistReadWithCleanup reads the resource with cleanup enabled
func resourceJamfProMacOSConfigurationProfilesPlistReadWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProMacOSConfigurationProfilesPlistRead(ctx, d, meta, true)
}

// resourceJamfProMacOSConfigurationProfilesPlistReadNoCleanup reads the resource with cleanup disabled
func resourceJamfProMacOSConfigurationProfilesPlistReadNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceJamfProMacOSConfigurationProfilesPlistRead(ctx, d, meta, false)
}

// resourceJamfProMacOSConfigurationProfilesPlistUpdate is responsible for updating an existing Jamf Pro config profile on the remote system.
func resourceJamfProMacOSConfigurationProfilesPlistUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Get existing profile to get current UUIDs
	existingProfile, err := client.GetMacOSConfigurationProfileByID(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get existing profile: %v", err))
	}

	// Extract existing UUIDs map from existing profile
	existingUUIDs := make(map[string]string)
	unescapedExistingPayload := html.UnescapeString(existingProfile.General.Payloads)
	var existingConfig pliststruct.ConfigurationProfile
	if err := plist.NewDecoder(strings.NewReader(unescapedExistingPayload)).Decode(&existingConfig); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode existing plist: %v", err))
	}

	// Store root UUID
	existingUUIDs["root"] = existingConfig.PayloadUUID

	// Store PayloadContent UUIDs
	for _, content := range existingConfig.PayloadContent {
		existingUUIDs[content.PayloadDisplayName] = content.PayloadUUID
	}

	// Construct new profile
	resource, err := constructJamfProMacOSConfigurationProfilePlist(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct profile for update: %v", err))
	}

	// Unescape the XML payload
	unescapedPayload := html.UnescapeString(resource.General.Payloads)

	// Parse new payload
	var newConfig pliststruct.ConfigurationProfile
	if err := plist.NewDecoder(strings.NewReader(unescapedPayload)).Decode(&newConfig); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode new plist: %v", err))
	}

	// Update root UUID
	if rootUUID, exists := existingUUIDs["root"]; exists {
		newConfig.PayloadUUID = rootUUID
	}

	// Update PayloadContent UUIDs
	for i, content := range newConfig.PayloadContent {
		if uuid, exists := existingUUIDs[content.PayloadDisplayName]; exists {
			newConfig.PayloadContent[i].PayloadUUID = uuid
			newConfig.PayloadContent[i].PayloadIdentifier = uuid
		}
	}

	// Encode back to plist
	var buf bytes.Buffer
	encoder := plist.NewEncoder(&buf)
	encoder.Indent("    ")
	if err := encoder.Encode(newConfig); err != nil {
		return diag.FromErr(fmt.Errorf("failed to encode updated plist: %v", err))
	}

	// Escape special characters for XML
	resource.General.Payloads = html.EscapeString(buf.String())

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateMacOSConfigurationProfileByID(resourceID, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update profile '%s' (ID: %s): %v", resource.General.Name, resourceID, err))
	}

	return append(diags, resourceJamfProMacOSConfigurationProfilesPlistReadNoCleanup(ctx, d, meta)...)
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
