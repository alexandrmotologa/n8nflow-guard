# BEST-002: default-node-name

## Metadata
- Rule ID: BEST-002
- Category: BestPractice
- Default Severity: LOW

## Description
Node retains a generic default name (e.g. 'HTTP Request', 'Code 1'), reducing readability.

## Why It Matters
Workflows in production handle customer data and sensitive integration points.
Violating this rule can result in unintended exposure, operational downtime, or runaway resource consumption.

## Remediation
Rename the node to reflect its specific domain role (e.g. 'Fetch Customer Profile', 'Format Invoice Items').

## Examples

### Non-compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Non-compliant configuration exhibiting default-node-name"
}
```

### Compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Remediated configuration conforming to security standards"
}
```
