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
	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/logging"

	"github.com/hashicorp/go-hclog"
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
	// Initialize the logging subsystem for the construction operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemConstruct, hclog.Debug)

	fileShareDistributionPoint := &jamfpro.ResourceFileShareDistributionPoint{
		Name:             util.GetStringFromInterface(d.Get("name")),
		IsMaster:         util.GetBoolFromInterface(d.Get("is_master")),
		IP_Address:       util.GetStringFromInterface(d.Get("ip_address")),
		IPAddress:        util.GetStringFromInterface(d.Get("ipaddress")),
		FailoverPoint:    util.GetStringFromInterface(d.Get("failover_point")),
		FailoverPointURL: util.GetStringFromInterface(d.Get("failover_point_url")),
		ConnectionType:   util.GetStringFromInterface(d.Get("connection_type")),
		ShareName:        util.GetStringFromInterface(d.Get("share_name")),
		SharePort:        util.GetIntFromInterface(d.Get("share_port")),

		EnableLoadBalancing:      util.GetBoolFromInterface(d.Get("enable_load_balancing")),
		LocalPath:                util.GetStringFromInterface(d.Get("local_path")),
		SSHUsername:              util.GetStringFromInterface(d.Get("ssh_username")),
		Password:                 util.GetStringFromInterface(d.Get("password")),
		WorkgroupOrDomain:        util.GetStringFromInterface(d.Get("workgroup_or_domain")),
		ReadOnlyUsername:         util.GetStringFromInterface(d.Get("read_only_username")),
		ReadOnlyPassword:         util.GetStringFromInterface(d.Get("read_only_password")),
		ReadWriteUsername:        util.GetStringFromInterface(d.Get("read_write_username")),
		ReadWritePassword:        util.GetStringFromInterface(d.Get("read_write_password")),
		HTTPDownloadsEnabled:     util.GetBoolFromInterface(d.Get("https_downloads_enabled")),
		HTTPURL:                  util.GetStringFromInterface(d.Get("http_url")),
		Context:                  util.GetStringFromInterface(d.Get("https_share_path")),
		Protocol:                 util.GetStringFromInterface(d.Get("protocol")),
		Port:                     util.GetIntFromInterface(d.Get("https_port")),
		NoAuthenticationRequired: util.GetBoolFromInterface(d.Get("no_authentication_required")),
		UsernamePasswordRequired: util.GetBoolFromInterface(d.Get("https_username_password_required")),
		HTTPUsername:             util.GetStringFromInterface(d.Get("https_username")),
		HTTPPassword:             util.GetStringFromInterface(d.Get("https_password")),
	}

	// Serialize and pretty-print the dockitem object as XML
	resourceXML, err := xml.MarshalIndent(fileShareDistributionPoint, "", "  ")
	if err != nil {
		logging.LogTFConstructResourceXMLMarshalFailure(subCtx, JamfProResourceDistributionPoint, err.Error())
		return nil, err
	}

	// Log the successful construction and serialization to XML
	logging.LogTFConstructedXMLResource(subCtx, JamfProResourceDistributionPoint, string(resourceXML))

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
	var creationResponse *jamfpro.ResourceFileShareDistributionPoint
	var apiErrorCode int

	// Extract values from the Terraform configuration for func useage
	resourceName := d.Get("name").(string)

	// Initialize the logging subsystem with the create operation context
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemCreate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Construct the dockitem object outside the retry loop to avoid reconstructing it on each retry
	fileShareDistributionPoint, err := constructJamfProFileShareDistributionPoint(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceDistributionPoint, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceDistributionPoint)

	// Retry the API call to create the dockitem in Jamf Pro
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		var apiErr error
		creationResponse, apiErr = conn.CreateDistributionPoint(fileShareDistributionPoint)
		if apiErr != nil {
			// Extract and log the API error code if available
			if apiError, ok := apiErr.(*.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPICreateFailedAfterRetry(subCtx, JamfProResourceDistributionPoint, resourceName, apiErr.Error(), apiErrorCode)
			// Return a non-retryable error to break out of the retry loop
			return retry.NonRetryableError(apiErr)
		}
		// No error, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final error and append it to the diagnostics
		logging.LogAPICreateFailure(subCtx, JamfProResourceDistributionPoint, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Log successful creation of the dockitem and set the resource ID in Terraform state
	logging.LogAPICreateSuccess(subCtx, JamfProResourceDistributionPoint, strconv.Itoa(creationResponse.ID))

	// set resource ID in the state
	d.SetId(strconv.Itoa(creationResponse.ID))

	// Retry reading the dockitem to ensure the Terraform state is up to date
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProFileShareDistributionPointsRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			// Log any read errors and return a retryable error to retry the read operation
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceDistributionPoint, d.Id(), readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		// Successfully read the dockitem, exit the retry loop
		return nil
	})

	if err != nil {
		// Log the final state sync failure and append it to the diagnostics
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceDistributionPoint, err.Error())
		diags = append(diags, diag.FromErr(err)...)
	} else {
		// Log successful state synchronization
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceDistributionPoint, d.Id())
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
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the read operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemRead, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	var apiErrorCode int
	var fileShareDistributionPoint *jamfpro.ResourceFileShareDistributionPoint

	resourceID := d.Id()

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceDistributionPoint, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Read operation with retry
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var apiErr error
		fileShareDistributionPoint, apiErr = conn.GetDistributionPointByID(resourceIDInt)
		if apiErr != nil {
			logging.LogFailedReadByID(subCtx, JamfProResourceDistributionPoint, resourceID, apiErr.Error(), apiErrorCode)
			// Convert any API error into a retryable error to continue retrying
			return retry.RetryableError(apiErr)
		}
		// Successfully read the script, exit the retry loop
		return nil
	})

	if err != nil {
		// Handle the final error after all retries have been exhausted
		d.SetId("") // Remove from Terraform state if unable to read after retries
		logging.LogTFStateRemovalWarning(subCtx, JamfProResourceDistributionPoint, resourceID)
		return diag.FromErr(err)
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
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the update operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemUpdate, hclog.Info)
	subSyncCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemSync, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceDistributionPoint, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Construct the resource object
	fileShareDistributionPoint, err := constructJamfProFileShareDistributionPoint(subCtx, d)
	if err != nil {
		logging.LogTFConstructResourceFailure(subCtx, JamfProResourceDistributionPoint, err.Error())
		return diag.FromErr(err)
	}
	logging.LogTFConstructResourceSuccess(subCtx, JamfProResourceDistributionPoint)

	// Update operations with retries
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		_, apiErr := conn.UpdateDistributionPointByID(resourceIDInt, fileShareDistributionPoint)
		if apiErr != nil {
			if apiError, ok := apiErr.(*.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}

			logging.LogAPIUpdateFailureByID(subCtx, JamfProResourceDistributionPoint, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			_, apiErrByName := conn.UpdateDistributionPointByName(resourceName, fileShareDistributionPoint)
			if apiErrByName != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErrByName.(*.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIUpdateFailureByName(subCtx, JamfProResourceDistributionPoint, resourceName, apiErrByName.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErrByName)
			}
		} else {
			logging.LogAPIUpdateSuccess(subCtx, JamfProResourceDistributionPoint, resourceID, resourceName)
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceDistributionPoint, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	// Retry reading the Site to synchronize the Terraform state
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		readDiags := ResourceJamfProFileShareDistributionPointsRead(subCtx, d, meta)
		if len(readDiags) > 0 {
			logging.LogTFStateSyncFailedAfterRetry(subSyncCtx, JamfProResourceDistributionPoint, resourceID, readDiags[0].Summary)
			return retry.RetryableError(fmt.Errorf(readDiags[0].Summary))
		}
		return nil
	})

	if err != nil {
		logging.LogTFStateSyncFailure(subSyncCtx, JamfProResourceDistributionPoint, err.Error())
		return diag.FromErr(err)
	} else {
		logging.LogTFStateSyncSuccess(subSyncCtx, JamfProResourceDistributionPoint, resourceID)
	}

	return nil
}

// ResourceJamfProFileShareDistributionPointsDeleteis responsible for deleting a Jamf Pro file share distribution point from the remote system.
func ResourceJamfProFileShareDistributionPointsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Initialize api client
	apiclient, ok := meta.(*client.APIClient)
	if !ok {
		return diag.Errorf("error asserting meta as *client.APIClient")
	}
	conn := apiclient.Conn

	// Initialize the logging subsystem for the delete operation
	subCtx := logging.NewSubsystemLogger(ctx, logging.SubsystemDelete, hclog.Info)

	// Initialize variables
	var diags diag.Diagnostics
	resourceID := d.Id()
	resourceName := d.Get("name").(string)
	var apiErrorCode int

	// Convert resourceID from string to int
	resourceIDInt, err := strconv.Atoi(resourceID)
	if err != nil {
		// Handle conversion error with structured logging
		logging.LogTypeConversionFailure(subCtx, "string", "int", JamfProResourceDistributionPoint, resourceID, err.Error())
		return diag.FromErr(err)
	}

	// Use the retry function for the delete operation with appropriate timeout
	err = retry.RetryContext(subCtx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		// Delete By ID
		apiErr := conn.DeleteDistributionPointByID(resourceIDInt)
		if apiErr != nil {
			if apiError, ok := apiErr.(*.APIError); ok {
				apiErrorCode = apiError.StatusCode
			}
			logging.LogAPIDeleteFailureByID(subCtx, JamfProResourceDistributionPoint, resourceID, resourceName, apiErr.Error(), apiErrorCode)

			// If Delete by ID fails then try Delete by Name
			apiErr = conn.DeleteDistributionPointByName(resourceName)
			if apiErr != nil {
				var apiErrByNameCode int
				if apiErrorByName, ok := apiErr.(*.APIError); ok {
					apiErrByNameCode = apiErrorByName.StatusCode
				}

				logging.LogAPIDeleteFailureByName(subCtx, JamfProResourceDistributionPoint, resourceName, apiErr.Error(), apiErrByNameCode)
				return retry.RetryableError(apiErr)
			}
		}
		return nil
	})

	// Send error to diag.diags
	if err != nil {
		logging.LogAPIDeleteFailedAfterRetry(subCtx, JamfProResourceDistributionPoint, resourceID, resourceName, err.Error(), apiErrorCode)
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	logging.LogAPIDeleteSuccess(subCtx, JamfProResourceDistributionPoint, resourceID, resourceName)

	// Clear the ID from the Terraform state as the resource has been deleted
	d.SetId("")

	return nil
}
