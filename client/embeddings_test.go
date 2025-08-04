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
	"github.com/andrejsstepanovs/go-litellm/response"
)

func TestEmbeddings_Functional(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}
	ctx := context.Background()
	modelMeta := models.ModelMeta{ModelId: testEmbeddingModel}
	inputText := "This is a test sentence."

	resp, err := clientInstance.Embeddings(ctx, modelMeta, inputText)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Data)
	assert.Equal(t, resp.Object, "list")
	assert.Greater(t, len(resp.Data), 0)
	assert.Greater(t, resp.Data[0].Index, -1)
	assert.Greater(t, len(resp.Data[0].Embedding.Float32()), 0)
}

func TestEmbeddings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"object":"list","data":[{"object":"embedding","index":0,"embedding":[0.1,0.2,0.3]}],"model":"test-embedding-model"}`))
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

	t.Run("success", func(t *testing.T) {
		resp, err := clientInstance.Embeddings(context.Background(), models.ModelMeta{ModelId: "test-embedding-model"}, "test input")
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Data)
		assert.Equal(t, resp.Object, "list")
		assert.Equal(t, 0, resp.Data[0].Index)
		assert.Equal(t, response.Embedding(response.Embedding{0.1, 0.2, 0.3}), resp.Data[0].Embedding)
	})

	t.Run("fail empty input", func(t *testing.T) {
		_, err := clientInstance.Embeddings(context.Background(), models.ModelMeta{ModelId: "test-embedding-model"}, "")
		assert.Error(t, err)
	})
}

func TestEmbeddings_ErrorScenarios(t *testing.T) {
	testCases := []struct {
		name              string
		mockServerHandler http.HandlerFunc
		input             string
		expectedErrStr    string
	}{
		{
			name: "5xx server error",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Internal Server Error"))
			},
			input:          "test input",
			expectedErrStr: "Internal Server Error",
		},
		{
			name: "malformed json response",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"object":"list","data":[{"object":"embedding","index":0,"embedding":[0.1,0.2,0.3]`))
			},
			input:          "test input",
			expectedErrStr: "failed to parse",
		},
		{
			name:              "empty input",
			mockServerHandler: nil, // No server call expected
			input:             "",  // Empty input
			expectedErrStr:    "inputText cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockServerHandler != nil {
				server := httptest.NewServer(tc.mockServerHandler)
				defer server.Close()

				testUrl, err := url.Parse(server.URL)
				assert.NoError(t, err)

				c := getConfig()
				conn := getConn()
				conn.URL = *testUrl
				clientInstance := client.Litellm{Config: c, Connection: conn}

				_, err = clientInstance.Embeddings(context.Background(), models.ModelMeta{ModelId: "test-model"}, tc.input)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrStr)
			} else {
				// For validation errors, no server call is made
				clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}
				_, err := clientInstance.Embeddings(context.Background(), models.ModelMeta{ModelId: "test-model"}, tc.input)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrStr)
			}
		})
	}
}
