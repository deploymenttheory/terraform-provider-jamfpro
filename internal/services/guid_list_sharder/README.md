## jamfpro_guid_list_sharder

This package provides a Jamf Pro-aware sharding helper based on https://github.com/deploymenttheory/go-jamf-guid-sharder:

- Fetch IDs from Jamf Pro (computers, mobile devices, users, and group membership)
- Apply exclusions and reservations
- Distribute IDs into deterministic shards using one of:
  - round-robin
  - percentage
  - size
  - rendezvous (HRW hashing)

Implementation goal: expose this as a Terraform data source that returns a map
of shard name -> IDs. Users can materialize these IDs into Jamf Pro static
groups using existing static group resources.
