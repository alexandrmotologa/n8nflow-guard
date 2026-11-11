package models

import (
	"fmt"
	"strings"
)

// Severity indicates the impact level of an identified issue.
type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityMedium   Severity = "MEDIUM"
	SeverityLow      Severity = "LOW"
	SeverityInfo     Severity = "INFO"
)

// Rank provides an ordinal rank for severity comparisons (higher is more severe).
func (s Severity) Rank() int {
	switch strings.ToUpper(string(s)) {
	case string(SeverityCritical):
		return 5
	case string(SeverityHigh):
		return 4
	case string(SeverityMedium):
		return 3
	case string(SeverityLow):
		return 2
	case string(SeverityInfo):
		return 1
	default:
		return 0
	}
}

// String returns the uppercase string representation.
func (s Severity) String() string {
	return string(s)
}

// ParseSeverity converts a string into a recognized Severity enum.
func ParseSeverity(raw string) (Severity, error) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "CRITICAL":
		return SeverityCritical, nil
	case "HIGH":
		return SeverityHigh, nil
	case "MEDIUM":
		return SeverityMedium, nil
	case "LOW":
		return SeverityLow, nil
	case "INFO":
		return SeverityInfo, nil
	default:
		return "", fmt.Errorf("unknown severity: %s", raw)
	}
}

// Category organizes rules into domain groups.
type Category string

const (
	CategorySecurity     Category = "Security"
	CategoryReliability  Category = "Reliability"
	CategoryPerformance  Category = "Performance"
	CategoryBestPractice Category = "BestPractice"
	CategoryGovernance   Category = "Governance"
)

// Issue represents a single linter violation identified during static analysis.
type Issue struct {
	RuleID       string   `json:"ruleId"`
	RuleName     string   `json:"ruleName"`
	Category     Category `json:"category"`
	Severity     Severity `json:"severity"`
	Message      string   `json:"message"`
	FilePath     string   `json:"filePath"`
	Line         int      `json:"line"`
	Column       int      `json:"column"`
	NodeID       string   `json:"nodeId,omitempty"`
	NodeName     string   `json:"nodeName,omitempty"`
	Field        string   `json:"field,omitempty"`
	Remediation  string   `json:"remediation,omitempty"`
	Fixable      bool     `json:"fixable"`
	FixSuggested string   `json:"fixSuggested,omitempty"`
}
