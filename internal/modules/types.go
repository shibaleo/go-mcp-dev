package modules

// Tool represents an MCP tool definition
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

// InputSchema defines the input parameters for a tool
type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

// Property defines a single property in the input schema
type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// ModuleDefinition defines a module with its tools and handlers
type ModuleDefinition struct {
	Name        string
	Description string
	Tools       []Tool
	Handlers    map[string]ToolHandler
}

// ToolHandler executes a tool with given parameters
type ToolHandler func(params map[string]interface{}) (string, error)

// ToolCallResult represents the result of a tool call
type ToolCallResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}

// ContentBlock represents a content block in the result
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
