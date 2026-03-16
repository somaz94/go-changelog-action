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
