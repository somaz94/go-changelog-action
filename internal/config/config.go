package config

import (
	"os"
	"strings"
)

// Config holds all configuration for the changelog action.
type Config struct {
	OutputFile     string
	TagPattern     string
	ExcludeTypes   []string
	IncludeBreaking bool
	DateFormat     string
	Header         string
	Unreleased     bool
	UnreleasedTitle string
	SkipCommits    string
	RepositoryURL  string
	DryRun         bool
}

// Load reads configuration from environment variables (INPUT_*).
func Load() *Config {
	cfg := &Config{
		OutputFile:      getEnvDefault("INPUT_OUTPUT_FILE", "CHANGELOG.md"),
		TagPattern:      getEnvDefault("INPUT_TAG_PATTERN", "v[0-9]*.[0-9]*.[0-9]*"),
		IncludeBreaking: getEnvDefault("INPUT_INCLUDE_BREAKING", "true") == "true",
		DateFormat:      getEnvDefault("INPUT_DATE_FORMAT", "2006-01-02"),
		Header:          getEnvDefault("INPUT_HEADER", "# Changelog"),
		Unreleased:      getEnvDefault("INPUT_UNRELEASED", "true") == "true",
		UnreleasedTitle: getEnvDefault("INPUT_UNRELEASED_TITLE", "Unreleased"),
		SkipCommits:     getEnvDefault("INPUT_SKIP_COMMITS", "^Merge"),
		RepositoryURL:   os.Getenv("INPUT_REPOSITORY_URL"),
		DryRun:          getEnvDefault("INPUT_DRY_RUN", "false") == "true",
	}

	if excludeStr := os.Getenv("INPUT_EXCLUDE_TYPES"); excludeStr != "" {
		for _, t := range strings.Split(excludeStr, ",") {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				cfg.ExcludeTypes = append(cfg.ExcludeTypes, trimmed)
			}
		}
	}

	return cfg
}

func getEnvDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
