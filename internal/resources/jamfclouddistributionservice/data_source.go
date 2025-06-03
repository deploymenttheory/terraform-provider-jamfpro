package jamfclouddistributionservice

import (
	"context"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceJamfProJamfCloudDistributionService returns a Terraform data source for Jamf Pro Jamf Cloud Distribution Service (JCDS).
func DataSourceJamfProJamfCloudDistributionService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"files": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"length": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"md5": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sha3": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"jcds2_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"file_stream_endpoint_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"max_chunk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

// dataSourceRead fetches the JCDS2 properties and files from Jamf Pro.
func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	files, err := client.GetJCDS2Packages()
	if err != nil {
		return diag.FromErr(err)
	}

	props, err := client.GetJCDS2Properties()
	if err != nil {
		return diag.FromErr(err)
	}

	fileList := make([]interface{}, len(files))
	for i, file := range files {
		fileMap := map[string]interface{}{
			"file_name": file.FileName,
			"length":    file.Length,
			"md5":       file.MD5,
			"region":    file.Region,
			"sha3":      file.SHA3,
		}
		fileList[i] = fileMap
	}

	if err := d.Set("jcds2_enabled", props.JCDS2Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("file_stream_endpoint_enabled", props.FileStreamEndpointEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("max_chunk_size", props.MaxChunkSize); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("files", fileList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().UTC().String())

	return diags
}
