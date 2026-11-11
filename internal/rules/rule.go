package rules

import (
	"github.com/alexandrmotologa/n8nflow-guard/internal/config"
	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
	"github.com/alexandrmotologa/n8nflow-guard/internal/parser"
)

// RuleContext supplies all workflow metadata, graphs, line mappers, and configuration to rule evaluators.
type RuleContext struct {
	Workflow   *models.Workflow
	Graph      *parser.Graph
	LineMapper *parser.LineMapper
	Config     *config.Config
}

// Rule defines the contract for an n8nflow-guard static analysis rule.
type Rule interface {
	ID() string
	Name() string
	Category() models.Category
	Severity() models.Severity
	Description() string
	Remediation() string
	Evaluate(ctx *RuleContext) []models.Issue
}
