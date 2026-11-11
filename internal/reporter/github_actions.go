package reporter

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// GitHubActionsReporter emits workflow commands that GitHub Actions translates into PR annotations.
type GitHubActionsReporter struct{}

func (r *GitHubActionsReporter) Report(w io.Writer, issues []models.Issue, summary Summary) error {
	for _, iss := range issues {
		command := "warning"
		if iss.Severity == models.SeverityCritical || iss.Severity == models.SeverityHigh {
			command = "error"
		} else if iss.Severity == models.SeverityLow || iss.Severity == models.SeverityInfo {
			command = "notice"
		}

		line := iss.Line
		if line <= 0 {
			line = 1
		}

		normalizedPath := filepath.ToSlash(iss.FilePath)
		normalizedPath = strings.TrimPrefix(normalizedPath, "./")

		title := fmt.Sprintf("[%s] %s", iss.RuleID, iss.RuleName)
		msg := iss.Message
		if iss.Remediation != "" {
			msg = fmt.Sprintf("%s (Hint: %s)", msg, iss.Remediation)
		}

		// ::error file={file},line={line},title={title}::{message}
		fmt.Fprintf(w, "::%s file=%s,line=%d,title=%s::%s\n",
			command,
			normalizedPath,
			line,
			escapeGHAProperty(title),
			escapeGHAMessage(msg),
		)
	}

	return nil
}

func escapeGHAProperty(s string) string {
	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "\r", "%0D")
	s = strings.ReplaceAll(s, "\n", "%0A")
	s = strings.ReplaceAll(s, ":", "%3A")
	s = strings.ReplaceAll(s, ",", "%2C")
	return s
}

func escapeGHAMessage(s string) string {
	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "\r", "%0D")
	s = strings.ReplaceAll(s, "\n", "%0A")
	return s
}
