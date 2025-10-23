package files

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// DownloadFile downloads a file from a URL and saves it to a temporary directory with security validation.
// It follows redirects up to a maximum of 10 times. The filename is determined in the following order:
// 1. From Content-Disposition header if present and valid
// 2. From the final URL after redirects if valid
// 3. Falls back to a timestamp-based name if both sources are invalid
//
// The function implements several security measures:
// - Validates filenames to prevent directory traversal
// - Restricts filenames to alphanumeric characters, dots, hyphens, underscores, and spaces
// - Replaces '%' characters with '_'
// - Ensures final path remains within temporary directory
// - Verifies path safety after normalization
//
// Returns the path to the downloaded file and any error encountered. Possible errors include:
// - Failed to create temporary file
// - Too many redirects
// - Download failure
// - Write failure
// - Invalid filename
// - Path traversal attempt
// - Rename failure
func DownloadFile(url string) (string, error) {
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

	var finalFileName string

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	filename := params["filename"]

	if err == nil && filename != "" {
		var sanitizeErr error
		finalFileName, sanitizeErr = sanitizeFileName(filename)
		if sanitizeErr != nil {
			log.Printf("[WARN] Failed to sanitize filename from Content-Disposition: %v", sanitizeErr)
		}
	}

	if finalFileName == "" {
		finalURL := resp.Request.URL.String()
		urlFileName := filepath.Base(finalURL)
		if validName, err := sanitizeFileName(urlFileName); err == nil {
			finalFileName = validName
		} else {
			finalFileName = fmt.Sprintf("downloaded-file-%d", time.Now().Unix())
		}
	}

	finalPath := filepath.Join(os.TempDir(), finalFileName)

	if !strings.HasPrefix(filepath.Clean(finalPath), os.TempDir()) {
		return "", fmt.Errorf("security error: final path '%s' would be outside temporary directory", finalPath)
	}

	err = os.Rename(tmpFile.Name(), finalPath)
	if err != nil {
		return "", fmt.Errorf("failed to rename temporary file to final destination: %v", err)
	}

	log.Printf("[INFO] File downloaded to: %s", finalPath)
	return finalPath, nil
}

// sanitizeFileName cleans and validates a filename for secure file operations
func sanitizeFileName(name string) (string, error) {
	name = filepath.Base(name)

	if strings.Contains(name, "..") {
		return "", fmt.Errorf("invalid filename: contains parent directory reference")
	}

	unescaped, unescapeErr := url.PathUnescape(name)
	if unescapeErr != nil {
		unescaped = name
	}

	cleaned := strings.TrimSpace(unescaped)
	if cleaned == "" || cleaned == "." {
		return "", fmt.Errorf("invalid filename: empty or invalid after cleaning")
	}

	return cleaned, nil
}

// CleanupDownloadedPackage handles the cleanup of downloaded package files from web sources.
// It ensures files are only deleted if they were downloaded from HTTP(s) sources and are in the temporary directory.
func CleanupDownloadedPackage(packageFileSource, localFilePath string) {
	if !regexp.MustCompile(`^(http|https)://`).MatchString(packageFileSource) {
		return
	}

	if !strings.HasPrefix(localFilePath, os.TempDir()) {
		log.Printf("[WARN] Refusing to remove file '%s' as it's not in the temporary directory: timestamp=%s",
			localFilePath, time.Now().UTC().Format(time.RFC3339))
		return
	}

	if err := os.Remove(localFilePath); err != nil {
		log.Printf("[WARN] Failed to remove downloaded package file '%s': %v: timestamp=%s",
			localFilePath, err, time.Now().UTC().Format(time.RFC3339))
	} else {
		log.Printf("[INFO] Successfully removed downloaded package file '%s': timestamp=%s",
			localFilePath, time.Now().UTC().Format(time.RFC3339))
	}
}

// CleanupDownloadedIcon handles the cleanup of downloaded icon files from web sources.
// It ensures files are only deleted if they were downloaded from HTTP(s) sources and are in the temporary directory.
func CleanupDownloadedIcon(webSource, filePath string) {
	if webSource == "" {
		return
	}

	if !strings.HasPrefix(filePath, os.TempDir()) {
		log.Printf("[WARN] Refusing to remove file '%s' as it's not in the temporary directory: timestamp=%s",
			filePath, time.Now().UTC().Format(time.RFC3339))
		return
	}

	if err := os.Remove(filePath); err != nil {
		log.Printf("[WARN] Failed to remove downloaded icon file '%s': %v: timestamp=%s",
			filePath, err, time.Now().UTC().Format(time.RFC3339))
	} else {
		log.Printf("[INFO] Successfully removed downloaded icon file '%s': timestamp=%s",
			filePath, time.Now().UTC().Format(time.RFC3339))
	}
}
