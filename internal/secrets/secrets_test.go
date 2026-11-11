package secrets

import (
	"testing"
)

func TestCalculateShannonEntropy(t *testing.T) {
	lowEntropy := "aaaaaaaaaaaaaaaa"
	if e := CalculateShannonEntropy(lowEntropy); e != 0.0 {
		t.Errorf("expected 0 entropy for uniform string, got %f", e)
	}

	highEntropy := "8fA7z9qLm2Kx1Vp0Q"
	if e := CalculateShannonEntropy(highEntropy); e < 3.5 {
		t.Errorf("expected high entropy (>3.5), got %f", e)
	}
}

func TestScanString(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		shouldMatch bool
	}{
		{
			input:       "Bearer ghp_1234567890abcdef1234567890abcdef1234",
			expected:    "GitHub Personal Access Token",
			shouldMatch: true,
		},
		{
			input:       "sk-proj-abcde12345ABCDE12345abcde12345ABCDE12345",
			expected:    "OpenAI API Key",
			shouldMatch: true,
		},
		{
			input:       "https://hooks.slack.com/services/T12345/B12345/abcde12345",
			expected:    "Slack Incoming Webhook",
			shouldMatch: true,
		},
		{
			input:       "{{ $env.OPENAI_API_KEY }}",
			expected:    "",
			shouldMatch: false,
		},
	}

	for _, tc := range tests {
		matches := ScanString(tc.input)
		if tc.shouldMatch {
			if len(matches) == 0 {
				t.Errorf("expected match for '%s', got 0", tc.input)
				continue
			}
			if matches[0].PatternName != tc.expected {
				t.Errorf("expected pattern '%s', got '%s'", tc.expected, matches[0].PatternName)
			}
		} else {
			if len(matches) > 0 {
				t.Errorf("expected no match for '%s', got %d", tc.input, len(matches))
			}
		}
	}
}

func TestMaskSecret(t *testing.T) {
	masked := MaskSecret("sk-proj-1234567890abcdef")
	if masked != "sk-p...cdef" {
		t.Errorf("unexpected masked output: %s", masked)
	}
}
