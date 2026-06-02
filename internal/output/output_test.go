package output

import (
	"os"
	"strings"
	"testing"
)

func TestSetOutputToFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "github_output")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	os.Setenv("GITHUB_OUTPUT", tmpFile.Name())
	defer os.Unsetenv("GITHUB_OUTPUT")

	err = SetOutput("test_key", "test_value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "test_key=test_value") {
		t.Errorf("expected output to contain 'test_key=test_value', got: %s", string(data))
	}
}

func TestSetOutputMultiline(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "github_output")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	os.Setenv("GITHUB_OUTPUT", tmpFile.Name())
	defer os.Unsetenv("GITHUB_OUTPUT")

	err = SetOutput("content", "line1\nline2\nline3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	if !strings.Contains(content, "content<<ghadelimiter_") {
		t.Errorf("expected multiline delimiter, got: %s", content)
	}
	if !strings.Contains(content, "line1\nline2\nline3") {
		t.Errorf("expected multiline content, got: %s", content)
	}
}

func TestSetOutputFallback(t *testing.T) {
	os.Unsetenv("GITHUB_OUTPUT")

	// Should not error when GITHUB_OUTPUT is not set (fallback mode)
	err := SetOutput("key", "value")
	if err != nil {
		t.Fatalf("unexpected error in fallback mode: %v", err)
	}
}

func TestSetOutputOpenFileError(t *testing.T) {
	// Point GITHUB_OUTPUT to a path that can't be opened
	os.Setenv("GITHUB_OUTPUT", "/nonexistent/dir/output")
	defer os.Unsetenv("GITHUB_OUTPUT")

	err := SetOutput("key", "value")
	if err == nil {
		t.Fatal("expected error when GITHUB_OUTPUT path is invalid")
	}
	if !strings.Contains(err.Error(), "failed to open GITHUB_OUTPUT") {
		t.Errorf("expected open error, got: %v", err)
	}
}

func TestLogInfo(t *testing.T) {
	// Capture stdout
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	LogInfo("test message")

	w.Close()
	os.Stdout = oldStdout

	buf := make([]byte, 256)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if !strings.Contains(output, "::notice::test message") {
		t.Errorf("expected '::notice::test message', got %q", output)
	}
}

func TestLogWarning(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	LogWarning("warn msg")

	w.Close()
	os.Stdout = oldStdout

	buf := make([]byte, 256)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if !strings.Contains(output, "::warning::warn msg") {
		t.Errorf("expected '::warning::warn msg', got %q", output)
	}
}

func TestLogError(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	LogError("error msg")

	w.Close()
	os.Stdout = oldStdout

	buf := make([]byte, 256)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	if !strings.Contains(output, "::error::error msg") {
		t.Errorf("expected '::error::error msg', got %q", output)
	}
}
