package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/engine"
	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
	"github.com/alexandrmotologa/n8nflow-guard/internal/reporter"
	"github.com/spf13/cobra"
)

type lintFlags struct {
	format           string
	minSeverity      string
	failOn           string
	outputPath       string
	strict           bool
	ignoreRules      string
	includeDisabled  bool
	staged           bool
	since            string
	baseline         string
	generateBaseline string
}

func newLintCmd() *cobra.Command {
	flags := &lintFlags{}

	cmd := &cobra.Command{
		Use:   "lint [paths...]",
		Short: "Analyze n8n workflow files for security and reliability violations",
		Long: `Scan one or more workflow JSON files, directories, or glob patterns against
all active security, reliability, performance, best practice, and governance rules.`,
		Example: `  # Scan a workflows directory with terminal output
  n8nflow-guard lint ./workflows/

  # Emit SARIF report for GitHub Code Scanning
  n8nflow-guard lint --format sarif --output results.sarif ./workflows/

  # Fail build on CRITICAL or HIGH severity issues
  n8nflow-guard lint --fail-on HIGH ./workflows/

  # Output GitHub Actions workflow commands
  n8nflow-guard lint --format github-actions ./workflows/`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				args = []string{"."}
			}

			// Parse severity filters
			var minSev models.Severity
			if flags.minSeverity != "" {
				s, err := models.ParseSeverity(flags.minSeverity)
				if err != nil {
					return err
				}
				minSev = s
			}

			var failOn models.Severity
			if flags.failOn != "" {
				s, err := models.ParseSeverity(flags.failOn)
				if err != nil {
					return err
				}
				failOn = s
			}

			if flags.ignoreRules != "" {
				for _, r := range strings.Split(flags.ignoreRules, ",") {
					trimmed := strings.TrimSpace(r)
					if trimmed != "" {
						appConfig.IgnoreRules = append(appConfig.IgnoreRules, trimmed)
					}
				}
			}

			if flags.includeDisabled {
				appConfig.IncludeDisabled = true
			}

			scanOpts := engine.ScanOptions{
				Paths:                args,
				MinSeverity:          minSev,
				FailOn:               failOn,
				Strict:               flags.strict,
				Staged:               flags.staged,
				Since:                flags.since,
				BaselinePath:         flags.baseline,
				GenerateBaselinePath: flags.generateBaseline,
				Config:               appConfig,
			}

			res, err := engine.RunScan(scanOpts)
			if err != nil {
				return err
			}

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

			if err := rep.Report(out, res.Issues, res.Summary); err != nil {
				return fmt.Errorf("error generating report: %w", err)
			}

			if outFile != nil && !quiet {
				fmt.Fprintf(cmd.OutOrStdout(), "Wrote report to %s\n", flags.outputPath)
			}

			if res.HasFailures {
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.format, "format", "f", "terminal", "output format (terminal, sarif, github-actions, json, markdown)")
	cmd.Flags().StringVarP(&flags.minSeverity, "min-severity", "s", "LOW", "minimum severity level to report (CRITICAL, HIGH, MEDIUM, LOW, INFO)")
	cmd.Flags().StringVar(&flags.failOn, "fail-on", "HIGH", "threshold severity to exit with code 1 (CRITICAL, HIGH, MEDIUM, LOW)")
	cmd.Flags().StringVarP(&flags.outputPath, "output", "o", "", "write report output to specified file path")
	cmd.Flags().BoolVar(&flags.strict, "strict", false, "treat all warnings as errors (fail on LOW or higher)")
	cmd.Flags().StringVar(&flags.ignoreRules, "ignore-rule", "", "comma-separated list of rule IDs to skip (e.g. BEST-002,REL-004)")
	cmd.Flags().BoolVar(&flags.includeDisabled, "include-disabled", false, "also evaluate disabled nodes")
	cmd.Flags().BoolVar(&flags.staged, "staged", false, "lint only workflow files staged in git index")
	cmd.Flags().StringVar(&flags.since, "since", "", "lint workflow files modified relative to a git ref (e.g. origin/main)")
	cmd.Flags().StringVar(&flags.baseline, "baseline", "", "path to baseline file to suppress previously recorded issues")
	cmd.Flags().StringVar(&flags.generateBaseline, "generate-baseline", "", "generate a baseline file recording current issues")

	return cmd
}
