package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIVersion(t *testing.T) {
	cmd := newVersionCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	cmd.Run(cmd, []string{})

	out := buf.String()
	if !strings.Contains(out, "n8nflow-guard version") {
		t.Errorf("expected version output, got: %s", out)
	}
}

func TestCLIRules(t *testing.T) {
	cmd := newRulesCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "SEC-001") || !strings.Contains(out, "Total rules: 34") {
		t.Errorf("expected 34 rules in output, got:\n%s", out)
	}
}

func TestCLIRulesSpecific(t *testing.T) {
	cmd := newRulesCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--rule", "SEC-001"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Rule ID:      SEC-001") {
		t.Errorf("expected rule SEC-001 info, got: %s", out)
	}
}

func TestCLIInit(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get wd: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	cmd := newInitCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, ".n8nflow-guard.yaml")); err != nil {
		t.Errorf("expected .n8nflow-guard.yaml to be created")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, ".pre-commit-config.yaml")); err != nil {
		t.Errorf("expected .pre-commit-config.yaml to be created")
	}
}

func TestCLIGraph(t *testing.T) {
	tmpDir := t.TempDir()
	workflowFile := filepath.Join(tmpDir, "flow.json")
	content := []byte(`{
		"name": "Simple Flow",
		"nodes": [{"name": "Start", "type": "n8n-nodes-base.scheduleTrigger"}]
	}`)
	if err := os.WriteFile(workflowFile, content, 0644); err != nil {
		t.Fatalf("failed writing test file: %v", err)
	}

	cmd := newGraphCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{workflowFile, "--format", "mermaid"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error executing graph cmd: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "flowchart TD") || !strings.Contains(out, "Start") {
		t.Errorf("expected mermaid graph output, got:\n%s", out)
	}
}
