// icons_object.go
package icons

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
		localPath, err := downloadFile(webSource)
		if err != nil {
			return "", fmt.Errorf("failed to download icon from %s: %v", webSource, err)
		}
		return localPath, nil
	}

	return "", fmt.Errorf("either icon_file_path or icon_file_web_source must be specified")
}

// downloadFile downloads a file from the given URL and saves it to a temporary file.
// If the Content-Disposition header is present in the response, it uses the filename
// from the header. Otherwise, if no filename is provided in the headers, it uses the
// final URL after any redirects to determine the filename. It also replaces any '%' characters
// in the filename with '_'.
func downloadFile(url string) (string, error) {
	tmpFile, err := os.CreateTemp("", "downloaded-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer tmpFile.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects when attempting to download file from %s", url)
			}
			return nil
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file from %s: %v", url, err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to temporary file: %v", err)
	}

	// Get the file name from the Content-Disposition header if available
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err == nil {
		if filename, ok := params["filename"]; ok {
			filename = strings.ReplaceAll(filename, "%", "_")
			finalPath := filepath.Join(os.TempDir(), filename)
			err = os.Rename(tmpFile.Name(), finalPath)
			if err != nil {
				return "", fmt.Errorf("failed to rename temporary file to final destination: %v", err)
			}
			log.Printf("[INFO] File downloaded to: %s", finalPath)
			return finalPath, nil
		}
	}

	finalURL := resp.Request.URL.String()
	fileName := filepath.Base(finalURL)
	fileName = strings.ReplaceAll(fileName, "%", "_")
	finalPath := filepath.Join(os.TempDir(), fileName)
	err = os.Rename(tmpFile.Name(), finalPath)
	if err != nil {
		return "", fmt.Errorf("failed to rename temporary file to final destination: %v", err)
	}
	log.Printf("[INFO] File downloaded to: %s", finalPath)
	return finalPath, nil
}
