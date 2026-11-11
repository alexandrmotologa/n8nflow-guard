package rules

import (
	"regexp"
	"strings"
)

var ignoreDirectiveRegex = regexp.MustCompile(`(?i)n8nflow-guard-(?:ignore|disable):\s*([a-zA-Z0-9_,\s-]+)`)

// ParseNodeIgnoredRules extracts rule IDs flagged for suppression in a node's notes field.
// Supported syntax:
//
//	// n8nflow-guard-ignore: SEC-002
//	n8nflow-guard-ignore: SEC-002, BEST-002 (public landing page endpoint)
//	/* n8nflow-guard-disable: SEC-001 */
func ParseNodeIgnoredRules(notes string) map[string]bool {
	ignored := make(map[string]bool)
	if strings.TrimSpace(notes) == "" {
		return ignored
	}

	matches := ignoreDirectiveRegex.FindAllStringSubmatch(notes, -1)
	for _, m := range matches {
		if len(m) > 1 {
			rawList := m[1]
			// Strip optional comment in parentheses
			if idx := strings.Index(rawList, "("); idx != -1 {
				rawList = rawList[:idx]
			}
			parts := strings.Split(rawList, ",")
			for _, p := range parts {
				ruleID := strings.ToUpper(strings.TrimSpace(p))
				if ruleID != "" {
					ignored[ruleID] = true
				}
			}
		}
	}

	return ignored
}
