package connections

import (
	"fmt"

	"github.com/andrejsstepanovs/go-litellm/pkg/conf/connections/litellm"
)

type Config struct {
	LiteLLM litellm.Connection
}

func (a *Config) Validate() error {
	var errs []error

	err := a.LiteLLM.Validate()
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		return nil
	}

	var finalErr error
	for _, e := range errs {
		finalErr = fmt.Errorf("%w err: %w", finalErr, e)
	}

	return finalErr
}

func New() (Config, error) {
	llm, err := litellm.New()
	if err != nil {
		return Config{}, fmt.Errorf("litellm connection error: %w", err)
	}

	return Config{
		LiteLLM: llm,
	}, nil
}
