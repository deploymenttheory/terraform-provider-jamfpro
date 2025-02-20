package enrollmentcustomizations

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Handle image upload first if image source is provided
	imagePath, err := constructImageUpload(d)
	if err == nil && imagePath != "" {
		err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			uploadResponse, apiErr := client.UploadEnrollmentCustomizationsImage(imagePath)
			if apiErr != nil {
				return retry.RetryableError(apiErr)
			}
			// Store the URL in the schema for the main resource construction
			brandingSettings := d.Get("branding_settings").([]interface{})
			if len(brandingSettings) > 0 {
				settings := brandingSettings[0].(map[string]interface{})
				settings["icon_url"] = uploadResponse.Url
				brandingSettingsList := []interface{}{settings}
				if err := d.Set("branding_settings", brandingSettingsList); err != nil {
					return retry.NonRetryableError(fmt.Errorf("failed to set icon_url in schema: %v", err))
				}
			}
			return nil
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to upload enrollment customization image: %v", err))
		}
	}

	// Construct the resource from schema data (now including the uploaded image URL)
	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct enrollment customization: %v", err))
	}

	// Create the enrollment customization
	var response *jamfpro.ResponseEnrollmentCustomizationCreate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.CreateEnrollmentCustomization(*resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create enrollment customization: %v", err))
	}

	// Set the resource ID
	d.SetId(response.Id)

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

func read(ctx context.Context, d *schema.ResourceData, meta interface{}, cleanup bool) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	var response *jamfpro.ResourceEnrollmentCustomization
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetEnrollmentCustomizationByID(d.Id())
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

func readWithCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, true)
}

func readNoCleanup(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return read(ctx, d, meta, false)
}

func update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	// Construct the resource from schema data
	resource, err := construct(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct enrollment customization: %v", err))
	}

	// Handle image upload if image source has changed
	if d.HasChange("enrollment_customization_image_source") {
		if imgSource := d.Get("enrollment_customization_image_source").(string); imgSource != "" {
			err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				uploadResponse, apiErr := client.UploadEnrollmentCustomizationsImage(imgSource)
				if apiErr != nil {
					return retry.RetryableError(apiErr)
				}
				// Use the uploaded image URL directly in the resource
				resource.BrandingSettings.IconUrl = uploadResponse.Url
				return nil
			})
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to upload enrollment customization image: %v", err))
			}
		}
	}

	// Update the enrollment customization with the potentially new icon URL
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateEnrollmentCustomizationByID(d.Id(), *resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update enrollment customization: %v", err))
	}

	return append(diags, readNoCleanup(ctx, d, meta)...)
}

func delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteEnrollmentCustomizationByID(d.Id())
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete enrollment customization: %v", err))
	}

	d.SetId("")
	return diags
}
