# Example 1: Look up script by ID
data "jamfpro_script" "by_id" {
  id = "1"
}

# Example 2: Look up script by name
data "jamfpro_script" "by_name" {
  name = "Install Software"
}

# Example 3: Using variables
variable "script_name" {
  type        = string
  description = "Name of the script to look up"
  default     = "System Configuration"
}

data "jamfpro_script" "dynamic" {
  name = var.script_name
}

# Example 4: Output examples
# Example data source outputs
output "script_details" {
  value = {
    id              = data.jamfpro_script.example.id
    name            = data.jamfpro_script.example.name
    category_id     = data.jamfpro_script.example.category_id
    info            = data.jamfpro_script.example.info
    notes           = data.jamfpro_script.example.notes
    os_requirements = data.jamfpro_script.example.os_requirements
    priority        = data.jamfpro_script.example.priority
    script_contents = data.jamfpro_script.example.script_contents
    parameter4      = data.jamfpro_script.example.parameter4
    parameter5      = data.jamfpro_script.example.parameter5
    parameter6      = data.jamfpro_script.example.parameter6
    parameter7      = data.jamfpro_script.example.parameter7
    parameter8      = data.jamfpro_script.example.parameter8
    parameter9      = data.jamfpro_script.example.parameter9
    parameter10     = data.jamfpro_script.example.parameter10
    parameter11     = data.jamfpro_script.example.parameter11
  }
}

# Example 5: Using with conditions
data "jamfpro_script" "deployment_script" {
  name = "Deployment Script"

  lifecycle {
    postcondition {
      condition     = self.priority == "BEFORE" && self.os_requirements != ""
      error_message = "Script must have BEFORE priority and defined OS requirements"
    }
  }
}

# Example 6: Using in another resource
resource "jamfpro_policy" "software_policy" {
  name            = "Software Installation Policy"
  enabled         = true
  script_id       = data.jamfpro_script.by_name.id
  script_priority = data.jamfpro_script.by_name.priority
}

# Example 7: Testing script parameters
data "jamfpro_script" "parameterized_script" {
  name = "Parameterized Installation"
}

output "script_parameters" {
  value = {
    param4  = data.jamfpro_script.parameterized_script.parameter4
    param5  = data.jamfpro_script.parameterized_script.parameter5
    param6  = data.jamfpro_script.parameterized_script.parameter6
    param7  = data.jamfpro_script.parameterized_script.parameter7
    param8  = data.jamfpro_script.parameterized_script.parameter8
    param9  = data.jamfpro_script.parameterized_script.parameter9
    param10 = data.jamfpro_script.parameterized_script.parameter10
    param11 = data.jamfpro_script.parameterized_script.parameter11
  }
}