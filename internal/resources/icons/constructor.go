// icons_object.go
package icons

import (
	"fmt"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func construct(d *schema.ResourceData) (string, error) {
	// Check if we have a local file path
	if filePath := d.Get("icon_file_path").(string); filePath != "" {
		return filePath, nil
	}

	// Check if we have a web source
	if webSource := d.Get("icon_file_web_source").(string); webSource != "" {
		// Download the file and return the local path
		localPath, err := common.DownloadFile(webSource)
		if err != nil {
			return "", fmt.Errorf("failed to download icon from %s: %v", webSource, err)
		}
		return localPath, nil
	}

	return "", fmt.Errorf("either icon_file_path or icon_file_web_source must be specified")
}
