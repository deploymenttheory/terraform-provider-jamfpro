package cloud_distribution_point

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// baseModel returns a model with every attribute null. Tests override only the
// fields they exercise.
func baseModel() cloudDistributionPointResourceModel {
	return cloudDistributionPointResourceModel{
		CdnType:                 types.StringNull(),
		Master:                  types.BoolNull(),
		Username:                types.StringNull(),
		Password:                types.StringNull(),
		Directory:               types.StringNull(),
		UploadURL:               types.StringNull(),
		DownloadURL:             types.StringNull(),
		SecondaryAuthRequired:   types.BoolNull(),
		SecondaryAuthStatusCode: types.Int64Null(),
		SecondaryAuthTimeToLive: types.Int64Null(),
		RequireSignedUrls:       types.BoolNull(),
		KeyPairID:               types.StringNull(),
		ExpirationSeconds:       types.Int64Null(),
		PrivateKey:              types.StringNull(),
	}
}

func TestValidateCloudDistributionPointPlan(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(m *cloudDistributionPointResourceModel)
		wantError bool
		// wantSummary, when set, must appear in at least one error summary.
		wantSummary string
	}{
		{
			// Regression: cdn_type sourced from a variable / for_each is unknown
			// at config-validation time and must not error. See issue #1110.
			name: "unknown cdn_type is deferred, not rejected",
			mutate: func(m *cloudDistributionPointResourceModel) {
				m.CdnType = types.StringUnknown()
				m.Master = types.BoolValue(true)
			},
			wantError: false,
		},
		{
			name: "unknown master is deferred, not rejected",
			mutate: func(m *cloudDistributionPointResourceModel) {
				m.CdnType = types.StringValue(cdnTypeJamfCloud)
				m.Master = types.BoolUnknown()
			},
			wantError: false,
		},
		{
			name: "null cdn_type is a config error",
			mutate: func(m *cloudDistributionPointResourceModel) {
				m.CdnType = types.StringNull()
				m.Master = types.BoolValue(true)
			},
			wantError:   true,
			wantSummary: "Missing CDN Type",
		},
		{
			name: "null master is a config error",
			mutate: func(m *cloudDistributionPointResourceModel) {
				m.CdnType = types.StringValue(cdnTypeJamfCloud)
				m.Master = types.BoolNull()
			},
			wantError:   true,
			wantSummary: "Missing Master Flag",
		},
		{
			name: "valid jamf cloud config passes",
			mutate: func(m *cloudDistributionPointResourceModel) {
				m.CdnType = types.StringValue(cdnTypeJamfCloud)
				m.Master = types.BoolValue(true)
			},
			wantError: false,
		},
		{
			// Akamai requires username/password/etc., but when those are unknown
			// (variable-driven) validation is deferred rather than failing.
			name: "akamai with unknown required attributes is deferred",
			mutate: func(m *cloudDistributionPointResourceModel) {
				m.CdnType = types.StringValue(cdnTypeAkamai)
				m.Master = types.BoolValue(true)
				m.Username = types.StringUnknown()
				m.Password = types.StringUnknown()
				m.Directory = types.StringUnknown()
				m.UploadURL = types.StringUnknown()
				m.DownloadURL = types.StringUnknown()
				m.SecondaryAuthRequired = types.BoolUnknown()
			},
			wantError: false,
		},
		{
			// Known-but-missing required attributes for a CDN type still error,
			// proving validation is not wholesale skipped.
			name: "akamai with known-missing required attributes still errors",
			mutate: func(m *cloudDistributionPointResourceModel) {
				m.CdnType = types.StringValue(cdnTypeAkamai)
				m.Master = types.BoolValue(true)
			},
			wantError:   true,
			wantSummary: "Missing username",
		},
		{
			// Attributes invalid for the chosen CDN type are still rejected when
			// they carry a known value.
			name: "attribute not valid for cdn type errors",
			mutate: func(m *cloudDistributionPointResourceModel) {
				m.CdnType = types.StringValue(cdnTypeJamfCloud)
				m.Master = types.BoolValue(true)
				m.Username = types.StringValue("someone")
			},
			wantError:   true,
			wantSummary: "Attribute username Is Not Valid For CDN Type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := baseModel()
			tt.mutate(&m)

			diags := validateCloudDistributionPointPlan(&m)

			if got := diags.HasError(); got != tt.wantError {
				t.Fatalf("HasError() = %v, want %v (diags: %v)", got, tt.wantError, diags)
			}

			if tt.wantSummary != "" {
				found := false
				for _, d := range diags.Errors() {
					if strings.Contains(d.Summary(), tt.wantSummary) {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("expected an error summary containing %q, got %v", tt.wantSummary, diags)
				}
			}
		})
	}
}
