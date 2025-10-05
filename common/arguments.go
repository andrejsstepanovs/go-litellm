package common

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Arguments map[string]any

// UnmarshalJSON implements custom unmarshaling for Arguments
// Preserves original types (string, number, bool, etc.)
func (a *Arguments) UnmarshalJSON(data []byte) error {
	// Unmarshal into a map[string]interface{} to preserve types
	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return fmt.Errorf("error unmarshalling arguments: %w", err)
	}

	// Initialize the map if it's nil
	if *a == nil {
		*a = make(Arguments)
	}

	// Copy values preserving their types
	for key, value := range rawMap {
		(*a)[key] = value
	}

	return nil
}

func (t *Arguments) SetStrArgument(name, value string) {
	(*t)[name] = value
}

func (t *Arguments) SetArgument(name string, value any) {
	(*t)[name] = value
}

func (t *Arguments) GetStrArgument(name string) (string, bool) {
	if !t.HasKey(name) {
		return "", false
	}
	value := (*t)[name]

	// Handle different types, converting to string
	switch v := value.(type) {
	case string:
		return v, true
	case float64:
		// Check if it's actually an integer
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10), true
		}
		return strconv.FormatFloat(v, 'f', -1, 64), true
	case bool:
		return strconv.FormatBool(v), true
	case nil:
		return "", true
	default:
		// For any other type, convert to JSON string representation
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return "", false
		}
		return string(jsonBytes), true
	}
}

func (t *Arguments) GetArgument(name string) (any, bool) {
	if !t.HasKey(name) {
		return nil, false
	}
	value := (*t)[name]
	return value, true
}

func (t *Arguments) HasKey(name string) bool {
	if t == nil {
		return false
	}
	_, ok := (*t)[name]
	return ok
}
