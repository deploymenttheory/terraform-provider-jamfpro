// icons_object.go
package icons

import (
	"fmt"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/resources/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// construct constructs a ResourceIcon object from the provided schema data.
func construct(d *schema.ResourceData) (string, error) {
	filePath := d.Get("icon_file_path").(string)
	webSource := d.Get("icon_file_web_source").(string)

	if filePath != "" && webSource != "" {
		return "", fmt.Errorf("cannot specify both icon_file_path and icon_file_web_source, choose one source only")
	}

	if filePath != "" {
		return filePath, nil
	}

	if webSource != "" {
		localPath, err := common.DownloadFile(webSource)
		if err != nil {
			return "", fmt.Errorf("failed to download icon from %s: %v", webSource, err)
		}
		return localPath, nil
	}
	return "", fmt.Errorf("either icon_file_path or icon_file_web_source must be specified")
}
