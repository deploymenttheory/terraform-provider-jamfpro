# Pull Request Description

## Summary

Policy scripts returned in API order caused repeated in-place updates even after
applying unchanged HCL. Convert the following repeated policy blocks from lists
to sets, with automatic state migration following
[PR #1077](https://github.com/deploymenttheory/terraform-provider-jamfpro/pull/1077):

- `payloads.scripts`
- `payloads.printers`
- `payloads.dock_items`
- `payloads.packages.package`
- `payloads.account_maintenance.local_accounts.account`
- `payloads.account_maintenance.directory_bindings.binding`
- `self_service.self_service_category`

Add schema version 2 and preserve the v0 → v1 → v2 migration chain. Complete
historical policy schema types retain unrelated fields when decoding legacy
flatmap state. Constructors read sets; account refresh retains passwords by
username instead of list position.

Remove the redundant `parameter5` equality diff suppressor: with set elements,
it retained empty values under removed element hashes and produced a spurious
empty script during parameter updates. A regression test reproduces this issue
before the fix and verifies the planned collection after the fix.

Single-configuration containers remain lists. Existing block syntax is unchanged,
but expressions indexing these collections numerically must select by identity.
Identical elements are deduplicated. No import, state removal, or policy
recreation is required.


### Issue Reference

N/A. This follows the state migration approach in [PR #1077](https://github.com/deploymenttheory/terraform-provider-jamfpro/pull/1077).

### Motivation and Context

The Jamf Pro API returns policy scripts in a different order from HCL declaration
order. With list-backed blocks, an unchanged configuration continues to plan
in-place updates after apply. Sets remove these order-only differences while
still detecting additions, removals, and actual field changes. The Testing
section below includes the reproduced drift and the same-state migration result.

### Dependencies

- No Go module or provider dependency updates.
- Live verification uses Terraform 1.14.6 and API client privileges for Create,
  Read, and Delete on Scripts, Categories, Printers, and Dock Items, plus Create,
  Read, Update, and Delete on Policies.
- Existing block syntax remains valid; numeric-index expressions for the changed
  collections must be updated to select elements by identity.

## Type of Change

Please mark the relevant option with an `x`:

- [ ] 🐛 Bug fix (non-breaking change which fixes an issue)
- [ ] ✨ New feature (non-breaking change which adds functionality)
- [x] 💥 Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [x] 📝 Documentation update (Wiki/README/Code comments)
- [ ] ♻️ Refactor (code improvement without functional changes)
- [ ] 🎨 Style update (formatting, renaming)
- [ ] 🔧 Configuration change
- [ ] 📦 Dependency update

Marked as a breaking change because numeric indexing into the changed collections
is no longer supported, despite automatic state migration preserving resources.

## Testing

- [x] I have added unit tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] I have added integration tests following the [testing implementation guide](../docs/testing-implementation.md)
- [x] I have tested this code in the following browsers/environments: Terraform 1.14.6 on darwin/arm64 against a live Jamf Pro tenant

The full-unit-suite checkbox remains unchecked because of the existing failure
listed below. The live fixture is a manual Terraform verification, not an
integration test added to the repository's automated integration harness.

### Automated verification

- `go test ./internal/services/policy ./internal/services/macos_onboarding_settings` passes.
- Protocol tests cover JSON and legacy flatmap migration from v0 and v1,
  preserving all policy attributes, including empty/null/omitted collections.
- Tests reproduce list ordering drift, confirm order-independent sets, and cover
  addition, removal, and field changes in all seven collections without replacement.
- Tests cover request construction, password retention after API reordering, and
  parameter updates without empty script elements.
- Provider build, Terraform 1.14.6 validation, HCL formatting, and `git diff --check` pass.
- The policy reference document matches `tfplugindocs` output from the new schema.
- `go test ./...` fails in the unchanged
  `macos_configuration_profile_plist/TestDiffSuppressEquivalentPayloads/Validation_disabled_should_not_suppress`
  test (`expected: false`, `actual: true`).


### Real tenant verification

Verified against `https://onishidev.jamfcloud.com` on September 19, 2026 (JST),
using Terraform 1.14.6 on darwin/arm64, a baseline provider built from `f8efc66a`,
and the modified provider built locally from this branch. Separate development
overrides select the two binaries. Credentials and full state are excluded.

The policy remained disabled with no scoped computers throughout. The fixture
uses three elements per collection. HCL declares them in `charlie → alpha → bravo`
order; the API returned scripts in `alpha → bravo → charlie` order.

#### HCL

The complete reproducible fixture is committed as
`testing/manual/policy-block-sets/main.tf`:

```hcl
terraform {
  required_providers {
    jamfpro = {
      source = "deploymenttheory/jamfpro"
    }
  }
}

variable "jamfpro_client_id" {
  type      = string
  sensitive = true
}

variable "jamfpro_client_secret" {
  type      = string
  sensitive = true
}

variable "suffix" {
  type    = string
  default = "policy-sets-20260919"
}

variable "item_order" {
  type    = list(string)
  default = ["charlie", "alpha", "bravo"]
}

variable "script_parameter" {
  type    = string
  default = "original"
}

variable "include_printers_and_dock_items" {
  type    = bool
  default = false
}

provider "jamfpro" {
  jamfpro_instance_fqdn                = "https://onishidev.jamfcloud.com"
  auth_method                          = "oauth2"
  client_id                            = var.jamfpro_client_id
  client_secret                        = var.jamfpro_client_secret
  token_refresh_buffer_period_seconds  = 5
  hide_sensitive_data                  = true
  jamfpro_load_balancer_lock           = true
  mandatory_request_delay_milliseconds = 100
}

resource "jamfpro_script" "test" {
  for_each        = toset(var.item_order)
  name            = "tf-${var.suffix}-${each.key}"
  priority        = "AFTER"
  script_contents = "#!/bin/sh\nexit 0\n"
}

resource "jamfpro_category" "test" {
  for_each = toset(var.item_order)
  name     = "tf-${var.suffix}-${each.key}"
  priority = 9
}

resource "jamfpro_dock_item" "test" {
  for_each = var.include_printers_and_dock_items ? toset(var.item_order) : toset([])
  name     = "tf-${var.suffix}-${each.key}"
  type     = "App"
  path     = "file://localhost/Applications/${each.key}.app"
}

resource "jamfpro_printer" "test" {
  for_each    = var.include_printers_and_dock_items ? toset(var.item_order) : toset([])
  name        = "tf-${var.suffix}-${each.key}"
  cups_name   = "tf_${each.key}"
  uri         = "ipp://192.0.2.1/${each.key}"
  use_generic = true
}

resource "jamfpro_policy" "test" {
  name          = "tf-${var.suffix}"
  enabled       = false
  trigger_other = "tf-${var.suffix}-never-run"
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
    dynamic "self_service_category" {
      for_each = var.item_order
      content {
        id         = jamfpro_category.test[self_service_category.value].id
        display_in = true
        feature_in = false
      }
    }
  }

  payloads {
    dynamic "scripts" {
      for_each = var.item_order
      content {
        id         = jamfpro_script.test[scripts.value].id
        priority   = "After"
        parameter4 = "${scripts.value}-${var.script_parameter}"
      }
    }
    dynamic "printers" {
      for_each = var.include_printers_and_dock_items ? var.item_order : []
      content {
        id           = jamfpro_printer.test[printers.value].id
        name         = jamfpro_printer.test[printers.value].name
        action       = "uninstall"
        make_default = false
      }
    }
    dynamic "dock_items" {
      for_each = var.include_printers_and_dock_items ? var.item_order : []
      content {
        id     = jamfpro_dock_item.test[dock_items.value].id
        name   = jamfpro_dock_item.test[dock_items.value].name
        action = "Remove"
      }
    }
    account_maintenance {
      local_accounts {
        dynamic "account" {
          for_each = var.item_order
          content {
            action                    = "Delete"
            username                  = "tf_${account.value}"
            archive_home_directory_to = "/tf_${account.value}.dmg"
          }
        }
      }
    }
  }
}
```

Default variables are used for baseline creation and migration. The archive path
is explicit because Jamf fills in that default even when omitted, which would
otherwise introduce an unrelated diff. The five-second token refresh buffer
supports the API client's short token lifetime.

#### 1. Create and apply baseline state

Three scripts and three categories were created by Terraform first. The baseline
then planned `1 to add, 0 to change, 0 to destroy` for the policy:

```console
$ TF_CLI_CONFIG_FILE=baseline.tfrc terraform apply -input=false -no-color baseline-policy.tfplan
jamfpro_policy.test: Creating...
jamfpro_policy.test: Creation complete after 1s [id=12]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

The policy ID was `12`, with `schema_version: 1`. After explicitly configuring
Jamf's local-account archive defaults, the old provider was planned and applied
again with the final baseline HCL:

```console
$ TF_CLI_CONFIG_FILE=baseline.tfrc terraform apply -input=false -no-color baseline-normalized.tfplan
jamfpro_policy.test: Modifying... [id=12]
jamfpro_policy.test: Modifications complete after 0s [id=12]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

#### 2. Reproduce the original drift after apply, with unchanged HCL

```console
$ TF_CLI_CONFIG_FILE=baseline.tfrc terraform plan -detailed-exitcode -input=false -no-color -out=baseline-drift.tfplan -var-file=<credentials>
Terraform will perform the following actions:

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

Exit code: **2**. This is an unchanged-HCL replan after apply, not a deliberate
configuration reorder. The same three scripts retain their IDs and parameter
values; only list position differs. Applying the old plan does not eliminate it.

#### 3. Upgrade the same state with the modified provider

No HCL edits, import, or state removal occurred between the preceding drift plan
and this migration plan:

```console
$ TF_CLI_CONFIG_FILE=current.tfrc terraform plan -detailed-exitcode -input=false -no-color -out=migration.tfplan -var-file=<credentials>
jamfpro_category.test["bravo"]: Refreshing state... [id=1]
jamfpro_script.test["charlie"]: Refreshing state... [id=1]
jamfpro_script.test["bravo"]: Refreshing state... [id=3]
jamfpro_script.test["alpha"]: Refreshing state... [id=2]
jamfpro_category.test["charlie"]: Refreshing state... [id=2]
jamfpro_category.test["alpha"]: Refreshing state... [id=3]
jamfpro_policy.test: Refreshing state... [id=12]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration
and found no differences, so no changes are needed.
```

Exit code: **0**.

```console
$ TF_CLI_CONFIG_FILE=current.tfrc terraform apply -input=false -no-color migration.tfplan
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
```

State assertions after apply:

```text
Policy ID: 12
Schema version: 1 -> 2
All policy attribute values preserved (ignoring collection order).
Scripts: 3
Local accounts: 3
Self Service categories: 3
```

#### 4. Extend live coverage to printers and Dock items

After the API role permissions were available, plan with
`-var=include_printers_and_dock_items=true` produced
`6 to add, 1 to change, 0 to destroy`. Apply:

```console
$ TF_CLI_CONFIG_FILE=current.tfrc terraform apply -input=false -no-color expanded.tfplan
jamfpro_dock_item.test["charlie"]: Creating...
jamfpro_dock_item.test["bravo"]: Creating...
jamfpro_printer.test["bravo"]: Creating...
jamfpro_printer.test["alpha"]: Creating...
jamfpro_dock_item.test["alpha"]: Creating...
jamfpro_printer.test["charlie"]: Creating...
jamfpro_dock_item.test["bravo"]: Creation complete after 0s [id=1]
jamfpro_printer.test["charlie"]: Creation complete after 0s [id=2]
jamfpro_printer.test["alpha"]: Creation complete after 0s [id=1]
jamfpro_dock_item.test["charlie"]: Creation complete after 1s [id=2]
jamfpro_printer.test["bravo"]: Creation complete after 1s [id=3]
jamfpro_dock_item.test["alpha"]: Creation complete after 1s [id=3]
jamfpro_policy.test: Modifying... [id=12]
jamfpro_policy.test: Modifications complete after 1s [id=12]

Apply complete! Resources: 6 added, 1 changed, 0 destroyed.
```

#### 5. Reorder all five populated collections

Keep `include_printers_and_dock_items=true`, and change only the order variable:

```console
$ TF_CLI_CONFIG_FILE=current.tfrc terraform plan -detailed-exitcode -input=false -no-color -var-file=<credentials> -var=include_printers_and_dock_items=true -var='item_order=["bravo","alpha","charlie"]'
jamfpro_category.test["charlie"]: Refreshing state... [id=2]
jamfpro_dock_item.test["charlie"]: Refreshing state... [id=2]
jamfpro_script.test["bravo"]: Refreshing state... [id=3]
jamfpro_printer.test["alpha"]: Refreshing state... [id=1]
jamfpro_dock_item.test["bravo"]: Refreshing state... [id=1]
jamfpro_script.test["charlie"]: Refreshing state... [id=1]
jamfpro_dock_item.test["alpha"]: Refreshing state... [id=3]
jamfpro_script.test["alpha"]: Refreshing state... [id=2]
jamfpro_category.test["alpha"]: Refreshing state... [id=3]
jamfpro_category.test["bravo"]: Refreshing state... [id=1]
jamfpro_printer.test["bravo"]: Refreshing state... [id=3]
jamfpro_printer.test["charlie"]: Refreshing state... [id=2]
jamfpro_policy.test: Refreshing state... [id=12]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration
and found no differences, so no changes are needed.
```

Exit code: **0**. Scripts, local accounts, Self Service categories, printers, and
Dock items each contain three entries, and declaration order causes no diff.

#### 6. Change actual script parameters

With the phase 5 variables, add `-var=script_parameter=updated`:

```console
$ TF_CLI_CONFIG_FILE=current.tfrc terraform plan -input=false -no-color -out=content-change-fixed.tfplan -var-file=<credentials> -var=include_printers_and_dock_items=true -var='item_order=["bravo","alpha","charlie"]' -var=script_parameter=updated
Terraform will perform the following actions:

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

Only the existing policy updates; script objects are not recreated. The nested
remove/add entries are Terraform's set representation of changed elements.

```console
$ TF_CLI_CONFIG_FILE=current.tfrc terraform apply -input=false -no-color content-change-fixed.tfplan
jamfpro_policy.test: Modifying... [id=12]
jamfpro_policy.test: Modifications complete after 0s [id=12]

Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

Assertions confirmed policy ID `12`, exactly three scripts, and parameters
`alpha-updated`, `bravo-updated`, and `charlie-updated`; the policy stayed disabled.

#### 7. Final unchanged-config refresh and plan

```console
$ TF_CLI_CONFIG_FILE=current.tfrc terraform plan -detailed-exitcode -input=false -no-color -var-file=<credentials> -var=include_printers_and_dock_items=true -var='item_order=["bravo","alpha","charlie"]' -var=script_parameter=updated
jamfpro_dock_item.test["charlie"]: Refreshing state... [id=2]
jamfpro_printer.test["bravo"]: Refreshing state... [id=3]
jamfpro_category.test["charlie"]: Refreshing state... [id=2]
jamfpro_script.test["alpha"]: Refreshing state... [id=2]
jamfpro_printer.test["alpha"]: Refreshing state... [id=1]
jamfpro_script.test["bravo"]: Refreshing state... [id=3]
jamfpro_dock_item.test["bravo"]: Refreshing state... [id=1]
jamfpro_dock_item.test["alpha"]: Refreshing state... [id=3]
jamfpro_printer.test["charlie"]: Refreshing state... [id=2]
jamfpro_script.test["charlie"]: Refreshing state... [id=1]
jamfpro_category.test["bravo"]: Refreshing state... [id=1]
jamfpro_category.test["alpha"]: Refreshing state... [id=3]
jamfpro_policy.test: Refreshing state... [id=12]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration
and found no differences, so no changes are needed.
```

Exit code: **0**.

#### 8. Cleanup

Using the same phase 7 variables, the saved destroy plan contained only the
13 temporary resources created for this test:

```console
$ TF_CLI_CONFIG_FILE=current.tfrc terraform plan -destroy -input=false -no-color -out=cleanup.tfplan -var-file=<credentials> -var=include_printers_and_dock_items=true -var='item_order=["bravo","alpha","charlie"]' -var=script_parameter=updated
Plan: 0 to add, 0 to change, 13 to destroy.

$ TF_CLI_CONFIG_FILE=current.tfrc terraform apply -input=false -no-color cleanup.tfplan
jamfpro_policy.test: Destroying... [id=12]
jamfpro_policy.test: Destruction complete after 0s
jamfpro_printer.test["bravo"]: Destroying... [id=3]
jamfpro_script.test["charlie"]: Destroying... [id=1]
jamfpro_dock_item.test["alpha"]: Destroying... [id=3]
jamfpro_category.test["charlie"]: Destroying... [id=2]
jamfpro_category.test["alpha"]: Destroying... [id=3]
jamfpro_dock_item.test["bravo"]: Destroying... [id=1]
jamfpro_printer.test["charlie"]: Destroying... [id=2]
jamfpro_category.test["bravo"]: Destroying... [id=1]
jamfpro_script.test["alpha"]: Destroying... [id=2]
jamfpro_dock_item.test["charlie"]: Destroying... [id=2]
jamfpro_category.test["bravo"]: Destruction complete after 1s
jamfpro_category.test["charlie"]: Destruction complete after 1s
jamfpro_script.test["charlie"]: Destruction complete after 1s
jamfpro_printer.test["alpha"]: Destroying... [id=1]
jamfpro_category.test["alpha"]: Destruction complete after 1s
jamfpro_script.test["bravo"]: Destroying... [id=3]
jamfpro_dock_item.test["bravo"]: Destruction complete after 1s
jamfpro_dock_item.test["charlie"]: Destruction complete after 1s
jamfpro_dock_item.test["alpha"]: Destruction complete after 1s
jamfpro_printer.test["charlie"]: Destruction complete after 1s
jamfpro_printer.test["bravo"]: Destruction complete after 1s
jamfpro_script.test["alpha"]: Destruction complete after 1s
jamfpro_printer.test["alpha"]: Destruction complete after 0s
jamfpro_script.test["bravo"]: Destruction complete after 1s

Apply complete! Resources: 0 added, 0 changed, 13 destroyed.
```

The final Terraform state contains **zero resources**. All 13 temporary
resources reported successful deletion.

#### Verification scope

Live v1 → v2 migration covered scripts, local accounts, and Self Service
categories. Printers and Dock items were subsequently created and exercised
with the new provider; all five populated collections passed the ordering-only
plan. Package and directory-binding collections, and the v0 migration path,
are covered by automated tests rather than live provisioning.

## Quality Checklist

- [x] I have reviewed my own code before requesting review
- [ ] I have verified there are no other open Pull Requests for the same update/change
- [ ] All CI/CD pipelines pass without errors or warnings
- [x] My code follows the established style guidelines of this project
- [x] My comments are used only when necessary, ideally where the codes purpose is not self explanatory (eg: necessary magic numbers)
- [x] I have added necessary documentation (if appropriate)
- [x] I have made corresponding changes to the README and other relevant documentation
- [ ] My changes generate no new warnings

The fixture README, migration guide, and generated policy reference are updated.
The local Terraform runs display the expected development-override warning;
CI completion and absence of all warnings are not claimed.

## Screenshots/Recordings (if appropriate)

Not applicable. The Testing section contains the actual HCL, pre-change drift
plan, post-change plan/apply output, and cleanup results.

## Additional Notes

The migration retains resource identity and automatically persists schema version
2 when applied. Users must update numeric-index expressions for changed
collections, and cannot downgrade a v2 state directly to an older provider.
The migration guide describes these changes.
