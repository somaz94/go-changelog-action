package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Commit represents a parsed git commit.
type Commit struct {
	Hash    string
	Message string
	Body    string
	Date    time.Time
	Author  string
}

// Tag represents a git tag with its associated commit.
type Tag struct {
	Name string
	Hash string
	Date time.Time
}

// RunCommand executes a git command and returns its output.
// This is a package-level variable so tests can override it.
var RunCommand = func(args ...string) ([]byte, error) {
	return exec.Command("git", args...).Output()
}

// commitFormat uses SOH (\x01) as field separator and NUL (\x00) as record separator.
// Fields: hash, subject, date, author, body
const commitFormat = "%H%x01%s%x01%aI%x01%an%x01%b%x00"

// GetTags returns all tags matching the given pattern, sorted by version (descending).
func GetTags(pattern string) ([]Tag, error) {
	out, err := RunCommand("tag", "-l", "--sort=-version:refname",
		"--format=%(refname:short)|%(objectname:short)|%(creatordate:iso8601-strict)")
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return parseTags(string(out), pattern)
}

// parseTags parses the raw tag output and filters by pattern.
func parseTags(output, pattern string) ([]Tag, error) {
	re, err := regexp.Compile(convertGlobToRegex(pattern))
	if err != nil {
		return nil, fmt.Errorf("invalid tag pattern %q: %w", pattern, err)
	}

	var tags []Tag
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 3)
		if len(parts) < 3 {
			continue
		}
		name := parts[0]
		if !re.MatchString(name) {
			continue
		}
		date, _ := time.Parse(time.RFC3339, parts[2])
		tags = append(tags, Tag{
			Name: name,
			Hash: parts[1],
			Date: date,
		})
	}

	return tags, nil
}

// GetCommitsBetween returns commits between two refs (exclusive from, inclusive to).
// If fromRef is empty, returns all commits up to toRef.
func GetCommitsBetween(fromRef, toRef string) ([]Commit, error) {
	var rangeArg string
	if fromRef == "" {
		rangeArg = toRef
	} else {
		rangeArg = fromRef + ".." + toRef
	}

	out, err := RunCommand("log", rangeArg,
		"--format="+commitFormat)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	return parseCommits(string(out))
}

// GetCommitsSinceTag returns all commits since the given tag.
// If tag is empty, returns all commits.
func GetCommitsSinceTag(tag string) ([]Commit, error) {
	var args []string
	if tag == "" {
		args = []string{"log", "--format=" + commitFormat}
	} else {
		args = []string{"log", tag + "..HEAD", "--format=" + commitFormat}
	}

	out, err := RunCommand(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	return parseCommits(string(out))
}

// GetRemoteURL returns the remote origin URL.
func GetRemoteURL() (string, error) {
	out, err := RunCommand("remote", "get-url", "origin")
	if err != nil {
		return "", fmt.Errorf("failed to get remote URL: %w", err)
	}
	return cleanRemoteURL(string(out)), nil
}

// cleanRemoteURL converts a raw remote URL to a clean HTTPS URL.
func cleanRemoteURL(raw string) string {
	url := strings.TrimSpace(raw)
	// Convert SSH URL to HTTPS
	if strings.HasPrefix(url, "git@") {
		url = strings.Replace(url, ":", "/", 1)
		url = strings.Replace(url, "git@", "https://", 1)
	}
	url = strings.TrimSuffix(url, ".git")
	return url
}

func parseCommits(output string) ([]Commit, error) {
	var commits []Commit
	// Split by NUL byte (record separator)
	records := strings.Split(output, "\x00")
	for _, record := range records {
		record = strings.TrimSpace(record)
		if record == "" {
			continue
		}

		// Split by SOH byte (field separator): hash, subject, date, author, body
		fields := strings.SplitN(record, "\x01", 5)
		if len(fields) < 4 {
			continue
		}

		hash := fields[0]
		subject := fields[1]
		dateStr := fields[2]
		author := fields[3]
		var body string
		if len(fields) > 4 {
			body = strings.TrimSpace(fields[4])
		}

		if hash == "" {
			continue
		}

		date, _ := time.Parse(time.RFC3339, dateStr)

		commits = append(commits, Commit{
			Hash:    hash,
			Message: subject,
			Body:    body,
			Date:    date,
			Author:  author,
		})
	}

	// Sort by date descending
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.After(commits[j].Date)
	})

	return commits, nil
}

// convertGlobToRegex converts a simple glob pattern to a regex.
func convertGlobToRegex(glob string) string {
	result := "^"
	for _, ch := range glob {
		switch ch {
		case '*':
			result += ".*"
		case '?':
			result += "."
		case '.':
			result += "\\."
		case '[', ']':
			result += string(ch)
		default:
			result += string(ch)
		}
	}
	result += "$"
	return result
}
