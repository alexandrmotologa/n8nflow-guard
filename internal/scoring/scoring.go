package scoring

import (
	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
	"github.com/alexandrmotologa/n8nflow-guard/internal/reporter"
)

// SecurityScore represents an evaluation grade and OWASP risk mapping for audited workflows.
type SecurityScore struct {
	Score          int            `json:"score"`
	Grade          string         `json:"grade"`
	OWASPBreakdown map[string]int `json:"owaspBreakdown"`
}

// OWASP Category definitions for API Security Top 10 (2023)
const (
	OWASPBrokenAuth         = "API2:2023 Broken Authentication"
	OWASPPropertyAuth       = "API3:2023 Broken Object Property Level Auth"
	OWASPResourceConsume    = "API4:2023 Unrestricted Resource Consumption"
	OWASPMisconfiguration   = "API8:2023 Security Misconfiguration"
	OWASPUnsafeAPIConsume   = "API10:2023 Unsafe Consumption of APIs"
	OWASPOther              = "Other Governance & Quality"
)

// MapRuleToOWASP maps an n8nflow-guard rule to OWASP API Security Top 10 categories.
func MapRuleToOWASP(ruleID string) string {
	switch ruleID {
	case "SEC-001", "SEC-002", "SEC-003":
		return OWASPBrokenAuth
	case "SEC-007":
		return OWASPPropertyAuth
	case "PERF-001", "PERF-002", "REL-004":
		return OWASPResourceConsume
	case "SEC-004", "SEC-005", "SEC-006", "SEC-008", "SEC-009", "SEC-010":
		return OWASPMisconfiguration
	case "REL-001", "REL-002", "REL-003", "PERF-004":
		return OWASPUnsafeAPIConsume
	default:
		return OWASPOther
	}
}

// CalculateScore computes a 0-100 score and letter grade based on issue counts.
func CalculateScore(summary reporter.Summary, issues []models.Issue) SecurityScore {
	if summary.TotalWorkflows == 0 {
		return SecurityScore{
			Score:          100,
			Grade:          "A+",
			OWASPBreakdown: make(map[string]int),
		}
	}

	deductions := (summary.CriticalIssues * 15) +
		(summary.HighIssues * 8) +
		(summary.MediumIssues * 3) +
		(summary.LowIssues * 1)

	score := 100 - deductions
	if score < 0 {
		score = 0
	}

	grade := "A+"
	switch {
	case score >= 95:
		grade = "A+"
	case score >= 90:
		grade = "A"
	case score >= 80:
		grade = "B"
	case score >= 70:
		grade = "C"
	case score >= 60:
		grade = "D"
	default:
		grade = "F"
	}

	breakdown := make(map[string]int)
	for _, iss := range issues {
		cat := MapRuleToOWASP(iss.RuleID)
		breakdown[cat]++
	}

	return SecurityScore{
		Score:          score,
		Grade:          grade,
		OWASPBreakdown: breakdown,
	}
}
