package parser

import (
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
)

// Graph represents the directed node connection structure of a workflow.
type Graph struct {
	Workflow   *models.Workflow
	Adjacency  map[string][]string // Outgoing edges: source -> targets
	Inbound    map[string][]string // Incoming edges: target -> sources
	NodeByName map[string]*models.Node
	NodeByID   map[string]*models.Node
}

// BuildGraph constructs an indexed directed graph from an n8n workflow.
func BuildGraph(w *models.Workflow) *Graph {
	g := &Graph{
		Workflow:   w,
		Adjacency:  make(map[string][]string),
		Inbound:    make(map[string][]string),
		NodeByName: make(map[string]*models.Node),
		NodeByID:   make(map[string]*models.Node),
	}

	for _, node := range w.Nodes {
		g.NodeByName[node.Name] = node
		g.NodeByID[node.ID] = node
		g.Adjacency[node.Name] = []string{}
		g.Inbound[node.Name] = []string{}
	}

	for sourceName, typeOutputs := range w.Connections {
		// Resolve source node canonical name
		canonicalSource := sourceName
		if node, exists := g.NodeByID[sourceName]; exists {
			canonicalSource = node.Name
		}

		for _, outputGroups := range typeOutputs {
			for _, targets := range outputGroups {
				for _, target := range targets {
					targetName := target.Node
					if targetNode, exists := g.NodeByID[target.Node]; exists {
						targetName = targetNode.Name
					}

					g.Adjacency[canonicalSource] = append(g.Adjacency[canonicalSource], targetName)
					g.Inbound[targetName] = append(g.Inbound[targetName], canonicalSource)
				}
			}
		}
	}

	return g
}

// FindOrphanNodes identifies active nodes with zero inbound and zero outbound connections.
// If the workflow contains only a single node, it is not marked as an orphan.
func (g *Graph) FindOrphanNodes() []*models.Node {
	if len(g.Workflow.Nodes) <= 1 {
		return nil
	}

	orphans := make([]*models.Node, 0)
	for _, node := range g.Workflow.Nodes {
		if node.Disabled {
			continue
		}
		inCount := len(g.Inbound[node.Name])
		outCount := len(g.Adjacency[node.Name])

		if inCount == 0 && outCount == 0 {
			orphans = append(orphans, node)
		}
	}
	return orphans
}

// DetectCycles finds directed cycles in the graph using three-color DFS traversal.
// Returns a slice of cycle paths, where each path is a slice of node names forming the cycle.
func (g *Graph) DetectCycles() [][]string {
	const (
		white = 0 // Unvisited
		gray  = 1 // Currently exploring
		black = 2 // Fully explored
	)

	colors := make(map[string]int)
	for _, node := range g.Workflow.Nodes {
		colors[node.Name] = white
	}

	var cycles [][]string
	var path []string

	var dfs func(current string)
	dfs = func(current string) {
		colors[current] = gray
		path = append(path, current)

		for _, neighbor := range g.Adjacency[current] {
			if _, exists := colors[neighbor]; !exists {
				continue
			}

			if colors[neighbor] == gray {
				// Found a back-edge indicating a cycle
				cycleStart := -1
				for i, name := range path {
					if name == neighbor {
						cycleStart = i
						break
					}
				}
				if cycleStart != -1 {
					cycle := make([]string, len(path)-cycleStart+1)
					copy(cycle, path[cycleStart:])
					cycle[len(cycle)-1] = neighbor
					cycles = append(cycles, cycle)
				}
			} else if colors[neighbor] == white {
				dfs(neighbor)
			}
		}

		path = path[:len(path)-1]
		colors[current] = black
	}

	for _, node := range g.Workflow.Nodes {
		if colors[node.Name] == white {
			dfs(node.Name)
		}
	}

	return cycles
}

// FindTriggerNodes returns nodes that initiate execution.
func (g *Graph) FindTriggerNodes() []*models.Node {
	triggers := make([]*models.Node, 0)
	for _, node := range g.Workflow.Nodes {
		if IsTriggerNodeType(node.Type) {
			triggers = append(triggers, node)
		}
	}
	return triggers
}

// IsTriggerNodeType checks whether a node type represents an entrypoint/trigger.
func IsTriggerNodeType(nodeType string) bool {
	lower := strings.ToLower(nodeType)
	triggers := []string{
		"trigger",
		"webhook",
		"poll",
		"interval",
		"cron",
		"manual",
		"onboarding",
		"form",
		"chat",
	}
	for _, term := range triggers {
		if strings.Contains(lower, term) {
			return true
		}
	}
	return false
}

// HasDownstreamErrorHandling checks if a node connects to an error branch or error handler.
func (g *Graph) HasDownstreamErrorHandling(nodeName string) bool {
	connections, exists := g.Workflow.Connections[nodeName]
	if !exists {
		return false
	}
	// Check for dedicated error output connection
	if errorOutputs, hasError := connections["error"]; hasError && len(errorOutputs) > 0 {
		return true
	}
	return false
}
