---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "jamfpro_script Resource - terraform-provider-jamfpro"
subcategory: ""
description: |-
  
---

# jamfpro_script (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Display name for the script.
- `priority` (String) Execution priority of the script (BEFORE, AFTER, AT_REBOOT).
- `script_contents` (String) Contents of the script. Must be non-compiled and in an accepted format.

### Optional

- `category_id` (String) Script Category
- `info` (String) Information to display to the administrator when the script is run.
- `notes` (String) Notes to display about the script (e.g., who created it and when it was created).
- `os_requirements` (String) The script can only be run on computers with these operating system versions. Each version must be separated by a comma (e.g., 10.11, 15, 16.1).
- `parameter10` (String) Script parameter label 10
- `parameter11` (String) Script parameter label 11
- `parameter4` (String) Script parameter label 4
- `parameter5` (String) Script parameter label 5
- `parameter6` (String) Script parameter label 6
- `parameter7` (String) Script parameter label 7
- `parameter8` (String) Script parameter label 8
- `parameter9` (String) Script parameter label 9
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The Jamf Pro unique identifier (ID) of the script.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)
