package ssocertificate

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProSSOCertificate provides information about Jamf Pro's SSO Certificate configuration
func DataSourceJamfProSSOCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"keystore": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key identifier for the certificate",
						},
						"keys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the certificate key",
									},
									"valid": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the certificate key is valid",
									},
								},
							},
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the keystore (e.g., PKCS12, JKS)",
						},
						"keystore_file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The filename of the keystore",
						},
						"keystore_setup_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The setup type of the keystore",
						},
					},
				},
				Description: "The keystore configuration for the SSO certificate",
			},
			"keystore_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keys": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of keys in the keystore",
						},
						"issuer": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The issuer of the certificate",
						},
						"subject": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subject of the certificate",
						},
						"expiration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expiration date of the certificate",
						},
						"serial_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The serial number of the certificate",
						},
					},
				},
				Description: "Detailed information about the SSO certificate",
			},
		},
	}
}

// dataSourceRead reads the SSO certificate settings from Jamf Pro
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	var response *jamfpro.ResourceSSOKeystoreResponse
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		response, apiErr = client.GetSSOCertificate()
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read SSO certificate settings: %v", err))
	}

	d.SetId(fmt.Sprintf("jamfpro_sso_certificate_%s", response.Keystore.KeystoreFileName))

	if err := d.Set("keystore", []interface{}{
		map[string]interface{}{
			"key":                 response.Keystore.Key,
			"type":                response.Keystore.Type,
			"keystore_file_name":  response.Keystore.KeystoreFileName,
			"keystore_setup_type": response.Keystore.KeystoreSetupType,
			"keys":                flattenCertKeys(response.Keystore.Keys),
		},
	}); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	if response.KeystoreDetails != nil {
		if err := d.Set("keystore_details", []interface{}{
			map[string]interface{}{
				"keys":          response.KeystoreDetails.Keys,
				"issuer":        response.KeystoreDetails.Issuer,
				"subject":       response.KeystoreDetails.Subject,
				"expiration":    response.KeystoreDetails.Expiration,
				"serial_number": response.KeystoreDetails.SerialNumber,
			},
		}); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}

	return diags
}
