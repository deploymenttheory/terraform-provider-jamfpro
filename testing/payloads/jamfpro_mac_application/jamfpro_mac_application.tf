resource "jamfpro_mac_application" "slack" {
  name      = "Slack for Desktop"
  version   = "4.37.101"
  bundle_id = "com.tinyspeck.slackmacgap"
  url       = "https://apps.apple.com/gb/app/slack-for-desktop/id803453959?mt=12"

  site_id     = -1
  category_id = -1

  deployment_type = "Make Available in Self Service"

  scope {
    all_computers = false
    all_jss_users = false
  }

  self_service {
    install_button_text             = "Install"
    self_service_description        = <<-EOT
Slack brings team communication and collaboration into one place so you can get more work done, whether you belong to a large enterprise or a small business. Check off your to-do list and move your projects forward by bringing the right people, conversations, tools, and information you need together. Slack is available on any device, so you can find and access your team and your work, whether you're at your desk or on the go.

Use Slack to: 
• Communicate with your team and organize your conversations by topics, projects, or anything else that matters to your work
• Message or call any person or group within your team
• Share and edit documents and collaborate with the right people all in Slack 
• Integrate into your workflow, the tools and services you already use including Google Drive, Salesforce, Dropbox, Asana, Twitter, Zendesk, and more
• Easily search a central knowledge base that automatically indexes and archives your team's past conversations and files
• Customize your notifications so you stay focused on what matters

Scientifically proven (or at least rumored) to make your working life simpler, more pleasant, and more productive. We hope you'll give Slack a try.

Stop by and learn more at: https://slack.com/
    EOT
    force_users_to_view_description = false
    feature_on_main_page            = false
    notification                    = "Self Service"
  }

  vpp {
    assign_vpp_device_based_licenses = false
    vpp_admin_account_id             = -1
  }
}
