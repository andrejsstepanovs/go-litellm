package common

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Arguments map[string]string

// UnmarshalJSON implements custom unmarshaling for Arguments
// It converts numeric values to strings automatically
func (a *Arguments) UnmarshalJSON(data []byte) error {
	// First unmarshal into a map[string]interface{} to handle mixed types
	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return fmt.Errorf("error unmarshalling arguments: %w", err)
	}

	// Initialize the map if it's nil
	if *a == nil {
		*a = make(Arguments)
	}

	// Convert all values to strings
	for key, value := range rawMap {
		switch v := value.(type) {
		case string:
			(*a)[key] = v
		case float64:
			// JSON numbers are unmarshaled as float64
			// Check if it's actually an integer
			if v == float64(int64(v)) {
				(*a)[key] = strconv.FormatInt(int64(v), 10)
			} else {
				(*a)[key] = strconv.FormatFloat(v, 'f', -1, 64)
			}
		case bool:
			(*a)[key] = strconv.FormatBool(v)
		case nil:
			(*a)[key] = ""
		default:
			// For any other type, convert to JSON string representation
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("error marshalling value for key %s: %w", key, err)
			}
			(*a)[key] = string(jsonBytes)
		}
	}

	return nil
}

func (t *Arguments) SetStrArgument(name, value string) {
	(*t)[name] = value
}

func (t *Arguments) GetStrArgument(name string) (string, bool) {
	if !t.HasKey(name) {
		return "", false
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
