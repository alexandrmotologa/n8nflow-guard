package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const starterConfigYAML = `# n8nflow-guard configuration file
version: "1.0"

# Minimum severity level to report in terminal and reports (CRITICAL, HIGH, MEDIUM, LOW, INFO)
minSeverity: LOW

# Severity threshold to fail CI builds with exit code 1 (CRITICAL, HIGH, MEDIUM, LOW)
failOn: HIGH

# Glob patterns to exclude from static analysis
ignorePaths:
  - ".git/**"
  - "node_modules/**"
  - "vendor/**"
  - "archive/**"

# Specific rules to ignore globally
ignoreRules:
  - "BEST-002" # Ignore default node names if team allows default labels

# Per-rule customizations
rules:
  SEC-001:
    severity: CRITICAL
    enabled: true
  REL-001:
    severity: HIGH
    enabled: true
  BEST-004:
    severity: LOW
    properties:
      maxNodes: 60

# Proprietary organizational API key patterns
customSecrets: []

# Set to true to also validate disabled canvas nodes
includeDisabled: false
`

const starterPreCommitYAML = `# .pre-commit-config.yaml hook definition for n8nflow-guard
repos:
  - repo: https://github.com/alexandrmotologa/n8nflow-guard
    rev: v1.0.0
    hooks:
      - id: n8nflow-guard-lint
        name: n8nflow-guard lint
        description: Static analysis and security linter for n8n workflow JSON files
        entry: n8nflow-guard lint
        language: golang
        files: \.json$
        types: [file]
`

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize starter configuration and pre-commit configuration",
		Long:  `Creates a default .n8nflow-guard.yaml configuration and .pre-commit-config.yaml in the working directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configTarget := ".n8nflow-guard.yaml"
			if _, err := os.Stat(configTarget); err == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "%s already exists, skipping.\n", configTarget)
			} else {
				if err := os.WriteFile(configTarget, []byte(starterConfigYAML), 0644); err != nil {
					return fmt.Errorf("failed creating %s: %w", configTarget, err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", configTarget)
			}

			preCommitTarget := ".pre-commit-config.yaml"
			if _, err := os.Stat(preCommitTarget); err == nil {
				fmt.Fprintf(cmd.OutOrStdout(), "%s already exists, skipping.\n", preCommitTarget)
			} else {
				if err := os.WriteFile(preCommitTarget, []byte(starterPreCommitYAML), 0644); err != nil {
					return fmt.Errorf("failed creating %s: %w", preCommitTarget, err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", preCommitTarget)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Initialization complete.")
			return nil
		},
	}
}
