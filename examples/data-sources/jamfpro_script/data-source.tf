# Example 1: Look up extension attribute by ID
data "jamfpro_computer_extension_attribute" "by_id" {
  id = "1"
}

# Example 2: Look up extension attribute by name
data "jamfpro_computer_extension_attribute" "custom_attribute" {
  name = "Custom Hardware Attribute"
}

# Example 3: Using variables
variable "attribute_name" {
  type        = string
  description = "Name of the computer extension attribute to look up"
  default     = "Software Version"
}

data "jamfpro_computer_extension_attribute" "dynamic" {
  name = var.attribute_name
}

# Example 4: Output examples
output "attribute_details" {
  value = {
    name        = data.jamfpro_computer_extension_attribute.custom_attribute.name
    description = data.jamfpro_computer_extension_attribute.custom_attribute.description
    data_type   = data.jamfpro_computer_extension_attribute.custom_attribute.data_type
    input_type  = data.jamfpro_computer_extension_attribute.custom_attribute.input_type
  }
}

# Example 5: Using with conditions
data "jamfpro_computer_extension_attribute" "script_attribute" {
  name = "Custom Script Attribute"

  lifecycle {
    postcondition {
      condition     = self.input_type == "SCRIPT" && self.enabled == true
      error_message = "Script attribute must be enabled and of type SCRIPT"
    }
  }
}

# Example 6: Using in another resource
resource "jamfpro_computer_group" "smart_group" {
  name         = "Computers with Custom Attribute"
  is_smart     = true
  attribute_id = data.jamfpro_computer_extension_attribute.custom_attribute.id

  criteria {
    name     = data.jamfpro_computer_extension_attribute.custom_attribute.name
    operator = "is"
    value    = "specific_value"
  }
}

# Example 7: Working with LDAP attributes
data "jamfpro_computer_extension_attribute" "ldap_attribute" {
  name = "LDAP User Details"

  lifecycle {
    postcondition {
      condition     = self.input_type == "DIRECTORY_SERVICE_ATTRIBUTE_MAPPING"
      error_message = "Attribute must be of type DIRECTORY_SERVICE_ATTRIBUTE_MAPPING"
    }
  }
}

output "ldap_mapping" {
  value = {
    attribute = data.jamfpro_computer_extension_attribute.ldap_attribute.ldap_attribute_mapping
    allowed   = data.jamfpro_computer_extension_attribute.ldap_attribute.ldap_extension_attribute_allowed
  }
}