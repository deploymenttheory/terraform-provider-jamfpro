
# Jamf Cloud Distribution Service
resource "jamfpro_cloud_distribution_point" "jamf_cloud" {
  cdn_type = "JAMF_CLOUD"
  master   = true
}

# Akamai
resource "jamfpro_cloud_distribution_point" "akamai" {
  cdn_type = "AKAMAI"
  master   = true

  username     = "akamai_user"
  password     = "akamai_password"
  directory    = "/tf/example"
  upload_url   = "sftp://upload.example.akamai.com/netstorage"
  download_url = "https://downloads.example.akamai.com/edge"

  secondary_auth_required     = true
  secondary_auth_status_code  = 200
  secondary_auth_time_to_live = 3600
}

# Amazon S3 with Signed URLs
resource "jamfpro_cloud_distribution_point" "amazon_signed" {
  cdn_type = "AMAZON_S3"
  master   = true

  username = "amazon_user"
  password = "amazon_password"

  require_signed_urls = true
  key_pair_id         = "APKAIEXAMPLE"
  expiration_seconds  = 3600
  private_key         = file("/path/to/private/key.pem")
}
