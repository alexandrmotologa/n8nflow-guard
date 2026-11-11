# BEST-006: pinned-test-data-leak

## Metadata
- Rule ID: BEST-006
- Category: BestPractice
- Default Severity: MEDIUM

## Description
Workflow includes committed pinData test payloads, which inflate repository size and can leak customer records.

## Why It Matters
Workflows in production handle customer data and sensitive integration points.
Violating this rule can result in unintended exposure, operational downtime, or runaway resource consumption.

## Remediation
Clear pinned test data in n8n or use 'n8nflow-guard fix --strip-pindata' before committing to git.

## Examples

### Non-compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Non-compliant configuration exhibiting pinned-test-data-leak"
}
```

### Compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Remediated configuration conforming to security standards"
}
```
