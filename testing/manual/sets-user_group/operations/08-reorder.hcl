resource "jamfpro_user_group" "operations" {
 name = "tf-sets-user-operations-20260919"
 is_smart = false
 site_id = -1
 user_additions { id = "1" }
 user_additions { id = "2" }

 user_additions { id = "3" }
 user_deletions { id = "1" }
 user_deletions { id = "2" }
}
