# BEST-001: orphan-node

## Metadata
- Rule ID: BEST-001
- Category: BestPractice
- Default Severity: LOW

## Description
Active node has zero incoming and zero outgoing connections, leaving dead logic on canvas.

## Why It Matters
Workflows in production handle customer data and sensitive integration points.
Violating this rule can result in unintended exposure, operational downtime, or runaway resource consumption.

## Remediation
Connect the orphan node into the workflow or remove it from the canvas.

## Examples

### Non-compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Non-compliant configuration exhibiting orphan-node"
}
```

### Compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Remediated configuration conforming to security standards"
}
```
