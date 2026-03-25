package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	envVars := []string{
		"INPUT_OUTPUT_FILE", "INPUT_TAG_PATTERN", "INPUT_EXCLUDE_TYPES",
		"INPUT_INCLUDE_BREAKING", "INPUT_DATE_FORMAT", "INPUT_HEADER",
		"INPUT_UNRELEASED", "INPUT_UNRELEASED_TITLE", "INPUT_SKIP_COMMITS",
		"INPUT_REPOSITORY_URL", "INPUT_DRY_RUN", "INPUT_INCLUDE_NON_CONVENTIONAL",
		"INPUT_SINCE_TAG", "INPUT_UNTIL_TAG", "INPUT_CUSTOM_TYPE_MAPPING",
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
	if cfg.IncludeNonConventional {
		t.Error("expected IncludeNonConventional=false")
	}
	if cfg.SinceTag != "" {
		t.Errorf("expected empty SinceTag, got %s", cfg.SinceTag)
	}
	if cfg.UntilTag != "" {
		t.Errorf("expected empty UntilTag, got %s", cfg.UntilTag)
	}
	if cfg.CustomTypeMapping != nil {
		t.Errorf("expected nil CustomTypeMapping, got %v", cfg.CustomTypeMapping)
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

func TestLoadNewFeatures(t *testing.T) {
	os.Setenv("INPUT_INCLUDE_NON_CONVENTIONAL", "true")
	os.Setenv("INPUT_SINCE_TAG", "v1.0.0")
	os.Setenv("INPUT_UNTIL_TAG", "v2.0.0")
	os.Setenv("INPUT_CUSTOM_TYPE_MAPPING", `{"feat": "New Features", "fix": "Bugfixes"}`)
	defer func() {
		os.Unsetenv("INPUT_INCLUDE_NON_CONVENTIONAL")
		os.Unsetenv("INPUT_SINCE_TAG")
		os.Unsetenv("INPUT_UNTIL_TAG")
		os.Unsetenv("INPUT_CUSTOM_TYPE_MAPPING")
	}()

	cfg := Load()

	if !cfg.IncludeNonConventional {
		t.Error("expected IncludeNonConventional=true")
	}
	if cfg.SinceTag != "v1.0.0" {
		t.Errorf("expected SinceTag=v1.0.0, got %s", cfg.SinceTag)
	}
	if cfg.UntilTag != "v2.0.0" {
		t.Errorf("expected UntilTag=v2.0.0, got %s", cfg.UntilTag)
	}
	if cfg.CustomTypeMapping == nil {
		t.Fatal("expected non-nil CustomTypeMapping")
	}
	if cfg.CustomTypeMapping["feat"] != "New Features" {
		t.Errorf("expected feat -> 'New Features', got %q", cfg.CustomTypeMapping["feat"])
	}
	if cfg.CustomTypeMapping["fix"] != "Bugfixes" {
		t.Errorf("expected fix -> 'Bugfixes', got %q", cfg.CustomTypeMapping["fix"])
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	os.Setenv("INPUT_CUSTOM_TYPE_MAPPING", "not valid json")
	defer os.Unsetenv("INPUT_CUSTOM_TYPE_MAPPING")

	cfg := Load()

	if cfg.CustomTypeMapping != nil {
		t.Errorf("expected nil CustomTypeMapping for invalid JSON, got %v", cfg.CustomTypeMapping)
	}
}

func TestLoadExcludeAuthors(t *testing.T) {
	os.Setenv("INPUT_EXCLUDE_AUTHORS", "GitHub Actions, dependabot[bot], , renovate[bot]")
	defer os.Unsetenv("INPUT_EXCLUDE_AUTHORS")

	cfg := Load()

	if len(cfg.ExcludeAuthors) != 3 {
		t.Fatalf("expected 3 ExcludeAuthors, got %d: %v", len(cfg.ExcludeAuthors), cfg.ExcludeAuthors)
	}
	expected := []string{"GitHub Actions", "dependabot[bot]", "renovate[bot]"}
	for i, v := range expected {
		if cfg.ExcludeAuthors[i] != v {
			t.Errorf("expected ExcludeAuthors[%d]=%s, got %s", i, v, cfg.ExcludeAuthors[i])
		}
	}
}

func TestLoadExcludeAuthorsEmpty(t *testing.T) {
	os.Unsetenv("INPUT_EXCLUDE_AUTHORS")

	cfg := Load()

	if len(cfg.ExcludeAuthors) != 0 {
		t.Errorf("expected empty ExcludeAuthors, got %v", cfg.ExcludeAuthors)
	}
}
