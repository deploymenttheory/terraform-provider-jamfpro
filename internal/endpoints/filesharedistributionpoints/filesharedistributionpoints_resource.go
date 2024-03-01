// filesharedistributionpoints_resource.go
package filesharedistributionpoints

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/client"
	util "github.com/deploymenttheory/terraform-provider-jamfpro/internal/helpers/type_assertion"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceJamfProFileShareDistributionPoints defines the schema and CRUD operations for managing Jamf Pro Distribution Point in Terraform.
func ResourceJamfProFileShareDistributionPoints() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceJamfProFileShareDistributionPointsCreate,
		ReadContext:   ResourceJamfProFileShareDistributionPointsRead,
		UpdateContext: ResourceJamfProFileShareDistributionPointsUpdate,
		DeleteContext: ResourceJamfProFileShareDistributionPointsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Read:   schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
			Delete: schema.DefaultTimeout(30 * time.Second),
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

// ResourceJamfProFileShareDistributionPointsCreate is responsible for creating a new file share
// distribution point object in the remote system.
// The function:
// 1. Constructs the dock item data using the provided Terraform configuration.
// 2. Calls the API to create the dock item in Jamf Pro.
// 3. Updates the Terraform state with the ID of the newly created dock item.
// 4. Initiates a read operation to synchronize the Terraform state with the actual state in Jamf Pro.
func ResourceJamfProFileShareDistributionPointsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Assert the meta interface to the expected APIClient type
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize variables
	var diags diag.Diagnostics

	// Construct the resource object
	resource, err := constructJamfProFileShareDistributionPoint(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro file share distribution point: %v", err))
	}

	// Retry the API call to create the resource in Jamf Pro
	var creationResponse *jamfpro.ResponseFileShareDistributionPointCreatedAndUpdated
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateDistributionPoint(resource)
		if apiErr != nil {
			return retry.RetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Jamf Pro file share distribution point '%s' after retries: %v", resource.Name, err))
	}

	// Set the resource ID in Terraform state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Read the site to ensure the Terraform state is up to date
	readDiags := ResourceJamfProFileShareDistributionPointsRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProFileShareDistributionPointsRead is responsible for reading the current state of a
// Jamf Pro File Share Distribution Point Resource from the remote system.
// The function:
// 1. Fetches the dock item's current state using its ID. If it fails then obtain dock item's current state using its Name.
// 2. Updates the Terraform state with the fetched data to ensure it accurately reflects the current state in Jamf Pro.
// 3. Handles any discrepancies, such as the dock item being deleted outside of Terraform, to keep the Terraform state synchronized.
func ResourceJamfProFileShareDistributionPointsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	var resource *jamfpro.ResourceFileShareDistributionPoint

	// Read operation with retry
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		resource, apiErr = conn.GetDistributionPointByID(resourceIDInt)
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
				Detail:   fmt.Sprintf("Jamf Pro Distribution Point with ID '%s' was not found on the server and is marked for deletion from terraform state.", resourceID),
			})
			return diags
		}

		// For other errors, return an error diagnostic
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Distribution Point with ID '%s' after retries: %v", resourceID, err))
	}

	// Check if fileShareDistributionPoint data exists
	if resource != nil {
		// Set the fields directly in the Terraform state
		if err := d.Set("id", strconv.Itoa(resource.ID)); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("name", resource.Name); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("ip_address", resource.IP_Address); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("ipaddress", resource.IPAddress); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("is_master", resource.IsMaster); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("is_master", resource.IsMaster); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("failover_point", resource.FailoverPoint); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("failover_point_url", resource.FailoverPointURL); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("enable_load_balancing", resource.EnableLoadBalancing); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("local_path", resource.LocalPath); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("ssh_username", resource.SSHUsername); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("password", resource.Password); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("connection_type", resource.ConnectionType); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("share_name", resource.ShareName); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("workgroup_or_domain", resource.WorkgroupOrDomain); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("share_port", resource.SharePort); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("read_only_username", resource.ReadOnlyUsername); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("https_downloads_enabled", resource.HTTPDownloadsEnabled); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("http_url", resource.HTTPURL); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("https_share_path", resource.Context); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("protocol", resource.Protocol); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("https_port", resource.Port); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("no_authentication_required", resource.NoAuthenticationRequired); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("https_username_password_required", resource.UsernamePasswordRequired); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("https_username", resource.HTTPUsername); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}

		// sensitive field handling
		// "read_only_password" , "read_write_password" and "https_password" are not stored in state
		// as they are sensitive fields and are not returned by the API.

	}

	return diags
}

// ResourceJamfProFileShareDistributionPointsUpdate is responsible for updating an existing Jamf Pro Site on the remote system.
func ResourceJamfProFileShareDistributionPointsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	// Construct the resource object
	resource, err := constructJamfProFileShareDistributionPoint(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to construct Jamf Pro file share distribution point for update: %v", err))
	}

	// Update operations with retries
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateDistributionPointByID(resourceIDInt, resource)
		if apiErr != nil {
			// If updating by ID fails, attempt to update by Name
			return retry.RetryableError(apiErr)
		}
		// Successfully updated the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update Jamf Pro file share distribution point '%s' (ID: %d) after retries: %v", resource.Name, resourceIDInt, err))
	}

	// Read the resource to ensure the Terraform state is up to date
	readDiags := ResourceJamfProFileShareDistributionPointsRead(ctx, d, meta)
	if len(readDiags) > 0 {
		diags = append(diags, readDiags...)
	}

	return diags
}

// ResourceJamfProFileShareDistributionPointsDeleteis responsible for deleting a Jamf Pro file share distribution point from the remote system.
func ResourceJamfProFileShareDistributionPointsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		apiErr := conn.DeleteDistributionPointByID(resourceIDInt)
		if apiErr != nil {
			// If deleting by ID fails, attempt to delete by Name
			resourceName := d.Get("name").(string)
			apiErrByName := conn.DeleteDistributionPointByName(resourceName)
			if apiErrByName != nil {
				// If deletion by name also fails, return a retryable error
				return retry.RetryableError(apiErrByName)
			}
		}
		// Successfully deleted the resource, exit the retry loop
		return nil
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Jamf Pro file share distribution point '%s' (ID: %d) after retries: %v", d.Get("name").(string), resourceIDInt, err))
	}

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return diags
}
