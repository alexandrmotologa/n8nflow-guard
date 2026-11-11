package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// ExpressionMatch holds a extracted expression and its position in the raw string.
type ExpressionMatch struct {
	Raw        string // The full {{ ... }} or ={{ ... }} string
	Inner      string // The extracted JavaScript content inside {{ ... }}
	StartIndex int
	EndIndex   int
}

var (
	// Regex matching standard n8n expression blocks {{ ... }}
	expressionRegex = regexp.MustCompile(`\{\{([^{}]*(?:\{[^{}]*\}[^{}]*)*)\}\}`)
	// Regex matching node reference syntax $('NodeName') or $("NodeName")
	nodeRefRegex = regexp.MustCompile(`\$\(['"]([^'"]+)['"]\)`)
	// Regex matching legacy node reference syntax $node["NodeName"]
	legacyNodeRefRegex = regexp.MustCompile(`\$node\[['"]([^'"]+)['"]\]`)
	// Regex matching environment variable reference $env.VARIABLE or $env['VARIABLE']
	envRefRegex = regexp.MustCompile(`\$env(?:\.([A-Za-z0-9_]+)|\[['"]([A-Za-z0-9_]+)['"]\])`)
)

// ExtractExpressions finds all n8n expression segments within a parameter string.
func ExtractExpressions(input string) []ExpressionMatch {
	var matches []ExpressionMatch

	// Handle leading equal sign (=) shorthand
	trimmed := strings.TrimSpace(input)
	if strings.HasPrefix(trimmed, "=") {
		trimmed = strings.TrimPrefix(trimmed, "=")
	}

	submatches := expressionRegex.FindAllStringSubmatchIndex(trimmed, -1)
	for _, sm := range submatches {
		if len(sm) >= 4 {
			raw := trimmed[sm[0]:sm[1]]
			inner := trimmed[sm[2]:sm[3]]
			matches = append(matches, ExpressionMatch{
				Raw:        raw,
				Inner:      strings.TrimSpace(inner),
				StartIndex: sm[0],
				EndIndex:   sm[1],
			})
		}
	}

	return matches
}

// ValidateExpressionSyntax checks for common expression syntax errors like unclosed braces.
func ValidateExpressionSyntax(input string) error {
	// Check for unclosed {{ without matching }}
	openCount := strings.Count(input, "{{")
	closeCount := strings.Count(input, "}}")

	if openCount != closeCount {
		return fmt.Errorf("mismatched expression braces: %d opening '{{' and %d closing '}}'", openCount, closeCount)
	}

	for _, expr := range ExtractExpressions(input) {
		// Basic check for balanced parentheses inside expression
		parens := 0
		for _, r := range expr.Inner {
			if r == '(' {
				parens++
			} else if r == ')' {
				parens--
				if parens < 0 {
					return fmt.Errorf("unexpected closing parenthesis in expression: %s", expr.Raw)
				}
			}
		}
		if parens != 0 {
			return fmt.Errorf("unclosed parenthesis in expression: %s", expr.Raw)
		}
	}

	return nil
}

// ExtractReferencedNodes returns the names of all other nodes referenced within an expression.
func ExtractReferencedNodes(expr string) []string {
	var nodes []string

	for _, submatch := range nodeRefRegex.FindAllStringSubmatch(expr, -1) {
		if len(submatch) > 1 {
			nodes = append(nodes, submatch[1])
		}
	}

	for _, submatch := range legacyNodeRefRegex.FindAllStringSubmatch(expr, -1) {
		if len(submatch) > 1 {
			nodes = append(nodes, submatch[1])
		}
	}

	return nodes
}

// IsEnvExpression returns true if the input consists exclusively of an $env reference.
func IsEnvExpression(input string) bool {
	trimmed := strings.TrimSpace(input)
	if strings.HasPrefix(trimmed, "=") {
		trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "="))
	}
	if strings.HasPrefix(trimmed, "{{") && strings.HasSuffix(trimmed, "}}") {
		inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "{{"), "}}"))
		return envRefRegex.MatchString(inner)
	}
	return false
}
