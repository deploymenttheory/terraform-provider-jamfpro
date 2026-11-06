# Smart computer group v3 live verification

Verified on 2026-11-06 against the authorized Jamf Pro lab running 11.33.3. The final provider build uses Jamf Pro SDK v1.54.0 and HTTP client v0.4.8, with no local dependency replacements. OAuth credentials were supplied through private environment variables and temporary files; no credentials are included here.

The test created uniquely named, unscoped `tf-testing-smart-v3-provider-` groups. A pre-existing string extension attribute was read without modification to exercise the known `has` validation defect.

| Check | Result |
| --- | --- |
| Native v3 create and Classic-fallback create with descriptions | Passed |
| Data source lookup by numeric ID and name | Passed |
| Clean plans after create | Passed |
| Rename, description update, and complete criteria updates on both write paths | Passed |
| Clean plans after updates | Passed |
| Explicit description clearing on both write paths | Passed |
| Clearing all native criteria, including data-source read-back | Passed |
| Out-of-band deletion plans recreation, then apply recreates the group | Passed |
| Import an existing numeric group ID and retain configured timeouts | Passed |
| Noncontiguous priority is rejected before an API write | Passed |
| Terraform destroy and independent API 404 verification | Passed |

An initial run exposed an untyped null list in data sources when criteria were empty. That run's two groups were deleted and verified absent. The corrected full run passed all checks and deleted all three group IDs it created. See `results.json` for the successful run's exit codes and cleanup IDs.

No patch software title was configured in the lab, so successful Patch Reporting fallback was not live-tested. The real extension-attribute fallback uses the same refusal classifier and Classic write sequence; unit tests cover both refusal messages.
