package planner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aescanero/dago-node-planner/pkg/models"
	"go.uber.org/zap"
)

// Prompter builds LLM prompts for graph planning.
type Prompter struct {
	promptPath         string
	systemPrompt       string
	planningTemplate   string
	errorFixingTemplate string
	logger             *zap.Logger
}

// NewPrompter creates a new prompter.
func NewPrompter(promptPath string, logger *zap.Logger) *Prompter {
	p := &Prompter{
		promptPath: promptPath,
		logger:     logger,
	}

	// Load prompts from files or use defaults
	p.loadPrompts()

	return p
}

// loadPrompts loads prompt templates from files or uses defaults.
func (p *Prompter) loadPrompts() {
	// Try to load from files
	if p.promptPath != "" {
		p.systemPrompt = p.loadPromptFile("system-prompt.txt", defaultSystemPrompt)
		p.planningTemplate = p.loadPromptFile("task-planning.txt", defaultPlanningTemplate)
		p.errorFixingTemplate = p.loadPromptFile("error-fixing.txt", defaultErrorFixingTemplate)
	} else {
		// Use defaults
		p.systemPrompt = defaultSystemPrompt
		p.planningTemplate = defaultPlanningTemplate
		p.errorFixingTemplate = defaultErrorFixingTemplate
	}
}

// loadPromptFile loads a prompt from a file, or returns the default if not found.
func (p *Prompter) loadPromptFile(filename, defaultContent string) string {
	if p.promptPath == "" {
		return defaultContent
	}

	filePath := filepath.Join(p.promptPath, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		p.logger.Debug("prompt file not found, using default",
			zap.String("file", filePath),
		)
		return defaultContent
	}

	return string(content)
}

// GetSystemPrompt returns the system prompt.
func (p *Prompter) GetSystemPrompt() string {
	return p.systemPrompt
}

// BuildPlanningPrompt builds the initial planning prompt.
func (p *Prompter) BuildPlanningPrompt(
	task string,
	analysis *models.TaskAnalysis,
	schemas map[string]string,
	constraints *models.Constraints,
) (string, error) {
	prompt := p.planningTemplate

	// Replace placeholders
	prompt = strings.ReplaceAll(prompt, "{{TASK}}", task)

	// Add analysis if available
	analysisStr := ""
	if analysis != nil {
		analysisStr = fmt.Sprintf(`
Task Analysis:
- Complexity: %s
- Requires Tools: %t
- Requires Routing: %t
- Suggested Node Types: %v
- Intent: %s
`,
			analysis.Complexity,
			analysis.RequiresTools,
			analysis.RequiresRouting,
			analysis.SuggestedNodeTypes,
			analysis.Intent,
		)
	}
	prompt = strings.ReplaceAll(prompt, "{{ANALYSIS}}", analysisStr)

	// Add constraints if available
	constraintsStr := ""
	if constraints != nil {
		constraintsStr = fmt.Sprintf(`
Constraints:
- Max Nodes: %d
- Preferred Modes: %v
- Available Tools: %v
`,
			constraints.MaxNodes,
			constraints.PreferredModes,
			constraints.AvailableTools,
		)
	}
	prompt = strings.ReplaceAll(prompt, "{{CONSTRAINTS}}", constraintsStr)

	// Add schemas (simplified for prompt)
	schemasStr := "See the graph schema documentation for the full schema specification."
	prompt = strings.ReplaceAll(prompt, "{{SCHEMAS}}", schemasStr)

	return prompt, nil
}

// BuildErrorFixingPrompt builds a prompt for fixing validation errors.
func (p *Prompter) BuildErrorFixingPrompt(
	task string,
	previousGraph string,
	validationErrors []string,
	attempt int,
) (string, error) {
	prompt := p.errorFixingTemplate

	// Replace placeholders
	prompt = strings.ReplaceAll(prompt, "{{TASK}}", task)
	prompt = strings.ReplaceAll(prompt, "{{PREVIOUS_GRAPH}}", previousGraph)
	prompt = strings.ReplaceAll(prompt, "{{ATTEMPT}}", fmt.Sprintf("%d", attempt))

	// Format validation errors
	errorsStr := ""
	for i, err := range validationErrors {
		errorsStr += fmt.Sprintf("%d. %s\n", i+1, err)
	}
	prompt = strings.ReplaceAll(prompt, "{{VALIDATION_ERRORS}}", errorsStr)

	return prompt, nil
}

// Default prompts
const defaultSystemPrompt = `You are an expert graph planning assistant for the DA Orchestrator workflow system.

Your role is to convert natural language task descriptions into valid execution graphs.

A graph consists of:
1. Executor nodes: Perform actions using LLMs, tools, or both
2. Router nodes: Make routing decisions (deterministic or LLM-based)
3. Edges: Connect nodes to define execution flow

You must respond with a valid JSON graph that conforms to the graph schema.

Always include:
- A clear "reasoning" field explaining your graph design
- Proper node IDs and edge connections
- Valid node configurations for each node type

Be concise and focus on creating minimal, effective graphs.`

const defaultPlanningTemplate = `Generate an execution graph for the following task:

{{TASK}}

{{ANALYSIS}}

{{CONSTRAINTS}}

{{SCHEMAS}}

Respond with:
1. A "reasoning" section explaining your graph design
2. A "graph" section containing the complete JSON graph

The graph must conform to the graph schema and be executable.`

const defaultErrorFixingTemplate = `The previous graph had validation errors. Please fix them.

Original Task: {{TASK}}

Previous Graph (attempt {{ATTEMPT}}):
{{PREVIOUS_GRAPH}}

Validation Errors:
{{VALIDATION_ERRORS}}

Generate a corrected graph that fixes these validation errors while maintaining the intent of the original task.

Respond with:
1. A "reasoning" section explaining your fixes
2. A "graph" section containing the corrected JSON graph`
