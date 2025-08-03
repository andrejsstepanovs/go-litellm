package litellm

import (
	"fmt"

	"github.com/spf13/viper"
)

type TargetName string

const (
	CLIENT_SYSTEM TargetName = "system"
	CLIENT_MCP    TargetName = "mcp"
	CLIENT_LLM    TargetName = "llm"
)

type Targets struct {
	System Target
	LLM    Target
	MCP    Target
}

func (t *Targets) Get(name TargetName) Target {
	switch name {
	case CLIENT_SYSTEM:
		return t.System
	case CLIENT_MCP:
		return t.MCP
	case CLIENT_LLM:
		return t.LLM
	}
	return Target{}
}

func (t *Targets) Validate() error {
	var errs []error

	err := t.System.Validate()
	if err != nil {
		errs = append(errs, fmt.Errorf("system target validation error: %w", err))
	}

	err = t.LLM.Validate()
	if err != nil {
		errs = append(errs, fmt.Errorf("llm target validation error: %w", err))
	}

	err = t.MCP.Validate()
	if err != nil {
		errs = append(errs, fmt.Errorf("mcp target validation error: %w", err))
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

func NewTargets() Targets {
	return Targets{
		LLM: Target{
			Timeout: viper.GetDuration("litellm.targets.llm.timeout"),
		},
		System: Target{
			Timeout: viper.GetDuration("litellm.targets.system.timeout"),
		},
		MCP: Target{
			Timeout: viper.GetDuration("litellm.targets.mcp.timeout"),
		},
	}
}
