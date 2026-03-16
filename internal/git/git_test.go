package git

import (
	"regexp"
	"testing"
)

func TestConvertGlobToRegex(t *testing.T) {
	tests := []struct {
		glob     string
		input    string
		expected bool
	}{
		{"v[0-9]*.[0-9]*.[0-9]*", "v1.0.0", true},
		{"v[0-9]*.[0-9]*.[0-9]*", "v12.34.56", true},
		{"v[0-9]*.[0-9]*.[0-9]*", "v1.0.0-rc1", true}, // glob * matches any chars including -rc1
		{"v[0-9]*.[0-9]*.[0-9]*", "release-1.0.0", false},
		{"v*", "v1", true},
		{"v*", "v1.0.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.glob+"_"+tt.input, func(t *testing.T) {
			pattern := convertGlobToRegex(tt.glob)
			re, err := regexp.Compile(pattern)
			if err != nil {
				t.Fatalf("regex error: %v", err)
			}
			matched := re.MatchString(tt.input)
			if matched != tt.expected {
				t.Errorf("glob=%q input=%q: expected %v, got %v (regex=%q)", tt.glob, tt.input, tt.expected, matched, pattern)
			}
		})
	}
}

func TestParseCommits(t *testing.T) {
	// Test empty input
	commits, err := parseCommits("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 0 {
		t.Errorf("expected 0 commits, got %d", len(commits))
	}
}
