package plist

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"howett.net/plist"
)

// corpusPayloads is a spread of realistic configuration-profile payload
// shapes: the home-screen layout that triggers the phantom-page PI, nested
// array-of-arrays, arrays of dicts, sibling arrays of strings, PPPC code
// requirements containing quotes and ampersand-adjacent characters, base64
// <data> blocks, dates/reals/integers/booleans, intentional empty
// <array/>/<string/> elements that must survive, XML comments, entity
// references, and multi-byte UTF-8 content.
var corpusPayloads = map[string]string{
	"home_screen_layout": `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
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
          <key>Pages</key>
          <array>
            <array>
              <dict>
                <key>Type</key>
                <string>Application</string>
                <key>BundleID</key>
                <string>com.apple.Preferences</string>
              </dict>
            </array>
          </array>
        </dict>
      </array>
      <array/>
    </array>
  </dict>
</plist>`,
	"restrictions": `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>PayloadType</key>
    <string>com.apple.applicationaccess</string>
    <key>allowEraseContentAndSettings</key>
    <false/>
    <key>allowDeviceNameModification</key>
    <false/>
    <key>whitelistedAppBundleIDs</key>
    <array>
      <string>com.apple.mobilesafari</string>
      <string>com.apple.camera</string>
      <string>com.jamfsoftware.selfservice</string>
    </array>
  </dict>
</plist>`,
	"pppc_code_requirement": `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>PayloadType</key>
    <string>com.apple.TCC.configuration-profile-policy</string>
    <key>Services</key>
    <dict>
      <key>ScreenCapture</key>
      <array>
        <dict>
          <key>Authorization</key>
          <string>AllowStandardUserToSetSystemService</string>
          <key>CodeRequirement</key>
          <string>identifier "com.microsoft.teams" and anchor apple generic and certificate 1[field.1.2.840.113635.100.6.2.6] /* exists */ and certificate leaf[subject.OU] = UBF8T346G9</string>
          <key>Comment</key>
          <string></string>
          <key>IdentifierType</key>
          <string>bundleID</string>
        </dict>
      </array>
    </dict>
  </dict>
</plist>`,
	"webclip_with_data_icon": `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>PayloadType</key>
    <string>com.apple.webClip.managed</string>
    <key>URL</key>
    <string>https://example.com/form?serial=$SERIALNUMBER&amp;source=clip</string>
    <key>Label</key>
    <string>Redeployment</string>
    <key>Icon</key>
    <data>
    iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR4
    2mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
    </data>
    <key>IsRemovable</key>
    <true/>
    <key>FullScreen</key>
    <true/>
  </dict>
</plist>`,
	"root_certificate": `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>PayloadType</key>
    <string>com.apple.security.root</string>
    <key>PayloadCertificateFileName</key>
    <string>ca.cer</string>
    <key>PayloadContent</key>
    <data>
    TUlJQjJUQ0NBVUtnQXdJQkFnSUJBVEFOQmdrcWhraUc5dzBCQVFzRkFE
    QVNNUkF3RGdZRFZRUUREQWRVWlhOMElFTkJNQjRYRFRJMk1ERXdNVEF3
    </data>
  </dict>
</plist>`,
	"system_preferences_sibling_string_arrays": `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>PayloadType</key>
    <string>com.apple.systempreferences</string>
    <key>DisabledPreferencePanes</key>
    <array>
      <string>com.apple.preferences.softwareupdate</string>
      <string>com.apple.preference.security</string>
    </array>
    <key>HiddenPreferencePanes</key>
    <array>
      <string>com.apple.preferences.icloud</string>
    </array>
  </dict>
</plist>`,
	"dock_nested_tile_data": `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>PayloadType</key>
    <string>com.apple.dock</string>
    <key>static-apps</key>
    <array>
      <dict>
        <key>tile-type</key>
        <string>file-tile</string>
        <key>tile-data</key>
        <dict>
          <key>file-data</key>
          <dict>
            <key>_CFURLString</key>
            <string>/Applications/Safari.app</string>
            <key>_CFURLStringType</key>
            <integer>0</integer>
          </dict>
        </dict>
      </dict>
    </array>
    <key>tilesize</key>
    <real>48.5</real>
  </dict>
</plist>`,
	"scalar_types_and_date": `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>PayloadType</key>
    <string>com.example.scalars</string>
    <key>RemovalDate</key>
    <date>2026-12-31T23:59:59Z</date>
    <key>GracePeriodHours</key>
    <integer>72</integer>
    <key>Threshold</key>
    <real>0.75</real>
    <key>Enabled</key>
    <true/>
    <key>Disabled</key>
    <false/>
  </dict>
</plist>`,
	"intentional_empty_containers": `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>PayloadType</key>
    <string>com.example.empties</string>
    <key>EmptyArray</key>
    <array/>
    <key>EmptyString</key>
    <string></string>
    <key>SelfClosedString</key>
    <string/>
    <key>EmptyDict</key>
    <dict/>
    <key>ArrayOfEmptyArrays</key>
    <array>
      <array/>
      <array/>
    </array>
  </dict>
</plist>`,
	"comments_and_entities": `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <!-- deployment note: managed by terraform -->
    <key>PayloadType</key>
    <string>com.example.entities</string>
    <key>Expression</key>
    <string>a &amp; b &lt; c &gt; d &quot;quoted&quot;</string>
    <key>Path</key>
    <string>C:\Program Files\Example &amp; Co</string>
  </dict>
</plist>`,
	"multibyte_unicode": `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>PayloadType</key>
    <string>com.example.unicode</string>
    <key>Organisation</key>
    <string>Müller &amp; Sørensen GmbH — Test Org</string>
    <key>Emoji</key>
    <string>📱 デバイス管理 устройство</string>
  </dict>
</plist>`,
}

// compactAndRequireSemanticEquality routes one payload variant through
// CompactStructuralWhitespace and asserts the three corpus invariants:
// compaction succeeds for well-formed input, the compacted bytes parse to the
// identical plist tree as the input, and compaction is idempotent (a second
// pass changes nothing, proving no strippable whitespace remains).
func compactAndRequireSemanticEquality(t *testing.T, label string, raw []byte) {
	t.Helper()

	compacted, err := CompactStructuralWhitespace(raw)
	require.NoError(t, err, "%s: compaction must succeed", label)

	var inputTree, compactTree any
	_, err = plist.Unmarshal(raw, &inputTree)
	require.NoError(t, err, "%s: input must parse", label)
	_, err = plist.Unmarshal(compacted, &compactTree)
	require.NoError(t, err, "%s: compacted output must parse", label)
	assert.True(t, reflect.DeepEqual(inputTree, compactTree),
		"%s: compaction changed the parsed plist structure", label)

	again, err := CompactStructuralWhitespace(compacted)
	require.NoError(t, err, "%s: re-compaction must succeed", label)
	assert.Equal(t, string(compacted), string(again), "%s: compaction must be idempotent", label)
}

// pipelineVariants derives the wire-shaped permutations of one payload that
// the provider can realistically produce: the pretty source as written, a
// CRLF-line-ending copy, howett re-encodes at both indent styles used in the
// codebase plus the compact style, and the update-path shape (pretty howett
// encode followed by the constructors' "&#34;" -> `"` unescape).
func pipelineVariants(t *testing.T, original string) map[string][]byte {
	t.Helper()

	var tree any
	_, err := plist.Unmarshal([]byte(original), &tree)
	require.NoError(t, err, "corpus payload must parse")

	encode := func(indent string) []byte {
		var buf bytes.Buffer
		encoder := plist.NewEncoder(&buf)
		if indent != "" {
			encoder.Indent(indent)
		}
		require.NoError(t, encoder.Encode(tree))
		return buf.Bytes()
	}

	prettyEncoded := encode("    ")
	return map[string][]byte{
		"original":           []byte(original),
		"crlf":               []byte(strings.ReplaceAll(original, "\n", "\r\n")),
		"howett_tab":         encode("\t"),
		"howett_four_spaces": prettyEncoded,
		"howett_compact":     encode(""),
		"update_path_shape":  []byte(strings.ReplaceAll(string(prettyEncoded), "&#34;", "\"")),
	}
}

// TestCompactStructuralWhitespace_PayloadCorpus runs every corpus payload
// through every pipeline variant and asserts the corpus invariants hold for
// all of them. This is the offline analogue of the 200-profile roundtrip
// corpus used to validate the equivalent fix in the jamfplatform provider:
// instead of live server echoes it proves, for a spread of payload shapes and
// whitespace permutations, that compaction never alters plist semantics.
func TestCompactStructuralWhitespace_PayloadCorpus(t *testing.T) {
	for name, payload := range corpusPayloads {
		for variant, raw := range pipelineVariants(t, payload) {
			compactAndRequireSemanticEquality(t, name+"/"+variant, raw)
		}
	}
}

// TestCompactStructuralWhitespace_RepoMobileconfigCorpus applies the same
// invariants to every .mobileconfig checked into testing/payloads — the exact
// files the integration workflow sends to a real Jamf instance.
func TestCompactStructuralWhitespace_RepoMobileconfigCorpus(t *testing.T) {
	root := filepath.Join("..", "..", "..", "testing", "payloads")
	var processed int
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".mobileconfig") {
			return nil
		}
		// #nosec G304,G122 -- path comes from a WalkDir over the repo's own
		// checked-in testing/payloads tree, not from user input, so neither
		// file inclusion nor symlink TOCTOU is reachable here.
		raw, readErr := os.ReadFile(path)
		require.NoError(t, readErr)
		compactAndRequireSemanticEquality(t, d.Name(), raw)
		processed++
		return nil
	})
	require.NoError(t, err)
	require.NotZero(t, processed, "expected .mobileconfig files under %s", root)
	t.Logf("processed %d repo payload files", processed)
}
