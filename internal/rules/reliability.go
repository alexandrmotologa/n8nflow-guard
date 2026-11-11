package rules

import (
	"fmt"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
	"github.com/alexandrmotologa/n8nflow-guard/internal/parser"
)

// -------------------------------------------------------------
// REL-001: missing-error-workflow
// -------------------------------------------------------------
type RuleREL001 struct{}

func (r *RuleREL001) ID() string             { return "REL-001" }
func (r *RuleREL001) Name() string           { return "missing-error-workflow" }
func (r *RuleREL001) Category() models.Category { return models.CategoryReliability }
func (r *RuleREL001) Severity() models.Severity { return models.SeverityHigh }
func (r *RuleREL001) Description() string {
	return "Workflow has no global error handler workflow assigned, risking unnotified silent failures."
}
func (r *RuleREL001) Remediation() string {
	return "Set errorWorkflow in workflow settings to designate a dedicated incident response workflow."
}
func (r *RuleREL001) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	if strings.TrimSpace(ctx.Workflow.Settings.ErrorWorkflow) == "" {
		line := ctx.LineMapper.FindFieldLine("settings", "errorWorkflow")
		issues = append(issues, models.Issue{
			RuleID:      r.ID(),
			RuleName:    r.Name(),
			Category:    r.Category(),
			Severity:    r.Severity(),
			Message:     "Workflow does not configure an errorWorkflow handler in settings",
			FilePath:    ctx.Workflow.FilePath,
			Line:        line,
			Field:       "settings.errorWorkflow",
			Remediation: r.Remediation(),
		})
	}

	return issues
}

// -------------------------------------------------------------
// REL-002: missing-retry-on-fail
// -------------------------------------------------------------
type RuleREL002 struct{}

func (r *RuleREL002) ID() string             { return "REL-002" }
func (r *RuleREL002) Name() string           { return "missing-retry-on-fail" }
func (r *RuleREL002) Category() models.Category { return models.CategoryReliability }
func (r *RuleREL002) Severity() models.Severity { return models.SeverityMedium }
func (r *RuleREL002) Description() string {
	return "External network call node lacks retry-on-failure configuration for handling transient errors."
}
func (r *RuleREL002) Remediation() string {
	return "Enable retryOnFail and set maxTries to 3 or more with exponential backoff."
}
func (r *RuleREL002) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if !strings.Contains(strings.ToLower(node.Type), "httprequest") {
			continue
		}
		if node.Disabled && (ctx.Config == nil || !ctx.Config.IncludeDisabled) {
			continue
		}

		if !node.RetryOnFail || node.MaxTries <= 1 {
			line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("HTTP Request node '%s' does not enable retry-on-failure", node.Name),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "retryOnFail",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// REL-003: unhandled-failure-continue
// -------------------------------------------------------------
type RuleREL003 struct{}

func (r *RuleREL003) ID() string             { return "REL-003" }
func (r *RuleREL003) Name() string           { return "unhandled-failure-continue" }
func (r *RuleREL003) Category() models.Category { return models.CategoryReliability }
func (r *RuleREL003) Severity() models.Severity { return models.SeverityMedium }
func (r *RuleREL003) Description() string {
	return "Node continues on failure without dedicated error output connections or handling."
}
func (r *RuleREL003) Remediation() string {
	return "Ensure nodes using continueOnFail inspect error states in downstream nodes or use error triggers."
}
func (r *RuleREL003) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if !node.ContinueOnFail {
			continue
		}
		if node.Disabled && (ctx.Config == nil || !ctx.Config.IncludeDisabled) {
			continue
		}

		hasErrorBranch := ctx.Graph.HasDownstreamErrorHandling(node.Name)
		if !hasErrorBranch {
			line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Node '%s' has continueOnFail=true but does not connect to an error branch", node.Name),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "continueOnFail",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// REL-004: missing-execution-timeout
// -------------------------------------------------------------
type RuleREL004 struct{}

func (r *RuleREL004) ID() string             { return "REL-004" }
func (r *RuleREL004) Name() string           { return "missing-execution-timeout" }
func (r *RuleREL004) Category() models.Category { return models.CategoryReliability }
func (r *RuleREL004) Severity() models.Severity { return models.SeverityMedium }
func (r *RuleREL004) Description() string {
	return "Workflow has no explicit execution timeout configured, risking hanging workers on hung sockets."
}
func (r *RuleREL004) Remediation() string {
	return "Specify an executionTimeout (in seconds) under workflow settings (e.g. 300 to 3600 seconds)."
}
func (r *RuleREL004) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	if ctx.Workflow.Settings.ExecutionTimeout <= 0 {
		line := ctx.LineMapper.FindFieldLine("settings", "executionTimeout")
		issues = append(issues, models.Issue{
			RuleID:      r.ID(),
			RuleName:    r.Name(),
			Category:    r.Category(),
			Severity:    r.Severity(),
			Message:     "Workflow does not configure an executionTimeout limit in settings",
			FilePath:    ctx.Workflow.FilePath,
			Line:        line,
			Field:       "settings.executionTimeout",
			Remediation: r.Remediation(),
		})
	}

	return issues
}

// -------------------------------------------------------------
// REL-005: malformed-expression
// -------------------------------------------------------------
type RuleREL005 struct{}

func (r *RuleREL005) ID() string             { return "REL-005" }
func (r *RuleREL005) Name() string           { return "malformed-expression" }
func (r *RuleREL005) Category() models.Category { return models.CategoryReliability }
func (r *RuleREL005) Severity() models.Severity { return models.SeverityHigh }
func (r *RuleREL005) Description() string {
	return "Parameter contains an invalid or unclosed n8n expression syntax."
}
func (r *RuleREL005) Remediation() string {
	return "Correct mismatched braces or parenthesis in the expression."
}
func (r *RuleREL005) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		walkParameters(node.Parameters, "", func(keyPath string, val string) {
			if strings.Contains(val, "{{") || strings.Contains(val, "}}") {
				if err := parser.ValidateExpressionSyntax(val); err != nil {
					line := ctx.LineMapper.FindFieldLine(node.Name, keyPath)
					issues = append(issues, models.Issue{
						RuleID:      r.ID(),
						RuleName:    r.Name(),
						Category:    r.Category(),
						Severity:    r.Severity(),
						Message:     fmt.Sprintf("Node '%s' parameter '%s' has syntax error: %v", node.Name, keyPath, err),
						FilePath:    ctx.Workflow.FilePath,
						Line:        line,
						NodeID:      node.ID,
						NodeName:    node.Name,
						Field:       keyPath,
						Remediation: r.Remediation(),
					})
				}
			}
		})
	}

	return issues
}

// -------------------------------------------------------------
// REL-006: reference-to-missing-node
// -------------------------------------------------------------
type RuleREL006 struct{}

func (r *RuleREL006) ID() string             { return "REL-006" }
func (r *RuleREL006) Name() string           { return "reference-to-missing-node" }
func (r *RuleREL006) Category() models.Category { return models.CategoryReliability }
func (r *RuleREL006) Severity() models.Severity { return models.SeverityHigh }
func (r *RuleREL006) Description() string {
	return "Expression references a node name that does not exist in the current workflow."
}
func (r *RuleREL006) Remediation() string {
	return "Update the expression to reference an existing active node in the workflow canvas."
}
func (r *RuleREL006) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		walkParameters(node.Parameters, "", func(keyPath string, val string) {
			if !strings.Contains(val, "{{") {
				return
			}
			referenced := parser.ExtractReferencedNodes(val)
			for _, ref := range referenced {
				if ctx.Workflow.FindNodeByName(ref) == nil && ctx.Workflow.FindNodeByID(ref) == nil {
					line := ctx.LineMapper.FindFieldLine(node.Name, keyPath)
					issues = append(issues, models.Issue{
						RuleID:      r.ID(),
						RuleName:    r.Name(),
						Category:    r.Category(),
						Severity:    r.Severity(),
						Message:     fmt.Sprintf("Node '%s' references nonexistent node '%s' in parameter '%s'", node.Name, ref, keyPath),
						FilePath:    ctx.Workflow.FilePath,
						Line:        line,
						NodeID:      node.ID,
						NodeName:    node.Name,
						Field:       keyPath,
						Remediation: r.Remediation(),
					})
				}
			}
		})
	}

	return issues
}

// -------------------------------------------------------------
// REL-007: disabled-critical-node
// -------------------------------------------------------------
type RuleREL007 struct{}

func (r *RuleREL007) ID() string             { return "REL-007" }
func (r *RuleREL007) Name() string           { return "disabled-critical-node" }
func (r *RuleREL007) Category() models.Category { return models.CategoryReliability }
func (r *RuleREL007) Severity() models.Severity { return models.SeverityMedium }
func (r *RuleREL007) Description() string {
	return "A disabled node is positioned between active connected nodes, breaking the execution chain."
}
func (r *RuleREL007) Remediation() string {
	return "Either re-enable the node, bypass its connections, or delete it from the canvas."
}
func (r *RuleREL007) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if !node.Disabled {
			continue
		}

		inCount := len(ctx.Graph.Inbound[node.Name])
		outCount := len(ctx.Graph.Adjacency[node.Name])

		// If a disabled node has incoming and outgoing connections, it interrupts a flow
		if inCount > 0 && outCount > 0 {
			line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Disabled node '%s' sits in the active execution path between other nodes", node.Name),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "disabled",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// REL-008: unhandled-webhook-error
// -------------------------------------------------------------
type RuleREL008 struct{}

func (r *RuleREL008) ID() string             { return "REL-008" }
func (r *RuleREL008) Name() string           { return "unhandled-webhook-error" }
func (r *RuleREL008) Category() models.Category { return models.CategoryReliability }
func (r *RuleREL008) Severity() models.Severity { return models.SeverityMedium }
func (r *RuleREL008) Description() string {
	return "Webhook node has responseMode set to responseNode without downstream Respond to Webhook node."
}
func (r *RuleREL008) Remediation() string {
	return "Add a 'Respond to Webhook' node to guarantee clients receive an HTTP response."
}
func (r *RuleREL008) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if !strings.Contains(strings.ToLower(node.Type), "webhook") {
			continue
		}

		responseMode, _ := node.Parameters["responseMode"].(string)
		if responseMode != "responseNode" {
			continue
		}

		hasResponder := false
		for _, other := range ctx.Workflow.Nodes {
			if strings.Contains(strings.ToLower(other.Type), "respondtowebhook") {
				hasResponder = true
				break
			}
		}

		if !hasResponder {
			line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Webhook node '%s' expects a response node, but none exists in the workflow", node.Name),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "parameters.responseMode",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// REL-009: missing-database-timeout
// -------------------------------------------------------------
type RuleREL009 struct{}

func (r *RuleREL009) ID() string             { return "REL-009" }
func (r *RuleREL009) Name() string           { return "missing-database-timeout" }
func (r *RuleREL009) Category() models.Category { return models.CategoryReliability }
func (r *RuleREL009) Severity() models.Severity { return models.SeverityLow }
func (r *RuleREL009) Description() string {
	return "Database query node lacks explicit statement timeout configuration."
}
func (r *RuleREL009) Remediation() string {
	return "Configure query statement timeout options to prevent unindexed queries from locking connections."
}
func (r *RuleREL009) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	dbKeywords := []string{"postgres", "mysql", "snowflake", "oracle", "mssql"}

	for _, node := range ctx.Workflow.Nodes {
		lowerType := strings.ToLower(node.Type)
		isDB := false
		for _, kw := range dbKeywords {
			if strings.Contains(lowerType, kw) {
				isDB = true
				break
			}
		}
		if !isDB {
			continue
		}

		options, ok := node.Parameters["options"].(map[string]interface{})
		hasTimeout := false
		if ok {
			if _, exists := options["timeout"]; exists {
				hasTimeout = true
			}
		}

		if !hasTimeout {
			line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Database node '%s' lacks explicit query timeout options", node.Name),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "parameters.options",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// REL-010: missing-conditional-fallback
// -------------------------------------------------------------
type RuleREL010 struct{}

func (r *RuleREL010) ID() string             { return "REL-010" }
func (r *RuleREL010) Name() string           { return "missing-conditional-fallback" }
func (r *RuleREL010) Category() models.Category { return models.CategoryReliability }
func (r *RuleREL010) Severity() models.Severity { return models.SeverityMedium }
func (r *RuleREL010) Description() string {
	return "If or Switch node leaves one or more branch outputs completely unconnected without fallback handling."
}
func (r *RuleREL010) Remediation() string {
	return "Connect unhandled branches to a logging node or error alert to avoid silent drop of non-matching items."
}
func (r *RuleREL010) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		lowerType := strings.ToLower(node.Type)
		if !strings.Contains(lowerType, ".if") && !strings.Contains(lowerType, ".switch") {
			continue
		}
		if node.Disabled && (ctx.Config == nil || !ctx.Config.IncludeDisabled) {
			continue
		}

		connections, exists := ctx.Workflow.Connections[node.Name]
		if !exists {
			continue
		}

		mainOutputs, hasMain := connections["main"]
		if !hasMain || len(mainOutputs) < 2 {
			line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Conditional node '%s' has unconnected branch outputs", node.Name),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "connections",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}
