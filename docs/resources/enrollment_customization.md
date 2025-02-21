---
page_title: "jamfpro_enrollment_customization"
description: |-
  
---

# jamfpro_enrollment_customization (Resource)


## Example Usage
```terraform
# enrollment customization with text prestage pane only
resource "jamfpro_enrollment_customization" "text_only" {
  site_id                               = "-1" # -1 for None
  display_name                          = "Corporate Enrollment - Text"
  description                           = "Corporate enrollment with welcome message"
  enrollment_customization_image_source = "/path/to/your/logo.png"

  branding_settings {
    text_color        = "000000" # Black text
    button_color      = "0066CC" # Blue buttons
    button_text_color = "FFFFFF" # White button text
    background_color  = "F5F5F5" # Light gray background
  }

  text_pane {
    display_name         = "Welcome Message"
    rank                 = 1
    title                = "Welcome to Our Company"
    body                 = "We're excited to get your device set up with all the tools you need to be productive."
    subtext              = "This process should take about 10 minutes to complete."
    back_button_text     = "Back"
    continue_button_text = "Continue"
  }
}

# enrollment customization with ldap prestage pane only
resource "jamfpro_enrollment_customization" "ldap_only" {
  site_id                               = "-1"
  display_name                          = "Corporate Enrollment - LDAP"
  description                           = "Corporate enrollment with LDAP authentication"
  enrollment_customization_image_source = "/path/to/your/logo.png"

  branding_settings {
    text_color        = "000000"
    button_color      = "0066CC"
    button_text_color = "FFFFFF"
    background_color  = "F5F5F5"
  }

  ldap_pane {
    display_name         = "Company Authentication"
    rank                 = 1
    title                = "Enter Your Network Credentials"
    username_label       = "Network Username"
    password_label       = "Network Password"
    back_button_text     = "Back"
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
  site_id                               = "-1"
  display_name                          = "Corporate Enrollment - SSO"
  description                           = "Corporate enrollment with Single Sign-On"
  enrollment_customization_image_source = "/path/to/your/logo.png"

  branding_settings {
    text_color        = "000000"
    button_color      = "0066CC"
    button_text_color = "FFFFFF"
    background_color  = "F5F5F5"
  }

  sso_pane {
    display_name                       = "Corporate SSO"
    rank                               = 1
    is_group_enrollment_access_enabled = true
    group_enrollment_access_name       = "All-Employees"
    is_use_jamf_connect                = true
    short_name_attribute               = "sAMAccountName"
    long_name_attribute                = "displayName"
  }
}

# enrollment customization with text and sso prestage panes
resource "jamfpro_enrollment_customization" "text_ldap" {
  site_id                               = "-1"
  display_name                          = "Corporate Enrollment - Text and LDAP"
  description                           = "Corporate enrollment with welcome message and LDAP authentication"
  enrollment_customization_image_source = "/path/to/your/logo.png"

  branding_settings {
    text_color        = "000000"
    button_color      = "0066CC"
    button_text_color = "FFFFFF"
    background_color  = "F5F5F5"
  }

  text_pane {
    display_name         = "Welcome Message"
    rank                 = 1
    title                = "Welcome to Our Company"
    body                 = "We're excited to get your device set up. Please authenticate using your network credentials on the next screen."
    subtext              = "This process should take about 10 minutes to complete."
    back_button_text     = "Back"
    continue_button_text = "Continue"
  }

  ldap_pane {
    display_name         = "Company Authentication"
    rank                 = 2
    title                = "Enter Your Network Credentials"
    username_label       = "Network Username"
    password_label       = "Network Password"
    back_button_text     = "Back"
    continue_button_text = "Authenticate"

    ldap_group_access {
      group_name     = "All-Users"
      ldap_server_id = 1
    }
  }
}

# enrollment customization with text and ldap prestage panes

resource "jamfpro_enrollment_customization" "text_sso" {
  site_id                               = "-1"
  display_name                          = "Corporate Enrollment - Text and SSO"
  description                           = "Corporate enrollment with welcome message and Single Sign-On"
  enrollment_customization_image_source = "/path/to/your/logo.png"

  branding_settings {
    text_color        = "000000"
    button_color      = "0066CC"
    button_text_color = "FFFFFF"
    background_color  = "F5F5F5"
  }

  text_pane {
    display_name         = "Welcome Message"
    rank                 = 1
    title                = "Welcome to Our Company"
    body                 = "We're excited to get your device set up. Please sign in with your corporate identity on the next screen."
    subtext              = "This process should take about 10 minutes to complete."
    back_button_text     = "Back"
    continue_button_text = "Continue"
  }

  sso_pane {
    display_name                       = "Corporate SSO"
    rank                               = 2
    is_group_enrollment_access_enabled = true
    group_enrollment_access_name       = "All-Employees"
    is_use_jamf_connect                = true
    short_name_attribute               = "sAMAccountName"
    long_name_attribute                = "displayName"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `branding_settings` (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--branding_settings))
- `description` (String) The description of the enrollment customization.
- `display_name` (String) The display name of the enrollment customization.
- `enrollment_customization_image_source` (String) The .png image source file for upload to the enrollment customization. Recommended: 180x180 pixels and GIF or PNG format

### Optional

- `ldap_pane` (Block List) (see [below for nested schema](#nestedblock--ldap_pane))
- `site_id` (String) The ID of the site associated with the enrollment customization.
- `sso_pane` (Block List) (see [below for nested schema](#nestedblock--sso_pane))
- `text_pane` (Block List) (see [below for nested schema](#nestedblock--text_pane))
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The unique identifier of the enrollment customization.

<a id="nestedblock--branding_settings"></a>
### Nested Schema for `branding_settings`

Required:

- `background_color` (String) The background color in hexadecimal format (6 characters, no # prefix).
- `button_color` (String) The button color in hexadecimal format (6 characters, no # prefix).
- `button_text_color` (String) The button text color in hexadecimal format (6 characters, no # prefix).
- `text_color` (String) The text color in hexadecimal format (6 characters, no # prefix).

Read-Only:

- `icon_url` (String) The URL of the icon image. the format must be 'https://your_jamfUrl/api/v2/enrollment-customizations/images/1'


<a id="nestedblock--ldap_pane"></a>
### Nested Schema for `ldap_pane`

Required:

- `back_button_text` (String) The text for the back button.
- `continue_button_text` (String) The text for the continue button.
- `display_name` (String) The display name of the LDAP pane.
- `password_label` (String) The label for the password field.
- `rank` (Number) The rank/order of the LDAP pane in the enrollment process.
- `title` (String) The title of the LDAP pane.
- `username_label` (String) The label for the username field.

Optional:

- `ldap_group_access` (Block List) (see [below for nested schema](#nestedblock--ldap_pane--ldap_group_access))

Read-Only:

- `id` (Number) The unique identifier of the LDAP pane.

<a id="nestedblock--ldap_pane--ldap_group_access"></a>
### Nested Schema for `ldap_pane.ldap_group_access`

Required:

- `group_name` (String) The name of the LDAP group.
- `ldap_server_id` (Number) The ID of the LDAP server.



<a id="nestedblock--sso_pane"></a>
### Nested Schema for `sso_pane`

Required:

- `display_name` (String) The display name of the SSO pane.
- `rank` (Number) The rank/order of the SSO pane in the enrollment process.

Optional:

- `group_enrollment_access_name` (String) The name of the group for enrollment access.
- `is_group_enrollment_access_enabled` (Boolean) Whether group enrollment access is enabled.
- `is_use_jamf_connect` (Boolean) Whether to use Jamf Connect.
- `long_name_attribute` (String) The attribute to use for long name.
- `short_name_attribute` (String) The attribute to use for short name.

Read-Only:

- `id` (Number) The unique identifier of the SSO pane.


<a id="nestedblock--text_pane"></a>
### Nested Schema for `text_pane`

Required:

- `back_button_text` (String) The text for the back button.
- `body` (String) The main content text of the pane.
- `continue_button_text` (String) The text for the continue button.
- `display_name` (String) The display name of the text pane.
- `rank` (Number) The rank/order of the text pane in the enrollment process.
- `title` (String) The title of the text pane.

Optional:

- `subtext` (String) The subtext content of the pane.

Read-Only:

- `id` (Number) The unique identifier of the text pane.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)