package computerinventorycollectionsettings

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceJamfProComputerInventoryCollectionSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Second),
			Read:   schema.DefaultTimeout(15 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"computer_inventory_collection_preferences": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"monitor_application_usage": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Collect the number of minutes applications are in the foreground.",
						},
						"include_fonts": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Include fonts in inventory collection.",
						},
						"include_plugins": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Include plugins in inventory collection.",
						},
						"include_packages": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Include packages in inventory collection.",
						},
						"include_software_updates": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Include software updates in inventory collection.",
						},
						"include_software_id": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Include software ID in inventory collection.",
						},
						"include_accounts": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Collect UIDs, usernames, full names, and home directory paths for local user accounts.",
						},
						"calculate_sizes": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Calculate sizes of items in inventory collection.",
						},
						"include_hidden_accounts": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Include hidden accounts in inventory collection.",
						},
						"include_printers": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Include printers in inventory collection.",
						},
						"include_services": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Include services in inventory collection.",
						},
						"collect_synced_mobile_device_info": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Collect information about synced mobile devices.",
						},
						"update_ldap_info_on_computer_inventory_submissions": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Update LDAP information when computer inventory is submitted.",
						},
						"monitor_beacons": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Monitor iBeacon regions and have computers submit information to Jamf Pro when they enter or exit a region.",
						},
						"allow_changing_user_and_location": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Collect user and location information from Directory Service.",
						},
						"use_unix_user_paths": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Allow local administrators to use the jamf binary recon verb to change User and Location inventory information in Jamf Pro.",
						},
						"collect_unmanaged_certificates": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Collect unmanaged certificates (managed certificates are always collected).",
						},
					},
				},
			},
			"application_paths": pathSchema("Custom paths to use for collecting applications. These paths will also be used for collecting Application Usage information (if enabled)"),
			"font_paths":        pathSchema("Custom paths to use when collecting fonts. Collect names, version numbers, and paths of installed fonts"),
			"plugin_paths":      pathSchema("Custom paths to use when collecting plug-ins. Collect names, version numbers, and paths of installed plug-ins"),
		},
	}
}

// pathSchema returns a schema definition for a path list with ID tracking
func pathSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Description: description,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"path": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Custom path to be included in inventory collection.",
				},
				"id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Unique ID for this inventory collection custom path.",
				},
			},
		},
	}
}
