package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/graphviz"
	"github.com/alexandrmotologa/n8nflow-guard/internal/parser"
	"github.com/spf13/cobra"
)

type graphFlags struct {
	format     string
	outputPath string
}

func newGraphCmd() *cobra.Command {
	flags := &graphFlags{}

	cmd := &cobra.Command{
		Use:   "graph <workflow.json>",
		Short: "Generate visual graph diagram of an n8n workflow (Mermaid, DOT, ASCII)",
		Long: `Visualizes the Directed Acyclic Graph (DAG) of an n8n workflow file.
Supports output in Mermaid flowchart syntax (for Markdown/GitHub PRs),
Graphviz DOT format, and ASCII console tree.`,
		Example: `  # Generate Mermaid flowchart for GitHub Markdown
  n8nflow-guard graph ./workflows/order-sync.json --format mermaid

  # Output Graphviz DOT file
  n8nflow-guard graph ./workflows/order-sync.json --format dot --output flow.dot

  # Display ASCII tree in terminal
  n8nflow-guard graph ./workflows/order-sync.json --format ascii`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			workflows, err := parser.ParseWorkflowsFromFile(filePath)
			if err != nil {
				return err
			}

			if len(workflows) == 0 {
				return fmt.Errorf("no workflows found in %s", filePath)
			}

			wf := workflows[0]
			g := parser.BuildGraph(wf)

			var output string
			switch strings.ToLower(flags.format) {
			case "mermaid", "mmd", "":
				output = graphviz.ExportMermaid(g)
			case "dot", "graphviz":
				output = graphviz.ExportDOT(g)
			case "ascii", "tree", "text":
				output = graphviz.ExportASCII(g)
			default:
				return fmt.Errorf("unsupported graph format: %s (supported: mermaid, dot, ascii)", flags.format)
			}

			var w io.Writer = cmd.OutOrStdout()
			if flags.outputPath != "" {
				f, err := os.Create(flags.outputPath)
				if err != nil {
					return fmt.Errorf("failed creating output file %s: %w", flags.outputPath, err)
				}
				defer f.Close()
				w = f
			}

			fmt.Fprint(w, output)
			return nil
		},
	}

	cmd.Flags().StringVarP(&flags.format, "format", "f", "mermaid", "diagram format (mermaid, dot, ascii)")
	cmd.Flags().StringVarP(&flags.outputPath, "output", "o", "", "file path to write the generated diagram")

	return cmd
}
