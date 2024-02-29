// printers_data_object.go
package printers

import (
	"encoding/xml"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// constructJamfProPrinter constructs a ResourcePrinter object from the provided schema data.
func constructJamfProPrinter(d *schema.ResourceData) (*jamfpro.ResourcePrinter, error) {
	printer := &jamfpro.ResourcePrinter{
		Name:        d.Get("name").(string),
		Category:    d.Get("category").(string),
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

	// Serialize and pretty-print the site object as XML
	resourceXML, err := xml.MarshalIndent(printer, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Jamf Pro printer '%s' to XML: %v", printer.Name, err)
	}
	fmt.Printf("Constructed Jamf Pro Printer XML:\n%s\n", string(resourceXML))

	return printer, nil
}
