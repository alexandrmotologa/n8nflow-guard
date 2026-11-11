package cli

import (
	"fmt"
	"os"

	"github.com/alexandrmotologa/n8nflow-guard/internal/config"
	"github.com/spf13/cobra"
)

var (
	configPath string
	quiet      bool
	noColor    bool
	appConfig  *config.Config

	// Version metadata set during compilation via -ldflags
	Version   = "1.0.0"
	GitCommit = "dev"
	BuildDate = "unknown"
)

// RootCmd represents the base command executed without subcommands.
var RootCmd = &cobra.Command{
	Use:   "n8nflow-guard",
	Short: "Static analysis and security linter for n8n workflow JSON files",
	Long: `n8nflow-guard is a fast static analysis tool and security validator for n8n workflows.
It parses n8n JSON documents, evaluates DAG graph topology, identifies hardcoded secrets,
uncovers missing error workflows, and emits SARIF reports for GitHub Code Scanning.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if configPath != "" {
			appConfig, err = config.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("error loading config from %s: %w", configPath, err)
			}
		} else {
			appConfig, err = config.FindAndLoadConfig(".")
			if err != nil {
				return fmt.Errorf("error finding config: %w", err)
			}
		}
		return nil
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to configuration file (.n8nflow-guard.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-error output")
	RootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable ANSI color formatting")

	RootCmd.AddCommand(newLintCmd())
	RootCmd.AddCommand(newFixCmd())
	RootCmd.AddCommand(newRulesCmd())
	RootCmd.AddCommand(newInitCmd())
	RootCmd.AddCommand(newVersionCmd())
}

// Execute handles CLI entry and error routing.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
