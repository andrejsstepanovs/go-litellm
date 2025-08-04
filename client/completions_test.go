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
	"github.com/andrejsstepanovs/go-litellm/response"
)

func TestCompletions_Functional(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	t.Run("success", func(t *testing.T) {
		clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}
		ctx := context.Background()

		messages := request.Messages{}
		messages.AddMessagePair(
			request.UserMessageSimple("Hello my name is Bobik"),
			request.UserMessageSimple("Hi!"),
		)
		messages.AddMessage(request.UserMessageSimple("What is my name again?"))

		modelMeta, err := clientInstance.Model(ctx, testModel)
		assert.NoError(t, err)

		var temp float32 = 0.2
		req := request.NewCompletionRequest(modelMeta, messages, request.LLMCallTools{}, &temp, 0.7)
		resp, err := clientInstance.Completion(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp)
		assert.Contains(t, resp.String(), "Bobik")
		assert.Equal(t, resp.Choice().FinishReason, response.FINISH_REASON_STOP)
		assert.Equal(t, resp.Object, "chat.completion")
		assert.Greater(t, resp.Usage.PromptTokens, 0)
		assert.Greater(t, resp.Usage.CompletionTokens, 0)
		assert.Greater(t, resp.Usage.TotalTokens, 0)
		assert.GreaterOrEqual(t, resp.Usage.QueueTime, 0.0)
		assert.GreaterOrEqual(t, resp.Usage.PromptTime, 0.0)
		assert.GreaterOrEqual(t, resp.Usage.CompletionTime, 0.0)
	})
}

func TestCompletions_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"error" : {"message" : "litellm.BadRequestError: GroqException - {\"error\":{\"message\":\"tool call validation failed: parameters for tool json_tool_call did not match schema: errors: [/result: expected object, but got boolean] LiteLLM Retried: 1 times, LiteLLM Max Retries: 2","type" : null,"param" : null,"code" : "400"}}`))
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
	req := request.NewCompletionRequest(models.ModelMeta{ModelId: "test"}, request.Messages{{
		Role: request.ROLE_USER,
		Contents: request.MessageContents{
			{
				Type: "text",
				Text: "test",
			},
		},
	}}, request.LLMCallTools{}, nil, 0.7)

	_, err = clientInstance.Completion(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "client error: litellm.BadRequestError: GroqException")
	assert.ErrorContains(t, err, "LiteLLM Retried: 1 times")
	assert.ErrorContains(t, err, "parameters for tool json_tool_call did not match schema")
}

func TestCompletions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"id":"chatcmpl-a32405f8-99a4-468f-8f43-ae543beb0d12","created":1747237215,"model":"llama-3.1-8b-instant","object":"chat.completion","system_fingerprint":"fp_5ac2c9fbe7","choices":[{"finish_reason":"stop","index":0,"message":{"content":"Your name is Bobik. How are you today?","role":"assistant","tool_calls":null,"function_call":null}}],"usage":{"completion_tokens":12,"prompt_tokens":59,"total_tokens":71,"completion_tokens_details":null,"prompt_tokens_details":null,"queue_time":0.16690093,"prompt_time":0.010938355,"completion_time":0.016,"total_time":0.026938355},"usage_breakdown":{"models":null},"x_groq":{"id":"req_01jv7q86yqf74877qm2y952k7w"}}`))
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
		req := request.NewCompletionRequest(models.ModelMeta{ModelId: "test"}, request.Messages{{
			Role: request.ROLE_USER,
			Contents: request.MessageContents{
				{
					Type: "text",
					Text: "test",
				},
			},
		}}, request.LLMCallTools{}, nil, 0.7)
		resp, err := clientInstance.Completion(context.Background(), req)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp)
		assert.Contains(t, resp.String(), "Bobik")
		assert.Equal(t, resp.Choice().FinishReason, response.FINISH_REASON_STOP)
		assert.Equal(t, resp.Object, "chat.completion")
		assert.Greater(t, resp.Usage.PromptTokens, 0)
		assert.Greater(t, resp.Usage.CompletionTokens, 0)
		assert.Greater(t, resp.Usage.TotalTokens, 0)
		assert.Greater(t, resp.Usage.QueueTime, 0.0)
		assert.Greater(t, resp.Usage.PromptTime, 0.0)
		assert.Greater(t, resp.Usage.CompletionTime, 0.0)
		assert.Greater(t, resp.Usage.TotalTime, 0.0)
	})

	t.Run("fail empty model", func(t *testing.T) {
		req := request.NewCompletionRequest(models.ModelMeta{}, request.Messages{{
			Role: request.ROLE_USER,
			Contents: request.MessageContents{
				{
					Type: "text",
					Text: "test",
				},
			},
		}}, request.LLMCallTools{}, nil, 0.7)

		_, err = clientInstance.Completion(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("fail empty messages", func(t *testing.T) {
		req := request.NewCompletionRequest(models.ModelMeta{ModelId: "test"}, request.Messages{}, request.LLMCallTools{}, nil, 0.7)
		_, err = clientInstance.Completion(context.Background(), req)
		assert.Error(t, err)
	})
}

func TestCompletion_ErrorScenarios(t *testing.T) {
	// Define a helper to create a basic valid request
	validRequest := func() *request.Request {
		return request.NewCompletionRequest(
			models.ModelMeta{ModelId: "test-model"},
			request.Messages{{Role: request.ROLE_USER, Contents: request.MessageContents{{Type: "text", Text: "test"}}}},
			request.LLMCallTools{},
			nil,
			0.7,
		)
	}

	testCases := []struct {
		name              string
		mockServerHandler http.HandlerFunc
		request           *request.Request
		expectedErrStr    string
	}{
		{
			name: "5xx server error",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Internal Server Error"))
			},
			request:        validRequest(),
			expectedErrStr: "failed to parse completions response",
		},
		{
			name: "malformed json response",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"id":"chatcmpl-bad-json","created":,"model":"test-model","object":"chat.completion"}`))
			},
			request:        validRequest(),
			expectedErrStr: "failed to parse completions response",
		},
		{
			name: "empty response body",
			mockServerHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
			},
			request:        validRequest(),
			expectedErrStr: "failed to parse completions response",
		},
		{
			name:              "empty model id",
			mockServerHandler: nil, // No server call expected
			request: request.NewCompletionRequest(
				models.ModelMeta{ModelId: ""}, // Empty model ID
				request.Messages{{Role: request.ROLE_USER, Contents: request.MessageContents{{Type: "text", Text: "test"}}}},
				request.LLMCallTools{},
				nil,
				0.7,
			),
			expectedErrStr: "modelID cannot be empty",
		},
		{
			name:              "empty messages",
			mockServerHandler: nil, // No server call expected
			request: request.NewCompletionRequest(
				models.ModelMeta{ModelId: "test-model"},
				request.Messages{}, // Empty messages
				request.LLMCallTools{},
				nil,
				0.7,
			),
			expectedErrStr: "messages cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			var server *httptest.Server
			if tc.mockServerHandler != nil {
				server = httptest.NewServer(tc.mockServerHandler)
				defer server.Close()

				testUrl, err := url.Parse(server.URL)
				assert.NoError(t, err)

				// Create client with mock server URL
				c := getConfig()
				conn := getConn()
				conn.URL = *testUrl
				clientInstance := client.Litellm{Config: c, Connection: conn}

				// Execute
				_, err = clientInstance.Completion(context.Background(), tc.request)

				// Assert
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrStr)
			} else {
				// For validation errors, no server call is made
				clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}
				_, err := clientInstance.Completion(context.Background(), tc.request)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrStr)
			}
		})
	}
}
