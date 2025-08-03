package json_schema

import (
	"fmt"
)

// JSONSchema represents the top-level JSON schema structure for LiteLLM API
type JSONSchema struct {
	Name   string                 `json:"name"`
	Schema map[string]interface{} `json:"schema"`
	Strict bool                   `json:"strict"`
}

// PropertyType represents the supported JSON schema types
type PropertyType string

const (
	TypeObject  PropertyType = "object"
	TypeArray   PropertyType = "array"
	TypeString  PropertyType = "string"
	TypeInteger PropertyType = "integer"
	TypeNumber  PropertyType = "number"
	TypeBoolean PropertyType = "boolean"
	TypeNull    PropertyType = "null"
)

// Property represents a single property in the JSON schema
type Property struct {
	Type        PropertyType        `json:"type"`
	Description string              `json:"description,omitempty"`
	Properties  map[string]Property `json:"properties,omitempty"`
	Items       *Property           `json:"items,omitempty"`
	Required    []string            `json:"required,omitempty"`
	Enum        []interface{}       `json:"enum,omitempty"`
	Default     interface{}         `json:"default,omitempty"`
}

// SchemaBuilder provides a fluent interface for building JSON schemas
type SchemaBuilder struct {
	name       string
	properties map[string]Property
	required   []string
	strict     bool
}

// NewSchemaBuilder creates a new schema builder with the given name
func NewSchemaBuilder(name string) *SchemaBuilder {
	return &SchemaBuilder{
		name:       name,
		properties: make(map[string]Property),
		required:   make([]string, 0),
		strict:     true,
	}
}

// AddProperty adds a property to the schema
func (sb *SchemaBuilder) AddProperty(name string, property Property) *SchemaBuilder {
	sb.properties[name] = property
	return sb
}

// AddStringProperty adds a string property to the schema
func (sb *SchemaBuilder) AddStringProperty(name, description string) *SchemaBuilder {
	sb.properties[name] = Property{
		Type:        TypeString,
		Description: description,
	}
	return sb
}

// AddIntegerProperty adds an integer property to the schema
func (sb *SchemaBuilder) AddIntegerProperty(name, description string) *SchemaBuilder {
	sb.properties[name] = Property{
		Type:        TypeInteger,
		Description: description,
	}
	return sb
}

// AddNumberProperty adds a number property to the schema
func (sb *SchemaBuilder) AddNumberProperty(name, description string) *SchemaBuilder {
	sb.properties[name] = Property{
		Type:        TypeNumber,
		Description: description,
	}
	return sb
}

// AddBooleanProperty adds a boolean property to the schema
func (sb *SchemaBuilder) AddBooleanProperty(name, description string) *SchemaBuilder {
	sb.properties[name] = Property{
		Type:        TypeBoolean,
		Description: description,
	}
	return sb
}

// AddObjectProperty adds an object property to the schema
func (sb *SchemaBuilder) AddObjectProperty(name, description string, properties map[string]Property, required []string) *SchemaBuilder {
	sb.properties[name] = Property{
		Type:        TypeObject,
		Description: description,
		Properties:  properties,
		Required:    required,
	}
	return sb
}

// AddArrayProperty adds an array property to the schema
func (sb *SchemaBuilder) AddArrayProperty(name, description string, items Property) *SchemaBuilder {
	sb.properties[name] = Property{
		Type:        TypeArray,
		Description: description,
		Items:       &items,
	}
	return sb
}

// SetRequired sets the required fields for the schema
func (sb *SchemaBuilder) SetRequired(required []string) *SchemaBuilder {
	sb.required = required
	return sb
}

// AddRequired adds a field to the required list
func (sb *SchemaBuilder) AddRequired(field string) *SchemaBuilder {
	sb.required = append(sb.required, field)
	return sb
}

// SetStrict sets the strict mode for the schema
func (sb *SchemaBuilder) SetStrict(strict bool) *SchemaBuilder {
	sb.strict = strict
	return sb
}

// Build constructs the final JSONSchema
func (sb *SchemaBuilder) Build() (*JSONSchema, error) {
	if sb.name == "" {
		return nil, fmt.Errorf("schema name is required")
	}

	if len(sb.properties) == 0 {
		return nil, fmt.Errorf("schema must have at least one property")
	}

	schema := map[string]interface{}{
		"type":                 "object",
		"properties":           sb.properties,
		"required":             sb.required,
		"additionalProperties": false,
	}

	return &JSONSchema{
		Name:   sb.name,
		Schema: schema,
		Strict: sb.strict,
	}, nil
}

// BuildFromMapping creates a JSON schema from a mapping structure
func BuildFromMapping(name string, mapping map[string]interface{}) (*JSONSchema, error) {
	if name == "" {
		return nil, fmt.Errorf("schema name is required")
	}

	if len(mapping) == 0 {
		return nil, fmt.Errorf("mapping cannot be empty")
	}

	builder := NewSchemaBuilder(name)

	for fieldName, fieldDef := range mapping {
		property, err := buildPropertyFromDefinition(fieldDef)
		if err != nil {
			return nil, fmt.Errorf("failed to build property %s: %w", fieldName, err)
		}
		builder.AddProperty(fieldName, property)
	}

	return builder.Build()
}

// buildPropertyFromDefinition converts a field definition to a Property
func buildPropertyFromDefinition(def interface{}) (Property, error) {
	switch v := def.(type) {
	case string:
		return buildPropertyFromString(v)
	case map[string]interface{}:
		return buildPropertyFromMap(v)
	default:
		return Property{}, fmt.Errorf("unsupported definition type: %T", def)
	}
}

// buildPropertyFromString creates a property from a string type definition
func buildPropertyFromString(typeDef string) (Property, error) {
	switch typeDef {
	case "string":
		return Property{Type: TypeString}, nil
	case "integer":
		return Property{Type: TypeInteger}, nil
	case "number":
		return Property{Type: TypeNumber}, nil
	case "boolean":
		return Property{Type: TypeBoolean}, nil
	case "null":
		return Property{Type: TypeNull}, nil
	default:
		return Property{}, fmt.Errorf("unsupported type: %s", typeDef)
	}
}

// buildPropertyFromMap creates a property from a map definition
func buildPropertyFromMap(def map[string]interface{}) (Property, error) {
	property := Property{}

	// Extract type
	if typeVal, ok := def["type"]; ok {
		if typeStr, ok := typeVal.(string); ok {
			property.Type = PropertyType(typeStr)
		} else {
			return Property{}, fmt.Errorf("type must be a string")
		}
	} else {
		return Property{}, fmt.Errorf("type is required")
	}

	// Extract description
	if desc, ok := def["description"]; ok {
		if descStr, ok := desc.(string); ok {
			property.Description = descStr
		}
	}

	// Handle object type
	if property.Type == TypeObject {
		if props, ok := def["properties"]; ok {
			if propsMap, ok := props.(map[string]interface{}); ok {
				properties := make(map[string]Property)
				for propName, propDef := range propsMap {
					prop, err := buildPropertyFromDefinition(propDef)
					if err != nil {
						return Property{}, fmt.Errorf("failed to build nested property %s: %w", propName, err)
					}
					properties[propName] = prop
				}
				property.Properties = properties
			}
		}

		if req, ok := def["required"]; ok {
			if reqSlice, ok := req.([]interface{}); ok {
				required := make([]string, len(reqSlice))
				for i, r := range reqSlice {
					if reqStr, ok := r.(string); ok {
						required[i] = reqStr
					} else {
						return Property{}, fmt.Errorf("required field names must be strings")
					}
				}
				property.Required = required
			}
		}
	}

	// Handle array type
	if property.Type == TypeArray {
		if items, ok := def["items"]; ok {
			itemProp, err := buildPropertyFromDefinition(items)
			if err != nil {
				return Property{}, fmt.Errorf("failed to build array items: %w", err)
			}
			property.Items = &itemProp
		}
	}

	// Handle enum
	if enum, ok := def["enum"]; ok {
		if enumSlice, ok := enum.([]interface{}); ok {
			property.Enum = enumSlice
		}
	}

	// Handle default
	if defaultVal, ok := def["default"]; ok {
		property.Default = defaultVal
	}

	return property, nil
}
