data "jamfpro_guid_list_sharder" "rollout" {
  source_type = "computer_group_membership"
  group_id    = "456" # "eligible" group

  strategy          = "percentage"
  shard_percentages = [10, 20, 70]
  seed              = "macos-14.5-rollout"

  exclude_group_id = "123" # "do-not-touch" group

  reserve_group_id_shard_0 = "789" # "pilot" group
}
