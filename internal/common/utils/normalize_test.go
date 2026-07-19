package utils

import "testing"

func TestNormalizeWhitespace(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "heredoc trailing newline is stripped",
			in:   "Line 1\nLine 2\n",
			want: "Line 1\nLine 2",
		},
		{
			name: "leading and trailing blank lines are stripped",
			in:   "\n\nLine 1\nLine 2\n\n",
			want: "Line 1\nLine 2",
		},
		{
			name: "CRLF line endings normalize to LF",
			in:   "Line 1\r\nLine 2\r\n",
			want: "Line 1\nLine 2",
		},
		{
			name: "bare CR line endings normalize to LF",
			in:   "Line 1\rLine 2",
			want: "Line 1\nLine 2",
		},
		{
			name: "per-line leading/trailing spaces are trimmed",
			in:   "  Line 1  \n  Line 2  ",
			want: "Line 1\nLine 2",
		},
		{
			name: "empty string stays empty",
			in:   "",
			want: "",
		},
		{
			name: "whitespace-only string collapses to empty",
			in:   "   \n\t\n   ",
			want: "",
		},
		{
			name: "single line unaffected",
			in:   "no trailing newline here",
			want: "no trailing newline here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeWhitespace(tt.in)
			if got != tt.want {
				t.Errorf("NormalizeWhitespace(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestNormalizeWhitespace_HeredocDriftEquivalence(t *testing.T) {
	// This is the exact shape from issue #1145: a heredoc string always carries
	// a trailing newline before EOT, but the Jamf Pro API strips it server-side.
	// Old (from API) and new (from HCL heredoc) must normalize equal.
	apiValue := "Line 1\nLine 2"
	heredocValue := "Line 1\nLine 2\n"

	if NormalizeWhitespace(apiValue) != NormalizeWhitespace(heredocValue) {
		t.Errorf("expected heredoc value with trailing newline to normalize equal to API-stripped value: got %q vs %q",
			NormalizeWhitespace(apiValue), NormalizeWhitespace(heredocValue))
	}

	if NormalizeWhitespace("Line 1\nLine 2") == NormalizeWhitespace("Line 1\nLine 3") {
		t.Error("expected genuinely different content to normalize unequal")
	}
}
