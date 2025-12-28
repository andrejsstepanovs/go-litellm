package response_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/common"
	"github.com/andrejsstepanovs/go-litellm/models"
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
		"key3": float64(123),
		"key4": true,
	}

	value, ok := args.GetStrArgument("key1")
	assert.Equal(t, "value1", value)
	assert.True(t, ok)

	value, ok = args.GetStrArgument("key3")
	assert.Equal(t, "123", value)
	assert.True(t, ok)

	value, ok = args.GetStrArgument("key4")
	assert.Equal(t, "true", value)
	assert.True(t, ok)

	value, ok = args.GetStrArgument("key5")
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

func Test_ToolResponses_NewFormat_Unmarshal_Unit(t *testing.T) {
	jsonStr := `{
		"meta": null,
		"content": [
			{
				"meta": null,
				"text": "Europe/London: 2025-12-28 15:14:23",
				"type": "text",
				"annotations": null
			}
		],
		"isError": false,
		"structuredContent": null
	}`

	var res response.ToolResponses
	err := json.Unmarshal([]byte(jsonStr), &res)
	assert.NoError(t, err)

	assert.Len(t, res, 1)
	assert.Equal(t, "text", res[0].Type)
	assert.Equal(t, "Europe/London: 2025-12-28 15:14:23", res[0].Text)
	assert.Nil(t, res[0].Annotations)
}

func Test_ToolResponses_OldFormat_Unmarshal_Unit(t *testing.T) {
	jsonStr := `[
		{
			"type": "text",
			"text": "Hello world",
			"annotations": null
		}
	]`

	var res response.ToolResponses
	err := json.Unmarshal([]byte(jsonStr), &res)
	assert.NoError(t, err)

	assert.Len(t, res, 1)
	assert.Equal(t, "text", res[0].Type)
	assert.Equal(t, "Hello world", res[0].Text)
	assert.Nil(t, res[0].Annotations)
}

func Test_ToolResponses_SingleFormat_Unmarshal_Unit(t *testing.T) {
	jsonStr := `{
		"type": "text",
		"text": "Single response",
		"annotations": null
	}`

	var res response.ToolResponses
	err := json.Unmarshal([]byte(jsonStr), &res)
	assert.NoError(t, err)

	assert.Len(t, res, 1)
	assert.Equal(t, "text", res[0].Type)
	assert.Equal(t, "Single response", res[0].Text)
	assert.Nil(t, res[0].Annotations)
}

func Test_Response_NewFields_Unmarshal_Unit(t *testing.T) {
	jsonStr := `{
		"id": "J0pRaYO0PJievdIP-P6LmAE",
		"model": "gemini-2.5-flash",
		"usage": {
			"total_tokens": 179,
			"prompt_tokens": 94,
			"completion_tokens": 85,
			"prompt_tokens_details": {
				"text_tokens": 94,
				"audio_tokens": null,
				"image_tokens": null,
				"cached_tokens": null
			},
			"completion_tokens_details": {
				"text_tokens": 17,
				"audio_tokens": null,
				"image_tokens": null,
				"reasoning_tokens": 68,
				"accepted_prediction_tokens": null,
				"rejected_prediction_tokens": null
			}
		},
		"object": "chat.completion",
		"choices": [{
			"index": 0,
			"message": {
				"role": "assistant",
				"images": [],
				"content": null,
				"tool_calls": [{
					"id": "call_3a095b3222434581ae6df726c290",
					"type": "function",
					"index": 0,
					"function": {
						"name": "current_time",
						"arguments": "{\"timezone\": \"Europe/Berlin\"}"
					},
					"provider_specific_fields": {
						"thought_signature": "Cs0CAXLI2nyD/fYpQMfsG4wTQG1XQYotxJQkU7MnOkyP8c4ZCm+sY7UdWumUj1NFDfHd5Vx7UfO0WfgZwyChbLV6+YBRbZos+N98/NkU2FL99I6hN/nr0jGGv5R3T637COpoyTfrb0r6cblXexBH+xQtSWuTyjOPG5MBLvqcdvyMWZ5M9wdhITG3INT0zTo5KbUNmVWX5KIvHISxe0wMIZ2OgAQDc/Mwg91oilXlvxsDlIN96z28eoZaKHcjvB+4cmLGJx/BrirBlCoXgAjUrNgpVD15lk9oGjlLJRtx8hgO/ppSm+hDGodh/0Oto0C1E1PdqamLxib6oy4wYdBNMVzy4JZ4wVEXmlUYh4Pi64zonrcrK8OgUZ3KRpWfVpVyWFHs3pedvHKGVc8bjaC391s/Fctu6D6Jlxnzz6lnOGNhSiHFHGdXgCgX6866g62w"
					}
				}],
				"function_call": null,
				"thinking_blocks": [],
				"provider_specific_fields": null
			},
			"finish_reason": "tool_calls"
		}],
		"created": 1766935079,
		"system_fingerprint": null,
		"vertex_ai_safety_results": [],
		"vertex_ai_citation_metadata": [],
		"vertex_ai_grounding_metadata": [],
		"vertex_ai_url_context_metadata": []
	}`

	r := response.Response{}
	err := json.Unmarshal([]byte(jsonStr), &r)
	assert.NoError(t, err)

	assert.Equal(t, "J0pRaYO0PJievdIP-P6LmAE", r.ID)
	assert.Equal(t, models.ModelID("gemini-2.5-flash"), r.Model)
	assert.Equal(t, 179, r.Usage.TotalTokens)
	assert.Equal(t, 94, r.Usage.PromptTokens)
	assert.Equal(t, 85, r.Usage.CompletionTokens)
	assert.Equal(t, 0, r.Usage.PromptTokensDetails.AudioTokens)
	assert.Equal(t, 0, r.Usage.PromptTokensDetails.ImageTokens)
	assert.Equal(t, 0, r.Usage.PromptTokensDetails.CachedTokens)
	assert.Equal(t, 0, r.Usage.CompletionTokensDetails.AudioTokens)
	assert.Equal(t, 0, r.Usage.CompletionTokensDetails.AcceptedPredictionTokens)
	assert.Equal(t, 0, r.Usage.CompletionTokensDetails.RejectedPredictionTokens)

	choice := r.Choice()
	assert.Equal(t, "assistant", choice.Message.Role)
	assert.Equal(t, response.FINISH_REASON_TOOL, choice.FinishReason)

	toolCall := choice.Message.ToolCalls[0]
	assert.Equal(t, "call_3a095b3222434581ae6df726c290", toolCall.ID)
	assert.Equal(t, "function", toolCall.Type)
	assert.Equal(t, 0, toolCall.Index)
	assert.Equal(t, "current_time", toolCall.Function.Name)
}
