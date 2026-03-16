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

// GetTags returns all tags matching the given pattern, sorted by version (descending).
func GetTags(pattern string) ([]Tag, error) {
	out, err := exec.Command("git", "tag", "-l", "--sort=-version:refname",
		"--format=%(refname:short)|%(objectname:short)|%(creatordate:iso8601-strict)").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	re, err := regexp.Compile(convertGlobToRegex(pattern))
	if err != nil {
		return nil, fmt.Errorf("invalid tag pattern %q: %w", pattern, err)
	}

	var tags []Tag
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
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

	out, err := exec.Command("git", "log", rangeArg,
		"--format=%H|%s|%b%x00|%aI|%an").Output()
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
		args = []string{"log", "--format=%H|%s|%b%x00|%aI|%an"}
	} else {
		args = []string{"log", tag + "..HEAD", "--format=%H|%s|%b%x00|%aI|%an"}
	}

	out, err := exec.Command("git", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	return parseCommits(string(out))
}

// GetRemoteURL returns the remote origin URL.
func GetRemoteURL() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get remote URL: %w", err)
	}
	url := strings.TrimSpace(string(out))
	// Convert SSH URL to HTTPS
	if strings.HasPrefix(url, "git@") {
		url = strings.Replace(url, ":", "/", 1)
		url = strings.Replace(url, "git@", "https://", 1)
	}
	url = strings.TrimSuffix(url, ".git")
	return url, nil
}

func parseCommits(output string) ([]Commit, error) {
	var commits []Commit
	// Split by null byte separator
	entries := strings.Split(output, "\x00")
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		// The entry format: body_continuation\n|date|author  OR  hash|subject|body\x00|date|author
		// Find the last occurrence of |date|author pattern
		lines := strings.Split(entry, "\n")
		if len(lines) == 0 {
			continue
		}

		// Last line contains |date|author
		lastLine := lines[len(lines)-1]
		metaParts := strings.SplitN(lastLine, "|", 3)
		if len(metaParts) < 3 {
			continue
		}

		date, _ := time.Parse(time.RFC3339, metaParts[1])
		author := metaParts[2]

		// First line contains hash|subject|body_start
		firstLine := lines[0]
		firstParts := strings.SplitN(firstLine, "|", 3)
		if len(firstParts) < 2 {
			continue
		}

		hash := firstParts[0]
		subject := firstParts[1]
		var bodyParts []string
		if len(firstParts) > 2 {
			bodyParts = append(bodyParts, firstParts[2])
		}
		// Add middle lines to body
		for i := 1; i < len(lines)-1; i++ {
			bodyParts = append(bodyParts, lines[i])
		}
		body := strings.TrimSpace(strings.Join(bodyParts, "\n"))

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
