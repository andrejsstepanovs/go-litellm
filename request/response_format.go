package request

// ResponseFormat represents the response_format field in LiteLLM API requests
type ResponseFormat struct {
	Type       string     `json:"type"`
	JSONSchema JSONSchema `json:"json_schema,omitempty"`
}
