package common

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArguments(t *testing.T) {
	t.Run("string types are preserved", func(t *testing.T) {
		const file = "testdata/args.json"
		dat, err := os.ReadFile(file)
		assert.NoError(t, err)

		toolCall := ToolCalls{}
		err = json.Unmarshal(dat, &toolCall)
		assert.NoError(t, err)

		assert.Len(t, toolCall, 1)
		args := toolCall[0].Function.Arguments
		val, exists := args.GetArgument("country")
		assert.True(t, exists)
		assert.Equal(t, "Latvia", val)
	})

	t.Run("int argument is preserved as string", func(t *testing.T) {
		const file = "testdata/args_int.json"
		dat, err := os.ReadFile(file)
		assert.NoError(t, err)

		toolCall := ToolCalls{}
		err = json.Unmarshal(dat, &toolCall)
		assert.NoError(t, err)

		assert.Len(t, toolCall, 1)
		args := toolCall[0].Function.Arguments
		val, exists := args.GetStrArgument("pageIdx")
		assert.True(t, exists)
		assert.Equal(t, "4", val)
	})

	t.Run("bool argument is preserved as bool", func(t *testing.T) {
		const file = "testdata/args_bool.json"
		dat, err := os.ReadFile(file)
		assert.NoError(t, err)

		toolCall := ToolCalls{}
		err = json.Unmarshal(dat, &toolCall)
		assert.NoError(t, err)

		assert.Len(t, toolCall, 1)
		args := toolCall[0].Function.Arguments
		val, exists := args.GetArgument("isTrue")
		assert.True(t, exists)
		assert.Equal(t, true, val)
	})
}

func TestToolCall_ThoughtSignature(t *testing.T) {
	t.Run("unmarshals and exposes thought_signature from provider_specific_fields", func(t *testing.T) {
		jsonStr := `{
			"id": "call_abc123",
			"type": "function",
			"index": 0,
			"function": {"name": "get_weather", "arguments": "{\"location\": \"Tokyo\"}"},
			"provider_specific_fields": {"thought_signature": "Cs0CAXLI2nyD"}
		}`

		var tc ToolCall
		err := json.Unmarshal([]byte(jsonStr), &tc)
		assert.NoError(t, err)
		assert.Equal(t, "Cs0CAXLI2nyD", tc.ThoughtSignature())
	})

	t.Run("round-trips thought_signature back into the marshaled JSON", func(t *testing.T) {
		tc := ToolCall{
			ID:   "call_abc123",
			Type: "function",
			Function: ToolCallFunction{
				Name:      "get_weather",
				Arguments: Arguments{"location": "Tokyo"},
			},
			ProviderSpecificFields: map[string]any{"thought_signature": "Cs0CAXLI2nyD"},
		}

		out, err := json.Marshal(&tc)
		assert.NoError(t, err)

		var roundTripped ToolCall
		err = json.Unmarshal(out, &roundTripped)
		assert.NoError(t, err)
		assert.Equal(t, "Cs0CAXLI2nyD", roundTripped.ThoughtSignature())
	})

	t.Run("returns empty string when absent", func(t *testing.T) {
		tc := ToolCall{ID: "call_abc123"}
		assert.Equal(t, "", tc.ThoughtSignature())
	})

	t.Run("returns empty string when value is not a string", func(t *testing.T) {
		tc := ToolCall{ProviderSpecificFields: map[string]any{"thought_signature": 42}}
		assert.Equal(t, "", tc.ThoughtSignature())
	})
}
