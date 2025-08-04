package common

import (
	"sort"
)

type ToolCalls []ToolCall

type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Index    int              `json:"index"`
	Function ToolCallFunction `json:"function"`
}

func (tc ToolCalls) SortASC() ToolCalls {
	sort.Slice(tc, func(i, j int) bool {
		return tc[i].Index < tc[j].Index
	})

	return tc
}
