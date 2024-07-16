package common

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Constructor[T any] func(*schema.ResourceData) (*T, error)

type SdkCreator[PayloadType any, ResponseType struct{ ID any }] func(Payload *PayloadType) (*ResponseType, error)

type ReaderFunc func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics

type ResourceInterface[T any] interface {
	Constructor(*schema.ResourceData) (*T, error)
}

func Create[SdkPayloadType any, SdkResponseType struct{ ID any }](
	ctx context.Context,
	d *schema.ResourceData,
	meta interface{},
	constructor Constructor[SdkPayloadType],
	SdkCreatorFunc SdkCreator[SdkPayloadType, SdkResponseType],
	reader ReaderFunc,

) diag.Diagnostics {

	var diags diag.Diagnostics

	resource, err := constructor(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Policy: %v", err))
	}

	var creationResponse *SdkResponseType
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = SdkCreatorFunc(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return append(diags, diag.FromErr(fmt.Errorf("failed to create Jamf Pro Policy after retries: %v", err))...)
	}

	IdField, err := getIDField(creationResponse)
	d.SetId(IdField)

	return append(diags, reader(ctx, d, meta)...)
}

func getIDField(response interface{}) (string, error) {
	v := reflect.ValueOf(response).Elem()
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return "", fmt.Errorf("ID field not found in response")
	}
	return idField.Interface().(string), nil
}
