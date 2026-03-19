// packages_state.go
package packages

import (
	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Helper function to update Terraform state
func updateState(d *schema.ResourceData, resource *jamfpro.ResourcePackage) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("package_name", resource.PackageName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("filename", resource.FileName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("category_id", resource.CategoryID); err != nil {
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
	if err := d.Set("os_requirements", resource.OSRequirements); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("fill_user_template", resource.FillUserTemplate); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("indexed", resource.Indexed); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("fill_existing_users", resource.FillExistingUsers); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("swu", resource.SWU); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("reboot_required", resource.RebootRequired); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("self_heal_notify", resource.SelfHealNotify); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("self_healing_action", resource.SelfHealingAction); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("os_install", resource.OSInstall); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("serial_number", resource.SerialNumber); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("parent_package_id", resource.ParentPackageID); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("base_path", resource.BasePath); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("suppress_updates", resource.SuppressUpdates); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("cloud_transfer_status", resource.CloudTransferStatus); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("ignore_conflicts", resource.IgnoreConflicts); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("suppress_from_dock", resource.SuppressFromDock); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("suppress_eula", resource.SuppressEula); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("suppress_registration", resource.SuppressRegistration); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("install_language", resource.InstallLanguage); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("md5", resource.MD5); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sha256", resource.SHA256); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("sha3512", resource.SHA3512); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("hash_type", resource.HashType); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("hash_value", resource.HashValue); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("size", resource.Size); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("os_installer_version", resource.OSInstallerVersion); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("manifest", resource.Manifest); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("manifest_file_name", resource.ManifestFileName); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}
	if err := d.Set("format", resource.Format); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}
