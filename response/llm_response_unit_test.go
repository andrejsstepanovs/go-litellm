package response_test

import (
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
			expected: response.ResponseChoice{
				FinishReason: response.FINISH_REASON_STOP,
				Index:        0,
				Message:      response.ResponseMessage{Content: "First", Role: "assistant"},
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
			expected: response.ResponseMessage{Content: "First", Role: "assistant"},
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
