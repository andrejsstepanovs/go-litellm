package request_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/json_schema"
	"github.com/andrejsstepanovs/go-litellm/models"
	"github.com/andrejsstepanovs/go-litellm/request"
)

func TestResponseFormat_Marshal(t *testing.T) {
	tests := []struct {
		name           string
		responseFormat *request.ResponseFormat
		expectedJSON   string
	}{
		{
			name: "json_schema with strict mode",
			responseFormat: &request.ResponseFormat{
				Type: "json_schema",
				JSONSchema: json_schema.JSONSchema{
					Schema: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name": map[string]interface{}{
								"type": "string",
							},
						},
						"required": []string{"name"},
					},
					Strict: true,
				},
			},
			expectedJSON: `{"type":"json_schema","json_schema":{"name":"","schema":{"type":"object","properties":{"name":{"type":"string"}},"required":["name"]},"strict":true}}`,
		},
		{
			name: "json_schema without strict mode",
			responseFormat: &request.ResponseFormat{
				Type: "json_schema",
				JSONSchema: json_schema.JSONSchema{
					Schema: map[string]interface{}{
						"type": "string",
					},
					Strict: false,
				},
			},
			expectedJSON: `{"type":"json_schema","json_schema":{"name":"","schema":{"type":"string"},"strict":false}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.responseFormat)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expectedJSON, string(jsonBytes))
		})
	}
}

func TestRequest_WithResponseFormat_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		request      *request.Request
		expectedJSON string
	}{
		{
			name: "request with response_format",
			request: &request.Request{
				Model: models.ModelID("test-model"),
				Messages: request.Messages{
					{
						Role: request.ROLE_USER,
						Contents: request.MessageContents{
							{
								Type: "text",
								Text: "test message",
							},
						},
					},
				},
				Stream: false,
				ResponseFormat: &request.ResponseFormat{
					Type: "json_schema",
					JSONSchema: json_schema.JSONSchema{
						Schema: map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"answer": map[string]interface{}{
									"type": "string",
								},
							},
						},
						Strict: true,
					},
				},
			},
			expectedJSON: `{"model":"test-model","messages":[{"role":"user","content":[{"type":"text","text":"test message"}]}],"stream":false,"response_format":{"type":"json_schema","json_schema":{"name":"","schema":{"properties":{"answer":{"type":"string"}},"type":"object"},"strict":true}}}`,
		},
		{
			name: "request without response_format",
			request: &request.Request{
				Model: models.ModelID("test-model"),
				Messages: request.Messages{
					{
						Role: request.ROLE_USER,
						Contents: request.MessageContents{
							{
								Type: "text",
								Text: "test message",
							},
						},
					},
				},
				Stream: false,
			},
			expectedJSON: `{"model":"test-model","messages":[{"role":"user","content":[{"type":"text","text":"test message"}]}],"stream":false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.request)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expectedJSON, string(jsonBytes))
		})
	}
}
