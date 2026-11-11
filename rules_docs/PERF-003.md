# PERF-003: excessive-execution-logging

## Metadata
- Rule ID: PERF-003
- Category: Performance
- Default Severity: MEDIUM

## Description
Workflow persists execution progress step by step, degrading throughput on production n8n instances.

## Why It Matters
Workflows in production handle customer data and sensitive integration points.
Violating this rule can result in unintended exposure, operational downtime, or runaway resource consumption.

## Remediation
Set saveExecutionProgress to false unless actively troubleshooting.

## Examples

### Non-compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Non-compliant configuration exhibiting excessive-execution-logging"
}
```

### Compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Remediated configuration conforming to security standards"
}
```
