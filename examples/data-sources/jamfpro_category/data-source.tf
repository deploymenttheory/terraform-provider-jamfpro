# Example 1: Look up by ID
data "jamfpro_category" "by_id" {
  id = "1"
}

# Example 2: Look up by name
data "jamfpro_category" "by_name" {
  name = "Applications"
}

# Example 3: Using variables
variable "category_name" {
  type        = string
  description = "Category to look up"
  default     = "Applications"
}

data "jamfpro_category" "dynamic" {
  name = var.category_name
}

# Output examples
output "category_details" {
  value = {
    id       = data.jamfpro_category.by_name.id
    name     = data.jamfpro_category.by_name.name
    priority = data.jamfpro_category.by_name.priority
  }
}

# Example 4: Using in another resource 
resource "jamfpro_policy" "app_policy" {
  name        = "Application Install Policy"
  category_id = data.jamfpro_category.by_name.id
}