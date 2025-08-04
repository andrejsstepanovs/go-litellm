package response_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/common"
	"github.com/andrejsstepanovs/go-litellm/response"
)

func Test_ToolCallFunction_Unmarshal_Unit(t *testing.T) {
	jsonStr := `{
  "id": "chatcmpl-b8d7e94b-c789-415c-b69d-cf88c9a10e21",
  "created": 1748012750,
  "model": "gemini-2.0-flash",
  "object": "chat.completion",
  "system_fingerprint": null,
  "choices": [
    {
      "finish_reason": "tool_calls",
      "index": 0,
      "message": {
        "content": null,
        "role": "assistant",
        "tool_calls": [
          {
            "index": 0,
            "function": {
              "arguments": "{\"timezone\": \"Europe/Berlin\"}",
              "name": "current_time"
            },
            "id": "call_eb3401da-114b-4bfa-b51c-361d0b574a83",
            "type": "function"
          }
        ],
        "function_call": null
      }
    }
  ],
  "usage": {
    "completion_tokens": 3,
    "prompt_tokens": 34,
    "total_tokens": 37,
    "completion_tokens_details": null,
    "prompt_tokens_details": {
      "audio_tokens": null,
      "cached_tokens": null
    }
  },
  "vertex_ai_grounding_metadata": [],
  "vertex_ai_safety_results": [],
  "vertex_ai_citation_metadata": []
}`

	r := response.Response{}
	err := json.Unmarshal([]byte(jsonStr), &r)
	assert.NoError(t, err)

	assert.Equal(t, "chatcmpl-b8d7e94b-c789-415c-b69d-cf88c9a10e21", r.ID)
	assert.Equal(t, "call_eb3401da-114b-4bfa-b51c-361d0b574a83", r.Choice().Message.ToolCalls[0].ID)
	assert.Equal(t, "current_time", r.Choice().Message.ToolCalls[0].Function.Name)
	assert.Equal(t, "function", r.Choice().Message.ToolCalls[0].Type)
}

func Test_SortASC_Unit(t *testing.T) {
	tc := common.ToolCalls{
		{ID: "call_2", Index: 2},
		{ID: "call_1", Index: 1},
		{ID: "call_0", Index: 0},
		{ID: "call_3", Index: 3},
	}

	sorted := tc.SortASC()

	assert.Equal(t, "call_0", sorted[0].ID)
	assert.Equal(t, "call_1", sorted[1].ID)
	assert.Equal(t, "call_2", sorted[2].ID)
	assert.Equal(t, "call_3", sorted[3].ID)
}

func Test_Arguments_GetStrArgument_Unit(t *testing.T) {
	args := common.Arguments{
		"key1": "value1",
		"key2": "value2",
	}

	value, ok := args.GetStrArgument("key1")
	assert.Equal(t, "value1", value)
	assert.True(t, ok)

	value, ok = args.GetStrArgument("key3")
	assert.Equal(t, "", value)
	assert.False(t, ok)
}

func Test_Arguments_HasKey_Unit(t *testing.T) {
	args := common.Arguments{
		"key1": "value1",
		"key2": "value2",
	}

	assert.True(t, args.HasKey("key1"))
	assert.False(t, args.HasKey("key"))
}
