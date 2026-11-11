package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// -------------------------------------------------------------
// BEST-001: orphan-node
// -------------------------------------------------------------
type RuleBEST001 struct{}

func (r *RuleBEST001) ID() string             { return "BEST-001" }
func (r *RuleBEST001) Name() string           { return "orphan-node" }
func (r *RuleBEST001) Category() models.Category { return models.CategoryBestPractice }
func (r *RuleBEST001) Severity() models.Severity { return models.SeverityLow }
func (r *RuleBEST001) Description() string {
	return "Active node has zero incoming and zero outgoing connections, leaving dead logic on canvas."
}
func (r *RuleBEST001) Remediation() string {
	return "Connect the orphan node into the workflow or remove it from the canvas."
}
func (r *RuleBEST001) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	orphans := ctx.Graph.FindOrphanNodes()
	for _, node := range orphans {
		line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
		issues = append(issues, models.Issue{
			RuleID:      r.ID(),
			RuleName:    r.Name(),
			Category:    r.Category(),
			Severity:    r.Severity(),
			Message:     fmt.Sprintf("Node '%s' is disconnected from all other nodes", node.Name),
			FilePath:    ctx.Workflow.FilePath,
			Line:        line,
			NodeID:      node.ID,
			NodeName:    node.Name,
			Field:       "connections",
			Remediation: r.Remediation(),
		})
	}

	return issues
}

// -------------------------------------------------------------
// BEST-002: default-node-name
// -------------------------------------------------------------
type RuleBEST002 struct{}

func (r *RuleBEST002) ID() string             { return "BEST-002" }
func (r *RuleBEST002) Name() string           { return "default-node-name" }
func (r *RuleBEST002) Category() models.Category { return models.CategoryBestPractice }
func (r *RuleBEST002) Severity() models.Severity { return models.SeverityLow }
func (r *RuleBEST002) Description() string {
	return "Node retains a generic default name (e.g. 'HTTP Request', 'Code 1'), reducing readability."
}
func (r *RuleBEST002) Remediation() string {
	return "Rename the node to reflect its specific domain role (e.g. 'Fetch Customer Profile', 'Format Invoice Items')."
}
func (r *RuleBEST002) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	defaultNamesRegex := regexp.MustCompile(`^(HTTP Request|Code|Webhook|Set|If|Switch|Schedule Trigger|Manual Trigger|Split In Batches)(\s+\d+)?$`)

	for _, node := range ctx.Workflow.Nodes {
		if defaultNamesRegex.MatchString(node.Name) {
			line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Node '%s' uses default generic name", node.Name),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "name",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// BEST-003: missing-node-description
// -------------------------------------------------------------
type RuleBEST003 struct{}

func (r *RuleBEST003) ID() string             { return "BEST-003" }
func (r *RuleBEST003) Name() string           { return "missing-node-description" }
func (r *RuleBEST003) Category() models.Category { return models.CategoryBestPractice }
func (r *RuleBEST003) Severity() models.Severity { return models.SeverityLow }
func (r *RuleBEST003) Description() string {
	return "Complex custom Code node lacks inline documentation or notes."
}
func (r *RuleBEST003) Remediation() string {
	return "Add descriptive notes explaining the transformations or business logic in the Code node."
}
func (r *RuleBEST003) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if !strings.Contains(strings.ToLower(node.Type), "code") {
			continue
		}

		if strings.TrimSpace(node.Notes) == "" {
			line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Code node '%s' contains no explanatory notes", node.Name),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "notes",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// BEST-004: excessive-canvas-complexity
// -------------------------------------------------------------
type RuleBEST004 struct{}

func (r *RuleBEST004) ID() string             { return "BEST-004" }
func (r *RuleBEST004) Name() string           { return "excessive-canvas-complexity" }
func (r *RuleBEST004) Category() models.Category { return models.CategoryBestPractice }
func (r *RuleBEST004) Severity() models.Severity { return models.SeverityLow }
func (r *RuleBEST004) Description() string {
	return "Workflow canvas exceeds 50 nodes without modular sub-workflow decomposition."
}
func (r *RuleBEST004) Remediation() string {
	return "Split monolithic workflows into smaller modular sub-workflows using Execute Workflow nodes."
}
func (r *RuleBEST004) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	threshold := 50
	if len(ctx.Workflow.Nodes) > threshold {
		issues = append(issues, models.Issue{
			RuleID:      r.ID(),
			RuleName:    r.Name(),
			Category:    r.Category(),
			Severity:    r.Severity(),
			Message:     fmt.Sprintf("Workflow contains %d nodes (recommended threshold is %d)", len(ctx.Workflow.Nodes), threshold),
			FilePath:    ctx.Workflow.FilePath,
			Line:        1,
			Field:       "nodes",
			Remediation: r.Remediation(),
		})
	}

	return issues
}

// -------------------------------------------------------------
// BEST-005: outdated-node-version
// -------------------------------------------------------------
type RuleBEST005 struct{}

func (r *RuleBEST005) ID() string             { return "BEST-005" }
func (r *RuleBEST005) Name() string           { return "outdated-node-version" }
func (r *RuleBEST005) Category() models.Category { return models.CategoryBestPractice }
func (r *RuleBEST005) Severity() models.Severity { return models.SeverityMedium }
func (r *RuleBEST005) Description() string {
	return "Node uses a legacy typeVersion (e.g. httpRequest v1) that lacks security and performance updates."
}
func (r *RuleBEST005) Remediation() string {
	return "Upgrade node to the latest version in the n8n UI canvas."
}
func (r *RuleBEST005) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if strings.Contains(strings.ToLower(node.Type), "httprequest") && node.TypeVersion < 4.0 {
			line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Node '%s' is on legacy version %.1f (current is >= 4.0)", node.Name, node.TypeVersion),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "typeVersion",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// BEST-006: pinned-test-data-leak
// -------------------------------------------------------------
type RuleBEST006 struct{}

func (r *RuleBEST006) ID() string             { return "BEST-006" }
func (r *RuleBEST006) Name() string           { return "pinned-test-data-leak" }
func (r *RuleBEST006) Category() models.Category { return models.CategoryBestPractice }
func (r *RuleBEST006) Severity() models.Severity { return models.SeverityMedium }
func (r *RuleBEST006) Description() string {
	return "Workflow includes committed pinData test payloads, which inflate repository size and can leak customer records."
}
func (r *RuleBEST006) Remediation() string {
	return "Clear pinned test data in n8n or use 'n8nflow-guard fix --strip-pindata' before committing to git."
}
func (r *RuleBEST006) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	if len(ctx.Workflow.PinData) > 0 {
		line := ctx.LineMapper.FindFieldLine("", "pinData")
		issues = append(issues, models.Issue{
			RuleID:       r.ID(),
			RuleName:     r.Name(),
			Category:     r.Category(),
			Severity:     r.Severity(),
			Message:      fmt.Sprintf("Workflow includes committed pinData for %d nodes", len(ctx.Workflow.PinData)),
			FilePath:     ctx.Workflow.FilePath,
			Line:         line,
			Field:        "pinData",
			Remediation:  r.Remediation(),
			Fixable:      true,
			FixSuggested: "Empty pinData payload",
		})
	}

	return issues
}
