# BEST-004: excessive-canvas-complexity

## Metadata
- Rule ID: BEST-004
- Category: BestPractice
- Default Severity: LOW

## Description
Workflow canvas exceeds 50 nodes without modular sub-workflow decomposition.

## Why It Matters
Workflows in production handle customer data and sensitive integration points.
Violating this rule can result in unintended exposure, operational downtime, or runaway resource consumption.

## Remediation
Split monolithic workflows into smaller modular sub-workflows using Execute Workflow nodes.

## Examples

### Non-compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Non-compliant configuration exhibiting excessive-canvas-complexity"
}
```

### Compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Remediated configuration conforming to security standards"
}
```
