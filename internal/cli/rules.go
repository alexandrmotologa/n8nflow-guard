package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
	"github.com/alexandrmotologa/n8nflow-guard/internal/rules"
	"github.com/spf13/cobra"
)

type rulesFlags struct {
	ruleID   string
	category string
	format   string
}

func newRulesCmd() *cobra.Command {
	flags := &rulesFlags{}

	cmd := &cobra.Command{
		Use:   "rules",
		Short: "List and inspect available static analysis rules",
		Long:  `Display documentation, category groupings, severity ratings, and remediation instructions for all 33 built-in rules.`,
		Example: `  # List all rules in an interactive table
  n8nflow-guard rules

  # View detailed guidance for SEC-001
  n8nflow-guard rules --rule SEC-001

  # Filter by category
  n8nflow-guard rules --category Security`,
		RunE: func(cmd *cobra.Command, args []string) error {
			reg := rules.DefaultRegistry()

			// Specific rule inspection
			if flags.ruleID != "" {
				r, found := reg.Get(flags.ruleID)
				if !found {
					return fmt.Errorf("unknown rule ID: %s", flags.ruleID)
				}

				fmt.Fprintf(cmd.OutOrStdout(), "Rule ID:      %s\n", r.ID())
				fmt.Fprintf(cmd.OutOrStdout(), "Name:         %s\n", r.Name())
				fmt.Fprintf(cmd.OutOrStdout(), "Category:     %s\n", r.Category())
				fmt.Fprintf(cmd.OutOrStdout(), "Severity:     %s\n", r.Severity())
				fmt.Fprintf(cmd.OutOrStdout(), "Description:  %s\n", r.Description())
				fmt.Fprintf(cmd.OutOrStdout(), "Remediation:  %s\n", r.Remediation())
				fmt.Fprintf(cmd.OutOrStdout(), "Docs:         https://github.com/alexandrmotologa/n8nflow-guard/tree/main/rules_docs/%s.md\n", r.ID())
				return nil
			}

			all := reg.All()
			var filtered []rules.Rule
			for _, r := range all {
				if flags.category != "" && !strings.EqualFold(string(r.Category()), flags.category) {
					continue
				}
				filtered = append(filtered, r)
			}

			switch strings.ToLower(flags.format) {
			case "json":
				type ruleInfo struct {
					ID          string          `json:"id"`
					Name        string          `json:"name"`
					Category    models.Category `json:"category"`
					Severity    models.Severity `json:"severity"`
					Description string          `json:"description"`
					Remediation string          `json:"remediation"`
				}
				var infoList []ruleInfo
				for _, r := range filtered {
					infoList = append(infoList, ruleInfo{
						ID:          r.ID(),
						Name:        r.Name(),
						Category:    r.Category(),
						Severity:    r.Severity(),
						Description: r.Description(),
						Remediation: r.Remediation(),
					})
				}
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(infoList)

			case "markdown":
				fmt.Fprintf(cmd.OutOrStdout(), "| ID | Category | Severity | Name | Description |\n")
				fmt.Fprintf(cmd.OutOrStdout(), "|---|---|---|---|---|\n")
				for _, r := range filtered {
					fmt.Fprintf(cmd.OutOrStdout(), "| `%s` | %s | %s | %s | %s |\n",
						r.ID(), r.Category(), r.Severity(), r.Name(), r.Description())
				}

			default: // Table
				fmt.Fprintf(cmd.OutOrStdout(), "%-8s %-14s %-10s %-28s %s\n", "ID", "CATEGORY", "SEVERITY", "NAME", "DESCRIPTION")
				fmt.Fprintln(cmd.OutOrStdout(), strings.Repeat("-", 100))
				for _, r := range filtered {
					desc := r.Description()
					if len(desc) > 45 {
						desc = desc[:42] + "..."
					}
					fmt.Fprintf(cmd.OutOrStdout(), "%-8s %-14s %-10s %-28s %s\n",
						r.ID(), r.Category(), r.Severity(), r.Name(), desc)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "\nTotal rules: %d\n", len(filtered))
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.ruleID, "rule", "r", "", "inspect detailed documentation for a specific rule ID")
	cmd.Flags().StringVar(&flags.category, "category", "", "filter rules by category (Security, Reliability, Performance, BestPractice, Governance)")
	cmd.Flags().StringVar(&flags.format, "format", "table", "display format (table, json, markdown)")

	return cmd
}
