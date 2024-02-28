// filesharedistributionpoints_resource.go
package filesharedistributionpoints

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
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

// cconstructJamfProFileShareDistributionPoint constructs a ResourceDockItem object from the provided schema data.
func constructJamfProFileShareDistributionPoint(ctx context.Context, d *schema.ResourceData) (*jamfpro.ResourceFileShareDistributionPoint, error) {
	fileShareDistributionPoint := &jamfpro.ResourceFileShareDistributionPoint{
		Name:                     d.Get("name").(string),
		IP_Address:               d.Get("ip_address").(string),
		IsMaster:                 d.Get("is_master").(bool),
		FailoverPoint:            d.Get("failover_point").(string),
		ConnectionType:           d.Get("connection_type").(string),
		ShareName:                d.Get("share_name").(string),
		SharePort:                d.Get("share_port").(int),
		EnableLoadBalancing:      d.Get("enable_load_balancing").(bool),
		WorkgroupOrDomain:        d.Get("workgroup_or_domain").(string),
		ReadOnlyUsername:         d.Get("read_only_username").(string),
		ReadOnlyPassword:         d.Get("read_only_password").(string),
		ReadWriteUsername:        d.Get("read_write_username").(string),
		ReadWritePassword:        d.Get("read_write_password").(string),
		NoAuthenticationRequired: d.Get("no_authentication_required").(bool),
		HTTPDownloadsEnabled:     d.Get("https_downloads_enabled").(bool),
		Port:                     d.Get("https_port").(int),
		Context:                  d.Get("https_share_path").(string),
		HTTPUsername:             d.Get("https_username").(string),
		HTTPPassword:             d.Get("https_password").(string),
		Protocol:                 d.Get("protocol").(string),
		HTTPURL:                  d.Get("http_url").(string),
	}

	// Serialize and pretty-print the file share distribution point object as XML
	resourceXML, err := xml.MarshalIndent(fileShareDistributionPoint, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro File Share Distribution Point '%s' to XML: %v", fileShareDistributionPoint.Name, err)
	}
	fmt.Printf("Constructed Jamf Pro File Share Distribution Point XML:\n%s\n", string(resourceXML))

	return fileShareDistributionPoint, nil
}

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
	resource, err := constructJamfProFileShareDistributionPoint(ctx, d)
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

	var fileShareDistributionPoint *jamfpro.ResourceFileShareDistributionPoint

	// Read operation with retry
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		fileShareDistributionPoint, apiErr = conn.GetDistributionPointByID(resourceIDInt)
		if apiErr != nil {
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the site, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		d.SetId("") // Remove from Terraform state if unable to read after retries
		return diag.FromErr(fmt.Errorf("failed to read Jamf Pro Printer with ID '%d' after retries: %v", resourceIDInt, err))
	}

	// Check if fileShareDistributionPoint data exists
	if fileShareDistributionPoint != nil {
		// Set the fields directly in the Terraform state
		if err := d.Set("id", strconv.Itoa(fileShareDistributionPoint.ID)); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("name", fileShareDistributionPoint.Name); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("ip_address", fileShareDistributionPoint.IP_Address); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("ipaddress", fileShareDistributionPoint.IPAddress); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("is_master", fileShareDistributionPoint.IsMaster); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("is_master", fileShareDistributionPoint.IsMaster); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("failover_point", fileShareDistributionPoint.FailoverPoint); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("failover_point_url", fileShareDistributionPoint.FailoverPointURL); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("enable_load_balancing", fileShareDistributionPoint.EnableLoadBalancing); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("local_path", fileShareDistributionPoint.LocalPath); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("ssh_username", fileShareDistributionPoint.SSHUsername); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("password", fileShareDistributionPoint.Password); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("connection_type", fileShareDistributionPoint.ConnectionType); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("share_name", fileShareDistributionPoint.ShareName); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("workgroup_or_domain", fileShareDistributionPoint.WorkgroupOrDomain); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("share_port", fileShareDistributionPoint.SharePort); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("read_only_username", fileShareDistributionPoint.ReadOnlyUsername); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("https_downloads_enabled", fileShareDistributionPoint.HTTPDownloadsEnabled); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("http_url", fileShareDistributionPoint.HTTPURL); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("https_share_path", fileShareDistributionPoint.Context); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("protocol", fileShareDistributionPoint.Protocol); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("https_port", fileShareDistributionPoint.Port); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("no_authentication_required", fileShareDistributionPoint.NoAuthenticationRequired); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("https_username_password_required", fileShareDistributionPoint.UsernamePasswordRequired); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
		if err := d.Set("https_username", fileShareDistributionPoint.HTTPUsername); err != nil {
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
	resource, err := constructJamfProFileShareDistributionPoint(ctx, d)
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
