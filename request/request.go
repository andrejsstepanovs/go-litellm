package request

import (
	"github.com/andrejsstepanovs/go-litellm/pkg/json_schema"
	"github.com/andrejsstepanovs/go-litellm/pkg/models"
)

// https://docs.litellm.ai/docs/completion/input
type Request struct {
	Model          models.ModelID  `json:"model"`
	Messages       Messages        `json:"messages"`
	Stream         bool            `json:"stream"`
	Temperature    float32         `json:"temperature,omitempty"`
	Tools          *LLMCallTools   `json:"tools,omitempty"`
	ToolChoice     string          `json:"tool_choice,omitempty"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
	Functions      string          `json:"functions,omitempty"`
}

// TokenCounterRequest represents the request body for the LiteLLM /utils/token_counter endpoint.
type TokenCounterRequest struct {
	Model    models.ModelID `json:"model"`
	Messages Messages       `json:"messages,omitempty"`
}

func NewRequest(model models.ModelMeta) *Request {
	r := &Request{
		Model:  model.ModelId,
		Stream: false,
	}

	return r
}

func (r *Request) SetMessages(messages Messages) *Request {
	messages.RemoveEmpty()
	r.Messages = messages
	return r
}

// SetTemperature sets the temperature for the request.
// Use value -1 to unset the temperature.
func (r *Request) SetTemperature(temp float32, supportedParams []string) *Request {
	if temp < 0 {
		return r
	}
	for _, param := range supportedParams {
		if param == "temperature" {
			r.Temperature = temp
			break
		}
	}
	return r
}

func (r *Request) SetAvailableTools(tools LLMCallTools) *Request {
	r.Tools = &tools

	return r
}

// SetJSONSchema sets the response format to use JSON schema for structured output
func (r *Request) SetJSONSchema(schema json_schema.JSONSchema) *Request {
	r.ResponseFormat = &ResponseFormat{
		Type:       "json_schema",
		JSONSchema: schema,
	}
	return r
}

func (r *Request) SetJSONMode() *Request {
	r.ResponseFormat = &ResponseFormat{
		Type: "json_object",
	}
	return r
}

func NewCompletionRequest(model models.ModelMeta, messages Messages, availableTools LLMCallTools, temperature *float32, defaultTemperature float32) *Request {
	r := NewRequest(model)
	r.SetMessages(messages)

	if len(availableTools) > 0 {
		r.SetAvailableTools(availableTools)
	}

	if temperature != nil {
		r.SetTemperature(*temperature, model.SupportedOpenAIParams)
	} else {
		r.SetTemperature(defaultTemperature, model.SupportedOpenAIParams)
	}

	return r
}
