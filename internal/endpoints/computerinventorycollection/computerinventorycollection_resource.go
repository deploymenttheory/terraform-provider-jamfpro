// computerinventorycollection_resource.go
package computerinventorycollection

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProComputerInventoryCollection defines the schema and RU operations for managing Jamf Pro computer checkin configuration in Terraform.
func ResourceJamfProComputerInventoryCollection() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProComputerInventoryCollectionCreate,
		ReadContext:   ResourceJamfProComputerInventoryCollectionRead,
		UpdateContext: ResourceJamfProComputerInventoryCollectionUpdate,
		DeleteContext: ResourceJamfProComputerInventoryCollectionDelete,
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
							Optional:    true,
							Description: "Path to the application.",
						},
						"platform": {
							Type:        schema.TypeString,
							Optional:    true,
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
							Optional:    true,
							Description: "Path to the font.",
						},
						"platform": {
							Type:        schema.TypeString,
							Optional:    true,
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
							Optional:    true,
							Description: "Path to the plugin.",
						},
						"platform": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Platform of the plugin.",
						},
					},
				},
			},
		},
	}
}

// ResourceJamfProComputerInventoryCollectionCreate is responsible for initializing the Jamf Pro Computer Inventory Collection configuration in Terraform.
// Since this resource is a configuration set and not a resource that is 'created' in the traditional sense,
// this function will simply set the initial state in Terraform.
// ResourceJamfProComputerInventoryCollectionCreate is responsible for initializing the Jamf Pro Computer Inventory Collection configuration in Terraform.
func ResourceJamfProComputerInventoryCollectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProComputerInventoryCollection(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Computer Inventory Collection for update: %v", err))
	}

	// Update (or effectively create) the check-in configuration with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		apiErr := conn.UpdateComputerInventoryCollectionInformation(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Computer Inventory Collection configuration after retries: %v", err))
	}

	// Since this resource is a singleton, use a fixed ID to represent it in the Terraform state
	d.SetId("jamfpro_computer_inventory_collection_singleton")

	// Read the site to ensure the Terraform state is up to date
	readDiags := ResourceJamfProComputerInventoryCollectionRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProComputerInventoryCollectionRead is responsible for reading the current state of the Jamf Pro Computer Inventory Collection configuration.
func ResourceJamfProComputerInventoryCollectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}

	// Initialize variables
	var diags diag.Diagnostics

	// Attempt to fetch the resource by ID
	resource, err := apiclient.Conn.GetComputerInventoryCollectionInformation()

	// The constant ID "jamfpro_computer_inventory_collection_singleton" is assigned to satisfy Terraform's requirement for an ID.
	d.SetId("jamfpro_computer_inventory_collection_singleton")

	if err != nil {
		// Handle not found error or other errors
		return state.HandleResourceNotFoundError(err, d)
	}

	// Update the Terraform state with the fetched data from the resource
	diags = updateTerraformState(d, resource)

	// Handle any errors and return diagnostics
	if len(diags) > 0 {
		return diags
	}
	return nil
}

// ResourceJamfProComputerInventoryCollectionUpdate is responsible for updating the Jamf Pro Computer Inventory Collection configuration.
func ResourceJamfProComputerInventoryCollectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	inventoryCollectionConfig, err := constructJamfProComputerInventoryCollection(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Computer Inventory Collection for update: %v", err))
	}

	// Update (or effectively create) the check-in configuration with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		apiErr := conn.UpdateComputerInventoryCollectionInformation(inventoryCollectionConfig)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to apply Jamf Pro Computer Inventory Collection configuration after retries: %v", err))
	}

	// Since this resource is a singleton, use a fixed ID to represent it in the Terraform state
	d.SetId("jamfpro_computer_checkin_singleton")

	// Read the site to ensure the Terraform state is up to date
	readDiags := ResourceJamfProComputerInventoryCollectionRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProComputerInventoryCollectionDelete is responsible for 'deleting' the Jamf Pro Computer Inventory Collection configuration.
// Since this resource represents a configuration and not an actual entity that can be deleted,
// this function will simply remove it from the Terraform state.
func ResourceJamfProComputerInventoryCollectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Simply remove the resource from the Terraform state by setting the ID to an empty string.
	d.SetId("")

	return nil
}
