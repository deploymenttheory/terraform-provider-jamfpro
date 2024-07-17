package common

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Func Definitions
type payoadConstructorFunc[T any] func(*schema.ResourceData) (*T, error)

type sdkCreateUpdateFunc[PayloadType any, ResponseType any] func(Payload *PayloadType) (*ResponseType, error)

type sdkUpdateFunc[PayloadType any, ResponseType any] func(resourceID string, Payload *PayloadType) (*ResponseType, error)

type sdkGetFunc[responseType any] func(resourceID string) (*responseType, error)

type sdkDeleteFunc func(resourceID string) error

type providerStateFunc[resourceType any] func(d *schema.ResourceData, resource *resourceType) diag.Diagnostics

type providerReadFunc func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics

// Create Update
func CreateUpdate[sdkPayloadType any, sdkResponseType any](
	ctx context.Context,
	d *schema.ResourceData,
	meta interface{},
	construct payoadConstructorFunc[sdkPayloadType],
	serverOutcomeFunc sdkCreateUpdateFunc[sdkPayloadType, sdkResponseType],
	reader providerReadFunc,
) diag.Diagnostics {

	var diags diag.Diagnostics

	payload, err := construct(d)
	payloadtypeName := reflect.TypeOf(payload).Name()

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct %s: %v", payloadtypeName, err))
	}

	var outcomeResponse *sdkResponseType
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		outcomeResponse, apiErr = serverOutcomeFunc(payload)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("failed to create %s after retries: %v", payloadtypeName, err))...)
	}

	idField, err := getIDField(outcomeResponse)
	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("Error getting ID field from response: %v", err))...)
	}

	d.SetId(idField)

	return append(diags, reader(ctx, d, meta)...)
}

// Read
func Read[sdkResponseType any](
	ctx context.Context,
	d *schema.ResourceData,
	meta interface{},
	removeDeleteResourcesFromState bool,
	serverOutcomeFunc sdkGetFunc[sdkResponseType],
	providerStateFunc providerStateFunc[sdkResponseType],
) diag.Diagnostics {

	var diags diag.Diagnostics
	resourceID := d.Id()

	var response *sdkResponseType
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = serverOutcomeFunc(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, HandleResourceNotFoundError(err, d, removeDeleteResourcesFromState)...)
	}

	return append(diags, providerStateFunc(d, response)...)
}

// Delete
func Delete(ctx context.Context, d *schema.ResourceData, meta interface{}, serverOutcomeFunc sdkDeleteFunc) diag.Diagnostics {
	var diags diag.Diagnostics
	resourceID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		var apiErr error
		apiErr = serverOutcomeFunc(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro resource '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
