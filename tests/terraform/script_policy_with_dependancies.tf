# Test Description: 
# - Create Independent computer extension attributes with no dependencies
# - Create Independent sites with no dependencies
# - Create Independent categories with no dependencies
# - Create Dependent scripts with categories as dependencies
# - Create Dependent policies with scripts and categories as dependencies

locals {
  count     = 10
  base_name = "tf-mutex"
}

# Create extension attributes separate in isolation from policies. No interdependancies
resource "jamfpro_computer_extension_attribute" "computer_extension_attribute_array" {
  count                  = local.count
  name                   = "${local.base_name}-ea-${format("%03d", count.index + 1)}"
  enabled                = true
  description            = "Extension attribute to track user migration state"
  inventory_display_type = "EXTENSION_ATTRIBUTES"
  data_type              = "STRING"
  input_type             = "SCRIPT"
  script_contents        = "#!/bin/bash\n\n# Script: Migration State Check\n# Purpose: Checks and reports the current state of user migration\n# Created: 2024-11-26\n\n################################\n# VARIABLES\n################################\nLOGGED_IN_USER=$(/usr/bin/stat -f%Su \"/dev/console\")\nMIGRATION_STATE=\"Not Started\"\n\n################################\n# FUNCTIONS\n################################\ncheck_migration_state() {\n    local user=\"$1\"\n    # Add your migration state check logic here\n    # This is a placeholder that you can customize based on your needs\n    if [ -f \"/Users/$user/.migration_complete\" ]; then\n        MIGRATION_STATE=\"Completed\"\n    elif [ -f \"/Users/$user/.migration_in_progress\" ]; then\n        MIGRATION_STATE=\"In Progress\"\n    fi\n}\n\n################################\n# MAIN LOGIC\n################################\nif [ ! -z \"$LOGGED_IN_USER\" ] && [ \"$LOGGED_IN_USER\" != \"root\" ]; then\n    check_migration_state \"$LOGGED_IN_USER\"\nfi\n\n################################\n# OUTPUT\n################################\necho \"<result>$MIGRATION_STATE</result>\"\nexit 0"
}

# Create sites
resource "jamfpro_site" "jamfpro_site_array" {
  count = local.count
  name  = "${local.base_name}-site-${format("%03d", count.index + 1)}"
}

# Create categories
resource "jamfpro_category" "jamfpro_category_array" {
  count    = local.count
  name     = "${local.base_name}-category-${format("%03d", count.index + 1)}"
  priority = 5
}

# Create scripts
resource "jamfpro_script" "jamfpro_script_array" {
  count           = local.count
  name            = "${local.base_name}-add-or-remove-group-membership-v4.0-${format("%03d", count.index + 1)}"
  script_contents = file("${path.module}/support_files/scripts/Add or Remove Group Membership.zsh")
  category_id     = jamfpro_category.jamfpro_category_array[count.index].id
  os_requirements = "13"
  priority        = "BEFORE"
  info            = "Adds target user or group to specified group membership, or removes said membership."
  notes           = "Jamf Pro script parameters: 4 -> 7"
  parameter4      = "100"           // targetID
  parameter5      = "group"         // Target Type - Must be either "user" or "group"
  parameter6      = "someGroupName" // targetMembership
  parameter7      = "add"           // Script Action - Must be either "add" or "remove"
}

# Create corresponding policies
resource "jamfpro_policy" "jamfpro_policy_array" {
  count                         = local.count
  name                          = "${local.base_name}-policy-script-${format("%03d", count.index + 1)}"
  enabled                       = false
  trigger_checkin               = false
  trigger_enrollment_complete   = false
  trigger_login                 = false
  trigger_network_state_changed = false
  trigger_startup               = false
  trigger_other                 = "EVENT"
  frequency                     = "Once per computer"
  retry_event                   = "none"
  retry_attempts                = -1
  notify_on_each_failed_retry   = false
  target_drive                  = "/"
  offline                       = false
  category_id                   = jamfpro_category.jamfpro_category_array[count.index].id
  site_id                       = jamfpro_site.jamfpro_site_array[count.index].id

  network_limitations {
    minimum_network_connection = "No Minimum"
    any_ip_address             = false
  }

  scope {
    all_computers = false
    all_jss_users = false
  }

  self_service {
    use_for_self_service            = true
    self_service_display_name       = ""
    install_button_text             = "Install"
    self_service_description        = ""
    force_users_to_view_description = false
    feature_on_main_page            = false
  }

  payloads {
    scripts {
      id          = jamfpro_script.jamfpro_script_array[count.index].id
      priority    = "After"
      parameter4  = "param_value_4"
      parameter5  = "param_value_5"
      parameter6  = "param_value_6"
      parameter7  = "param_value_7"
      parameter8  = "param_value_8"
      parameter9  = "param_value_9"
      parameter10 = "param_value_10"
      parameter11 = "param_value_11"
    }
  }
}

# Outputs for verification
output "site_details" {
  description = "Details of all created sites"
  value = [
    for i in range(local.count) : {
      name = jamfpro_site.jamfpro_site_array[i].name
      id   = jamfpro_site.jamfpro_site_array[i].id
    }
  ]
}

output "category_details" {
  description = "Details of all created categories"
  value = [
    for i in range(local.count) : {
      name     = jamfpro_category.jamfpro_category_array[i].name
      id       = jamfpro_category.jamfpro_category_array[i].id
      priority = jamfpro_category.jamfpro_category_array[i].priority
    }
  ]
}

output "script_details" {
  description = "Details of all created scripts"
  value = [
    for i in range(local.count) : {
      name        = jamfpro_script.jamfpro_script_array[i].name
      id          = jamfpro_script.jamfpro_script_array[i].id
      category_id = jamfpro_script.jamfpro_script_array[i].category_id
    }
  ]
}

output "policy_details" {
  description = "Details of all created policies"
  value = [
    for i in range(local.count) : {
      name        = jamfpro_policy.jamfpro_policy_array[i].name
      category_id = jamfpro_policy.jamfpro_policy_array[i].category_id
      script_id   = jamfpro_script.jamfpro_script_array[i].id
      site_id     = jamfpro_policy.jamfpro_policy_array[i].site_id
    }
  ]
}