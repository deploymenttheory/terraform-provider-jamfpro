---
page_title: "jamfpro_guid_list_sharder Data Source"
description: |-
  Deterministically shard Jamf Pro IDs (computers, mobile devices, users, or group membership) for progressive rollouts.
---

# jamfpro_guid_list_sharder (Data Source)

Deterministically shards Jamf Pro IDs into named shards (e.g. `shard_0`, `shard_1`, ...) using one of:

- `round-robin`
- `percentage`
- `size`
- `rendezvous` (HRW hashing; best stability as fleet size changes)

This data source is intended to be combined with static group resources (e.g. `jamfpro_static_computer_group`) to materialize rollout waves.

## Example Usage

```hcl
data "jamfpro_guid_list_sharder" "computers" {
  source_type = "computer_inventory"

  strategy    = "rendezvous"
  shard_count = 4
  seed        = "macos-update-wave-2026-q1"

  exclude_ids = ["1001", "1002"]

  reserved_ids = {
    shard_0 = ["42"]
  }
}

resource "jamfpro_static_computer_group" "rollout" {
  for_each = data.jamfpro_guid_list_sharder.computers.shards

  name = "rollout-${each.key}"

  # static group expects ints
  assigned_computer_ids = [for id in each.value : tonumber(id)]
}

## Example: Shard A Group, Excluding Another Group

This example shards the membership of one computer group, but excludes any computers that
are members of another group (e.g. break/fix devices, executives, lab machines, etc.).

```hcl
data "jamfpro_guid_list_sharder" "rollout" {
  source_type = "computer_group_membership"
  group_id    = "456" # "eligible" group

  strategy          = "percentage"
  shard_percentages = [10, 20, 70]
  seed              = "macos-14.5-rollout"

  exclude_group_id = "123" # "do-not-touch" group
}
```

## Example: Force A Group Into The First Shard

To guarantee a set of devices are always in the first shard, reserve those IDs for `shard_0`.

```hcl
data "jamfpro_guid_list_sharder" "rollout" {
  source_type = "computer_inventory"

  strategy          = "percentage"
  shard_percentages = [10, 20, 70]
  seed              = "macos-14.5-rollout"

  reserve_group_id_shard_0 = "789" # "pilot" group
}
```
