package models

type ModelID string

type Model struct {
	ID      ModelID `json:"id"`
	Object  string  `json:"object"`
	OwnedBy string  `json:"owned_by"`
}

type Models []Model

type ModelMeta struct {
	ModelId                         ModelID  `json:"model_group"`
	MaxInputTokens                  float64  `json:"max_input_tokens"`
	MaxOutputTokens                 float64  `json:"max_output_tokens"`
	InputCostPerToken               float64  `json:"input_cost_per_token"`
	OutputCostPerToken              float64  `json:"output_cost_per_token"`
	Providers                       []string `json:"providers"`
	Mode                            string   `json:"mode"`
	RPM                             int      `json:"rpm"`
	TPM                             int      `json:"tpm"`
	SupportsVision                  bool     `json:"supports_vision"`
	SupportsWebSearch               bool     `json:"supports_web_search"`
	SupportsReasoning               bool     `json:"supports_reasoning"`
	SupportsFunctionCalling         bool     `json:"supports_function_calling"`
	SupportsParallelFunctionCalling bool     `json:"supports_parallel_function_calling"`
	SupportedOpenAIParams           []string `json:"supported_openai_params"`
}

func (m Models) Get(id ModelID) (Model, bool) {
	for _, model := range m {
		if model.ID == id {
			return model, true
		}
	}
	return Model{}, false
}
