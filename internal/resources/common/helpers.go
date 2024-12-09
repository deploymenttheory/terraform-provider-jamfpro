// hash.go
// This package contains shared / common hash functions
package common

import (
	"crypto/sha256"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// HashString calculates the SHA-256 hash of a string and returns it as a hexadecimal string.
func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hash := fmt.Sprintf("%x", h.Sum(nil))
	log.Printf("Computed hash: %s", hash)
	return hash
}

// SerializeAndRedactXML serializes a resource to XML and redacts specified fields.
func SerializeAndRedactXML(resource interface{}, redactFields []string) (string, error) {
	v := reflect.ValueOf(resource)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("resource must be a pointer to a struct")
	}

	resourceCopy := reflect.New(v.Elem().Type()).Elem()
	resourceCopy.Set(v.Elem())

	for _, field := range redactFields {
		if f := resourceCopy.FieldByName(field); f.IsValid() && f.CanSet() {
			if f.Kind() == reflect.String {
				f.SetString("***REDACTED***")
			}
		}
	}

	if marshaledXML, err := xml.MarshalIndent(resourceCopy.Interface(), "", "  "); err != nil {
		return "", fmt.Errorf("failed to marshal %s to XML: %v", v.Elem().Type(), err)
	} else {
		return string(marshaledXML), nil
	}
}

// SerializeAndRedactJSON serializes a resource to JSON and redacts specified fields.
func SerializeAndRedactJSON(resource interface{}, redactFields []string) (string, error) {
	v := reflect.ValueOf(resource)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return "", fmt.Errorf("resource must be a pointer to a struct")
	}

	resourceCopy := reflect.New(v.Elem().Type()).Elem()
	resourceCopy.Set(v.Elem())

	for _, field := range redactFields {
		if f := resourceCopy.FieldByName(field); f.IsValid() && f.CanSet() {
			if f.Kind() == reflect.String {
				f.SetString("***REDACTED***")
			}
		}
	}

	marshaledJSON, err := json.MarshalIndent(resourceCopy.Interface(), "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal %s to JSON: %v", v.Elem().Type(), err)
	}

	return string(marshaledJSON), nil
}

func getIDField(response interface{}) (any, error) {
	v := reflect.ValueOf(response).Elem()

	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return "", fmt.Errorf("ID field not found in response")
	}

	str, ok := idField.Interface().(string)
	if ok {
		return str, nil
	}

	integer, ok := idField.Interface().(int)
	if ok {
		return strconv.Itoa(integer), nil
	}

	return nil, fmt.Errorf("unsupported type")
}

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

	// Default is 10 anyway TODO remove this for now and test if it's needed to be specified
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

	validateFileName := func(name string) (string, error) {
		name = filepath.Base(name)

		if strings.Contains(name, "..") {
			return "", fmt.Errorf("invalid filename '%s': contains parent directory reference", name)
		}

		name = strings.ReplaceAll(name, "%", "_")

		if !regexp.MustCompile(`^[\w\-\. ]+$`).MatchString(name) {
			return "", fmt.Errorf("invalid filename '%s': must only contain alphanumeric characters, dots, hyphens, underscores, and spaces", name)
		}

		return name, nil
	}

	var finalFileName string

	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	filename := params["filename"]

	if err == nil && filename != "" {
		if finalFileName, err = validateFileName(filename); err == nil {
		}
	}

	if finalFileName == "" {
		finalURL := resp.Request.URL.String()
		urlFileName := filepath.Base(finalURL)
		if validName, err := validateFileName(urlFileName); err == nil {
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
