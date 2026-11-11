package reporter

import (
	"encoding/json"
	"io"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// JSONReport represents the structured machine-readable payload.
type JSONReport struct {
	Summary Summary        `json:"summary"`
	Issues  []models.Issue `json:"issues"`
}

// JSONReporter emits raw structured JSON output.
type JSONReporter struct{}

func (r *JSONReporter) Report(w io.Writer, issues []models.Issue, summary Summary) error {
	report := JSONReport{
		Summary: summary,
		Issues:  issues,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
