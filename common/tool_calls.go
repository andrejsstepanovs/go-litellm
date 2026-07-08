package common

import (
	"sort"
)

type ToolCalls []ToolCall

type ToolCall struct {
	ID                     string           `json:"id"`
	Type                   string           `json:"type"`
	Index                  int              `json:"index"`
	Function               ToolCallFunction `json:"function"`
	ProviderSpecificFields map[string]any   `json:"provider_specific_fields,omitempty"`
}

// ThoughtSignature returns the Gemini "thought_signature" value stored in
// ProviderSpecificFields, if present. It returns an empty string when the
// field is absent, not a string, or ProviderSpecificFields is nil.
func (tc ToolCall) ThoughtSignature() string {
	if tc.ProviderSpecificFields == nil {
		return ""
	}
	sig, _ := tc.ProviderSpecificFields["thought_signature"].(string)
	return sig
}

func (tc ToolCalls) SortASC() ToolCalls {
	sort.Slice(tc, func(i, j int) bool {
		return tc[i].Index < tc[j].Index
	})

	return tc
}
