package rules

import (
	"sort"
	"strings"
	"sync"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// Registry manages the set of available static analysis rules.
type Registry struct {
	mu    sync.RWMutex
	rules map[string]Rule
}

// NewRegistry instantiates an empty rule registry.
func NewRegistry() *Registry {
	return &Registry{
		rules: make(map[string]Rule),
	}
}

// Register adds a rule to the registry.
func (r *Registry) Register(rule Rule) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[rule.ID()] = rule
}

// Get looks up a rule by its unique identifier (e.g., "SEC-001").
func (r *Registry) Get(id string) (Rule, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rule, exists := r.rules[strings.ToUpper(id)]
	return rule, exists
}

// All returns a sorted slice of all registered rules.
func (r *Registry) All() []Rule {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Rule, 0, len(r.rules))
	for _, rule := range r.rules {
		result = append(result, rule)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID() < result[j].ID()
	})

	return result
}

// Evaluate runs all active rules against the supplied workflow context.
func (r *Registry) Evaluate(ctx *RuleContext) []models.Issue {
	allRules := r.All()
	var issues []models.Issue

	for _, rule := range allRules {
		if ctx.Config != nil && !ctx.Config.IsRuleEnabled(rule.ID()) {
			continue
		}

		ruleIssues := rule.Evaluate(ctx)

		// Apply severity overrides and minimum severity filters
		for i := range ruleIssues {
			iss := ruleIssues[i]

			// Check for inline node note suppression
			if iss.NodeName != "" || iss.NodeID != "" {
				var targetNode *models.Node
				if iss.NodeName != "" {
					targetNode = ctx.Workflow.FindNodeByName(iss.NodeName)
				}
				if targetNode == nil && iss.NodeID != "" {
					targetNode = ctx.Workflow.FindNodeByID(iss.NodeID)
				}
				if targetNode != nil && targetNode.Notes != "" {
					ignoredRules := ParseNodeIgnoredRules(targetNode.Notes)
					if ignoredRules[rule.ID()] {
						continue
					}
				}
			}

			if ctx.Config != nil {
				iss.Severity = ctx.Config.GetRuleSeverity(rule.ID(), iss.Severity)
			}
			if ctx.Config != nil && ctx.Config.MinSeverity != "" {
				if iss.Severity.Rank() < ctx.Config.MinSeverity.Rank() {
					continue
				}
			}
			issues = append(issues, iss)
		}
	}

	return issues
}

var (
	defaultRegistry     *Registry
	defaultRegistryOnce sync.Once
)

// DefaultRegistry returns the singleton registry with all built-in rules registered.
func DefaultRegistry() *Registry {
	defaultRegistryOnce.Do(func() {
		defaultRegistry = NewRegistry()
		RegisterAllRules(defaultRegistry)
	})
	return defaultRegistry
}

// RegisterAllRules adds all 34 built-in rules to the given registry.
func RegisterAllRules(r *Registry) {
	// Security Suite (SEC-001 to SEC-010)
	r.Register(&RuleSEC001{})
	r.Register(&RuleSEC002{})
	r.Register(&RuleSEC003{})
	r.Register(&RuleSEC004{})
	r.Register(&RuleSEC005{})
	r.Register(&RuleSEC006{})
	r.Register(&RuleSEC007{})
	r.Register(&RuleSEC008{})
	r.Register(&RuleSEC009{})
	r.Register(&RuleSEC010{})

	// Reliability Suite (REL-001 to REL-011)
	r.Register(&RuleREL001{})
	r.Register(&RuleREL002{})
	r.Register(&RuleREL003{})
	r.Register(&RuleREL004{})
	r.Register(&RuleREL005{})
	r.Register(&RuleREL006{})
	r.Register(&RuleREL007{})
	r.Register(&RuleREL008{})
	r.Register(&RuleREL009{})
	r.Register(&RuleREL010{})
	r.Register(&RuleREL011{})

	// Performance Suite (PERF-001 to PERF-005)
	r.Register(&RulePERF001{})
	r.Register(&RulePERF002{})
	r.Register(&RulePERF003{})
	r.Register(&RulePERF004{})
	r.Register(&RulePERF005{})

	// Best Practice Suite (BEST-001 to BEST-006)
	r.Register(&RuleBEST001{})
	r.Register(&RuleBEST002{})
	r.Register(&RuleBEST003{})
	r.Register(&RuleBEST004{})
	r.Register(&RuleBEST005{})
	r.Register(&RuleBEST006{})

	// Governance Suite (GOV-001 to GOV-002)
	r.Register(&RuleGOV001{})
	r.Register(&RuleGOV002{})
}
