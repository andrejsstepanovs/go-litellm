package request_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/request"
)

func TestUserMessage(t *testing.T) {
	tests := []struct {
		name     string
		contents request.MessageContents
		expected request.Message
	}{
		{
			name: "creates user message with content",
			contents: request.MessageContents{
				{
					Type: "text",
					Text: "Hello, how are you?",
				},
			},
			expected: request.Message{
				Role: request.ROLE_USER,
				Contents: request.MessageContents{
					{
						Type: "text",
						Text: "Hello, how are you?",
					},
				},
			},
		},
		{
			name: "creates user message with empty content",
			contents: request.MessageContents{
				{
					Type: "text",
					Text: "",
				},
			},
			expected: request.Message{
				Role: request.ROLE_USER,
				Contents: request.MessageContents{
					{
						Type: "text",
						Text: "",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := request.UserMessage(tt.contents)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAssistantMessage(t *testing.T) {
	tests := []struct {
		name     string
		contents request.MessageContents
		expected request.Message
	}{
		{
			name: "creates assistant message with content",
			contents: request.MessageContents{
				{
					Type: "text",
					Text: "I'm doing well, thank you!",
				},
			},
			expected: request.Message{
				Role: request.ROLE_ASSISTANT,
				Contents: request.MessageContents{
					{
						Type: "text",
						Text: "I'm doing well, thank you!",
					},
				},
			},
		},
		{
			name: "creates assistant message with empty content",
			contents: request.MessageContents{
				{
					Type: "text",
					Text: "",
				},
			},
			expected: request.Message{
				Role: request.ROLE_ASSISTANT,
				Contents: request.MessageContents{
					{
						Type: "text",
						Text: "",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := request.AssistantMessage(tt.contents)
			assert.Equal(t, tt.expected, result)
		})
	}
}
