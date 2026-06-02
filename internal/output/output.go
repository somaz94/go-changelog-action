package output

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// SetOutput sets a GitHub Actions output variable.
func SetOutput(name, value string) (err error) {
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		// Fallback for local testing
		fmt.Printf("::set-output name=%s::%s\n", name, value)
		return nil
	}

	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open GITHUB_OUTPUT: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Use a collision-resistant delimiter for multiline values so a commit body
	// line equal to the delimiter token cannot terminate the heredoc early.
	if strings.Contains(value, "\n") {
		delimiter := fmt.Sprintf("ghadelimiter_%d", time.Now().UnixNano())
		_, err = fmt.Fprintf(f, "%s<<%s\n%s\n%s\n", name, delimiter, value, delimiter)
	} else {
		_, err = fmt.Fprintf(f, "%s=%s\n", name, value)
	}

	return err
}

// LogInfo prints an info message in GitHub Actions format.
func LogInfo(msg string) {
	fmt.Printf("::notice::%s\n", msg)
}

// LogWarning prints a warning message in GitHub Actions format.
func LogWarning(msg string) {
	fmt.Printf("::warning::%s\n", msg)
}

// LogError prints an error message in GitHub Actions format.
func LogError(msg string) {
	fmt.Printf("::error::%s\n", msg)
}
