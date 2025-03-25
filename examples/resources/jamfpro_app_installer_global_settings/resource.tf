resource "jamfpro_app_installer_global_settings" "jamfpro_app_installer_singleton" {
  notification_message = "A new version of this app is ready to install." # This message is shown to users when an app update becomes available.
  notification_interval = 24 # Interval (in hours) between repeated user notifications about the update.
  deadline_message = "Please install the update to continue using the app." # Message displayed when the deadline to update has passed.
  deadline = 48 # Deadline (in hours) after which the app is forcibly quit and updated.
  quit_delay = 5 # Additional time (in minutes) for users to save their work before the app quits.
  complete_message = "Installation complete. You may now relaunch the app." # Message shown after successful installation.
  relaunch = true # If true, the app will automatically relaunch after installation.
  suppress = false # If true, no notifications will be shown to the end user.
}
