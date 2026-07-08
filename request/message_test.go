package request_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andrejsstepanovs/go-litellm/common"
	"github.com/andrejsstepanovs/go-litellm/request"
	"github.com/andrejsstepanovs/go-litellm/response"
)

func TestMessage_CachePoint(t *testing.T) {
	message := request.SystemMessageSimple("cache this").CachePoint()

	require.Len(t, message.Contents, 1)
	require.NotNil(t, message.Contents[0].CacheControl)
	assert.Equal(t, request.CacheControl{
		Type: request.CacheControlEphemeral,
	}, *message.Contents[0].CacheControl)
}

func TestMessage_CachePointMutatesMessage(t *testing.T) {
	message := request.SystemMessageSimple("cache this")

	message.CachePoint()

	require.NotNil(t, message.Contents[0].CacheControl)
	assert.Equal(t, request.CacheControl{
		Type: request.CacheControlEphemeral,
	}, *message.Contents[0].CacheControl)
}

func TestMessage_CachePointWithNoContents(t *testing.T) {
	message := request.Message{
		Role: request.ROLE_SYSTEM,
	}
	cached := message.CachePoint()

	assert.Empty(t, message.Contents)
	assert.Equal(t, message, cached)
}

func TestMessage_LastContent(t *testing.T) {
	message := request.SystemMessage(request.MessageContents{
		{
			Type: "text",
			Text: "first",
		},
		{
			Type: "text",
			Text: "last",
		},
	})

	content := message.LastContent()

	require.NotNil(t, content)
	assert.Equal(t, request.MessageContent{
		Type: "text",
		Text: "last",
	}, *content)
}

func TestMessage_LastContentReturnsNilForEmptyMessage(t *testing.T) {
	message := request.Message{
		Role: request.ROLE_SYSTEM,
	}

	assert.Nil(t, message.LastContent())
}

func TestMessage_LastContentCanModifyContent(t *testing.T) {
	message := request.SystemMessageSimple("cache this")
	content := message.LastContent()
	require.NotNil(t, content)

	*content = content.Cache(request.CacheControlEphemeral)

	require.NotNil(t, message.Contents[0].CacheControl)
	assert.Equal(t, request.CacheControl{
		Type: request.CacheControlEphemeral,
	}, *message.Contents[0].CacheControl)
}

func TestMessageConstructors(t *testing.T) {
	contents := request.MessageContents{{Type: "text", Text: "hello"}}

	assert.Equal(t, request.Message{
		Role:     request.ROLE_SYSTEM,
		Contents: request.MessageContents{{Type: "text", Text: "system"}},
	}, request.SystemMessageSimple("system"))
	assert.Equal(t, request.Message{
		Role:     request.ROLE_USER,
		Contents: request.MessageContents{{Type: "text", Text: "user"}},
	}, request.UserMessageSimple("user"))
	assert.Equal(t, request.Message{
		Role:     request.ROLE_SYSTEM,
		Contents: contents,
	}, request.SystemMessage(contents))
	assert.Equal(t, request.Message{
		Role:     request.ROLE_USER,
		Contents: contents,
	}, request.UserMessage(contents))
	assert.Equal(t, request.Message{
		Role:     request.ROLE_ASSISTANT,
		Contents: request.MessageContents{{Type: "text", Text: "assistant"}},
	}, request.AssistantMessageSimple("assistant"))
	assert.Equal(t, request.Message{
		Role:     request.ROLE_ASSISTANT,
		Contents: request.MessageContents{{Type: "text", Text: "-"}},
	}, request.AssistantMessageSimple(""))
	assert.Equal(t, request.Message{
		Role:     request.ROLE_ASSISTANT,
		Contents: contents,
	}, request.AssistantMessage(contents))
}

func TestMessageImageHelpers(t *testing.T) {
	image := request.MessageImage("https://example.com/image.png")

	assert.Equal(t, request.ImageUrl{URL: "https://example.com/image.png"}, image)
	assert.Equal(t, request.Message{
		Role: request.ROLE_USER,
		Contents: request.MessageContents{
			{
				Type: "text",
				Text: "describe",
			},
			{
				Type:     "image_url",
				ImageUrl: &image,
			},
		},
	}, request.UserMessageImage("describe", image))
	assert.Equal(t, request.Message{
		Role: request.ROLE_USER,
		Contents: request.MessageContents{
			{
				Type:     "image_url",
				ImageUrl: &image,
			},
		},
	}, request.UserMessageImage("", image))
}

func TestToolCallMessage(t *testing.T) {
	toolCall := common.ToolCall{
		ID: "call-1",
		Function: common.ToolCallFunction{
			Name: "lookup",
		},
	}

	assert.Equal(t, request.Message{
		Role:       request.ROLE_TOOL,
		Name:       "lookup",
		ToolCallID: "call-1",
		Contents: request.MessageContents{
			{
				Type: "text",
				Text: "tool result",
			},
		},
	}, request.ToolCallMessage(toolCall, response.ToolResponse{Text: "tool result"}))
	assert.Equal(t, "-", request.ToolCallMessage(toolCall, response.ToolResponse{}).Contents[0].Text)
}

func TestAIMessage(t *testing.T) {
	toolCalls := common.ToolCalls{
		{ID: "second", Index: 2},
		{ID: "first", Index: 1},
	}

	message := request.AIMessage(response.ResponseMessage{
		Role:      "assistant",
		Content:   "done",
		ToolCalls: toolCalls,
	})

	assert.Equal(t, request.ROLE_ASSISTANT, message.Role)
	assert.Equal(t, "done", message.Contents[0].Text)
	assert.Equal(t, "first", message.ToolCalls[0].ID)
	assert.Equal(t, "second", message.ToolCalls[1].ID)
	assert.Equal(t, "-", request.AIMessage(response.ResponseMessage{Role: "assistant"}).Contents[0].Text)
}

// TestAIMessage_PreservesThoughtSignature ensures that when a response's
// tool call carries a Gemini thought_signature in provider_specific_fields,
// AIMessage() round-trips it into the assistant request message unchanged
// so it can be echoed back to Gemini on the next turn. Without this,
// multi-turn function calling with Gemini fails with a 400 "missing
// thought_signature" error.
func TestAIMessage_PreservesThoughtSignature(t *testing.T) {
	toolCalls := common.ToolCalls{
		{
			ID:                     "call-1",
			Index:                  0,
			Function:               common.ToolCallFunction{Name: "lookup"},
			ProviderSpecificFields: map[string]any{"thought_signature": "Cs0CAXLI2nyD"},
		},
	}

	message := request.AIMessage(response.ResponseMessage{
		Role:      "assistant",
		ToolCalls: toolCalls,
	})

	assert.Equal(t, "Cs0CAXLI2nyD", message.ToolCalls[0].ThoughtSignature())
}

func TestMessagesHelpers(t *testing.T) {
	user := request.UserMessageSimple("hello")
	assistant := request.AssistantMessageSimple("hi")
	messages := request.Messages{}

	messages.AddMessage(user)
	messages.AddMessagePair(user, assistant)

	assert.Equal(t, request.Messages{user, user, assistant}, messages)
	assert.Equal(t, "user: hello\nuser: hello\nassistant: hi", messages.String())
}

func TestMessageString(t *testing.T) {
	message := request.SystemMessageSimple("system instructions")

	assert.Equal(t, "system: system instructions", message.String())
}
