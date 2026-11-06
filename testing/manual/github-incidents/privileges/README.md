# Privilege catalog verification — 2026-11-06

All comparisons use actual curl responses from the authorized Jamf Pro 11.33.3 lab. No existing accounts, groups, API roles, or client permissions were changed.

- `/api/v1/api-role-privileges`: HTTP 200, 525 unique names.
- Search responses for PKI, User-Initiated Enrollment, License Information, Allowed File Extension, and blueprints exactly matched the corresponding values in the full catalog.
- Excluding the two exposed backend SaaS permissions (`Read Knobs`, `Update Knobs`), the catalog contains 523 user-facing names. These exactly match the reference Administrator account and three full-privilege account groups.
- The SDK lookup exactly matched the fresh curl capture. The validator accepted current user-facing names and rejected a fabricated name.

`catalog.json` preserves the catalog response. `comparison.json` contains the compared value sets and differences, with account/client identities omitted. The two Knobs permissions are not used as evidence of an end-user privilege defect.

The validator retains strict name checking and requires **Read API Roles** to read this catalog. Name validity is distinct from authority to grant the privilege; Jamf checks the latter on the write.

To repeat the direct request with a privately supplied token:

```sh
curl --silent --show-error \
  --header "Authorization: Bearer ${JAMFPRO_TOKEN}" \
  --header "Accept: application/json" \
  "${JAMFPRO_INSTANCE_FQDN}/api/v1/api-role-privileges"
```

The optional integration test compares the real SDK lookup with a fresh curl capture; supply the usual OAuth environment variables and `JAMFPRO_CATALOG_CAPTURE`, then run:

```sh
go test -tags=integration ./testing/manual/github-incidents -run TestLivePrivilegeValidationMatchesCurl -count=1
```
