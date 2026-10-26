resource "jamfpro_user_group" "test" {
 name = "tf-sets-user-group-20260919"
 is_smart = false
 assigned_user_ids = [4,1,2]
 site_id = -1
}
