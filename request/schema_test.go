package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSchemaBuilder(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "creates builder with name",
			input:    "test_schema",
			expected: "test_schema",
		},
		{
			name:     "creates builder with empty name",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSchemaBuilder(tt.input)
			assert.Equal(t, tt.expected, builder.name)
			assert.NotNil(t, builder.properties)
			assert.NotNil(t, builder.required)
			assert.True(t, builder.strict)
		})
	}
}

func TestSchemaBuilder_AddStringProperty(t *testing.T) {
	tests := []struct {
		name        string
		propName    string
		description string
	}{
		{
			name:        "adds string property",
			propName:    "username",
			description: "User's name",
		},
		{
			name:        "adds string property without description",
			propName:    "email",
			description: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSchemaBuilder("test")
			result := builder.AddStringProperty(tt.propName, tt.description)

			assert.Equal(t, builder, result) // AI: Check fluent interface
			assert.Contains(t, builder.properties, tt.propName)
			assert.Equal(t, TypeString, builder.properties[tt.propName].Type)
			assert.Equal(t, tt.description, builder.properties[tt.propName].Description)
		})
	}
}

func TestSchemaBuilder_AddIntegerProperty(t *testing.T) {
	builder := NewSchemaBuilder("test")
	result := builder.AddIntegerProperty("age", "User's age")

	assert.Equal(t, builder, result)
	assert.Contains(t, builder.properties, "age")
	assert.Equal(t, TypeInteger, builder.properties["age"].Type)
	assert.Equal(t, "User's age", builder.properties["age"].Description)
}

func TestSchemaBuilder_AddNumberProperty(t *testing.T) {
	builder := NewSchemaBuilder("test")
	result := builder.AddNumberProperty("price", "Product price")

	assert.Equal(t, builder, result)
	assert.Contains(t, builder.properties, "price")
	assert.Equal(t, TypeNumber, builder.properties["price"].Type)
	assert.Equal(t, "Product price", builder.properties["price"].Description)
}

func TestSchemaBuilder_AddBooleanProperty(t *testing.T) {
	builder := NewSchemaBuilder("test")
	result := builder.AddBooleanProperty("active", "Is active")

	assert.Equal(t, builder, result)
	assert.Contains(t, builder.properties, "active")
	assert.Equal(t, TypeBoolean, builder.properties["active"].Type)
	assert.Equal(t, "Is active", builder.properties["active"].Description)
}

func TestSchemaBuilder_AddObjectProperty(t *testing.T) {
	properties := map[string]Property{
		"street": {Type: TypeString, Description: "Street name"},
		"city":   {Type: TypeString, Description: "City name"},
	}
	required := []string{"street", "city"}

	builder := NewSchemaBuilder("test")
	result := builder.AddObjectProperty("address", "User address", properties, required)

	assert.Equal(t, builder, result)
	assert.Contains(t, builder.properties, "address")
	assert.Equal(t, TypeObject, builder.properties["address"].Type)
	assert.Equal(t, "User address", builder.properties["address"].Description)
	assert.Equal(t, properties, builder.properties["address"].Properties)
	assert.Equal(t, required, builder.properties["address"].Required)
}

func TestSchemaBuilder_AddArrayProperty(t *testing.T) {
	items := Property{Type: TypeString, Description: "Tag name"}

	builder := NewSchemaBuilder("test")
	result := builder.AddArrayProperty("tags", "User tags", items)

	assert.Equal(t, builder, result)
	assert.Contains(t, builder.properties, "tags")
	assert.Equal(t, TypeArray, builder.properties["tags"].Type)
	assert.Equal(t, "User tags", builder.properties["tags"].Description)
	assert.Equal(t, &items, builder.properties["tags"].Items)
}

func TestSchemaBuilder_SetRequired(t *testing.T) {
	required := []string{"name", "email"}

	builder := NewSchemaBuilder("test")
	result := builder.SetRequired(required)

	assert.Equal(t, builder, result)
	assert.Equal(t, required, builder.required)
}

func TestSchemaBuilder_AddRequired(t *testing.T) {
	builder := NewSchemaBuilder("test")
	result := builder.AddRequired("name").AddRequired("email")

	assert.Equal(t, builder, result)
	assert.Equal(t, []string{"name", "email"}, builder.required)
}

func TestSchemaBuilder_SetStrict(t *testing.T) {
	tests := []struct {
		name   string
		strict bool
	}{
		{
			name:   "sets strict to true",
			strict: true,
		},
		{
			name:   "sets strict to false",
			strict: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSchemaBuilder("test")
			result := builder.SetStrict(tt.strict)

			assert.Equal(t, builder, result)
			assert.Equal(t, tt.strict, builder.strict)
		})
	}
}

func TestSchemaBuilder_Build(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() *SchemaBuilder
		expectError bool
		errorMsg    string
	}{
		{
			name: "builds valid schema",
			setupFunc: func() *SchemaBuilder {
				return NewSchemaBuilder("user_schema").
					AddStringProperty("name", "User name").
					AddIntegerProperty("age", "User age").
					SetRequired([]string{"name"})
			},
			expectError: false,
		},
		{
			name: "fails with empty name",
			setupFunc: func() *SchemaBuilder {
				return NewSchemaBuilder("").
					AddStringProperty("name", "User name")
			},
			expectError: true,
			errorMsg:    "schema name is required",
		},
		{
			name: "fails with no properties",
			setupFunc: func() *SchemaBuilder {
				return NewSchemaBuilder("test_schema")
			},
			expectError: true,
			errorMsg:    "schema must have at least one property",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.setupFunc()
			schema, err := builder.Build()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, schema)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, schema)
				assert.Equal(t, builder.name, schema.Name)
				assert.Equal(t, builder.strict, schema.Strict)
				assert.Contains(t, schema.Schema, "type")
				assert.Contains(t, schema.Schema, "properties")
				assert.Contains(t, schema.Schema, "required")
				assert.Contains(t, schema.Schema, "additionalProperties")
				assert.Equal(t, "object", schema.Schema["type"])
				assert.Equal(t, false, schema.Schema["additionalProperties"])
			}
		})
	}
}

func TestBuildFromMapping(t *testing.T) {
	tests := []struct {
		name        string
		schemaName  string
		mapping     map[string]interface{}
		expectError bool
		errorMsg    string
		validate    func(t *testing.T, schema *JSONSchema)
	}{
		{
			name:       "builds schema from simple string types",
			schemaName: "user",
			mapping: map[string]interface{}{
				"name":  "string",
				"email": "string",
				"age":   "integer",
			},
			expectError: false,
			validate: func(t *testing.T, schema *JSONSchema) {
				assert.Equal(t, "user", schema.Name)
				assert.True(t, schema.Strict)

				properties := schema.Schema["properties"].(map[string]Property)
				assert.Equal(t, TypeString, properties["name"].Type)
				assert.Equal(t, TypeString, properties["email"].Type)
				assert.Equal(t, TypeInteger, properties["age"].Type)
			},
		},
		{
			name:       "builds schema from detailed definitions",
			schemaName: "product",
			mapping: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Product name",
				},
				"price": map[string]interface{}{
					"type":        "number",
					"description": "Product price",
				},
			},
			expectError: false,
			validate: func(t *testing.T, schema *JSONSchema) {
				properties := schema.Schema["properties"].(map[string]Property)
				assert.Equal(t, TypeString, properties["name"].Type)
				assert.Equal(t, "Product name", properties["name"].Description)
				assert.Equal(t, TypeNumber, properties["price"].Type)
				assert.Equal(t, "Product price", properties["price"].Description)
			},
		},
		{
			name:       "builds schema with nested object",
			schemaName: "user",
			mapping: map[string]interface{}{
				"name": "string",
				"address": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"street": "string",
						"city":   "string",
					},
					"required": []interface{}{"street", "city"},
				},
			},
			expectError: false,
			validate: func(t *testing.T, schema *JSONSchema) {
				properties := schema.Schema["properties"].(map[string]Property)
				assert.Equal(t, TypeObject, properties["address"].Type)
				assert.Contains(t, properties["address"].Properties, "street")
				assert.Contains(t, properties["address"].Properties, "city")
				assert.Equal(t, []string{"street", "city"}, properties["address"].Required)
			},
		},
		{
			name:       "builds schema with array",
			schemaName: "user",
			mapping: map[string]interface{}{
				"name": "string",
				"tags": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
			expectError: false,
			validate: func(t *testing.T, schema *JSONSchema) {
				properties := schema.Schema["properties"].(map[string]Property)
				assert.Equal(t, TypeArray, properties["tags"].Type)
				assert.NotNil(t, properties["tags"].Items)
				assert.Equal(t, TypeString, properties["tags"].Items.Type)
			},
		},
		{
			name:        "fails with empty name",
			schemaName:  "",
			mapping:     map[string]interface{}{"name": "string"},
			expectError: true,
			errorMsg:    "schema name is required",
		},
		{
			name:        "fails with empty mapping",
			schemaName:  "test",
			mapping:     map[string]interface{}{},
			expectError: true,
			errorMsg:    "mapping cannot be empty",
		},
		{
			name:       "fails with unsupported type",
			schemaName: "test",
			mapping: map[string]interface{}{
				"field": "unsupported_type",
			},
			expectError: true,
			errorMsg:    "unsupported type: unsupported_type",
		},
		{
			name:       "fails with invalid definition type",
			schemaName: "test",
			mapping: map[string]interface{}{
				"field": 123,
			},
			expectError: true,
			errorMsg:    "unsupported definition type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema, err := BuildFromMapping(tt.schemaName, tt.mapping)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, schema)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, schema)
				if tt.validate != nil {
					tt.validate(t, schema)
				}
			}
		})
	}
}

func TestBuildPropertyFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    PropertyType
		expectError bool
	}{
		{
			name:        "string type",
			input:       "string",
			expected:    TypeString,
			expectError: false,
		},
		{
			name:        "integer type",
			input:       "integer",
			expected:    TypeInteger,
			expectError: false,
		},
		{
			name:        "number type",
			input:       "number",
			expected:    TypeNumber,
			expectError: false,
		},
		{
			name:        "boolean type",
			input:       "boolean",
			expected:    TypeBoolean,
			expectError: false,
		},
		{
			name:        "null type",
			input:       "null",
			expected:    TypeNull,
			expectError: false,
		},
		{
			name:        "unsupported type",
			input:       "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			property, err := buildPropertyFromString(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, property.Type)
			}
		})
	}
}

func TestBuildPropertyFromMap(t *testing.T) {
	tests := []struct {
		name        string
		input       map[string]interface{}
		expectError bool
		errorMsg    string
		validate    func(t *testing.T, property Property)
	}{
		{
			name: "simple string property",
			input: map[string]interface{}{
				"type":        "string",
				"description": "A string field",
			},
			expectError: false,
			validate: func(t *testing.T, property Property) {
				assert.Equal(t, TypeString, property.Type)
				assert.Equal(t, "A string field", property.Description)
			},
		},
		{
			name: "property with enum",
			input: map[string]interface{}{
				"type": "string",
				"enum": []interface{}{"option1", "option2", "option3"},
			},
			expectError: false,
			validate: func(t *testing.T, property Property) {
				assert.Equal(t, TypeString, property.Type)
				assert.Equal(t, []interface{}{"option1", "option2", "option3"}, property.Enum)
			},
		},
		{
			name: "property with default value",
			input: map[string]interface{}{
				"type":    "string",
				"default": "default_value",
			},
			expectError: false,
			validate: func(t *testing.T, property Property) {
				assert.Equal(t, TypeString, property.Type)
				assert.Equal(t, "default_value", property.Default)
			},
		},
		{
			name: "missing type",
			input: map[string]interface{}{
				"description": "A field without type",
			},
			expectError: true,
			errorMsg:    "type is required",
		},
		{
			name: "invalid type format",
			input: map[string]interface{}{
				"type": 123,
			},
			expectError: true,
			errorMsg:    "type must be a string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			property, err := buildPropertyFromMap(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, property)
				}
			}
		})
	}
}
