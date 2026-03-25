package changelog

import (
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/somaz94/go-changelog-action/internal/git"
)

func mockGitRunner(output []byte, err error) func() {
	original := git.RunCommand
	git.RunCommand = func(args ...string) ([]byte, error) {
		return output, err
	}
	return func() { git.RunCommand = original }
}

// mockGitRunnerFunc allows different responses based on command args.
func mockGitRunnerFunc(fn func(args ...string) ([]byte, error)) func() {
	original := git.RunCommand
	git.RunCommand = fn
	return func() { git.RunCommand = original }
}

func TestBuildEntry(t *testing.T) {
	commits := []git.Commit{
		{Hash: "abc1234567890", Message: "feat: add new feature", Date: time.Now(), Author: "alice"},
		{Hash: "def1234567890", Message: "fix: resolve bug", Date: time.Now(), Author: "bob"},
		{Hash: "ghi1234567890", Message: "Merge pull request #1", Date: time.Now(), Author: "alice"},
		{Hash: "jkl1234567890", Message: "chore: update deps", Date: time.Now(), Author: "charlie"},
	}

	skipRegex := regexp.MustCompile("^Merge")
	excludeSet := map[string]bool{"chore": true}

	entry := buildEntry("v1.0.0", "v0.9.0", time.Now(), commits, skipRegex, excludeSet, true, false, nil, nil)

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

	entry := buildEntry("v2.0.0", "v1.0.0", time.Now(), commits, nil, nil, true, false, nil, nil)

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
	entry := buildEntry("v1.0.0", "", time.Now(), commits, nil, nil, true, false, nil, nil)
	if len(entry.Sections) != 1 {
		t.Errorf("expected 1 section without non-conventional, got %d", len(entry.Sections))
	}

	// With includeNonConventional
	entry2 := buildEntry("v1.0.0", "", time.Now(), commits, nil, nil, true, true, nil, nil)
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
	entry := buildEntry("v1.0.0", "", time.Now(), commits, nil, nil, true, false, customMapping, nil)

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

func TestHasSections(t *testing.T) {
	// Empty entry
	empty := Entry{Sections: make(map[string][]ConventionalCommit)}
	if hasSections(empty) {
		t.Error("expected hasSections=false for empty entry")
	}

	// Entry with sections
	withSections := Entry{
		Sections: map[string][]ConventionalCommit{
			"Features": {{Type: "feat", Description: "test"}},
		},
	}
	if !hasSections(withSections) {
		t.Error("expected hasSections=true for entry with sections")
	}

	// Entry with only breaking changes
	withBreaking := Entry{
		Sections: make(map[string][]ConventionalCommit),
		Breaking: []ConventionalCommit{{Type: "feat", Description: "break", Breaking: true}},
	}
	if !hasSections(withBreaking) {
		t.Error("expected hasSections=true for entry with breaking changes")
	}
}

func TestRenderMarkdownUnreleased(t *testing.T) {
	entries := []Entry{
		{
			Version:     "Unreleased",
			PrevVersion: "",
			Date:        time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
			Sections: map[string][]ConventionalCommit{
				"Features": {
					{Type: "feat", Description: "new thing", Hash: "abc1234567890", Author: "test"},
				},
			},
		},
	}

	content := renderMarkdown(entries, "# Changelog", "2006-01-02", "https://github.com/owner/repo")

	// Unreleased should not have a compare link
	if strings.Contains(content, "compare/") {
		t.Error("unreleased should not have compare link")
	}
	if !strings.Contains(content, "## Unreleased") {
		t.Error("expected Unreleased header")
	}
}

func TestRenderMarkdownFirstVersionNoPrev(t *testing.T) {
	// When PrevVersion is empty, first version should link to release tag
	entries := []Entry{
		{
			Version:     "v1.0.0",
			PrevVersion: "",
			Date:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Sections: map[string][]ConventionalCommit{
				"Features": {
					{Type: "feat", Description: "initial", Hash: "abc1234567890", Author: "test"},
				},
			},
		},
	}

	content := renderMarkdown(entries, "# Changelog", "2006-01-02", "https://github.com/owner/repo")

	// First version with no PrevVersion should link to release tag
	if !strings.Contains(content, "releases/tag/v1.0.0") {
		t.Error("expected release tag link for first version")
	}
	if strings.Contains(content, "compare/") {
		t.Error("first version should not have compare link")
	}
}

func TestRenderMarkdownBreakingChanges(t *testing.T) {
	entries := []Entry{
		{
			Version:     "v2.0.0",
			PrevVersion: "v1.0.0",
			Date:        time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			Sections: map[string][]ConventionalCommit{
				"Features": {
					{Type: "feat", Description: "breaking thing", Hash: "abc1234567890", Author: "test", Breaking: true},
				},
			},
			Breaking: []ConventionalCommit{
				{Type: "feat", Description: "breaking thing", Hash: "abc1234567890", Author: "test", Breaking: true},
			},
		},
	}

	content := renderMarkdown(entries, "# Changelog", "2006-01-02", "https://github.com/owner/repo")

	if !strings.Contains(content, "### BREAKING CHANGES") {
		t.Error("expected BREAKING CHANGES section")
	}
}

func TestRenderMarkdownCustomSections(t *testing.T) {
	entries := []Entry{
		{
			Version:     "v1.0.0",
			PrevVersion: "v0.9.0",
			Date:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Sections: map[string][]ConventionalCommit{
				"Features": {
					{Type: "feat", Description: "feature", Hash: "abc1234567890", Author: "test"},
				},
				"Custom Section": {
					{Type: "custom", Description: "custom thing", Hash: "def1234567890", Author: "test"},
				},
			},
		},
	}

	content := renderMarkdown(entries, "# Changelog", "2006-01-02", "https://github.com/owner/repo")

	if !strings.Contains(content, "### Custom Section") {
		t.Error("expected Custom Section")
	}
}

func TestRenderMarkdownNoRepoURL(t *testing.T) {
	entries := []Entry{
		{
			Version:     "v1.0.0",
			PrevVersion: "v0.9.0",
			Date:        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Sections: map[string][]ConventionalCommit{
				"Features": {
					{Type: "feat", Description: "feature", Hash: "abc1234567890", Author: "test",
						PRNumbers: []string{"42"}, Issues: []string{"99"}},
				},
			},
		},
	}

	content := renderMarkdown(entries, "# Changelog", "2006-01-02", "")

	// Without repo URL, should just show hash without link
	if strings.Contains(content, "github.com") {
		t.Error("should not contain github links without repo URL")
	}
	if !strings.Contains(content, "(abc1234)") {
		t.Error("expected short hash without link")
	}
}

func TestRenderCommitLineWithScope(t *testing.T) {
	var sb strings.Builder
	cc := ConventionalCommit{
		Type:        "fix",
		Scope:       "core",
		Description: "fix issue",
		Hash:        "abc1234567890",
		Author:      "test",
	}
	renderCommitLine(&sb, cc, "https://github.com/owner/repo")
	output := sb.String()
	if !strings.Contains(output, "**core:**") {
		t.Errorf("expected scope in output, got %q", output)
	}
}

func TestRenderCommitLineShortHash(t *testing.T) {
	var sb strings.Builder
	// Hash shorter than 7 chars
	cc := ConventionalCommit{
		Type:        "fix",
		Description: "fix",
		Hash:        "abc",
		Author:      "test",
	}
	renderCommitLine(&sb, cc, "")
	output := sb.String()
	if !strings.Contains(output, "(abc)") {
		t.Errorf("expected short hash abc, got %q", output)
	}
}

func TestBuildEntryBreakingDisabled(t *testing.T) {
	commits := []git.Commit{
		{Hash: "abc1234567890", Message: "feat!: breaking change", Date: time.Now(), Author: "test"},
	}

	entry := buildEntry("v2.0.0", "v1.0.0", time.Now(), commits, nil, nil, false, false, nil, nil)

	if len(entry.Breaking) != 0 {
		t.Errorf("expected 0 breaking changes when disabled, got %d", len(entry.Breaking))
	}
}

func TestGenerate(t *testing.T) {
	callCount := 0
	restore := mockGitRunnerFunc(func(args ...string) ([]byte, error) {
		callCount++
		if len(args) > 0 && args[0] == "tag" {
			// Return two tags
			return []byte("v2.0.0|abc1234|2024-02-01T00:00:00Z\nv1.0.0|def5678|2024-01-01T00:00:00Z\n"), nil
		}
		if len(args) > 0 && args[0] == "remote" {
			return []byte("https://github.com/owner/repo.git\n"), nil
		}
		if len(args) > 0 && args[0] == "log" {
			// Return a commit
			return []byte("aaa111\x01feat: new feature\x012024-01-15T10:00:00Z\x01alice\x01\x00"), nil
		}
		return []byte(""), nil
	})
	defer restore()

	cfg := GeneratorConfig{
		TagPattern:      "v[0-9]*.[0-9]*.[0-9]*",
		IncludeBreaking: true,
		DateFormat:      "2006-01-02",
		Header:          "# Changelog",
		Unreleased:      true,
		UnreleasedTitle: "Unreleased",
	}

	result, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.EntriesCount == 0 {
		t.Error("expected at least 1 entry")
	}
	if result.LatestVersion != "v2.0.0" {
		t.Errorf("expected latest version v2.0.0, got %s", result.LatestVersion)
	}
	if !strings.Contains(result.Content, "# Changelog") {
		t.Error("expected header in content")
	}
}

func TestGenerateNoTags(t *testing.T) {
	restore := mockGitRunnerFunc(func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "tag" {
			return []byte(""), nil
		}
		if len(args) > 0 && args[0] == "remote" {
			return []byte("https://github.com/owner/repo\n"), nil
		}
		if len(args) > 0 && args[0] == "log" {
			return []byte("aaa111\x01feat: initial\x012024-01-15T10:00:00Z\x01alice\x01\x00"), nil
		}
		return []byte(""), nil
	})
	defer restore()

	cfg := GeneratorConfig{
		TagPattern:      "v*",
		IncludeBreaking: true,
		DateFormat:      "2006-01-02",
		Header:          "# Changelog",
		Unreleased:      true,
		UnreleasedTitle: "Unreleased",
	}

	result, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.LatestVersion != "" {
		t.Errorf("expected empty latest version, got %s", result.LatestVersion)
	}
}

func TestGenerateError(t *testing.T) {
	restore := mockGitRunner(nil, errors.New("git error"))
	defer restore()

	cfg := GeneratorConfig{
		TagPattern: "v*",
		DateFormat: "2006-01-02",
		Header:     "# Changelog",
	}

	_, err := Generate(cfg)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGenerateInvalidSkipPattern(t *testing.T) {
	restore := mockGitRunnerFunc(func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "tag" {
			return []byte(""), nil
		}
		return []byte(""), nil
	})
	defer restore()

	cfg := GeneratorConfig{
		TagPattern:  "v*",
		DateFormat:  "2006-01-02",
		Header:      "# Changelog",
		SkipCommits: "[invalid",
	}

	_, err := Generate(cfg)
	if err == nil {
		t.Fatal("expected error for invalid skip_commits pattern")
	}
}

func TestGenerateWithSinceUntil(t *testing.T) {
	restore := mockGitRunnerFunc(func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "tag" {
			return []byte("v3.0.0|aaa|2024-03-01T00:00:00Z\nv2.0.0|bbb|2024-02-01T00:00:00Z\nv1.0.0|ccc|2024-01-01T00:00:00Z\n"), nil
		}
		if len(args) > 0 && args[0] == "remote" {
			return []byte("https://github.com/owner/repo\n"), nil
		}
		if len(args) > 0 && args[0] == "log" {
			return []byte("aaa111\x01feat: change\x012024-02-15T10:00:00Z\x01alice\x01\x00"), nil
		}
		return []byte(""), nil
	})
	defer restore()

	cfg := GeneratorConfig{
		TagPattern:      "v[0-9]*.[0-9]*.[0-9]*",
		IncludeBreaking: true,
		DateFormat:      "2006-01-02",
		Header:          "# Changelog",
		SinceTag:        "v1.0.0",
		UntilTag:        "v2.0.0",
	}

	result, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.LatestVersion != "v2.0.0" {
		t.Errorf("expected latest version v2.0.0, got %s", result.LatestVersion)
	}
}

func TestFilterTagsInvalidRange(t *testing.T) {
	tags := []git.Tag{
		{Name: "v3.0.0"},
		{Name: "v2.0.0"},
		{Name: "v1.0.0"},
	}

	// until=v1.0.0 (index 2), since=v3.0.0 (index 0+1=1) → startIdx(2) >= endIdx(1) → nil
	result := filterTags(tags, "v3.0.0", "v1.0.0")
	if result != nil {
		t.Errorf("expected nil for invalid range, got %v", tagNames(result))
	}
}

func TestGenerateWithExcludeTypes(t *testing.T) {
	restore := mockGitRunnerFunc(func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "tag" {
			return []byte("v1.0.0|abc1234|2024-01-01T00:00:00Z\n"), nil
		}
		if len(args) > 0 && args[0] == "remote" {
			return []byte("https://github.com/owner/repo\n"), nil
		}
		if len(args) > 0 && args[0] == "log" {
			return []byte("aaa111\x01feat: feature\x012024-01-15T10:00:00Z\x01alice\x01\x00bbb222\x01chore: cleanup\x012024-01-14T10:00:00Z\x01bob\x01\x00"), nil
		}
		return []byte(""), nil
	})
	defer restore()

	cfg := GeneratorConfig{
		TagPattern:      "v*",
		ExcludeTypes:    []string{"chore"},
		IncludeBreaking: true,
		DateFormat:      "2006-01-02",
		Header:          "# Changelog",
	}

	result, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result.Content, "Chores") {
		t.Error("expected chore commits to be excluded")
	}
}

func TestGenerateRemoteURLError(t *testing.T) {
	restore := mockGitRunnerFunc(func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "tag" {
			return []byte("v1.0.0|abc1234|2024-01-01T00:00:00Z\n"), nil
		}
		if len(args) > 0 && args[0] == "remote" {
			return nil, errors.New("no remote")
		}
		if len(args) > 0 && args[0] == "log" {
			return []byte("aaa111\x01feat: feature\x012024-01-15T10:00:00Z\x01alice\x01\x00"), nil
		}
		return []byte(""), nil
	})
	defer restore()

	cfg := GeneratorConfig{
		TagPattern:      "v*",
		IncludeBreaking: true,
		DateFormat:      "2006-01-02",
		Header:          "# Changelog",
	}

	result, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should still generate content, just without links
	if result.EntriesCount == 0 {
		t.Error("expected entries even without remote URL")
	}
}

func TestGenerateGetCommitsBetweenError(t *testing.T) {
	callCount := 0
	restore := mockGitRunnerFunc(func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "tag" {
			return []byte("v1.0.0|abc1234|2024-01-01T00:00:00Z\n"), nil
		}
		if len(args) > 0 && args[0] == "remote" {
			return []byte("https://github.com/owner/repo\n"), nil
		}
		if len(args) > 0 && args[0] == "log" {
			callCount++
			return nil, errors.New("log error")
		}
		return []byte(""), nil
	})
	defer restore()

	cfg := GeneratorConfig{
		TagPattern:      "v*",
		IncludeBreaking: true,
		DateFormat:      "2006-01-02",
		Header:          "# Changelog",
	}

	result, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should produce result with 0 entries since commits failed
	if result.EntriesCount != 0 {
		t.Errorf("expected 0 entries when commits fail, got %d", result.EntriesCount)
	}
}

func TestIsExcludedAuthor(t *testing.T) {
	excludeList := []string{"GitHub Action", "GitHub Actions", "dependabot[bot]", "renovate[bot]", "github-actions[bot]"}

	tests := []struct {
		author   string
		excluded bool
	}{
		{"GitHub Action", true},
		{"GitHub Actions", true},
		{"dependabot[bot]", true},
		{"renovate[bot]", true},
		{"github-actions[bot]", true},
		{"somaz", false},
		{"alice", false},
	}

	for _, tt := range tests {
		if got := isExcludedAuthor(tt.author, excludeList); got != tt.excluded {
			t.Errorf("isExcludedAuthor(%q) = %v, want %v", tt.author, got, tt.excluded)
		}
	}
}

func TestIsExcludedAuthorCaseInsensitive(t *testing.T) {
	excludeList := []string{"GitHub Actions"}
	if !isExcludedAuthor("github actions", excludeList) {
		t.Error("expected case-insensitive match for 'github actions'")
	}
	if !isExcludedAuthor("GITHUB ACTIONS", excludeList) {
		t.Error("expected case-insensitive match for 'GITHUB ACTIONS'")
	}
}

func TestIsExcludedAuthorWildcard(t *testing.T) {
	excludeList := []string{"bot*"}
	tests := []struct {
		author   string
		excluded bool
	}{
		{"dependabot", true},
		{"renovatebot", true},
		{"mybot", true},
		{"bot", true},
		{"alice", false},
		{"somaz", false},
	}
	for _, tt := range tests {
		if got := isExcludedAuthor(tt.author, excludeList); got != tt.excluded {
			t.Errorf("isExcludedAuthor(%q, [bot*]) = %v, want %v", tt.author, got, tt.excluded)
		}
	}
}

func TestIsExcludedAuthorEmptyList(t *testing.T) {
	if isExcludedAuthor("anyone", nil) {
		t.Error("expected no exclusion with nil list")
	}
	if isExcludedAuthor("anyone", []string{}) {
		t.Error("expected no exclusion with empty list")
	}
}

func TestBuildEntryExcludeAuthors(t *testing.T) {
	commits := []git.Commit{
		{Hash: "abc123", Message: "feat: feature", Date: time.Now(), Author: "somaz"},
		{Hash: "def456", Message: "ci: update workflow", Date: time.Now(), Author: "GitHub Actions"},
		{Hash: "ghi789", Message: "chore(deps): bump deps", Date: time.Now(), Author: "dependabot[bot]"},
	}

	excludeAuthors := []string{"GitHub Actions", "dependabot[bot]"}
	entry := buildEntry("v1.0.0", "", time.Now(), commits, nil, nil, false, false, nil, excludeAuthors)

	if len(entry.Contributors) != 1 || entry.Contributors[0] != "somaz" {
		t.Errorf("expected only [somaz], got %v", entry.Contributors)
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
