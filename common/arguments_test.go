package common

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgumentsString(t *testing.T) {
	const file = "testdata/args.json"
	dat, err := os.ReadFile(file)
	assert.NoError(t, err)

	toolCall := ToolCalls{}
	err = json.Unmarshal(dat, &toolCall)
	assert.NoError(t, err)

	// Verify that the string value works as expected
	assert.Len(t, toolCall, 1)
	args := toolCall[0].Function.Arguments
	country, exists := args.GetStrArgument("country")
	assert.True(t, exists)
	assert.Equal(t, "Latvia", country)
}

func TestParseArgumentsInt(t *testing.T) {
	const file = "testdata/args_int.json"
	dat, err := os.ReadFile(file)
	assert.NoError(t, err)

	toolCall := ToolCalls{}
	err = json.Unmarshal(dat, &toolCall)
	assert.NoError(t, err)

	// Verify that the integer value was converted to string
	assert.Len(t, toolCall, 1)
	args := toolCall[0].Function.Arguments
	pageIdx, exists := args.GetStrArgument("pageIdx")
	assert.True(t, exists)
	assert.Equal(t, "4", pageIdx) // Should be string "4", not integer 4
}
