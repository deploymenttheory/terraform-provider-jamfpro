# Resource Migration Guide

This guide explains how to migrate from resources that have been superseded in this provider using current HashiCorp Terraform guidance. Resources are often superseded with new versions when there are significant/complex changes to their schemas such that a state upgrade is deemed to pose significant risk.

A common example is when a resource is refactored to use a different version of the underlying API (for example changing from the Classic to Jamf Pro API). The new resource is typically named with a `_v{int}` suffix to indicate the major version change and signify the API endpoint version it was initially created for. A deprecation warning is added to the old resource to encourage users to migrate to the new version.

Users should plan to migrate to the new version as soon as possible to avoid running into issues when the old resource is eventually removed from the provider or support from the API it uses is removed.

## References (HashiCorp)

- Import blocks: <https://developer.hashicorp.com/terraform/language/block/import>
- Removed blocks: <https://developer.hashicorp.com/terraform/language/block/removed>
- Refactoring and moved blocks: <https://developer.hashicorp.com/terraform/language/modules/develop/refactoring>

## Recommended workflow (import + removed)

This is the safest option for changing resource types because it keeps the existing Jamf Pro objects in place and re-associates state to the resource type.

The example below shows how to migrate from `jamfpro_smart_computer_group` to `jamfpro_smart_computer_group_v2`. Adjust the resource types and arguments as needed for your specific migration.

### 1) Add new resources

Create new `jamfpro_smart_computer_group_v2` resources that mirror your existing configuration. This is typically a straightforward copy/paste of your current resource blocks with the resource type renamed and any v2-specific arguments adjusted.

### 2) Import existing objects into the v2 resources

Use `import` blocks to attach the existing Jamf Pro objects to the new v2 resource addresses.

Single resource:

```hcl
import {
  to = jamfpro_smart_computer_group_v2.example
  id = "123"
}

resource "jamfpro_smart_computer_group_v2" "example" {
  # v2 configuration
}
```

If you used `for_each`, use an `import` block with `for_each` as well. The `to` argument must be a resource address, not an ID. The ID goes in `id` (or `identity` if the provider requires it).

```hcl
locals {
  groups = {
    "finance" = { id = "101" }
    "hr"      = { id = "102" }
  }
}

resource "jamfpro_smart_computer_group_v2" "this" {
  for_each = local.groups
  # v2 configuration using each.value
}

import {
  for_each = local.groups
  to       = jamfpro_smart_computer_group_v2.this[each.key]
  id       = each.value.id
}
```

Run `terraform apply` to perform the imports.

### 3) Remove the old resources from state without destroying

After the v2 resources are imported and stable, remove the old resources from state using a `removed` block. This ensures Terraform stops managing the old resource addresses without destroying the actual Jamf Pro objects.

```hcl
removed {
  from = jamfpro_smart_computer_group.example

  lifecycle {
    destroy = false
  }
}
```

Run `terraform apply` again to remove the legacy addresses from state.

### 4) Clean up

Once the state is updated, remove the legacy `jamfpro_smart_computer_group` blocks and any temporary `import` or `removed` blocks you no longer need.

## Dealing with dependent resources

In practice, `jamfpro_smart_computer_group` resources are commonly referenced by `jamfpro_policy` via `scope.computer_group_ids`. During a migration to `jamfpro_smart_computer_group_v2`, update those references carefully so Terraform doesn’t plan an unintended replace of the policy.

### Recommended ordering

1. Add the new `jamfpro_smart_computer_group_v2` resource.
2. Import the existing Jamf Pro smart group into the v2 address.
3. Update `jamfpro_policy` to reference the v2 smart group ID.
4. Apply.
5. Only then remove the legacy smart group address from state using `removed { ... destroy = false }`.

### Example: policy scope referencing a smart computer group

If your policy currently references the legacy smart group:

```hcl
resource "jamfpro_smart_computer_group" "example" {
  # legacy configuration
}

resource "jamfpro_policy" "install_app" {
  # ...
  scope {
    computer_group_ids = [jamfpro_smart_computer_group.example.id]
    # ...
  }
}
```

During the migration window (when both resource blocks may temporarily exist), you can keep the policy stable by computing the group ID once and using it in `scope`:

```hcl
resource "jamfpro_smart_computer_group" "example" {
  # legacy configuration
}

resource "jamfpro_smart_computer_group_v2" "example" {
  # v2 configuration
}

locals {
  smart_group_id = coalesce(
    try(jamfpro_smart_computer_group_v2.example.id, null),
    try(jamfpro_smart_computer_group.example.id, null),
  )
}

resource "jamfpro_policy" "install_app" {
  # ...
  scope {
    computer_group_ids = [local.smart_group_id]
    # ...
  }
}
```

Once the migration is complete, simplify the policy to reference only `jamfpro_smart_computer_group_v2.example.id`.

## Quick verification checklist

- `terraform plan` shows no destroy actions for smart computer groups.
- State contains `jamfpro_smart_computer_group_v2` addresses only.
- Old `jamfpro_smart_computer_group` addresses are gone from state.

## `jamfpro_mac_application`: unordered collections (v0 to v1)

Schema version 1 automatically upgrades existing version 0 state. The collections
listed below are unordered sets; repeated identical elements are deduplicated.
HCL block syntax is unchanged, but numeric indexing into these collections is no
longer supported. Select by ID/key with a `for` expression instead.

`self_service_category`.

Back up state before upgrading. Do not use upgraded state with an older provider;
restore the pre-upgrade state and provider together if a rollback is necessary.

The upgrader retains the complete v0 schema, including unrelated nested values,
and accepts JSON and legacy flatmap state. No resource replacement, state removal,
or import is required. Other ordered lists and singleton blocks are unchanged.

## `jamfpro_enrollment_customization`: unordered collections (v0 to v1)

Schema version 1 automatically upgrades existing version 0 state. The collections
listed below are unordered sets; repeated identical elements are deduplicated.
HCL block syntax is unchanged, but numeric indexing into these collections is no
longer supported. Select by ID/key with a `for` expression instead.

`ldap_group_access`.

Back up state before upgrading. Do not use upgraded state with an older provider;
restore the pre-upgrade state and provider together if a rollback is necessary.

The upgrader retains the complete v0 schema, including unrelated nested values,
and accepts JSON and legacy flatmap state. No resource replacement, state removal,
or import is required. Other ordered lists and singleton blocks are unchanged.

## `jamfpro_static_computer_group`: unordered collections (v0 to v1)

Schema version 1 automatically upgrades existing version 0 state. The collections
listed below are unordered sets; repeated identical elements are deduplicated.
HCL block syntax is unchanged, but numeric indexing into these collections is no
longer supported. Select by ID/key with a `for` expression instead.

`assigned_computer_ids`.

Back up state before upgrading. Do not use upgraded state with an older provider;
restore the pre-upgrade state and provider together if a rollback is necessary.

The upgrader retains the complete v0 schema, including unrelated nested values,
and accepts JSON and legacy flatmap state. No resource replacement, state removal,
or import is required. Other ordered lists and singleton blocks are unchanged.

## `jamfpro_static_mobile_device_group`: unordered collections (v0 to v1)

Schema version 1 automatically upgrades existing version 0 state. The collections
listed below are unordered sets; repeated identical elements are deduplicated.
HCL block syntax is unchanged, but numeric indexing into these collections is no
longer supported. Select by ID/key with a `for` expression instead.

`assigned_mobile_device_ids`.

Back up state before upgrading. Do not use upgraded state with an older provider;
restore the pre-upgrade state and provider together if a rollback is necessary.

The upgrader retains the complete v0 schema, including unrelated nested values,
and accepts JSON and legacy flatmap state. No resource replacement, state removal,
or import is required. Other ordered lists and singleton blocks are unchanged.

## `jamfpro_account_group`: unordered collections (v0 to v1)

Schema version 1 automatically upgrades existing version 0 state. The collections
listed below are unordered sets; repeated identical elements are deduplicated.
HCL block syntax is unchanged, but numeric indexing into these collections is no
longer supported. Select by ID/key with a `for` expression instead.

`member_ids`.

Back up state before upgrading. Do not use upgraded state with an older provider;
restore the pre-upgrade state and provider together if a rollback is necessary.

The upgrader retains the complete v0 schema, including unrelated nested values,
and accepts JSON and legacy flatmap state. No resource replacement, state removal,
or import is required. Other ordered lists and singleton blocks are unchanged.

## Policy collection blocks: schema version 2

The following `jamfpro_policy` blocks now use unordered sets:

- `payloads.scripts`
- `payloads.packages.package`
- `payloads.account_maintenance.local_accounts.account`
- `payloads.account_maintenance.directory_bindings.binding`
- `self_service.self_service_category`

The provider automatically upgrades existing policy state from schema version 1
to version 2. Version 0 state first passes through the existing user-interaction
field rename, then the collection migration. No `terraform state rm`, import,
or policy recreation is required. Run `terraform plan`, review the result, and
apply the saved plan to persist the new state version.

Repeated HCL block syntax remains unchanged. Reordering these blocks no longer
causes a policy update. Identical elements are deduplicated; changing a field,
adding an element, or removing an element still produces an in-place update.
Script execution priority remains controlled by each script's `priority` field.
Single-configuration container blocks remain lists. `payloads.printers` and
`payloads.dock_items` also remain ordered lists, because Jamf Pro stores them in
the order sent and dock item order can affect "Add To End" and
"Add To Beginning" actions.

Expressions that select collection elements by numeric index must be updated.
For example, replace `jamfpro_policy.example.payloads[0].scripts[0].id` with a
selection by identity:

```hcl
output "selected_script_parameter" {
  value = one([
    for script in jamfpro_policy.example.payloads[0].scripts : script.parameter4
    if script.id == jamfpro_script.example.id
  ])
}
```

Back up state before upgrading. After applying with schema version 2, an older
provider cannot read the upgraded state; do not downgrade that state in place.

## `jamfpro_macos_configuration_profile_plist_generator`: unordered collections (v0 to v1)

Schema version 1 automatically upgrades existing version 0 state. The collections
listed below are unordered sets; repeated identical elements are deduplicated.
HCL block syntax is unchanged, but numeric indexing into these collections is no
longer supported. Select by ID/key with a `for` expression instead.

`setting`, `dictionary`.

Back up state before upgrading. Do not use upgraded state with an older provider;
restore the pre-upgrade state and provider together if a rollback is necessary.

Only `payloads.payload_content.setting` and its nested dictionary entry blocks
are sets. `payloads`, `payload_content`, and all plist arrays keep their existing
order semantics. Dictionary keys must be unique.

Use `array_json` instead of `value` or `dictionary` for an ordered array:

```hcl
setting {
  key        = "OrderedValues"
  array_json = jsonencode(["second", "first", "second"])
}
```

`array_json` is an optional string on settings and dictionary entries (up to the
existing six nested block levels). It supports JSON-compatible plist array
values and preserves element order and duplicates. Nested dictionaries inside
an array remain array elements. Reordering an array is a real configuration
change. The legacy string value conversion continues to recognize booleans and
integers. The terminal dictionary remains a string map.

The upgrader retains the complete v0 schema, including unrelated nested values,
and accepts JSON and legacy flatmap state. No resource replacement, state removal,
or import is required. Other ordered lists and singleton blocks are unchanged.

Arrays containing plist dates or binary data are rejected with an explicit error;
they are not coerced to JSON strings. Use the raw plist profile resource for
those payload types.

Array JSON numbers are normalized without reordering or deduplicating elements.
Integers retain their signed/unsigned 64-bit values; real numbers use the plist
64-bit floating-point representation and retain a decimal/exponent marker when
read back (for example, `1.0` remains a real). JSON `null`, integers outside the
supported 64-bit range, and non-finite/out-of-range real values are rejected
before apply rather than silently removed or coerced.

### Behaviour changes

**Nested dictionary scalar values are now sent.** In v0.42.0 and earlier, a
`dictionary` entry with a scalar `value` was sent to Jamf Pro as an empty
`<dict/>`, so the value was lost. This version sends the configured value. After
the upgrade, profiles that contain such entries show an in-place update on the
first plan, and applying it changes the payload deployed by Jamf Pro. Review the
diff before applying. For example, this entry was previously sent as
`<key>Enabled</key><dict/>` and is now sent as `<key>Enabled</key><true/>`:

```hcl
setting {
  key = "Options"
  dictionary {
    key   = "Enabled"
    value = "true"
  }
}
```

**Ambiguous settings are now rejected.** v0.42.0 accepted the following
configurations and silently discarded part of them. This version returns an
error at apply time instead:

- Two `setting` blocks in the same `payload_content`, or two `dictionary` entries
  in the same dictionary, with the same `key` (previously the last one won).
- A `setting` or `dictionary` entry that sets both `value` and `dictionary`
  (previously `value` was discarded).

Remove the duplicate key or the unused `value` before upgrading.

## `jamfpro_user_group`: unordered collections (v0 to v1)

Schema version 1 automatically upgrades existing version 0 state. The collections
listed below are unordered sets; repeated identical elements are deduplicated.
HCL block syntax is unchanged, but numeric indexing into these collections is no
longer supported. Select by ID/key with a `for` expression instead.

`assigned_user_ids`, `user_additions`, `user_deletions`.

Back up state before upgrading. Do not use upgraded state with an older provider;
restore the pre-upgrade state and provider together if a rollback is necessary.

`user_additions` and `user_deletions` are command sets, not desired membership.
Only newly added or changed commands are submitted. Removing a command does not
undo it, and refresh does not erase it when the API omits command fields.
Computed contact fields do not participate in command identity. Use
`assigned_user_ids` to manage complete membership; do not combine it (including an explicitly empty set)
with operation sets. A metadata-only update preserves current
membership and does not replay prior deletions. The constructor accepts numeric
string IDs without a type assertion panic.

The upgrader retains the complete v0 schema, including unrelated nested values,
and accepts JSON and legacy flatmap state. No resource replacement, state removal,
or import is required. Other ordered lists and singleton blocks are unchanged.

When `assigned_user_ids` is omitted, the provider reports the observed membership
as a computed set. Changes to operation blocks recompute that set. An explicit
empty set still manages the group as empty.
