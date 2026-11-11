package engine

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/baseline"
	"github.com/alexandrmotologa/n8nflow-guard/internal/config"
	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
	"github.com/alexandrmotologa/n8nflow-guard/internal/parser"
	"github.com/alexandrmotologa/n8nflow-guard/internal/reporter"
	"github.com/alexandrmotologa/n8nflow-guard/internal/rules"
)

// ScanOptions configures the static analysis execution.
type ScanOptions struct {
	Paths                []string
	MinSeverity          models.Severity
	FailOn               models.Severity
	Strict               bool
	Staged               bool
	Since                string
	BaselinePath         string
	GenerateBaselinePath string
	Config               *config.Config
	Registry             *rules.Registry
}

// ScanResult contains all findings, summaries, and exit criteria.
type ScanResult struct {
	Issues      []models.Issue
	Summary     reporter.Summary
	HasFailures bool
}

// RunScan scans the supplied paths and evaluates all enabled rules.
func RunScan(opts ScanOptions) (*ScanResult, error) {
	if opts.Config == nil {
		opts.Config = config.DefaultConfig()
	}
	if opts.Registry == nil {
		opts.Registry = rules.DefaultRegistry()
	}

	// Apply severity overrides from options if provided
	if opts.MinSeverity != "" {
		opts.Config.MinSeverity = opts.MinSeverity
	}
	if opts.FailOn != "" {
		opts.Config.FailOn = opts.FailOn
	}

	var files []string
	var err error
	if opts.Staged {
		files, err = GetStagedWorkflowFiles()
		if err != nil {
			return nil, err
		}
	} else if opts.Since != "" {
		files, err = GetDiffWorkflowFiles(opts.Since)
		if err != nil {
			return nil, err
		}
	} else {
		files, err = resolveFilePaths(opts.Paths, opts.Config)
		if err != nil {
			return nil, fmt.Errorf("failed resolving input paths: %w", err)
		}
	}

	if len(files) == 0 {
		return &ScanResult{
			Issues:      nil,
			Summary:     reporter.Summary{},
			HasFailures: false,
		}, nil
	}

	var allIssues []models.Issue
	totalWorkflows := 0

	for _, filePath := range files {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed reading file %s: %w", filePath, err)
		}

		workflows, err := parser.ParseWorkflows(data, filePath)
		if err != nil {
			// If JSON is invalid or not an n8n workflow, report as a syntax error issue
			allIssues = append(allIssues, models.Issue{
				RuleID:      "REL-005",
				RuleName:    "malformed-json",
				Category:    models.CategoryReliability,
				Severity:    models.SeverityHigh,
				Message:     fmt.Sprintf("Failed parsing JSON file: %v", err),
				FilePath:    filePath,
				Line:        1,
				Remediation: "Fix invalid JSON formatting.",
			})
			continue
		}

		totalWorkflows += len(workflows)
		lineMapper := parser.NewLineMapper(data)

		for _, wf := range workflows {
			graph := parser.BuildGraph(wf)
			ctx := &rules.RuleContext{
				Workflow:   wf,
				Graph:      graph,
				LineMapper: lineMapper,
				Config:     opts.Config,
			}

			issues := opts.Registry.Evaluate(ctx)
			allIssues = append(allIssues, issues...)
		}
	}

	// Generate baseline file if requested
	if opts.GenerateBaselinePath != "" {
		if err := baseline.GenerateBaseline(allIssues, opts.GenerateBaselinePath); err != nil {
			return nil, fmt.Errorf("failed generating baseline file: %w", err)
		}
	}

	// Filter against baseline if supplied
	if opts.BaselinePath != "" {
		_, baselineIndex, err := baseline.LoadBaseline(opts.BaselinePath)
		if err != nil {
			return nil, fmt.Errorf("failed loading baseline: %w", err)
		}
		allIssues = baseline.FilterNewIssues(allIssues, baselineIndex)
	}

	// Sort issues by file path, line, and severity rank descending
	sort.Slice(allIssues, func(i, j int) bool {
		if allIssues[i].FilePath != allIssues[j].FilePath {
			return allIssues[i].FilePath < allIssues[j].FilePath
		}
		if allIssues[i].Line != allIssues[j].Line {
			return allIssues[i].Line < allIssues[j].Line
		}
		return allIssues[i].Severity.Rank() > allIssues[j].Severity.Rank()
	})

	summary := reporter.ComputeSummary(len(files), totalWorkflows, allIssues)

	// Determine if exit criteria should fail the build
	hasFailures := false
	failRank := opts.Config.FailOn.Rank()
	if failRank == 0 {
		failRank = models.SeverityHigh.Rank()
	}
	if opts.Strict {
		failRank = models.SeverityLow.Rank()
	}

	for _, iss := range allIssues {
		if iss.Severity.Rank() >= failRank {
			hasFailures = true
			break
		}
	}

	return &ScanResult{
		Issues:      allIssues,
		Summary:     summary,
		HasFailures: hasFailures,
	}, nil
}

// resolveFilePaths traverses directories and expands globs into concrete file lists.
func resolveFilePaths(paths []string, cfg *config.Config) ([]string, error) {
	if len(paths) == 0 {
		paths = []string{"."}
	}

	foundMap := make(map[string]bool)

	for _, p := range paths {
		// Handle glob patterns
		if strings.ContainsAny(p, "*?[]") {
			matches, err := filepath.Glob(p)
			if err != nil {
				return nil, err
			}
			for _, m := range matches {
				if shouldIncludeFile(m, cfg) {
					foundMap[m] = true
				}
			}
			continue
		}

		info, err := os.Stat(p)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			if shouldIncludeFile(p, cfg) {
				foundMap[p] = true
			}
			continue
		}

		// Traverse directory recursively
		err = filepath.WalkDir(p, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				name := d.Name()
				if name == ".git" || name == "node_modules" || name == "vendor" {
					return filepath.SkipDir
				}
				if cfg != nil && cfg.ShouldIgnorePath(path) {
					return filepath.SkipDir
				}
				return nil
			}

			if shouldIncludeFile(path, cfg) {
				foundMap[path] = true
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	var results []string
	for f := range foundMap {
		results = append(results, f)
	}
	sort.Strings(results)

	return results, nil
}

func shouldIncludeFile(path string, cfg *config.Config) bool {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".json" {
		return false
	}
	if cfg != nil && cfg.ShouldIgnorePath(path) {
		return false
	}
	return true
}
