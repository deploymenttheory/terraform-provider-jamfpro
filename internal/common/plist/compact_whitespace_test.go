package plist

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"howett.net/plist"
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

// A <string> whose entire value is whitespace is content, not formatting, and
// must survive even though its parent chain is structural.
func TestCompactStructuralWhitespace_PreservesWhitespaceOnlyStringValue(t *testing.T) {
	in := `<plist version="1.0"><dict><key>Indent</key><string>   </string></dict></plist>`
	out, err := CompactStructuralWhitespace([]byte(in))
	require.NoError(t, err)
	assert.Equal(t, in, string(out), "a whitespace-only <string> value must not be stripped")
}

// Whitespace expressed as a character reference between structural children is
// not whitespace at the byte level, so the conservative byte-check must leave
// it untouched (it round-trips identically through any plist parser anyway).
func TestCompactStructuralWhitespace_PreservesCharacterReferenceWhitespace(t *testing.T) {
	in := `<plist version="1.0"><dict><key>a</key>&#x20;<string>b</string></dict></plist>`
	out, err := CompactStructuralWhitespace([]byte(in))
	require.NoError(t, err)
	assert.Contains(t, string(out), "&#x20;", "character-reference whitespace must be preserved verbatim")
}

func TestCompactStructuralWhitespace_Idempotent(t *testing.T) {
	in := `<plist version="1.0"><dict><key>a</key><string>b</string></dict></plist>`
	out, err := CompactStructuralWhitespace([]byte(in))
	require.NoError(t, err)
	assert.Equal(t, in, string(out), "already-compact input must be unchanged")
}

func TestCompactStructuralWhitespace_MalformedReturnsInputAndError(t *testing.T) {
	in := []byte(`<plist><dict><key>oops</dict></plist>`)
	out, err := CompactStructuralWhitespace(in)
	assert.Error(t, err, "expected an error for malformed XML")
	assert.Equal(t, string(in), string(out), "malformed input must be returned unchanged")
}

// A realistic, pretty-printed com.apple.homescreenlayout payload must parse to
// the exact same plist structure after compaction: compaction removes
// formatting only and never changes semantics, and the compacted Pages array
// must not gain a phantom empty leading page.
func TestCompactStructuralWhitespace_HomeScreenLayoutRoundTripsEqual(t *testing.T) {
	pretty := `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>PayloadType</key>
    <string>com.apple.homescreenlayout</string>
    <key>Dock</key>
    <array>
      <dict>
        <key>Type</key>
        <string>Application</string>
        <key>BundleID</key>
        <string>com.apple.mobilesafari</string>
      </dict>
    </array>
    <key>Pages</key>
    <array>
      <array>
        <dict>
          <key>DisplayName</key>
          <string>Apps</string>
          <key>Type</key>
          <string>Folder</string>
        </dict>
      </array>
    </array>
  </dict>
</plist>`

	compacted, err := CompactStructuralWhitespace([]byte(pretty))
	require.NoError(t, err)
	require.NotContains(t, string(compacted), "</array>\n", "expected inter-array whitespace removed")

	var prettyTree, compactTree map[string]any
	_, err = plist.Unmarshal([]byte(pretty), &prettyTree)
	require.NoError(t, err, "pretty payload must parse")
	_, err = plist.Unmarshal(compacted, &compactTree)
	require.NoError(t, err, "compacted payload must parse")

	assert.True(t, reflect.DeepEqual(prettyTree, compactTree),
		"compaction changed the parsed plist structure:\npretty:   %#v\ncompacted:%#v", prettyTree, compactTree)

	pages, _ := compactTree["Pages"].([]any)
	require.Len(t, pages, 1, "Pages must still have exactly one page (no phantom empty page)")
}
