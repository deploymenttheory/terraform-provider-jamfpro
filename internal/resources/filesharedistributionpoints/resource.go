// filesharedistributionpoints_resource.go
package filesharedistributionpoints

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
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hostname or IP address of the distribution point server.",
			},
			"is_master": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the distribution point is the principal distribution point, used  as the authoritative source for all files",
			},
			"failover_point": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The failover point for the distribution point.Can be ",
			},
			"ipaddress": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hostname or IP address of the distribution point server.",
			},
			"failover_point_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the failover point.",
			},
			// Page 2
			"connection_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of connection protocol to the distribution point. Can be either 'AFP', or 'SMB'.",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					validTypes := map[string]bool{
						"SMB": true,
						"AFP": true,
					}
					if _, valid := validTypes[v]; !valid {

						errs = append(errs, fmt.Errorf("%q must be one of 'SMB', or 'AFP', got: %s", key, v))
					}
					return warns, errs
				},
			},
			"share_name": {
				Type:     schema.TypeString,
				Optional: true,

				Description: "The name of the network share.",
			},
			"share_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The port number used for the fileshare distribution point.",
			},
			"enable_load_balancing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if load balancing is enabled.",
			},
			"local_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The local path to the distribution point.",
			},

			"ssh_username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SSH username for the distribution point.",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The password for the distribution point. This field is marked as sensitive and will not be displayed in logs or console output.",
			},

			"workgroup_or_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The workgroup or domain of the distribution point.",
			},
			"read_only_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username for read-only access to the distribution point.",
			},
			"read_only_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password for read-only access. This field is marked as sensitive.",
			},
			"read_write_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username for read-write access to the distribution point.",
			},
			"read_write_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password for read-write access. This field is marked as sensitive.",
			},
			"no_authentication_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if no authentication is required for accessing the distribution point.",
			},
			// Page 3
			"https_downloads_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if HTTP downloads are enabled.",
			},
			"https_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The port number for the https share.",
			},
			"https_share_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to the https share (e.g. if the share is accessible at http://192.168.10.10/JamfShare, the context is 'JamfShare').",
			},
			"https_username_password_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if username/password authentication is required for accessing the distribution point.",
			},
			"https_username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username for HTTP access, if username/password authentication is required.",
			},
			"https_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The password for HTTP access, if username/password authentication is required. This field is marked as sensitive.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol used if HTTPS is enabled for the  distribution point. Result will always be 'https' if enabled.",
			},
			"http_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for HTTP downloads.Constructed from the protocol, IP address, and port.",
			},
		},
	}
}
