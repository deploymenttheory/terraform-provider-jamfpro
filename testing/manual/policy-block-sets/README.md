# Policy collection migration verification

`main.tf` creates a disabled, unscoped policy with three scripts, local-account
entries, and Self Service categories. Set `include_printers_and_dock_items=true`
to add three printers and three Dock items. The dependent objects are temporary;
the scripts exit without doing any work. Package and directory-binding blocks
are covered by the unit tests, not this live fixture.

Use an isolated working directory outside the repository for state, saved plans,
credentials, and complete logs. Supply credentials through a separate `-var-file`
or sensitive `TF_VAR_jamfpro_client_id` / `TF_VAR_jamfpro_client_secret`
environment variables. Never commit credentials or state.

The API client needs Create, Read, and Delete privileges for Scripts, Categories,
Printers, and Dock Items, plus Create, Read, Update, and Delete privileges for
Policies. The fixture uses a five-second token refresh buffer to support short
API client token lifetimes.

Build the baseline provider from commit `f8efc66a` and the modified provider into
separate directories. Configure two Terraform CLI configuration files with
`provider_installation.dev_overrides` for `deploymenttheory/jamfpro`, pointing at
the respective directories. With development overrides, skip `terraform init`.

Run the following phases, recording each command's exit code and output:

1. With the baseline provider and default variables, plan with
   `-out=baseline.tfplan`, review the plan, then apply it. Save a state backup and
   confirm policy `schema_version` is `1`.
2. Run another baseline `plan -detailed-exitcode` without changing HCL. The API
   returns scripts in alphabetical name order rather than declaration order,
   producing the repeated ordering diff. Apply the plan and replan to confirm
   that applying does not resolve the drift.
3. Switch to the modified provider using the same HCL and state. Plan and apply
   without import or state removal. Confirm `schema_version` becomes `2`, the
   policy ID is unchanged, and all attribute values are preserved.
4. Plan and apply with `-var=include_printers_and_dock_items=true` to extend the
   live check to printers and Dock items.
5. Run `plan -detailed-exitcode` with the previous flag and
   `-var='item_order=["bravo","alpha","charlie"]'`. The ordering-only check
   should exit `0` with no changes.
6. Add `-var=script_parameter=updated`, plan, and apply the saved plan. Confirm
   the policy updates in place. Run a follow-up plan with the same variables and
   confirm exit `0`.
7. Plan with `-destroy` and the same variables, then apply the saved destroy
   plan. Confirm the final state contains zero managed resources.

`archive_home_directory_to` is explicitly configured because Jamf otherwise
populates a default path, producing a separate diff unrelated to collection
ordering. `enabled=false` and an empty scope remain in effect for every phase.

See `PR_DESCRIPTION.md` for actual command outputs and the verification scope.

## Recorded result

The September 19, 2026 run reproduced the old unchanged-HCL script drift
(exit `2`), upgraded the same state without remote changes (exit `0`), preserved
all policy attributes and ID `12`, ignored reordering of five populated
collections, applied real script parameter changes in place, and ended with a
clean plan (exit `0`). All 13 temporary resources were deleted and the final
state was empty. See the PR description for exact outputs and coverage limits.
