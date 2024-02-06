package macosconfigurationprofiles

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/http_client"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const JamfProResourceMacOSConfigurationProfile = "macos_configuration_profile"

func ResourceJamfProMacOSConfigurationProfiles() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProMacOSConfigurationProfilesCreate,
		ReadContext:   ResourceJamfProMacOSConfigurationProfilesRead,
		UpdateContext: ResourceJamfProMacOSConfigurationProfilesUpdate,
		DeleteContext: ResourceJamfProMacOSConfigurationProfilesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the macOS configuration profile.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the configuration profile.",
			},
		},
	}
}

func constructJamfProMacOSConfigurationProfile(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceMacOSConfigurationProfile, error) {

	out := jamfpro.ResourceMacOSConfigurationProfile{
		General: jamfpro.MacOSConfigurationProfileSubsetGeneral{
			Name: d.Get("name").(string),
		},
	}

	return &out, nil
}

func ResourceJamfProMacOSConfigurationProfilesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics
	var creationResponse *jamfpro.ResponseMacOSConfigurationProfileCreationUpdate
	var apiErrorCode int
	resourceName := d.Get("name").(string)

	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemCreate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	out, err := constructJamfProMacOSConfigurationProfile(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceMacOSConfigurationProfile, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceMacOSConfigurationProfile)

	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateMacOSConfigurationProfile(out)
		if apiErr != nil {

			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPICreateFailedAfterRetry(subCtx, JamfProResourceMacOSConfigurationProfile, resourceName, apiErr.Error(), apiErrorCode)

			return retry.NonRetryableError(apiErr)
		}

		return nil
	})

	if err != nil {

		logging.LogAPICreateFailure(subCtx, JamfProResourceMacOSConfigurationProfile, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	idString := strconv.Itoa(creationResponse.ID)
	logging.LogAPICreateSuccess(subCtx, JamfProResourceMacOSConfigurationProfile, idString)
	d.SetId(idString)

	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProMacOSConfigurationProfilesRead(subCtx, d, meta)
		if len(readDiags) > 0 {

			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceMacOSConfigurationProfile, d.Id(), readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}

		return nil
	})

	if err != nil {

		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceMacOSConfigurationProfile, err.Error())
		diags = append(diags, diag.FromErr(err)...)
	} else {

		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceMacOSConfigurationProfile, d.Id())
	}

	return diags
}

func ResourceJamfProMacOSConfigurationProfilesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)

	var diags diag.Diagnostics
	var apiErrorCode int
	var resp *jamfpro.ResourceMacOSConfigurationProfile
	resourceID := d.Id()
	resourceIDString, convErr := strconv.Atoi(resourceID)
	if convErr != nil {
		return diag.FromErr(convErr)

	}

	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resp, apiErr = conn.GetMacOSConfigurationProfileByID(resourceIDString)
		if apiErr != nil {
			logging.LogFailedReadByID(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, apiErr.Error(), apiErrorCode)

			return retry.RetryableError(apiErr)
		}

		return nil
	})

	if err != nil {

		logging.LogTFStateRemovalWarning(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID)
		return diag.FromErr(err)
	}

	logging.LogAPIReadSuccess(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID)

	if err := d.Set("id", resourceID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if err := d.Set("name", resp.General.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

func ResourceJamfProMacOSConfigurationProfilesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemUpdate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceIDString, convErr := strconv.Atoi(resourceID)
	if convErr != nil {
		return diag.FromErr(convErr)

	}
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	constructedPayload, err := constructJamfProMacOSConfigurationProfile(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceMacOSConfigurationProfile, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceMacOSConfigurationProfile)

	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateMacOSConfigurationProfileByID(resourceIDString, constructedPayload)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.LogAPIUpdateFailureByID(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			_, apiErrByName := conn.UpdateMacOSConfigurationProfileByName(resourceName, constructedPayload)
			if apiErrByName != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErrByName.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIUpdateFailureByName(subCtx, JamfProResourceMacOSConfigurationProfile, resourceName, apiErrByName.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErrByName)
			}
		} else {
			logging.LogAPIUpdateSuccess(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, resourceName)
		}
		return nil
	})

	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProMacOSConfigurationProfilesRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceMacOSConfigurationProfile, resourceID, readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceMacOSConfigurationProfile, err.Error())
		return diag.FromErr(err)
	} else {
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceMacOSConfigurationProfile, resourceID)
	}

	return nil
}

func ResourceJamfProMacOSConfigurationProfilesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemDelete, hclog.Info)

	err := retry.RetryContext(subCtx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {

		apiErr := conn.DeleteDepartmentByID(resourceID)
		if apiErr != nil {
			if apiError, ok := apiErr.(*http_client.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPIDeleteFailureByID(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			apiErr = conn.DeleteDepartmentByName(resourceName)
			if apiErr != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErr.(*http_client.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIDeleteFailureByName(subCtx, JamfProResourceMacOSConfigurationProfile, resourceName, apiErr.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	logging.LogAPIDeleteSuccess(subCtx, JamfProResourceMacOSConfigurationProfile, resourceID, resourceName)

	d.SetId("")

	return nil
}
