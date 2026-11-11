package parser

import (
	"bytes"
	"fmt"
	"regexp"
)

// LineMapper translates JSON keys, node names, and fields to 1-based line numbers.
type LineMapper struct {
	lines [][]byte
}

// NewLineMapper initializes an index from raw file bytes.
func NewLineMapper(content []byte) *LineMapper {
	// Normalize CRLF to LF for consistent line indexing
	normalized := bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
	return &LineMapper{
		lines: bytes.Split(normalized, []byte("\n")),
	}
}

// FindNodeLine locates the line number where a node with the given name or id is defined.
func (m *LineMapper) FindNodeLine(nodeName, nodeID string) int {
	if len(m.lines) == 0 {
		return 1
	}

	namePattern := fmt.Sprintf(`"name"\s*:\s*"%s"`, regexp.QuoteMeta(nodeName))
	nameRe := regexp.MustCompile(namePattern)

	for idx, line := range m.lines {
		if nameRe.Match(line) {
			return idx + 1
		}
	}

	if nodeID != "" {
		idPattern := fmt.Sprintf(`"id"\s*:\s*"%s"`, regexp.QuoteMeta(nodeID))
		idRe := regexp.MustCompile(idPattern)
		for idx, line := range m.lines {
			if idRe.Match(line) {
				return idx + 1
			}
		}
	}

	return 1
}

// FindFieldLine finds the line of a specific field within or near a node's definition.
func (m *LineMapper) FindFieldLine(nodeName, fieldName string) int {
	startLine := m.FindNodeLine(nodeName, "")
	if len(m.lines) == 0 {
		return 1
	}

	fieldPattern := fmt.Sprintf(`"%s"\s*:`, regexp.QuoteMeta(fieldName))
	fieldRe := regexp.MustCompile(fieldPattern)

	// Search forward up to 80 lines from node start
	maxLookahead := startLine + 80
	if maxLookahead > len(m.lines) {
		maxLookahead = len(m.lines)
	}

	for idx := startLine - 1; idx < maxLookahead; idx++ {
		if idx >= 0 && fieldRe.Match(m.lines[idx]) {
			return idx + 1
		}
	}

	// Fallback to anywhere in the document
	for idx, line := range m.lines {
		if fieldRe.Match(line) {
			return idx + 1
		}
	}

	return startLine
}

// TotalLines returns the total count of lines indexed.
func (m *LineMapper) TotalLines() int {
	return len(m.lines)
}
