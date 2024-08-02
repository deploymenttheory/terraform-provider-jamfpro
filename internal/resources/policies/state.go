package policies

import (
	"strconv"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Parent func for invdividual stating functions
func updateState(d *schema.ResourceData, resp *jamfpro.ResourcePolicy) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("id", strconv.Itoa(resp.General.ID)); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// General/Root level
	stateGeneral(d, resp, &diags)

	// Scope
	stateScope(d, resp, &diags)

	// Self Service
	stateSelfService(d, resp, &diags)

	// Payloads
	statePayloads(d, resp, &diags)

	return diags
}
