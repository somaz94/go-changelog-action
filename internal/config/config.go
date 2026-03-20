package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/somaz94/go-changelog-action/internal/output"
)

// Config holds all configuration for the changelog action.
type Config struct {
	OutputFile             string
	TagPattern             string
	ExcludeTypes           []string
	IncludeBreaking        bool
	DateFormat             string
	Header                 string
	Unreleased             bool
	UnreleasedTitle        string
	SkipCommits            string
	RepositoryURL          string
	DryRun                 bool
	IncludeNonConventional bool
	SinceTag               string
	UntilTag               string
	CustomTypeMapping      map[string]string
	ExcludeAuthors         []string
}

// Load reads configuration from environment variables (INPUT_*).
func Load() *Config {
	cfg := &Config{
		OutputFile:             getEnvDefault("INPUT_OUTPUT_FILE", "CHANGELOG.md"),
		TagPattern:             getEnvDefault("INPUT_TAG_PATTERN", "v[0-9]*.[0-9]*.[0-9]*"),
		IncludeBreaking:        getEnvDefault("INPUT_INCLUDE_BREAKING", "true") == "true",
		DateFormat:             getEnvDefault("INPUT_DATE_FORMAT", "2006-01-02"),
		Header:                 getEnvDefault("INPUT_HEADER", "# Changelog"),
		Unreleased:             getEnvDefault("INPUT_UNRELEASED", "true") == "true",
		UnreleasedTitle:        getEnvDefault("INPUT_UNRELEASED_TITLE", "Unreleased"),
		SkipCommits:            getEnvDefault("INPUT_SKIP_COMMITS", "^Merge|^docs: update changelog|^docs: update CONTRIBUTORS"),
		RepositoryURL:          os.Getenv("INPUT_REPOSITORY_URL"),
		DryRun:                 getEnvDefault("INPUT_DRY_RUN", "false") == "true",
		IncludeNonConventional: getEnvDefault("INPUT_INCLUDE_NON_CONVENTIONAL", "false") == "true",
		SinceTag:               os.Getenv("INPUT_SINCE_TAG"),
		UntilTag:               os.Getenv("INPUT_UNTIL_TAG"),
	}

	if excludeAuthors := os.Getenv("INPUT_EXCLUDE_AUTHORS"); excludeAuthors != "" {
		for _, a := range strings.Split(excludeAuthors, ",") {
			trimmed := strings.TrimSpace(a)
			if trimmed != "" {
				cfg.ExcludeAuthors = append(cfg.ExcludeAuthors, trimmed)
			}
		}
	}

	if excludeStr := os.Getenv("INPUT_EXCLUDE_TYPES"); excludeStr != "" {
		for _, t := range strings.Split(excludeStr, ",") {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				cfg.ExcludeTypes = append(cfg.ExcludeTypes, trimmed)
			}
		}
	}

	if mappingStr := os.Getenv("INPUT_CUSTOM_TYPE_MAPPING"); mappingStr != "" {
		mapping := make(map[string]string)
		if err := json.Unmarshal([]byte(mappingStr), &mapping); err != nil {
			output.LogWarning(fmt.Sprintf("Failed to parse custom_type_mapping JSON: %v. Using defaults.", err))
		} else {
			cfg.CustomTypeMapping = mapping
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
