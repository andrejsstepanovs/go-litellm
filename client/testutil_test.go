package client_test

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/andrejsstepanovs/go-litellm/client"
	"github.com/andrejsstepanovs/go-litellm/conf/connections/litellm"
	"github.com/andrejsstepanovs/go-litellm/models"
)

const testModelGood = models.ModelID("active")
const testModel = models.ModelID("active")
const testEmbeddingModel = models.ModelID("gemma-embed")
const testEmbeddingMediaModel = models.ModelID("gemma-embed")

// Files must exist on the machine where the LiteLLM server runs; paths are resolved server-side.
const embeddingTestImage = "/tmp/red.png"
const embeddingTestAudio = "/tmp/tone.wav"

// URLs are fetched by the embedding server directly; no local file needed.
// NASA imagery is public domain.
const embeddingTestImageURL = "https://images-assets.nasa.gov/image/PIA12235/PIA12235~small.jpg"
const embeddingTestAudioURL = "https://github.com/pdx-cs-sound/wavs/raw/refs/heads/main/voice-note.wav"
const testSTTOne = models.ModelID("parakeet-tdt")
const testSTTTwo = models.ModelID("parakeet-tdt")
const testTTSOne = models.ModelID("tts-openai")
const testTTSTwo = models.ModelID("tts-gemini")

func getConn() litellm.Connection {
	return litellm.Connection{
		URL: getTestURL(),
		Targets: litellm.Targets{
			System: litellm.Target{
				Timeout:          10 * time.Second,
				RetryInterval:    0,
				RetryMaxAttempts: 0,
				RetryBackoffRate: 0,
				MaxRetry:         0,
			},
			MCP: litellm.Target{
				Timeout:          10 * time.Second,
				RetryInterval:    0,
				RetryMaxAttempts: 0,
				RetryBackoffRate: 0,
				MaxRetry:         0,
			},
			LLM: litellm.Target{
				Timeout:          10 * time.Second,
				RetryInterval:    0,
				RetryMaxAttempts: 0,
				RetryBackoffRate: 0,
				MaxRetry:         0,
			},
		},
	}
}

func getConfig() client.Config {
	return client.Config{
		APIKey:      getTestKey(),
		Temperature: 0,
	}
}

// mustAvailableModels returns the model groups accessible with the configured API key.
func mustAvailableModels(t *testing.T, c client.Litellm) models.Models {
	t.Helper()
	ms, err := c.Models(context.Background())
	if err != nil {
		t.Fatalf("failed to list models: %v", err)
	}
	return ms
}

// getTestURL returns the LiteLLM address from LITELLM_TEST_URL, defaulting to localhost:4000.
func getTestURL() url.URL {
	if raw := os.Getenv("LITELLM_TEST_URL"); raw != "" {
		u, err := url.Parse(raw)
		if err == nil && u.Host != "" {
			return *u
		}
	}
	return url.URL{Scheme: "http", Host: "localhost:4000"}
}

// getTestKey returns the API key from LITELLM_TEST_KEY, defaulting to the local proxy key.
func getTestKey() string {
	if key := os.Getenv("LITELLM_TEST_KEY"); key != "" {
		return key
	}
	return "sk-1234"
}
