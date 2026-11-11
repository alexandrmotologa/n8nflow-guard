package scoring

import (
	"testing"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
	"github.com/alexandrmotologa/n8nflow-guard/internal/reporter"
)

func TestCalculateScore(t *testing.T) {
	// Clean run
	cleanSummary := reporter.Summary{TotalWorkflows: 1, TotalFilesScanned: 1}
	cleanScore := CalculateScore(cleanSummary, nil)
	if cleanScore.Score != 100 || cleanScore.Grade != "A+" {
		t.Errorf("expected 100/A+ for clean run, got %d/%s", cleanScore.Score, cleanScore.Grade)
	}

	// Run with critical issue
	issues := []models.Issue{
		{RuleID: "SEC-001", Severity: models.SeverityCritical},
		{RuleID: "SEC-002", Severity: models.SeverityHigh},
	}
	summary := reporter.Summary{
		TotalWorkflows: 1,
		CriticalIssues: 1,
		HighIssues:     1,
		TotalIssues:    2,
	}
	score := CalculateScore(summary, issues)
	// 100 - 15 - 8 = 77 -> C
	if score.Score != 77 || score.Grade != "C" {
		t.Errorf("expected 77/C, got %d/%s", score.Score, score.Grade)
	}

	if score.OWASPBreakdown[OWASPBrokenAuth] != 2 {
		t.Errorf("expected 2 Broken Auth issues in OWASP breakdown, got %d", score.OWASPBreakdown[OWASPBrokenAuth])
	}
}
