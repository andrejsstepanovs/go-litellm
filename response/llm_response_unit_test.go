package response_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/response"
)

func Test_Choice_Unit(t *testing.T) {
	tests := []struct {
		name     string
		response *response.Response
		expected response.ResponseChoice
	}{
		{
			name: "Non-empty Choices (single element)",
			response: &response.Response{
				Choices: response.ResponseChoices{
					{
						FinishReason: response.FINISH_REASON_STOP,
						Index:        0,
						Message:      response.ResponseMessage{Content: "Hello", Role: "assistant"},
					},
				},
			},
			expected: response.ResponseChoice{
				FinishReason: response.FINISH_REASON_STOP,
				Index:        0,
				Message:      response.ResponseMessage{Content: "Hello", Role: "assistant"},
			},
		},
		{
			name: "Non-empty Choices (multiple elements)",
			response: &response.Response{
				Choices: response.ResponseChoices{
					{
						FinishReason: response.FINISH_REASON_TOOL,
						Index:        0,
						Message:      response.ResponseMessage{Content: "First", Role: "assistant"},
					},
					{
						FinishReason: response.FINISH_REASON_STOP,
						Index:        1,
						Message:      response.ResponseMessage{Content: "Second", Role: "assistant"},
					},
				},
			},
			expected: response.ResponseChoice{
				FinishReason: response.FINISH_REASON_STOP,
				Index:        1,
				Message:      response.ResponseMessage{Content: "Second", Role: "assistant"},
			},
		},
		{
			name:     "Empty Choices",
			response: &response.Response{Choices: response.ResponseChoices{}},
			expected: response.ResponseChoice{},
		},
		{
			name:     "Nil Choices slice",
			response: &response.Response{Choices: nil},
			expected: response.ResponseChoice{},
		},
		{
			name:     "Nil Response pointer",
			response: nil,
			expected: response.ResponseChoice{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.response.Choice()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_String_Unit(t *testing.T) {
	tests := []struct {
		name     string
		response *response.Response
		expected string
	}{
		{
			name: "Valid Response with non-empty content",
			response: &response.Response{
				Choices: response.ResponseChoices{
					{
						Message: response.ResponseMessage{Content: "Hello, world!", Role: "assistant"},
					},
				},
			},
			expected: "Hello, world!",
		},
		{
			name: "Valid Response with empty content",
			response: &response.Response{
				Choices: response.ResponseChoices{
					{
						Message: response.ResponseMessage{Content: "", Role: "assistant"},
					},
				},
			},
			expected: "",
		},
		{
			name: "Valid Response with zero-value Message",
			response: &response.Response{
				Choices: response.ResponseChoices{
					{
						Message: response.ResponseMessage{},
					},
				},
			},
			expected: "",
		},
		{
			name:     "Empty Choices slice",
			response: &response.Response{Choices: response.ResponseChoices{}},
			expected: "",
		},
		{
			name:     "Nil Choices slice",
			response: &response.Response{Choices: nil},
			expected: "",
		},
		{
			name:     "Nil Response pointer",
			response: nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.response.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_Message_Unit(t *testing.T) {
	tests := []struct {
		name     string
		response *response.Response
		expected response.ResponseMessage
	}{
		{
			name: "Non-empty Choices (single element)",
			response: &response.Response{
				Choices: response.ResponseChoices{
					{
						FinishReason: response.FINISH_REASON_STOP,
						Index:        0,
						Message:      response.ResponseMessage{Content: "Hello", Role: "assistant"},
					},
				},
			},
			expected: response.ResponseMessage{Content: "Hello", Role: "assistant"},
		},
		{
			name: "Non-empty Choices (multiple elements)",
			response: &response.Response{
				Choices: response.ResponseChoices{
					{
						FinishReason: response.FINISH_REASON_STOP,
						Index:        0,
						Message:      response.ResponseMessage{Content: "First", Role: "assistant"},
					},
					{
						FinishReason: response.FINISH_REASON_TOOL,
						Index:        1,
						Message:      response.ResponseMessage{Content: "Second", Role: "assistant"},
					},
				},
			},
			expected: response.ResponseMessage{Content: "Second", Role: "assistant"},
		},
		{
			name:     "Empty Choices",
			response: &response.Response{Choices: response.ResponseChoices{}},
			expected: response.ResponseMessage{},
		},
		{
			name:     "Nil Choices slice",
			response: &response.Response{Choices: nil},
			expected: response.ResponseMessage{},
		},
		{
			name:     "Nil Response pointer",
			response: nil,
			expected: response.ResponseMessage{},
		},
		{
			name: "Zero-value Message in first Choice",
			response: &response.Response{
				Choices: response.ResponseChoices{
					{
						FinishReason: response.FINISH_REASON_STOP,
						Index:        0,
						Message:      response.ResponseMessage{},
					},
				},
			},
			expected: response.ResponseMessage{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.response.Message()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func Test_SetText_Unit(t *testing.T) {
	tests := []struct {
		name          string
		response      *response.Response
		text          string
		expectedState *response.Response
	}{
		{
			name:     "Set text on empty Choices slice",
			response: &response.Response{},
			text:     "Hello world",
			expectedState: &response.Response{
				Choices: response.ResponseChoices{
					{
						FinishReason: response.FINISH_REASON_STOP,
						Message: response.ResponseMessage{
							Content: "Hello world",
							Role:    "assistant",
						},
					},
				},
			},
		},
		{
			name: "Set text on non-empty Choices slice",
			response: &response.Response{
				Choices: response.ResponseChoices{
					{
						FinishReason: response.FINISH_REASON_STOP,
						Message: response.ResponseMessage{
							Content: "Original text",
							Role:    "assistant",
						},
					},
				},
			},
			text: "Updated text",
			expectedState: &response.Response{
				Choices: response.ResponseChoices{
					{
						FinishReason: response.FINISH_REASON_STOP,
						Message: response.ResponseMessage{
							Content: "Updated text",
							Role:    "assistant",
						},
					},
				},
			},
		},
		{
			name:     "Set empty text",
			response: &response.Response{},
			text:     "",
			expectedState: &response.Response{
				Choices: response.ResponseChoices{
					{
						FinishReason: response.FINISH_REASON_STOP,
						Message: response.ResponseMessage{
							Content: "",
							Role:    "assistant",
						},
					},
				},
			},
		},
		{
			name: "Set text multiple times",
			response: &response.Response{
				Choices: response.ResponseChoices{
					{
						FinishReason: response.FINISH_REASON_STOP,
						Message: response.ResponseMessage{
							Content: "First text",
							Role:    "assistant",
						},
					},
				},
			},
			text: "Final text",
			expectedState: &response.Response{
				Choices: response.ResponseChoices{
					{
						FinishReason: response.FINISH_REASON_STOP,
						Message: response.ResponseMessage{
							Content: "Final text",
							Role:    "assistant",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.response.SetText(tt.text)
			assert.Equal(t, tt.expectedState, tt.response)
		})
	}
}

func Test_MessageReasoningString_Unit(t *testing.T) {
	tests := []struct {
		name string
		msg  response.ResponseMessage
		want string
	}{
		{
			name: "nil message returns empty string",
			msg:  response.ResponseMessage{},
			want: "",
		},
		{
			name: "no reasoning field returns empty string",
			msg:  response.ResponseMessage{Content: "hello", Role: "assistant"},
			want: "",
		},
		{
			name: "reasoning as JSON string is returned as-is",
			msg:  response.ResponseMessage{Reasoning: json.RawMessage(`"thinking then answer"`)},
			want: "thinking then answer",
		},
		{
			name: "reasoning as JSON object is re-marshalled",
			msg:  response.ResponseMessage{Reasoning: json.RawMessage(`{"text":"structured reasoning"}`)},
			want: `{"text":"structured reasoning"}`,
		},
		{
			name: "reasoning as JSON array of details is re-marshalled",
			msg: response.ResponseMessage{Reasoning: json.RawMessage(`[
				{"type":"reasoning.summary","summary":"sum"},
				{"type":"reasoning.text","text":"details"}
			]`)},
			want: "[\n\t\t\t\t{\"type\":\"reasoning.summary\",\"summary\":\"sum\"},\n\t\t\t\t{\"type\":\"reasoning.text\",\"text\":\"details\"}\n\t\t\t]",
		},
		{
			name: "empty reasoning payload returns empty string",
			msg:  response.ResponseMessage{Reasoning: json.RawMessage(`""`)},
			want: "",
		},
		{
			name: "invalid JSON bytes are returned verbatim",
			msg:  response.ResponseMessage{Reasoning: json.RawMessage(`not-json`)},
			want: "not-json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.msg.ReasoningString())
		})
	}
}

func Test_MessageReasoningString_NilReceiver(t *testing.T) {
	var msg *response.ResponseMessage
	assert.Equal(t, "", msg.ReasoningString())
}

func Test_MessageIsEmpty_Unit(t *testing.T) {
	tests := []struct {
		name string
		msg  response.ResponseMessage
		want bool
	}{
		{
			name: "all fields empty",
			msg:  response.ResponseMessage{},
			want: true,
		},
		{
			name: "content populated",
			msg:  response.ResponseMessage{Content: "hi"},
			want: false,
		},
		{
			name: "reasoning_content populated",
			msg:  response.ResponseMessage{ReasoningContent: "thinking"},
			want: false,
		},
		{
			name: "reasoning as string populated",
			msg:  response.ResponseMessage{Reasoning: json.RawMessage(`"answer"`)},
			want: false,
		},
		{
			name: "reasoning as object populated",
			msg:  response.ResponseMessage{Reasoning: json.RawMessage(`{"text":"x"}`)},
			want: false,
		},
		{
			name: "empty reasoning payload still counts as empty",
			msg:  response.ResponseMessage{Reasoning: json.RawMessage(`""`)},
			want: true,
		},
		{
			name: "all three reasoning shapes populated at once",
			msg: response.ResponseMessage{
				Content:          "c",
				ReasoningContent: "rc",
				Reasoning:        json.RawMessage(`"r"`),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.msg.IsEmpty())
		})
	}
}

func Test_MessageIsEmpty_NilReceiver(t *testing.T) {
	var msg *response.ResponseMessage
	assert.True(t, msg.IsEmpty())
}

func Test_ResponseReasoningString_Unit(t *testing.T) {
	resp := &response.Response{
		Choices: response.ResponseChoices{
			{
				Message: response.ResponseMessage{
					ReasoningContent: "from reasoning_content",
					Reasoning:        json.RawMessage(`"from reasoning field"`),
				},
			},
		},
	}
	// Response.ReasoningString preserves the existing behaviour: it returns
	// ReasoningContent only. The new OpenRouter-style field is exposed via
	// ResponseMessage.ReasoningString() for callers that want the merged
	// view.
	assert.Equal(t, "from reasoning_content", resp.ReasoningString())
}

func Test_ResponseUnmarshalCapturesReasoningField(t *testing.T) {
	raw := []byte(`{
		"id": "resp-1",
		"model": "openai/gpt-oss-20b",
		"choices": [
			{
				"index": 0,
				"finish_reason": "stop",
				"message": {
					"role": "assistant",
					"content": "",
					"reasoning": "the answer lives here"
				}
			}
		]
	}`)

	var resp response.Response
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	msg := resp.Message()
	assert.Equal(t, "", msg.Content)
	assert.Equal(t, "", msg.ReasoningContent)
	assert.Equal(t, "the answer lives here", msg.ReasoningString())
	assert.False(t, msg.IsEmpty())
}

func Test_ResponseUnmarshalCapturesStructuredReasoningField(t *testing.T) {
	raw := []byte(`{
		"id": "resp-2",
		"model": "openai/gpt-oss-20b",
		"choices": [
			{
				"index": 0,
				"finish_reason": "stop",
				"message": {
					"role": "assistant",
					"content": "",
					"reasoning": [
						{"type": "reasoning.text", "text": "step one"},
						{"type": "reasoning.text", "text": "step two"}
					]
				}
			}
		]
	}`)

	var resp response.Response
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	msg := resp.Message()
	assert.Equal(t, "", msg.Content)
	got := msg.ReasoningString()
	// When reasoning is an object/array, the helper re-emits the raw JSON so
	// callers can introspect the structured form.
	assert.Contains(t, got, "step one")
	assert.Contains(t, got, "step two")
	assert.False(t, msg.IsEmpty())
}
