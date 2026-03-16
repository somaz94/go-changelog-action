package changelog

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/somaz94/go-changelog-action/internal/git"
)

func TestBuildEntry(t *testing.T) {
	commits := []git.Commit{
		{Hash: "abc1234567890", Message: "feat: add new feature", Date: time.Now(), Author: "test"},
		{Hash: "def1234567890", Message: "fix: resolve bug", Date: time.Now(), Author: "test"},
		{Hash: "ghi1234567890", Message: "Merge pull request #1", Date: time.Now(), Author: "test"},
		{Hash: "jkl1234567890", Message: "chore: update deps", Date: time.Now(), Author: "test"},
	}

	skipRegex := regexp.MustCompile("^Merge")
	excludeSet := map[string]bool{"chore": true}

	entry := buildEntry("v1.0.0", time.Now(), commits, skipRegex, excludeSet, true)

	if entry.Version != "v1.0.0" {
		t.Errorf("expected version v1.0.0, got %s", entry.Version)
	}

	if len(entry.Sections) != 2 {
		t.Errorf("expected 2 sections, got %d: %v", len(entry.Sections), getSectionNames(entry))
	}
	if _, ok := entry.Sections["Features"]; !ok {
		t.Error("expected Features section")
	}
	if _, ok := entry.Sections["Bug Fixes"]; !ok {
		t.Error("expected Bug Fixes section")
	}
}

func TestBuildEntryBreaking(t *testing.T) {
	commits := []git.Commit{
		{Hash: "abc1234567890", Message: "feat!: breaking change", Date: time.Now(), Author: "test"},
	}

	entry := buildEntry("v2.0.0", time.Now(), commits, nil, nil, true)

	if len(entry.Breaking) != 1 {
		t.Errorf("expected 1 breaking change, got %d", len(entry.Breaking))
	}
}

func TestRenderMarkdown(t *testing.T) {
	entries := []Entry{
		{
			Version: "v1.0.0",
			Date:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Sections: map[string][]ConventionalCommit{
				"Features": {
					{Type: "feat", Description: "add login", Hash: "abc1234567890", Author: "test"},
				},
				"Bug Fixes": {
					{Type: "fix", Description: "fix crash", Hash: "def1234567890", Author: "test", Scope: "core"},
				},
			},
		},
	}

	content := renderMarkdown(entries, "# Changelog", "2006-01-02", "https://github.com/somaz94/test")

	if !strings.Contains(content, "# Changelog") {
		t.Error("expected header")
	}
	if !strings.Contains(content, "## [v1.0.0]") {
		t.Error("expected version header with link")
	}
	if !strings.Contains(content, "### Features") {
		t.Error("expected Features section")
	}
	if !strings.Contains(content, "- add login") {
		t.Error("expected feature commit")
	}
	if !strings.Contains(content, "**core:**") {
		t.Error("expected scoped commit")
	}
	if !strings.Contains(content, "2024-01-15") {
		t.Error("expected formatted date")
	}
}

func getSectionNames(entry Entry) []string {
	var names []string
	for name := range entry.Sections {
		names = append(names, name)
	}
	return names
}
