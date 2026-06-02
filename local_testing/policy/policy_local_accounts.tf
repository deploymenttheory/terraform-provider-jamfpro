resource "jamfpro_policy" "jamfpro_policy_local_account_example" {
  name                          = "acc-test-policy-local-account"
  enabled                       = false
  trigger_checkin               = false
  trigger_enrollment_complete   = false
  trigger_login                 = false
  trigger_network_state_changed = false
  trigger_startup               = false
  trigger_other                 = "EVENT" // "USER_INITIATED" for self service trigger , "EVENT" for an event trigger
  frequency                     = "Once per computer"
  retry_event                   = "none"
  retry_attempts                = -1
  notify_on_each_failed_retry   = false
  target_drive                  = "/"
  offline                       = false
  category_id                   = -1
  site_id                       = -1

  network_limitations {
    minimum_network_connection = "No Minimum"
    any_ip_address             = false
  }

  scope {
    all_computers = false
    all_jss_users = false
  }

  payloads {
    account_maintenance {
      local_accounts {
        account {
          action                    = "Create"
          username                  = "testuser"
          realname                  = "Test User"
          archive_home_directory    = false
          archive_home_directory_to = ""
          home                      = "/Users/testuser"
          password                  = "SuperSecretPassword123!"
          admin                     = true
          secure_token_allowed      = true
        }
      }
    }
  }

}
