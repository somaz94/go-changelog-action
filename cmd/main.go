package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/somaz94/go-changelog-action/internal/changelog"
	"github.com/somaz94/go-changelog-action/internal/config"
	"github.com/somaz94/go-changelog-action/internal/output"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		output.LogWarning("Received shutdown signal, cleaning up...")
		cancel()
	}()

	if err := run(ctx); err != nil {
		output.LogError(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg := config.Load()

	output.LogInfo("Starting changelog generation...")
	output.LogInfo(fmt.Sprintf("Output file: %s", cfg.OutputFile))
	output.LogInfo(fmt.Sprintf("Tag pattern: %s", cfg.TagPattern))

	// Configure working directory
	workDir := os.Getenv("GITHUB_WORKSPACE")
	if workDir == "" {
		workDir = "/app"
	}
	if err := os.Chdir(workDir); err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", workDir, err)
	}

	// Configure git safe directory to avoid ownership issues in containers
	if err := exec.Command("git", "config", "--global", "--add", "safe.directory", workDir).Run(); err != nil {
		output.LogWarning(fmt.Sprintf("Failed to set git safe.directory: %v", err))
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("cancelled")
	default:
	}

	genCfg := changelog.GeneratorConfig{
		TagPattern:             cfg.TagPattern,
		ExcludeTypes:           cfg.ExcludeTypes,
		IncludeBreaking:        cfg.IncludeBreaking,
		DateFormat:             cfg.DateFormat,
		Header:                 cfg.Header,
		Unreleased:             cfg.Unreleased,
		UnreleasedTitle:        cfg.UnreleasedTitle,
		SkipCommits:            cfg.SkipCommits,
		RepositoryURL:          cfg.RepositoryURL,
		IncludeNonConventional: cfg.IncludeNonConventional,
		SinceTag:               cfg.SinceTag,
		UntilTag:               cfg.UntilTag,
		CustomTypeMapping:      cfg.CustomTypeMapping,
		ExcludeAuthors:         cfg.ExcludeAuthors,
	}

	result, err := changelog.Generate(genCfg)
	if err != nil {
		return fmt.Errorf("failed to generate changelog: %w", err)
	}

	if cfg.DryRun {
		output.LogInfo("Dry run mode - changelog preview:")
		fmt.Println(result.Content)
	} else {
		if err := os.WriteFile(cfg.OutputFile, []byte(result.Content), 0644); err != nil {
			return fmt.Errorf("failed to write changelog: %w", err)
		}
		output.LogInfo(fmt.Sprintf("Changelog written to %s", cfg.OutputFile))
	}

	// Set outputs
	if err := output.SetOutput("changelog_file", cfg.OutputFile); err != nil {
		output.LogWarning(fmt.Sprintf("Failed to set changelog_file output: %v", err))
	}
	if err := output.SetOutput("changelog_content", result.Content); err != nil {
		output.LogWarning(fmt.Sprintf("Failed to set changelog_content output: %v", err))
	}
	if err := output.SetOutput("entries_count", strconv.Itoa(result.EntriesCount)); err != nil {
		output.LogWarning(fmt.Sprintf("Failed to set entries_count output: %v", err))
	}
	if err := output.SetOutput("latest_version", result.LatestVersion); err != nil {
		output.LogWarning(fmt.Sprintf("Failed to set latest_version output: %v", err))
	}

	output.LogInfo(fmt.Sprintf("Changelog generated successfully: %d entries, latest version: %s", result.EntriesCount, result.LatestVersion))
	return nil
}
