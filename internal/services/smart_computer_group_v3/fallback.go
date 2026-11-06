package smart_computer_group_v3

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	apiresponse "github.com/deploymenttheory/go-api-http-client/response"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
)

func apiError(err error) *apiresponse.APIError {
	if err == nil {
		return nil
	}
	var result *apiresponse.APIError
	if errors.As(err, &result) {
		return result
	}
	return nil
}

var (
	errMissingGroupID  = errors.New("jamf pro returned no valid group ID after creating the group")
	errMissingResponse = errors.New("jamf pro returned an empty response")
)

func isNotFound(err error) bool {
	apiErr := apiError(err)
	return apiErr != nil && (apiErr.StatusCode == http.StatusNotFound || apiErr.StatusCode == http.StatusGone)
}

func modernWriteRefused(err error) bool {
	apiErr := apiError(err)
	if apiErr == nil || apiErr.StatusCode != http.StatusBadRequest {
		return false
	}
	for _, detail := range apiErr.Errors {
		criterionRefused := strings.Contains(detail.Description, "The criterion ")
		operatorRefused := strings.Contains(detail.Description, "The operator ")
		if strings.EqualFold(detail.Code, "INVALID_FIELD") && strings.Contains(detail.Description, "is not valid") &&
			(criterionRefused || operatorRefused) {
			return true
		}
	}
	return false
}

func (r *smartComputerGroupV3FrameworkResource) createGroup(
	ctx context.Context, input *jamfpro.ResourceSmartComputerGroupV3,
) (string, bool, error) {
	if err := ctx.Err(); err != nil {
		return "", false, fmt.Errorf("create smart computer group: %w", err)
	}
	created, modernErr := r.client.CreateSmartComputerGroupV3WithContext(ctx, *input)
	if modernErr == nil {
		if created == nil || !positiveID.MatchString(created.ID) {
			return "", false, errMissingGroupID
		}
		return created.ID, false, nil
	}
	if !modernWriteRefused(modernErr) {
		return "", false, fmt.Errorf("create smart computer group: %w", modernErr)
	}
	classic, err := constructClassic(input)
	if err != nil {
		return "", true, err
	}
	if err := ctx.Err(); err != nil {
		return "", true, fmt.Errorf("create smart computer group through Classic: %w", err)
	}
	result, err := r.client.CreateComputerGroupWithContext(ctx, classic)
	if err != nil {
		return "", true, fmt.Errorf("v3 refused the criteria (%w); Classic create failed: %w", modernErr, err)
	}
	if result == nil || result.ID <= 0 {
		return "", true, errMissingGroupID
	}
	return strconv.Itoa(result.ID), true, nil
}

func (r *smartComputerGroupV3FrameworkResource) updateGroup(
	ctx context.Context, id string, input *jamfpro.ResourceSmartComputerGroupV3,
) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, fmt.Errorf("update smart computer group: %w", err)
	}
	_, modernErr := r.client.UpdateSmartComputerGroupByIDV3WithContext(ctx, id, *input)
	if modernErr == nil {
		return false, nil
	}
	if !modernWriteRefused(modernErr) {
		return false, fmt.Errorf("update smart computer group: %w", modernErr)
	}
	classic, err := constructClassic(input)
	if err != nil {
		return true, err
	}
	if err := ctx.Err(); err != nil {
		return true, fmt.Errorf("update smart computer group through Classic: %w", err)
	}
	if _, err := r.client.UpdateComputerGroupByIDWithContext(ctx, id, classic); err != nil {
		return true, fmt.Errorf("v3 refused the criteria (%w); Classic update failed: %w", modernErr, err)
	}
	return true, nil
}
