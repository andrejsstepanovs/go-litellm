package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/client"
	"github.com/andrejsstepanovs/go-litellm/common"
	"github.com/andrejsstepanovs/go-litellm/response"
)

func Test_ToolCall_Functional(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	t.Run("success", func(t *testing.T) {
		clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}

		tool := common.ToolCallFunction{
			Name: "current_time",
			Arguments: common.Arguments{
				"timezone": "Europe/Riga",
			},
		}

		res, err := clientInstance.ToolCall(context.Background(), tool)

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Len(t, res, 1)
		assert.NotContains(t, res.String(), "must be")
		assert.NotEmpty(t, res.String())
	})
}

func Test_ToolCall(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`[{"type":"text","text":"2025-05-14 16:27:17","annotations":null}]`))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}))
		defer server.Close()

		testUrl, err := url.Parse(server.URL)
		assert.NoError(t, err)

		c := getConfig()
		conn := getConn()
		conn.URL = *testUrl
		clientInstance := client.Litellm{Config: c, Connection: conn}

		res, err := clientInstance.ToolCall(context.Background(), common.ToolCallFunction{Name: "current_time"})

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Len(t, res, 1)
		assert.Equal(t, "2025-05-14 16:27:17", res.String())
		assert.Equal(t, response.ToolResponses{
			{
				Type:        "text",
				Text:        "2025-05-14 16:27:17",
				Annotations: nil,
			},
		}, res)
	})

	t.Run("success with multiple responses - should combine with newlines", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			// Simulate API returning multiple responses (edge case that causes duplicate IDs)
			_, err := w.Write([]byte(`[
				{"type":"text","text":"First response","annotations":null},
				{"type":"text","text":"Second response","annotations":null},
				{"type":"text","text":"Third response","annotations":null}
			]`))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}))
		defer server.Close()

		testUrl, err := url.Parse(server.URL)
		assert.NoError(t, err)

		c := getConfig()
		conn := getConn()
		conn.URL = *testUrl
		clientInstance := client.Litellm{Config: c, Connection: conn}

		res, err := clientInstance.ToolCall(context.Background(), common.ToolCallFunction{Name: "test_tool"})

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		// Should return one response with all text combined by newlines
		assert.Len(t, res, 1)
		assert.Equal(t, "First response\nSecond response\nThird response", res.String())
		assert.Equal(t, response.ToolResponses{
			{
				Type:        "text",
				Text:        "First response\nSecond response\nThird response",
				Annotations: nil,
			},
		}, res)
	})

	t.Run("success with multiple responses including empty ones", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			// Simulate API returning multiple responses with some empty ones
			_, err := w.Write([]byte(`[
				{"type":"text","text":"First response","annotations":null},
				{"type":"text","text":"","annotations":null},
				{"type":"text","text":"Third response","annotations":null}
			]`))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}))
		defer server.Close()

		testUrl, err := url.Parse(server.URL)
		assert.NoError(t, err)

		c := getConfig()
		conn := getConn()
		conn.URL = *testUrl
		clientInstance := client.Litellm{Config: c, Connection: conn}

		res, err := clientInstance.ToolCall(context.Background(), common.ToolCallFunction{Name: "test_tool"})

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Len(t, res, 1)
		// Should preserve empty responses as blank lines
		assert.Equal(t, "First response\n\nThird response", res.String())
		assert.Equal(t, response.ToolResponses{
			{
				Type:        "text",
				Text:        "First response\n\nThird response",
				Annotations: nil,
			},
		}, res)
	})
}
