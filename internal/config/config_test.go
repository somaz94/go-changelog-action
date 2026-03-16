package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	// Clear all relevant env vars
	envVars := []string{
		"INPUT_OUTPUT_FILE", "INPUT_TAG_PATTERN", "INPUT_EXCLUDE_TYPES",
		"INPUT_INCLUDE_BREAKING", "INPUT_DATE_FORMAT", "INPUT_HEADER",
		"INPUT_UNRELEASED", "INPUT_UNRELEASED_TITLE", "INPUT_SKIP_COMMITS",
		"INPUT_REPOSITORY_URL", "INPUT_DRY_RUN",
	}
	for _, v := range envVars {
		os.Unsetenv(v)
	}

	cfg := Load()

	if cfg.OutputFile != "CHANGELOG.md" {
		t.Errorf("expected OutputFile=CHANGELOG.md, got %s", cfg.OutputFile)
	}
	if cfg.TagPattern != "v[0-9]*.[0-9]*.[0-9]*" {
		t.Errorf("expected default TagPattern, got %s", cfg.TagPattern)
	}
	if !cfg.IncludeBreaking {
		t.Error("expected IncludeBreaking=true")
	}
	if !cfg.Unreleased {
		t.Error("expected Unreleased=true")
	}
	if cfg.DryRun {
		t.Error("expected DryRun=false")
	}
	if len(cfg.ExcludeTypes) != 0 {
		t.Errorf("expected empty ExcludeTypes, got %v", cfg.ExcludeTypes)
	}
}

func TestLoadCustomValues(t *testing.T) {
	os.Setenv("INPUT_OUTPUT_FILE", "CHANGES.md")
	os.Setenv("INPUT_EXCLUDE_TYPES", "chore, style, ci")
	os.Setenv("INPUT_DRY_RUN", "true")
	os.Setenv("INPUT_INCLUDE_BREAKING", "false")
	defer func() {
		os.Unsetenv("INPUT_OUTPUT_FILE")
		os.Unsetenv("INPUT_EXCLUDE_TYPES")
		os.Unsetenv("INPUT_DRY_RUN")
		os.Unsetenv("INPUT_INCLUDE_BREAKING")
	}()

	cfg := Load()

	if cfg.OutputFile != "CHANGES.md" {
		t.Errorf("expected OutputFile=CHANGES.md, got %s", cfg.OutputFile)
	}
	if !cfg.DryRun {
		t.Error("expected DryRun=true")
	}
	if cfg.IncludeBreaking {
		t.Error("expected IncludeBreaking=false")
	}
	if len(cfg.ExcludeTypes) != 3 {
		t.Errorf("expected 3 ExcludeTypes, got %d: %v", len(cfg.ExcludeTypes), cfg.ExcludeTypes)
	}
	expected := []string{"chore", "style", "ci"}
	for i, v := range expected {
		if cfg.ExcludeTypes[i] != v {
			t.Errorf("expected ExcludeTypes[%d]=%s, got %s", i, v, cfg.ExcludeTypes[i])
		}
	}
}
