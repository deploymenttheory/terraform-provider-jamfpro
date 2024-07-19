// printers_data_object.go
package printers

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProPrinter constructs a ResourcePrinter object from the provided schema data.
func construct(d *schema.ResourceData) (*jamfpro.ResourcePrinter, error) {
	var resource *jamfpro.ResourcePrinter

	resource = &jamfpro.ResourcePrinter{
		Name:        d.Get("name").(string),
		Category:    d.Get("category_name").(string),
		URI:         d.Get("uri").(string),
		CUPSName:    d.Get("cups_name").(string),
		Location:    d.Get("location").(string),
		Model:       d.Get("model").(string),
		Info:        d.Get("info").(string),
		Notes:       d.Get("notes").(string),
		MakeDefault: d.Get("make_default").(bool),
		UseGeneric:  d.Get("use_generic").(bool),
		PPD:         d.Get("ppd").(string),
		PPDPath:     d.Get("ppd_path").(string),
		PPDContents: d.Get("ppd_contents").(string),
	}

	resourceXML, err := xml.MarshalIndent(resource, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro Printer '%s' to XML: %v", resource.Name, err)
	}

	log.Printf("[DEBUG] Constructed Jamf Pro Printer XML:\n%s\n", string(resourceXML))

	return resource, nil
}
