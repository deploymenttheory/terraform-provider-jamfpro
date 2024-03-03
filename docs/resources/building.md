---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "jamfpro_building Resource - terraform-provider-jamfpro"
subcategory: ""
description: |-
  
---

# jamfpro_building (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the building.

### Optional

- `city` (String) The city in which the building is located.
- `country` (String) The country in which the building is located.
- `state_province` (String) The state or province in which the building is located.
- `street_address1` (String) The first line of the street address of the building.
- `street_address2` (String) The second line of the street address of the building.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `zip_postal_code` (String) The ZIP or postal code of the building.

### Read-Only

- `id` (String) The unique identifier of the building.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)