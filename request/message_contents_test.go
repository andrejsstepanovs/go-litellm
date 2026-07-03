package request_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andrejsstepanovs/go-litellm/request"
)

func TestMessageContent_CacheMarshal(t *testing.T) {
	tests := []struct {
		name         string
		content      request.MessageContent
		expectedJSON string
	}{
		{
			name: "ephemeral cache control",
			content: request.MessageContent{
				Type: "text",
				Text: "cache this",
			}.Cache(request.CacheControlEphemeral),
			expectedJSON: `{"type":"text","text":"cache this","cache_control":{"type":"ephemeral"}}`,
		},
		{
			name: "ephemeral cache control with ttl",
			content: request.MessageContent{
				Type: "text",
				Text: "cache this for one hour",
			}.Cache(request.CacheControlEphemeral, request.CacheTTL("1h")),
			expectedJSON: `{"type":"text","text":"cache this for one hour","cache_control":{"type":"ephemeral","ttl":"1h"}}`,
		},
		{
			name: "arbitrary cache control type with ttl",
			content: request.MessageContent{
				Type: "text",
				Text: "cache with future provider type",
			}.Cache(request.CacheControlType("provider-specific"), request.CacheTTL("30m")),
			expectedJSON: `{"type":"text","text":"cache with future provider type","cache_control":{"type":"provider-specific","ttl":"30m"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.content)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expectedJSON, string(data))
		})
	}
}

func TestMessageContent_Cache(t *testing.T) {
	tests := []struct {
		name          string
		content       request.MessageContent
		controlType   request.CacheControlType
		options       []request.CacheOption
		expected      request.MessageContent
		originalCache *request.CacheControl
	}{
		{
			name: "adds default ephemeral cache control",
			content: request.MessageContent{
				Type: "text",
				Text: "cache this",
			},
			controlType: request.CacheControlEphemeral,
			expected: request.MessageContent{
				Type: "text",
				Text: "cache this",
				CacheControl: &request.CacheControl{
					Type: request.CacheControlEphemeral,
				},
			},
		},
		{
			name: "adds ttl option",
			content: request.MessageContent{
				Type: "text",
				Text: "cache this for one hour",
			},
			controlType: request.CacheControlEphemeral,
			options:     []request.CacheOption{request.CacheTTL("1h")},
			expected: request.MessageContent{
				Type: "text",
				Text: "cache this for one hour",
				CacheControl: &request.CacheControl{
					Type: request.CacheControlEphemeral,
					TTL:  "1h",
				},
			},
		},
		{
			name: "supports custom cache control type",
			content: request.MessageContent{
				Type: "text",
				Text: "cache with provider-specific type",
			},
			controlType: request.CacheControlType("provider-specific"),
			expected: request.MessageContent{
				Type: "text",
				Text: "cache with provider-specific type",
				CacheControl: &request.CacheControl{
					Type: request.CacheControlType("provider-specific"),
				},
			},
		},
		{
			name: "replaces existing cache control",
			content: request.MessageContent{
				Type: "text",
				Text: "replace cache control",
				CacheControl: &request.CacheControl{
					Type: request.CacheControlType("old-type"),
					TTL:  "5m",
				},
			},
			controlType: request.CacheControlEphemeral,
			expected: request.MessageContent{
				Type: "text",
				Text: "replace cache control",
				CacheControl: &request.CacheControl{
					Type: request.CacheControlEphemeral,
				},
			},
			originalCache: &request.CacheControl{
				Type: request.CacheControlType("old-type"),
				TTL:  "5m",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.content.Cache(tt.controlType, tt.options...)

			assert.Equal(t, tt.expected, actual)
			if tt.originalCache != nil {
				require.NotNil(t, tt.content.CacheControl)
				assert.Equal(t, *tt.originalCache, *tt.content.CacheControl)
			}
		})
	}
}

func TestMessageContent_CacheDoesNotMutateOriginal(t *testing.T) {
	content := request.MessageContent{
		Type: "text",
		Text: "cache this",
	}

	cached := content.Cache(request.CacheControlEphemeral)

	assert.Nil(t, content.CacheControl)
	require.NotNil(t, cached.CacheControl)
	assert.Equal(t, request.CacheControl{Type: request.CacheControlEphemeral}, *cached.CacheControl)
}

func TestMessageContents_String(t *testing.T) {
	tests := []struct {
		name     string
		contents request.MessageContents
		expected string
	}{
		{
			name: "empty",
		},
		{
			name: "text content",
			contents: request.MessageContents{
				{
					Type: "text",
					Text: "hello",
				},
				{
					Type: "text",
					Text: "world",
				},
			},
			expected: "hello world",
		},
		{
			name: "image content",
			contents: request.MessageContents{
				{
					Type: "image_url",
					ImageUrl: &request.ImageUrl{
						URL: "https://example.com/image.png",
					},
				},
			},
			expected: "[Image: https://example.com/image.png]",
		},
		{
			name: "mixed supported content",
			contents: request.MessageContents{
				{
					Type: "text",
					Text: "describe",
				},
				{
					Type: "image_url",
					ImageUrl: &request.ImageUrl{
						URL: "https://example.com/image.png",
					},
				},
			},
			expected: "describe [Image: https://example.com/image.png]",
		},
		{
			name: "unsupported and nil image content skipped",
			contents: request.MessageContents{
				{
					Type: "unsupported",
					Text: "hidden",
				},
				{
					Type: "image_url",
				},
				{
					Type: "text",
					Text: "visible",
				},
			},
			expected: "visible",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.contents.String())
		})
	}
}
