# Collection migration verification

Convert `ldap_pane.ldap_group_access` to sets in `enrollment_customization`. Add schema version 1 with the complete frozen v0 schema and a state upgrader; preserve all unrelated attributes. Singleton and ordered lists remain lists.

Draft: the tenant has no LDAP servers. Creating an LDAP pane with an empty access collection returned HTTP 400 (INVALID_FIELD: ldapGroupAccess). The unattached parent was deleted. Populated LDAP access migration, API readback, mutation and reapply remain unverified. Offline fixtures are synthetic and were never applied.

Baseline: `f8efc66ae7a7242ac41ff25727bc3e555dc26c85`. Tests used separate locally built old/new provider binaries through Terraform development overrides. Files pair exact resource HCL with sanitized CLI output; provider authentication is omitted. Logs include exploratory baseline runs where retained; the PR identifies the final verification sequence. Offline fixtures were not applied. Empty membership reorder is only a no-op, not a nonempty ordering test. No state or credentials are included.

Commands: `terraform plan -input=false -no-color -out=change.tfplan`, `terraform apply -input=false -no-color change.tfplan`, and unchanged-HCL `terraform plan -input=false -no-color -detailed-exitcode`. Cleanup uses `terraform plan -destroy` followed by apply. Unit commands and outputs are in `unit-tests.txt`.

Environment: 2026-09-19, Terraform 1.14.6, darwin_arm64, Go 1.26.4. All temporary resource IDs were verified absent through API HTTP 404; dedicated Terraform states contain zero resources.
