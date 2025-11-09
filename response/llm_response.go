package response

import (
	"github.com/andrejsstepanovs/go-litellm/common"
	"github.com/andrejsstepanovs/go-litellm/models"
)

type FinishReasonType string

const FINISH_REASON_STOP FinishReasonType = "stop"
const FINISH_REASON_TOOL FinishReasonType = "tool_calls"

// https://platform.openai.com/docs/guides/function-calling?api-mode=responses

type ResponseChoice struct {
	Index        int              `json:"index"`
	FinishReason FinishReasonType `json:"finish_reason"`
	Message      ResponseMessage  `json:"message"`
}

type ResponseChoices []ResponseChoice

type ResponseMessage struct {
	Content   string           `json:"content"`
	Role      string           `json:"role"`
	ToolCalls common.ToolCalls `json:"tool_calls,omitempty"`
}

func (rm *ResponseMessage) IsEmpty() bool {
	return rm == nil || (rm.Content == "" && len(rm.ToolCalls) == 0)
}

type ResponseUsage struct {
	CompletionTokens        int                     `json:"completion_tokens"`
	PromptTokens            int                     `json:"prompt_tokens"`
	TotalTokens             int                     `json:"total_tokens"`
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details"`
	PromptTokensDetails     PromptTokensDetails     `json:"prompt_tokens_details"` // TODO NOT A STRING
	QueueTime               float64                 `json:"queue_time"`
	PromptTime              float64                 `json:"prompt_time"`
	CompletionTime          float64                 `json:"completion_time"`
	TotalTime               float64                 `json:"total_time"`
}

type CompletionTokensDetails struct {
	AudioTokens              int `json:"audio_tokens"`
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
	RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
	ReasoningTokens          int `json:"reasoning_tokens"`
}

type PromptTokensDetails struct {
	AudioTokens  int `json:"audio_tokens"`
	CachedTokens int `json:"cached_tokens"`
}

type Response struct {
	ID                string          `json:"id"`
	Created           int             `json:"created"`
	Model             models.ModelID  `json:"model"`
	Object            string          `json:"object"`
	SystemFingerprint string          `json:"system_fingerprint"`
	Choices           ResponseChoices `json:"choices"`
	Usage             ResponseUsage   `json:"usage"`
}

func (r *Response) Choice() ResponseChoice {
	if r == nil || len(r.Choices) == 0 {
		return ResponseChoice{}
	}

	return r.Choices[len(r.Choices)-1]
}

func (r *Response) Message() ResponseMessage {
	if r == nil {
		return ResponseMessage{}
	}
	return r.Choice().Message
}

func (r *Response) String() string {
	return r.Message().Content
}

func (r *Response) SetText(text string) {
	if len(r.Choices) == 0 {
		r.Choices = ResponseChoices{
			{
				FinishReason: FINISH_REASON_STOP,
				Message: ResponseMessage{
					Content: text,
					Role:    "assistant",
				},
			},
		}
	}

	r.Choices[len(r.Choices)-1].Message.Content = text
}

func (r *Response) Bytes() []byte {
	return []byte(r.String())
}

type TokenCounterResponse struct {
	TotalTokens   float64 `json:"total_tokens"`
	ModelUsed     string  `json:"model_used"`
	RequestModel  string  `json:"request_model"`
	TokenizerType string  `json:"tokenizer_type"`
}
