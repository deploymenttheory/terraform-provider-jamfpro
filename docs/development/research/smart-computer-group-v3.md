# Smart computer group v3 implementation research

Research date: 2026-11-06. Scope: the new `jamfpro_smart_computer_group_v3` Framework resource and matching data source, including the required Classic API workaround. Jamf's provider source was inspected at commit `21c3969a3ba29111ca58a59c6b5187acf2e5a9bb`. Live checks used only uniquely named disposable groups in the explicitly authorized lab; their results and cleanup are recorded below.

## Endpoint and version decision

The intended object is a **smart** computer group. Use `/api/v3/computer-groups/smart-groups`; the supplied static-group documentation describes a different object. Jamf officially introduced the complete v3 computer-group surface in **Jamf Pro 11.28.0**, replacing v2. The published release entry takes precedence over the upstream provider's older comment saying no minimum version was documented. [Jamf release entry](https://developer.jamf.com/jamf-pro/changelog/11280-additions), [smart-group create reference](https://developer.jamf.com/jamf-pro/reference/post_v3-computer-groups-smart-groups).

| Operation | Route | Success | Required Jamf Pro privilege |
| --- | --- | --- | --- |
| Search | `GET /api/v3/computer-groups/smart-groups` | 200 | Read Smart Computer Groups |
| Create | `POST /api/v3/computer-groups/smart-groups` | 201, `{id, href}` | Create Smart Computer Groups |
| Read | `GET /api/v3/computer-groups/smart-groups/{id}` | 200 | Read Smart Computer Groups |
| Update | `PUT /api/v3/computer-groups/smart-groups/{id}` | 202, group body | Update Smart Computer Groups |
| Delete | `DELETE /api/v3/computer-groups/smart-groups/{id}` | 204 | Delete Smart Computer Groups |

The operation-specific privileges are identified by the official API references and Jamf's SDK method registry usage. The Classic fallback needs corresponding computer-group create/update privileges; its description patch needs smart-group update privilege. The upstream Platform provider has additional integration permissions for its Platform identifier bridge; those are not a reason to introduce a separate Platform client into this direct Jamf Pro provider. [Create](https://developer.jamf.com/jamf-pro/reference/post_v3-computer-groups-smart-groups), [read](https://developer.jamf.com/jamf-pro/reference/get_v3-computer-groups-smart-groups-id), [search](https://developer.jamf.com/jamf-pro/reference/get_v3-computer-groups-smart-groups), [upstream permission declarations](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/permissions.go).

## Wire fields and Terraform state

The create/update/read group body has these fields. **GET and PUT do not return an identifier**: retain the POST, import, or prior-state ID. There is no user-managed membership field on a smart group. [Pinned SDK v3 definitions](https://github.com/deploymenttheory/go-api-sdk-jamfpro/blob/v1.53.1/sdk/jamfpro/jamfproapi_smart_computer_groups.go), [upstream resource wire notes](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/resource.go).

| JSON field | Type | Meaning/default |
| --- | --- | --- |
| `name` | string | Required, maximum 255 characters. |
| `description` | string | Unset is stored as `""`. |
| `siteId` | string or null | `"-1"` represents no site. |
| `criteria` | array | Ordered criteria; omitted or empty removes all criteria. |
| `criteria[].name` | string | Server criterion name. |
| `criteria[].priority` | integer | Exactly the zero-based list position. |
| `criteria[].andOr` | string | Documented enum `"and"` or `"or"`; first entry has no preceding criterion. |
| `criteria[].searchType` | string | Operator valid for that criterion. |
| `criteria[].value` | string | Operator comparison value; may legitimately be empty. |
| `criteria[].openingParen` | boolean | Defaults false. |
| `criteria[].closingParen` | boolean | Defaults false. |

POST and PUT replace the complete group representation. Omitting description clears it, omitting site removes its assignment, and omitting criteria clears the list. The SDK's `omitempty` tags on description and criteria therefore do not prevent clearing through v3. Always supply the planned site, defaulting to `"-1"`, and normalize nil/empty site responses to that sentinel. A PUT response may express no site as null while GET expresses `"-1"`. [Upstream input builders](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/input_builders.go), [site normalization](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/common/progroups/site.go).

Criteria must be a list, since order changes evaluation. Priorities must be contiguous `0..n-1` in list order; reject an explicitly conflicting priority instead of silently changing a configured value. Default `and_or` to `"and"` and both parenthesis flags to false. The upstream provider returns null for an empty criterion list and consequently rejects configured `[]`; that is its chosen state representation. A provider can also support configured `[]` by preserving the planned empty-versus-null representation while treating both as no criteria. The lab returned JSON `[]` for both omitted and explicitly empty criteria. [Priority and criterion builders](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/common/criteria/pro_groups.go), [upstream schema choices](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/resource.go).

Names are unique case-insensitively across smart, static, and Platform COMPUTER groups. The v3 path stores surrounding whitespace verbatim; Classic removes it. This difference matters when a fallback succeeds: a known planned name must not silently become a different string in state. Jamf's duplicate-name check can race under concurrent creates; upstream collision acceptance tests serialize duplicate creates for this reason. [Collision implementation](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/collision.go), [collision tests](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/collision_acceptance_test.go), [state builder whitespace handling](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/state_builders.go).

## Required Classic fallback

Jamf's v3 validation rejects combinations supported by its own UI and Classic API: patch-reporting criteria (PI-1286, PI-1183, PI-626), and `has` on string extension attributes (PI-1032). Rejected v3 writes apply nothing, so only the classified refusal can be retried through Classic. The fallback classifier is deliberately broader than hardcoding a particular criterion name, and narrower than any HTTP 400: [authoritative classifier and reproductions](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/common/progroups/fallback.go).

```text
HTTP status == 400
AND an errors[] entry has code equal-fold "INVALID_FIELD"
AND its description contains "is not valid"
AND its description contains "The criterion " OR "The operator "
```

The phrase checks are case-sensitive in upstream. Do not trigger fallback for malformed priorities, invalid `andOr`, duplicate names, site refusal, authorization failures, transport errors, or arbitrary server errors. Because genuine misspelled criterion names also match this classifier, Classic must remain authoritative: its 409 containing `Problem with criteria` means both interfaces refused the criteria. Other fallback failures require their own diagnosis. Classic does not validate every criterion/operator pairing; report a warning after successful fallback. [Classifier and diagnostics](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/common/progroups/fallback.go), [classifier regression tests](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/fallback_test.go).

Classic create uses `POST /JSSResource/computergroups/id/0`; update uses `PUT /JSSResource/computergroups/id/{id}`. Send **name, site, and all criteria in one request**, with `is_smart=true` and no membership payload. Never create a temporary half-built criterion list: membership recalculation is synchronous and scoped actions could observe it. Classic has no description field and preserves the stored description. [Pinned SDK Classic implementation](https://github.com/deploymenttheory/go-api-sdk-jamfpro/blob/v1.53.1/sdk/jamfpro/classicapi_computer_groups.go), [fallback request builder](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/input_builders.go).

After fallback create with a nonempty description, or fallback update changing description, patch only `{ "groupDescription": "..." }` through `PATCH /api/v2/groups/{platformUUID}`. Explicit `""` clears description. This endpoint takes the Platform UUID; a Jamf Pro numeric group ID returns 404. Never include criteria, assignments, or name in that patch. Resolve the UUID by listing `/api/v2/groups` with `groupType=="COMPUTER" and isSmart=="true"`, then matching `groupJamfProId` locally. Pagination is required. Filtering directly by `groupJamfProId` is rejected; advertised `groupId` and `siteId` filters can return 500. [Official UUID patch reference](https://developer.jamf.com/jamf-pro/reference/patch_v2-groups-id), [identifier bridge and verified filter pitfalls](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/common/progroups/platformid.go).

## Consistency and failure handling

Jamf 11.32 can return 202 before reads reflect a write; upstream observed a 150–300 ms delay. Poll the read after create/update against the complete intended name, description, normalized site, and every criterion field. Upstream uses 250 ms intervals and a 10 s bound; a fresh first read needs no delay. Preserve the created ID in state before reporting a subsequent read/description-patch failure, so an already-created group is not orphaned. [Read-back helper and evidence](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/common/progroups/readback.go), [create recovery](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/crud.go).

**Wait for Classic name/site/criteria to be readable before sending the description patch.** On 11.32 that patch can read-modify-write a stale representation and permanently undo the preceding Classic change. Upstream reproduced the revert twice without a gap and zero times in four runs with at least 300 ms separation; use a settled read instead of a fixed sleep. [Description-patch ordering evidence](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/crud.go#L363).

Resource Read should remove state on an actual 404. A data source should report a missing group. Delete 404 is successful idempotent deletion. After a successful 204, reads may still see the group for 150–250 ms, so optionally wait for 404 before completing. Delete 422 `HAS_DEPENDENCIES` is permanent while the listed references remain: surface their names and do not repeatedly issue DELETE. [Read/delete implementation](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/crud.go), [dependency diagnostics](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/common/progroups/errors.go).

Import passes the numeric Jamf Pro group ID through to state and hydrates the ordinary GET fields; no create, name-based adoption, or second identifier is needed. Search is paginated (default page 0, page size 100) and supports `id`, `name`, and `siteId` filters. Only full-access admins may filter by site; sited admins are automatically constrained. A name filter is not sufficient to establish an exact unique match: compare returned names and detect ambiguity. [Import implementation](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/resource.go), [official search reference](https://developer.jamf.com/jamf-pro/reference/get_v3-computer-groups-smart-groups), [upstream name-filter caveat](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/common/progroups/platformid.go).

## Baseline SDK assessment (v1.53.1)

ADR-0001 requires the existing Jamf SDK. Version 1.53.1 already supplies v3 CRUD/list, Classic computer-group CRUD, and unified v2 group list/read/update. The existing SDK remains the integration layer; missing capabilities are added in [SDK PR #1172](https://github.com/deploymenttheory/go-api-sdk-jamfpro/pull/1172). Its v3 request fields match the researched group body. [Repository ADR](../adr/0001-go-api-sdk-jamfpro-usage.md), [v3 SDK](https://github.com/deploymenttheory/go-api-sdk-jamfpro/blob/v1.53.1/sdk/jamfpro/jamfproapi_smart_computer_groups.go), [unified-group SDK](https://github.com/deploymenttheory/go-api-sdk-jamfpro/blob/v1.53.1/sdk/jamfpro/jamfproapi_groups.go).

Important gaps:

- `GetSmartComputerGroupsV3` and `GetGroupsV2` fetch one page, unlike older SDK helpers that aggregate pages.
- `RequestGroupUpdateV2.GroupDescription` is a string with `omitempty`, so the typed patch method cannot express clearing description. The SDK must provide a description-only update method carrying `groupDescription` without omission. No transport adapter or raw HTTP request belongs in the provider.
- SDK endpoint errors format the underlying error with **`%v`, not `%w`**, losing the error chain. The new SDK context-aware methods preserve the transport APIError with `%w`, allowing the provider to use typed errors without parsing serialized error strings. Its real fields are `status_code`, `httpStatus`, and `errors[]` containing `code`, `field`, and `description`.
- Endpoint methods accept no context, and the pinned HTTP client starts requests with `context.Background()`. Framework timeout values therefore cannot cancel an already-running request through those methods. Do not promise physical cancellation or abandon a mutating request in a detached goroutine. The SDK transport's own request/retry timeout still matters.
- Delete and unified-group patch insist on 204. Preserve success handling for the documented 202 smart-group update, and test an empty update response defensively if supporting it.

[SDK error-format constants](https://github.com/deploymenttheory/go-api-sdk-jamfpro/blob/v1.53.1/sdk/jamfpro/shared_error_messages.go), [transport APIError](https://github.com/deploymenttheory/go-api-http-client/blob/v0.4.7/response/error.go), [transport request/context handling](https://github.com/deploymenttheory/go-api-http-client/blob/v0.4.7/httpclient/request.go).

The upstream provider uses five-minute create/update budgets and one-minute read/delete budgets because its **Platform gateway** allows writes up to 180 seconds. That gateway-specific fact does not establish an identical limit for this provider's direct Jamf Pro connection. [Upstream timeout rationale](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/resource.go).

## Other compatibility caveats

Upstream blocks 11.29.0 through 11.30.0 because PI-1394 can change another-Jamf-group criterion values from authored names into numeric IDs. Its own notes explicitly say the v3 smart-computer endpoint was not directly verified in that window; treat this as a documented precaution rather than a proven v3 defect. [Version gate and its limitations](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/common/progroups/version.go).

Directory-service group criteria on 11.29+ use base64-encoded JSON `{"uuid":"opaque-directory-id","serverId":"numeric-string"}` as the value. Passing an already encoded value is sufficient for wire-level support; accepting a human group name additionally requires directory lookup, encoding, and preserving the authored form on refresh. Exact names/casing matter, including `User last logged in - Mdm directory service group`. This conversion is a separate feature from the required Classic fallback. [Directory criterion wire contract](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/common/criteria/dsgroup.go).

## Authorized live curl evidence

Lab version: `11.33.3-t1793127642444`. Requests used OAuth client credentials supplied privately by the user, private curl configuration files, and no credentials in command arguments or logs. Existing extension attributes and patch-title configuration were read only. No existing group, device, setting, extension attribute, or patch title was changed. The following table records direct observations, not guarantees for every Jamf version.

| Probe | Observed result |
| --- | --- |
| v3 create | 201; response exactly had `id` and `href`; GET had name, description, criteria, siteId, no ID. |
| PUT only name | 202; description became `""`, criteria `[]`, siteId `"-1"`; GET agreed immediately. |
| Explicit empty description and criteria | 202; GET retained `""` and `[]`. |
| Priority 5 for first entry | 400 `INVALID_REQUEST_PARAMETER_VALUE`, field `priority`. |
| Uppercase `AND` / `OR` | Accepted (201 create / 202 update), then GET canonicalized to lowercase `and` / `or`, including the first criterion. Keep the documented lowercase schema vocabulary. |
| Invalid `xor` | 400 `INVALID_FIELD` for criteria/andOr. Its description does not satisfy the fallback classifier. |
| Rejected v3 update changing name, description, and criterion | 400; subsequent GET confirmed all three previous values unchanged. |
| Smart name reused by uppercase static create | 422 `DUPLICATE_FIELD`, field `name`; no static group created. |
| `has` on an existing STRING extension attribute | v3 400 `INVALID_FIELD` with the classifier's operator phrase. |
| Same full criterion via Classic | Create 201; v3 GET retained the criterion. |
| UUID bridge | Filtered `/api/v2/groups` returned the created group and its Platform UUID. |
| Description-only patch by UUID | 204; GET showed the description and retained the criterion. |
| Same patch by numeric group ID | 404 `NOT_FOUND`. |
| One Classic update carrying renamed group and two criteria | 201; GET retained both criteria and prior description. |
| UUID description patch with explicit `""` | 204; description cleared, name and both criteria retained. |
| Synthetic patch-reporting criterion without a configured patch title | v3 400 `INVALID_FIELD`; not evidence of successful Classic fallback for a real configured title. |
| Read consistency on this 11.33.3 lab | Each successful tested write was visible in the first GET. |
| Cleanup | Created groups `26989`, `26990`, and `26991`: DELETE 204, subsequent GET 404 for every group. |

The patch-title catalog was empty, so no configured patch-reporting fallback was exercised. The successful text-extension-attribute case covers the same classified v3 refusal and complete Classic/description-patch flow without creating dependencies.

## Regression checks to retain

1. Exercise real SDK HTTP serialization for v3 CRUD, Classic atomic fallback, UUID resolution, and description-only PATCH, including clearing description.
2. Check every fallback classifier boundary: known criterion/operator refusal, unrelated INVALID_FIELD, priority failure, authorization, duplicate, transport, and server failures; codes in arbitrary free text must not activate a workaround.
3. Verify a rejected modern update followed by one Classic update sends the full name/site/criteria together and never sends computer membership.
4. Reproduce a stale read after update and a stale Classic read before description patch; assert the description patch waits for the intended Classic fields.
5. Confirm normal GET 404 removes resource state, data-source absence errors, delete 404 succeeds, dependency errors retain state, and post-delete stale reads settle.
6. Cover criteria reorder/removal, defaults, valid empty comparison values, first priority zero, conflicting explicit priority, parentheses, no-site normalization, and empty-versus-null criteria consistency.
7. Cover multi-page data-source/UUID lookups, exact name matching, unknown selectors, import, missing POST ID, and recording an existing created ID when a later operation fails.

These cases follow the actual upstream failure modes and the lab observations above. Upstream provides focused examples in its [fallback tests](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/fallback_test.go), [read-back tests](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/common/progroups/readback_test.go), and [acceptance lifecycle tests](https://github.com/jamf/terraform-provider-jamfplatform/blob/21c3969a3ba29111ca58a59c6b5187acf2e5a9bb/internal/resources/pro/smart_computer_group/resource_acceptance_test.go).

## Implementation and live Terraform verification

The provider delegates all HTTP requests, request encoding, pagination, and Platform ID resolution to the Jamf Pro SDK. HTTP client v0.4.8 supplies caller-context propagation and bodyless success handling. The SDK additions were merged in [SDK PR #1172](https://github.com/deploymenttheory/go-api-sdk-jamfpro/pull/1172) and released in [v1.54.0](https://github.com/deploymenttheory/go-api-sdk-jamfpro/releases/tag/v1.54.0). The provider consumes that release without local dependency replacements. The complete live Terraform lifecycle results, including cleanup, are recorded in [the verification report](../../../testing/manual/smart-computer-group-v3/README.md).
