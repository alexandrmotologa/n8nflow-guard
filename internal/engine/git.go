package engine

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetStagedWorkflowFiles finds JSON files staged for commit in git index.
func GetStagedWorkflowFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--cached")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git diff --cached failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}

	return filterJSONFiles(stdout.String())
}

// GetDiffWorkflowFiles finds JSON files changed relative to a git ref (e.g. "origin/main" or "HEAD~1").
func GetDiffWorkflowFiles(ref string) ([]string, error) {
	if ref == "" {
		ref = "HEAD"
	}
	cmd := exec.Command("git", "diff", "--name-only", ref)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git diff %s failed: %v (%s)", ref, err, strings.TrimSpace(stderr.String()))
	}

	return filterJSONFiles(stdout.String())
}

func filterJSONFiles(rawOutput string) ([]string, error) {
	lines := strings.Split(strings.TrimSpace(rawOutput), "\n")
	var files []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if !strings.EqualFold(filepath.Ext(trimmed), ".json") {
			continue
		}
		// Confirm file exists on disk (was not deleted)
		if _, err := os.Stat(trimmed); err == nil {
			files = append(files, trimmed)
		}
	}

	return files, nil
}
