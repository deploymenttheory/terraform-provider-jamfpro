// ========================================================================== //
// Restricted Software
// ========================================================================== //

// ========================================================================== //
// Supporting Resources for Scope Testing
// ========================================================================== //

resource "jamfpro_building" "restricted_software_building_001" {
  name = "tf-testing-${var.testing_id}-rs-building-001-${random_id.rng.hex}"
}

resource "jamfpro_building" "restricted_software_building_002" {
  name = "tf-testing-${var.testing_id}-rs-building-002-${random_id.rng.hex}"
}

resource "jamfpro_department" "restricted_software_department_001" {
  name = "tf-testing-${var.testing_id}-rs-department-001-${random_id.rng.hex}"
}

resource "jamfpro_department" "restricted_software_department_002" {
  name = "tf-testing-${var.testing_id}-rs-department-002-${random_id.rng.hex}"
}

resource "jamfpro_smart_computer_group" "restricted_software_group_001" {
  name = "tf-testing-${var.testing_id}-rs-group-001-${random_id.rng.hex}"
  criteria {
    name        = "Operating System Version"
    search_type = "like"
    value       = "14."
    priority    = 0
  }
}

resource "jamfpro_smart_computer_group" "restricted_software_group_002" {
  name = "tf-testing-${var.testing_id}-rs-group-002-${random_id.rng.hex}"
  criteria {
    name        = "Operating System Version"
    search_type = "like"
    value       = "15."
    priority    = 0
  }
}

// ========================================================================== //
// Minimal Restricted Software - Only required fields
// ========================================================================== //

resource "jamfpro_restricted_software" "restricted_software_min" {
  name         = "tf-testing-${var.testing_id}-rs-min-${random_id.rng.hex}"
  process_name = "SomeUnwantedApp.app"
}

// ========================================================================== //
// Restricted Software with All Computers Scope
// ========================================================================== //

resource "jamfpro_restricted_software" "restricted_software_all_computers" {
  name                     = "tf-testing-${var.testing_id}-rs-all-computers-${random_id.rng.hex}"
  process_name             = "Install macOS High Sierra.app"
  match_exact_process_name = true
  send_notification        = true
  kill_process             = true
  delete_executable        = false
  display_message          = "This macOS installer is restricted. Please contact IT for assistance."

  scope {
    all_computers = true
  }
}

// ========================================================================== //
// Maximal Restricted Software - All fields with scope and exclusions
// ========================================================================== //

resource "jamfpro_restricted_software" "restricted_software_max" {
  name                     = "tf-testing-${var.testing_id}-rs-max-${random_id.rng.hex}"
  process_name             = "BitTorrent.app"
  match_exact_process_name = true
  send_notification        = true
  kill_process             = true
  delete_executable        = true
  display_message          = "BitTorrent is not permitted on corporate devices. The application has been terminated."

  scope {
    all_computers      = false
    computer_group_ids = [jamfpro_smart_computer_group.restricted_software_group_001.id]
    building_ids       = [jamfpro_building.restricted_software_building_001.id]
    department_ids     = [jamfpro_department.restricted_software_department_001.id]

    exclusions {
      computer_group_ids                   = [jamfpro_smart_computer_group.restricted_software_group_002.id]
      building_ids                         = [jamfpro_building.restricted_software_building_002.id]
      department_ids                       = [jamfpro_department.restricted_software_department_002.id]
      directory_service_or_local_usernames = ["admin", "itadmin"]
    }
  }
}

// ========================================================================== //
// Restricted Software with Multiple Scope Targets
// ========================================================================== //

resource "jamfpro_restricted_software" "restricted_software_multi_scope" {
  name                     = "tf-testing-${var.testing_id}-rs-multi-${random_id.rng.hex}"
  process_name             = "uTorrent"
  match_exact_process_name = false
  send_notification        = true
  kill_process             = true
  delete_executable        = false
  display_message          = "Torrent applications are not allowed. Please uninstall this software."

  scope {
    all_computers = false
    computer_group_ids = [
      jamfpro_smart_computer_group.restricted_software_group_001.id,
      jamfpro_smart_computer_group.restricted_software_group_002.id
    ]
    building_ids = [
      jamfpro_building.restricted_software_building_001.id,
      jamfpro_building.restricted_software_building_002.id
    ]
    department_ids = [
      jamfpro_department.restricted_software_department_001.id,
      jamfpro_department.restricted_software_department_002.id
    ]
  }
}
