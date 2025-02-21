# enrollment customization with text prestage pane only
resource "jamfpro_enrollment_customization" "text_only" {
  site_id      = "-1"  # -1 for None
  display_name = "Corporate Enrollment - Text"
  description  = "Corporate enrollment with welcome message"
  enrollment_customization_image_source = "/path/to/your/logo.png"

  branding_settings {
    text_color        = "000000"  # Black text
    button_color      = "0066CC"  # Blue buttons
    button_text_color = "FFFFFF"  # White button text
    background_color  = "F5F5F5"  # Light gray background
  }

  text_pane {
    display_name        = "Welcome Message"
    rank                = 1
    title               = "Welcome to Our Company"
    body                = "We're excited to get your device set up with all the tools you need to be productive."
    subtext             = "This process should take about 10 minutes to complete."
    back_button_text    = "Back"
    continue_button_text = "Continue"
  }
}

# enrollment customization with ldap prestage pane only
resource "jamfpro_enrollment_customization" "ldap_only" {
  site_id      = "-1"
  display_name = "Corporate Enrollment - LDAP"
  description  = "Corporate enrollment with LDAP authentication"
  enrollment_customization_image_source = "/path/to/your/logo.png"

  branding_settings {
    text_color        = "000000"
    button_color      = "0066CC"
    button_text_color = "FFFFFF"
    background_color  = "F5F5F5"
  }

  ldap_pane {
    display_name        = "Company Authentication"
    rank                = 1
    title               = "Enter Your Network Credentials"
    username_label      = "Network Username"
    password_label      = "Network Password"
    back_button_text    = "Back"
    continue_button_text = "Authenticate"
    
    ldap_group_access {
      group_name     = "IT-Department"
      ldap_server_id = 1
    }
    
    ldap_group_access {
      group_name     = "Engineering"
      ldap_server_id = 1
    }
  }
}

# enrollment customization with sso prestage pane only
resource "jamfpro_enrollment_customization" "sso_only" {
  site_id      = "-1"
  display_name = "Corporate Enrollment - SSO"
  description  = "Corporate enrollment with Single Sign-On"
  enrollment_customization_image_source = "/path/to/your/logo.png"

  branding_settings {
    text_color        = "000000"
    button_color      = "0066CC"
    button_text_color = "FFFFFF"
    background_color  = "F5F5F5"
  }

  sso_pane {
    display_name                     = "Corporate SSO"
    rank                             = 1
    is_group_enrollment_access_enabled = true
    group_enrollment_access_name     = "All-Employees"
    is_use_jamf_connect              = true
    short_name_attribute             = "sAMAccountName"
    long_name_attribute              = "displayName"
  }
}

# enrollment customization with text and sso prestage panes
resource "jamfpro_enrollment_customization" "text_ldap" {
  site_id      = "-1"
  display_name = "Corporate Enrollment - Text and LDAP"
  description  = "Corporate enrollment with welcome message and LDAP authentication"
  enrollment_customization_image_source = "/path/to/your/logo.png"

  branding_settings {
    text_color        = "000000"
    button_color      = "0066CC"
    button_text_color = "FFFFFF"
    background_color  = "F5F5F5"
  }

  text_pane {
    display_name        = "Welcome Message"
    rank                = 1
    title               = "Welcome to Our Company"
    body                = "We're excited to get your device set up. Please authenticate using your network credentials on the next screen."
    subtext             = "This process should take about 10 minutes to complete."
    back_button_text    = "Back"
    continue_button_text = "Continue"
  }

  ldap_pane {
    display_name        = "Company Authentication"
    rank                = 2
    title               = "Enter Your Network Credentials"
    username_label      = "Network Username"
    password_label      = "Network Password"
    back_button_text    = "Back"
    continue_button_text = "Authenticate"
    
    ldap_group_access {
      group_name     = "All-Users"
      ldap_server_id = 1
    }
  }
}

# enrollment customization with text and ldap prestage panes

resource "jamfpro_enrollment_customization" "text_sso" {
  site_id      = "-1"
  display_name = "Corporate Enrollment - Text and SSO"
  description  = "Corporate enrollment with welcome message and Single Sign-On"
  enrollment_customization_image_source = "/path/to/your/logo.png"

  branding_settings {
    text_color        = "000000"
    button_color      = "0066CC"
    button_text_color = "FFFFFF"
    background_color  = "F5F5F5"
  }

  text_pane {
    display_name        = "Welcome Message"
    rank                = 1
    title               = "Welcome to Our Company"
    body                = "We're excited to get your device set up. Please sign in with your corporate identity on the next screen."
    subtext             = "This process should take about 10 minutes to complete."
    back_button_text    = "Back"
    continue_button_text = "Continue"
  }

  sso_pane {
    display_name                     = "Corporate SSO"
    rank                             = 2
    is_group_enrollment_access_enabled = true
    group_enrollment_access_name     = "All-Employees"
    is_use_jamf_connect              = true
    short_name_attribute             = "sAMAccountName"
    long_name_attribute              = "displayName"
  }
}