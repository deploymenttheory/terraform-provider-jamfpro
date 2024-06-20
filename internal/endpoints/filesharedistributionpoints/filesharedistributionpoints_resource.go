// filesharedistributionpoints_resource.go
package filesharedistributionpoints

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProFileShareDistributionPoints defines the schema and CRUD operations for managing Jamf Pro Distribution Point in Terraform.
func ResourceJamfProFileShareDistributionPoints() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProFileShareDistributionPointsCreate,
		ReadContext:   resourceJamfProFileShareDistributionPointsRead,
		UpdateContext: resourceJamfProFileShareDistributionPointsUpdate,
		DeleteContext: resourceJamfProFileShareDistributionPointsDelete,
		CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(15 * time.Second),
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
					v := util.GetString(val)
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

const (
	JamfProResourceDistributionPoint = "Distribution Point"
)

// resourceJamfProFileShareDistributionPointsCreate is responsible for creating a new file share
// distribution point object in the remote system.
// The function:
// 1. Constructs the dock item data using the provided Terraform configuration.
// 2. Calls the API to create the dock item in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created dock item.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func resourceJamfProFileShareDistributionPointsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProFileShareDistributionPoint(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Fileshare Distribution Point: %v", err))
	}

	var creationResponse *jamfpro.ResponseFileShareDistributionPointCreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateDistributionPoint(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Fileshare Distribution Point '%s' after retries: %v", resource.Name, err))
	}

	d.SetId(strconv.Itoa(creationResponse.ID))

	// checkResourceExists := func(id interface{}) (interface{}, error) {
	// 	intID, err := strconv.Atoi(id.(string))
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
	// 	}
	// 	return client.GetDistributionPointByID(intID)
	// }

	// _, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Fileshare Distribution Point", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second)
	// if waitDiags.HasError() {
	// 	return waitDiags
	// }

	return append(diags, resourceJamfProFileShareDistributionPointsRead(ctx, d, meta)...)
}

// resourceJamfProFileShareDistributionPointsRead is responsible for reading the current state of a
// Jamf Pro File Share Distribution Point Resource from the remote system.
// The function:
// 1. Fetches the dock item's current state using its ID. If it fails then obtain dock item's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the dock item being deleted outside of Terraform, to keep the Terraform state synchronized.
func resourceJamfProFileShareDistributionPointsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	resourceID := d.Id()
	var diags diag.Diagnostics

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := client.GetDistributionPointByID(resourceIDInt)

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// resourceJamfProFileShareDistributionPointsUpdate is responsible for updating an existing Jamf Pro Site on the remote system.
func resourceJamfProFileShareDistributionPointsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	resource, err := constructJamfProFileShareDistributionPoint(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro file share distribution point for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateDistributionPointByID(resourceIDInt, resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro file share distribution point '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	return append(diags, resourceJamfProFileShareDistributionPointsRead(ctx, d, meta)...)
}

// resourceJamfProFileShareDistributionPointsDeleteis responsible for deleting a Jamf Pro file share distribution point from the remote system.
func resourceJamfProFileShareDistributionPointsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics
	resourceID := d.Id()

	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		apiErr := client.DeleteDistributionPointByID(resourceIDInt)
		if apiErr != nil {
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteDistributionPointByName(resourceName)
			if apiErrByName != nil {
				return retry.RetryableError(apiErrByName)
			}
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro file share distribution point '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	d.SetId("")

	return diags
}
