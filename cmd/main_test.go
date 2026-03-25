package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/somaz94/go-changelog-action/internal/git"
)

func TestRunDryRun(t *testing.T) {
	// Mock git commands
	original := git.RunCommand
	git.RunCommand = func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "tag" {
			return []byte("v1.0.0|abc1234|2024-01-01T00:00:00Z\n"), nil
		}
		if len(args) > 0 && args[0] == "remote" {
			return []byte("https://github.com/owner/repo\n"), nil
		}
		if len(args) > 0 && args[0] == "log" {
			return []byte("aaa111\x01feat: test feature\x012024-01-15T10:00:00Z\x01alice\x01\x00"), nil
		}
		return []byte(""), nil
	}
	defer func() { git.RunCommand = original }()

	// Set env vars for dry run
	os.Setenv("INPUT_DRY_RUN", "true")
	os.Setenv("GITHUB_WORKSPACE", t.TempDir())
	defer func() {
		os.Unsetenv("INPUT_DRY_RUN")
		os.Unsetenv("GITHUB_WORKSPACE")
	}()

	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWriteFile(t *testing.T) {
	original := git.RunCommand
	git.RunCommand = func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "tag" {
			return []byte("v1.0.0|abc1234|2024-01-01T00:00:00Z\n"), nil
		}
		if len(args) > 0 && args[0] == "remote" {
			return []byte("https://github.com/owner/repo\n"), nil
		}
		if len(args) > 0 && args[0] == "log" {
			return []byte("aaa111\x01feat: feature\x012024-01-15T10:00:00Z\x01alice\x01\x00"), nil
		}
		return []byte(""), nil
	}
	defer func() { git.RunCommand = original }()

	tmpDir := t.TempDir()
	os.Setenv("INPUT_DRY_RUN", "false")
	os.Setenv("INPUT_OUTPUT_FILE", tmpDir+"/CHANGELOG.md")
	os.Setenv("GITHUB_WORKSPACE", tmpDir)
	defer func() {
		os.Unsetenv("INPUT_DRY_RUN")
		os.Unsetenv("INPUT_OUTPUT_FILE")
		os.Unsetenv("GITHUB_WORKSPACE")
	}()

	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file was created
	data, err := os.ReadFile(tmpDir + "/CHANGELOG.md")
	if err != nil {
		t.Fatalf("expected changelog file to exist: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty changelog file")
	}
}

func TestRunWriteFileWithGitHubOutput(t *testing.T) {
	original := git.RunCommand
	git.RunCommand = func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "tag" {
			return []byte("v1.0.0|abc1234|2024-01-01T00:00:00Z\n"), nil
		}
		if len(args) > 0 && args[0] == "remote" {
			return []byte("https://github.com/owner/repo\n"), nil
		}
		if len(args) > 0 && args[0] == "log" {
			return []byte("aaa111\x01feat: feature\x012024-01-15T10:00:00Z\x01alice\x01\x00"), nil
		}
		return []byte(""), nil
	}
	defer func() { git.RunCommand = original }()

	tmpDir := t.TempDir()
	outputFile := tmpDir + "/github_output"
	os.WriteFile(outputFile, []byte{}, 0644)

	os.Setenv("INPUT_DRY_RUN", "false")
	os.Setenv("INPUT_OUTPUT_FILE", tmpDir+"/CHANGELOG.md")
	os.Setenv("GITHUB_WORKSPACE", tmpDir)
	os.Setenv("GITHUB_OUTPUT", outputFile)
	defer func() {
		os.Unsetenv("INPUT_DRY_RUN")
		os.Unsetenv("INPUT_OUTPUT_FILE")
		os.Unsetenv("GITHUB_WORKSPACE")
		os.Unsetenv("GITHUB_OUTPUT")
	}()

	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify GITHUB_OUTPUT was written
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read GITHUB_OUTPUT: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected GITHUB_OUTPUT to have content")
	}
}

func TestRunInvalidWorkDir(t *testing.T) {
	original := git.RunCommand
	git.RunCommand = func(args ...string) ([]byte, error) {
		return []byte(""), nil
	}
	defer func() { git.RunCommand = original }()

	os.Setenv("GITHUB_WORKSPACE", "/nonexistent/path/12345")
	defer os.Unsetenv("GITHUB_WORKSPACE")

	ctx := context.Background()
	err := run(ctx)
	if err == nil {
		t.Fatal("expected error for invalid workspace directory")
	}
	if !strings.Contains(err.Error(), "failed to change directory") {
		t.Errorf("expected chdir error, got: %v", err)
	}
}

func TestRunGenerateError(t *testing.T) {
	original := git.RunCommand
	git.RunCommand = func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "tag" {
			return nil, fmt.Errorf("git not found")
		}
		return []byte(""), nil
	}
	defer func() { git.RunCommand = original }()

	tmpDir := t.TempDir()
	os.Setenv("GITHUB_WORKSPACE", tmpDir)
	defer os.Unsetenv("GITHUB_WORKSPACE")

	ctx := context.Background()
	err := run(ctx)
	if err == nil {
		t.Fatal("expected error when generate fails")
	}
	if !strings.Contains(err.Error(), "failed to generate changelog") {
		t.Errorf("expected generate error, got: %v", err)
	}
}

func TestRunDefaultWorkDir(t *testing.T) {
	original := git.RunCommand
	git.RunCommand = func(args ...string) ([]byte, error) {
		return []byte(""), nil
	}
	defer func() { git.RunCommand = original }()

	// No GITHUB_WORKSPACE set → defaults to /app which doesn't exist
	os.Unsetenv("GITHUB_WORKSPACE")

	ctx := context.Background()
	err := run(ctx)
	if err == nil {
		t.Fatal("expected error for default /app directory")
	}
	if !strings.Contains(err.Error(), "failed to change directory") {
		t.Errorf("expected chdir error, got: %v", err)
	}
}

func TestRunWriteFileError(t *testing.T) {
	original := git.RunCommand
	git.RunCommand = func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "tag" {
			return []byte("v1.0.0|abc1234|2024-01-01T00:00:00Z\n"), nil
		}
		if len(args) > 0 && args[0] == "remote" {
			return []byte("https://github.com/owner/repo\n"), nil
		}
		if len(args) > 0 && args[0] == "log" {
			return []byte("aaa111\x01feat: feature\x012024-01-15T10:00:00Z\x01alice\x01\x00"), nil
		}
		return []byte(""), nil
	}
	defer func() { git.RunCommand = original }()

	tmpDir := t.TempDir()
	// Point output to a non-existent directory so WriteFile fails
	os.Setenv("INPUT_DRY_RUN", "false")
	os.Setenv("INPUT_OUTPUT_FILE", tmpDir+"/nonexistent/CHANGELOG.md")
	os.Setenv("GITHUB_WORKSPACE", tmpDir)
	defer func() {
		os.Unsetenv("INPUT_DRY_RUN")
		os.Unsetenv("INPUT_OUTPUT_FILE")
		os.Unsetenv("GITHUB_WORKSPACE")
	}()

	ctx := context.Background()
	err := run(ctx)
	if err == nil {
		t.Fatal("expected error when write fails")
	}
	if !strings.Contains(err.Error(), "failed to write changelog") {
		t.Errorf("expected write error, got: %v", err)
	}
}

func TestRunPathTraversal(t *testing.T) {
	original := git.RunCommand
	git.RunCommand = func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "tag" {
			return []byte("v1.0.0|abc1234|2024-01-01T00:00:00Z\n"), nil
		}
		if len(args) > 0 && args[0] == "remote" {
			return []byte("https://github.com/owner/repo\n"), nil
		}
		if len(args) > 0 && args[0] == "log" {
			return []byte("aaa111\x01feat: feature\x012024-01-15T10:00:00Z\x01alice\x01\x00"), nil
		}
		return []byte(""), nil
	}
	defer func() { git.RunCommand = original }()

	tmpDir := t.TempDir()
	os.Setenv("INPUT_DRY_RUN", "false")
	os.Setenv("INPUT_OUTPUT_FILE", "../../etc/passwd")
	os.Setenv("GITHUB_WORKSPACE", tmpDir)
	defer func() {
		os.Unsetenv("INPUT_DRY_RUN")
		os.Unsetenv("INPUT_OUTPUT_FILE")
		os.Unsetenv("GITHUB_WORKSPACE")
	}()

	ctx := context.Background()
	err := run(ctx)
	if err == nil {
		t.Fatal("expected error for path traversal")
	}
	if !strings.Contains(err.Error(), "outside workspace") {
		t.Errorf("expected 'outside workspace' error, got: %v", err)
	}
}

func TestRunCancelled(t *testing.T) {
	original := git.RunCommand
	git.RunCommand = func(args ...string) ([]byte, error) {
		return []byte(""), nil
	}
	defer func() { git.RunCommand = original }()

	tmpDir := t.TempDir()
	os.Setenv("GITHUB_WORKSPACE", tmpDir)
	defer os.Unsetenv("GITHUB_WORKSPACE")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := run(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
	if err.Error() != "cancelled" {
		t.Errorf("expected 'cancelled' error, got %q", err.Error())
	}
}
