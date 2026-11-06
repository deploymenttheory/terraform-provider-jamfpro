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
