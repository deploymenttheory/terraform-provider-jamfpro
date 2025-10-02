package helpers

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseJSONFile_EmptyPath tests the empty path validation unit
func TestParseJSONFile_EmptyPath(t *testing.T) {
	content, err := ParseJSONFile("")
	
	assert.Error(t, err, "Expected error for empty file path")
	assert.Empty(t, content, "Expected empty content for empty path")
	assert.Contains(t, err.Error(), "file path for json file cannot be empty")
}

// TestParseJSONFile_ValidExtensions tests the file extension validation unit
func TestParseJSONFile_ValidExtensions(t *testing.T) {
	testCases := []struct {
		name        string
		filename    string
		shouldPass  bool
		expectedErr string
	}{
		{
			name:       "valid json extension",
			filename:   "test.json",
			shouldPass: true,
		},
		{
			name:       "valid uppercase json extension",
			filename:   "test.JSON",
			shouldPass: true,
		},
		{
			name:        "invalid txt extension",
			filename:    "test.txt",
			shouldPass:  false,
			expectedErr: "invalid file extension",
		},
		{
			name:        "invalid js extension",
			filename:    "test.js",
			shouldPass:  false,
			expectedErr: "invalid file extension",
		},
		{
			name:        "invalid yaml extension",
			filename:    "test.yaml",
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
				testContent := `{"name": "test", "value": 42}`
				require.NoError(t, os.WriteFile(tc.filename, []byte(testContent), 0644))
				defer func() {
					assert.NoError(t, os.Remove(tc.filename))
				}()

				content, err := ParseJSONFile(tc.filename)
				assert.NoError(t, err, "Expected success for %s", tc.filename)
				assert.Equal(t, testContent, content, "Content mismatch for %s", tc.filename)
			} else {
				// Test invalid extension (file doesn't need to exist)
				content, err := ParseJSONFile(tc.filename)
				assert.Error(t, err, "Expected error for %s", tc.filename)
				assert.Empty(t, content, "Expected empty content for %s", tc.filename)
				assert.Contains(t, err.Error(), tc.expectedErr, "Expected error containing '%s' for %s", tc.expectedErr, tc.filename)
			}
		})
	}
}

// TestParseJSONFile_FileExistence tests the file existence validation unit
func TestParseJSONFile_FileExistence(t *testing.T) {
	t.Run("non-existent file", func(t *testing.T) {
		nonExistentFile := "nonexistent.json"

		content, err := ParseJSONFile(nonExistentFile)
		
		assert.Error(t, err, "Expected error for non-existent file")
		assert.Empty(t, content, "Expected empty content for non-existent file")
		assert.Contains(t, err.Error(), "json file does not exist")
	})
}

// TestParseJSONFile_FileTypeValidation tests the regular file validation unit
func TestParseJSONFile_FileTypeValidation(t *testing.T) {
	t.Run("directory instead of file", func(t *testing.T) {
		dirPath := "fake_dir.json"
		require.NoError(t, os.MkdirAll(dirPath, 0755))
		defer func() {
			assert.NoError(t, os.RemoveAll(dirPath))
		}()

		content, err := ParseJSONFile(dirPath)
		
		assert.Error(t, err, "Expected error when trying to read directory")
		assert.Empty(t, content, "Expected empty content when reading directory")
		assert.Contains(t, err.Error(), "supplied path does not resolve to a file")
	})
}

// TestParseJSONFile_FileSizeValidation tests the file size limit validation unit
func TestParseJSONFile_FileSizeValidation(t *testing.T) {
	t.Run("file size within limit", func(t *testing.T) {
		testFile := "small.json"
		testContent := `{"comments": ["` + strings.Repeat("small json content ", 100) + `"]}`
		
		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		content, err := ParseJSONFile(testFile)
		
		assert.NoError(t, err, "ParseJSONFile should not fail for small file")
		assert.Equal(t, testContent, content, "Content should match for small file")
	})

	t.Run("file size exceeds limit", func(t *testing.T) {
		testFile := "large.json"
		// Create file larger than 10MB limit (create ~12MB file)
		largeArray := strings.Repeat(`"This is a long string that will be repeated to create a large JSON file",`, 200000)
		largeContent := `{"data": [` + largeArray[:len(largeArray)-1] + `]}`
		
		require.NoError(t, os.WriteFile(testFile, []byte(largeContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		content, err := ParseJSONFile(testFile)
		
		assert.Error(t, err, "Expected error for oversized file")
		assert.Empty(t, content, "Expected empty content for oversized file")
		assert.Contains(t, err.Error(), "file too large")
	})
}

// TestParseJSONFile_PathTraversalSecurity tests the path traversal prevention unit
func TestParseJSONFile_PathTraversalSecurity(t *testing.T) {
	testCases := []struct {
		name string
		path string
	}{
		{
			name: "simple path traversal",
			path: "../../../etc/config.json",
		},
		{
			name: "complex path traversal",
			path: "../../../../../../usr/local/config.json",
		},
		{
			name: "mixed path traversal",
			path: "./../../sensitive/data.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content, err := ParseJSONFile(tc.path)
			
			assert.Error(t, err, "Expected error for path traversal attempt: %s", tc.path)
			assert.Empty(t, content, "Expected empty content for path traversal")

			// Should either fail at file existence or security validation
			hasExpectedError := strings.Contains(err.Error(), "json file does not exist") ||
				strings.Contains(err.Error(), "access denied") ||
				strings.Contains(err.Error(), "path outside project boundaries")

			assert.True(t, hasExpectedError, "Expected security-related error for %s, got: %v", tc.path, err)
		})
	}
}

// TestParseJSONFile_EmptyFileContent tests handling of empty files unit
func TestParseJSONFile_EmptyFileContent(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		emptyFile := "empty.json"
		require.NoError(t, os.WriteFile(emptyFile, []byte(""), 0644))
		defer func() {
			assert.NoError(t, os.Remove(emptyFile))
		}()

		content, err := ParseJSONFile(emptyFile)
		
		assert.NoError(t, err, "ParseJSONFile should not fail for empty file")
		assert.Empty(t, content, "Expected empty content for empty file")
	})
}

// TestParseJSONFile_ValidFileContent tests successful file reading unit
func TestParseJSONFile_ValidFileContent(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		content  string
	}{
		{
			name:     "simple json object",
			filename: "simple.json",
			content: `{
  "name": "test-resource",
  "type": "example",
  "value": 42,
  "enabled": true
}`,
		},
		{
			name:     "json array",
			filename: "array.json",
			content: `[
  {"id": 1, "name": "first"},
  {"id": 2, "name": "second"},
  {"id": 3, "name": "third"}
]`,
		},
		{
			name:     "nested json structure",
			filename: "nested.json",
			content: `{
  "metadata": {
    "version": "1.0",
    "author": "test"
  },
  "config": {
    "settings": {
      "debug": true,
      "timeout": 30
    },
    "features": ["auth", "logging"]
  }
}`,
		},
		{
			name:     "file with special characters",
			filename: "special.json",
			content:  `{"description": "JSON with special chars: åæø, 中文, 🎉", "test": true}`,
		},
		{
			name:     "compact json",
			filename: "compact.json",
			content:  `{"compact":true,"noSpaces":42,"array":[1,2,3]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.WriteFile(tc.filename, []byte(tc.content), 0644))
			defer func() {
				assert.NoError(t, os.Remove(tc.filename))
			}()

			content, err := ParseJSONFile(tc.filename)
			
			assert.NoError(t, err, "ParseJSONFile should not fail for valid content")
			assert.Equal(t, tc.content, content, "Content should match exactly")
		})
	}
}

// TestParseJSONFile_MalformedJSON tests that the parser loads malformed JSON (validation is not its responsibility)
func TestParseJSONFile_MalformedJSON(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		content  string
	}{
		{
			name:     "invalid json syntax",
			filename: "invalid.json",
			content:  `{"name": "test", "value": 42,}`, // trailing comma
		},
		{
			name:     "unclosed bracket",
			filename: "unclosed.json",
			content:  `{"name": "test"`, // missing closing brace
		},
		{
			name:     "not json at all",
			filename: "notjson.json",
			content:  `This is not JSON content at all`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.WriteFile(tc.filename, []byte(tc.content), 0644))
			defer func() {
				assert.NoError(t, os.Remove(tc.filename))
			}()

			// The parser should load the content successfully - JSON validation is not its responsibility
			content, err := ParseJSONFile(tc.filename)
			
			assert.NoError(t, err, "ParseJSONFile should load malformed JSON (validation is not its responsibility)")
			assert.Equal(t, tc.content, content, "Content should match exactly even if malformed")
		})
	}
}
