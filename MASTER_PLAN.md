# n8nflow-guard: Architecture & Engineering Master Plan
> Static Analysis, Security Linter & CI/CD Pipeline Validator for n8n Workflows  
> Target Footprint: Sub-20MB RAM Single Go Binary CLI, Sub-100ms Scan, GitHub Action & Pre-Commit Hook

---

## 1. Executive Summary & Market Problem

### 1.1 The Operational Blindspot in Enterprise n8n Workflows
As organizations adopt n8n for mission-critical automation, workflow definitions transition from visual canvas experiments to large JSON documents stored in Git repositories (`workflow.json`).

Reviewing a 4,000-line n8n JSON file in a pull request is error-prone for human reviewers. Common production issues include:
1. Hardcoded Secrets: API keys, database passwords, and tokens embedded in node parameters during local testing and committed directly to source control.
2. Missing Error Handling: Workflows executing without a configured error handler, causing silent failures that drop business transactions without alerts.
3. Graph Topology Defects: Disconnected orphan nodes, broken expression references, and unintended recursion cycles that trigger runaway execution loops.
4. Operational Data Leaks: Production workflows committed with execution data retention and pinned test data containing customer information.
5. Inconsistent Governance: Lack of automated guardrails in CI pipelines to enforce organizational standards before deployment to n8n production instances.

### 1.2 The n8nflow-guard Solution
n8nflow-guard is a standalone static analysis engine and security linter written in Go. It inspects n8n workflow JSON files before deployment:
- Zero Runtime Dependencies: Compiles into a single binary that scans hundreds of workflows in under 50 milliseconds with minimal memory footprint.
- 33 Built-in Rules: Covers Security, Reliability, Performance, Best Practices, and Governance.
- Graph & Expression Analysis: Builds an in-memory directed graph of nodes and connections to evaluate cycles, reachability, and expression syntax.
- Shannon Entropy & Pattern Scanning: Detects credentials across 50+ provider formats and high-entropy secrets in string fields.
- Automated Fixer: Sanitizes hardcoded credentials into environment variable expressions and generates matching `.env.example` templates.
- Native CI/CD Integration: Emits SARIF v2.1.0 for GitHub Code Scanning, GitHub Actions workflow commands for PR diff annotations, and formatted console tables.

---

## 2. High-Level Architecture

```
                          +----------------------------------------------------+
                          |              n8n Workflow JSON Files               |
                          |     (Single Files, Directories, Glob Patterns)     |
                          +----------------------------------------------------+
                                                   |
                                                   v
                          +----------------------------------------------------+
                          |                n8nflow-guard CLI                   |
                          |           (Cobra / Go Engine Core)                 |
                          +----------------------------------------------------+
                                                   |
                   +-------------------------------+-------------------------------+
                   |                               |                               |
                   v                               v                               v
    +------------------------------+ +------------------------------+ +------------------------------+
    |      JSON & Line Mapper      | |      DAG & Topology Engine   | |       Secret Scanner         |
    | - Validates JSON structure   | | - Builds directed graph      | | - Shannon entropy analysis   |
    | - Preserves exact node and   | | - Detects recursion cycles   | | - 50+ provider regex rules   |
    |   field line coordinates     | | - Finds orphan nodes         | | - Parameter tree inspection  |
    | - Evaluates expressions      | | - Traces failure branches    | | - Plaintext credential audit |
    +------------------------------+ +------------------------------+ +------------------------------+
                   |                               |                               |
                   +-------------------------------+-------------------------------+
                                                   |
                                                   v
                          +----------------------------------------------------+
                          |                Rule Evaluation Engine              |
                          |  - Security Suite (SEC-001 to SEC-010)             |
                          |  - Reliability Suite (REL-001 to REL-010)          |
                          |  - Performance Suite (PERF-001 to PERF-005)        |
                          |  - Best Practices Suite (BEST-001 to BEST-006)     |
                          |  - Governance Suite (GOV-001 to GOV-002)           |
                          +----------------------------------------------------+
                                                   |
                   +-------------------------------+-------------------------------+
                   |                                                               |
                   v                                                               v
    +------------------------------+                                +------------------------------+
    |       Reporter Layer         |                                |         Fixer Engine         |
    | - Terminal (ANSI Tables)     |                                | - Replaces secrets with      |
    | - SARIF v2.1.0 (GitHub)      |                                |   {{ $env.VAR_NAME }} syntax |
    | - GitHub Action Annotations  |                                | - Strips pinData             |
    | - JSON / Markdown / CSV      |                                | - Exports .env.example       |
    +------------------------------+                                +------------------------------+
```

---

## 3. Comprehensive Rule Taxonomy (33 Rules)

### 3.1 Security Rules (SEC)
| Rule ID | Severity | Name | Description |
|---|---|---|---|
| SEC-001 | CRITICAL | hardcoded-secret | Hardcoded API keys, bearer tokens, or credentials found in node parameters. |
| SEC-002 | HIGH | unauthenticated-webhook | Webhook trigger configured with authentication set to none. |
| SEC-003 | HIGH | plaintext-credential-object | Plaintext credential payload embedded in parameters instead of n8n credentials vault. |
| SEC-004 | MEDIUM | insecure-http-url | HTTP Request node uses unencrypted http:// transport instead of https://. |
| SEC-005 | HIGH | disabled-ssl-verification | HTTP Request node has allowUnauthorizedCerts or ignoreSSL enabled. |
| SEC-006 | HIGH | dangerous-command-injection | Execute Command or shell execution node contains unparameterized inputs. |
| SEC-007 | MEDIUM | production-log-data-leak | Workflow configured to persist all execution data in production (saveDataSuccessExecution=all). |
| SEC-008 | MEDIUM | permissive-cors-header | Webhook node responds with wildcard Access-Control-Allow-Origin: * header. |
| SEC-009 | HIGH | unparameterized-sql-query | Database node uses string concatenation or raw expression interpolation in query fields. |
| SEC-010 | MEDIUM | weak-cryptographic-hash | Crypto node uses deprecated algorithms such as MD5 or SHA1. |

### 3.2 Reliability Rules (REL)
| Rule ID | Severity | Name | Description |
|---|---|---|---|
| REL-001 | HIGH | missing-error-workflow | Production workflow missing an assigned errorWorkflow handler. |
| REL-002 | MEDIUM | missing-retry-on-fail | External HTTP or API node missing retry-on-failure configuration. |
| REL-003 | MEDIUM | unhandled-failure-continue | Node configured with continueOnFail=true without downstream failure branch. |
| REL-004 | MEDIUM | missing-execution-timeout | Workflow executionTimeout setting is missing or set to 0. |
| REL-005 | HIGH | malformed-expression | Malformed n8n expression syntax or unclosed curly braces. |
| REL-006 | HIGH | reference-to-missing-node | Expression references a node name not present in the workflow graph. |
| REL-007 | MEDIUM | disabled-critical-node | Disabled node sits along the execution path between trigger and terminal nodes. |
| REL-008 | MEDIUM | unhandled-webhook-error | Webhook node lacks an explicit error response or error handler connection. |
| REL-009 | LOW | missing-database-timeout | Database query node lacks explicit connection or query timeout settings. |
| REL-010 | MEDIUM | missing-conditional-fallback | If or Switch node branches terminate without default or fallback handlers. |
| REL-011 | LOW | unsafe-property-access | Deeply nested property access without optional chaining (?.) risking runtime exceptions. |

### 3.3 Performance Rules (PERF)
| Rule ID | Severity | Name | Description |
|---|---|---|---|
| PERF-001 | HIGH | infinite-loop-cycle | Graph contains a directed cycle without an explicit counter or termination node. |
| PERF-002 | MEDIUM | unbounded-batch-processing | Split In Batches node configured without batch size limit or throttling delay. |
| PERF-003 | MEDIUM | excessive-execution-logging | saveExecutionProgress enabled on high-throughput workflow. |
| PERF-004 | MEDIUM | n-plus-one-request-pattern | HTTP Request or database node executes inside a single-item loop. |
| PERF-005 | LOW | excessive-polling-frequency | Interval or Poll trigger interval set below 10 seconds. |

### 3.4 Best Practices (BEST)
| Rule ID | Severity | Name | Description |
|---|---|---|---|
| BEST-001 | LOW | orphan-node | Disconnected node with zero incoming and zero outgoing connections. |
| BEST-002 | LOW | default-node-name | Node retains default n8n name (e.g., "HTTP Request", "Code 1"). |
| BEST-003 | LOW | missing-node-description | Complex Code or HTTP Request node contains empty notes or documentation. |
| BEST-004 | LOW | excessive-canvas-complexity | Workflow contains more than 50 nodes without sub-workflow decomposition. |
| BEST-005 | MEDIUM | outdated-node-version | Node typeVersion is significantly behind current n8n version. |
| BEST-006 | MEDIUM | pinned-test-data-leak | Workflow contains committed pinData payloads in version control. |

### 3.5 Governance (GOV)
| Rule ID | Severity | Name | Description |
|---|---|---|---|
| GOV-001 | LOW | missing-workflow-tags | Workflow JSON lacks metadata tags for organizational categorization. |
| GOV-002 | MEDIUM | missing-trigger-node | Workflow contains no trigger or webhook node and cannot run automatically. |

---

## 4. Repository Structure

```
n8nflow-guard/
├── .github/
│   ├── workflows/
│   │   ├── ci.yml                  # Build, test, lint, and coverage
│   │   └── release.yml             # Multi-arch binary releases via GoReleaser
│   └── pull_request_template.md
├── cmd/
│   └── n8nflow-guard/
│       └── main.go                 # Application entrypoint
├── internal/
│   ├── baseline/                   # Fingerprinting and suppression of baseline technical debt
│   ├── cli/                        # Cobra commands (lint, fix, rules, graph, remote, init, version)
│   ├── config/                     # Configuration loader (.n8nflow-guard.yaml)
│   ├── engine/                     # Traversal engine, git diff/staged resolver
│   ├── fixer/                      # Automated remediation engine
│   ├── graphviz/                   # Mermaid, Graphviz DOT, and ASCII visualizers
│   ├── models/                     # Workflow, Node, DAG, and Issue data structures
│   ├── parser/                     # JSON parser, line mapper, DAG, and expressions
│   ├── remote/                     # Live n8n instance audit REST client
│   ├── reporter/                   # Terminal, HTML dashboard, SARIF, GitHub Actions, JSON, Markdown
│   ├── rules/                      # Rule registry, inline notes suppression, and rule implementations
│   ├── scoring/                    # Security score calculation and OWASP API Top 10 mapping
│   └── secrets/                    # Shannon entropy and regex patterns
├── rules_docs/                     # Markdown documentation for all 34 rules
├── testdata/                       # Sample valid and invalid workflow fixtures
├── action.yml                      # GitHub Action definition
├── .pre-commit-hooks.yaml          # Pre-commit hook configuration
├── Makefile                        # Build, test, lint targets
├── go.mod
├── go.sum
├── LICENSE                         # MIT License
└── README.md                       # Documentation
```

---

## 5. CLI Command Reference

```bash
# Scan workflows with terminal output
n8nflow-guard lint ./workflows/

# Fail build only on HIGH or CRITICAL issues
n8nflow-guard lint --min-severity HIGH ./workflows/

# Export SARIF for GitHub Code Scanning
n8nflow-guard lint --format sarif --output results.sarif ./workflows/

# Emit GitHub Actions workflow commands
n8nflow-guard lint --format github-actions ./workflows/

# Automatically sanitize hardcoded secrets and create .env.example
n8nflow-guard fix --in-place --env-file .env.example ./workflows/

# Preview fixes without modifying files
n8nflow-guard fix --dry-run ./workflows/

# List all available rules
n8nflow-guard rules

# View documentation for a specific rule
n8nflow-guard rules --rule SEC-001

# Generate starter configuration file
n8nflow-guard init
```

---

## 6. Detailed Implementation Roadmap

### Phase 1: Core Foundation & Parser
- Set up Go module, directory hierarchy, Makefile, and CI workflows.
- Implement models for n8n workflows, nodes, connections, settings, and findings.
- Build JSON parser with line-number tracking for precise source location.
- Implement Directed Acyclic Graph (DAG) builder and cycle detection algorithm.
- Build n8n expression parser for `{{ $json... }}` syntax validation.

### Phase 2: Secret Scanner & Security Suite (SEC-001 to SEC-010)
- Implement Shannon entropy calculator with character distribution heuristics.
- Implement regex pattern database covering 50+ API key and credential formats.
- Implement all 10 security rules with comprehensive unit test coverage.

### Phase 3: Reliability, Performance, Best Practice & Governance Suites
- Implement reliability rules (REL-001 to REL-010) including error workflows and retry policies.
- Implement performance rules (PERF-001 to PERF-005) including graph cycle verification.
- Implement best practice rules (BEST-001 to BEST-006) including orphan node detection.
- Implement governance rules (GOV-001 to GOV-002) for trigger and tag compliance.

### Phase 4: Auto-Fixer, Config System & Multi-Format Reporters
- Build configuration loader with YAML and JSON support and rule overrides.
- Build `--fix` engine for credential sanitization and `.env.example` generation.
- Build terminal reporter with colored tables and severity badges.
- Build SARIF v2.1.0 reporter for GitHub Security tab integration.
- Build GitHub Actions, JSON, and Markdown reporters.

### Phase 5: CLI Surface, GitHub Action, Test Fixtures & Documentation
- Implement Cobra CLI commands: `lint`, `fix`, `rules`, `init`, and `version`.
- Create comprehensive test fixtures in `testdata/`.
- Author documentation files for all 33 rules in `rules_docs/`.
- Create `action.yml` and `.pre-commit-hooks.yaml`.
- Write humanized `README.md` and project documentation.
- Run end-to-end integration tests, memory benchmarks, and CI checks.
