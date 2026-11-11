package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
	"gopkg.in/yaml.v3"
)

// RuleConfig contains per-rule customization options.
type RuleConfig struct {
	Enabled    *bool           `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Severity   models.Severity `json:"severity,omitempty" yaml:"severity,omitempty"`
	Properties map[string]any  `json:"properties,omitempty" yaml:"properties,omitempty"`
}

// CustomSecret allows defining proprietary organizational secret regexes.
type CustomSecret struct {
	Name        string `json:"name" yaml:"name"`
	Pattern     string `json:"pattern" yaml:"pattern"`
	Description string `json:"description" yaml:"description"`
}

// Config represents the complete runtime settings for n8nflow-guard.
type Config struct {
	Version         string                `json:"version" yaml:"version"`
	MinSeverity     models.Severity       `json:"minSeverity" yaml:"minSeverity"`
	FailOn          models.Severity       `json:"failOn" yaml:"failOn"`
	IgnorePaths     []string              `json:"ignorePaths" yaml:"ignorePaths"`
	IgnoreRules     []string              `json:"ignoreRules" yaml:"ignoreRules"`
	Rules           map[string]RuleConfig `json:"rules" yaml:"rules"`
	CustomSecrets   []CustomSecret        `json:"customSecrets" yaml:"customSecrets"`
	IncludeDisabled bool                  `json:"includeDisabled" yaml:"includeDisabled"`
}

// DefaultConfig returns safe production defaults.
func DefaultConfig() *Config {
	return &Config{
		Version:         "1.0",
		MinSeverity:     models.SeverityLow,
		FailOn:          models.SeverityHigh,
		IgnorePaths:     []string{".git/**", "node_modules/**", "vendor/**"},
		IgnoreRules:     []string{},
		Rules:           make(map[string]RuleConfig),
		CustomSecrets:   make([]CustomSecret, 0),
		IncludeDisabled: false,
	}
}

// LoadConfig reads and parses configuration from a specified file.
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	cfg := DefaultConfig()
	ext := strings.ToLower(filepath.Ext(filePath))

	if ext == ".json" {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("invalid JSON config in %s: %w", filePath, err)
		}
	} else {
		// Default to YAML
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("invalid YAML config in %s: %w", filePath, err)
		}
	}

	return cfg, nil
}

// FindAndLoadConfig searches common config locations in the workspace directory.
func FindAndLoadConfig(workDir string) (*Config, error) {
	candidates := []string{
		".n8nflow-guard.yaml",
		".n8nflow-guard.yml",
		".n8nflow-guard.json",
		"n8nflow-guard.yaml",
		"n8nflow-guard.yml",
		"n8nflow-guard.json",
	}

	for _, name := range candidates {
		fullPath := filepath.Join(workDir, name)
		if _, err := os.Stat(fullPath); err == nil {
			return LoadConfig(fullPath)
		}
	}

	return DefaultConfig(), nil
}

// IsRuleEnabled verifies if a rule is active according to ignore lists and overrides.
func (c *Config) IsRuleEnabled(ruleID string) bool {
	for _, ignored := range c.IgnoreRules {
		if strings.EqualFold(ignored, ruleID) {
			return false
		}
	}

	if rc, exists := c.Rules[ruleID]; exists && rc.Enabled != nil {
		return *rc.Enabled
	}

	return true
}

// GetRuleSeverity returns the configured severity or falls back to default.
func (c *Config) GetRuleSeverity(ruleID string, defaultSeverity models.Severity) models.Severity {
	if rc, exists := c.Rules[ruleID]; exists && rc.Severity != "" {
		return rc.Severity
	}
	return defaultSeverity
}

// ShouldIgnorePath checks if a relative file path matches any ignore pattern.
func (c *Config) ShouldIgnorePath(targetPath string) bool {
	normalized := filepath.ToSlash(targetPath)
	for _, pattern := range c.IgnorePaths {
		match, err := filepath.Match(pattern, normalized)
		if err == nil && match {
			return true
		}
		// Check substring / directory prefix match
		trimmedPattern := strings.TrimSuffix(pattern, "/**")
		if strings.HasPrefix(normalized, trimmedPattern+"/") || normalized == trimmedPattern {
			return true
		}
	}
	return false
}
