package rules

import (
	"fmt"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// -------------------------------------------------------------
// PERF-001: infinite-loop-cycle
// -------------------------------------------------------------
type RulePERF001 struct{}

func (r *RulePERF001) ID() string             { return "PERF-001" }
func (r *RulePERF001) Name() string           { return "infinite-loop-cycle" }
func (r *RulePERF001) Category() models.Category { return models.CategoryPerformance }
func (r *RulePERF001) Severity() models.Severity { return models.SeverityHigh }
func (r *RulePERF001) Description() string {
	return "Workflow graph contains a directed loop/cycle that risks infinite execution."
}
func (r *RulePERF001) Remediation() string {
	return "Break the cycle or introduce a Split In Batches node with an explicit loop counter."
}
func (r *RulePERF001) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	cycles := ctx.Graph.DetectCycles()
	for _, cycle := range cycles {
		if len(cycle) == 0 {
			continue
		}
		firstNode := cycle[0]
		line := ctx.LineMapper.FindNodeLine(firstNode, "")
		cycleStr := strings.Join(cycle, " -> ")

		issues = append(issues, models.Issue{
			RuleID:      r.ID(),
			RuleName:    r.Name(),
			Category:    r.Category(),
			Severity:    r.Severity(),
			Message:     fmt.Sprintf("Detected cyclic loop in workflow: %s", cycleStr),
			FilePath:    ctx.Workflow.FilePath,
			Line:        line,
			NodeName:    firstNode,
			Field:       "connections",
			Remediation: r.Remediation(),
		})
	}

	return issues
}

// -------------------------------------------------------------
// PERF-002: unbounded-batch-processing
// -------------------------------------------------------------
type RulePERF002 struct{}

func (r *RulePERF002) ID() string             { return "PERF-002" }
func (r *RulePERF002) Name() string           { return "unbounded-batch-processing" }
func (r *RulePERF002) Category() models.Category { return models.CategoryPerformance }
func (r *RulePERF002) Severity() models.Severity { return models.SeverityMedium }
func (r *RulePERF002) Description() string {
	return "Split In Batches node does not define a batch size or has batchSize <= 0."
}
func (r *RulePERF002) Remediation() string {
	return "Configure batchSize to an appropriate finite number (e.g. 50 to 500)."
}
func (r *RulePERF002) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if !strings.Contains(strings.ToLower(node.Type), "splitinbatches") {
			continue
		}

		batchSize, ok := node.Parameters["batchSize"]
		invalid := false
		if !ok {
			invalid = true
		} else {
			switch v := batchSize.(type) {
			case float64:
				if v <= 0 {
					invalid = true
				}
			case int:
				if v <= 0 {
					invalid = true
				}
			}
		}

		if invalid {
			line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Split In Batches node '%s' lacks valid batchSize configuration", node.Name),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "parameters.batchSize",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// PERF-003: excessive-execution-logging
// -------------------------------------------------------------
type RulePERF003 struct{}

func (r *RulePERF003) ID() string             { return "PERF-003" }
func (r *RulePERF003) Name() string           { return "excessive-execution-logging" }
func (r *RulePERF003) Category() models.Category { return models.CategoryPerformance }
func (r *RulePERF003) Severity() models.Severity { return models.SeverityMedium }
func (r *RulePERF003) Description() string {
	return "Workflow persists execution progress step by step, degrading throughput on production n8n instances."
}
func (r *RulePERF003) Remediation() string {
	return "Set saveExecutionProgress to false unless actively troubleshooting."
}
func (r *RulePERF003) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	if ctx.Workflow.Settings.SaveExecutionProgress {
		line := ctx.LineMapper.FindFieldLine("settings", "saveExecutionProgress")
		issues = append(issues, models.Issue{
			RuleID:      r.ID(),
			RuleName:    r.Name(),
			Category:    r.Category(),
			Severity:    r.Severity(),
			Message:     "Workflow has saveExecutionProgress enabled, which adds database write overhead to every node execution",
			FilePath:    ctx.Workflow.FilePath,
			Line:        line,
			Field:       "settings.saveExecutionProgress",
			Remediation: r.Remediation(),
		})
	}

	return issues
}

// -------------------------------------------------------------
// PERF-004: n-plus-one-request-pattern
// -------------------------------------------------------------
type RulePERF004 struct{}

func (r *RulePERF004) ID() string             { return "PERF-004" }
func (r *RulePERF004) Name() string           { return "n-plus-one-request-pattern" }
func (r *RulePERF004) Category() models.Category { return models.CategoryPerformance }
func (r *RulePERF004) Severity() models.Severity { return models.SeverityMedium }
func (r *RulePERF004) Description() string {
	return "HTTP Request node executes in a per-item loop without batching or request aggregation."
}
func (r *RulePERF004) Remediation() string {
	return "Use batch API endpoints or aggregate items into bulk requests where possible."
}
func (r *RulePERF004) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	cycles := ctx.Graph.DetectCycles()
	for _, cycle := range cycles {
		for _, nodeName := range cycle {
			node := ctx.Workflow.FindNodeByName(nodeName)
			if node != nil && strings.Contains(strings.ToLower(node.Type), "httprequest") {
				line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
				issues = append(issues, models.Issue{
					RuleID:      r.ID(),
					RuleName:    r.Name(),
					Category:    r.Category(),
					Severity:    r.Severity(),
					Message:     fmt.Sprintf("HTTP Request node '%s' is placed inside a loop, causing potential N+1 API calls", node.Name),
					FilePath:    ctx.Workflow.FilePath,
					Line:        line,
					NodeID:      node.ID,
					NodeName:    node.Name,
					Field:       "connections",
					Remediation: r.Remediation(),
				})
			}
		}
	}

	return issues
}

// -------------------------------------------------------------
// PERF-005: excessive-polling-frequency
// -------------------------------------------------------------
type RulePERF005 struct{}

func (r *RulePERF005) ID() string             { return "PERF-005" }
func (r *RulePERF005) Name() string           { return "excessive-polling-frequency" }
func (r *RulePERF005) Category() models.Category { return models.CategoryPerformance }
func (r *RulePERF005) Severity() models.Severity { return models.SeverityLow }
func (r *RulePERF005) Description() string {
	return "Interval or polling trigger executes more frequently than once every 10 seconds."
}
func (r *RulePERF005) Remediation() string {
	return "Increase polling interval to >= 60 seconds or switch to webhook triggers."
}
func (r *RulePERF005) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if !strings.Contains(strings.ToLower(node.Type), "interval") {
			continue
		}

		interval, ok := node.Parameters["interval"].(float64)
		unit, _ := node.Parameters["unit"].(string)

		if ok && unit == "seconds" && interval < 10 {
			line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Interval trigger node '%s' polls every %.0f seconds, which can cause excessive load", node.Name, interval),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "parameters.interval",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}
