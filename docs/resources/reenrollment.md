---
page_title: "jamfpro_reenrollment"
description: |-
  
---

# jamfpro_reenrollment (Resource)


## Example Usage
```terraform
resource "jamfpro_reenrollment" "settings" {
  flush_location_information         = true
  flush_location_information_history = true
  flush_policy_history               = true
  flush_extension_attributes         = true
  flush_software_update_plans        = true
  flush_mdm_queue                    = "DELETE_EVERYTHING" // required, valid values: DELETE_NOTHING, DELETE_ERRORS, DELETE_EVERYTHING_EXCEPT_ACKNOWLEDGED, DELETE_EVERYTHING
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `flush_mdm_queue` (String) Clears computer and mobile device information from the Management History category on the History tab in inventory information during re-enrollment. Valid values are DELETE_NOTHING, DELETE_ERRORS, DELETE_EVERYTHING_EXCEPT_ACKNOWLEDGED, or DELETE_EVERYTHING.

### Optional

- `flush_extension_attributes` (Boolean) Clears all values for extension attributes from computer and mobile device inventory information during re-enrollment. This does not apply to extension attributes populated by scripts or Directory Service Attribute Mapping.
- `flush_location_information` (Boolean) Clears computer and mobile device information from the User and Location category on the Inventory tab in inventory information during re-enrollment.
- `flush_location_information_history` (Boolean) Clears computer and mobile device information from the User and Location History category on the History tab in inventory information during re-enrollment.
- `flush_policy_history` (Boolean) Clears the logs for policies that ran on the computer and clears computer information from the Policy Logs category on the History tab in inventory information during re-enrollment.
- `flush_software_update_plans` (Boolean) Clears all values for software update plans from computer and mobile device inventory information during re-enrollment.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)