// packages_resource.go
package packages

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
			Create: schema.DefaultTimeout(20 * time.Minute),
			Read:   schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the package.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the Jamf Pro package.",
			},
			"package_file_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file path of the Jamf Pro package.",
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The category of the Jamf Pro package.",
			},
			"filename": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The filename of the Jamf Pro package.",
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
			"reboot_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether a reboot is required after installing the Jamf Pro package.",
			},
			"fill_user_template": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to fill the user template.",
			},
			"fill_existing_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to fill existing users.",
			},
			"boot_volume_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether a boot volume is required.",
				Default:     false,
			},
			"allow_uninstalled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to allow the package to be uninstalled.",
			},
			"os_requirements": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The OS requirements for the Jamf Pro package.",
			},
			"required_processor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The required processor for the Jamf Pro package.",
				Default:     "None",
			},
			"switch_with_package": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The package to switch with.",
				Default:     "Do Not Install",
			},
			"install_if_reported_available": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to install the package if it's reported as available.",
			},
			"reinstall_option": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The reinstall option for the Jamf Pro package.",
				Default:     "Do Not Reinstall",
			},
			"triggering_files": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The triggering files for the Jamf Pro package.",
			},
			"send_notification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to send a notification for the Jamf Pro package.",
			},
		},
	}
}

// ResourceJamfProPackagesCreate is responsible for creating a new Jamf Pro Package in the remote system.
// The function:
// 1. Constructs the attribute data using the provided Terraform configuration.
// 2. Calls the API to create the attribute in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created attribute.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProPackagesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize diagnostics
	var diags diag.Diagnostics

	// Extract the file path for the package
	filePath := d.Get("package_file_path").(string)

	// Step 1: Call CreateJCDS2PackageV2 to upload the file to JCDS 2.0
	fileUploadResponse, err := conn.CreateJCDS2PackageV2(filePath)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to upload file to JCDS 2.0 with file path '%s': %v", filePath, err))
	}
	fmt.Printf("File uploaded successfully, URI: %s\n", fileUploadResponse.URI)

	// Pause for 10 seconds
	time.Sleep(10 * time.Second)

	// Construct the resource object
	packageResourcePointer, err := constructJamfProPackageCreate(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package: %v", err))
	}

	// Dereference the pointer to get the value
	packageResource := *packageResourcePointer

	// Step 2: Call CreatePackage to create the package metadata in Jamf Pro
	creationResponse, err := conn.CreatePackage(packageResource)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Package '%s': %v", packageResource.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Read the site to ensure the Terraform state is up to date
	readDiags := ResourceJamfProPackagesRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProPackagesRead is responsible for reading the current state of a Jamf Pro Site Resource from the remote system.
// The function:
// 1. Fetches the attribute's current state using its ID. If it fails then obtain attribute's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the attribute being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProPackagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error converting resource ID '%s' to int: %v", resourceID, err))
	}

	var resource *jamfpro.ResourcePackage

	// Read operation with retry
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = conn.GetPackageByID(resourceIDInt)
		if apiErr != nil {
			if strings.Contains(apiErr.Error(), "404") || strings.Contains(apiErr.Error(), "410") {
				// Return non-retryable error with a message to avoid SDK issues
				return retry.NonRetryableError(fmt.Errorf("resource not found, marked for deletion"))
			}
			// Retry for other types of errors
			return retry.RetryableError(apiErr)
		}
		return nil
	})

	// If err is not nil, check if it's due to the resource being not found
	if err != nil {
		if err.Error() == "resource not found, marked for deletion" {
			// Resource not found, remove from Terraform state
			d.SetId("")
			// Append a warning diagnostic and return
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("Jamf Pro Package with ID '%s' was not found on the server and is marked for deletion from terraform state.", resourceID),
			})
			return diags
		}

		// For other errors, return an error diagnostic
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Package with ID '%s' after retries: %v", resourceID, err))
	}

	// Update Terraform state with the resource information
	if err := d.Set("id", strconv.Itoa(resource.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("name", resource.Name); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("category", resource.Category); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("filename", resource.Filename); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("info", resource.Info); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("notes", resource.Notes); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("priority", resource.Priority); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("reboot_required", resource.RebootRequired); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("fill_user_template", resource.FillUserTemplate); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("fill_existing_users", resource.FillExistingUsers); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("boot_volume_required", resource.BootVolumeRequired); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("allow_uninstalled", resource.AllowUninstalled); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("os_requirements", resource.OSRequirements); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("required_processor", resource.RequiredProcessor); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("switch_with_package", resource.SwitchWithPackage); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("install_if_reported_available", resource.InstallIfReportedAvailable); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("reinstall_option", resource.ReinstallOption); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("triggering_files", resource.TriggeringFiles); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("send_notification", resource.SendNotification); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

// ResourceJamfProPackagesUpdate is responsible for updating an existing Jamf Pro Site on the remote system.
func ResourceJamfProPackagesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize diagnostics
	var diags diag.Diagnostics

	// Construct the package resource object
	packageData, err := constructJamfProPackage(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro Package: %v", err))
	}

	// Extract the file path for the package
	filePath := d.Get("package_file_path").(string)

	// Call DoPackageUpload to upload the package and create its metadata
	_, resource, err := conn.DoPackageUpload(filePath, packageData)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro Package with file path '%s': %v", filePath, err))
	}

	// Assuming the ID from packageCreationResponse is a suitable unique identifier for the Terraform resource
	if resource != nil {
		d.SetId(strconv.Itoa(resource.ID))
	} else {
		return diag.FromErr(fmt.Errorf("package creation response is nil"))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProPackagesRead(ctx, d, meta)
	if readDiags.HasError() {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProPackagesDelete is responsible for deleting a Jamf Pro Package.
func ResourceJamfProPackagesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize API client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

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
		apiErr := conn.DeletePackageByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeletePackageByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro Package '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
