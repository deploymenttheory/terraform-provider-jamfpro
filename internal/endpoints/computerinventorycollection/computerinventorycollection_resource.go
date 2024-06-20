// computerinventorycollection_resource.go
package computerinventorycollection

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProComputerInventoryCollection defines the schema and RU operations for managing Jamf Pro computer checkin configuration in Terraform.
func ResourceJamfProComputerInventoryCollection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProComputerInventoryCollectionCreate,
		ReadContext:   resourceJamfProComputerInventoryCollectionRead,
		UpdateContext: resourceJamfProComputerInventoryCollectionUpdate,
		DeleteContext: resourceJamfProComputerInventoryCollectionDelete,
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
			"local_user_accounts": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on local user accounts.",
			},
			"home_directory_sizes": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on home directory sizes.",
			},
			"hidden_accounts": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on hidden accounts.",
			},
			"printers": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on printers.",
			},
			"active_services": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on active services.",
			},
			"mobile_device_app_purchasing_info": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on mobile device app purchasing.",
			},
			"computer_location_information": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect computer location information.",
			},
			"package_receipts": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect package receipts.",
			},
			"available_software_updates": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Collect information on available software updates.",
			},
			"include_applications": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Include applications in the inventory collection.",
			},
			"include_fonts": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Include fonts in the inventory collection.",
			},
			"include_plugins": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Include plugins in the inventory collection.",
			},
			"applications": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of applications.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path to the application.",
						},
						"platform": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Platform of the application.",
						},
					},
				},
			},
			"fonts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of fonts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path to the font.",
						},
						"platform": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Platform of the font.",
						},
					},
				},
			},
			"plugins": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of plugins.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path to the plugin.",
						},
						"platform": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Platform of the plugin.",
						},
					},
				},
			},
		},
	}
}

// resourceJamfProComputerInventoryCollectionCreate is responsible for initializing the Jamf Pro Computer Inventory Collection configuration in Terraform.
// Since this resource is a configuration set and not a resource that is 'created' in the traditional sense,
// this function will simply set the initial state in Terraform.
// resourceJamfProComputerInventoryCollectionCreate is responsible for initializing the Jamf Pro Computer Inventory Collection configuration in Terraform.
func resourceJamfProComputerInventoryCollectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := constructJamfProComputerInventoryCollection(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Computer Inventory Collection for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		apiErr := client.UpdateComputerInventoryCollectionInformation(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Computer Inventory Collection configuration after retries: %v", err))
	}

	d.SetId("jamfpro_computer_inventory_collection_singleton")

	return append(diags, resourceJamfProComputerInventoryCollectionRead(ctx, d, meta)...)
}

// resourceJamfProComputerInventoryCollectionRead is responsible for reading the current state of the Jamf Pro Computer Inventory Collection configuration.
func resourceJamfProComputerInventoryCollectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	resource, err := client.GetComputerInventoryCollectionInformation()

	d.SetId("jamfpro_computer_inventory_collection_singleton")

	if err != nil {
		return state.HandleResourceNotFoundError(err, d)
	}

	return append(diags, updateTerraformState(d, resource)...)
}

// resourceJamfProComputerInventoryCollectionUpdate is responsible for updating the Jamf Pro Computer Inventory Collection configuration.
func resourceJamfProComputerInventoryCollectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*jamfpro.Client)
	var diags diag.Diagnostics

	inventoryCollectionConfig, err := constructJamfProComputerInventoryCollection(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Computer Inventory Collection for update: %v", err))
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		apiErr := client.UpdateComputerInventoryCollectionInformation(inventoryCollectionConfig)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Computer Inventory Collection configuration after retries: %v", err))
	}

	d.SetId("jamfpro_computer_checkin_singleton")

	return append(diags, resourceJamfProComputerInventoryCollectionRead(ctx, d, meta)...)
}

// resourceJamfProComputerInventoryCollectionDelete is responsible for 'deleting' the Jamf Pro Computer Inventory Collection configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func resourceJamfProComputerInventoryCollectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}
