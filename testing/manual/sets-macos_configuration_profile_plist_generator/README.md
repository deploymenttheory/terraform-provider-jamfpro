# Collection migration verification

Convert `payloads.payload_content.setting and its nested dictionary collections` to sets in `macos_configuration_profile_plist_generator`. Add schema version 1 with the complete frozen v0 schema and a state upgrader; preserve all unrelated attributes. Singleton and ordered lists remain lists. Preserve plist array order and duplicates through explicit `array_json`; accept both legacy list and set dictionary input in the shared converter. Only this generator calls the changed conversion entry points.

Draft: top-level settings migrated live. Legacy nested-dictionary migration is covered by JSON/flatmap tests but was not separately seeded with the old provider live. New nested dictionaries and ordered JSON arrays were applied and read back live. Array-only updates can also render unchanged sibling setting blocks as remove/add in the SDK plan; subsequent plans are clean. array_json accepts JSON-compatible arrays; plist date/binary array values are rejected rather than silently converted.

Baseline: `f8efc66ae7a7242ac41ff25727bc3e555dc26c85`. Tests used separate locally built old/new provider binaries through Terraform development overrides. Files pair exact resource HCL with sanitized CLI output; provider authentication is omitted. Logs include exploratory baseline runs where retained; the PR identifies the final verification sequence. Offline fixtures were not applied. Empty membership reorder is only a no-op, not a nonempty ordering test. No state or credentials are included.

Commands: `terraform plan -input=false -no-color -out=change.tfplan`, `terraform apply -input=false -no-color change.tfplan`, and unchanged-HCL `terraform plan -input=false -no-color -detailed-exitcode`. Cleanup uses `terraform plan -destroy` followed by apply. Unit commands and outputs are in `unit-tests.txt`.

Environment: 2026-09-19, Terraform 1.14.6, darwin_arm64, Go 1.26.4. All temporary resource IDs were verified absent through API HTTP 404; dedicated Terraform states contain zero resources.
