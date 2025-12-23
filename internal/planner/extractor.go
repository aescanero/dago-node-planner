package planner

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

// Extractor extracts graph JSON from LLM responses.
type Extractor struct {
	logger *zap.Logger
}

// NewExtractor creates a new extractor.
func NewExtractor(logger *zap.Logger) *Extractor {
	return &Extractor{
		logger: logger,
	}
}

// Extract extracts graph JSON and reasoning from an LLM response.
func (e *Extractor) Extract(content string) (graphJSON string, reasoning string, err error) {
	e.logger.Debug("extracting graph from LLM response")

	// Try to extract reasoning
	reasoning = e.extractReasoning(content)

	// Extract JSON
	graphJSON = extractJSON(content)
	if graphJSON == "" {
		return "", "", fmt.Errorf("no JSON found in response")
	}

	// Validate it's valid JSON
	var temp interface{}
	if err := json.Unmarshal([]byte(graphJSON), &temp); err != nil {
		return "", "", fmt.Errorf("invalid JSON: %w", err)
	}

	return graphJSON, reasoning, nil
}

// extractReasoning extracts the reasoning section from the response.
func (e *Extractor) extractReasoning(content string) string {
	// Look for reasoning in various formats
	patterns := []string{
		`(?s)reasoning["\s:]+([^{]+)`,
		`(?s)Reasoning[:\s]+([^{]+)`,
		`(?s)## Reasoning\s+([^#]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(content); len(matches) > 1 {
			reasoning := strings.TrimSpace(matches[1])
			// Clean up common artifacts
			reasoning = strings.Trim(reasoning, `"'`)
			return reasoning
		}
	}

	return ""
}

// ParseGraph parses a graph JSON string into a map.
func (e *Extractor) ParseGraph(graphJSON string) (map[string]any, error) {
	var graph map[string]any
	if err := json.Unmarshal([]byte(graphJSON), &graph); err != nil {
		return nil, fmt.Errorf("failed to parse graph JSON: %w", err)
	}
	return graph, nil
}

// extractJSON extracts JSON from text content.
// It looks for JSON objects and arrays, preferring complete, well-formed structures.
func extractJSON(content string) string {
	// Try to find JSON in code blocks first
	codeBlockPatterns := []string{
		"```json\\s*\\n([^`]+)```",
		"```\\s*\\n([^`]+)```",
	}

	for _, pattern := range codeBlockPatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(content); len(matches) > 1 {
			jsonStr := strings.TrimSpace(matches[1])
			if isValidJSON(jsonStr) {
				return jsonStr
			}
		}
	}

	// Try to find raw JSON objects
	// Look for { ... } that spans multiple lines
	re := regexp.MustCompile(`(?s)\{[^}]*"nodes"[^}]*\}`)
	if matches := re.FindStringSubmatch(content); len(matches) > 0 {
		jsonStr := matches[0]
		if isValidJSON(jsonStr) {
			return jsonStr
		}
	}

	// Last resort: try to find any valid JSON object
	var depth int
	var start int
	var inString bool
	var escape bool

	for i, ch := range content {
		if escape {
			escape = false
			continue
		}

		if ch == '\\' {
			escape = true
			continue
		}

		if ch == '"' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		switch ch {
		case '{':
			if depth == 0 {
				start = i
			}
			depth++
		case '}':
			depth--
			if depth == 0 && start > 0 {
				jsonStr := content[start : i+1]
				if isValidJSON(jsonStr) {
					return jsonStr
				}
			}
		}
	}

	return ""
}

// isValidJSON checks if a string is valid JSON.
func isValidJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
