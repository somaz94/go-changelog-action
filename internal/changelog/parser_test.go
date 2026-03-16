package changelog

import (
	"testing"
)

func TestParseConventionalCommit(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		body     string
		expected *ConventionalCommit
	}{
		{
			name:    "simple feat",
			message: "feat: add user authentication",
			body:    "",
			expected: &ConventionalCommit{
				Type:        "feat",
				Description: "add user authentication",
			},
		},
		{
			name:    "feat with scope",
			message: "feat(auth): add OAuth2 support",
			body:    "",
			expected: &ConventionalCommit{
				Type:        "feat",
				Scope:       "auth",
				Description: "add OAuth2 support",
			},
		},
		{
			name:    "breaking change with !",
			message: "feat!: remove deprecated API",
			body:    "",
			expected: &ConventionalCommit{
				Type:        "feat",
				Description: "remove deprecated API",
				Breaking:    true,
			},
		},
		{
			name:    "breaking change in body",
			message: "feat: update API",
			body:    "BREAKING CHANGE: removed v1 endpoints",
			expected: &ConventionalCommit{
				Type:        "feat",
				Description: "update API",
				Breaking:    true,
			},
		},
		{
			name:    "fix commit",
			message: "fix: resolve null pointer exception",
			body:    "",
			expected: &ConventionalCommit{
				Type:        "fix",
				Description: "resolve null pointer exception",
			},
		},
		{
			name:     "non-conventional commit",
			message:  "Update README",
			body:     "",
			expected: nil,
		},
		{
			name:     "merge commit",
			message:  "Merge pull request #123",
			body:     "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseConventionalCommit(tt.message, tt.body, "abc123", "testuser")
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %+v", result)
				}
				return
			}
			if result == nil {
				t.Fatal("expected non-nil result")
			}
			if result.Type != tt.expected.Type {
				t.Errorf("Type: expected %q, got %q", tt.expected.Type, result.Type)
			}
			if result.Scope != tt.expected.Scope {
				t.Errorf("Scope: expected %q, got %q", tt.expected.Scope, result.Scope)
			}
			if result.Description != tt.expected.Description {
				t.Errorf("Description: expected %q, got %q", tt.expected.Description, result.Description)
			}
			if result.Breaking != tt.expected.Breaking {
				t.Errorf("Breaking: expected %v, got %v", tt.expected.Breaking, result.Breaking)
			}
		})
	}
}

func TestPRAndIssueDetection(t *testing.T) {
	// PR number in description
	cc := ParseConventionalCommit("feat: add login (#42)", "", "abc123", "user")
	if cc == nil {
		t.Fatal("expected non-nil")
	}
	if len(cc.PRNumbers) != 1 || cc.PRNumbers[0] != "42" {
		t.Errorf("expected PR #42, got %v", cc.PRNumbers)
	}

	// Issue reference in body
	cc2 := ParseConventionalCommit("fix: resolve crash", "closes #99\nfixes #100", "def456", "user")
	if cc2 == nil {
		t.Fatal("expected non-nil")
	}
	if len(cc2.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d: %v", len(cc2.Issues), cc2.Issues)
	}

	// Multiple PRs
	cc3 := ParseConventionalCommit("feat: big feature (#10) (#20)", "", "ghi789", "user")
	if cc3 == nil {
		t.Fatal("expected non-nil")
	}
	if len(cc3.PRNumbers) != 2 {
		t.Errorf("expected 2 PRs, got %d: %v", len(cc3.PRNumbers), cc3.PRNumbers)
	}
}

func TestParseNonConventionalCommit(t *testing.T) {
	cc := ParseNonConventionalCommit("Update README (#5)", "fixes #10", "abc123", "user")
	if cc.Type != "other" {
		t.Errorf("expected type 'other', got %q", cc.Type)
	}
	if cc.Description != "Update README (#5)" {
		t.Errorf("expected full message as description, got %q", cc.Description)
	}
	if len(cc.PRNumbers) != 1 || cc.PRNumbers[0] != "5" {
		t.Errorf("expected PR #5, got %v", cc.PRNumbers)
	}
	if len(cc.Issues) != 1 || cc.Issues[0] != "10" {
		t.Errorf("expected issue #10, got %v", cc.Issues)
	}
}

func TestTypeDisplayName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"feat", "Features"},
		{"fix", "Bug Fixes"},
		{"docs", "Documentation"},
		{"ci", "Continuous Integration"},
		{"chore", "Chores"},
		{"other", "Other Changes"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := TypeDisplayName(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTypeDisplayNameWithCustom(t *testing.T) {
	custom := map[string]string{
		"feat": "New Features",
		"fix":  "Bugfixes",
	}

	if got := TypeDisplayNameWithCustom("feat", custom); got != "New Features" {
		t.Errorf("expected 'New Features', got %q", got)
	}
	if got := TypeDisplayNameWithCustom("fix", custom); got != "Bugfixes" {
		t.Errorf("expected 'Bugfixes', got %q", got)
	}
	// Falls back to default
	if got := TypeDisplayNameWithCustom("docs", custom); got != "Documentation" {
		t.Errorf("expected 'Documentation', got %q", got)
	}
}
