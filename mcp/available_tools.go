package mcp

// AvailableToolsResponse represents the complete response structure
type AvailableToolsResponse struct {
	Tools   AvailableTools `json:"tools"`
	Error   *string        `json:"error"`
	Message string         `json:"message"`
}

type AvailableTools []AvailableTool

type AvailableTool struct {
	Type              string `json:"type,omitempty"`                // "function" for function calls
	SearchContextSize string `json:"search_context_size,omitempty"` // special param for LiteLLM web_search_preview Options: "low", "medium" (default), "high"

	Name        string                   `json:"name,omitempty"`
	Description string                   `json:"description,omitempty"`
	InputSchema AvailableToolInputSchema `json:"inputSchema,omitempty"`
	McpInfo     AvailableToolMcpInfo     `json:"mcp_info,omitempty"`
}

type AvailableToolInputSchema struct {
	Properties map[string]AvailableToolProperty `json:"properties"`
	Required   []string                         `json:"required"`
	Type       string                           `json:"type"`
}

type AvailableToolProperty struct {
	Description string `json:"description"`
	Type        string `json:"type"`
}

type AvailableToolMcpInfo struct {
	ServerName string `json:"server_name"`
}
