# Pull Request Description

## Summary

- Convert seven repeated policy blocks from lists to sets: `payloads.scripts`,
  `payloads.printers`, `payloads.dock_items`, `payloads.packages.package`,
  `payloads.account_maintenance.local_accounts.account`,
  `payloads.account_maintenance.directory_bindings.binding`, and
  `self_service.self_service_category`.
- Add automatic state migration to schema version 2, retaining the v0 → v1 → v2
  path and existing policy identity. No import or state removal is required.
- Read sets in constructors, retain account passwords by username during refresh,
  and remove a redundant script diff suppressor that created empty set elements
  during parameter updates. Single-configuration containers remain lists.

### Issue Reference

N/A. Follows the migration approach in [PR #1077](https://github.com/deploymenttheory/terraform-provider-jamfpro/pull/1077).

### Motivation and Context

After applying a policy, an **unchanged HCL configuration still plans an update**
because Jamf returns scripts in a different order. HCL declared charlie → alpha
→ bravo, while the API returned alpha → bravo → charlie.

The following is actual output from the baseline provider (`f8efc66a`) against
`https://onishidev.jamfcloud.com`, immediately after another successful apply.
Only refresh messages and the development-override warning are omitted:

```console
$ terraform plan -detailed-exitcode -input=false -no-color -out=baseline-drift.tfplan -var-file=<credentials>
  # jamfpro_policy.test will be updated in-place
  ~ resource "jamfpro_policy" "test" {
        id                            = "12"
        name                          = "tf-policy-sets-20260919"
        # (17 unchanged attributes hidden)

      ~ payloads {
            # (1 unchanged attribute hidden)

          ~ scripts {
              ~ id          = "2" -> "1"
              ~ parameter4  = "alpha-original" -> "charlie-original"
                # (8 unchanged attributes hidden)
            }
          ~ scripts {
              ~ id          = "3" -> "2"
              ~ parameter4  = "bravo-original" -> "alpha-original"
                # (8 unchanged attributes hidden)
            }
          ~ scripts {
              ~ id          = "1" -> "3"
              ~ parameter4  = "charlie-original" -> "bravo-original"
                # (8 unchanged attributes hidden)
            }

            # (1 unchanged block hidden)
        }

        # (3 unchanged blocks hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

Exit code: **2**. The script IDs and parameter values are unchanged as a group;
only their list positions differ. Applying again does not resolve the drift.
Switching to this PR's provider with the **same HCL and state** produces
`No changes` (exit **0**). Sets eliminate this ordering diff while preserving
real additions, removals, and field changes.

### Dependencies

No dependency/version updates. Live testing requires API client Create/Read/Delete
permissions for Scripts, Categories, Printers, and Dock Items, and
Create/Read/Update/Delete permissions for Policies.

## Type of Change

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)
- [x] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [x] 📝 Documentation update (Wiki/README/Code comments)
- [ ] ♻️ Refactor (code improvement without functional changes)
- [ ] 🎨 Style update (formatting, renaming)
- [ ] 🔧 Configuration change
- [ ] 📦 Dependency update

Automatic state migration preserves resources, but expressions using numeric
indexes into these collections must change to identity-based selection.

## Testing

- [x] I have added unit tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] I have added integration tests following the [testing implementation guide](../docs/testing-implementation.md)
- [x] I have tested this code in the following browsers/environments: Terraform 1.14.6, darwin/arm64, live Jamf Pro tenant

### Automated checks

`go test ./internal/services/policy ./internal/services/macos_onboarding_settings`
and provider build pass. Tests cover v0/v1 JSON and legacy flatmap migration,
attribute preservation, empty/null collections, ordering, add/remove/edit for all
seven collections, password retention, and parameter updates without empty scripts.
HCL formatting, validation, generated policy docs, and `git diff --check` pass.

**Existing full-suite failure:** `go test ./...` fails in unchanged test
`macos_configuration_profile_plist/TestDiffSuppressEquivalentPayloads/Validation_disabled_should_not_suppress`
(`expected: false`, `actual: true`). The integration checkbox is unchecked because
the live fixture was run manually, outside the automated integration harness.

### Live HCL and results

Tested September 19, 2026 (JST), using separate baseline and PR provider development
overrides. The policy stayed **disabled with no scoped computers** throughout.

<details>
<summary>Baseline policy HCL — literal values, no variables or dynamic blocks</summary>

The policy below expands the recorded fixture's inputs into their actual values;
it is not a separate run. Script IDs 1/2/3 are charlie/alpha/bravo, and category
IDs 2/3/1 are charlie/alpha/bravo. Dependencies were created first; provider and
credential setup are omitted. Archive paths match Jamf's defaults to exclude
unrelated drift.

```hcl
resource "jamfpro_policy" "test" {
  name          = "tf-policy-sets-20260919"
  enabled       = false
  trigger_other = "tf-policy-sets-20260919-never-run"
  frequency     = "Ongoing"

  network_limitations {
    minimum_network_connection = "No Minimum"
    any_ip_address             = false
  }

  scope {
    all_computers = false
    all_jss_users = false
  }

  self_service {
    use_for_self_service      = false
    self_service_display_name = "Terraform policy set migration test"

    self_service_category {
      id         = 2 # charlie
      display_in = true
      feature_in = false
    }
    self_service_category {
      id         = 3 # alpha
      display_in = true
      feature_in = false
    }
    self_service_category {
      id         = 1 # bravo
      display_in = true
      feature_in = false
    }
  }

  payloads {
    scripts {
      id         = "1"
      priority   = "After"
      parameter4 = "charlie-original"
    }
    scripts {
      id         = "2"
      priority   = "After"
      parameter4 = "alpha-original"
    }
    scripts {
      id         = "3"
      priority   = "After"
      parameter4 = "bravo-original"
    }

    account_maintenance {
      local_accounts {
        account {
          action                    = "Delete"
          username                  = "tf_charlie"
          archive_home_directory_to = "/tf_charlie.dmg"
        }
        account {
          action                    = "Delete"
          username                  = "tf_alpha"
          archive_home_directory_to = "/tf_alpha.dmg"
        }
        account {
          action                    = "Delete"
          username                  = "tf_bravo"
          archive_home_directory_to = "/tf_bravo.dmg"
        }
      }
    }
  }
}
```

</details>

| Step | Actual result |
| --- | --- |
| Baseline apply | Policy ID `12` created; schema version `1`; three scripts, accounts, and categories. Reapplying still produces the drift shown above. |
| Same-state migration plan/apply | `No changes`; exit `0`; `0 added, 0 changed, 0 destroyed`. Schema version becomes `2`; ID and all attribute values are preserved. |
| Add printers and Dock items | Three of each created and attached: `6 added, 1 changed, 0 destroyed`. |
| Reorder all five populated collections | charlie/alpha/bravo → bravo/alpha/charlie, with identical field values: `No changes`; exit `0`. |
| Change script parameters | `alpha-original`, `bravo-original`, `charlie-original` → corresponding `*-updated` values: `0 added, 1 changed, 0 destroyed`; policy ID stays `12`, with exactly three scripts. |
| Final unchanged-config plan | `No changes`; exit `0`. |
| Cleanup | `0 added, 0 changed, 13 destroyed`; final state contains zero resources. |

Selected output from migration, content update, final plan, and cleanup
(refresh messages and development-override warnings omitted):

```console
$ terraform plan -detailed-exitcode -input=false -no-color -out=migration.tfplan -var-file=<credentials>
No changes. Your infrastructure matches the configuration.

$ terraform apply -input=false -no-color migration.tfplan
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

$ terraform apply -input=false -no-color content-change-fixed.tfplan
jamfpro_policy.test: Modifying... [id=12]
jamfpro_policy.test: Modifications complete after 0s [id=12]
Apply complete! Resources: 0 added, 1 changed, 0 destroyed.

$ terraform plan -detailed-exitcode -input=false -no-color -var-file=<credentials> -var=include_printers_and_dock_items=true -var='item_order=["bravo","alpha","charlie"]' -var=script_parameter=updated
No changes. Your infrastructure matches the configuration.

$ terraform apply -input=false -no-color cleanup.tfplan
Apply complete! Resources: 0 added, 0 changed, 13 destroyed.
```

**Coverage:** Live v1 → v2 migration covers scripts, accounts, and categories.
Printers and Dock items were added with the PR provider and included in the
reordering check. Packages, directory bindings, and v0 migration are covered by
unit tests. The original runner and reproduction steps are in
`testing/manual/policy-block-sets/`.

## Quality Checklist

- [x] I have reviewed my own code before requesting review
- [ ] I have verified there are no other open Pull Requests for the same update/change
- [ ] All CI/CD pipelines pass without errors or warnings
- [x] My code follows the established style guidelines of this project
- [x] My comments are used only when necessary, ideally where the codes purpose is not self explanatory (eg: necessary magic numbers)
- [x] I have added necessary documentation (if appropriate)
- [x] I have made corresponding changes to the README and other relevant documentation
- [ ] My changes generate no new warnings

Migration guide, policy reference, and fixture README are updated. CI completion
and absence of all warnings are not claimed; local Terraform runs include the
expected development-override warning.

## Screenshots/Recordings (if appropriate)

Not applicable; actual Terraform output is included above.

## Additional Notes

Identical set elements are deduplicated. A migrated v2 state cannot be used
directly with an older provider. The migration guide documents indexing changes
and downgrade limitations. All 13 temporary tenant resources were deleted.
