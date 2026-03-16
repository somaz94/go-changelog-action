package changelog

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/somaz94/go-changelog-action/internal/git"
)

// Entry represents a single version entry in the changelog.
type Entry struct {
	Version  string
	Date     time.Time
	Sections map[string][]ConventionalCommit
	Breaking []ConventionalCommit
}

// GeneratorConfig holds configuration for changelog generation.
type GeneratorConfig struct {
	TagPattern     string
	ExcludeTypes   []string
	IncludeBreaking bool
	DateFormat     string
	Header         string
	Unreleased     bool
	UnreleasedTitle string
	SkipCommits    string
	RepositoryURL  string
}

// Result holds the output of changelog generation.
type Result struct {
	Content      string
	EntriesCount int
	LatestVersion string
}

// Generate creates a changelog from git history.
func Generate(cfg GeneratorConfig) (*Result, error) {
	tags, err := git.GetTags(cfg.TagPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	var skipRegex *regexp.Regexp
	if cfg.SkipCommits != "" {
		skipRegex, err = regexp.Compile(cfg.SkipCommits)
		if err != nil {
			return nil, fmt.Errorf("invalid skip_commits pattern: %w", err)
		}
	}

	excludeSet := make(map[string]bool)
	for _, t := range cfg.ExcludeTypes {
		excludeSet[t] = true
	}

	repoURL := cfg.RepositoryURL
	if repoURL == "" {
		repoURL, _ = git.GetRemoteURL()
	}

	var entries []Entry

	// Generate unreleased section
	if cfg.Unreleased && len(tags) > 0 {
		commits, err := git.GetCommitsSinceTag(tags[0].Name)
		if err == nil && len(commits) > 0 {
			entry := buildEntry(cfg.UnreleasedTitle, time.Now(), commits, skipRegex, excludeSet, cfg.IncludeBreaking)
			if hasSections(entry) {
				entries = append(entries, entry)
			}
		}
	} else if cfg.Unreleased && len(tags) == 0 {
		// No tags yet - all commits are unreleased
		commits, err := git.GetCommitsSinceTag("")
		if err == nil && len(commits) > 0 {
			entry := buildEntry(cfg.UnreleasedTitle, time.Now(), commits, skipRegex, excludeSet, cfg.IncludeBreaking)
			if hasSections(entry) {
				entries = append(entries, entry)
			}
		}
	}

	// Generate entries for each tag
	for i, tag := range tags {
		var fromRef string
		if i+1 < len(tags) {
			fromRef = tags[i+1].Name
		}

		commits, err := git.GetCommitsBetween(fromRef, tag.Name)
		if err != nil {
			continue
		}

		entry := buildEntry(tag.Name, tag.Date, commits, skipRegex, excludeSet, cfg.IncludeBreaking)
		entries = append(entries, entry)
	}

	latestVersion := ""
	if len(tags) > 0 {
		latestVersion = tags[0].Name
	}

	content := renderMarkdown(entries, cfg.Header, cfg.DateFormat, repoURL)

	return &Result{
		Content:       content,
		EntriesCount:  len(entries),
		LatestVersion: latestVersion,
	}, nil
}

func buildEntry(version string, date time.Time, commits []git.Commit, skipRegex *regexp.Regexp, excludeSet map[string]bool, includeBreaking bool) Entry {
	entry := Entry{
		Version:  version,
		Date:     date,
		Sections: make(map[string][]ConventionalCommit),
	}

	for _, commit := range commits {
		if skipRegex != nil && skipRegex.MatchString(commit.Message) {
			continue
		}

		cc := ParseConventionalCommit(commit.Message, commit.Body, commit.Hash, commit.Author)
		if cc == nil {
			continue
		}

		if excludeSet[cc.Type] {
			continue
		}

		if includeBreaking && cc.Breaking {
			entry.Breaking = append(entry.Breaking, *cc)
		}

		typeName := TypeDisplayName(cc.Type)
		entry.Sections[typeName] = append(entry.Sections[typeName], *cc)
	}

	return entry
}

func hasSections(entry Entry) bool {
	return len(entry.Sections) > 0 || len(entry.Breaking) > 0
}

func renderMarkdown(entries []Entry, header, dateFormat, repoURL string) string {
	var sb strings.Builder

	sb.WriteString(header)
	sb.WriteString("\n\n")
	sb.WriteString("All notable changes to this project will be documented in this file.\n\n")

	// Sort section names for consistent output
	sectionOrder := []string{
		"Features", "Bug Fixes", "Performance Improvements",
		"Code Refactoring", "Documentation", "Tests",
		"Builds", "Continuous Integration", "Styles", "Chores", "Reverts",
	}

	for _, entry := range entries {
		// Version header
		dateStr := entry.Date.Format(dateFormat)
		if repoURL != "" && entry.Version != "Unreleased" {
			sb.WriteString(fmt.Sprintf("## [%s](%s/releases/tag/%s) (%s)\n\n", entry.Version, repoURL, entry.Version, dateStr))
		} else {
			sb.WriteString(fmt.Sprintf("## %s (%s)\n\n", entry.Version, dateStr))
		}

		// Breaking changes
		if len(entry.Breaking) > 0 {
			sb.WriteString("### BREAKING CHANGES\n\n")
			for _, cc := range entry.Breaking {
				renderCommitLine(&sb, cc, repoURL)
			}
			sb.WriteString("\n")
		}

		// Regular sections
		rendered := make(map[string]bool)
		for _, sectionName := range sectionOrder {
			commits, ok := entry.Sections[sectionName]
			if !ok {
				continue
			}
			rendered[sectionName] = true
			sb.WriteString(fmt.Sprintf("### %s\n\n", sectionName))
			for _, cc := range commits {
				renderCommitLine(&sb, cc, repoURL)
			}
			sb.WriteString("\n")
		}

		// Remaining sections not in the predefined order
		var remaining []string
		for name := range entry.Sections {
			if !rendered[name] {
				remaining = append(remaining, name)
			}
		}
		sort.Strings(remaining)
		for _, sectionName := range remaining {
			commits := entry.Sections[sectionName]
			sb.WriteString(fmt.Sprintf("### %s\n\n", sectionName))
			for _, cc := range commits {
				renderCommitLine(&sb, cc, repoURL)
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func renderCommitLine(sb *strings.Builder, cc ConventionalCommit, repoURL string) {
	shortHash := cc.Hash
	if len(shortHash) > 7 {
		shortHash = shortHash[:7]
	}

	if cc.Scope != "" {
		if repoURL != "" {
			sb.WriteString(fmt.Sprintf("- **%s:** %s ([%s](%s/commit/%s))\n", cc.Scope, cc.Description, shortHash, repoURL, cc.Hash))
		} else {
			sb.WriteString(fmt.Sprintf("- **%s:** %s (%s)\n", cc.Scope, cc.Description, shortHash))
		}
	} else {
		if repoURL != "" {
			sb.WriteString(fmt.Sprintf("- %s ([%s](%s/commit/%s))\n", cc.Description, shortHash, repoURL, cc.Hash))
		} else {
			sb.WriteString(fmt.Sprintf("- %s (%s)\n", cc.Description, shortHash))
		}
	}
}
