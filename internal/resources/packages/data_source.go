// packages_data_source.go
package packages

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceJamfProPackages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the package metadata.",
			},
			"package_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique name of the Jamf Pro package.This doesn't have to match the filename of the package.",
			},
			"filename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The package filename reference of the Jamf Pro package. This is used to associate the package metadata with the file uploaded to the Jamf Pro server.",
			},
			"package_file_source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The file path or the URL source of the Jamf Pro package to be uploaded. Supports HTTP/HTTPS URLs, and local filepaths.",
			},
			"category_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The category ID of the Jamf Pro package. Defaults to -1 if not specified.",
			},
			"info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Information to display to the administrator when the package is deployed or uninstalled.",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Notes to display about the package (e.g., who built it and when it was built)",
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The package priority to use for deploying or uninstalling the package (e.g., A package with a priority of '1' is deployed or uninstalled before other packages)",
			},
			"os_requirements": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OS requirements for the Jamf Pro package. The package can only be deployed to computers with these operating system versions. Each version must be separated by a comma (e.g., '10.6.8, 10.7.x, 10.8')",
			},
			"fill_user_template": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Define whether to fill new home directories with the contents of the home directory in the package's Users folder. Applies to DMGs only. This setting can be changed when deploying or uninstalling the package using a policy.",
			},
			"indexed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the package has completed indexing within the jamf content delivery service.",
			},
			"fill_existing_users": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to fill existing home directories with the contents of the home directory in the package's Users folder. Applies to DMGs only. This setting can be changed when deploying or uninstalling the package using a policy.",
			},
			"swu": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Install the package only if it is available as an update. For this to work, the display name of the package must match the name in the command-line version of Software Update. Applies to PKGs only",
			},
			"reboot_required": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Defines whether a computer must be restarted after installing the package",
			},
			"self_heal_notify": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to notify for self-heal.",
			},
			"self_healing_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The self-healing action for the package. Defaults to 'nothing' if not specified.",
			},
			"os_install": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the package is an OS install.",
			},
			"serial_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The serial number of the package.",
			},
			"parent_package_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The parent package ID. Defaults to -1 if not specified.",
			},
			"base_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The base path for the package.",
			},
			"suppress_updates": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to suppress updates.",
			},
			"cloud_transfer_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cloud transfer status.",
			},
			"ignore_conflicts": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to ignore conflicts.",
			},
			"suppress_from_dock": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to suppress from dock.",
			},
			"suppress_eula": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to suppress EULA.",
			},
			"suppress_registration": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to suppress registration.",
			},
			"install_language": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The install language of the package.",
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
				Computed:    true,
				Description: "The size of the package.",
			},
			"os_installer_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OS installer version.",
			},
			"manifest": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The manifest of the package.",
			},
			"manifest_file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The manifest file name.",
			},
			"format": {
				Type:        schema.TypeString,
				Computed:    true,
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
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)

	resourceID := d.Get("id").(string)
	name := d.Get("package_name").(string)

	if resourceID == "" && name == "" {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return diag.FromErr(fmt.Errorf("either 'id' or 'package_name' must be provided"))
	}

	var resource *jamfpro.ResourcePackage
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error

		if name != "" {
			resource, apiErr = client.GetPackageByName(name)
		} else {
			resource, apiErr = client.GetPackageByID(resourceID)
		}

		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		lookupMethod := "ID"
		lookupValue := resourceID
		if name != "" {
			lookupMethod = "name"
			lookupValue = name
		}
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Package with %s '%s' after retries: %v", lookupMethod, lookupValue, err))
	}

	if resource == nil {
		d.SetId("")
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return diag.FromErr(fmt.Errorf(" the Jamf Pro Package cannot be found"))
	}

	d.SetId(resource.ID)
	return updateState(d, resource)
}
