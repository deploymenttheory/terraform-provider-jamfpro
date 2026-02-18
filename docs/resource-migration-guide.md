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
