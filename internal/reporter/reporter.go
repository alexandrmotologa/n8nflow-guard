package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// Summary aggregates counts and statistics across a completed scan.
type Summary struct {
	TotalFilesScanned int `json:"totalFilesScanned"`
	TotalWorkflows    int `json:"totalWorkflows"`
	TotalIssues       int `json:"totalIssues"`
	CriticalIssues    int `json:"criticalIssues"`
	HighIssues        int `json:"highIssues"`
	MediumIssues      int `json:"mediumIssues"`
	LowIssues         int `json:"lowIssues"`
	InfoIssues        int `json:"infoIssues"`
	FixableIssues     int `json:"fixableIssues"`
}

// ComputeSummary compiles statistics from a slice of issues.
func ComputeSummary(filesCount, workflowsCount int, issues []models.Issue) Summary {
	s := Summary{
		TotalFilesScanned: filesCount,
		TotalWorkflows:    workflowsCount,
		TotalIssues:       len(issues),
	}

	for _, iss := range issues {
		if iss.Fixable {
			s.FixableIssues++
		}
		switch iss.Severity {
		case models.SeverityCritical:
			s.CriticalIssues++
		case models.SeverityHigh:
			s.HighIssues++
		case models.SeverityMedium:
			s.MediumIssues++
		case models.SeverityLow:
			s.LowIssues++
		case models.SeverityInfo:
			s.InfoIssues++
		}
	}

	return s
}

// Reporter outputs findings to an io.Writer in a specific format.
type Reporter interface {
	Report(w io.Writer, issues []models.Issue, summary Summary) error
}

// NewReporter creates a reporter for the requested format name.
func NewReporter(format string, noColor bool) (Reporter, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "terminal", "text", "console", "":
		return &TerminalReporter{NoColor: noColor}, nil
	case "sarif":
		return &SARIFReporter{}, nil
	case "github-actions", "gha", "github":
		return &GitHubActionsReporter{}, nil
	case "json":
		return &JSONReporter{}, nil
	case "markdown", "md":
		return &MarkdownReporter{}, nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s (supported: terminal, sarif, github-actions, json, markdown)", format)
	}
}
