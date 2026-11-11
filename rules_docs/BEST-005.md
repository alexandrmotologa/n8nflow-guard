# BEST-005: outdated-node-version

## Metadata
- Rule ID: BEST-005
- Category: BestPractice
- Default Severity: MEDIUM

## Description
Node uses a legacy typeVersion (e.g. httpRequest v1) that lacks security and performance updates.

## Why It Matters
Workflows in production handle customer data and sensitive integration points.
Violating this rule can result in unintended exposure, operational downtime, or runaway resource consumption.

## Remediation
Upgrade node to the latest version in the n8n UI canvas.

## Examples

### Non-compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Non-compliant configuration exhibiting outdated-node-version"
}
```

### Compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Remediated configuration conforming to security standards"
}
```
