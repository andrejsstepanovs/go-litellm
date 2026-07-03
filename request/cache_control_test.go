package request_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andrejsstepanovs/go-litellm/request"
)

func TestCacheControl_Marshal(t *testing.T) {
	tests := []struct {
		name         string
		cacheControl request.CacheControl
		expectedJSON string
	}{
		{
			name: "ephemeral",
			cacheControl: request.CacheControl{
				Type: request.CacheControlEphemeral,
			},
			expectedJSON: `{"type":"ephemeral"}`,
		},
		{
			name: "custom type with ttl",
			cacheControl: request.CacheControl{
				Type: request.CacheControlType("provider-specific"),
				TTL:  "30m",
			},
			expectedJSON: `{"type":"provider-specific","ttl":"30m"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.cacheControl)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expectedJSON, string(data))
		})
	}
}

func TestMessages_CacheControlCount(t *testing.T) {
	messages := request.Messages{
		request.SystemMessage(request.MessageContents{
			request.MessageContent{
				Type: "text",
				Text: "A",
			}.Cache(request.CacheControlEphemeral),
			{
				Type: "text",
				Text: "uncached",
			},
		}),
		request.UserMessage(request.MessageContents{
			request.MessageContent{
				Type: "text",
				Text: "B",
			}.Cache(request.CacheControlEphemeral, request.CacheTTL("1h")),
			request.MessageContent{
				Type: "text",
				Text: "C",
			}.Cache(request.CacheControlType("custom-type")),
		}),
	}

	assert.Equal(t, 3, messages.CacheControlCount())
}

func TestMessages_RemoveEmptyPreservesCacheControl(t *testing.T) {
	messages := request.Messages{
		{
			Role: request.ROLE_SYSTEM,
			Contents: request.MessageContents{
				request.MessageContent{
					Type: "text",
					Text: "keep cache marker",
				}.Cache(request.CacheControlEphemeral, request.CacheTTL("1h")),
				{
					Type: "text",
					Text: "",
				},
			},
		},
	}

	messages.RemoveEmpty()

	require.Len(t, messages, 1)
	require.Len(t, messages[0].Contents, 1)
	require.NotNil(t, messages[0].Contents[0].CacheControl)
	assert.Equal(t, request.CacheControl{
		Type: request.CacheControlEphemeral,
		TTL:  "1h",
	}, *messages[0].Contents[0].CacheControl)
	assert.Equal(t, 1, messages.CacheControlCount())
}
