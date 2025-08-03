package litellm

import (
	"fmt"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var validate = validator.New()

type Connection struct {
	URL     url.URL `validate:"required"`
	Targets Targets
}

func (a *Connection) Validate() error {
	if a == nil {
		return fmt.Errorf("litellm is required")
	}

	err := validate.Struct(a)
	if err != nil {
		return fmt.Errorf("litellm validation error: %w", err)
	}

	// Check if URL is not empty
	if a.URL.String() == "" {
		return fmt.Errorf("litellm validation error: url is empty")
	}

	if a.URL.Host == "" || a.URL.Scheme == "" {
		return fmt.Errorf("litellm validation error: url is invalid or empty")
	}

	err = a.Targets.Validate()
	if err != nil {
		return fmt.Errorf("litellm targets validation error: %w", err)
	}

	return nil
}

func New() (Connection, error) {
	baseURL, err := url.Parse(viper.GetString("litellm.url"))
	if err != nil || baseURL == nil {
		return Connection{}, fmt.Errorf("litellm base_url parse error: %w", err)
	}

	return Connection{
		URL:     *baseURL,
		Targets: NewTargets(),
	}, nil
}
