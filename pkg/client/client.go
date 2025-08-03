package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/andrejsstepanovs/go-litellm/pkg/audio"
	"github.com/andrejsstepanovs/go-litellm/pkg/common"
	cfg "github.com/andrejsstepanovs/go-litellm/pkg/conf/connections/litellm"
	"github.com/andrejsstepanovs/go-litellm/pkg/httpresp"
	"github.com/andrejsstepanovs/go-litellm/pkg/mcp"
	"github.com/andrejsstepanovs/go-litellm/pkg/models"
	"github.com/andrejsstepanovs/go-litellm/pkg/request"
	"github.com/andrejsstepanovs/go-litellm/pkg/response"
	"github.com/andrejsstepanovs/go-litellm/pkg/users"
	"github.com/go-playground/validator/v10"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/opus-domini/fast-shot/constant/mime"
)

var validate = validator.New()

type Config struct {
	APIKey      string  `validate:"required"`
	Temperature float32 `validate:"required"`
}

func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("litellm config is required")
	}

	err := validate.Struct(c)
	if err != nil {
		return fmt.Errorf("litellm validation error: %w", err)
	}
	return nil
}

type Litellm struct {
	Config     Config
	Connection cfg.Connection
}

func New(config Config, connection cfg.Connection) (*Litellm, error) {
	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf("litellm config validation error: %w", err)
	}

	return &Litellm{Config: config, Connection: connection}, nil
}

func (l *Litellm) client(name cfg.TargetName) fastshot.ClientHttpMethods {
	timeout := l.Connection.Targets.Get(name).Timeout

	return fastshot.NewClient(l.Connection.URL.String()).
		Auth().BearerToken(l.Config.APIKey).
		Config().SetTimeout(timeout).
		Config().SetFollowRedirects(true).
		Header().AddUserAgent(string(name)).
		Header().AddContentType(mime.JSON).
		Build()
}

func (l *Litellm) Model(ctx context.Context, modelID models.ModelID) (models.ModelMeta, error) {
	target := l.Connection.Targets.Get(cfg.CLIENT_SYSTEM)

	resp, err := l.client(cfg.CLIENT_SYSTEM).
		GET("model_group/info").
		Query().AddParam("model_group", string(modelID)).
		Context().Set(ctx).
		Header().AddAccept(mime.JSON).
		Retry().SetExponentialBackoff(
		target.RetryInterval,
		target.RetryMaxAttempts,
		target.RetryBackoffRate).
		Send()

	if err != nil {
		return models.ModelMeta{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body().Close()

	res := struct {
		ModelInfo []models.ModelMeta `json:"data"`
	}{}
	err = httpresp.ParseHTTPResponse(*resp, &res)
	if err != nil {
		return models.ModelMeta{}, fmt.Errorf("failed to parse model response: %w", err)
	}

	if len(res.ModelInfo) == 1 {
		return res.ModelInfo[0], nil
	}

	return models.ModelMeta{}, fmt.Errorf("multiple or no models found for %q", modelID)
}

// ModelInfoMap model name => litellm model key (openrouter-qwen3-235b-a22b: openrouter/qwen/qwen3-235b-a22b)
func (l *Litellm) ModelInfoMap(ctx context.Context) (map[string]string, error) {
	target := l.Connection.Targets.Get(cfg.CLIENT_SYSTEM)

	resp, err := l.client(cfg.CLIENT_SYSTEM).
		GET("/v2/model/info").
		Context().Set(ctx).
		Header().AddAccept(mime.JSON).
		Retry().SetExponentialBackoff(
		target.RetryInterval,
		target.RetryMaxAttempts,
		target.RetryBackoffRate).
		Send()

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body().Close()

	type ModelInfo struct {
		Key string `json:"key"`
	}
	type InfoResponse struct {
		ModelName string    `json:"model_name"`
		Params    ModelInfo `json:"model_info"`
	}
	type Data struct {
		Data []InfoResponse `json:"data"`
	}

	data := &Data{}
	err = httpresp.ParseHTTPResponse(*resp, data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse model info response: %w", err)
	}

	ret := make(map[string]string)
	for _, info := range data.Data {
		ret[info.ModelName] = info.Params.Key
	}
	return ret, nil
}

func (l *Litellm) Models(ctx context.Context) (models.Models, error) {
	target := l.Connection.Targets.Get(cfg.CLIENT_SYSTEM)

	resp, err := l.client(cfg.CLIENT_SYSTEM).
		GET("/models").
		Context().Set(ctx).
		Header().AddAccept(mime.JSON).
		Retry().SetExponentialBackoff(
		target.RetryInterval,
		target.RetryMaxAttempts,
		target.RetryBackoffRate).
		Send()

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	res := struct {
		Models models.Models `json:"data"`
		Object string        `json:"object"`
	}{}
	err = httpresp.ParseHTTPResponse(*resp, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse models response: %w", err)
	}
	defer resp.Body().Close()

	return res.Models, nil
}

func (l *Litellm) ToolCall(ctx context.Context, user users.User, tool common.ToolCallFunction) (response.ToolResponses, error) {
	target := l.Connection.Targets.Get(cfg.CLIENT_MCP)

	resp, err := l.client(cfg.CLIENT_MCP).
		POST("/mcp-rest/tools/call").
		Context().Set(ctx).
		Header().AddAccept(mime.JSON).
		Header().Add("X-User-ID", strconv.FormatInt(user.ID, 10)).
		Retry().SetExponentialBackoff(
		target.RetryInterval,
		target.RetryMaxAttempts,
		target.RetryBackoffRate).
		Body().AsJSON(tool).
		Send()

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body().Close()

	var res response.ToolResponses
	err = httpresp.ParseHTTPResponse(*resp, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tool-call response: %w", err)
	}

	return res, nil
}

func (l *Litellm) Tools(ctx context.Context) (mcp.AvailableTools, error) {
	target := l.Connection.Targets.Get(cfg.CLIENT_MCP)

	resp, err := l.client(cfg.CLIENT_MCP).
		GET("/mcp-rest/tools/list").
		Context().Set(ctx).
		Header().AddAccept(mime.JSON).
		Retry().SetExponentialBackoff(
		target.RetryInterval,
		target.RetryMaxAttempts,
		target.RetryBackoffRate).
		Send()

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body().Close()

	var res mcp.AvailableToolsResponse
	err = httpresp.ParseHTTPResponse(*resp, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tools response: %w", err)
	}

	return res.Tools, nil
}

func (l *Litellm) Completion(ctx context.Context, req *request.Request) (response.Response, error) {
	if req.Model == "" {
		return response.Response{}, fmt.Errorf("modelID cannot be empty")
	}
	if len(req.Messages) == 0 {
		return response.Response{}, fmt.Errorf("messages cannot be empty")
	}

	target := l.Connection.Targets.Get(cfg.CLIENT_LLM)
	resp, err := l.client(cfg.CLIENT_LLM).
		POST("/chat/completions").
		Context().Set(ctx).
		Header().AddAccept(mime.JSON).
		Retry().SetExponentialBackoff(
		target.RetryInterval,
		target.RetryMaxAttempts,
		target.RetryBackoffRate).
		Body().AsJSON(req).
		Send()

	if err != nil {
		return response.Response{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body().Close()

	if resp.Status().Is4xxClientError() {
		var errValue response.ErrorResponse
		err = resp.Body().AsJSON(&errValue)
		if err != nil {
			return response.Response{}, fmt.Errorf("failed to parse error response: %w", err)
		}
		return response.Response{}, fmt.Errorf("client error: %s", errValue.Error.Message)
	}

	var res response.Response
	err = httpresp.ParseHTTPResponse(*resp, &res)
	if err != nil {
		return response.Response{}, fmt.Errorf("failed to parse completions response: %w", err)
	}

	return res, nil
}

func (l *Litellm) SpeechToText(ctx context.Context, model models.ModelMeta, audioFile string) (audio.AudioResponse, error) {
	url := fmt.Sprintf("%s/audio/transcriptions", l.Connection.URL.String())
	resp, err := audio.TranscribeAudio(url, l.Config.APIKey, audioFile, string(model.ModelId))
	if err != nil {
		return audio.AudioResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		return audio.AudioResponse{}, fmt.Errorf("failed to read error response: %w", err)
	}
	if resp.StatusCode != 200 {
		return audio.AudioResponse{}, errors.New(string(msg))
	}

	var audioResponse audio.AudioResponse
	err = json.Unmarshal(msg, &audioResponse)
	if err != nil {
		return audio.AudioResponse{}, fmt.Errorf("failed to parse speech-to-text response: %w", err)
	}

	return audioResponse, nil
}

// Embeddings retrieves text embeddings from the LiteLLM service.
func (l *Litellm) Embeddings(ctx context.Context, model models.ModelMeta, inputText string) (response.EmbeddingResponse, error) {
	if inputText == "" {
		return response.EmbeddingResponse{}, fmt.Errorf("inputText cannot be empty")
	}

	req := request.EmbeddingRequest{
		Model: string(model.ModelId),
		Input: inputText,
	}

	target := l.Connection.Targets.Get(cfg.CLIENT_LLM)
	resp, err := l.client(cfg.CLIENT_LLM).
		POST("/v1/embeddings").
		Context().Set(ctx).
		Header().AddAccept(mime.JSON).
		Retry().SetExponentialBackoff(
		target.RetryInterval,
		target.RetryMaxAttempts,
		target.RetryBackoffRate).
		Body().AsJSON(req).
		Send()

	if err != nil {
		return response.EmbeddingResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body().Close()

	var res response.EmbeddingResponse
	err = httpresp.ParseHTTPResponse(*resp, &res)
	if err != nil {
		return response.EmbeddingResponse{}, err
	}

	return res, nil
}

func (l *Litellm) TokenCounter(ctx context.Context, req *request.TokenCounterRequest) (*response.TokenCounterResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("TokenCounter request cannot be nil")
	}

	target := l.Connection.Targets.Get(cfg.CLIENT_LLM)
	resp, err := l.client(cfg.CLIENT_LLM).
		POST("/utils/token_counter").
		Context().Set(ctx).
		Header().AddAccept(mime.JSON).
		Retry().SetExponentialBackoff(
		target.RetryInterval,
		target.RetryMaxAttempts,
		target.RetryBackoffRate).
		Body().AsJSON(req).
		Send()

	if err != nil {
		return nil, fmt.Errorf("failed to send token counter request: %w", err)
	}
	defer resp.Body().Close()

	var res response.TokenCounterResponse
	err = httpresp.ParseHTTPResponse(*resp, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token counter response: %w", err)
	}

	return &res, nil
}
