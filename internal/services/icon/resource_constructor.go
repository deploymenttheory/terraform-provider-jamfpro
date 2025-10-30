// icons_object.go
package icon

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/files"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ErrMultipleIconSources = errors.New("cannot specify multiple icon sources, choose only one: icon_file_path, icon_file_web_source, or icon_file_base64")
	ErrNoIconSource        = errors.New("one of icon_file_path, icon_file_web_source, or icon_file_base64 must be specified")
	ErrDecodeBase64        = errors.New("failed to decode base64 icon data")
	ErrCreateTempFile      = errors.New("failed to create temporary file for base64 icon")
	ErrWriteTempFile       = errors.New("failed to write base64 icon data to temporary file")
	ErrDownloadIcon        = errors.New("failed to download icon")
)

// construct constructs a ResourceIcon object from the provided schema data.
func construct(d *schema.ResourceData) (string, error) {
	filePath := d.Get("icon_file_path").(string)
	webSource := d.Get("icon_file_web_source").(string)
	base64Data := d.Get("icon_file_base64").(string)

	sourcesCount := 0
	if filePath != "" {
		sourcesCount++
	}
	if webSource != "" {
		sourcesCount++
	}
	if base64Data != "" {
		sourcesCount++
	}

	if sourcesCount > 1 {
		return "", ErrMultipleIconSources
	}

	if sourcesCount == 0 {
		return "", ErrNoIconSource
	}

	if filePath != "" {
		return filePath, nil
	}

	if webSource != "" {
		localPath, err := files.DownloadFile(webSource)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrDownloadIcon, err)
		}
		return localPath, nil
	}

	if base64Data != "" {
		decodedData, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrDecodeBase64, err)
		}

		tmpFile, err := os.CreateTemp("", "icon-*.png")
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrCreateTempFile, err)
		}
		defer func() {
			_ = tmpFile.Close()
		}()

		if _, err := tmpFile.Write(decodedData); err != nil {
			_ = os.Remove(tmpFile.Name())
			return "", fmt.Errorf("%w: %w", ErrWriteTempFile, err)
		}

		return tmpFile.Name(), nil
	}

	return "", ErrNoIconSource
}
