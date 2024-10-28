---
page_title: "jamfpro_smart_mobile_device_group"
description: |-
  
---

# jamfpro_smart_mobile_device_group (Data Source)


## Example Usage
```terraform
data "jamfpro_smart_mobile_device_group" "jamfpro_smart_mobile_device_group_001_data" {
  id = jamfpro_smart_mobile_device_group.jamfpro_smart_mobile_device_group_001.id
}

output "jamfpro_jamfpro_smart_mobile_device_group_001_id" {
  value = data.jamfpro_smart_mobile_device_group.jamfpro_smart_mobile_device_group_001_data.id
}

output "jamfpro_jamfpro_smart_mobile_device_groups_001_name" {
  value = data.jamfpro_smart_mobile_device_group.jamfpro_smart_mobile_device_group_001_data.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The unique identifier of the Jamf Pro Smart mobile group.

### Read-Only

- `name` (String) The unique name of the Jamf Pro Smart mobile group.