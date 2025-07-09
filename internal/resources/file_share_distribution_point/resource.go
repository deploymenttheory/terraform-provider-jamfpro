// filesharedistributionpoints_resource.go
package file_share_distribution_point

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProFileShareDistributionPoints defines the schema and CRUD operations for managing Jamf Pro Distribution Point in Terraform.
func ResourceJamfProFileShareDistributionPoints() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(70 * time.Second),
			Update: schema.DefaultTimeout(70 * time.Second),
			Delete: schema.DefaultTimeout(70 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the distribution point.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the distribution point.",
			},
			"serverName": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hostname of the distribution point server.",
			},
			"principal": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if the distribution point is the principal distribution point, used as the authoritative source for all files. Defaults to false.",
			},
			"backupDistributionPointID": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "-1",
				Description: "The ID of the failover point. Defaults to -1.",
			},
			"localPathToShare": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The local path to the share.",
			},
			"fileSharingConnectionType": {
				Type:        schema.TypeString,
				Required:    true,
				Default:     "NONE",
				Description: "The type of connection protocol to the distribution point. Can be either 'SMB', 'AFP', or 'NONE'. Required. Defaults to 'NONE'. Either this or httpsEnabled must be set.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validTypes := map[string]bool{
						"SMB":  true,
						"AFP":  true,
						"NONE": true,
					}
					if _, valid := validTypes[v]; !valid {
						errs = append(errs, fmt.Errorf("%q must be one of 'SMB', 'AFP', or 'NONE', got: %s", key, v))
					}
					return warns, errs
				},
			},
			"shareName": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the network share. Required if fileSharingConnectionType is either AFP or SMB.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     139,
				Description: "The port number used for the fileshare distribution point. Defaults to 139. Required if fileSharingConnectionType is either AFP or SMB.",
			},
			"enableLoadBalancing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if load balancing is enabled. Defaults to false. Cannot be enabled when the backup distribution point configured is cloud.",
			},
			"sshUsername": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SSH username for the distribution point.",
			},
			"sshPassword": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The password for the distribution point. This field is marked as sensitive and will not be displayed in logs or console output.",
			},

			"workgroup": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The workgroup or domain of the distribution point.",
			},
			"readOnlyUsername": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username for read-only access to the distribution point.",
			},
			"readOnlyPassword": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password for read-only access. This field is marked as sensitive.",
			},
			"readWriteUsername": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username for read-write access to the distribution point.",
			},
			"readWritePassword": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password for read-write access. This field is marked as sensitive.",
			},
			"httpsEnabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if HTTP downloads are enabled. Defaults to false. Allow downloads over HTTPS - requires installation of a valid SSL certificate.",
			},
			"httpsPort": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     443,
				Description: "The port number for the https share. Defaults to 443. Required if HTTPS enabled.",
			},
			"httpsContext": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to the https share (e.g. if the share is accessible at http://192.168.10.10/JamfShare, the context is 'JamfShare'). Required if HTTPS enabled.",
			},
			"httpsSecurityType": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NONE",
				Description: "Type of authentication required to download files from the distribution point. Can be 'USERNAME_PASSWORD' or 'NONE'. Defaults to 'NONE'. Required if HTTPS enabled.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validTypes := map[string]bool{
						"NONE":              true,
						"USERNAME_PASSWORD": true,
					}
					if _, valid := validTypes[v]; !valid {
						errs = append(errs, fmt.Errorf("%q must be one of 'USERNAME_PASSWORD' or 'NONE', got: %s", key, v))
					}
					return warns, errs
				},
			},
			"httpsUsername": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username for HTTP access, if username/password authentication is required. Required if httpsSecurityType is USERNAME_PASSWORD.",
			},
			"httpsPassword": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password for HTTP access, if username/password authentication is required. This field is marked as sensitive. Required if httpsSecurityType is USERNAME_PASSWORD.",
			},
		},
	}
}
