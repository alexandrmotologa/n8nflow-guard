package graphviz

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/parser"
)

var idSanitizer = regexp.MustCompile(`[^a-zA-Z0-9_]`)

func sanitizeID(name string) string {
	s := idSanitizer.ReplaceAllString(name, "_")
	if len(s) == 0 {
		return "node"
	}
	return s
}

// ExportMermaid generates a Mermaid.js flowchart (flowchart TD).
func ExportMermaid(g *parser.Graph) string {
	var sb strings.Builder
	sb.WriteString("flowchart TD\n")

	// Collect and sort node names for deterministic output
	var nodeNames []string
	for name := range g.NodeByName {
		nodeNames = append(nodeNames, name)
	}
	sort.Strings(nodeNames)

	// Declare node definitions
	for _, name := range nodeNames {
		node := g.NodeByName[name]
		sanitized := sanitizeID(name)
		nodeLabel := strings.ReplaceAll(name, `"`, `\"`)

		if parser.IsTriggerNodeType(node.Type) {
			// Hexagon for triggers
			sb.WriteString(fmt.Sprintf("    %s{{\"%s\"}}\n", sanitized, nodeLabel))
		} else {
			// Rounded rectangle
			sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitized, nodeLabel))
		}
	}

	sb.WriteString("\n")

	// Declare connections
	for _, source := range nodeNames {
		targets := g.Adjacency[source]
		for _, target := range targets {
			srcID := sanitizeID(source)
			tgtID := sanitizeID(target)
			sb.WriteString(fmt.Sprintf("    %s --> %s\n", srcID, tgtID))
		}
	}

	// Add styles for triggers and disabled nodes
	sb.WriteString("\n    classDef trigger fill:#1b5e20,stroke:#81c784,color:#fff;\n")
	sb.WriteString("    classDef disabled fill:#424242,stroke:#757575,color:#bbb,stroke-dasharray: 5 5;\n")

	for _, name := range nodeNames {
		node := g.NodeByName[name]
		sanitized := sanitizeID(name)
		if node.Disabled {
			sb.WriteString(fmt.Sprintf("    class %s disabled;\n", sanitized))
		} else if parser.IsTriggerNodeType(node.Type) {
			sb.WriteString(fmt.Sprintf("    class %s trigger;\n", sanitized))
		}
	}

	return sb.String()
}

// ExportDOT generates a Graphviz DOT diagram.
func ExportDOT(g *parser.Graph) string {
	var sb strings.Builder
	sb.WriteString("digraph Workflow {\n")
	sb.WriteString("    rankdir=LR;\n")
	sb.WriteString("    node [shape=box, style=\"rounded,filled\", fontname=\"Helvetica\", fillcolor=\"#f5f5f5\"];\n\n")

	var nodeNames []string
	for name := range g.NodeByName {
		nodeNames = append(nodeNames, name)
	}
	sort.Strings(nodeNames)

	for _, name := range nodeNames {
		node := g.NodeByName[name]
		sanitized := sanitizeID(name)
		label := strings.ReplaceAll(name, `"`, `\"`)
		fillColor := "#e8f5e9"
		if !parser.IsTriggerNodeType(node.Type) {
			fillColor = "#ede7f6"
		}
		if node.Disabled {
			fillColor = "#eeeeee"
		}
		sb.WriteString(fmt.Sprintf("    %s [label=\"%s\", fillcolor=\"%s\"];\n", sanitized, label, fillColor))
	}

	sb.WriteString("\n")
	for _, source := range nodeNames {
		targets := g.Adjacency[source]
		for _, target := range targets {
			srcID := sanitizeID(source)
			tgtID := sanitizeID(target)
			sb.WriteString(fmt.Sprintf("    %s -> %s;\n", srcID, tgtID))
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

// ExportASCII generates an indented ASCII text summary of graph connectivity.
func ExportASCII(g *parser.Graph) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Workflow: %s (%d nodes)\n", g.Workflow.Name, len(g.Workflow.Nodes)))
	sb.WriteString(strings.Repeat("-", 60) + "\n")

	var nodeNames []string
	for name := range g.NodeByName {
		nodeNames = append(nodeNames, name)
	}
	sort.Strings(nodeNames)

	for _, name := range nodeNames {
		node := g.NodeByName[name]
		nodeType := "Action"
		if parser.IsTriggerNodeType(node.Type) {
			nodeType = "Trigger"
		}
		if node.Disabled {
			nodeType += " (Disabled)"
		}

		targets := g.Adjacency[name]
		inbound := g.Inbound[name]

		sb.WriteString(fmt.Sprintf("[%s] %s (%s)\n", nodeType, name, node.Type))
		if len(inbound) > 0 {
			sb.WriteString(fmt.Sprintf("    <- from: %s\n", strings.Join(inbound, ", ")))
		} else if !parser.IsTriggerNodeType(node.Type) && len(targets) == 0 {
			sb.WriteString("    ** ORPHAN NODE (0 in, 0 out) **\n")
		}

		if len(targets) > 0 {
			sb.WriteString(fmt.Sprintf("    -> to:   %s\n", strings.Join(targets, ", ")))
		} else {
			sb.WriteString("    -> [Terminal Output]\n")
		}
	}

	return sb.String()
}
