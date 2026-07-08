package request_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/models"
	"github.com/andrejsstepanovs/go-litellm/request"
)

func TestRequest_SetReasoningEffort(t *testing.T) {
	tests := []struct {
		name     string
		effort   string
		model    models.ModelMeta
		expected string
	}{
		{
			name:     "sets effort when model supports reasoning",
			effort:   "low",
			model:    models.ModelMeta{ModelId: "gemini/gemini-3-flash", SupportsReasoning: true},
			expected: "low",
		},
		{
			name:     "no-op when model does not support reasoning",
			effort:   "low",
			model:    models.ModelMeta{ModelId: "gpt-4o", SupportsReasoning: false},
			expected: "",
		},
		{
			name:     "no-op when effort is empty",
			effort:   "",
			model:    models.ModelMeta{ModelId: "gemini/gemini-3-flash", SupportsReasoning: true},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := request.NewRequest(tt.model)
			r.SetReasoningEffort(tt.effort, tt.model)
			assert.Equal(t, tt.expected, r.ReasoningEffort)
		})
	}
}
