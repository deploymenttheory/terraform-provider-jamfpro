package mobile_device_inventory

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestSetMobileDeviceInventoryData_ExtensionAttributes(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DataSourceJamfProMobileDeviceInventory().Schema, map[string]interface{}{})

	device := &jamfpro.ResourceMobileDeviceInventory{
		MobileDeviceId: "123",
		ExtensionAttributes: []jamfpro.SharedResourceExtensionAttributeValue{
			{
				DisplayName: "Multiple Values",
				Value:       []string{"value1", "value2", "value3"},
			},
			{
				DisplayName: "Single Value",
				Value:       []string{"single"},
			},
			{
				DisplayName: "Empty Values",
				Value:       []string{},
			},
			{
				DisplayName: "Nil Values",
				Value:       nil,
			},
		},
	}

	diags := setMobileDeviceInventoryData(d, device)

	assert.False(t, diags.HasError(), "Expected no errors but got: %v", diags)

	extAttrs := d.Get("extension_attributes").([]interface{})
	assert.Equal(t, 4, len(extAttrs), "Expected 4 extension attributes")

	attr0 := extAttrs[0].(map[string]interface{})
	assert.Equal(t, "Multiple Values", attr0["display_name"])
	values0 := attr0["values"].([]interface{})
	assert.Equal(t, 3, len(values0))
	assert.Equal(t, "value1", values0[0])
	assert.Equal(t, "value2", values0[1])
	assert.Equal(t, "value3", values0[2])

	attr1 := extAttrs[1].(map[string]interface{})
	assert.Equal(t, "Single Value", attr1["display_name"])
	values1 := attr1["values"].([]interface{})
	assert.Equal(t, 1, len(values1))
	assert.Equal(t, "single", values1[0])

	attr2 := extAttrs[2].(map[string]interface{})
	assert.Equal(t, "Empty Values", attr2["display_name"])
	values2 := attr2["values"].([]interface{})
	assert.Equal(t, 0, len(values2))

	attr3 := extAttrs[3].(map[string]interface{})
	assert.Equal(t, "Nil Values", attr3["display_name"])
	values3 := attr3["values"]
	if values3 != nil {
		assert.Equal(t, 0, len(values3.([]interface{})))
	}
}

func TestSetMobileDeviceInventoryData_NoExtensionAttributes(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DataSourceJamfProMobileDeviceInventory().Schema, map[string]interface{}{})

	device := &jamfpro.ResourceMobileDeviceInventory{
		MobileDeviceId:      "123",
		ExtensionAttributes: []jamfpro.SharedResourceExtensionAttributeValue{},
	}

	diags := setMobileDeviceInventoryData(d, device)

	assert.False(t, diags.HasError(), "Expected no errors but got: %v", diags)

	extAttrs := d.Get("extension_attributes")
	if extAttrs != nil {
		extAttrsList := extAttrs.([]interface{})
		assert.Equal(t, 0, len(extAttrsList), "Expected no extension attributes")
	}
}
