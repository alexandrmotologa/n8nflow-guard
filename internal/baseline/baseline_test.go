package baseline

import (
	"path/filepath"
	"testing"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

func TestBaselineGenerateAndFilter(t *testing.T) {
	tmpDir := t.TempDir()
	baselinePath := filepath.Join(tmpDir, "baseline.json")

	initialIssues := []models.Issue{
		{
			RuleID:   "SEC-002",
			RuleName: "unauthenticated-webhook",
			FilePath: "workflows/legacy.json",
			NodeName: "Webhook",
			Field:    "authentication",
		},
		{
			RuleID:   "REL-001",
			RuleName: "missing-error-workflow",
			FilePath: "workflows/legacy.json",
			Field:    "errorWorkflow",
		},
	}

	if err := GenerateBaseline(initialIssues, baselinePath); err != nil {
		t.Fatalf("failed generating baseline: %v", err)
	}

	bFile, index, err := LoadBaseline(baselinePath)
	if err != nil {
		t.Fatalf("failed loading baseline: %v", err)
	}

	if bFile.TotalIssues != 2 {
		t.Errorf("expected 2 total issues in baseline, got %d", bFile.TotalIssues)
	}

	currentIssues := []models.Issue{
		initialIssues[0], // Already in baseline
		initialIssues[1], // Already in baseline
		{
			RuleID:   "SEC-001",
			RuleName: "hardcoded-secret",
			FilePath: "workflows/legacy.json",
			NodeName: "HTTP Request",
			Field:    "apiKey",
		}, // Brand new issue!
	}

	newIssues := FilterNewIssues(currentIssues, index)
	if len(newIssues) != 1 {
		t.Fatalf("expected 1 new issue, got %d", len(newIssues))
	}
	if newIssues[0].RuleID != "SEC-001" {
		t.Errorf("expected new issue to be SEC-001, got %s", newIssues[0].RuleID)
	}
}
