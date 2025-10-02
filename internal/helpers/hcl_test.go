package helpers

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseHCLFile_EmptyPath tests the empty path validation unit
func TestParseHCLFile_EmptyPath(t *testing.T) {
	content, err := ParseHCLFile("")
	
	assert.Error(t, err, "Expected error for empty file path")
	assert.Empty(t, content, "Expected empty content for empty path")
	assert.Contains(t, err.Error(), "file path for terraform file cannot be empty")
}

// TestParseHCLFile_ValidExtensions tests the file extension validation unit
func TestParseHCLFile_ValidExtensions(t *testing.T) {
	testCases := []struct {
		name        string
		filename    string
		shouldPass  bool
		expectedErr string
	}{
		{
			name:       "valid tf extension",
			filename:   "test.tf",
			shouldPass: true,
		},
		{
			name:       "valid hcl extension",
			filename:   "test.hcl",
			shouldPass: true,
		},
		{
			name:       "valid uppercase tf extension",
			filename:   "test.TF",
			shouldPass: true,
		},
		{
			name:       "valid uppercase hcl extension",
			filename:   "test.HCL",
			shouldPass: true,
		},
		{
			name:        "invalid txt extension",
			filename:    "test.txt",
			shouldPass:  false,
			expectedErr: "invalid file extension",
		},
		{
			name:        "invalid json extension",
			filename:    "test.json",
			shouldPass:  false,
			expectedErr: "invalid file extension",
		},
		{
			name:        "no extension",
			filename:    "test",
			shouldPass:  false,
			expectedErr: "invalid file extension",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldPass {
				// Create a valid test file for valid extensions
				testContent := "resource \"test\" \"example\" {}"
				require.NoError(t, os.WriteFile(tc.filename, []byte(testContent), 0644))
				defer func() {
					assert.NoError(t, os.Remove(tc.filename))
				}()

				content, err := ParseHCLFile(tc.filename)
				assert.NoError(t, err, "Expected success for %s", tc.filename)
				assert.Equal(t, testContent, content, "Content mismatch for %s", tc.filename)
			} else {
				// Test invalid extension (file doesn't need to exist)
				content, err := ParseHCLFile(tc.filename)
				assert.Error(t, err, "Expected error for %s", tc.filename)
				assert.Empty(t, content, "Expected empty content for %s", tc.filename)
				assert.Contains(t, err.Error(), tc.expectedErr, "Expected error containing '%s' for %s", tc.expectedErr, tc.filename)
			}
		})
	}
}

// TestParseHCLFile_FileExistence tests the file existence validation unit
func TestParseHCLFile_FileExistence(t *testing.T) {
	t.Run("non-existent file", func(t *testing.T) {
		nonExistentFile := "nonexistent.tf"

		content, err := ParseHCLFile(nonExistentFile)
		
		assert.Error(t, err, "Expected error for non-existent file")
		assert.Empty(t, content, "Expected empty content for non-existent file")
		assert.Contains(t, err.Error(), "terraform file does not exist")
	})
}

// TestParseHCLFile_FileTypeValidation tests the regular file validation unit
func TestParseHCLFile_FileTypeValidation(t *testing.T) {
	t.Run("directory instead of file", func(t *testing.T) {
		dirPath := "fake_dir.tf"
		require.NoError(t, os.MkdirAll(dirPath, 0755))
		defer func() {
			assert.NoError(t, os.RemoveAll(dirPath))
		}()

		content, err := ParseHCLFile(dirPath)
		
		assert.Error(t, err, "Expected error when trying to read directory")
		assert.Empty(t, content, "Expected empty content when reading directory")
		assert.Contains(t, err.Error(), "supplied path does not resolve to a file")
	})
}

// TestParseHCLFile_FileSizeValidation tests the file size limit validation unit
func TestParseHCLFile_FileSizeValidation(t *testing.T) {
	t.Run("file size within limit", func(t *testing.T) {
		testFile := "small.tf"
		testContent := strings.Repeat("# comment\n", 100) // Small file
		
		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		content, err := ParseHCLFile(testFile)
		
		assert.NoError(t, err, "ParseHCLFile should not fail for small file")
		assert.Equal(t, testContent, content, "Content should match for small file")
	})

	t.Run("file size exceeds limit", func(t *testing.T) {
		testFile := "large.tf"
		// Create file larger than 1MB limit (create ~2MB file)
		largeContent := strings.Repeat("# This is a comment line that will be repeated to create a large file\n", 30000)
		
		require.NoError(t, os.WriteFile(testFile, []byte(largeContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		content, err := ParseHCLFile(testFile)
		
		assert.Error(t, err, "Expected error for oversized file")
		assert.Empty(t, content, "Expected empty content for oversized file")
		assert.Contains(t, err.Error(), "file too large")
	})
}

// TestParseHCLFile_PathTraversalSecurity tests the path traversal prevention unit
func TestParseHCLFile_PathTraversalSecurity(t *testing.T) {
	testCases := []struct {
		name string
		path string
	}{
		{
			name: "simple path traversal",
			path: "../../../etc/passwd.tf",
		},
		{
			name: "complex path traversal",
			path: "../../../../../../usr/bin/evil.tf",
		},
		{
			name: "mixed path traversal",
			path: "./../../sensitive/file.tf",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content, err := ParseHCLFile(tc.path)
			
			assert.Error(t, err, "Expected error for path traversal attempt: %s", tc.path)
			assert.Empty(t, content, "Expected empty content for path traversal")

			// Should either fail at file existence or security validation
			hasExpectedError := strings.Contains(err.Error(), "terraform file does not exist") ||
				strings.Contains(err.Error(), "access denied") ||
				strings.Contains(err.Error(), "path outside project boundaries")

			assert.True(t, hasExpectedError, "Expected security-related error for %s, got: %v", tc.path, err)
		})
	}
}

// TestParseHCLFile_EmptyFileContent tests handling of empty files unit
func TestParseHCLFile_EmptyFileContent(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		emptyFile := "empty.tf"
		require.NoError(t, os.WriteFile(emptyFile, []byte(""), 0644))
		defer func() {
			assert.NoError(t, os.Remove(emptyFile))
		}()

		content, err := ParseHCLFile(emptyFile)
		
		assert.NoError(t, err, "ParseHCLFile should not fail for empty file")
		assert.Empty(t, content, "Expected empty content for empty file")
	})
}

// TestParseHCLFile_ValidFileContent tests successful file reading unit
func TestParseHCLFile_ValidFileContent(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		content  string
	}{
		{
			name:     "simple terraform resource",
			filename: "simple.tf",
			content: `resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}`,
		},
		{
			name:     "hcl variable",
			filename: "vars.hcl",
			content: `variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}`,
		},
		{
			name:     "file with special characters",
			filename: "special.tf",
			content:  `# Comment with special chars: åæø, 中文, 🎉\nresource "test" "special" {}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.WriteFile(tc.filename, []byte(tc.content), 0644))
			defer func() {
				assert.NoError(t, os.Remove(tc.filename))
			}()

			content, err := ParseHCLFile(tc.filename)
			
			assert.NoError(t, err, "ParseHCLFile should not fail for valid content")
			assert.Equal(t, tc.content, content, "Content should match exactly")
		})
	}
}