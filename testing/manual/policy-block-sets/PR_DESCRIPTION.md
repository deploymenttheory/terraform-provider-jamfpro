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

### Live Terraform lifecycle verification

Verified September 19, 2026 (JST), against `https://onishidev.jamfcloud.com`, using
Terraform 1.14.6 and separate baseline (`f8efc66a`) and PR provider binaries.
The policy remained **disabled with no scoped computers**. All test objects were
removed afterward.

HCL below expands the recorded runner's inputs into literal values; it is not a
separate execution. Only the relevant policy resource/blocks are shown, with
actual dependency IDs. Provider setup, credentials, refresh messages, and
expected development-override warnings are omitted. The original runner is in
`testing/manual/policy-block-sets/`.

#### 1. Apply the baseline policy with the old provider

Three scripts and three categories had been created by Terraform, followed by
policy ID `12`. The following HCL was then applied with schema version 1.
Archive paths explicitly match Jamf's defaults, excluding unrelated drift.

HCL applied:

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

Plan (summary; the script-order diff is shown in Motivation and Context):

```text
Plan: 0 to add, 1 to change, 0 to destroy.
```

Apply (`baseline-normalized.tfplan`):

```text
jamfpro_policy.test: Modifying... [id=12]
jamfpro_policy.test: Modifications complete after 0s [id=12]
Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

Replanning **without editing this HCL** still produced the script-order diff
shown above, with exit code **2**. This established the v1 state used below.

#### 2. Migrate the same state using the PR provider

HCL: **unchanged from step 1**. Only the provider binary changed; no import or
state removal was used.

Plan (`migration.tfplan`, exit **0**):

```text
No changes. Your infrastructure matches the configuration.
```

Apply:

```text
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
```

State assertions:

```text
Policy ID: 12
Schema version: 1 -> 2
All policy attribute values preserved (ignoring collection order).
Scripts: 3
Local accounts: 3
Self Service categories: 3
```

#### 3. Add printers and Dock items

Three printers and three Dock items were created and attached to the same policy.
HCL added inside `payloads` (all step 1 blocks remain):

```hcl
printers {
  id           = 2
  name         = "tf-policy-sets-20260919-charlie"
  action       = "uninstall"
  make_default = false
}
printers {
  id           = 1
  name         = "tf-policy-sets-20260919-alpha"
  action       = "uninstall"
  make_default = false
}
printers {
  id           = 3
  name         = "tf-policy-sets-20260919-bravo"
  action       = "uninstall"
  make_default = false
}

dock_items {
  id     = 2
  name   = "tf-policy-sets-20260919-charlie"
  action = "Remove"
}
dock_items {
  id     = 3
  name   = "tf-policy-sets-20260919-alpha"
  action = "Remove"
}
dock_items {
  id     = 1
  name   = "tf-policy-sets-20260919-bravo"
  action = "Remove"
}
```

Plan (`expanded.tfplan`):

```text
Plan: 6 to add, 1 to change, 0 to destroy.
```

Apply (dependency creation messages omitted):

```text
jamfpro_policy.test: Modifying... [id=12]
jamfpro_policy.test: Modifications complete after 1s [id=12]
Apply complete! Resources: 6 added, 1 changed, 0 destroyed.
```

#### 4. Reorder the blocks without changing their values

All five populated collections were reordered from charlie → alpha → bravo to
bravo → alpha → charlie. For example, the `scripts` blocks inside `payloads`
became the following; accounts, categories, printers, and Dock items were reordered
in the same way. No field values changed.

HCL planned:

```hcl
scripts {
  id         = "3"
  priority   = "After"
  parameter4 = "bravo-original"
}
scripts {
  id         = "2"
  priority   = "After"
  parameter4 = "alpha-original"
}
scripts {
  id         = "1"
  priority   = "After"
  parameter4 = "charlie-original"
}
```

Plan (`reordered.tfplan`, exit **0**):

```text
No changes. Your infrastructure matches the configuration.
```

Apply: not run for this ordering-only check.

#### 5. Update script parameters in place

Keep the step 4 order and all other settings. Change only the three `parameter4`
values:

HCL applied inside `payloads`:

```hcl
scripts {
  id         = "3"
  priority   = "After"
  parameter4 = "bravo-updated"
}
scripts {
  id         = "2"
  priority   = "After"
  parameter4 = "alpha-updated"
}
scripts {
  id         = "1"
  priority   = "After"
  parameter4 = "charlie-updated"
}
```

Plan (`content-change-fixed.tfplan`):

```text
  # jamfpro_policy.test will be updated in-place
  ~ resource "jamfpro_policy" "test" {
        id                            = "12"
        name                          = "tf-policy-sets-20260919"
        # (17 unchanged attributes hidden)

      ~ payloads {
            # (1 unchanged attribute hidden)

          - scripts {
              - id          = "1" -> null
              - parameter4  = "charlie-original" -> null
              - priority    = "After" -> null
                # (7 unchanged attributes hidden)
            }
          - scripts {
              - id          = "2" -> null
              - parameter4  = "alpha-original" -> null
              - priority    = "After" -> null
                # (7 unchanged attributes hidden)
            }
          - scripts {
              - id          = "3" -> null
              - parameter4  = "bravo-original" -> null
              - priority    = "After" -> null
                # (7 unchanged attributes hidden)
            }
          + scripts {
              + id          = "1"
              + parameter4  = "charlie-updated"
              + priority    = "After"
                # (7 unchanged attributes hidden)
            }
          + scripts {
              + id          = "2"
              + parameter4  = "alpha-updated"
              + priority    = "After"
                # (7 unchanged attributes hidden)
            }
          + scripts {
              + id          = "3"
              + parameter4  = "bravo-updated"
              + priority    = "After"
                # (7 unchanged attributes hidden)
            }

            # (7 unchanged blocks hidden)
        }

        # (3 unchanged blocks hidden)
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

Apply:

```text
jamfpro_policy.test: Modifying... [id=12]
jamfpro_policy.test: Modifications complete after 0s [id=12]
Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

Assertions confirmed policy ID `12`, exactly three scripts, and the three
`*-updated` parameters. Nested remove/add entries represent changed set elements;
the policy and script resources were not recreated.

#### 6. Verify the final no-op plan

HCL: **unchanged from step 5**.

Plan after refresh (exit **0**):

```text
No changes. Your infrastructure matches the configuration.
```

Apply: not run; the preceding apply had already persisted these values.

#### 7. Destroy the temporary resources

HCL: **unchanged from step 5**. Run `terraform plan -destroy`, then apply the
saved plan (`cleanup.tfplan`).

Plan:

```text
Plan: 0 to add, 0 to change, 13 to destroy.
```

Apply (dependency deletion messages omitted):

```text
jamfpro_policy.test: Destroying... [id=12]
jamfpro_policy.test: Destruction complete after 0s
Apply complete! Resources: 0 added, 0 changed, 13 destroyed.
```

All 13 resources reported successful deletion; the final Terraform state was empty.

**Coverage:** Live v1 → v2 migration covers scripts, accounts, and categories.
Printers and Dock items were added afterward and included in the reordering check.
Packages, directory bindings, and v0 migration are covered by unit tests.

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
