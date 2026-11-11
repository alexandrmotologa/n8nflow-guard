BINARY_NAME=n8nflow-guard
VERSION?=1.0.0
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE?=$(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-s -w -X github.com/alexandrmotologa/n8nflow-guard/internal/cli.Version=$(VERSION) -X github.com/alexandrmotologa/n8nflow-guard/internal/cli.GitCommit=$(COMMIT) -X github.com/alexandrmotologa/n8nflow-guard/internal/cli.BuildDate=$(BUILD_DATE)"

.PHONY: all build test clean lint run install

all: test build

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/n8nflow-guard

install:
	go install $(LDFLAGS) ./cmd/n8nflow-guard

test:
	go test -v -race ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

clean:
	rm -rf bin/ coverage.out coverage.txt

lint:
	go run ./cmd/n8nflow-guard lint ./testdata/clean_workflow.json
