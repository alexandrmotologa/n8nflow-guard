# PERF-002: unbounded-batch-processing

## Metadata
- Rule ID: PERF-002
- Category: Performance
- Default Severity: MEDIUM

## Description
Split In Batches node does not define a batch size or has batchSize <= 0.

## Why It Matters
Workflows in production handle customer data and sensitive integration points.
Violating this rule can result in unintended exposure, operational downtime, or runaway resource consumption.

## Remediation
Configure batchSize to an appropriate finite number (e.g. 50 to 500).

## Examples

### Non-compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Non-compliant configuration exhibiting unbounded-batch-processing"
}
```

### Compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Remediated configuration conforming to security standards"
}
```
