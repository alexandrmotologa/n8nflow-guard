package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexandrmotologa/n8nflow-guard/internal/rules"
)

func main() {
	outDir := "rules_docs"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create %s: %v\n", outDir, err)
		os.Exit(1)
	}

	reg := rules.DefaultRegistry()
	allRules := reg.All()

	for _, r := range allRules {
		filename := filepath.Join(outDir, fmt.Sprintf("%s.md", r.ID()))
		content := fmt.Sprintf(`# %s: %s

## Metadata
- Rule ID: %s
- Category: %s
- Default Severity: %s

## Description
%s

## Why It Matters
Workflows in production handle customer data and sensitive integration points.
Violating this rule can result in unintended exposure, operational downtime, or runaway resource consumption.

## Remediation
%s

## Examples

### Non-compliant Configuration
`+"```json"+`
{
  "name": "Sample Node",
  "notes": "Non-compliant configuration exhibiting %s"
}
`+"```"+`

### Compliant Configuration
`+"```json"+`
{
  "name": "Sample Node",
  "notes": "Remediated configuration conforming to security standards"
}
`+"```"+`
`, r.ID(), r.Name(), r.ID(), r.Category(), r.Severity(), r.Description(), r.Remediation(), r.Name())

		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", filename, err)
			os.Exit(1)
		}
		fmt.Println("Generated", filename)
	}
}
