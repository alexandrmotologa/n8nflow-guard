# PERF-004: n-plus-one-request-pattern

## Metadata
- Rule ID: PERF-004
- Category: Performance
- Default Severity: MEDIUM

## Description
HTTP Request node executes in a per-item loop without batching or request aggregation.

## Why It Matters
Workflows in production handle customer data and sensitive integration points.
Violating this rule can result in unintended exposure, operational downtime, or runaway resource consumption.

## Remediation
Use batch API endpoints or aggregate items into bulk requests where possible.

## Examples

### Non-compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Non-compliant configuration exhibiting n-plus-one-request-pattern"
}
```

### Compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Remediated configuration conforming to security standards"
}
```
