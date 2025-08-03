package litellm

import (
	"fmt"
	"time"
)

type Target struct {
	Timeout          time.Duration `validate:"required"`
	RetryInterval    time.Duration `validate:"required"`
	RetryMaxAttempts uint          `validate:"required"`
	RetryBackoffRate float64       `validate:"required"`
	MaxRetry         uint          `validate:"required"`
}

func (t *Target) Validate() error {
	if t == nil {
		return fmt.Errorf("target is required")
	}

	err := validate.Struct(t)
	if err != nil {
		return fmt.Errorf("target validation error: %w", err)
	}
	return nil
}
