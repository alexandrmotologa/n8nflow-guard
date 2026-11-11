# Contributing to n8nflow-guard

Thank you for your interest in contributing to n8nflow-guard. This document outlines development guidelines, code standards, and testing procedures.

## Prerequisites

- Go 1.27 or higher
- Git
- Make (optional, for convenience)

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/alexandrmotologa/n8nflow-guard.git
   cd n8nflow-guard
   ```

2. Run test suites:
   ```bash
   go test -v -race ./...
   ```

3. Build the binary locally:
   ```bash
   go build -o bin/n8nflow-guard ./cmd/n8nflow-guard
   ```

## Adding a New Rule

When creating a new static analysis rule:
1. Define the rule struct implementing the `Rule` interface in `internal/rules/`.
2. Assign a unique rule ID following the category taxonomy (`SEC-xxx`, `REL-xxx`, `PERF-xxx`, `BEST-xxx`, `GOV-xxx`).
3. Register the rule in `RegisterAllRules` inside `internal/rules/registry.go`.
4. Add unit test coverage in `internal/rules/rules_test.go`.
5. Create a documentation file in `rules_docs/<RULE_ID>.md` containing non-compliant and compliant JSON snippets.
6. Update rule counts and documentation in `README.md` and `MASTER_PLAN.md`.

## Coding Guidelines

- Keep code free of unnecessary third-party dependencies. Prefer the Go standard library.
- Keep execution time fast and memory usage below 20MB.
- Ensure all exported types, functions, and methods have clear Go doc comments.
- Run `go test -v -race ./...` before submitting a pull request.
