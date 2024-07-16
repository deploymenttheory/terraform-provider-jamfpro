package common

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type sdkStructPayoadConstructorFunc[T any] func(*schema.ResourceData) (*T, error)

type sdkCreateUpdateFunc[PayloadType any, ResponseType any] func(Payload *PayloadType) (*ResponseType, error)

type sdkGetFunc[responseType any] func(resourceID string) (*responseType, error)

type sdkDeleteFunc func(resourceID string) error

type providerStateFunc[resourceType any] func(d *schema.ResourceData, resource *resourceType) diag.Diagnostics

type providerReadFunc func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics

func CreateUpdate[sdkPayloadType any, sdkResponseType any](
	ctx context.Context,
	d *schema.ResourceData,
	meta interface{},
	construct sdkStructPayoadConstructorFunc[sdkPayloadType],
	serverOutcomeFunc sdkCreateUpdateFunc[sdkPayloadType, sdkResponseType],
	reader providerReadFunc,
) diag.Diagnostics {

	var diags diag.Diagnostics

	payload, err := construct(d)
	loggingTypeName := reflect.TypeOf(payload).Name()

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct %s: %v", loggingTypeName, err))
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
		return append(diags, diag.FromErr(fmt.Errorf("failed to create %s after retries: %v", loggingTypeName, err))...)
	}

	IdField, err := getIDField(outcomeResponse)
	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("Error getting ID field from response: %v", err))...)
	}

	d.SetId(IdField)

	return append(diags, reader(ctx, d, meta)...)
}

func Read[sdkResponseType any](
	ctx context.Context,
	d *schema.ResourceData,
	meta interface{},
	cleanup bool,
	getter sdkGetFunc[sdkResponseType],
	stator providerStateFunc[sdkResponseType],
) diag.Diagnostics {

	var diags diag.Diagnostics
	resourceID := d.Id()

	var response *sdkResponseType
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = getter(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, HandleResourceNotFoundError(err, d, cleanup)...)
	}

	return append(diags, stator(d, response)...)
}

func Delete(
	ctx context.Context,
	d *schema.ResourceData,
	meta interface{},
	deleter sdkDeleteFunc,

) diag.Diagnostics {

	var diags diag.Diagnostics
	resourceID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		var apiErr error
		apiErr = deleter(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Building '%s' (ID: %s) after retries: %v", d.Get("name").(string), resourceID, err))
	}

	d.SetId("")

	return diags
}
