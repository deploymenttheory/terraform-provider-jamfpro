package macosconfigurationprofilesplist

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestDiffSuppressEquivalentPayloads(t *testing.T) {
	tests := []struct {
		name            string
		old             string
		new             string
		payloadValidate bool
		wantSuppressed  bool
	}{
		{
			name: "Different indentation",
			old: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadContent</key>
    <array>
        <dict>
            <key>PayloadType</key>
            <string>com.apple.security.pkcs1</string>
            <key>PayloadVersion</key>
            <integer>1</integer>
        </dict>
    </array>
</dict>
</plist>`,
			new: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict><key>PayloadContent</key><array><dict><key>PayloadType</key><string>com.apple.security.pkcs1</string><key>PayloadVersion</key><integer>1</integer></dict></array></dict>
</plist>`,
			payloadValidate: true,
			wantSuppressed:  true,
		},
		{
			name: "Different newlines and spaces in boolean tags",
			old: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>AllowAllAppsAccess</key>
    <true />
    <key>KeyIsExtractable</key>
    <false />
</dict>
</plist>`,
			new: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>AllowAllAppsAccess</key>
    <true/>
    <key>KeyIsExtractable</key>
    <false/>
</dict>
</plist>`,
			payloadValidate: true,
			wantSuppressed:  true,
		},
		{
			name: "Different base64 data formatting",
			old: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadContent</key>
    <data>
        SGVsbG8gV29ybGQ=
    </data>
</dict>
</plist>`,
			new: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadContent</key>
    <data>SGVsbG8gV29ybGQ=</data>
</dict>
</plist>`,
			payloadValidate: true,
			wantSuppressed:  true,
		},
		{
			name: "Different ordering of dictionary keys",
			old: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadType</key>
    <string>com.apple.security.pkcs1</string>
    <key>PayloadVersion</key>
    <integer>1</integer>
</dict>
</plist>`,
			new: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadVersion</key>
    <integer>1</integer>
    <key>PayloadType</key>
    <string>com.apple.security.pkcs1</string>
</dict>
</plist>`,
			payloadValidate: true,
			wantSuppressed:  true,
		},
		{
			name: "Different whitespace in empty strings",
			old: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadDescription</key>
    <string>
    </string>
</dict>
</plist>`,
			new: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadDescription</key>
    <string></string>
</dict>
</plist>`,
			payloadValidate: true,
			wantSuppressed:  true,
		},
		{
			name: "Complex nested structure with different formatting",
			old: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadContent</key>
    <array>
        <dict>
            <key>PayloadType</key>
            <string>com.apple.security.pkcs1</string>
            <key>PayloadVersion</key>
            <integer>1</integer>
            <key>PayloadContent</key>
            <data>
                SGVsbG8gV29ybGQ=
            </data>
        </dict>
    </array>
    <key>PayloadDescription</key>
    <string></string>
    <key>PayloadDisplayName</key>
    <string>Test Profile</string>
    <key>PayloadIdentifier</key>
    <string>com.example.profile</string>
    <key>PayloadOrganization</key>
    <string>Example Org</string>
    <key>PayloadUUID</key>
    <string>123e4567-e89b-12d3-a456-426614174000</string>
    <key>PayloadVersion</key>
    <integer>1</integer>
</dict>
</plist>`,
			new: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadVersion</key>
    <integer>1</integer>
    <key>PayloadUUID</key>
    <string>987fcdeb-43a1-12d3-a456-426614174000</string>
    <key>PayloadOrganization</key>
    <string>Example Org</string>
    <key>PayloadIdentifier</key>
    <string>com.example.profile</string>
    <key>PayloadDisplayName</key>
    <string>Test Profile</string>
    <key>PayloadDescription</key>
    <string></string>
    <key>PayloadContent</key>
    <array><dict>
        <key>PayloadContent</key><data>SGVsbG8gV29ybGQ=</data>
        <key>PayloadVersion</key><integer>1</integer>
        <key>PayloadType</key><string>com.apple.security.pkcs1</string>
    </dict></array>
</dict>
</plist>`,
			payloadValidate: true,
			wantSuppressed:  true,
		},
		{
			name: "Actually different content should not be suppressed",
			old: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadContent</key>
    <string>Content1</string>
</dict>
</plist>`,
			new: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadContent</key>
    <string>Content2</string>
</dict>
</plist>`,
			payloadValidate: true,
			wantSuppressed:  false,
		},
		{
			name: "Validation disabled should not suppress",
			old: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadContent</key>
    <string>Test</string>
</dict>
</plist>`,
			new: `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadContent</key>
    <string>Test</string>
</dict>
</plist>`,
			payloadValidate: false,
			wantSuppressed:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new ResourceData with the payloadValidate field
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"payload_validate": {
					Type:     schema.TypeBool,
					Optional: true,
				},
			}, map[string]interface{}{
				"payload_validate": tt.payloadValidate,
			})

			got := DiffSuppressPayloads("payloads", tt.old, tt.new, d)
			assert.Equal(t, tt.wantSuppressed, got,
				"DiffSuppressPayloads() with payloadValidate=%v returned %v, want %v",
				tt.payloadValidate, got, tt.wantSuppressed)
		})
	}
}
