# PERF-005: excessive-polling-frequency

## Metadata
- Rule ID: PERF-005
- Category: Performance
- Default Severity: LOW

## Description
Interval or polling trigger executes more frequently than once every 10 seconds.

## Why It Matters
Workflows in production handle customer data and sensitive integration points.
Violating this rule can result in unintended exposure, operational downtime, or runaway resource consumption.

## Remediation
Increase polling interval to >= 60 seconds or switch to webhook triggers.

## Examples

### Non-compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Non-compliant configuration exhibiting excessive-polling-frequency"
}
```

### Compliant Configuration
```json
{
  "name": "Sample Node",
  "notes": "Remediated configuration conforming to security standards"
}
```
