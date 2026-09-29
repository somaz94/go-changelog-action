package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/somaz94/go-changelog-action/internal/git"
)

// isolateGitConfigEnv starts the test with no GIT_CONFIG_COUNT/KEY_<n>/VALUE_<n>
// set; t.Setenv restores the originals and the cleanup drops entries run() added.
func isolateGitConfigEnv(t *testing.T) {
	t.Helper()
	for _, name := range append(gitConfigEntryNames(), "GIT_CONFIG_COUNT") {
		t.Setenv(name, "")
		os.Unsetenv(name)
	}
	t.Cleanup(func() {
		for _, name := range gitConfigEntryNames() {
			os.Unsetenv(name)
		}
	})
}

func gitConfigEntryNames() []string {
	var names []string
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		if strings.HasPrefix(name, "GIT_CONFIG_KEY_") || strings.HasPrefix(name, "GIT_CONFIG_VALUE_") {
			names = append(names, name)
		}
	}
	return names
}

func TestRunDryRun(t *testing.T) {
	isolateGitConfigEnv(t)
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
	isolateGitConfigEnv(t)
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

	data, err := os.ReadFile(tmpDir + "/CHANGELOG.md")
	if err != nil {
		t.Fatalf("expected changelog file to exist: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty changelog file")
	}
}

func TestRunWriteFileWithGitHubOutput(t *testing.T) {
	isolateGitConfigEnv(t)
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

	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read GITHUB_OUTPUT: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected GITHUB_OUTPUT to have content")
	}
}

func TestRunInvalidWorkDir(t *testing.T) {
	isolateGitConfigEnv(t)
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
	isolateGitConfigEnv(t)
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
	isolateGitConfigEnv(t)
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
	isolateGitConfigEnv(t)
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

func TestRunAbsoluteOutputFile(t *testing.T) {
	isolateGitConfigEnv(t)
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
	outputPath := tmpDir + "/CHANGELOG.md"
	os.Setenv("INPUT_DRY_RUN", "false")
	os.Setenv("INPUT_OUTPUT_FILE", outputPath)
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

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty file")
	}
}

func TestRunPathTraversal(t *testing.T) {
	isolateGitConfigEnv(t)
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
	isolateGitConfigEnv(t)
	original := git.RunCommand
	git.RunCommand = func(args ...string) ([]byte, error) {
		return []byte(""), nil
	}
	defer func() { git.RunCommand = original }()

	tmpDir := t.TempDir()
	os.Setenv("GITHUB_WORKSPACE", tmpDir)
	defer os.Unsetenv("GITHUB_WORKSPACE")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := run(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
	if err.Error() != "cancelled" {
		t.Errorf("expected 'cancelled' error, got %q", err.Error())
	}
}

func mockDryRunGit(t *testing.T) {
	t.Helper()
	original := git.RunCommand
	git.RunCommand = func(args ...string) ([]byte, error) {
		if len(args) > 0 && args[0] == "tag" {
			return []byte("v1.0.0|abc1234|2024-01-01T00:00:00Z\n"), nil
		}
		if len(args) > 0 && args[0] == "log" {
			return []byte("aaa111\x01feat: feature\x012024-01-15T10:00:00Z\x01alice\x01\x00"), nil
		}
		return []byte(""), nil
	}
	t.Cleanup(func() { git.RunCommand = original })
	t.Setenv("INPUT_DRY_RUN", "true")
}

func TestRunInjectsSafeDirectoryEnv(t *testing.T) {
	isolateGitConfigEnv(t)
	mockDryRunGit(t)
	tmpDir := t.TempDir()
	t.Setenv("GITHUB_WORKSPACE", tmpDir)

	if err := run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for name, want := range map[string]string{
		"GIT_CONFIG_COUNT":   "1",
		"GIT_CONFIG_KEY_0":   "safe.directory",
		"GIT_CONFIG_VALUE_0": tmpDir,
	} {
		if got := os.Getenv(name); got != want {
			t.Errorf("expected %s=%q, got %q", name, want, got)
		}
	}
}

func TestRunInvalidGitConfigCountOnlyWarns(t *testing.T) {
	isolateGitConfigEnv(t)
	mockDryRunGit(t)
	t.Setenv("GITHUB_WORKSPACE", t.TempDir())
	t.Setenv("GIT_CONFIG_COUNT", "bogus")

	if err := run(context.Background()); err != nil {
		t.Fatalf("expected an invalid GIT_CONFIG_COUNT to only warn, got: %v", err)
	}
	if got := os.Getenv("GIT_CONFIG_COUNT"); got != "bogus" {
		t.Errorf("expected GIT_CONFIG_COUNT to stay %q, got %q", "bogus", got)
	}
}
