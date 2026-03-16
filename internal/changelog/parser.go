package changelog

import (
	"regexp"
	"strings"
	"unicode"
)

// ConventionalCommit represents a parsed conventional commit.
type ConventionalCommit struct {
	Type        string
	Scope       string
	Description string
	Body        string
	Breaking    bool
	Hash        string
	Author      string
	PRNumbers   []string
	Issues      []string
}

var (
	conventionalRegex = regexp.MustCompile(`^(\w+)(?:\(([^)]*)\))?(!)?:\s*(.+)$`)
	prRegex           = regexp.MustCompile(`\(#(\d+)\)`)
	issueRefRegex     = regexp.MustCompile(`(?i)(?:close[sd]?|fix(?:e[sd])?|resolve[sd]?)\s+#(\d+)`)
	issueHashRegex    = regexp.MustCompile(`#(\d+)`)
)

// ParseConventionalCommit parses a commit message into a ConventionalCommit.
// Returns nil if the message does not follow conventional commit format.
func ParseConventionalCommit(message, body, hash, author string) *ConventionalCommit {
	matches := conventionalRegex.FindStringSubmatch(message)
	if matches == nil {
		return nil
	}

	cc := &ConventionalCommit{
		Type:        matches[1],
		Scope:       matches[2],
		Description: matches[4],
		Body:        body,
		Hash:        hash,
		Author:      author,
		Breaking:    matches[3] == "!",
	}

	// Check for BREAKING CHANGE in body
	if !cc.Breaking && strings.Contains(body, "BREAKING CHANGE") {
		cc.Breaking = true
	}

	// Extract PR numbers from description like (#123)
	for _, m := range prRegex.FindAllStringSubmatch(message, -1) {
		cc.PRNumbers = append(cc.PRNumbers, m[1])
	}

	// Extract issue references from body (closes #123, fixes #456, etc.)
	fullText := message + "\n" + body
	for _, m := range issueRefRegex.FindAllStringSubmatch(fullText, -1) {
		cc.Issues = append(cc.Issues, m[1])
	}

	return cc
}

// ParseNonConventionalCommit creates a commit entry for non-conventional messages.
func ParseNonConventionalCommit(message, body, hash, author string) *ConventionalCommit {
	cc := &ConventionalCommit{
		Type:        "other",
		Description: message,
		Body:        body,
		Hash:        hash,
		Author:      author,
	}

	for _, m := range prRegex.FindAllStringSubmatch(message, -1) {
		cc.PRNumbers = append(cc.PRNumbers, m[1])
	}

	fullText := message + "\n" + body
	for _, m := range issueRefRegex.FindAllStringSubmatch(fullText, -1) {
		cc.Issues = append(cc.Issues, m[1])
	}

	return cc
}

// DefaultTypeNames returns the default mapping of commit types to display names.
func DefaultTypeNames() map[string]string {
	return map[string]string{
		"feat":     "Features",
		"fix":      "Bug Fixes",
		"docs":     "Documentation",
		"style":    "Styles",
		"refactor": "Code Refactoring",
		"perf":     "Performance Improvements",
		"test":     "Tests",
		"build":    "Builds",
		"ci":       "Continuous Integration",
		"chore":    "Chores",
		"revert":   "Reverts",
		"other":    "Other Changes",
	}
}

// TypeDisplayNameWithCustom returns the display name using custom mappings first, then defaults.
func TypeDisplayNameWithCustom(commitType string, customMapping map[string]string) string {
	if customMapping != nil {
		if name, ok := customMapping[commitType]; ok {
			return name
		}
	}
	defaults := DefaultTypeNames()
	if name, ok := defaults[commitType]; ok {
		return name
	}
	if commitType == "" {
		return ""
	}
	runes := []rune(commitType)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// TypeDisplayName returns the display name for a commit type (default mapping).
func TypeDisplayName(commitType string) string {
	return TypeDisplayNameWithCustom(commitType, nil)
}
