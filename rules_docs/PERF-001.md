# PERF-001: infinite-loop-cycle

## Metadata
- Rule ID: PERF-001
- Category: Performance
- Default Severity: HIGH

## Description
Workflow graph contains a directed loop/cycle that risks infinite execution.

## Why It Matters
Workflows in production handle customer data and sensitive integration points.
Violating this rule can result in unintended exposure, operational downtime, or runaway resource consumption.

## Remediation
Break the cycle or introduce a Split In Batches node with an explicit loop counter.

## Examples

### Non-compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Non-compliant configuration exhibiting infinite-loop-cycle"
}
```

### Compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Remediated configuration conforming to security standards"
}
```
