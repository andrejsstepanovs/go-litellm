package response

type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type Embedding []float64

type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding Embedding `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  EmbeddingUsage  `json:"usage"`
}

func (e *Embedding) Float32() []float32 {
	float32s := make([]float32, len(*e))
	for i, v := range *e {
		float32s[i] = float32(v)
	}
	return float32s
}
