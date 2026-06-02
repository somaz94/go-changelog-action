package changelog

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/somaz94/go-changelog-action/internal/git"
	"github.com/somaz94/go-changelog-action/internal/output"
)

// Entry represents a single version entry in the changelog.
type Entry struct {
	Version      string
	PrevVersion  string
	Date         time.Time
	Sections     map[string][]ConventionalCommit
	Breaking     []ConventionalCommit
	Contributors []string
}

// GeneratorConfig holds configuration for changelog generation.
type GeneratorConfig struct {
	TagPattern             string
	ExcludeTypes           []string
	IncludeBreaking        bool
	DateFormat             string
	Header                 string
	Unreleased             bool
	UnreleasedTitle        string
	SkipCommits            string
	RepositoryURL          string
	IncludeNonConventional bool
	SinceTag               string
	UntilTag               string
	CustomTypeMapping      map[string]string
	ExcludeAuthors         []string
}

// Result holds the output of changelog generation.
type Result struct {
	Content       string
	EntriesCount  int
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
		var err error
		repoURL, err = git.GetRemoteURL()
		if err != nil {
			output.LogWarning(fmt.Sprintf("Failed to detect repository URL: %v. Commit links will be omitted.", err))
		}
	}

	// Apply since/until tag filtering
	tags = filterTags(tags, cfg.SinceTag, cfg.UntilTag)

	var entries []Entry

	// Generate unreleased section
	if cfg.Unreleased {
		var latestTag string
		if len(tags) > 0 {
			latestTag = tags[0].Name
		}
		commits, err := git.GetCommitsSinceTag(latestTag)
		if err == nil && len(commits) > 0 {
			entry := buildEntry(cfg.UnreleasedTitle, "", time.Now(), commits, skipRegex, excludeSet, cfg.IncludeBreaking, cfg.IncludeNonConventional, cfg.CustomTypeMapping, cfg.ExcludeAuthors)
			if hasSections(entry) {
				entries = append(entries, entry)
			}
		}
	}

	// Generate entries for each tag
	for i, tag := range tags {
		var fromRef string
		var prevVersion string
		if i+1 < len(tags) {
			fromRef = tags[i+1].Name
			prevVersion = tags[i+1].Name
		}

		commits, err := git.GetCommitsBetween(fromRef, tag.Name)
		if err != nil {
			output.LogWarning(fmt.Sprintf("Failed to get commits for %s: %v. Skipping.", tag.Name, err))
			continue
		}

		entry := buildEntry(tag.Name, prevVersion, tag.Date, commits, skipRegex, excludeSet, cfg.IncludeBreaking, cfg.IncludeNonConventional, cfg.CustomTypeMapping, cfg.ExcludeAuthors)
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

func filterTags(tags []git.Tag, sinceTag, untilTag string) []git.Tag {
	if sinceTag == "" && untilTag == "" {
		return tags
	}

	// Tags are sorted descending: [v3, v2, v1]
	// until_tag: start collecting from this tag (skip newer)
	// since_tag: stop collecting after this tag (skip older)
	startIdx := 0
	endIdx := len(tags)

	if untilTag != "" {
		for i, tag := range tags {
			if tag.Name == untilTag {
				startIdx = i
				break
			}
		}
	}

	if sinceTag != "" {
		for i, tag := range tags {
			if tag.Name == sinceTag {
				endIdx = i + 1
				break
			}
		}
	}

	if startIdx >= endIdx {
		return nil
	}

	return tags[startIdx:endIdx]
}

func isExcludedAuthor(author string, excludeAuthors []string) bool {
	for _, pattern := range excludeAuthors {
		if strings.EqualFold(author, pattern) {
			return true
		}
		// Wildcard matching: a trailing "*" turns the remainder into a substring
		// match, e.g. "bot*" matches "dependabot", "mybot", and "renovatebot".
		if strings.HasSuffix(pattern, "*") {
			prefix := strings.ToLower(strings.TrimSuffix(pattern, "*"))
			if strings.Contains(strings.ToLower(author), prefix) {
				return true
			}
		}
	}
	return false
}

func buildEntry(version, prevVersion string, date time.Time, commits []git.Commit, skipRegex *regexp.Regexp, excludeSet map[string]bool, includeBreaking, includeNonConventional bool, customMapping map[string]string, excludeAuthors []string) Entry {
	entry := Entry{
		Version:     version,
		PrevVersion: prevVersion,
		Date:        date,
		Sections:    make(map[string][]ConventionalCommit),
	}

	contributorSet := make(map[string]bool)

	for _, commit := range commits {
		if skipRegex != nil && skipRegex.MatchString(commit.Message) {
			continue
		}

		// Track contributors (excluding bots)
		if commit.Author != "" && !isExcludedAuthor(commit.Author, excludeAuthors) {
			contributorSet[commit.Author] = true
		}

		cc := ParseConventionalCommit(commit.Message, commit.Body, commit.Hash, commit.Author)
		if cc == nil {
			if includeNonConventional {
				cc = ParseNonConventionalCommit(commit.Message, commit.Body, commit.Hash, commit.Author)
			} else {
				continue
			}
		}

		if excludeSet[cc.Type] {
			continue
		}

		if includeBreaking && cc.Breaking {
			entry.Breaking = append(entry.Breaking, *cc)
		}

		typeName := TypeDisplayNameWithCustom(cc.Type, customMapping)
		entry.Sections[typeName] = append(entry.Sections[typeName], *cc)
	}

	// Sort contributors
	for author := range contributorSet {
		entry.Contributors = append(entry.Contributors, author)
	}
	sort.Strings(entry.Contributors)

	return entry
}

func hasSections(entry Entry) bool {
	return len(entry.Sections) > 0 || len(entry.Breaking) > 0
}

// sectionOrder defines the display order for changelog sections.
var sectionOrder = []string{
	"Features", "Bug Fixes", "Performance Improvements",
	"Code Refactoring", "Documentation", "Tests",
	"Builds", "Continuous Integration", "Styles", "Chores", "Reverts",
	"Other Changes",
}

func renderMarkdown(entries []Entry, header, dateFormat, repoURL string) string {
	var sb strings.Builder

	sb.WriteString(header)
	sb.WriteString("\n\n")
	sb.WriteString("All notable changes to this project will be documented in this file.\n\n")

	for _, entry := range entries {
		renderVersionHeader(&sb, entry, dateFormat, repoURL)
		renderBreakingChanges(&sb, entry, repoURL)
		renderSections(&sb, entry, repoURL)
		renderContributors(&sb, entry)
	}

	return sb.String()
}

func renderVersionHeader(sb *strings.Builder, entry Entry, dateFormat, repoURL string) {
	dateStr := entry.Date.Format(dateFormat)
	isUnreleased := entry.Version == "Unreleased"

	if repoURL != "" && !isUnreleased {
		if entry.PrevVersion != "" {
			sb.WriteString(fmt.Sprintf("## [%s](%s/compare/%s...%s) (%s)\n\n",
				entry.Version, repoURL, entry.PrevVersion, entry.Version, dateStr))
		} else {
			sb.WriteString(fmt.Sprintf("## [%s](%s/releases/tag/%s) (%s)\n\n",
				entry.Version, repoURL, entry.Version, dateStr))
		}
	} else {
		sb.WriteString(fmt.Sprintf("## %s (%s)\n\n", entry.Version, dateStr))
	}
}

func renderBreakingChanges(sb *strings.Builder, entry Entry, repoURL string) {
	if len(entry.Breaking) == 0 {
		return
	}
	sb.WriteString("### BREAKING CHANGES\n\n")
	for _, cc := range entry.Breaking {
		renderCommitLine(sb, cc, repoURL)
	}
	sb.WriteString("\n")
}

func renderSections(sb *strings.Builder, entry Entry, repoURL string) {
	rendered := make(map[string]bool)

	// Render sections in predefined order
	for _, sectionName := range sectionOrder {
		commits, ok := entry.Sections[sectionName]
		if !ok {
			continue
		}
		rendered[sectionName] = true
		sb.WriteString(fmt.Sprintf("### %s\n\n", sectionName))
		for _, cc := range commits {
			renderCommitLine(sb, cc, repoURL)
		}
		sb.WriteString("\n")
	}

	// Render remaining sections alphabetically
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
			renderCommitLine(sb, cc, repoURL)
		}
		sb.WriteString("\n")
	}
}

func renderContributors(sb *strings.Builder, entry Entry) {
	if len(entry.Contributors) > 0 {
		sb.WriteString("### Contributors\n\n")
		for _, author := range entry.Contributors {
			sb.WriteString(fmt.Sprintf("- %s\n", author))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("<br/>\n\n")
}

func renderCommitLine(sb *strings.Builder, cc ConventionalCommit, repoURL string) {
	shortHash := cc.Hash
	if len(shortHash) > 7 {
		shortHash = shortHash[:7]
	}

	var line strings.Builder

	// Build description with scope
	if cc.Scope != "" {
		line.WriteString(fmt.Sprintf("- **%s:** %s", cc.Scope, cc.Description))
	} else {
		line.WriteString(fmt.Sprintf("- %s", cc.Description))
	}

	// Append PR links
	if repoURL != "" && len(cc.PRNumbers) > 0 {
		var prLinks []string
		for _, pr := range cc.PRNumbers {
			prLinks = append(prLinks, fmt.Sprintf("[#%s](%s/pull/%s)", pr, repoURL, pr))
		}
		line.WriteString(fmt.Sprintf(" (%s)", strings.Join(prLinks, ", ")))
	}

	// Append issue links
	if repoURL != "" && len(cc.Issues) > 0 {
		var issueLinks []string
		for _, issue := range cc.Issues {
			issueLinks = append(issueLinks, fmt.Sprintf("[#%s](%s/issues/%s)", issue, repoURL, issue))
		}
		line.WriteString(fmt.Sprintf(", closes %s", strings.Join(issueLinks, ", ")))
	}

	// Append commit hash link
	if repoURL != "" {
		line.WriteString(fmt.Sprintf(" ([%s](%s/commit/%s))", shortHash, repoURL, cc.Hash))
	} else {
		line.WriteString(fmt.Sprintf(" (%s)", shortHash))
	}

	line.WriteString("\n")
	sb.WriteString(line.String())
}
