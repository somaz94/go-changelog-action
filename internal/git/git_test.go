package git

import (
	"fmt"
	"regexp"
	"testing"
)

// Helper to set a mock command runner and restore it after the test.
func mockRunner(output []byte, err error) func() {
	original := RunCommand
	RunCommand = func(args ...string) ([]byte, error) {
		return output, err
	}
	return func() { RunCommand = original }
}

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

func TestConvertGlobToRegexQuestionMark(t *testing.T) {
	pattern := convertGlobToRegex("v?.0.0")
	re, err := regexp.Compile(pattern)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}
	if !re.MatchString("v1.0.0") {
		t.Error("expected v1.0.0 to match v?.0.0")
	}
	if re.MatchString("v12.0.0") {
		t.Error("expected v12.0.0 NOT to match v?.0.0")
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

func TestParseCommitsSingleCommit(t *testing.T) {
	input := "abc1234567890|feat: add feature|some body\n|2024-01-15T10:00:00Z|alice"
	commits, err := parseCommits(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}
	if commits[0].Hash != "abc1234567890" {
		t.Errorf("expected hash abc1234567890, got %s", commits[0].Hash)
	}
	if commits[0].Message != "feat: add feature" {
		t.Errorf("expected message 'feat: add feature', got %s", commits[0].Message)
	}
	if commits[0].Author != "alice" {
		t.Errorf("expected author alice, got %s", commits[0].Author)
	}
	if commits[0].Body != "some body" {
		t.Errorf("expected body 'some body', got %q", commits[0].Body)
	}
}

func TestParseCommitsMultipleCommits(t *testing.T) {
	input := "abc123|feat: first|\n|2024-01-15T10:00:00Z|alice\x00def456|fix: second|\n|2024-01-16T10:00:00Z|bob"
	commits, err := parseCommits(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(commits))
	}
	if commits[0].Author != "bob" {
		t.Errorf("expected first commit author bob (newer), got %s", commits[0].Author)
	}
	if commits[1].Author != "alice" {
		t.Errorf("expected second commit author alice (older), got %s", commits[1].Author)
	}
}

func TestParseCommitsMultilineBody(t *testing.T) {
	input := "abc123|feat: feature|body line 1\nbody line 2\n|2024-01-15T10:00:00Z|alice"
	commits, err := parseCommits(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}
	if commits[0].Message != "feat: feature" {
		t.Errorf("expected message 'feat: feature', got %q", commits[0].Message)
	}
	if commits[0].Body != "body line 1\nbody line 2" {
		t.Errorf("expected multiline body, got %q", commits[0].Body)
	}
}

func TestParseCommitsNoBody(t *testing.T) {
	input := "abc123|feat: no body|\n|2024-01-15T10:00:00Z|alice"
	commits, err := parseCommits(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}
	if commits[0].Message != "feat: no body" {
		t.Errorf("expected message, got %q", commits[0].Message)
	}
}

func TestParseCommitsMalformedEntry(t *testing.T) {
	input := "malformed_entry_no_pipes"
	commits, err := parseCommits(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 0 {
		t.Errorf("expected 0 commits for malformed input, got %d", len(commits))
	}
}

func TestParseCommitsThreeCommits(t *testing.T) {
	input := "aaa|feat: a|\n|2024-01-01T10:00:00Z|alice" +
		"\x00bbb|fix: b|\n|2024-01-02T10:00:00Z|bob" +
		"\x00ccc|docs: c|\n|2024-01-03T10:00:00Z|charlie"
	commits, err := parseCommits(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 3 {
		t.Fatalf("expected 3 commits, got %d", len(commits))
	}
	if commits[0].Author != "charlie" {
		t.Errorf("expected charlie first, got %s", commits[0].Author)
	}
	if commits[2].Author != "alice" {
		t.Errorf("expected alice last, got %s", commits[2].Author)
	}
}

// --- Tests for parseTags ---

func TestParseTags(t *testing.T) {
	output := "v2.0.0|abc1234|2024-02-01T00:00:00Z\nv1.0.0|def5678|2024-01-01T00:00:00Z\nrelease-1.0|ghi9012|2024-01-15T00:00:00Z"
	tags, err := parseTags(output, "v[0-9]*.[0-9]*.[0-9]*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0].Name != "v2.0.0" {
		t.Errorf("expected v2.0.0, got %s", tags[0].Name)
	}
	if tags[1].Name != "v1.0.0" {
		t.Errorf("expected v1.0.0, got %s", tags[1].Name)
	}
}

func TestParseTagsEmpty(t *testing.T) {
	tags, err := parseTags("", "v*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(tags))
	}
}

func TestParseTagsMalformedLine(t *testing.T) {
	output := "v1.0.0|abc1234\nv2.0.0|def5678|2024-01-01T00:00:00Z"
	tags, err := parseTags(output, "v*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// First line has < 3 parts, should be skipped
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
	if tags[0].Name != "v2.0.0" {
		t.Errorf("expected v2.0.0, got %s", tags[0].Name)
	}
}

// --- Tests for cleanRemoteURL ---

func TestCleanRemoteURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://github.com/owner/repo.git\n", "https://github.com/owner/repo"},
		{"https://github.com/owner/repo\n", "https://github.com/owner/repo"},
		{"git@github.com:owner/repo.git\n", "https://github.com/owner/repo"},
		{"git@github.com:owner/repo\n", "https://github.com/owner/repo"},
		{"  https://github.com/owner/repo.git  ", "https://github.com/owner/repo"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := cleanRemoteURL(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// --- Tests for GetTags with mock ---

func TestGetTags(t *testing.T) {
	restore := mockRunner([]byte("v2.0.0|abc1234|2024-02-01T00:00:00Z\nv1.0.0|def5678|2024-01-01T00:00:00Z\n"), nil)
	defer restore()

	tags, err := GetTags("v[0-9]*.[0-9]*.[0-9]*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
}

func TestGetTagsError(t *testing.T) {
	restore := mockRunner(nil, fmt.Errorf("git error"))
	defer restore()

	_, err := GetTags("v*")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Tests for GetCommitsBetween with mock ---

func TestGetCommitsBetween(t *testing.T) {
	mockOutput := "abc123|feat: feature|\n|2024-01-15T10:00:00Z|alice"
	restore := mockRunner([]byte(mockOutput), nil)
	defer restore()

	commits, err := GetCommitsBetween("v1.0.0", "v2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}
}

func TestGetCommitsBetweenEmptyFrom(t *testing.T) {
	mockOutput := "abc123|feat: initial|\n|2024-01-15T10:00:00Z|alice"
	restore := mockRunner([]byte(mockOutput), nil)
	defer restore()

	commits, err := GetCommitsBetween("", "v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}
}

func TestGetCommitsBetweenError(t *testing.T) {
	restore := mockRunner(nil, fmt.Errorf("git error"))
	defer restore()

	_, err := GetCommitsBetween("v1", "v2")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Tests for GetCommitsSinceTag with mock ---

func TestGetCommitsSinceTag(t *testing.T) {
	mockOutput := "abc123|feat: new|\n|2024-01-15T10:00:00Z|alice"
	restore := mockRunner([]byte(mockOutput), nil)
	defer restore()

	commits, err := GetCommitsSinceTag("v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}
}

func TestGetCommitsSinceTagEmpty(t *testing.T) {
	mockOutput := "abc123|feat: all|\n|2024-01-15T10:00:00Z|alice"
	restore := mockRunner([]byte(mockOutput), nil)
	defer restore()

	commits, err := GetCommitsSinceTag("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}
}

func TestGetCommitsSinceTagError(t *testing.T) {
	restore := mockRunner(nil, fmt.Errorf("git error"))
	defer restore()

	_, err := GetCommitsSinceTag("v1.0.0")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Tests for GetRemoteURL with mock ---

func TestGetRemoteURL(t *testing.T) {
	restore := mockRunner([]byte("https://github.com/owner/repo.git\n"), nil)
	defer restore()

	url, err := GetRemoteURL()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url != "https://github.com/owner/repo" {
		t.Errorf("expected https://github.com/owner/repo, got %s", url)
	}
}

func TestGetRemoteURLSSH(t *testing.T) {
	restore := mockRunner([]byte("git@github.com:owner/repo.git\n"), nil)
	defer restore()

	url, err := GetRemoteURL()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url != "https://github.com/owner/repo" {
		t.Errorf("expected https://github.com/owner/repo, got %s", url)
	}
}

func TestGetRemoteURLError(t *testing.T) {
	restore := mockRunner(nil, fmt.Errorf("git error"))
	defer restore()

	_, err := GetRemoteURL()
	if err == nil {
		t.Fatal("expected error")
	}
}
