// scripts_object.go
package scripts

import (
	"encoding/base64"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProScript constructs a ResourceScript object from the provided schema data.
func constructJamfProScript(d *schema.ResourceData) (*jamfpro.ResourceScript, error) {
	script := &jamfpro.ResourceScript{
		Name:           d.Get("name").(string),
		CategoryName:   d.Get("category_name").(string),
		CategoryId:     d.Get("category_id").(string),
		Info:           d.Get("info").(string),
		Notes:          d.Get("notes").(string),
		OSRequirements: d.Get("os_requirements").(string),
		Priority:       d.Get("priority").(string),
		Parameter4:     d.Get("parameter4").(string),
		Parameter5:     d.Get("parameter5").(string),
		Parameter6:     d.Get("parameter6").(string),
		Parameter7:     d.Get("parameter7").(string),
		Parameter8:     d.Get("parameter8").(string),
		Parameter9:     d.Get("parameter9").(string),
		Parameter10:    d.Get("parameter10").(string),
		Parameter11:    d.Get("parameter11").(string),
	}

	// Handle script_contents
	if scriptContent, ok := d.GetOk("script_contents"); ok && scriptContent.(string) != "" {
		script.ScriptContents = scriptContent.(string)
	} else {
		// Decode script contents from the state if not directly modified
		encodedScriptContents := d.Get("script_contents_encoded").(string)
		decodedBytes, err := base64.StdEncoding.DecodeString(encodedScriptContents)
		if err != nil {
			return nil, fmt.Errorf("error decoding script contents: %s", err)
		}
		script.ScriptContents = string(decodedBytes)
	}

	return script, nil
}
