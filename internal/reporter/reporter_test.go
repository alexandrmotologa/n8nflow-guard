package reporter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

func sampleData() ([]models.Issue, Summary) {
	issues := []models.Issue{
		{
			RuleID:      "SEC-001",
			RuleName:    "hardcoded-secret",
			Category:    models.CategorySecurity,
			Severity:    models.SeverityCritical,
			Message:     "Hardcoded secret detected",
			FilePath:    "workflows/payment.json",
			Line:        15,
			Fixable:     true,
			Remediation: "Use environment variable",
		},
		{
			RuleID:      "REL-001",
			RuleName:    "missing-error-workflow",
			Category:    models.CategoryReliability,
			Severity:    models.SeverityHigh,
			Message:     "No error workflow",
			FilePath:    "workflows/payment.json",
			Line:        1,
			Fixable:     false,
			Remediation: "Configure errorWorkflow",
		},
	}

	summary := ComputeSummary(1, 1, issues)
	return issues, summary
}

func TestTerminalReporter(t *testing.T) {
	issues, summary := sampleData()
	rep := &TerminalReporter{NoColor: true}

	var buf bytes.Buffer
	if err := rep.Report(&buf, issues, summary); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "SEC-001") || !strings.Contains(out, "REL-001") {
		t.Errorf("expected output to contain rule IDs, got:\n%s", out)
	}
	if !strings.Contains(out, "Critical:        1") {
		t.Errorf("expected summary with 1 critical, got:\n%s", out)
	}
}

func TestSARIFReporter(t *testing.T) {
	issues, summary := sampleData()
	rep := &SARIFReporter{}

	var buf bytes.Buffer
	if err := rep.Report(&buf, issues, summary); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"version": "2.1.0"`) {
		t.Errorf("expected SARIF 2.1.0 version, got:\n%s", out)
	}
	if !strings.Contains(out, `"ruleId": "SEC-001"`) {
		t.Errorf("expected SEC-001 in results, got:\n%s", out)
	}
}

func TestGitHubActionsReporter(t *testing.T) {
	issues, summary := sampleData()
	rep := &GitHubActionsReporter{}

	var buf bytes.Buffer
	if err := rep.Report(&buf, issues, summary); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "::error file=workflows/payment.json,line=15,title=[SEC-001] hardcoded-secret::") {
		t.Errorf("expected GHA error command, got:\n%s", out)
	}
}

func TestJSONReporter(t *testing.T) {
	issues, summary := sampleData()
	rep := &JSONReporter{}

	var buf bytes.Buffer
	if err := rep.Report(&buf, issues, summary); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"totalIssues": 2`) {
		t.Errorf("expected totalIssues in JSON report, got:\n%s", out)
	}
}

func TestMarkdownReporter(t *testing.T) {
	issues, summary := sampleData()
	rep := &MarkdownReporter{}

	var buf bytes.Buffer
	if err := rep.Report(&buf, issues, summary); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "## n8nflow-guard Analysis Report") {
		t.Errorf("expected markdown header, got:\n%s", out)
	}
}

func TestHTMLReporter(t *testing.T) {
	issues, summary := sampleData()
	rep := &HTMLReporter{}

	var buf bytes.Buffer
	if err := rep.Report(&buf, issues, summary); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "<title>n8nflow-guard Security Audit Report</title>") {
		t.Errorf("expected HTML title in report, got:\n%s", out)
	}
	if !strings.Contains(out, "SEC-001") {
		t.Errorf("expected SEC-001 in HTML report, got:\n%s", out)
	}
}
