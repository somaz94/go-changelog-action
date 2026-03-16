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
	if !strings.Contains(content, "content<<EOF") {
		t.Errorf("expected multiline delimiter, got: %s", content)
	}
	if !strings.Contains(content, "line1\nline2\nline3") {
		t.Errorf("expected multiline content, got: %s", content)
	}
}
