package common

import (
	"context"
	"fmt"
	"reflect"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Type definitions for function signatures used in CRUD operations

// payoadConstructorFunc defines a function type that builds an SDK payload from ResourceData
type payoadConstructorFunc[T any] func(*schema.ResourceData) (*T, error)

// sdkCreateUpdateFunc defines a function type for SDK create/update operations
type sdkCreateUpdateFunc[PayloadType any, ResponseType any] func(Payload *PayloadType) (*ResponseType, error)

// sdkUpdateFunc defines a function type for SDK update operations requiring a resource ID
type sdkUpdateFunc[PayloadType any, ResponseType any] func(resourceID string, Payload *PayloadType) (*ResponseType, error)

// sdkGetFunc defines a function type for SDK read operations
type sdkGetFunc[responseType any] func(resourceID string) (*responseType, error)

// sdkDeleteFunc defines a function type for SDK delete operations
type sdkDeleteFunc func(resourceID string) error

// providerStateFunc defines a function type for updating Terraform state from SDK response
type providerStateFunc[resourceType any] func(d *schema.ResourceData, resource *resourceType) diag.Diagnostics

// providerReadFunc defines a function type for reading resource state
type providerReadFunc func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics

// Create is a shared helper that handles the creation of a new resource in Jamf Pro with retry logic and state management.
// It accepts generic types for the SDK payload and response to maintain type safety while being reusable.
//
// Parameters:
// - ctx: The context for the operation, used for timeouts and cancellation
// - d: The ResourceData containing the desired state from Terraform
// - meta: The provider meta object containing the authenticated client
// - construct: A function that builds the SDK payload from ResourceData
// - serverOutcomeFunc: The SDK function that performs the actual creation API call
// - reader: A function that reads back the resource state after creation
//
// Returns:
// - diag.Diagnostics containing any errors or warnings from the operation
func Create[sdkPayloadType any, sdkResponseType any](
	ctx context.Context,
	d *schema.ResourceData,
	meta any,
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
		return append(diags, diag.FromErr(fmt.Errorf("error getting ID field from response: %v", err))...)
	}

	d.SetId(idField.(string))

	return append(diags, reader(ctx, d, meta)...)
}

// Update is a shared helper that handles updating an existing resource in Jamf Pro with retry logic.
// It constructs the update payload, sends it to the API, and refreshes the state.
//
// Parameters:
// - ctx: The context for the operation, used for timeouts and cancellation
// - d: The ResourceData containing the desired state from Terraform
// - meta: The provider meta object containing the authenticated client
// - constructor: A function that builds the SDK payload from ResourceData
// - outcomeFunc: The SDK function that performs the actual update API call
// - reader: A function that reads back the resource state after update
//
// Returns:
// - diag.Diagnostics containing any errors or warnings from the operation
func Update[sdkPayloadType any, sdkResponseType any](
	ctx context.Context,
	d *schema.ResourceData,
	meta any,
	constructor payoadConstructorFunc[sdkPayloadType],
	outcomeFunc sdkUpdateFunc[sdkPayloadType, sdkResponseType],
	reader providerReadFunc,

) diag.Diagnostics {

	var diags diag.Diagnostics
	resourceID := d.Id()

	payload, err := constructor(d)
	payloadtypeName := reflect.TypeOf(payload).Name()

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro %s for update: %v", payloadtypeName, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := outcomeFunc(resourceID, payload)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf pro %s (ID: %s) after retries: %v", payloadtypeName, resourceID, err))
	}

	return append(diags, reader(ctx, d, meta)...)
}

// Read is a shared helper that retrieves the current state of a resource from Jamf Pro and updates the Terraform state.
// It includes retry logic and can optionally remove deleted resources from state.
//
// Parameters:
// - ctx: The context for the operation, used for timeouts and cancellation
// - d: The ResourceData to be updated with the current state
// - meta: The provider meta object containing the authenticated client
// - removeDeleteResourcesFromState: Whether to remove the resource from state if not found
// - serverOutcomeFunc: The SDK function that performs the actual read API call
// - providerStateFunc: A function that updates ResourceData with the API response
//
// Returns:
// - diag.Diagnostics containing any errors or warnings from the operation
func Read[sdkResponseType any](
	ctx context.Context,
	d *schema.ResourceData,
	meta any,
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
		return append(diags, errors.HandleResourceNotFoundError(err, d, removeDeleteResourcesFromState)...)
	}

	return append(diags, providerStateFunc(d, response)...)
}

// Delete is a shared helper that removes a resource from Jamf Pro with retry logic.
// It handles both the API deletion and clearing the resource ID from state.
//
// Parameters:
// - ctx: The context for the operation, used for timeouts and cancellation
// - d: The ResourceData containing the resource to delete
// - meta: The provider meta object containing the authenticated client
// - serverOutcomeFunc: The SDK function that performs the actual delete API call
//
// Returns:
// - diag.Diagnostics containing any errors or warnings from the operation
func Delete(ctx context.Context, d *schema.ResourceData, meta any, serverOutcomeFunc sdkDeleteFunc) diag.Diagnostics {
	var diags diag.Diagnostics
	resourceID := d.Id()

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := serverOutcomeFunc(resourceID)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		resourceName := ""
		if nameVal := d.Get("name"); nameVal != nil {
			resourceName = nameVal.(string)
		}
		if resourceName != "" {
			return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro resource '%s' (ID: %s) after retries: %v", resourceName, resourceID, err))
		}
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro resource (ID: %s) after retries: %v", resourceID, err))
	}

	d.SetId("")

	return diags
}
