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
}

var conventionalRegex = regexp.MustCompile(`^(\w+)(?:\(([^)]*)\))?(!)?:\s*(.+)$`)

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

	return cc
}

// TypeDisplayName returns the display name for a commit type.
func TypeDisplayName(commitType string) string {
	names := map[string]string{
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
	}
	if name, ok := names[commitType]; ok {
		return name
	}
	if commitType == "" {
		return ""
	}
	runes := []rune(commitType)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
