package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/client"
	"github.com/andrejsstepanovs/go-litellm/models"
	"github.com/andrejsstepanovs/go-litellm/request"
)

func TestExtraHeaders_AppliedToAllRequestTypes(t *testing.T) {
	extraHeaders := map[string]string{
		"X-App-Name": "YourAppName",
		"X-User-Id":  "user_123",
	}

	assertExtraHeaders := func(t *testing.T, captured http.Header) {
		t.Helper()
		for key, value := range extraHeaders {
			assert.Equal(t, value, captured.Get(key), "expected header %q to be %q", key, value)
		}
	}

	t.Run("Completion", func(t *testing.T) {
		var captured http.Header
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			captured = r.Header.Clone()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"1","object":"chat.completion","choices":[{"finish_reason":"stop","index":0,"message":{"role":"assistant","content":"hi"}}]}`))
		}))
		defer server.Close()

		clientInstance := newTestClientWithHeaders(t, server.URL, extraHeaders)
		req := request.NewCompletionRequest(models.ModelMeta{ModelId: "test"}, request.Messages{{
			Role: request.ROLE_USER,
			Contents: request.MessageContents{
				{Type: "text", Text: "test"},
			},
		}}, request.LLMCallTools{}, nil, 0.7)

		_, err := clientInstance.Completion(context.Background(), req)
		assert.NoError(t, err)
		assertExtraHeaders(t, captured)
	})

	t.Run("Embeddings", func(t *testing.T) {
		var captured http.Header
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			captured = r.Header.Clone()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"object":"list","data":[{"object":"embedding","index":0,"embedding":[0.1,0.2,0.3]}],"model":"test"}`))
		}))
		defer server.Close()

		clientInstance := newTestClientWithHeaders(t, server.URL, extraHeaders)
		_, err := clientInstance.Embeddings(context.Background(), models.ModelMeta{ModelId: "test"}, "test input")
		assert.NoError(t, err)
		assertExtraHeaders(t, captured)
	})

	t.Run("TokenCounter", func(t *testing.T) {
		var captured http.Header
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			captured = r.Header.Clone()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"total_tokens":3,"model_used":"test","request_model":"test","tokenizer_type":"test"}`))
		}))
		defer server.Close()

		clientInstance := newTestClientWithHeaders(t, server.URL, extraHeaders)
		_, err := clientInstance.TokenCounter(context.Background(), &request.TokenCounterRequest{Model: "test"})
		assert.NoError(t, err)
		assertExtraHeaders(t, captured)
	})

	t.Run("Models", func(t *testing.T) {
		var captured http.Header
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			captured = r.Header.Clone()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"object":"list","data":[{"id":"test","object":"model","owned_by":"test"}]}`))
		}))
		defer server.Close()

		clientInstance := newTestClientWithHeaders(t, server.URL, extraHeaders)
		_, err := clientInstance.Models(context.Background())
		assert.NoError(t, err)
		assertExtraHeaders(t, captured)
	})
}

func TestExtraHeaders_NotSetWhenEmpty(t *testing.T) {
	var captured http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"object":"list","data":[{"id":"test","object":"model","owned_by":"test"}]}`))
	}))
	defer server.Close()

	clientInstance := newTestClientWithHeaders(t, server.URL, nil)
	_, err := clientInstance.Models(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, captured.Get("X-App-Name"))
	assert.Empty(t, captured.Get("X-User-Id"))
}

func TestExtraHeaders_TrimmedAndEmptyEntriesSkipped(t *testing.T) {
	var captured http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"object":"list","data":[{"id":"test","object":"model","owned_by":"test"}]}`))
	}))
	defer server.Close()

	clientInstance := newTestClientWithHeaders(t, server.URL, map[string]string{
		"  X-App-Name  ": "  YourAppName  ",
		"":               "orphan-value",
		"X-Empty-Value":  "   ",
		"   ":            "orphan-key",
	})
	_, err := clientInstance.Models(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "YourAppName", captured.Get("X-App-Name"))
	assert.Empty(t, captured.Get("X-Empty-Value"))
	assert.Empty(t, captured.Get(""))
}

func newTestClientWithHeaders(t *testing.T, serverURL string, extraHeaders map[string]string) client.Litellm {
	t.Helper()

	testURL, err := url.Parse(serverURL)
	assert.NoError(t, err)

	c := getConfig()
	c.ExtraHeaders = extraHeaders
	conn := getConn()
	conn.URL = *testURL

	return client.Litellm{Config: c, Connection: conn}
}
