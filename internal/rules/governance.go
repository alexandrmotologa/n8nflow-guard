package rules

import (
	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// -------------------------------------------------------------
// GOV-001: missing-workflow-tags
// -------------------------------------------------------------
type RuleGOV001 struct{}

func (r *RuleGOV001) ID() string             { return "GOV-001" }
func (r *RuleGOV001) Name() string           { return "missing-workflow-tags" }
func (r *RuleGOV001) Category() models.Category { return models.CategoryGovernance }
func (r *RuleGOV001) Severity() models.Severity { return models.SeverityLow }
func (r *RuleGOV001) Description() string {
	return "Workflow is missing organizational taxonomy tags in version control."
}
func (r *RuleGOV001) Remediation() string {
	return "Assign tags (e.g. environment, department, system) in n8n before exporting."
}
func (r *RuleGOV001) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	if len(ctx.Workflow.Tags) == 0 {
		issues = append(issues, models.Issue{
			RuleID:      r.ID(),
			RuleName:    r.Name(),
			Category:    r.Category(),
			Severity:    r.Severity(),
			Message:     "Workflow does not specify any categorization tags",
			FilePath:    ctx.Workflow.FilePath,
			Line:        1,
			Field:       "tags",
			Remediation: r.Remediation(),
		})
	}

	return issues
}

// -------------------------------------------------------------
// GOV-002: missing-trigger-node
// -------------------------------------------------------------
type RuleGOV002 struct{}

func (r *RuleGOV002) ID() string             { return "GOV-002" }
func (r *RuleGOV002) Name() string           { return "missing-trigger-node" }
func (r *RuleGOV002) Category() models.Category { return models.CategoryGovernance }
func (r *RuleGOV002) Severity() models.Severity { return models.SeverityMedium }
func (r *RuleGOV002) Description() string {
	return "Workflow has no entrypoint trigger node (webhook, schedule, poll, manual trigger)."
}
func (r *RuleGOV002) Remediation() string {
	return "Add a trigger node to ensure the workflow can receive invocation events."
}
func (r *RuleGOV002) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	triggers := ctx.Graph.FindTriggerNodes()
	// If workflow has nodes but no trigger
	if len(ctx.Workflow.Nodes) > 0 && len(triggers) == 0 {
		issues = append(issues, models.Issue{
			RuleID:      r.ID(),
			RuleName:    r.Name(),
			Category:    r.Category(),
			Severity:    r.Severity(),
			Message:     "Workflow has no trigger node to initiate execution",
			FilePath:    ctx.Workflow.FilePath,
			Line:        1,
			Field:       "nodes",
			Remediation: r.Remediation(),
		})
	}

	return issues
}
