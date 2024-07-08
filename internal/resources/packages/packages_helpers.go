// packages_helpers.go
package packages

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// generateMD5FileHash accepts a file path and returns an MD5 hash of the file's contents.
// It opens the file, creates a new MD5 hash object, copies the file content into the hash object, and computes the MD5 checksum of the file.
func generateMD5FileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file contents of %s: %v", filePath, err)
	}

	hashBytes := hash.Sum(nil)

	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}
