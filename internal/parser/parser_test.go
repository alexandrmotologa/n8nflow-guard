package parser

import (
	"testing"
)

func TestParseWorkflows(t *testing.T) {
	rawJSON := []byte(`{
		"name": "Order Processing",
		"nodes": [
			{
				"id": "node-1",
				"name": "Webhook Trigger",
				"type": "n8n-nodes-base.webhook",
				"typeVersion": 1,
				"position": [100, 200],
				"parameters": {
					"path": "order",
					"authentication": "none"
				}
			},
			{
				"id": "node-2",
				"name": "HTTP Request",
				"type": "n8n-nodes-base.httpRequest",
				"typeVersion": 4.2,
				"position": [300, 200],
				"parameters": {
					"url": "https://api.payment.com/charge"
				}
			}
		],
		"connections": {
			"Webhook Trigger": {
				"main": [
					[
						{
							"node": "HTTP Request",
							"type": "main",
							"index": 0
						}
					]
				]
			}
		},
		"settings": {
			"executionTimeout": 3600
		}
	}`)

	workflows, err := ParseWorkflows(rawJSON, "order_workflow.json")
	if err != nil {
		t.Fatalf("unexpected error parsing workflow: %v", err)
	}

	if len(workflows) != 1 {
		t.Fatalf("expected 1 workflow, got %d", len(workflows))
	}

	wf := workflows[0]
	if wf.Name != "Order Processing" {
		t.Errorf("expected workflow name 'Order Processing', got '%s'", wf.Name)
	}

	if len(wf.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(wf.Nodes))
	}

	g := BuildGraph(wf)
	if len(g.Adjacency["Webhook Trigger"]) != 1 {
		t.Errorf("expected 1 outgoing connection from Webhook Trigger, got %d", len(g.Adjacency["Webhook Trigger"]))
	}

	cycles := g.DetectCycles()
	if len(cycles) != 0 {
		t.Errorf("expected 0 cycles, got %d", len(cycles))
	}

	orphans := g.FindOrphanNodes()
	if len(orphans) != 0 {
		t.Errorf("expected 0 orphans, got %d", len(orphans))
	}
}

func TestDetectCycles(t *testing.T) {
	rawJSON := []byte(`{
		"name": "Cyclic Workflow",
		"nodes": [
			{"id": "n1", "name": "Node A", "type": "n8n-nodes-base.set"},
			{"id": "n2", "name": "Node B", "type": "n8n-nodes-base.set"},
			{"id": "n3", "name": "Node C", "type": "n8n-nodes-base.set"}
		],
		"connections": {
			"Node A": {"main": [[{"node": "Node B", "type": "main", "index": 0}]]},
			"Node B": {"main": [[{"node": "Node C", "type": "main", "index": 0}]]},
			"Node C": {"main": [[{"node": "Node A", "type": "main", "index": 0}]]}
		}
	}`)

	workflows, err := ParseWorkflows(rawJSON, "cyclic.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	g := BuildGraph(workflows[0])
	cycles := g.DetectCycles()
	if len(cycles) == 0 {
		t.Fatalf("expected cycles to be detected, got 0")
	}

	cycle := cycles[0]
	if len(cycle) != 4 || cycle[0] != cycle[len(cycle)-1] {
		t.Errorf("unexpected cycle path: %v", cycle)
	}
}

func TestExpressionValidation(t *testing.T) {
	valid := "Bearer {{ $json.token }}"
	if err := ValidateExpressionSyntax(valid); err != nil {
		t.Errorf("expected valid expression, got error: %v", err)
	}

	invalid := "Bearer {{ $json.token"
	if err := ValidateExpressionSyntax(invalid); err == nil {
		t.Errorf("expected syntax error for unclosed brace, got nil")
	}

	refs := ExtractReferencedNodes("{{ $('Database Query').item.json.id }} + {{ $node['Legacy API'].data }}")
	if len(refs) != 2 {
		t.Fatalf("expected 2 referenced nodes, got %d", len(refs))
	}
	if refs[0] != "Database Query" || refs[1] != "Legacy API" {
		t.Errorf("unexpected extracted nodes: %v", refs)
	}

	if !IsEnvExpression("={{ $env.SECRET_KEY }}") {
		t.Errorf("expected IsEnvExpression to return true for ={{ $env.SECRET_KEY }}")
	}
	if IsEnvExpression("sk-proj-1234567890") {
		t.Errorf("expected IsEnvExpression to return false for raw secret")
	}
}

func TestLineMapper(t *testing.T) {
	content := []byte(`{
  "name": "Test",
  "nodes": [
    {
      "name": "First Node",
      "type": "test"
    },
    {
      "name": "Second Node",
      "type": "test",
      "parameters": {
        "apiKey": "12345"
      }
    }
  ]
}`)

	lm := NewLineMapper(content)
	line := lm.FindNodeLine("Second Node", "")
	if line != 9 {
		t.Errorf("expected node line 9, got %d", line)
	}

	fieldLine := lm.FindFieldLine("Second Node", "apiKey")
	if fieldLine != 12 {
		t.Errorf("expected field line 12, got %d", fieldLine)
	}
}
