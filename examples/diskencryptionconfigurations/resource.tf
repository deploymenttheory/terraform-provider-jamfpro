
// jamfpro Institutional Recovery Key config tf example 

resource "jamfpro_disk_encryption_configurations" "disk_encryption_configuration_01" {
  name                      = "jamfpro-tf-example-InstitutionalRecoveryKey-config"
  key_type                  = "Institutional" # Or "Individual and Institutional"
  file_vault_enabled_users  = "Management Account" # Or "Current or Next User"

  institutional_recovery_key {
    certificate_type  = "PKCS12" # For .p12 certificate types
    password          = "secretThing"
    data              = filebase64("/Users/dafyddwatkins/localtesting/support_files/filevaultcertificate/FileVaultMaster-sdk.p12")
  }
  // add / remove this lifecycle block when you want to make updates to this resource.
  // the req will fail if this block is not present during a resource change
  // block is required once resource has been created to suppress false positives in tf plan.
  
  lifecycle {
    ignore_changes = [institutional_recovery_key[0].password]
  }
  
}

// jamfpro Individual Recovery Key config tf example 

resource "jamfpro_disk_encryption_configurations" "disk_encryption_configuration_02" {
  name                      = "jamfpro-tf-example-IndividualRecoveryKey-config"
  key_type                  = "Individual" 
  file_vault_enabled_users  = "Management Account" # Or "Current or Next User"

}