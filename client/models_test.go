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
)

func Test_Models_Functional(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	t.Run("success", func(t *testing.T) {
		clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}
		res, err := clientInstance.Models(context.Background())

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}

func TestModels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"data":[{"id":"sonnet-3-7","object":"model","created":1677610602,"owned_by":"openai"},{"id":"groq-llama-3.3-70b","object":"model","created":1677610602,"owned_by":"openai"}],"object":"list"}`))
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

		res, err := clientInstance.Models(context.Background())

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Len(t, res, 2)
		assert.Equal(t, models.Models{
			{ID: "sonnet-3-7", Object: "model", OwnedBy: "openai"},
			{ID: "groq-llama-3.3-70b", Object: "model", OwnedBy: "openai"},
		}, res)
	})
}

func TestModels_ErrorScenarios(t *testing.T) {
	testCases := []struct {
		name           string
		handler        http.HandlerFunc
		expectedErrStr string
	}{
		{
			name: "5xx server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Internal Server Error"))
			},
			expectedErrStr: "failed to parse models response",
		},
		{
			name: "empty data array",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":[],"object":"list"}`))
			},
			expectedErrStr: "",
		},
		{
			name: "malformed json",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":[{"id":"sonnet-3-7","object":"model","created":1677610602,"owned_by":"openai"`))
			},
			expectedErrStr: "failed to parse models response",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			testUrl, err := url.Parse(server.URL)
			assert.NoError(t, err)

			c := getConfig()
			conn := getConn()
			conn.URL = *testUrl
			clientInstance := client.Litellm{Config: c, Connection: conn}

			res, err := clientInstance.Models(context.Background())

			if tc.expectedErrStr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrStr)
			} else {
				assert.NoError(t, err)
				if tc.name == "empty data array" {
					assert.Empty(t, res)
				}
			}
		})
	}
}

func Test_Model_Functional(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	t.Run("success", func(t *testing.T) {
		clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}
		ctx := context.Background()

		info, err := clientInstance.Model(ctx, "claude-4")
		assert.NoError(t, err)
		assert.NotEmpty(t, info)
		assert.Equal(t, string(info.ModelId), "claude-4")
		assert.Equal(t, 1e+06, info.MaxInputTokens)
		assert.True(t, info.SupportsVision)
	})
}

func TestModel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"data":[{"model_group":"claude-4","providers":["anthropic"],"max_input_tokens":200000.0,"max_output_tokens":128000.0,"input_cost_per_token":3e-06,"output_cost_per_token":1.5e-05,"mode":"chat","tpm":null,"rpm":50,"supports_parallel_function_calling":false,"supports_vision":true,"supports_web_search":false,"supports_reasoning":true,"supports_function_calling":true,"supported_openai_params":["stream","stop","temperature","top_p","max_tokens","max_completion_tokens","tools","tool_choice","extra_headers","parallel_tool_calls","response_format","user","reasoning_effort","thinking"],"configurable_clientside_auth_params":null}]}`))
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

		info, err := clientInstance.Model(context.Background(), "claude-4")
		assert.NoError(t, err)
		assert.NotEmpty(t, info)
		assert.Equal(t, string(info.ModelId), "claude-4")
		assert.Equal(t, info.MaxInputTokens, 200000.0)
		assert.True(t, info.SupportsVision)
	})
}

func TestModel_ErrorScenarios(t *testing.T) {
	testCases := []struct {
		name           string
		handler        http.HandlerFunc
		modelID        string
		expectedErrStr string
	}{
		{
			name: "5xx server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Internal Server Error"))
			},
			modelID:        "test-model",
			expectedErrStr: "failed to parse model response",
		},
		{
			name: "empty model info array",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":[]}`))
			},
			modelID:        "test-model",
			expectedErrStr: `multiple or no models found for "test-model"`,
		},
		{
			name: "multiple models in response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":[{"model_group":"test-model"},{"model_group":"another-model"}]}`))
			},
			modelID:        "test-model",
			expectedErrStr: `multiple or no models found for "test-model"`,
		},
		{
			name: "malformed json",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":[{"model_group":"test-model","providers":["anthropic"],"max_input_tokens":200000.0`))
			},
			modelID:        "test-model",
			expectedErrStr: "failed to parse model response",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()

			testUrl, err := url.Parse(server.URL)
			assert.NoError(t, err)

			c := getConfig()
			conn := getConn()
			conn.URL = *testUrl
			clientInstance := client.Litellm{Config: c, Connection: conn}

			_, err = clientInstance.Model(context.Background(), models.ModelID(tc.modelID))

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErrStr)
		})
	}
}

func Test_Model_SupportedParams_Functional(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	t.Run("success", func(t *testing.T) {
		clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}
		model, err := clientInstance.Model(context.Background(), "claude-4")
		assert.NoError(t, err)

		assert.Contains(t, model.SupportedOpenAIParams, "stream")
		assert.Contains(t, model.SupportedOpenAIParams, "temperature")
		assert.Contains(t, model.SupportedOpenAIParams, "tools")
		assert.Contains(t, model.SupportedOpenAIParams, "thinking")
		assert.Contains(t, model.SupportedOpenAIParams, "reasoning_effort")
		assert.Contains(t, model.SupportedOpenAIParams, "response_format")
		assert.Contains(t, model.SupportedOpenAIParams, "parallel_tool_calls")
		assert.Contains(t, model.SupportedOpenAIParams, "tool_choice")
		assert.Contains(t, model.SupportedOpenAIParams, "user")
		assert.Contains(t, model.SupportedOpenAIParams, "max_completion_tokens")
		assert.Contains(t, model.SupportedOpenAIParams, "max_tokens")
		assert.Contains(t, model.SupportedOpenAIParams, "top_p")
	})
}
