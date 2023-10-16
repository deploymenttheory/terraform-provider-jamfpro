package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceJamfProDepartments() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfProDepartmentsCreate,
		ReadContext:   resourceJamfProDepartmentsRead,
		UpdateContext: resourceJamfProDepartmentsUpdate,
		DeleteContext: resourceJamfProDepartmentsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique identifier of the department.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The unique name of the Jamf Pro department.",
			},
		},
	}
}

func resourceJamfProDepartmentsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*APIClient).conn

	departmentName := d.Get("name").(string)
	department, err := conn.CreateDepartment(departmentName)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the ID of the department in the Terraform state
	d.SetId(fmt.Sprintf("%d", department.Id))

	return resourceJamfProDepartmentsRead(ctx, d, meta)
}

func resourceJamfProDepartmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*APIClient).conn

	var diags diag.Diagnostics

	// Initially attempt to get the department by its ID.
	departmentID, err := strconv.Atoi(d.Id())
	if err == nil {
		department, err := conn.GetDepartmentByID(departmentID)
		if err != nil {
			// If there's an error fetching by ID, log a warning and continue.
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to fetch department by ID",
				Detail:   fmt.Sprintf("Failed to fetch department with ID %d: %v", departmentID, err),
			})
		} else {
			// If successfully fetched the department by ID, set the details and return.
			if err := d.Set("name", department.Name); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Error setting name",
					Detail:   err.Error(),
				})
			}
			return diags
		}
	}

	// If fetching by ID failed or wasn't possible, try to fetch by the name.
	departmentName := d.Get("name").(string)
	department, err := conn.GetDepartmentByName(departmentName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to fetch department by name",
			Detail:   fmt.Sprintf("Failed to fetch department with name %s: %v", departmentName, err),
		})
		return diags
	}

	// Set the fetched department details into the state.
	if err := d.Set("name", department.Name); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error setting name",
			Detail:   err.Error(),
		})
	}
	if err := d.Set("id", fmt.Sprintf("%d", department.Id)); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error setting ID",
			Detail:   err.Error(),
		})
	}

	return diags
}

func resourceJamfProDepartmentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*APIClient).conn

	var diags diag.Diagnostics

	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")

		// Initially attempt to update the department by its ID.
		departmentID, err := strconv.Atoi(d.Id())
		if err == nil {
			_, err = conn.UpdateDepartmentByID(departmentID, newName.(string))
		}

		// If updating by ID failed or wasn't possible, try to update by the old name.
		if err != nil {
			_, err = conn.UpdateDepartmentByName(oldName.(string), newName.(string))
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Failed to update department with old name %s to new name %s", oldName.(string), newName.(string)),
					Detail:   err.Error(),
				})
				return diags
			}
		}
	}

	// Even if the update was successful, we run the Read function to get the latest state and verify the update.
	readDiags := resourceJamfProDepartmentsRead(ctx, d, meta)
	diags = append(diags, readDiags...)

	return diags
}

func resourceJamfProDepartmentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*APIClient).conn

	var diags diag.Diagnostics

	// Initially attempt to delete the department by its ID.
	departmentID, err := strconv.Atoi(d.Id())
	if err == nil {
		err := conn.DeleteDepartmentByID(departmentID)
		if err == nil {
			// Successfully deleted the department by ID.
			return diags
		}

		// If there's an error deleting by ID, log a warning and continue.
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to delete department by ID",
			Detail:   fmt.Sprintf("Failed to delete department with ID %d: %v", departmentID, err),
		})
	}

	// If deleting by ID failed or wasn't possible, try to delete by the name.
	departmentName := d.Get("name").(string)
	err = conn.DeleteDepartmentByName(departmentName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete department by name",
			Detail:   fmt.Sprintf("Failed to delete department with name %s: %v", departmentName, err),
		})
	}

	return diags
}
