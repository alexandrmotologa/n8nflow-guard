package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alexandrmotologa/n8nflow-guard/internal/models"
	"github.com/alexandrmotologa/n8nflow-guard/internal/secrets"
)

// Helper to recursively traverse parameter maps looking for string values
func walkParameters(params map[string]interface{}, prefix string, fn func(keyPath string, val string)) {
	for k, v := range params {
		fullPath := k
		if prefix != "" {
			fullPath = prefix + "." + k
		}
		switch val := v.(type) {
		case string:
			fn(fullPath, val)
		case map[string]interface{}:
			walkParameters(val, fullPath, fn)
		case []interface{}:
			for i, item := range val {
				itemPath := fmt.Sprintf("%s[%d]", fullPath, i)
				if strItem, ok := item.(string); ok {
					fn(itemPath, strItem)
				} else if mapItem, ok := item.(map[string]interface{}); ok {
					walkParameters(mapItem, itemPath, fn)
				}
			}
		}
	}
}

// -------------------------------------------------------------
// SEC-001: hardcoded-secret
// -------------------------------------------------------------
type RuleSEC001 struct{}

func (r *RuleSEC001) ID() string             { return "SEC-001" }
func (r *RuleSEC001) Name() string           { return "hardcoded-secret" }
func (r *RuleSEC001) Category() models.Category { return models.CategorySecurity }
func (r *RuleSEC001) Severity() models.Severity { return models.SeverityCritical }
func (r *RuleSEC001) Description() string {
	return "Detects hardcoded API keys, passwords, or high-entropy tokens in node parameters."
}
func (r *RuleSEC001) Remediation() string {
	return "Move the sensitive value into an n8n credential vault or environment variable ({{ $env.VAR_NAME }})."
}
func (r *RuleSEC001) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if node.Disabled && (ctx.Config == nil || !ctx.Config.IncludeDisabled) {
			continue
		}

		walkParameters(node.Parameters, "", func(keyPath string, val string) {
			trimmed := strings.TrimSpace(val)
			if strings.HasPrefix(trimmed, "=") {
				trimmed = strings.TrimPrefix(trimmed, "=")
			}
			if strings.HasPrefix(trimmed, "{{") && strings.HasSuffix(trimmed, "}}") {
				return
			}

			// Scan against known patterns
			matches := secrets.ScanString(val)
			for _, m := range matches {
				line := ctx.LineMapper.FindFieldLine(node.Name, keyPath)
				issues = append(issues, models.Issue{
					RuleID:       r.ID(),
					RuleName:     r.Name(),
					Category:     r.Category(),
					Severity:     r.Severity(),
					Message:      fmt.Sprintf("Node '%s' contains hardcoded %s (%s)", node.Name, m.PatternName, secrets.MaskSecret(m.MatchedText)),
					FilePath:     ctx.Workflow.FilePath,
					Line:         line,
					NodeID:       node.ID,
					NodeName:     node.Name,
					Field:        keyPath,
					Remediation:  r.Remediation(),
					Fixable:      true,
					FixSuggested: "{{ $env." + suggestEnvVarName(node.Name, keyPath) + " }}",
				})
			}

			// Check high entropy on sensitive key names
			if len(matches) == 0 && secrets.IsSensitiveKeyName(keyPath) && secrets.IsHighEntropySecret(val) {
				line := ctx.LineMapper.FindFieldLine(node.Name, keyPath)
				issues = append(issues, models.Issue{
					RuleID:       r.ID(),
					RuleName:     r.Name(),
					Category:     r.Category(),
					Severity:     r.Severity(),
					Message:      fmt.Sprintf("Node '%s' parameter '%s' contains high-entropy credential (%s)", node.Name, keyPath, secrets.MaskSecret(val)),
					FilePath:     ctx.Workflow.FilePath,
					Line:         line,
					NodeID:       node.ID,
					NodeName:     node.Name,
					Field:        keyPath,
					Remediation:  r.Remediation(),
					Fixable:      true,
					FixSuggested: "{{ $env." + suggestEnvVarName(node.Name, keyPath) + " }}",
				})
			}
		})
	}

	return issues
}

func suggestEnvVarName(nodeName, keyPath string) string {
	cleanNode := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(strings.ToUpper(nodeName), "_")
	cleanKey := regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(strings.ToUpper(keyPath), "_")
	return strings.Trim(fmt.Sprintf("%s_%s", cleanNode, cleanKey), "_")
}

// -------------------------------------------------------------
// SEC-002: unauthenticated-webhook
// -------------------------------------------------------------
type RuleSEC002 struct{}

func (r *RuleSEC002) ID() string             { return "SEC-002" }
func (r *RuleSEC002) Name() string           { return "unauthenticated-webhook" }
func (r *RuleSEC002) Category() models.Category { return models.CategorySecurity }
func (r *RuleSEC002) Severity() models.Severity { return models.SeverityHigh }
func (r *RuleSEC002) Description() string {
	return "Webhook trigger is exposed to the public internet without authentication."
}
func (r *RuleSEC002) Remediation() string {
	return "Configure authentication (headerAuth, basicAuth, or jwtAuth) in the webhook parameters."
}
func (r *RuleSEC002) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if !strings.Contains(strings.ToLower(node.Type), "webhook") {
			continue
		}
		if node.Disabled && (ctx.Config == nil || !ctx.Config.IncludeDisabled) {
			continue
		}

		auth, ok := node.Parameters["authentication"].(string)
		if !ok || auth == "none" || strings.TrimSpace(auth) == "" {
			line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Webhook node '%s' has authentication disabled ('none')", node.Name),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "parameters.authentication",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// SEC-003: plaintext-credential-object
// -------------------------------------------------------------
type RuleSEC003 struct{}

func (r *RuleSEC003) ID() string             { return "SEC-003" }
func (r *RuleSEC003) Name() string           { return "plaintext-credential-object" }
func (r *RuleSEC003) Category() models.Category { return models.CategorySecurity }
func (r *RuleSEC003) Severity() models.Severity { return models.SeverityHigh }
func (r *RuleSEC003) Description() string {
	return "Node contains inline plaintext credentials object instead of referencing an n8n credential ID."
}
func (r *RuleSEC003) Remediation() string {
	return "Link the node to a managed credential stored in the n8n credential store."
}
func (r *RuleSEC003) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		for credType, credVal := range node.Credentials {
			credMap, ok := credVal.(map[string]interface{})
			if !ok {
				continue
			}
			// Safe credentials in n8n only store "id" and "name". If "data", "password", or "token" exist, it's plaintext
			for key := range credMap {
				if key == "data" || key == "password" || key == "secret" || key == "token" {
					line := ctx.LineMapper.FindFieldLine(node.Name, "credentials")
					issues = append(issues, models.Issue{
						RuleID:      r.ID(),
						RuleName:    r.Name(),
						Category:    r.Category(),
						Severity:    r.Severity(),
						Message:     fmt.Sprintf("Node '%s' contains raw credential data for '%s' instead of a vault ID reference", node.Name, credType),
						FilePath:    ctx.Workflow.FilePath,
						Line:        line,
						NodeID:      node.ID,
						NodeName:    node.Name,
						Field:       "credentials." + credType,
						Remediation: r.Remediation(),
					})
					break
				}
			}
		}
	}

	return issues
}

// -------------------------------------------------------------
// SEC-004: insecure-http-url
// -------------------------------------------------------------
type RuleSEC004 struct{}

func (r *RuleSEC004) ID() string             { return "SEC-004" }
func (r *RuleSEC004) Name() string           { return "insecure-http-url" }
func (r *RuleSEC004) Category() models.Category { return models.CategorySecurity }
func (r *RuleSEC004) Severity() models.Severity { return models.SeverityMedium }
func (r *RuleSEC004) Description() string {
	return "HTTP Request node transmits unencrypted data across the network using http:// protocol."
}
func (r *RuleSEC004) Remediation() string {
	return "Update the URL protocol from http:// to https://."
}
func (r *RuleSEC004) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if !strings.Contains(strings.ToLower(node.Type), "httprequest") {
			continue
		}

		rawURL, ok := node.Parameters["url"].(string)
		if !ok {
			continue
		}

		lower := strings.ToLower(strings.TrimSpace(rawURL))
		if strings.HasPrefix(lower, "http://") {
			// Allow local loopback testing
			if strings.Contains(lower, "localhost") || strings.Contains(lower, "127.0.0.1") {
				continue
			}

			line := ctx.LineMapper.FindFieldLine(node.Name, "url")
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("HTTP Request node '%s' uses insecure cleartext protocol: %s", node.Name, rawURL),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "parameters.url",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// SEC-005: disabled-ssl-verification
// -------------------------------------------------------------
type RuleSEC005 struct{}

func (r *RuleSEC005) ID() string             { return "SEC-005" }
func (r *RuleSEC005) Name() string           { return "disabled-ssl-verification" }
func (r *RuleSEC005) Category() models.Category { return models.CategorySecurity }
func (r *RuleSEC005) Severity() models.Severity { return models.SeverityHigh }
func (r *RuleSEC005) Description() string {
	return "HTTP Request node disables TLS/SSL certificate verification, exposing connections to man-in-the-middle attacks."
}
func (r *RuleSEC005) Remediation() string {
	return "Remove allowUnauthorizedCerts or set it to false."
}
func (r *RuleSEC005) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if !strings.Contains(strings.ToLower(node.Type), "httprequest") {
			continue
		}

		options, ok := node.Parameters["options"].(map[string]interface{})
		if !ok {
			continue
		}

		if val, exists := options["allowUnauthorizedCerts"]; exists && val == true {
			line := ctx.LineMapper.FindFieldLine(node.Name, "allowUnauthorizedCerts")
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Node '%s' has allowUnauthorizedCerts enabled", node.Name),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "parameters.options.allowUnauthorizedCerts",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// SEC-006: dangerous-command-injection
// -------------------------------------------------------------
type RuleSEC006 struct{}

func (r *RuleSEC006) ID() string             { return "SEC-006" }
func (r *RuleSEC006) Name() string           { return "dangerous-command-injection" }
func (r *RuleSEC006) Category() models.Category { return models.CategorySecurity }
func (r *RuleSEC006) Severity() models.Severity { return models.SeverityHigh }
func (r *RuleSEC006) Description() string {
	return "Execute Command node incorporates unparameterized expressions, posing a remote code execution risk."
}
func (r *RuleSEC006) Remediation() string {
	return "Sanitize and validate all incoming inputs before concatenating them into shell commands."
}
func (r *RuleSEC006) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if !strings.Contains(strings.ToLower(node.Type), "executecommand") {
			continue
		}

		cmd, ok := node.Parameters["command"].(string)
		if !ok {
			continue
		}

		if strings.Contains(cmd, "{{") && strings.Contains(cmd, "$json") {
			line := ctx.LineMapper.FindFieldLine(node.Name, "command")
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Execute Command node '%s' evaluates raw JSON expression in shell string", node.Name),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "parameters.command",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// SEC-007: production-log-data-leak
// -------------------------------------------------------------
type RuleSEC007 struct{}

func (r *RuleSEC007) ID() string             { return "SEC-007" }
func (r *RuleSEC007) Name() string           { return "production-log-data-leak" }
func (r *RuleSEC007) Category() models.Category { return models.CategorySecurity }
func (r *RuleSEC007) Severity() models.Severity { return models.SeverityMedium }
func (r *RuleSEC007) Description() string {
	return "Workflow persists all execution data on success in production, creating potential data leaks and DB bloat."
}
func (r *RuleSEC007) Remediation() string {
	return "Set saveDataSuccessExecution to 'none' in workflow settings."
}
func (r *RuleSEC007) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	if strings.EqualFold(ctx.Workflow.Settings.SaveDataSuccessExecution, "all") {
		line := ctx.LineMapper.FindFieldLine("settings", "saveDataSuccessExecution")
		issues = append(issues, models.Issue{
			RuleID:      r.ID(),
			RuleName:    r.Name(),
			Category:    r.Category(),
			Severity:    r.Severity(),
			Message:     "Workflow is configured to persist all successful execution payload data",
			FilePath:    ctx.Workflow.FilePath,
			Line:        line,
			Field:       "settings.saveDataSuccessExecution",
			Remediation: r.Remediation(),
		})
	}

	return issues
}

// -------------------------------------------------------------
// SEC-008: permissive-cors-header
// -------------------------------------------------------------
type RuleSEC008 struct{}

func (r *RuleSEC008) ID() string             { return "SEC-008" }
func (r *RuleSEC008) Name() string           { return "permissive-cors-header" }
func (r *RuleSEC008) Category() models.Category { return models.CategorySecurity }
func (r *RuleSEC008) Severity() models.Severity { return models.SeverityMedium }
func (r *RuleSEC008) Description() string {
	return "Webhook node sends a wildcard Access-Control-Allow-Origin: * CORS header."
}
func (r *RuleSEC008) Remediation() string {
	return "Restrict the allowed origins to explicit trusted domain names."
}
func (r *RuleSEC008) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if !strings.Contains(strings.ToLower(node.Type), "webhook") {
			continue
		}

		options, ok := node.Parameters["options"].(map[string]interface{})
		if !ok {
			continue
		}

		// Inspect custom headers
		headers, ok := options["responseHeaders"].(map[string]interface{})
		if !ok {
			continue
		}

		entries, ok := headers["entries"].([]interface{})
		if !ok {
			continue
		}

		for _, item := range entries {
			entryMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			name, _ := entryMap["name"].(string)
			value, _ := entryMap["value"].(string)

			if strings.EqualFold(name, "Access-Control-Allow-Origin") && strings.TrimSpace(value) == "*" {
				line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
				issues = append(issues, models.Issue{
					RuleID:      r.ID(),
					RuleName:    r.Name(),
					Category:    r.Category(),
					Severity:    r.Severity(),
					Message:     fmt.Sprintf("Webhook node '%s' responds with wildcard CORS origin (*)", node.Name),
					FilePath:    ctx.Workflow.FilePath,
					Line:        line,
					NodeID:      node.ID,
					NodeName:    node.Name,
					Field:       "parameters.options.responseHeaders",
					Remediation: r.Remediation(),
				})
			}
		}
	}

	return issues
}

// -------------------------------------------------------------
// SEC-009: unparameterized-sql-query
// -------------------------------------------------------------
type RuleSEC009 struct{}

func (r *RuleSEC009) ID() string             { return "SEC-009" }
func (r *RuleSEC009) Name() string           { return "unparameterized-sql-query" }
func (r *RuleSEC009) Category() models.Category { return models.CategorySecurity }
func (r *RuleSEC009) Severity() models.Severity { return models.SeverityHigh }
func (r *RuleSEC009) Description() string {
	return "Database node interpolates raw expressions inside SQL strings, risking SQL injection."
}
func (r *RuleSEC009) Remediation() string {
	return "Use parameterized queries ($1, $2 or :param) and n8n query parameters instead of raw string interpolation."
}
func (r *RuleSEC009) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	dbKeywords := []string{"postgres", "mysql", "mariadb", "snowflake", "oracle", "mssql", "sqlite"}

	for _, node := range ctx.Workflow.Nodes {
		lowerType := strings.ToLower(node.Type)
		isDB := false
		for _, kw := range dbKeywords {
			if strings.Contains(lowerType, kw) {
				isDB = true
				break
			}
		}
		if !isDB {
			continue
		}

		query, ok := node.Parameters["query"].(string)
		if !ok {
			continue
		}

		if strings.Contains(query, "{{") && strings.Contains(query, "$json") {
			line := ctx.LineMapper.FindFieldLine(node.Name, "query")
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Database node '%s' uses direct expression interpolation in raw SQL query", node.Name),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "parameters.query",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}

// -------------------------------------------------------------
// SEC-010: weak-cryptographic-hash
// -------------------------------------------------------------
type RuleSEC010 struct{}

func (r *RuleSEC010) ID() string             { return "SEC-010" }
func (r *RuleSEC010) Name() string           { return "weak-cryptographic-hash" }
func (r *RuleSEC010) Category() models.Category { return models.CategorySecurity }
func (r *RuleSEC010) Severity() models.Severity { return models.SeverityMedium }
func (r *RuleSEC010) Description() string {
	return "Crypto node employs deprecated hash algorithms such as MD5 or SHA1."
}
func (r *RuleSEC010) Remediation() string {
	return "Switch to modern cryptographic hash functions such as SHA-256 or SHA-512."
}
func (r *RuleSEC010) Evaluate(ctx *RuleContext) []models.Issue {
	var issues []models.Issue

	for _, node := range ctx.Workflow.Nodes {
		if !strings.Contains(strings.ToLower(node.Type), "crypto") {
			continue
		}

		algorithm, ok := node.Parameters["type"].(string)
		if !ok {
			algorithm, _ = node.Parameters["algorithm"].(string)
		}

		lower := strings.ToLower(algorithm)
		if lower == "md5" || lower == "sha1" {
			line := ctx.LineMapper.FindNodeLine(node.Name, node.ID)
			issues = append(issues, models.Issue{
				RuleID:      r.ID(),
				RuleName:    r.Name(),
				Category:    r.Category(),
				Severity:    r.Severity(),
				Message:     fmt.Sprintf("Crypto node '%s' uses weak or collision-prone hashing algorithm '%s'", node.Name, algorithm),
				FilePath:    ctx.Workflow.FilePath,
				Line:        line,
				NodeID:      node.ID,
				NodeName:    node.Name,
				Field:       "parameters.algorithm",
				Remediation: r.Remediation(),
			})
		}
	}

	return issues
}
