package plist

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessConfigurationProfileForDiffSuppression(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		fieldsToRemove []string
		want           string
		wantErr        bool
	}{
		{
			name: "Remove multiple standard fields",
			input: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
			<key>PayloadUUID</key>
			<string>12345678-1234-1234-1234-123456789012</string>
			<key>PayloadIdentifier</key>
			<string>com.example.profile</string>
			<key>PayloadOrganization</key>
			<string>Example Org</string>
			<key>PayloadDisplayName</key>
			<string>Test Profile</string>
			<key>PayloadDescription</key>
			<string>Test Description</string>
	</dict>
</plist>`,
			fieldsToRemove: []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization", "PayloadDisplayName"},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
			<key>PayloadDescription</key>
			<string>Test Description</string>
	</dict>
</plist>`,
			wantErr: false,
		},
		{
			name: "Remove fields from nested array of dictionaries",
			input: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
			<key>PayloadUUID</key>
			<string>12345678-1234-1234-1234-123456789012</string>
			<key>PayloadContent</key>
			<array>
					<dict>
							<key>PayloadUUID</key>
							<string>abcdef12-3456-7890-abcd-ef1234567890</string>
							<key>PayloadIdentifier</key>
							<string>com.example.profile.1</string>
							<key>PayloadOrganization</key>
							<string>Example Org</string>
							<key>PayloadDisplayName</key>
							<string>Test Profile 1</string>
							<key>PayloadType</key>
							<string>com.apple.security.pkcs1</string>
					</dict>
					<dict>
							<key>PayloadUUID</key>
							<string>fedcba98-7654-3210-fedc-ba9876543210</string>
							<key>PayloadIdentifier</key>
							<string>com.example.profile.2</string>
							<key>PayloadOrganization</key>
							<string>Example Org</string>
							<key>PayloadDisplayName</key>
							<string>Test Profile 2</string>
							<key>PayloadType</key>
							<string>com.apple.security.pkcs1</string>
					</dict>
			</array>
	</dict>
</plist>`,
			fieldsToRemove: []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization", "PayloadDisplayName"},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
			<key>PayloadContent</key>
			<array>
					<dict>
							<key>PayloadType</key>
							<string>com.apple.security.pkcs1</string>
					</dict>
					<dict>
							<key>PayloadType</key>
							<string>com.apple.security.pkcs1</string>
					</dict>
			</array>
	</dict>
</plist>`,
			wantErr: false,
		},
		{
			name: "Remove fields from deeply nested structure",
			input: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
			<key>PayloadUUID</key>
			<string>12345678-1234-1234-1234-123456789012</string>
			<key>PayloadContent</key>
			<array>
					<dict>
							<key>PayloadUUID</key>
							<string>abcdef12-3456-7890-abcd-ef1234567890</string>
							<key>NestedContent</key>
							<dict>
									<key>PayloadUUID</key>
									<string>11111111-2222-3333-4444-555555555555</string>
									<key>PayloadIdentifier</key>
									<string>com.example.nested</string>
									<key>PayloadOrganization</key>
									<string>Example Org Nested</string>
									<key>PayloadDisplayName</key>
									<string>Nested Profile</string>
									<key>Settings</key>
									<dict>
											<key>PayloadUUID</key>
											<string>99999999-8888-7777-6666-555555555555</string>
											<key>ConfigurationOption</key>
											<string>SomeValue</string>
									</dict>
							</dict>
					</dict>
			</array>
	</dict>
</plist>`,
			fieldsToRemove: []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization", "PayloadDisplayName"},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
			<key>PayloadContent</key>
			<array>
					<dict>
							<key>NestedContent</key>
							<dict>
									<key>Settings</key>
									<dict>
											<key>ConfigurationOption</key>
											<string>SomeValue</string>
									</dict>
							</dict>
					</dict>
			</array>
	</dict>
</plist>`,
			wantErr: false,
		},
		{
			name: "Complex profile with mixed content and multiple nesting levels",
			input: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
			<key>PayloadUUID</key>
			<string>12345678-1234-1234-1234-123456789012</string>
			<key>PayloadContent</key>
			<array>
					<dict>
							<key>PayloadUUID</key>
							<string>abcdef12-3456-7890-abcd-ef1234567890</string>
							<key>PayloadIdentifier</key>
							<string>com.example.profile.1</string>
							<key>PayloadOrganization</key>
							<string>Example Org</string>
							<key>PayloadContent</key>
							<data>
									SGVsbG8gV29ybGQ=
							</data>
					</dict>
					<dict>
							<key>PayloadDisplayName</key>
							<string>Test Profile 2</string>
							<key>SubSettings</key>
							<array>
									<dict>
											<key>PayloadUUID</key>
											<string>11111111-2222-3333-4444-555555555555</string>
											<key>PayloadIdentifier</key>
											<string>com.example.subsetting</string>
											<key>Setting</key>
											<true/>
									</dict>
							</array>
					</dict>
			</array>
			<key>PayloadOrganization</key>
			<string>Example Org</string>
	</dict>
</plist>`,
			fieldsToRemove: []string{"PayloadUUID", "PayloadIdentifier", "PayloadOrganization", "PayloadDisplayName"},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
			<key>PayloadContent</key>
			<array>
					<dict>
							<key>PayloadContent</key>
							<data>SGVsbG8gV29ybGQ=</data>
					</dict>
					<dict>
							<key>SubSettings</key>
							<array>
									<dict>
											<key>Setting</key>
											<true/>
									</dict>
							</array>
					</dict>
			</array>
	</dict>
</plist>`,
			wantErr: false,
		},
		{
			name: "XML Tag Normalization Test",
			input: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Test1</key>
    <true />
    <key>Test2</key>
    <true  />
    <key>Test3</key>
    <true   />
    <key>Test4</key>
    <false />
    <key>Test5</key>
    <false    />
  </dict>
</plist>`,
			fieldsToRemove: []string{},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Test1</key>
    <true/>
    <key>Test2</key>
    <true/>
    <key>Test3</key>
    <true/>
    <key>Test4</key>
    <false/>
    <key>Test5</key>
    <false/>
  </dict>
</plist>`,
			wantErr: false,
		},
		{
			name: "Base64 Data Normalization",
			input: `<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
		<dict>
			<key>Certificate1</key>
			<data>
						MIIFYjCCBEqgAwIBAgIQd70NbNs2
						+RrqIQ/E8FjTDTANBgkqhkiG9w0BAQsF
			</data>
			<key>Certificate2</key>
			<data>
						MIIFYjCCBEqgAwIBAgIQd70NbNs2+RrqIQ/E8FjTDTANBgkqhkiG9w0BAQsF
			</data>
		</dict>
	</plist>`,
			fieldsToRemove: []string{},
			want: `<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
		<dict>
			<key>Certificate1</key>
			<data>MIIFYjCCBEqgAwIBAgIQd70NbNs2+RrqIQ/E8FjTDTANBgkqhkiG9w0BAQsF</data>
			<key>Certificate2</key>
			<data>MIIFYjCCBEqgAwIBAgIQd70NbNs2+RrqIQ/E8FjTDTANBgkqhkiG9w0BAQsF</data>
		</dict>
	</plist>`,
			wantErr: false,
		},
		{
			name: "Integer Values Test",
			input: `<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
		<dict>
			<key>PayloadContent</key>
			<array>
				<dict>
					<key>PayloadVersion</key>
					<integer>1</integer>
					<key>PayloadVersion2</key>
					<integer>1</integer>
				</dict>
			</array>
		</dict>
	</plist>`,
			fieldsToRemove: []string{},
			want: `<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
		<dict>
			<key>PayloadContent</key>
			<array>
				<dict>
					<key>PayloadVersion</key>
					<integer>1</integer>
					<key>PayloadVersion2</key>
					<integer>1</integer>
				</dict>
			</array>
		</dict>
	</plist>`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessConfigurationProfileForDiffSuppression(tt.input, tt.fieldsToRemove)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Normalize both expected and actual output by removing all whitespace
			normalizedGot := normalizeWhitespace(got)
			normalizedWant := normalizeWhitespace(tt.want)

			if normalizedGot != normalizedWant {
				t.Errorf("ProcessConfigurationProfileForDiffSuppression() \nGot:  %v\nWant: %v", got, tt.want)
			}
		})
	}
}

// normalizeWhitespace removes all whitespace from a string to make comparison easier
func normalizeWhitespace(s string) string {
	// Replace all whitespace (including newlines and tabs) with a single space
	space := regexp.MustCompile(`\s+`)
	s = space.ReplaceAllString(s, " ")
	// Trim leading/trailing space
	return strings.TrimSpace(s)
}

func TestRemoveSpecifiedXMLFields(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string]interface{}
		fieldsToRemove []string
		want           map[string]interface{}
	}{
		{
			name: "Remove single field",
			input: map[string]interface{}{
				"PayloadUUID": "test-uuid",
				"PayloadType": "test-type",
			},
			fieldsToRemove: []string{"PayloadUUID"},
			want: map[string]interface{}{
				"PayloadType": "test-type",
			},
		},
		{
			name: "Remove nested field",
			input: map[string]interface{}{
				"PayloadContent": map[string]interface{}{
					"PayloadUUID": "nested-uuid",
					"PayloadType": "nested-type",
				},
			},
			fieldsToRemove: []string{"PayloadUUID"},
			want: map[string]interface{}{
				"PayloadContent": map[string]interface{}{
					"PayloadType": "nested-type",
				},
			},
		},
		{
			name: "Remove field from array",
			input: map[string]interface{}{
				"PayloadContent": []interface{}{
					map[string]interface{}{
						"PayloadUUID": "array-uuid",
						"PayloadType": "array-type",
					},
				},
			},
			fieldsToRemove: []string{"PayloadUUID"},
			want: map[string]interface{}{
				"PayloadContent": []interface{}{
					map[string]interface{}{
						"PayloadType": "array-type",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeSpecifiedXMLFields(tt.input, tt.fieldsToRemove, "")
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNormalizeBase64Content(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{
			name:  "Normalize deeply indented <data> tag content",
			input: "<data>\n\t\t\t\t\t\t\t\t\tMIIFYjCCBEqgAwIBAgIQd70NbNs2+RrqIQ/E8FjTDTANBgkqhkiG9w0BAQsF\n\t\t\t\t\t\t\t\t\tADBXMQswCQYDVQQGEwJCRTEZMBcGA1UEChMQR2xvYmFsU2lnbiBudi1zYTEQ\n\t\t\t\t\t\t\t\t\tMA4GA1UECxMHUm9vdCBDQTEbMBkGA1UEAxMSR2xvYmFsU2lnbiBSb290IENBMB4X\n\t\t\t\t\t\t\t\t\t</data>",
			want:  "<data>MIIFYjCCBEqgAwIBAgIQd70NbNs2+RrqIQ/E8FjTDTANBgkqhkiG9w0BAQsFADBXMQswCQYDVQQGEwJCRTEZMBcGA1UEChMQR2xvYmFsU2lnbiBudi1zYTEQMA4GA1UECxMHUm9vdCBDQTEbMBkGA1UEAxMSR2xvYmFsU2lnbiBSb290IENBMB4X</data>",
		},
		{
			name: "Handle nested map with indented <data> block",
			input: map[string]interface{}{
				"payload": "<data>\n\t\t\t\tMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA\n\t\t\t\t</data>",
				"other":   "regular string",
			},
			want: map[string]interface{}{
				"payload": "<data>MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA</data>",
				"other":   "regular string",
			},
		},
		{
			name: "Handle array with mixed content",
			input: []interface{}{
				"<data>\n\t\t\t\tMIICIjAN\n\t\t\t\t</data>",
				"regular string",
				map[string]interface{}{
					"nested": "<data> \n\t\tSGVsbG8= \n\t\t</data>",
				},
			},
			want: []interface{}{
				"<data>MIICIjAN</data>",
				"regular string",
				map[string]interface{}{
					"nested": "<data>SGVsbG8=</data>",
				},
			},
		},
		{
			name: "Handle multiple <data> tags in string blob",
			input: `<dict>
				<key>cert1</key>
				<data>
					MIICIjAN
					BgkqhkiG
				</data>
				<key>cert2</key>
				<data>
					SGVsbG8=
				</data>
			</dict>`,
			want: `<dict>
				<key>cert1</key>
				<data>MIICIjANBgkqhkiG</data>
				<key>cert2</key>
				<data>SGVsbG8=</data>
			</dict>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeBase64Content(tt.input)
			assert.Equal(t, tt.want, got, "normalizeBase64Content() = %v, want %v", got, tt.want)
		})
	}
}

func TestNormalizeBase64(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "Remove whitespace from base64",
			input: `MIIFYjCCBEqgAwIBAgIQd70NbNs2+RrqIQ/E8FjTDTANBgkqhkiG9w0BAQsF
								 ADBXMQswCQYDVQQGEwJCRTEZMBcGA1UEChMQR2xvYmFsU2lnbiBudi1zYTEQ`,
			want: "MIIFYjCCBEqgAwIBAgIQd70NbNs2+RrqIQ/E8FjTDTANBgkqhkiG9w0BAQsFADBXMQswCQYDVQQGEwJCRTEZMBcGA1UEChMQR2xvYmFsU2lnbiBudi1zYTEQ",
		},
		{
			name:  "Preserve XML content",
			input: "<dict><key>test</key><string>value</string></dict>",
			want:  "<dict><key>test</key><string>value</string></dict>",
		},
		{
			name:  "Deep indentation",
			input: "\n\t\t\t\tMIICIjAN\n\t\t\t\tBgkqhkiG9w0BAQ\n\t\t\t\t",
			want:  "MIICIjANBgkqhkiG9w0BAQ",
		},
		{
			name:  "Empty string",
			input: "",
			want:  "",
		},
		{
			name:  "Only whitespace",
			input: "\n\t\r ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeBase64(tt.input)
			assert.Equal(t, tt.want, got, "NormalizeBase64() = %v, want %v", got, tt.want)
		})
	}
}

func TestNormalizeXMLTags(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{
			name:  "Normalize true tag",
			input: "<true    />",
			want:  "<true/>",
		},
		{
			name:  "Normalize false tag",
			input: "<false\t/>",
			want:  "<false/>",
		},
		{
			name:  "Handle normal string",
			input: "regular string",
			want:  "regular string",
		},
		{
			name: "Handle nested map",
			input: map[string]interface{}{
				"enabled": "<true   />",
			},
			want: map[string]interface{}{
				"enabled": "<true/>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeXMLTags(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNormalizeHTMLEntitiesForDiff(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{
			name:  "Simple unescape of amp",
			input: "Hello &amp; World",
			want:  "Hello & World",
		},
		{
			name:  "Preserve already escaped &lt;br/&gt;",
			input: "Line break here: &lt;br/&gt;More text",
			want:  "Line break here: &lt;br/&gt;More text", // should not be unescaped
		},
		{
			name:  "Ignore &amp;amp; to avoid double unescape",
			input: "Weird escape: &amp;amp;data",
			want:  "Weird escape: &amp;amp;data", // don’t unescape
		},
		{
			name: "Nested map with mixed values",
			input: map[string]interface{}{
				"html_safe":   "Click here &lt;a href='https://example.com'&gt;",
				"double_amp":  "Double escaped &amp;amp; stuff",
				"just_amp":    "Some &amp; thing",
				"no_entities": "Nothing here",
			},
			want: map[string]interface{}{
				"html_safe":   "Click here &lt;a href='https://example.com'&gt;", // unchanged
				"double_amp":  "Double escaped &amp;amp; stuff",                  // unchanged
				"just_amp":    "Some & thing",                                    // unescaped
				"no_entities": "Nothing here",                                    // unchanged
			},
		},
		{
			name: "Array of values",
			input: []interface{}{
				"Normal text",
				"Text with &amp;",
				"Text with &amp;amp;",
				"Text with &lt;br/&gt;",
			},
			want: []interface{}{
				"Normal text",
				"Text with &",
				"Text with &amp;amp;",
				"Text with &lt;br/&gt;",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeHTMLEntitiesForDiff(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTrimTrailingWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Trim spaces and tabs",
			input: "line1  \t \nline2\t  \nline3",
			want:  "line1\nline2\nline3",
		},
		{
			name:  "Handle no trailing whitespace",
			input: "line1\nline2\nline3",
			want:  "line1\nline2\nline3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trimTrailingWhitespace(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEmptyStringHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "Empty string in PayloadDescription",
			input: `<key>PayloadDescription</key>
<string></string>`,
			want: `<key>PayloadDescription</key>
<string></string>`,
		},
		{
			name: "String with only whitespace",
			input: `<key>PayloadDescription</key>
<string>    </string>`,
			want: `<key>PayloadDescription</key>
<string>    </string>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeXMLTags(tt.input)
			assert.Equal(t, tt.want, result)
		})
	}
}

// Add new test for boolean value normalization
func TestBooleanValueNormalization(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Space after true",
			input: "<true />",
			want:  "<true/>",
		},
		{
			name:  "Multiple spaces after true",
			input: "<true     />",
			want:  "<true/>",
		},
		{
			name:  "Space after false",
			input: "<false />",
			want:  "<false/>",
		},
		{
			name:  "Tab after false",
			input: "<false\t/>",
			want:  "<false/>",
		},
		{
			name:  "Mixed whitespace",
			input: "<true \t  />",
			want:  "<true/>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeXMLTags(tt.input)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestNormalizeEmptyStrings(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
		want  map[string]interface{}
	}{
		{
			name: "Empty and whitespace strings",
			input: map[string]interface{}{
				"empty":    "",
				"spaces":   "   ",
				"newlines": "\n    \n",
				"tabs":     "\t\t",
				"mixed":    "  \n\t  ",
				"content":  "actual content",
			},
			want: map[string]interface{}{
				"empty":    "",
				"spaces":   "",
				"newlines": "",
				"tabs":     "",
				"mixed":    "",
				"content":  "actual content",
			},
		},
		{
			name: "Nested dictionary",
			input: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner_empty":   "",
					"inner_spaces":  "   ",
					"inner_content": "test",
				},
			},
			want: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner_empty":   "",
					"inner_spaces":  "",
					"inner_content": "test",
				},
			},
		},
		{
			name: "Array with empty strings",
			input: map[string]interface{}{
				"items": []interface{}{
					"   ",
					"\n\t",
					"content",
					"",
				},
			},
			want: map[string]interface{}{
				"items": []interface{}{
					"",
					"",
					"content",
					"",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeEmptyStrings(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProcessConfigurationProfileEmptyStrings(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		fieldsToRemove []string
		want           string
		wantErr        bool
	}{
		{
			name: "Empty string variations",
			input: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>empty</key>
	<string></string>
	<key>spaces</key>
	<string>   </string>
	<key>newlines</key>
	<string>
	
	</string>
	<key>content</key>
	<string>actual content</string>
</dict>
</plist>`,
			fieldsToRemove: []string{},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>content</key>
	<string>actual content</string>
	<key>empty</key>
	<string/>
	<key>newlines</key>
	<string/>
	<key>spaces</key>
	<string/>
</dict>
</plist>`,
			wantErr: false,
		},
		{
			name: "Mixed empty strings and content",
			input: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>PayloadDescription</key>
	<string>
	</string>
	<key>PayloadDisplayName</key>
	<string>Test Profile</string>
	<key>PayloadContent</key>
	<array>
			<dict>
					<key>EmptyField</key>
					<string>   </string>
					<key>ContentField</key>
					<string>actual content</string>
			</dict>
	</array>
</dict>
</plist>`,
			fieldsToRemove: []string{},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>PayloadContent</key>
	<array>
			<dict>
					<key>ContentField</key>
					<string>actual content</string>
					<key>EmptyField</key>
					<string/>
			</dict>
	</array>
	<key>PayloadDescription</key>
	<string/>
	<key>PayloadDisplayName</key>
	<string>Test Profile</string>
</dict>
</plist>`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessConfigurationProfileForDiffSuppression(tt.input, tt.fieldsToRemove)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Normalize both expected and actual output by removing all whitespace
			normalizedGot := normalizeWhitespace(got)
			normalizedWant := normalizeWhitespace(tt.want)

			assert.Equal(t, normalizedWant, normalizedGot)

			if normalizedGot != normalizedWant {
				t.Errorf("ProcessConfigurationProfileForDiffSuppression() got = %v, want %v", got, tt.want)
			}
		})
	}
}
