package client_test

import (
	"net/url"
	"time"

	"github.com/andrejsstepanovs/go-litellm/client"
	"github.com/andrejsstepanovs/go-litellm/conf/connections/litellm"
	"github.com/andrejsstepanovs/go-litellm/models"
)

const testModelGood = models.ModelID("claude-4")
const testModel = models.ModelID("groq-llama-3.1-8b")
const testEmbeddingModel = models.ModelID("mistral-embed")
const testSTTOne = models.ModelID("whisper-1")
const testSTTTwo = models.ModelID("deepgram-nova-2")

func getConn() litellm.Connection {
	return litellm.Connection{
		URL: url.URL{Scheme: "http", Host: "localhost:4000"},
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
		APIKey:      "sk-1234",
		Temperature: 0,
	}
}
