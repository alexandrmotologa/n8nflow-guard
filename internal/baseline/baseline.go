package baseline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// Entry records a single baseline issue fingerprint.
type Entry struct {
	Fingerprint string `json:"fingerprint"`
	RuleID      string `json:"ruleId"`
	RuleName    string `json:"ruleName"`
	FilePath    string `json:"filePath"`
	NodeName    string `json:"nodeName,omitempty"`
	Field       string `json:"field,omitempty"`
}

// File holds the serialized baseline representation.
type File struct {
	Version     string    `json:"version"`
	GeneratedAt time.Time `json:"generatedAt"`
	TotalIssues int       `json:"totalIssues"`
	Issues      []Entry   `json:"issues"`
}

// ComputeFingerprint calculates a stable content hash for an issue independent of line numbers.
func ComputeFingerprint(iss models.Issue) string {
	normalizedPath := filepath.ToSlash(iss.FilePath)
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%s:%s", iss.RuleID, normalizedPath, iss.NodeName, iss.Field)))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// GenerateBaseline creates a baseline file containing fingerprints of all supplied issues.
func GenerateBaseline(issues []models.Issue, outputPath string) error {
	entries := make([]Entry, 0, len(issues))
	for _, iss := range issues {
		entries = append(entries, Entry{
			Fingerprint: ComputeFingerprint(iss),
			RuleID:      iss.RuleID,
			RuleName:    iss.RuleName,
			FilePath:    filepath.ToSlash(iss.FilePath),
			NodeName:    iss.NodeName,
			Field:       iss.Field,
		})
	}

	bFile := File{
		Version:     "1.0",
		GeneratedAt: time.Now().UTC(),
		TotalIssues: len(entries),
		Issues:      entries,
	}

	data, err := json.MarshalIndent(bFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed serializing baseline: %w", err)
	}

	return os.WriteFile(outputPath, append(data, '\n'), 0644)
}

// LoadBaseline reads and indexes a baseline file from disk.
func LoadBaseline(baselinePath string) (*File, map[string]bool, error) {
	data, err := os.ReadFile(baselinePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed reading baseline file %s: %w", baselinePath, err)
	}

	var bFile File
	if err := json.Unmarshal(data, &bFile); err != nil {
		return nil, nil, fmt.Errorf("failed parsing baseline file %s: %w", baselinePath, err)
	}

	index := make(map[string]bool)
	for _, entry := range bFile.Issues {
		index[entry.Fingerprint] = true
	}

	return &bFile, index, nil
}

// FilterNewIssues removes issues already recorded in the baseline index.
func FilterNewIssues(issues []models.Issue, baselineIndex map[string]bool) []models.Issue {
	if len(baselineIndex) == 0 {
		return issues
	}

	var newIssues []models.Issue
	for _, iss := range issues {
		fp := ComputeFingerprint(iss)
		if !baselineIndex[fp] {
			newIssues = append(newIssues, iss)
		}
	}

	return newIssues
}
