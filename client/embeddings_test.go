package client_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andrejsstepanovs/go-litellm/client"
	"github.com/andrejsstepanovs/go-litellm/models"
	"github.com/andrejsstepanovs/go-litellm/request"
	"github.com/andrejsstepanovs/go-litellm/response"
)

const embeddingResponseFixture = `{"object":"list","data":[{"object":"embedding","index":0,"embedding":[0.1,0.2,0.3]}],"model":"test-embedding-model"}`

func newEmbeddingsTestClient(t *testing.T, handler http.HandlerFunc) client.Litellm {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	testURL, err := url.Parse(server.URL)
	require.NoError(t, err)

	conn := getConn()
	conn.URL = *testURL
	return client.Litellm{Config: getConfig(), Connection: conn}
}

func captureEmbeddingRequestBody(t *testing.T, r *http.Request) map[string]any {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	require.NoError(t, err)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(body, &parsed))
	return parsed
}

func TestEmbeddings_Functional(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}
	ctx := context.Background()
	modelMeta := models.ModelMeta{ModelId: testEmbeddingModel}

	t.Run("text", func(t *testing.T) {
		resp, err := clientInstance.Embeddings(ctx, modelMeta, request.NewTextInput("This is a test sentence."))
		require.NoError(t, err)
		require.Equal(t, "list", resp.Object)
		require.NotEmpty(t, resp.Data)
		assert.Greater(t, resp.Data[0].Index, -1)
		assert.Greater(t, len(resp.Data[0].Embedding.Float32()), 0)
		assert.Greater(t, resp.Usage.TotalTokens, 0)
	})

	t.Run("text batch", func(t *testing.T) {
		resp, err := clientInstance.Embeddings(ctx, modelMeta, request.NewTextsInput([]string{"one", "two", "three"}))
		require.NoError(t, err)
		require.Len(t, resp.Data, 3)
	})

	t.Run("media text and image", func(t *testing.T) {
		input := request.NewMediaInput(request.EmbeddingMedia{
			Text:  "a red square",
			Image: []string{embeddingTestImage},
		})
		resp, err := clientInstance.Embeddings(ctx, models.ModelMeta{ModelId: testEmbeddingMediaModel}, input)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Data)
		assert.Greater(t, len(resp.Data[0].Embedding.Float32()), 0)
	})

	t.Run("image only", func(t *testing.T) {
		input := request.NewMediaInput(request.EmbeddingMedia{
			Image: []string{embeddingTestImage},
		})
		resp, err := clientInstance.Embeddings(ctx, models.ModelMeta{ModelId: testEmbeddingMediaModel}, input)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Data)
		assert.Greater(t, len(resp.Data[0].Embedding.Float32()), 0)
	})

	t.Run("audio only", func(t *testing.T) {
		input := request.NewMediaInput(request.EmbeddingMedia{
			Audio: []string{embeddingTestAudio},
		})
		resp, err := clientInstance.Embeddings(ctx, models.ModelMeta{ModelId: testEmbeddingMediaModel}, input)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Data)
		assert.Greater(t, len(resp.Data[0].Embedding.Float32()), 0)
	})

	t.Run("audio url", func(t *testing.T) {
		input := request.NewMediaInput(request.EmbeddingMedia{
			Audio: []string{embeddingTestAudioURL},
		})
		resp, err := clientInstance.Embeddings(ctx, models.ModelMeta{ModelId: testEmbeddingMediaModel}, input)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Data)
		assert.Greater(t, len(resp.Data[0].Embedding.Float32()), 0)
	})

	t.Run("image url", func(t *testing.T) {
		input := request.NewMediaInput(request.EmbeddingMedia{
			Image: []string{embeddingTestImageURL},
		})
		resp, err := clientInstance.Embeddings(ctx, models.ModelMeta{ModelId: testEmbeddingMediaModel}, input)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Data)
		assert.Greater(t, len(resp.Data[0].Embedding.Float32()), 0)
	})
}

func TestEmbeddings(t *testing.T) {
	clientInstance := newEmbeddingsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(embeddingResponseFixture))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	t.Run("success", func(t *testing.T) {
		resp, err := clientInstance.Embeddings(context.Background(), models.ModelMeta{ModelId: "test-embedding-model"}, request.NewTextInput("test input"))
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Data)
		assert.Equal(t, "list", resp.Object)
		assert.Equal(t, 0, resp.Data[0].Index)
		assert.Equal(t, response.Embedding(response.Embedding{0.1, 0.2, 0.3}), resp.Data[0].Embedding)
	})

	testCases := []struct {
		name  string
		input request.EmbeddingInput
		want  any // expected `input` value in the request body
	}{
		{
			name:  "single text",
			input: request.NewTextInput("hello"),
			want:  "hello",
		},
		{
			name:  "text batch",
			input: request.NewTextsInput([]string{"a", "b", "c"}),
			want:  []any{"a", "b", "c"},
		},
		{
			name:  "media text and image",
			input: request.NewMediaInput(request.EmbeddingMedia{Text: "a red square", Image: []string{"/tmp/red.png"}}),
			want:  map[string]any{"text": "a red square", "image": []any{"/tmp/red.png"}},
		},
		{
			name:  "image only",
			input: request.NewMediaInput(request.EmbeddingMedia{Image: []string{"/tmp/red.png"}}),
			want:  map[string]any{"image": []any{"/tmp/red.png"}},
		},
		{
			name:  "audio only",
			input: request.NewMediaInput(request.EmbeddingMedia{Audio: []string{"/tmp/tone.wav"}}),
			want:  map[string]any{"audio": []any{"/tmp/tone.wav"}},
		},
		{
			name:  "media batch",
			input: request.NewMediaListInput([]request.EmbeddingMedia{{Text: "a"}, {Image: []string{"/b.png"}}}),
			want:  []any{map[string]any{"text": "a"}, map[string]any{"image": []any{"/b.png"}}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotBody map[string]any
			clientInstance := newEmbeddingsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				gotBody = captureEmbeddingRequestBody(t, r)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(embeddingResponseFixture))
			})

			_, err := clientInstance.Embeddings(context.Background(), models.ModelMeta{ModelId: "test-embedding-model"}, tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, gotBody["input"])
			assert.Equal(t, "test-embedding-model", gotBody["model"])
		})
	}

	t.Run("dimensions and prompt name", func(t *testing.T) {
		var gotBody map[string]any
		clientInstance := newEmbeddingsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			gotBody = captureEmbeddingRequestBody(t, r)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(embeddingResponseFixture))
		})

		req := request.EmbeddingRequest{
			Model:      "test-embedding-model",
			Input:      request.NewTextInput("hello"),
			Dimensions: 256,
			PromptName: "Retrieval-query",
		}
		_, err := clientInstance.EmbeddingsRequest(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, float64(256), gotBody["dimensions"])
		assert.Equal(t, "Retrieval-query", gotBody["prompt_name"])
	})

	t.Run("validation errors", func(t *testing.T) {
		testCases := []struct {
			name  string
			input request.EmbeddingInput
		}{
			{name: "empty input", input: request.EmbeddingInput{}},
			{name: "empty texts", input: request.NewTextsInput([]string{"", " "})},
			{name: "empty media", input: request.NewMediaInput(request.EmbeddingMedia{})},
			{name: "media with empty entries", input: request.NewMediaInput(request.EmbeddingMedia{Image: []string{""}})},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := clientInstance.Embeddings(context.Background(), models.ModelMeta{ModelId: "test-embedding-model"}, tc.input)
				assert.Error(t, err)
			})
		}
	})
}

func TestEmbeddings_RequestValidation(t *testing.T) {
	clientInstance := newEmbeddingsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(embeddingResponseFixture))
	})

	testCases := []struct {
		name string
		req  request.EmbeddingRequest
		err  string
	}{
		{
			name: "missing model",
			req:  request.EmbeddingRequest{Input: request.NewTextInput("a")},
			err:  "model must not be empty",
		},
		{
			name: "missing input",
			req:  request.EmbeddingRequest{Model: "m"},
			err:  "input must not be empty",
		},
		{
			name: "prompt name with media",
			req: request.EmbeddingRequest{
				Model:      "m",
				Input:      request.NewMediaInput(request.EmbeddingMedia{Image: []string{"/a.png"}}),
				PromptName: "Retrieval-query",
			},
			err: "prompt_name is only valid with plain string inputs",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := clientInstance.EmbeddingsRequest(context.Background(), tc.req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.err)
		})
	}

	t.Run("no server call on validation error", func(t *testing.T) {
		calls := 0
		clientInstance := newEmbeddingsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			calls++
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(embeddingResponseFixture))
		})

		_, err := clientInstance.Embeddings(context.Background(), models.ModelMeta{ModelId: "m"}, request.EmbeddingInput{})
		assert.Error(t, err)
		assert.Equal(t, 0, calls)
	})
}

func TestEmbeddings_ErrorScenarios(t *testing.T) {
	testCases := []struct {
		name              string
		mockServerHandler http.HandlerFunc
		input             request.EmbeddingInput
		expectedErrStr    string
	}{
		{
			name: "5xx server error",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Internal Server Error"))
			},
			input:          request.NewTextInput("test input"),
			expectedErrStr: "Internal Server Error",
		},
		{
			name: "malformed json response",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"object":"list","data":[{"object":"embedding","index":0,"embedding":[0.1,0.2,0.3]`))
			},
			input:          request.NewTextInput("test input"),
			expectedErrStr: "failed to parse",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clientInstance := newEmbeddingsTestClient(t, tc.mockServerHandler)

			_, err := clientInstance.Embeddings(context.Background(), models.ModelMeta{ModelId: "test-model"}, tc.input)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErrStr)
		})
	}
}
