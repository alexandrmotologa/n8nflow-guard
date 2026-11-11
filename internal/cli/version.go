package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print n8nflow-guard version and runtime information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "n8nflow-guard version %s (commit: %s, built: %s, %s/%s)\n",
				Version, GitCommit, BuildDate, runtime.GOOS, runtime.GOARCH)
		},
	}
}
