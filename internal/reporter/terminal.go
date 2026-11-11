package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// TerminalReporter prints colorized, human-readable findings and tables to the console.
type TerminalReporter struct {
	NoColor bool
}

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiCyan   = "\033[36m"
	ansiGray   = "\033[90m"
)

func (r *TerminalReporter) color(code, text string) string {
	if r.NoColor {
		return text
	}
	return code + text + ansiReset
}

func (r *TerminalReporter) severityBadge(sev models.Severity) string {
	switch sev {
	case models.SeverityCritical:
		return r.color(ansiRed+ansiBold, "[CRITICAL]")
	case models.SeverityHigh:
		return r.color(ansiRed, "[HIGH]    ")
	case models.SeverityMedium:
		return r.color(ansiYellow, "[MEDIUM]  ")
	case models.SeverityLow:
		return r.color(ansiCyan, "[LOW]     ")
	case models.SeverityInfo:
		return r.color(ansiGray, "[INFO]    ")
	default:
		return fmt.Sprintf("[%s]", sev)
	}
}

// Report renders the findings to the provided writer.
func (r *TerminalReporter) Report(w io.Writer, issues []models.Issue, summary Summary) error {
	if len(issues) == 0 {
		fmt.Fprintf(w, "%s %s\n",
			r.color(ansiGreen+ansiBold, "PASS"),
			r.color(ansiGreen, fmt.Sprintf("No issues found across %d workflow(s) in %d file(s).", summary.TotalWorkflows, summary.TotalFilesScanned)),
		)
		return nil
	}

	// Group issues by FilePath
	fileMap := make(map[string][]models.Issue)
	for _, iss := range issues {
		fileMap[iss.FilePath] = append(fileMap[iss.FilePath], iss)
	}

	var files []string
	for f := range fileMap {
		files = append(files, f)
	}
	sort.Strings(files)

	for _, file := range files {
		fmt.Fprintf(w, "\n%s %s\n", r.color(ansiBold, "FILE:"), file)
		fileIssues := fileMap[file]
		for _, iss := range fileIssues {
			location := fmt.Sprintf("%d", iss.Line)
			if iss.Column > 0 {
				location = fmt.Sprintf("%d:%d", iss.Line, iss.Column)
			}

			fmt.Fprintf(w, "  %s %s:%s  %s  %s\n",
				r.severityBadge(iss.Severity),
				file,
				location,
				r.color(ansiBold, iss.RuleID),
				iss.Message,
			)

			if iss.Remediation != "" {
				fmt.Fprintf(w, "    %s %s\n", r.color(ansiGray, "Hint:"), iss.Remediation)
			}
			if iss.Fixable && iss.FixSuggested != "" {
				fmt.Fprintf(w, "    %s %s (auto-fixable with --fix)\n", r.color(ansiCyan, "Suggested:"), iss.FixSuggested)
			}
		}
	}

	// Summary statistics footer
	fmt.Fprintf(w, "\n%s\n", r.color(ansiBold, "--- Scan Summary ---"))
	fmt.Fprintf(w, "Files scanned:     %d\n", summary.TotalFilesScanned)
	fmt.Fprintf(w, "Workflows checked: %d\n", summary.TotalWorkflows)
	fmt.Fprintf(w, "Total findings:    %d\n", summary.TotalIssues)
	if summary.CriticalIssues > 0 {
		fmt.Fprintf(w, "  Critical:        %s\n", r.color(ansiRed+ansiBold, fmt.Sprintf("%d", summary.CriticalIssues)))
	}
	if summary.HighIssues > 0 {
		fmt.Fprintf(w, "  High:            %s\n", r.color(ansiRed, fmt.Sprintf("%d", summary.HighIssues)))
	}
	if summary.MediumIssues > 0 {
		fmt.Fprintf(w, "  Medium:          %s\n", r.color(ansiYellow, fmt.Sprintf("%d", summary.MediumIssues)))
	}
	if summary.LowIssues > 0 {
		fmt.Fprintf(w, "  Low:             %s\n", r.color(ansiCyan, fmt.Sprintf("%d", summary.LowIssues)))
	}
	if summary.InfoIssues > 0 {
		fmt.Fprintf(w, "  Info:            %d\n", summary.InfoIssues)
	}
	if summary.FixableIssues > 0 {
		fmt.Fprintf(w, "Fixable issues:    %s\n", r.color(ansiGreen, fmt.Sprintf("%d (run 'n8nflow-guard fix' to resolve)", summary.FixableIssues)))
	}

	return nil
}
