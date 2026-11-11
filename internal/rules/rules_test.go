package rules

import (
	"testing"

	"github.com/alexandrmotologa/n8nflow-guard/internal/config"
	"github.com/alexandrmotologa/n8nflow-guard/internal/parser"
)

func TestDefaultRegistryHas33Rules(t *testing.T) {
	reg := DefaultRegistry()
	all := reg.All()

	if len(all) != 33 {
		t.Fatalf("expected 33 registered rules, found %d", len(all))
	}

	// Verify all 33 IDs exist
	expectedIDs := []string{
		"SEC-001", "SEC-002", "SEC-003", "SEC-004", "SEC-005",
		"SEC-006", "SEC-007", "SEC-008", "SEC-009", "SEC-010",
		"REL-001", "REL-002", "REL-003", "REL-004", "REL-005",
		"REL-006", "REL-007", "REL-008", "REL-009", "REL-010",
		"PERF-001", "PERF-002", "PERF-003", "PERF-004", "PERF-005",
		"BEST-001", "BEST-002", "BEST-003", "BEST-004", "BEST-005", "BEST-006",
		"GOV-001", "GOV-002",
	}

	for _, id := range expectedIDs {
		rule, found := reg.Get(id)
		if !found {
			t.Errorf("rule %s is missing from registry", id)
			continue
		}
		if rule.ID() != id {
			t.Errorf("expected rule ID %s, got %s", id, rule.ID())
		}
		if rule.Name() == "" {
			t.Errorf("rule %s has empty name", id)
		}
		if rule.Description() == "" {
			t.Errorf("rule %s has empty description", id)
		}
		if rule.Remediation() == "" {
			t.Errorf("rule %s has empty remediation", id)
		}
	}
}

func TestEvaluateVulnerableWorkflow(t *testing.T) {
	rawJSON := []byte(`{
		"name": "Vulnerable Pipeline",
		"nodes": [
			{
				"id": "node-1",
				"name": "Webhook",
				"type": "n8n-nodes-base.webhook",
				"typeVersion": 1,
				"parameters": {
					"authentication": "none"
				}
			},
			{
				"id": "node-2",
				"name": "HTTP Request",
				"type": "n8n-nodes-base.httpRequest",
				"typeVersion": 1.0,
				"parameters": {
					"url": "http://api.insecure.com/users",
					"headerParameters": {
						"parameters": [
							{"name": "Authorization", "value": "Bearer ghp_1234567890abcdef1234567890abcdef1234"}
						]
					},
					"options": {
						"allowUnauthorizedCerts": true
					}
				},
				"retryOnFail": false
			}
		],
		"connections": {
			"Webhook": {
				"main": [[{"node": "HTTP Request", "type": "main", "index": 0}]]
			}
		},
		"settings": {
			"saveDataSuccessExecution": "all"
		},
		"pinData": {
			"HTTP Request": [{"user": "alice"}]
		}
	}`)

	workflows, err := parser.ParseWorkflows(rawJSON, "vulnerable.json")
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	wf := workflows[0]
	graph := parser.BuildGraph(wf)
	lineMapper := parser.NewLineMapper(rawJSON)
	cfg := config.DefaultConfig()

	ctx := &RuleContext{
		Workflow:   wf,
		Graph:      graph,
		LineMapper: lineMapper,
		Config:     cfg,
	}

	reg := DefaultRegistry()
	issues := reg.Evaluate(ctx)

	if len(issues) == 0 {
		t.Fatalf("expected multiple issues in vulnerable workflow, got 0")
	}

	ruleHits := make(map[string]bool)
	for _, issue := range issues {
		ruleHits[issue.RuleID] = true
	}

	// Should flag hardcoded secret
	if !ruleHits["SEC-001"] {
		t.Errorf("expected SEC-001 violation for GitHub PAT")
	}
	// Should flag unauthenticated webhook
	if !ruleHits["SEC-002"] {
		t.Errorf("expected SEC-002 violation for webhook authentication: none")
	}
	// Should flag insecure HTTP URL
	if !ruleHits["SEC-004"] {
		t.Errorf("expected SEC-004 violation for http:// URL")
	}
	// Should flag disabled SSL verification
	if !ruleHits["SEC-005"] {
		t.Errorf("expected SEC-005 violation for allowUnauthorizedCerts")
	}
	// Should flag production log data leak
	if !ruleHits["SEC-007"] {
		t.Errorf("expected SEC-007 violation for saveDataSuccessExecution: all")
	}
	// Should flag missing error workflow
	if !ruleHits["REL-001"] {
		t.Errorf("expected REL-001 violation for missing errorWorkflow")
	}
	// Should flag missing retry
	if !ruleHits["REL-002"] {
		t.Errorf("expected REL-002 violation for retryOnFail: false")
	}
	// Should flag default node name
	if !ruleHits["BEST-002"] {
		t.Errorf("expected BEST-002 violation for default node name 'HTTP Request'")
	}
	// Should flag outdated node version
	if !ruleHits["BEST-005"] {
		t.Errorf("expected BEST-005 violation for typeVersion 1.0")
	}
	// Should flag pinned test data
	if !ruleHits["BEST-006"] {
		t.Errorf("expected BEST-006 violation for pinData")
	}
}
