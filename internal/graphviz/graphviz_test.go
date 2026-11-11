package graphviz

import (
	"strings"
	"testing"

	"github.com/alexandrmotologa/n8nflow-guard/internal/parser"
)

func TestExportMermaidAndDOT(t *testing.T) {
	rawJSON := []byte(`{
		"name": "Order Pipeline",
		"nodes": [
			{"name": "Webhook Trigger", "type": "n8n-nodes-base.webhook"},
			{"name": "Process Data", "type": "n8n-nodes-base.code"}
		],
		"connections": {
			"Webhook Trigger": {
				"main": [[{"node": "Process Data", "type": "main", "index": 0}]]
			}
		}
	}`)

	workflows, err := parser.ParseWorkflows(rawJSON, "order.json")
	if err != nil {
		t.Fatalf("failed parsing: %v", err)
	}

	g := parser.BuildGraph(workflows[0])

	// Test Mermaid export
	mermaid := ExportMermaid(g)
	if !strings.Contains(mermaid, "flowchart TD") {
		t.Errorf("expected flowchart TD, got:\n%s", mermaid)
	}
	if !strings.Contains(mermaid, "Webhook_Trigger") || !strings.Contains(mermaid, "Process_Data") {
		t.Errorf("expected sanitized node IDs in mermaid, got:\n%s", mermaid)
	}
	if !strings.Contains(mermaid, "Webhook_Trigger --> Process_Data") {
		t.Errorf("expected connection arrow in mermaid, got:\n%s", mermaid)
	}

	// Test DOT export
	dot := ExportDOT(g)
	if !strings.Contains(dot, "digraph Workflow") {
		t.Errorf("expected digraph Workflow in DOT, got:\n%s", dot)
	}
	if !strings.Contains(dot, "Webhook_Trigger -> Process_Data;") {
		t.Errorf("expected edge in DOT, got:\n%s", dot)
	}

	// Test ASCII export
	ascii := ExportASCII(g)
	if !strings.Contains(ascii, "[Trigger] Webhook Trigger") {
		t.Errorf("expected trigger in ASCII, got:\n%s", ascii)
	}
	if !strings.Contains(ascii, "-> to:   Process Data") {
		t.Errorf("expected outbound connection in ASCII, got:\n%s", ascii)
	}
}
