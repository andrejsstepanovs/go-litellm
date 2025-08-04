package request

import "github.com/andrejsstepanovs/go-litellm/mcp"

const FunctionToolType = "function"

// ``
//{
//  "model": "gpt-4o",
//  "messages": [
//    {"role": "user", "content": "What is the weather like in Paris today?"}
//  ],
//  "tools": [
//    {
//      "type": "function",
//      "function": {
//        "name": "get_weather",
//        "description": "Get current temperature for a given location.",
//        "parameters": {
//          "type": "object",
//          "properties": {
//            "location": {
//              "type": "string",
//              "description": "City and country e.g. Paris, France"
//            }
//          },
//          "required": ["location"]
//        }
//      }
//    }
//  ]
//}
//```
//
//
//
//----
//{
//  "model": "gpt-3.5-turbo-1106",
//  "messages": [
//    {"role": "user", "content": "What is the weather like in Boston?"}
//  ],
//  "tools": [
//    {
//      "type": "function",
//      "function": {
//        "name": "get_current_weather",
//        "description": "Get the current weather in a given location",
//        "parameters": {
//          "type": "object",
//          "properties": {
//            "location": {"type": "string", "description": "The city and state, e.g. San Francisco, CA"},
//            "unit": {"type": "string", "description": "Temperature unit", "enum": ["fahrenheit", "celsius"]}
//          },
//          "required": ["location", "unit"]
//        }
//      }
//    }
//  ]
//}

type LLMCallTools []LLMCallTool

type LLMCallTool struct {
	Type     string               `json:"type"`               // "function"
	Function *LLMCallToolFunction `json:"function,omitempty"` // {"name": "get_current_weather", "description": "Get the current weather in a given location", "parameters": {"type": "object", "properties": {"location": {"type": "string", "description": "The city and state, e.g. San Francisco, CA"}, "unit": {"type": "string", "description": "Temperature unit", "enum": ["fahrenheit", "celsius"]}}, "required": ["location", "unit"]}}
}

type LLMCallToolFunction struct {
	Name        string                         `json:"name,omitempty"`        // "get_current_weather"
	Description string                         `json:"description,omitempty"` // "Get the current weather in a given location"
	Parameters  *LLMCallToolFunctionParameters `json:"parameters,omitempty"`  // {"type": "object", "properties": {"location": {"type": "string", "description": "The city and state, e.g. San Francisco, CA"}, "unit": {"type": "string", "description": "Temperature unit", "enum": ["fahrenheit", "celsius"]}}, "required": ["location", "unit"]}
}

type LLMCallToolFunctionParameters struct {
	Type       string                                 `json:"type,omitempty"`       // "object"
	Properties map[string]LLMCallToolFunctionProperty `json:"properties,omitempty"` // {"location": {"type": "string", "description": "The city and state, e.g. San Francisco, CA"}, "unit": {"type": "string", "description": "Temperature unit", "enum": ["fahrenheit", "celsius"]}}
	Required   []string                               `json:"required,omitempty"`   // ["location", "unit"]
}

type LLMCallToolFunctionProperty struct {
	Description string   `json:"description,omitempty"` // "The city and state, e.g. San Francisco, CA"
	Type        string   `json:"type,omitempty"`        // "string"
	Enum        []string `json:"enum,omitempty"`        // ["fahrenheit", "celsius"] (optional, only if applicable)
	Format      string   `json:"format,omitempty"`      // "date-time" (optional, only if applicable)
	Example     string   `json:"example,omitempty"`     // "2023-10-01T12:00:00Z" (optional, only if applicable)
}

func ToLLMCallTools(availableTools mcp.AvailableTools) LLMCallTools {
	if len(availableTools) == 0 {
		return LLMCallTools{}
	}
	tools := make(LLMCallTools, 0)
	for _, tool := range availableTools {
		properties := make(map[string]LLMCallToolFunctionProperty)
		for propName, prop := range tool.InputSchema.Properties {
			properties[propName] = LLMCallToolFunctionProperty{
				Description: prop.Description,
				Type:        prop.Type,
			}
		}
		toolType := FunctionToolType
		if tool.Type != "" {
			toolType = tool.Type
		}

		addTool := LLMCallTool{Type: toolType}
		if tool.Name != "" {
			addTool.Function = &LLMCallToolFunction{
				Name:        tool.Name,
				Description: tool.Description,
			}
			if len(properties) > 0 {
				props := LLMCallToolFunctionParameters{
					Type:       "object",
					Properties: properties,
					Required:   tool.InputSchema.Required,
				}
				addTool.Function.Parameters = &props
			}
		}

		tools = append(tools, addTool)
	}

	return tools
}
