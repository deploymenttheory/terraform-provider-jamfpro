// packages_resource.go
package packages

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProPackages defines the schema and CRUD operations for managing Jamf Pro Packages in Terraform.
func ResourceJamfProPackages() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProPackagesCreate,
		ReadContext:   ResourceJamfProPackagesRead,
		UpdateContext: ResourceJamfProPackagesUpdate,
		DeleteContext: ResourceJamfProPackagesDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		CustomizeDiff: customValidateFilePath,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the package.",
			},
			"package_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro package.",
			},
			"package_file_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file path of the Jamf Pro package.",
			},
			"category_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The category ID of the Jamf Pro package.",
				Default:     "-1",
			},
			"info": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Information about the Jamf Pro package.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notes associated with the Jamf Pro package.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The priority of the Jamf Pro package.",
			},
			"os_requirements": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OS requirements for the Jamf Pro package.",
			},
			"fill_user_template": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to fill the user template.",
			},
			"indexed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the package is indexed.",
			},
			"fill_existing_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to fill existing users.",
			},
			"swu": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the package is a software update.",
			},
			"reboot_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether a reboot is required after installing the Jamf Pro package.",
			},
			"self_heal_notify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to notify for self-heal.",
			},
			"self_healing_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The self-healing action for the package.",
			},
			"os_install": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the package is an OS install.",
			},
			"serial_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The serial number of the package.",
			},
			"parent_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The parent package ID.",
			},
			"base_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The base path for the package.",
			},
			"suppress_updates": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to suppress updates.",
			},
			"cloud_transfer_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cloud transfer status.",
			},
			"ignore_conflicts": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to ignore conflicts.",
			},
			"suppress_from_dock": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to suppress from dock.",
			},
			"suppress_eula": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to suppress EULA.",
			},
			"suppress_registration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to suppress registration.",
			},
			"install_language": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The install language.",
			},
			"md5": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The MD5 hash of the package.",
			},
			"sha256": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SHA256 hash of the package.",
			},
			"hash_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hash type of the package.",
			},
			"hash_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hash value of the package.",
			},
			"size": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The size of the package.",
			},
			"os_installer_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OS installer version.",
			},
			"manifest": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The manifest of the package.",
			},
			"manifest_file_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The manifest file name.",
			},
			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The format of the package.",
			},
			"package_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URI of the package in the Jamf Cloud Distribution Service (JCDS).",
			},
			"md5_file_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "md5 hash of the package file for integrity comparison.",
			},
			"filename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The filename of the Jamf Pro package.",
			},
		},
	}
}
