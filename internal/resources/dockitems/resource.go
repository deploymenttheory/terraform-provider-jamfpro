// dockitems_resource.go
package dockitems

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceJamfProDockItems defines the schema and CRUD operations for managing Jamf Pro Dock Items in Terraform.
func ResourceJamfProDockItems() *schema.Resource {
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
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the dock item.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the dock item.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the dock item (App/File/Folder).",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v, ok := val.(string)
					if !ok {
						errs = append(errs, fmt.Errorf("expected a string for %q but got a different type", key))
						return
					}
					validTypes := map[string]bool{
						"App":    true,
						"File":   true,
						"Folder": true,
					}
					if !validTypes[v] {
						errs = append(errs, fmt.Errorf("%q must be one of 'App', 'File', or 'Folder', got: %s", key, v))
					}
					return
				},
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path of the dock item.",
			},
			"contents": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contents of the dock item.",
			},
		},
	}
}
