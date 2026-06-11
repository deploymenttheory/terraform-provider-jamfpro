package plist

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompactStructuralWhitespace_RemovesInterArrayWhitespace(t *testing.T) {
	pretty := `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>Pages</key>
    <array>
      <array>
        <dict>
          <key>Type</key>
          <string>Application</string>
        </dict>
      </array>
    </array>
  </dict>
</plist>`

	out, err := CompactStructuralWhitespace([]byte(pretty))
	require.NoError(t, err)
	got := string(out)

	// No whitespace may remain between structural tags — this is the exact
	// pattern the Classic API mis-parses into a phantom empty <array/>.
	for _, bad := range []string{">\n", ">  <", "<array>\n", "<array>  "} {
		assert.NotContains(t, got, bad, "compacted output still contains inter-tag whitespace")
	}
	assert.Contains(t, got, "<array><array><dict>", "expected structural tags to be adjacent")
	assert.Contains(t, got, "<string>Application</string>", "leaf <string> content must survive")
}

func TestCompactStructuralWhitespace_PreservesLeafContentAndComments(t *testing.T) {
	in := `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <!-- keep me -->
    <key>Label With Spaces</key>
    <string>two  spaces and
a newline</string>
    <key>Icon</key>
    <data>
    QUJD
    REVG
    </data>
  </dict>
</plist>`

	out, err := CompactStructuralWhitespace([]byte(in))
	require.NoError(t, err)
	got := string(out)

	assert.Contains(t, got, "<!-- keep me -->", "comment must be preserved")
	assert.Contains(t, got, "<string>two  spaces and\na newline</string>", "significant <string> whitespace must be preserved verbatim")
	assert.Contains(t, got, "<data>\n    QUJD\n    REVG\n    </data>", "<data> leaf content must be left untouched")
}

func TestCompactStructuralWhitespace_Idempotent(t *testing.T) {
	in := `<plist version="1.0"><dict><key>a</key><string>b</string></dict></plist>`
	out, err := CompactStructuralWhitespace([]byte(in))
	require.NoError(t, err)
	assert.Equal(t, in, string(out), "already-compact input must be unchanged")
}

func TestCompactStructuralWhitespace_MalformedReturnsInputAndError(t *testing.T) {
	in := []byte(`<plist><dict><key>oops</dict></plist>`) // unbalanced
	out, err := CompactStructuralWhitespace(in)
	assert.Error(t, err, "expected an error for malformed XML")
	assert.Equal(t, string(in), string(out), "malformed input must be returned unchanged")
}
