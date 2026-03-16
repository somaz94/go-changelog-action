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
		{Hash: "abc1234567890", Message: "feat: add new feature", Date: time.Now(), Author: "alice"},
		{Hash: "def1234567890", Message: "fix: resolve bug", Date: time.Now(), Author: "bob"},
		{Hash: "ghi1234567890", Message: "Merge pull request #1", Date: time.Now(), Author: "alice"},
		{Hash: "jkl1234567890", Message: "chore: update deps", Date: time.Now(), Author: "charlie"},
	}

	skipRegex := regexp.MustCompile("^Merge")
	excludeSet := map[string]bool{"chore": true}

	entry := buildEntry("v1.0.0", "v0.9.0", time.Now(), commits, skipRegex, excludeSet, true, false, nil)

	if entry.Version != "v1.0.0" {
		t.Errorf("expected version v1.0.0, got %s", entry.Version)
	}
	if entry.PrevVersion != "v0.9.0" {
		t.Errorf("expected prev version v0.9.0, got %s", entry.PrevVersion)
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

	// Contributors should include alice and bob (chore excluded but author still tracked, merge skipped)
	if len(entry.Contributors) < 2 {
		t.Errorf("expected at least 2 contributors, got %d: %v", len(entry.Contributors), entry.Contributors)
	}
}

func TestBuildEntryBreaking(t *testing.T) {
	commits := []git.Commit{
		{Hash: "abc1234567890", Message: "feat!: breaking change", Date: time.Now(), Author: "test"},
	}

	entry := buildEntry("v2.0.0", "v1.0.0", time.Now(), commits, nil, nil, true, false, nil)

	if len(entry.Breaking) != 1 {
		t.Errorf("expected 1 breaking change, got %d", len(entry.Breaking))
	}
}

func TestBuildEntryNonConventional(t *testing.T) {
	commits := []git.Commit{
		{Hash: "abc1234567890", Message: "feat: add feature", Date: time.Now(), Author: "test"},
		{Hash: "def1234567890", Message: "Update README", Date: time.Now(), Author: "test"},
	}

	// Without includeNonConventional
	entry := buildEntry("v1.0.0", "", time.Now(), commits, nil, nil, true, false, nil)
	if len(entry.Sections) != 1 {
		t.Errorf("expected 1 section without non-conventional, got %d", len(entry.Sections))
	}

	// With includeNonConventional
	entry2 := buildEntry("v1.0.0", "", time.Now(), commits, nil, nil, true, true, nil)
	if len(entry2.Sections) != 2 {
		t.Errorf("expected 2 sections with non-conventional, got %d: %v", len(entry2.Sections), getSectionNames(entry2))
	}
	if _, ok := entry2.Sections["Other Changes"]; !ok {
		t.Error("expected Other Changes section")
	}
}

func TestBuildEntryCustomMapping(t *testing.T) {
	commits := []git.Commit{
		{Hash: "abc1234567890", Message: "feat: add feature", Date: time.Now(), Author: "test"},
	}

	customMapping := map[string]string{"feat": "New Features"}
	entry := buildEntry("v1.0.0", "", time.Now(), commits, nil, nil, true, false, customMapping)

	if _, ok := entry.Sections["New Features"]; !ok {
		t.Errorf("expected 'New Features' section with custom mapping, got: %v", getSectionNames(entry))
	}
}

func TestRenderMarkdown(t *testing.T) {
	entries := []Entry{
		{
			Version:     "v1.0.0",
			PrevVersion: "v0.9.0",
			Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Sections: map[string][]ConventionalCommit{
				"Features": {
					{Type: "feat", Description: "add login", Hash: "abc1234567890", Author: "test",
						PRNumbers: []string{"42"}},
				},
				"Bug Fixes": {
					{Type: "fix", Description: "fix crash", Hash: "def1234567890", Author: "test", Scope: "core",
						Issues: []string{"99"}},
				},
			},
			Contributors: []string{"alice", "bob"},
		},
	}

	content := renderMarkdown(entries, "# Changelog", "2006-01-02", "https://github.com/somaz94/test")

	if !strings.Contains(content, "# Changelog") {
		t.Error("expected header")
	}
	// Compare link
	if !strings.Contains(content, "compare/v0.9.0...v1.0.0") {
		t.Error("expected compare link between versions")
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
	// PR link
	if !strings.Contains(content, "[#42]") {
		t.Error("expected PR link")
	}
	// Issue link
	if !strings.Contains(content, "closes") {
		t.Error("expected issue close reference")
	}
	// Contributors
	if !strings.Contains(content, "### Contributors") {
		t.Error("expected Contributors section")
	}
	if !strings.Contains(content, "- alice") {
		t.Error("expected contributor alice")
	}
}

func TestFilterTags(t *testing.T) {
	tags := []git.Tag{
		{Name: "v3.0.0"},
		{Name: "v2.0.0"},
		{Name: "v1.0.0"},
	}

	// No filter
	result := filterTags(tags, "", "")
	if len(result) != 3 {
		t.Errorf("expected 3 tags, got %d", len(result))
	}

	// Since tag
	result = filterTags(tags, "v2.0.0", "")
	if len(result) != 2 {
		t.Errorf("expected 2 tags with since=v2.0.0, got %d: %v", len(result), tagNames(result))
	}

	// Until tag
	result = filterTags(tags, "", "v2.0.0")
	if len(result) != 2 {
		t.Errorf("expected 2 tags with until=v2.0.0, got %d: %v", len(result), tagNames(result))
	}

	// Both
	result = filterTags(tags, "v2.0.0", "v2.0.0")
	if len(result) != 1 || result[0].Name != "v2.0.0" {
		t.Errorf("expected only v2.0.0, got %v", tagNames(result))
	}
}

func getSectionNames(entry Entry) []string {
	var names []string
	for name := range entry.Sections {
		names = append(names, name)
	}
	return names
}

func tagNames(tags []git.Tag) []string {
	var names []string
	for _, t := range tags {
		names = append(names, t.Name)
	}
	return names
}
