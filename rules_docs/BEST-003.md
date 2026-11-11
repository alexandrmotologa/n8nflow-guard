# BEST-003: missing-node-description

## Metadata
- Rule ID: BEST-003
- Category: BestPractice
- Default Severity: LOW

## Description
Complex custom Code node lacks inline documentation or notes.

## Why It Matters
Workflows in production handle customer data and sensitive integration points.
Violating this rule can result in unintended exposure, operational downtime, or runaway resource consumption.

## Remediation
Add descriptive notes explaining the transformations or business logic in the Code node.

## Examples

### Non-compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Non-compliant configuration exhibiting missing-node-description"
}
```

### Compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Remediated configuration conforming to security standards"
}
```
