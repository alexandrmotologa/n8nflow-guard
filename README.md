# n8nflow-guard

Static analysis engine and security linter for n8n workflow JSON files.

n8nflow-guard checks n8n workflow definitions for hardcoded secrets, syntax errors in expressions, graph anomalies, missing error handlers, and reliability defects. It runs as a standalone binary in local environments, as a pre-commit hook, or in CI/CD pipelines with SARIF reports and GitHub Actions annotations.

## Features

- High throughput: Scans hundreds of workflow definitions in milliseconds with minimal memory footprint.
- 33 Built-in Rules: Covers security (SEC), reliability (REL), performance (PERF), best practices (BEST), and governance (GOV).
- Secret Detection: Combines Shannon entropy checks with regex signatures for 50+ cloud providers and credential formats.
- Graph & Expression Validation: Identifies cycles, orphan nodes, dead branches, and broken expression references.
- Automated Fixer (`--fix`): Sanitizes hardcoded credentials into `{{ $env.VAR_NAME }}` expressions and generates an `.env.example` template.
- CI/CD Ready: Native SARIF v2.1.0 output for GitHub Code Scanning, GitHub Actions workflow commands, JSON, and Markdown summaries.

## Installation

### Using Go Install

```bash
go install github.com/alexandrmotologa/n8nflow-guard/cmd/n8nflow-guard@latest
```

### From Pre-built Binaries

Download the appropriate binary for your platform (Linux, macOS, Windows) from the GitHub Releases page.

## Quick Start

```bash
# Lint all workflows in a directory
n8nflow-guard lint ./workflows/

# Fail build only on HIGH or CRITICAL severity findings
n8nflow-guard lint --fail-on HIGH ./workflows/

# Emit SARIF report for GitHub Code Scanning
n8nflow-guard lint --format sarif --output results.sarif ./workflows/

# Preview automated secret fixes without modifying files
n8nflow-guard fix --dry-run ./workflows/

# Apply automated fixes in place and export .env.example
n8nflow-guard fix --in-place --env-file .env.example ./workflows/

# List all 33 available rules
n8nflow-guard rules

# Initialize configuration in the current repository
n8nflow-guard init
```

## CLI Reference

### `lint`

```bash
n8nflow-guard lint [flags] [paths...]
```

Flags:
- `-f, --format string`: Output format: `terminal` (default), `sarif`, `github-actions`, `json`, `markdown`.
- `-s, --min-severity string`: Minimum severity level to report: `CRITICAL`, `HIGH`, `MEDIUM`, `LOW`, `INFO`.
- `--fail-on string`: Threshold severity that triggers non-zero exit code: `CRITICAL`, `HIGH`, `MEDIUM`, `LOW`.
- `-o, --output string`: Write report to file path instead of standard output.
- `--strict`: Treat warnings as errors (fail on `LOW` or higher).
- `--ignore-rule string`: Comma-separated list of rule IDs to skip (e.g. `BEST-002,REL-004`).
- `--include-disabled`: Evaluate disabled nodes on the canvas.

### `fix`

```bash
n8nflow-guard fix [flags] [paths...]
```

Flags:
- `-i, --in-place`: Apply modifications directly to workflow files on disk.
- `-b, --backup`: Create `.bak` backup files before modifying.
- `--dry-run`: Output a unified diff without changing files.
- `--strip-pindata`: Clear pinned execution test payloads (default `true`).
- `--env-file string`: Write extracted environment variables to a file (e.g. `.env.example`).

### `rules`

```bash
n8nflow-guard rules [flags]
```

Flags:
- `-r, --rule string`: Display detailed information and remediation steps for a specific rule.
- `--category string`: Filter rules by category (`Security`, `Reliability`, `Performance`, `BestPractice`, `Governance`).
- `--format string`: Display format: `table` (default), `json`, `markdown`.

## Built-in Rule Taxonomy

### Security (SEC)
| Rule ID | Severity | Name | Description |
|---|---|---|---|
| `SEC-001` | CRITICAL | hardcoded-secret | Detects hardcoded API keys, passwords, or high-entropy tokens in node parameters. |
| `SEC-002` | HIGH | unauthenticated-webhook | Webhook trigger is exposed to the public internet without authentication. |
| `SEC-003` | HIGH | plaintext-credential-object | Node contains inline plaintext credentials object instead of vault reference. |
| `SEC-004` | MEDIUM | insecure-http-url | HTTP Request node transmits unencrypted data using http:// protocol. |
| `SEC-005` | HIGH | disabled-ssl-verification | HTTP Request node disables TLS/SSL certificate verification. |
| `SEC-006` | HIGH | dangerous-command-injection | Execute Command node incorporates unparameterized expressions. |
| `SEC-007` | MEDIUM | production-log-data-leak | Workflow persists all execution data on success in production. |
| `SEC-008` | MEDIUM | permissive-cors-header | Webhook node sends a wildcard Access-Control-Allow-Origin: * header. |
| `SEC-009` | HIGH | unparameterized-sql-query | Database node interpolates raw expressions inside SQL strings. |
| `SEC-010` | MEDIUM | weak-cryptographic-hash | Crypto node employs deprecated hash algorithms such as MD5 or SHA1. |

### Reliability (REL)
| Rule ID | Severity | Name | Description |
|---|---|---|---|
| `REL-001` | HIGH | missing-error-workflow | Workflow has no global error handler workflow assigned. |
| `REL-002` | MEDIUM | missing-retry-on-fail | External network call node lacks retry-on-failure configuration. |
| `REL-003` | MEDIUM | unhandled-failure-continue | Node continues on failure without dedicated error output connections. |
| `REL-004` | MEDIUM | missing-execution-timeout | Workflow has no explicit execution timeout configured. |
| `REL-005` | HIGH | malformed-expression | Parameter contains an invalid or unclosed n8n expression syntax. |
| `REL-006` | HIGH | reference-to-missing-node | Expression references a node name that does not exist in the workflow. |
| `REL-007` | MEDIUM | disabled-critical-node | A disabled node sits between active connected nodes, breaking flow. |
| `REL-008` | MEDIUM | unhandled-webhook-error | Webhook node expects a response node, but none exists in the workflow. |
| `REL-009` | LOW | missing-database-timeout | Database query node lacks explicit statement timeout configuration. |
| `REL-010` | MEDIUM | missing-conditional-fallback | If or Switch node leaves branch outputs completely unconnected. |

### Performance (PERF)
| Rule ID | Severity | Name | Description |
|---|---|---|---|
| `PERF-001` | HIGH | infinite-loop-cycle | Workflow graph contains a directed cycle risking infinite execution. |
| `PERF-002` | MEDIUM | unbounded-batch-processing | Split In Batches node does not define a valid batch size. |
| `PERF-003` | MEDIUM | excessive-execution-logging | Workflow persists execution progress step by step in production. |
| `PERF-004` | MEDIUM | n-plus-one-request-pattern | HTTP Request node executes in a per-item loop without batching. |
| `PERF-005` | LOW | excessive-polling-frequency | Interval or polling trigger executes more frequently than once every 10s. |

### Best Practices (BEST)
| Rule ID | Severity | Name | Description |
|---|---|---|---|
| `BEST-001` | LOW | orphan-node | Active node has zero incoming and zero outgoing connections. |
| `BEST-002` | LOW | default-node-name | Node retains a generic default name (e.g. 'HTTP Request', 'Code 1'). |
| `BEST-003` | LOW | missing-node-description | Complex custom Code node lacks inline documentation or notes. |
| `BEST-004` | LOW | excessive-canvas-complexity | Workflow canvas exceeds 50 nodes without modular sub-workflow decomposition. |
| `BEST-005` | MEDIUM | outdated-node-version | Node uses a legacy typeVersion that lacks security and performance updates. |
| `BEST-006` | MEDIUM | pinned-test-data-leak | Workflow includes committed pinData test payloads in version control. |

### Governance (GOV)
| Rule ID | Severity | Name | Description |
|---|---|---|---|
| `GOV-001` | LOW | missing-workflow-tags | Workflow is missing organizational taxonomy tags in version control. |
| `GOV-002` | MEDIUM | missing-trigger-node | Workflow has no entrypoint trigger node to initiate execution. |

## CI/CD Integration

### GitHub Actions

Add the following workflow to `.github/workflows/n8n-lint.yml`:

```yaml
name: n8n Workflow Security & Reliability

on:
  pull_request:
    paths:
      - '**.json'

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Run n8nflow-guard
        uses: alexandrmotologa/n8nflow-guard@main
        with:
          path: './workflows'
          fail-on: 'HIGH'
          format: 'github-actions'
          sarif-file: 'results.sarif'

      - name: Upload SARIF to GitHub Code Scanning
        uses: github/codeql-action/upload-sarif@v3
        if: always()
        with:
          sarif_file: 'results.sarif'
```

### Pre-commit Hook

Add to your `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/alexandrmotologa/n8nflow-guard
    rev: v1.0.0
    hooks:
      - id: n8nflow-guard-lint
        files: \.json$
```

## Configuration File

Customize rule severities and exclusions in `.n8nflow-guard.yaml`:

```yaml
version: "1.0"
minSeverity: LOW
failOn: HIGH
ignorePaths:
  - ".git/**"
  - "archive/**"
ignoreRules:
  - "BEST-002"
rules:
  SEC-001:
    severity: CRITICAL
    enabled: true
  REL-001:
    severity: HIGH
    enabled: true
```

## Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on code style, testing requirements, and development setup.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
