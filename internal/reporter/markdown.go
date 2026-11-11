package reporter

import (
	"fmt"
	"io"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// MarkdownReporter generates formatted Markdown tables for PR review comments.
type MarkdownReporter struct{}

func (r *MarkdownReporter) Report(w io.Writer, issues []models.Issue, summary Summary) error {
	fmt.Fprintf(w, "## n8nflow-guard Analysis Report\n\n")

	if len(issues) == 0 {
		fmt.Fprintf(w, "No issues detected across %d workflow(s) in %d file(s).\n", summary.TotalWorkflows, summary.TotalFilesScanned)
		return nil
	}

	fmt.Fprintf(w, "| Metric | Count |\n|---|---|\n")
	fmt.Fprintf(w, "| Files Scanned | %d |\n", summary.TotalFilesScanned)
	fmt.Fprintf(w, "| Workflows Checked | %d |\n", summary.TotalWorkflows)
	fmt.Fprintf(w, "| Total Issues | %d |\n", summary.TotalIssues)
	fmt.Fprintf(w, "| Critical | %d |\n", summary.CriticalIssues)
	fmt.Fprintf(w, "| High | %d |\n", summary.HighIssues)
	fmt.Fprintf(w, "| Medium | %d |\n", summary.MediumIssues)
	fmt.Fprintf(w, "| Low | %d |\n", summary.LowIssues)
	fmt.Fprintf(w, "| Auto-Fixable | %d |\n\n", summary.FixableIssues)

	fmt.Fprintf(w, "### Findings\n\n")
	fmt.Fprintf(w, "| Severity | Rule | Location | Description | Fixable |\n")
	fmt.Fprintf(w, "|---|---|---|---|---|\n")

	for _, iss := range issues {
		fixableStr := "No"
		if iss.Fixable {
			fixableStr = "Yes (`--fix`)"
		}
		loc := fmt.Sprintf("`%s:%d`", iss.FilePath, iss.Line)
		ruleLink := fmt.Sprintf("[%s](https://github.com/alexandrmotologa/n8nflow-guard/tree/main/rules_docs/%s.md)", iss.RuleID, iss.RuleID)

		fmt.Fprintf(w, "| **%s** | %s | %s | %s | %s |\n",
			iss.Severity,
			ruleLink,
			loc,
			iss.Message,
			fixableStr,
		)
	}

	return nil
}
