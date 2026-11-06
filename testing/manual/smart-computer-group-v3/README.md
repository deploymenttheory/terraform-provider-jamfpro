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

## Pre-merge lifecycle verification

A fresh run on 2026-11-06 used local Terraform **1.16.4** and the current provider binary against Jamf Pro **11.33.3**. Both native v3 and Classic fallback completed the requested sequence:

| Phase | Changes and verification |
| --- | --- |
| Create | Create both groups with descriptions and two criteria. A separate refused v3 request confirmed the live extension-attribute `has` defect returns `400 INVALID_FIELD`. |
| Update | Rename both groups, change descriptions and criterion values, and add a third criterion. Saved plans contain only in-place updates. |
| Different update | Retain the new names, reorder criteria, change values and `and_or`, and add balanced parentheses. Omit the native description to clear it; keep the fallback description unchanged to verify preservation. Saved plans contain only in-place updates. |
| Delete | Apply a saved destroy plan for both groups. Confirm empty Terraform state, remove the group declarations, and confirm a clean plan. |

After every create/update, API responses, Terraform resource state, both ID/name data sources for each group, and configured timeouts matched the expected values. IDs **27002** and **27003** remained unchanged through both updates. Two ordinary `terraform plan -detailed-exitcode` commands with refresh enabled returned **0** after each phase and after removing the destroyed groups from configuration.

Additional refresh-only plans printed “No changes” but returned **2** on Terraform 1.16.4 while the state contained groups. Their saved JSON contained no resource drift, resource changes, or output changes, and the refreshed resource/data-source values equaled saved state. These auxiliary plans were not applied. An initial attempt stopped at this unexpected exit code; its groups **27000** and **27001** were also deleted. The ordinary refreshed plans returned **0** throughout the completed run.

The integration fixture initially failed CI because its shared `provider.tf` symlink was missing. Restoring `provider.tf -> ../../provider.tf` follows the convention in `testing/README.md` and makes the test module require `terraform.local/local/jamfpro`. The repaired fixture passed local `terraform init`, `terraform validate`, and `terraform test -parallelism=1` (**1 passed, 0 failed**), including teardown, against the lab.

Independent API checks returned **404** for all four lifecycle attempt IDs. A Classic API list contained no remaining groups matching either this lifecycle prefix or the local CI fixture prefix. Existing lab objects were only read.

See [the detailed results](lifecycle/results.json) and the sanitized configurations in `lifecycle/`. To repeat the sequence, copy `provider.tf.example` to a private workspace as `provider.tf`, then copy one phase file at a time to `groups.tf`. Supply a unique `test_prefix` beginning with `tf-testing-`, the name of an existing string extension attribute, the usual `JAMFPRO_*` OAuth environment variables, and a Terraform development override pointing at the local provider binary. Apply each phase, check repeated refreshed plans, then destroy and independently verify cleanup. The `.example` suffix prevents the alternative configurations from loading together.
