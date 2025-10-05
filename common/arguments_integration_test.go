package common_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/common"
)

// Test that demonstrates Arguments preserves types correctly for MCP tool calls
func Test_Arguments_PreservesTypes_Integration(t *testing.T) {
	// This simulates the exact scenario from the user's error:
	// The AI returns: {"uid": "1_146", "dblClick": true}
	// MCP expects: dblClick as boolean, not string

	jsonStr := `{
		"name": "click",
		"arguments": "{\"uid\": \"1_146\", \"dblClick\": true}"
	}`

	var toolFunc common.ToolCallFunction
	err := json.Unmarshal([]byte(jsonStr), &toolFunc)
	assert.NoError(t, err)

	// Verify dblClick is preserved as bool
	dblClick, ok := toolFunc.Arguments.GetArgument("dblClick")
	assert.True(t, ok)
	assert.IsType(t, true, dblClick, "dblClick should be boolean type")
	assert.Equal(t, true, dblClick)

	// Verify uid is preserved as string
	uid, ok := toolFunc.Arguments.GetArgument("uid")
	assert.True(t, ok)
	assert.IsType(t, "", uid, "uid should be string type")
	assert.Equal(t, "1_146", uid)

	// When we marshal this back (for MCP call), types should be preserved
	argsJSON, err := json.Marshal(toolFunc.Arguments)
	assert.NoError(t, err)

	// Unmarshal to verify structure
	var argsMap map[string]interface{}
	err = json.Unmarshal(argsJSON, &argsMap)
	assert.NoError(t, err)

	// Verify types are correct after round-trip
	assert.Equal(t, "1_146", argsMap["uid"])
	assert.Equal(t, true, argsMap["dblClick"])
	assert.IsType(t, true, argsMap["dblClick"], "dblClick must remain boolean after marshal/unmarshal")
}

// Test various types to ensure they're all preserved correctly
func Test_Arguments_AllTypes_Integration(t *testing.T) {
	jsonStr := `{
		"name": "test_tool",
		"arguments": "{\"str\": \"hello\", \"int\": 42, \"float\": 3.14, \"bool\": true, \"null\": null}"
	}`

	var toolFunc common.ToolCallFunction
	err := json.Unmarshal([]byte(jsonStr), &toolFunc)
	assert.NoError(t, err)

	// String
	val, ok := toolFunc.Arguments.GetArgument("str")
	assert.True(t, ok)
	assert.Equal(t, "hello", val)
	assert.IsType(t, "", val)

	// Integer (unmarshaled as float64 in JSON)
	val, ok = toolFunc.Arguments.GetArgument("int")
	assert.True(t, ok)
	assert.Equal(t, float64(42), val)
	assert.IsType(t, float64(0), val)

	// Float
	val, ok = toolFunc.Arguments.GetArgument("float")
	assert.True(t, ok)
	assert.Equal(t, 3.14, val)
	assert.IsType(t, float64(0), val)

	// Boolean
	val, ok = toolFunc.Arguments.GetArgument("bool")
	assert.True(t, ok)
	assert.Equal(t, true, val)
	assert.IsType(t, true, val)

	// Null
	val, ok = toolFunc.Arguments.GetArgument("null")
	assert.True(t, ok)
	assert.Nil(t, val)

	// Test GetStrArgument still works for all types
	strVal, ok := toolFunc.Arguments.GetStrArgument("int")
	assert.True(t, ok)
	assert.Equal(t, "42", strVal)

	strVal, ok = toolFunc.Arguments.GetStrArgument("bool")
	assert.True(t, ok)
	assert.Equal(t, "true", strVal)
}
