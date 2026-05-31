// ========================================================================== //
// Cloud Distribution Point
// ========================================================================== //

// JCDS (Jamf Cloud Distribution Service) backed cloud distribution point set as
// the principal (master) distribution point. JAMF_CLOUD requires no CDN
// credentials, so this is the minimal create/teardown fixture; terraform test
// applies it and destroys it at the end of the run.

resource "jamfpro_cloud_distribution_point" "jcds" {
  cdn_type = "JAMF_CLOUD"
  master   = true
}
