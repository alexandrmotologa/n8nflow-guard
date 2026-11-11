package engine

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

func TestRunScan(t *testing.T) {
	tmpDir := t.TempDir()

	validWorkflow := []byte(`{
		"name": "Production Order Flow",
		"nodes": [
			{
				"name": "Schedule Trigger",
				"type": "n8n-nodes-base.scheduleTrigger",
				"typeVersion": 1,
				"position": [100, 200]
			},
			{
				"name": "Process Order",
				"type": "n8n-nodes-base.code",
				"typeVersion": 2,
				"notes": "Handles business transformation",
				"position": [300, 200]
			}
		],
		"connections": {
			"Schedule Trigger": {
				"main": [[{"node": "Process Order", "type": "main", "index": 0}]]
			}
		},
		"settings": {
			"errorWorkflow": "wf-error-alert-id",
			"executionTimeout": 3600
		},
		"tags": [{"id": "t1", "name": "production"}]
	}`)

	vulnerableWorkflow := []byte(`{
		"name": "Vulnerable Flow",
		"nodes": [
			{
				"name": "HTTP Request",
				"type": "n8n-nodes-base.httpRequest",
				"parameters": {
					"apiKey": "sk-proj-1234567890abcdef1234567890abcdef"
				}
			}
		]
	}`)

	validFile := filepath.Join(tmpDir, "valid.json")
	vulnFile := filepath.Join(tmpDir, "vuln.json")

	if err := os.WriteFile(validFile, validWorkflow, 0644); err != nil {
		t.Fatalf("failed to write valid workflow: %v", err)
	}
	if err := os.WriteFile(vulnFile, vulnerableWorkflow, 0644); err != nil {
		t.Fatalf("failed to write vulnerable workflow: %v", err)
	}

	opts := ScanOptions{
		Paths:       []string{tmpDir},
		MinSeverity: models.SeverityLow,
		FailOn:      models.SeverityHigh,
	}

	res, err := RunScan(opts)
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}

	if res.Summary.TotalFilesScanned != 2 {
		t.Errorf("expected 2 files scanned, got %d", res.Summary.TotalFilesScanned)
	}

	if !res.HasFailures {
		t.Errorf("expected scan to fail due to critical secret in vuln.json")
	}

	// Test scanning only valid workflow
	optsValid := ScanOptions{
		Paths:       []string{validFile},
		MinSeverity: models.SeverityHigh,
		FailOn:      models.SeverityHigh,
	}
	resValid, err := RunScan(optsValid)
	if err != nil {
		t.Fatalf("scan error on valid file: %v", err)
	}
	if resValid.HasFailures {
		t.Errorf("expected valid workflow to pass without high severity failures")
	}
}
