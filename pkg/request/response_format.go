package request

import "github.com/andrejsstepanovs/go-litellm/pkg/json_schema"

// ResponseFormat represents the response_format field in LiteLLM API requests
type ResponseFormat struct {
	Type       string                 `json:"type"`
	JSONSchema json_schema.JSONSchema `json:"json_schema,omitempty"`
}
