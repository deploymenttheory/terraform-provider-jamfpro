# Example 1: Look up building by ID
data "jamfpro_building" "by_id" {
  id = "1"
}

# Example 2: Look up building by name
data "jamfpro_building" "by_name" {
  name = "Corporate HQ"
}

# Output examples
output "building_address" {
  value = "${data.jamfpro_building.by_name.street_address1}, ${data.jamfpro_building.by_name.city}"
}

output "building_details" {
  value = {
    name = data.jamfpro_building.by_name.name
    full_address = join("\n", [
      data.jamfpro_building.by_name.street_address1,
      data.jamfpro_building.by_name.street_address2,
      data.jamfpro_building.by_name.city,
      data.jamfpro_building.by_name.state_province,
      data.jamfpro_building.by_name.zip_postal_code,
      data.jamfpro_building.by_name.country
    ])
  }
}

# Example 3: Using with variables
variable "building_name" {
  type        = string
  description = "The name of the building to look up"
  default     = "Corporate HQ"
}

data "jamfpro_building" "dynamic" {
  name = var.building_name
}

# Example 4: Using in another resource
resource "jamfpro_computer" "office_computer" {
  name         = "Office Workstation"
  building_id  = data.jamfpro_building.by_name.id
  description  = "Workstation located at ${data.jamfpro_building.by_name.name}"
}