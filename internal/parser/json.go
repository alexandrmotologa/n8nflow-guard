package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// ParseWorkflowsFromFile reads and decodes one or more n8n workflows from disk.
func ParseWorkflowsFromFile(filePath string) ([]*models.Workflow, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return ParseWorkflows(data, filePath)
}

// ParseWorkflows decodes raw JSON bytes into Workflow structs.
// It supports both single workflow objects and bulk workflow arrays.
func ParseWorkflows(data []byte, filePath string) ([]*models.Workflow, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("empty file content: %s", filePath)
	}

	if trimmed[0] == '[' {
		// Array of workflows
		var workflows []*models.Workflow
		if err := json.Unmarshal(trimmed, &workflows); err != nil {
			return nil, fmt.Errorf("failed to parse workflow array in %s: %w", filePath, err)
		}
		for _, w := range workflows {
			w.FilePath = filePath
			w.RawBytes = trimmed
			normalizeWorkflow(w)
		}
		return workflows, nil
	}

	// Single workflow object
	var workflow models.Workflow
	if err := json.Unmarshal(trimmed, &workflow); err != nil {
		return nil, fmt.Errorf("failed to parse workflow JSON in %s: %w", filePath, err)
	}
	workflow.FilePath = filePath
	workflow.RawBytes = trimmed
	normalizeWorkflow(&workflow)

	return []*models.Workflow{&workflow}, nil
}

// normalizeWorkflow ensures collections and maps are non-nil for safe traversal.
func normalizeWorkflow(w *models.Workflow) {
	if w.Nodes == nil {
		w.Nodes = make([]*models.Node, 0)
	}
	if w.Connections == nil {
		w.Connections = make(map[string]map[string][][]models.ConnectionItem)
	}
	if w.PinData == nil {
		w.PinData = make(map[string]interface{})
	}
	if w.Tags == nil {
		w.Tags = make([]models.WorkflowTag, 0)
	}

	for _, node := range w.Nodes {
		if node.Parameters == nil {
			node.Parameters = make(map[string]interface{})
		}
		if node.Credentials == nil {
			node.Credentials = make(map[string]interface{})
		}
	}
}
