// restrictedsoftware_resource.go
package restrictedsoftware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/endpoints/common/state"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/waitfor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProRestrictedSoftwares defines the schema and CRUD operations for managing Jamf Pro Restricted Software in Terraform.
func ResourceJamfProRestrictedSoftwares() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProRestrictedSoftwareCreate,
		ReadContext:   ResourceJamfProRestrictedSoftwareRead,
		UpdateContext: ResourceJamfProRestrictedSoftwareUpdate,
		DeleteContext: ResourceJamfProRestrictedSoftwareDelete,
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
				Description: "The unique identifier of the restricted software.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the restricted software.",
			},
			"process_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The process name of the restricted software.",
			},
			"match_exact_process_name": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the process name should be matched exactly.",
			},
			"send_notification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if a notification should be sent.",
			},
			"kill_process": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the process should be killed.",
			},
			"delete_executable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the executable should be deleted.",
			},
			"display_message": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The message to display when the software is restricted.",
			},
			"site": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The unique identifier of the site.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name of the site.",
						},
					},
				},
				Description: "The site associated with the restricted software.",
			},
			"scope": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The scope of the restricted software.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_computers": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if the restricted software applies to all computers.",
						},
						"computer_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "A list of computer IDs associated with the restricted software.",
						},
						"computer_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "A list of computer group IDs associated with the restricted software.",
						},
						"building_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "A list of building IDs associated with the restricted software.",
						},
						"department_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "A list of department IDs associated with the restricted software.",
						},
						"limitations": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Limitations for the restricted software.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_segment_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "A list of network segment IDs for limitations.",
									},
									"ibeacon_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "A list of iBeacon IDs for limitations.",
									},
								},
							},
						},
						"exclusions": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Exclusions for the restricted software.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"computer_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "A list of computer IDs for exclusions.",
									},
									"computer_group_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "A list of computer group IDs for exclusions.",
									},
									"building_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "A list of building IDs for exclusions.",
									},
									"department_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "A list of department IDs for exclusions.",
									},
									"directory_service_or_local_usernames": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "A list of directory service / local usernames for scoping exclusions.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// scopeEntitySchema returns the schema for scope entities.
func scopeEntitySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The unique identifier of the scope entity.",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the scope entity.",
		},
	}
}

// ResourceJamfProRestrictedSoftwareCreate is responsible for creating a new Jamf Pro Restricted Software in the remote system.
// The function:
// 1. Constructs the User Group data using the provided Terraform configuration.
// 2. Calls the API to create the User Group in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created User Group.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProRestrictedSoftwareCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected client type
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProRestrictedSoftware(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Restricted Software: %v", err))
	}

	// Retry the API call to create the resource in Jamf Pro
	var creationResponse *jamfpro.ResponseRestrictedSoftwareCreateAndUpdate
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = client.CreateRestrictedSoftware(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Restricted Software '%s' after retries: %v", resource.General.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Wait for the resource to be fully available before reading it
	checkResourceExists := func(id interface{}) (interface{}, error) {
		intID, err := strconv.Atoi(id.(string))
		if err != nil {
			return nil, fmt.Errorf("error converting ID '%v' to integer: %v", id, err)
		}
		return client.GetRestrictedSoftwareByID(intID)
	}

	_, waitDiags := waitfor.ResourceIsAvailable(ctx, d, "Jamf Pro Restricted Software", strconv.Itoa(creationResponse.ID), checkResourceExists, time.Duration(common.DefaultPropagationTime)*time.Second, client.EnableCookieJar)
	if waitDiags.HasError() {
		return waitDiags
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProRestrictedSoftwareRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProRestrictedSoftwareRead is responsible for reading the current state of a Jamf Pro Restricted Software Resource from the remote system.
// The function:
// 1. Fetches the user group's current state using its ID. If it fails, it tries to obtain the user group's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the user group being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProRestrictedSoftwareRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Attempt to fetch the resource by ID
	resource, err := client.GetRestrictedSoftwareByID(resourceIDInt)

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

// ResourceJamfProRestrictedSoftwareUpdate is responsible for updating an existing Jamf Pro Printer on the remote system.
func ResourceJamfProRestrictedSoftwareUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Construct the resource object
	resource, err := constructJamfProRestrictedSoftware(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Restricted Software for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := client.UpdateRestrictedSoftwareByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro Restricted Software '%s' (ID: %d) after retries: %v", resource.General.Name, resourceIDInt, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProRestrictedSoftwareRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProRestrictedSoftwareDelete is responsible for deleting a Jamf Pro Restricted Software.
func ResourceJamfProRestrictedSoftwareDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	client, ok := meta.(*jamfpro.Client)
	if !ok {
		return diag.Errorf("error asserting meta as *client.client")
	}

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Attempt to delete by ID
		apiErr := client.DeleteRestrictedSoftwareByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := client.DeleteRestrictedSoftwareByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Restricted Software '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
