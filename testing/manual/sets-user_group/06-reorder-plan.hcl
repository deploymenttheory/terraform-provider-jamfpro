resource "jamfpro_user_group" "test" {
 name = "tf-sets-user-group-20260919"
 is_smart = false
 assigned_user_ids = [2,3,1]
 site_id = -1
}
