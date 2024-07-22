// department_data_object.go
package departments

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Test_construct(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ResourceName: "jamfpro_department",
				ImportState:  true,
			},
		},
	})
	// type args struct {
	// 	d *schema.ResourceData
	// }
	// tests := []struct {
	// 	name    string
	// 	args    args
	// 	want    *jamfpro.ResourceDepartment
	// 	wantErr bool
	// }{
	// 	{
	// 		name:    "test 1",
	// 		args:    args{d: test_getResourceData()},
	// 		want:    test_getWant(),
	// 		wantErr: false,
	// 	},
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		got, err := construct(tt.args.d)
	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("construct() error = %v, wantErr %v", err, tt.wantErr)
	// 			return
	// 		}
	// 		if !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("construct() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}

func test_getWant() *jamfpro.ResourceDepartment {
	return &jamfpro.ResourceDepartment{
		Name: "Test Success",
	}
}

func test_getResourceData() *schema.ResourceData {
	out := schema.ResourceData{}
	out.Set("name", "success")
	return &out
}
