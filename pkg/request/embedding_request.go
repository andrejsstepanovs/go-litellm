package request

type EmbeddingRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}
