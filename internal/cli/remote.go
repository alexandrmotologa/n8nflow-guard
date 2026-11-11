package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
	"github.com/alexandrmotologa/n8nflow-guard/internal/parser"
	"github.com/alexandrmotologa/n8nflow-guard/internal/remote"
	"github.com/alexandrmotologa/n8nflow-guard/internal/reporter"
	"github.com/alexandrmotologa/n8nflow-guard/internal/rules"
	"github.com/spf13/cobra"
)

type remoteAuditFlags struct {
	url         string
	apiKey      string
	activeOnly  bool
	format      string
	minSeverity string
	failOn      string
	outputPath  string
}

func newRemoteCmd() *cobra.Command {
	flags := &remoteAuditFlags{}

	cmd := &cobra.Command{
		Use:   "remote",
		Short: "Audit live n8n instances directly via n8n REST API",
	}

	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Fetch and audit active workflows from a remote n8n server",
		Long: `Connects to an n8n instance using its REST API, downloads all active workflow definitions,
and runs static analysis against all 34 security, reliability, and governance rules.`,
		Example: `  # Audit production n8n server with terminal output
  n8nflow-guard remote audit --url https://n8n.example.com --api-key $N8N_API_KEY

  # Generate interactive HTML dashboard of remote instance
  n8nflow-guard remote audit --url https://n8n.example.com --api-key $N8N_API_KEY --format html --output audit.html`,
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL := flags.url
			if baseURL == "" {
				baseURL = os.Getenv("N8N_BASE_URL")
			}
			if baseURL == "" {
				return fmt.Errorf("n8n instance URL required via --url or N8N_BASE_URL environment variable")
			}

			apiKey := flags.apiKey
			if apiKey == "" {
				apiKey = os.Getenv("N8N_API_KEY")
			}
			if apiKey == "" {
				return fmt.Errorf("n8n API key required via --api-key or N8N_API_KEY environment variable")
			}

			var minSev models.Severity
			if flags.minSeverity != "" {
				s, err := models.ParseSeverity(flags.minSeverity)
				if err != nil {
					return err
				}
				minSev = s
				appConfig.MinSeverity = minSev
			}

			var failOn models.Severity
			if flags.failOn != "" {
				s, err := models.ParseSeverity(flags.failOn)
				if err != nil {
					return err
				}
				failOn = s
			}

			if !quiet {
				fmt.Fprintf(cmd.OutOrStdout(), "Connecting to n8n instance at %s...\n", baseURL)
			}

			client := remote.NewClient(baseURL, apiKey)
			workflows, err := client.FetchWorkflows(flags.activeOnly)
			if err != nil {
				return fmt.Errorf("remote audit failed: %w", err)
			}

			if !quiet {
				fmt.Fprintf(cmd.OutOrStdout(), "Retrieved %d workflow(s). Evaluating rules...\n", len(workflows))
			}

			reg := rules.DefaultRegistry()
			var allIssues []models.Issue

			for _, wf := range workflows {
				graph := parser.BuildGraph(wf)
				lineMapper := parser.NewLineMapper(wf.RawBytes)
				ctx := &rules.RuleContext{
					Workflow:   wf,
					Graph:      graph,
					LineMapper: lineMapper,
					Config:     appConfig,
				}
				allIssues = append(allIssues, reg.Evaluate(ctx)...)
			}

			summary := reporter.ComputeSummary(len(workflows), len(workflows), allIssues)

			rep, err := reporter.NewReporter(flags.format, noColor)
			if err != nil {
				return err
			}

			var out io.Writer = cmd.OutOrStdout()
			var outFile *os.File
			if flags.outputPath != "" {
				f, err := os.Create(flags.outputPath)
				if err != nil {
					return fmt.Errorf("failed creating output file %s: %w", flags.outputPath, err)
				}
				defer f.Close()
				out = f
				outFile = f
			}

			if err := rep.Report(out, allIssues, summary); err != nil {
				return err
			}

			if outFile != nil && !quiet {
				fmt.Fprintf(cmd.OutOrStdout(), "Wrote remote audit report to %s\n", flags.outputPath)
			}

			failRank := failOn.Rank()
			if failRank == 0 {
				failRank = models.SeverityHigh.Rank()
			}
			for _, iss := range allIssues {
				if iss.Severity.Rank() >= failRank {
					os.Exit(1)
				}
			}

			return nil
		},
	}

	auditCmd.Flags().StringVar(&flags.url, "url", "", "base URL of n8n instance (e.g. https://n8n.example.com)")
	auditCmd.Flags().StringVar(&flags.apiKey, "api-key", "", "n8n user or instance API key")
	auditCmd.Flags().BoolVar(&flags.activeOnly, "active-only", true, "audit only active production workflows")
	auditCmd.Flags().StringVarP(&flags.format, "format", "f", "terminal", "output format (terminal, html, sarif, github-actions, json, markdown)")
	auditCmd.Flags().StringVarP(&flags.minSeverity, "min-severity", "s", "LOW", "minimum severity level to report")
	auditCmd.Flags().StringVar(&flags.failOn, "fail-on", "HIGH", "threshold severity to exit with code 1")
	auditCmd.Flags().StringVarP(&flags.outputPath, "output", "o", "", "file path to write the generated report")

	cmd.AddCommand(auditCmd)
	return cmd
}
