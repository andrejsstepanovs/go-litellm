package common

import (
	"encoding/json"
	"fmt"
)

type ToolCallFunction struct {
	Name      string    `json:"name"`
	Arguments Arguments `json:"arguments"`
}

func (t *ToolCallFunction) UnmarshalJSON(data []byte) error {
	type Alias ToolCallFunction
	aux := &struct {
		Arguments string `json:"arguments"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("error unmarshalling ToolCallFunction: %w", err)
	}

	if aux.Arguments != "" {
		var args Arguments
		if err := json.Unmarshal([]byte(aux.Arguments), &args); err != nil {
			return fmt.Errorf("error unmarshalling arguments: %w", err)
		}
		t.Arguments = args
	}

	return nil
}

func (t *ToolCallFunction) MarshalJSON() ([]byte, error) {
	args, err := json.Marshal(t.Arguments)
	if err != nil {
		return nil, fmt.Errorf("error marshalling arguments: %w", err)
	}

	result := map[string]interface{}{
		"name":      t.Name,
		"arguments": string(args),
	}

	return json.Marshal(result)
}
