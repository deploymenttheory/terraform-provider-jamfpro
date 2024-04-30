resource "jamfpro_macos_configuration_profile" "jamfpro_macos_configuration_profile_001" {
  name                = "your-macos_configuration_profile-name"
  distribution_method = "Install Automatically" // Can be "Make Available in Self Service" or "Install Automatically"
  user_removeable     = false
  level               = "User" // Can be "Computer", "User", "System"
  payload             = file("${path.module}/path/to/your.mobileconfig")
  category {
    id = -1 // -1 for no category or use the category resource id reference. 
  }
  site { // optional block
    id   = 1
    name = "site name"
  }
  scope {
    all_computers      = false
    computer_ids       = sort([17, 18]) // uses computer resource id
    computer_group_ids = sort([53])     // uses computer group resource id
    jss_user_ids       = [1, 2, 3]
    jss_user_group_ids = [1, 2, 3]
    building_ids       = sort([1, 2, 3])
    department_ids     = sort([1, 2, 3])

    limitations {
      network_segment_ids = [1, 2, 3]
      ibeacon_ids = [1, 2, 3]
    }

    exclusions {
      network_segment_ids = [1, 2, 3]
      ibeacon_ids = [1, 2, 3]
      department_ids = [1, 2, 3]
    }
  }
}

resource "jamfpro_macos_configuration_profile" "jamfpro_macos_configuration_profile_multi" {
  count               = 5
  name                = "tf-localtest-macosconfigprofile-dockitems-multi-${count.index}"
  distribution_method = "Install Automatically"
  payload             = file("${path.module}/support_files/configurationprofiles/dockitems-chara-nosub-test.mobileconfig")
  category {
    id = -1
  }
  scope {
    all_computers      = false
    computer_ids       = sort([17, 18])
    computer_group_ids = sort([53])
    jss_user_ids       = [4]
    jss_user_group_ids = [4]

    exclusions {
      department_ids = [27653]
    }
  }
}


resource "jamfpro_macos_configuration_profile" "jamfpro_macos_configuration_profile_001" {
  name                = "your-macos_configuration_profile-name"
  distribution_method = "Install Automatically" // Can be "Make Available in Self Service" or "Install Automatically"
  user_removeable     = false
  level               = "User" // Can be "Computer", "User", "System"
  payload             = file("${path.module}/path/to/your.mobileconfig")
  category {
    id = -1 // -1 for no category or use the category resource id reference. 
  }
  site { // optional block
    id   = 1
    name = "site name"
  }
  scope {
    all_computers      = false
    computer_ids       = sort([17, 18]) // uses computer resource id
    computer_group_ids = sort([53])     // uses computer group resource id
    jss_user_ids       = [1, 2, 3]
    jss_user_group_ids = [1, 2, 3]
    building_ids       = sort([1, 2, 3])
    department_ids     = sort([1, 2, 3])

    limitations {
      network_segment_ids = [1, 2, 3]
      ibeacon_ids = [1, 2, 3]
    }

    exclusions {
      network_segment_ids = [1, 2, 3]
      ibeacon_ids = [1, 2, 3]
      department_ids = [1, 2, 3]
    }
  }
}