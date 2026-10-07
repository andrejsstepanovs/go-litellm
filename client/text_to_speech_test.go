package client_test

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/client"
	"github.com/andrejsstepanovs/go-litellm/models"
	"github.com/andrejsstepanovs/go-litellm/request"
)

func TestTextToSpeech_Functional(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping functional test in short mode")
		return
	}

	clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}

	available := map[string]bool{}
	for _, m := range mustAvailableModels(t, clientInstance) {
		available[string(m.ID)] = true
	}

	testCases := []struct {
		name              string
		modelName         string
		voice             string
		expectedText      string
		expectedWordCount int
		input             string
		format            string
	}{
		{
			name:      "openai",
			input:     "I like cats",
			modelName: string(testTTSOne),
			voice:     "alloy",
			format:    "mp3",
		},
		{
			name:      "gemini",
			modelName: string(testTTSTwo),
			input:     "I like dogs",
			voice:     "alloy",
			format:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !available[tc.modelName] {
				t.Skipf("model %q is not accessible with the configured API key", tc.modelName)
			}

			req := request.Speech{
				Input:          tc.input,
				Model:          models.ModelID(tc.modelName),
				Voice:          tc.voice,
				ResponseFormat: tc.format,
			}
			res, err := clientInstance.TextToSpeech(
				context.Background(),
				req,
			)

			assert.NoError(t, err)
			assert.NotEmpty(t, res.Full)
			assert.NotEmpty(t, res.Name)
			assert.NotEmpty(t, res.Directory)

			log.Printf("%v+", res)
		})
	}
}
