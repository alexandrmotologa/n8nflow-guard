package reporter

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// SARIFLog represents the top-level SARIF v2.1.0 JSON schema.
type SARIFLog struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []SARIFRun `json:"runs"`
}

type SARIFRun struct {
	Tool    SARIFTool     `json:"tool"`
	Results []SARIFResult `json:"results"`
}

type SARIFTool struct {
	Driver SARIFDriver `json:"driver"`
}

type SARIFDriver struct {
	Name           string      `json:"name"`
	InformationURI string      `json:"informationUri"`
	Version        string      `json:"semanticVersion"`
	Rules          []SARIFRule `json:"rules"`
}

type SARIFRule struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	ShortDescription SARIFMultiformatText `json:"shortDescription"`
	HelpURI          string               `json:"helpUri,omitempty"`
}

type SARIFResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   SARIFMessage    `json:"message"`
	Locations []SARIFLocation `json:"locations"`
}

type SARIFMessage struct {
	Text string `json:"text"`
}

type SARIFLocation struct {
	PhysicalLocation SARIFPhysicalLocation `json:"physicalLocation"`
}

type SARIFPhysicalLocation struct {
	ArtifactLocation SARIFArtifactLocation `json:"artifactLocation"`
	Region           SARIFRegion           `json:"region"`
}

type SARIFArtifactLocation struct {
	URI string `json:"uri"`
}

type SARIFRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn,omitempty"`
}

type SARIFMultiformatText struct {
	Text string `json:"text"`
}

// SARIFReporter outputs findings in standard SARIF v2.1.0 JSON format for GitHub Security.
type SARIFReporter struct{}

func (r *SARIFReporter) Report(w io.Writer, issues []models.Issue, summary Summary) error {
	rulesMap := make(map[string]SARIFRule)
	var results []SARIFResult

	for _, iss := range issues {
		if _, exists := rulesMap[iss.RuleID]; !exists {
			rulesMap[iss.RuleID] = SARIFRule{
				ID:   iss.RuleID,
				Name: iss.RuleName,
				ShortDescription: SARIFMultiformatText{
					Text: iss.RuleName,
				},
				HelpURI: "https://github.com/alexandrmotologa/n8nflow-guard/tree/main/rules_docs/" + iss.RuleID + ".md",
			}
		}

		level := "warning"
		if iss.Severity == models.SeverityCritical || iss.Severity == models.SeverityHigh {
			level = "error"
		} else if iss.Severity == models.SeverityLow || iss.Severity == models.SeverityInfo {
			level = "note"
		}

		line := iss.Line
		if line <= 0 {
			line = 1
		}

		normalizedPath := filepath.ToSlash(iss.FilePath)
		// Strip leading ./ if present
		normalizedPath = strings.TrimPrefix(normalizedPath, "./")

		result := SARIFResult{
			RuleID: iss.RuleID,
			Level:  level,
			Message: SARIFMessage{
				Text: iss.Message,
			},
			Locations: []SARIFLocation{
				{
					PhysicalLocation: SARIFPhysicalLocation{
						ArtifactLocation: SARIFArtifactLocation{
							URI: normalizedPath,
						},
						Region: SARIFRegion{
							StartLine:   line,
							StartColumn: iss.Column,
						},
					},
				},
			},
		}

		results = append(results, result)
	}

	var rules []SARIFRule
	for _, rule := range rulesMap {
		rules = append(rules, rule)
	}

	log := SARIFLog{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs: []SARIFRun{
			{
				Tool: SARIFTool{
					Driver: SARIFDriver{
						Name:           "n8nflow-guard",
						InformationURI: "https://github.com/alexandrmotologa/n8nflow-guard",
						Version:        "1.0.0",
						Rules:          rules,
					},
				},
				Results: results,
			},
		},
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(log)
}
