resource "jamfpro_smart_computer_group_v3" "test" {
  name        = "tf-testing-${var.testing_id}-smart-v3-${random_id.rng.hex}"
  description = "Framework smart computer group v3 lifecycle test."
  site_id     = "-1"

  criteria {
    name        = "Serial Number"
    priority    = 0
    search_type = "is"
    value       = "tf-testing-${var.testing_id}"
  }

  criteria {
    name        = "Operating System Version"
    priority    = 1
    and_or      = "and"
    search_type = "like"
    value       = "26."
  }

  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}

data "jamfpro_smart_computer_group_v3" "by_id" {
  id = jamfpro_smart_computer_group_v3.test.id
}

data "jamfpro_smart_computer_group_v3" "by_name" {
  name       = jamfpro_smart_computer_group_v3.test.name
  depends_on = [jamfpro_smart_computer_group_v3.test]
}
