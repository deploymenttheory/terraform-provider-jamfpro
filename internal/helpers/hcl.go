package helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ParseHCLFile securely loads a Terraform configuration file from the specified path.
// Security measures:
// - Validates file extensions (.tf, .hcl)
// - Prevents path traversal attacks
// - Restricts access to files within the project directory
// - Comprehensive error logging and validation
func ParseHCLFile(filePath string) (string, error) {
	if filePath == "" {
		err := fmt.Errorf("file path for terraform file cannot be empty")
		return "", err
	}

	// Security: Validate file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".tf" && ext != ".hcl" {
		err := fmt.Errorf("invalid file extension: %s (only .tf and .hcl files allowed)", ext)
		return "", err
	}

	// Clean the path to normalize separators and remove redundant elements
	cleanPath := filepath.Clean(filePath)

	_, callerFile, _, ok := runtime.Caller(1)
	if !ok {
		err := fmt.Errorf("unable to determine calling file's directory to establish a safe base directory")
		return "", err
	}
	callerDir := filepath.Dir(callerFile)

	// Resolve absolute path and validate it's within expected boundaries
	var absPath string
	var err error

	if filepath.IsAbs(cleanPath) {
		absPath = cleanPath

		// For absolute paths, ensure they're within the project root
		// Find the project root by looking for go.mod file
		projectRoot := findProjectRoot(callerDir)
		if projectRoot == "" {
			err := fmt.Errorf("unable to determine project root for security validation")
			return "", err
		}

		// Ensure absolute path is within project boundaries
		relPath, err := filepath.Rel(projectRoot, absPath)
		if err != nil || strings.HasPrefix(relPath, "..") {
			err := fmt.Errorf("access denied: path outside project boundaries: %s", absPath)
			return "", err
		}
	} else {
		// For relative paths, resolve from caller's directory
		absPath, err = filepath.Abs(filepath.Join(callerDir, cleanPath))
		if err != nil {
			err := fmt.Errorf("failed to resolve absolute path for %s: %w", cleanPath, err)
			return "", err
		}
	}

	if strings.Contains(absPath, "..") {
		err := fmt.Errorf("path traversal detected in resolved path: %s", absPath)
		return "", err
	}

	// Verify file exists and get info
	fileInfo, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		err := fmt.Errorf("terraform file does not exist: %s", absPath)
		return "", err
	}
	if err != nil {
		err := fmt.Errorf("failed to access terraform file %s: %w", absPath, err)
		return "", err
	}

	if !fileInfo.Mode().IsRegular() {
		err := fmt.Errorf("supplied path does not resolve to a file: %s", absPath)
		return "", err
	}

	// Limit file size to prevent memory exhaustion attacks
	maxSize := int64(1 * 1024 * 1024) // 1MB limit
	if fileInfo.Size() > maxSize {
		err := fmt.Errorf("file too large: %d bytes (max: %d bytes)", fileInfo.Size(), maxSize)
		return "", err
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		err := fmt.Errorf("failed to read terraform file %s: %w", absPath, err)
		return "", err
	}

	if len(content) == 0 {
		log.Printf("[WARN] terraform file is empty: %s", absPath)
	}

	log.Printf("[DEBUG] terraform file successfully loaded: %s (%d bytes)", absPath, len(content))
	return string(content), nil
}

// findProjectRoot searches upward from the given directory to find the project root
// by looking for go.mod file or .git directory
func findProjectRoot(startDir string) string {
	dir := startDir
	for {
		// Check for go.mod file
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		// Check for .git directory as fallback
		gitPath := filepath.Join(dir, ".git")
		if info, err := os.Stat(gitPath); err == nil && info.IsDir() {
			return dir
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}
	return ""
}
