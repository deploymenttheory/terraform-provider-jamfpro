resource "jamfpro_account_group" "test" {
 name = "tf-sets-account-group-20260919"
 access_level = "Full Access"
 privilege_set = "Custom"
 member_ids = [3,4,2]
 site_id = -1
}
