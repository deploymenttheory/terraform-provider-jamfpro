package smart_mobile_device_group_v1

import (
	"context"
	"strconv"
	"testing"

	apiresponse "github.com/deploymenttheory/go-api-http-client/response"
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	schemahelpers "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema/helpers"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

type fakeGroupSDK struct{ err error }

func (f *fakeGroupSDK) CreateSmartMobileDeviceGroupV1(jamfpro.ResourceSmartMobileDeviceGroupV1) (*jamfpro.ResponseSmartMobileDeviceGroupCreateV1, error) {
	panic("unexpected create")
}

func (f *fakeGroupSDK) GetSmartMobileDeviceGroupByIDV1WithContext(context.Context, string) (*jamfpro.ResourceSmartMobileDeviceGroupV1, error) {
	return nil, f.err
}

func (f *fakeGroupSDK) UpdateSmartMobileDeviceGroupByIDV1(string, jamfpro.ResourceSmartMobileDeviceGroupV1) (*jamfpro.ResourceSmartMobileDeviceGroupV1, error) {
	panic("unexpected update")
}

func (f *fakeGroupSDK) DeleteSmartMobileDeviceGroupByIDV1WithContext(context.Context, string) error {
	return f.err
}

func testResource(t *testing.T, code int) (*smartMobileDeviceGroupV1FrameworkResource, tfsdk.State) {
	t.Helper()
	r := NewSmartMobileDeviceGroupV1FrameworkResource().(*smartMobileDeviceGroupV1FrameworkResource)
	r.client = &fakeGroupSDK{err: &apiresponse.APIError{StatusCode: code}}
	var response resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &response)
	plan := tfsdk.Plan{Schema: response.Schema}
	model := smartMobileDeviceGroupV1ResourceModel{
		ID: types.StringValue("42"), Name: types.StringValue("tf-testing-mobile-missing"), Description: types.StringNull(),
		SiteID: types.StringValue("-1"), Criteria: types.ListNull(schemahelpers.ObjectType),
		Timeouts: timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType,
		})},
	}
	require.False(t, plan.Set(context.Background(), &model).HasError())
	return r, tfsdk.State{Schema: response.Schema, Raw: plan.Raw}
}

func TestReadMissingGroupAndPermissionErrors(t *testing.T) {
	for _, code := range []int{404, 410, 403} {
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			r, initial := testResource(t, code)
			response := resource.ReadResponse{State: initial}
			r.Read(context.Background(), resource.ReadRequest{State: initial}, &response)
			require.Equal(t, code != 403, response.State.Raw.IsNull())
			require.Equal(t, code == 403, response.Diagnostics.HasError())
		})
	}
}

func TestDeleteMissingGroup(t *testing.T) {
	r, initial := testResource(t, 404)
	response := resource.DeleteResponse{State: initial}
	r.Delete(context.Background(), resource.DeleteRequest{State: initial}, &response)
	require.False(t, response.Diagnostics.HasError(), "%v", response.Diagnostics)
	require.True(t, response.State.Raw.IsNull())
}
