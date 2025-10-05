package common

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArguments(t *testing.T) {
	t.Run("string types are preserved", func(t *testing.T) {
		const file = "testdata/args.json"
		dat, err := os.ReadFile(file)
		assert.NoError(t, err)

		toolCall := ToolCalls{}
		err = json.Unmarshal(dat, &toolCall)
		assert.NoError(t, err)

		assert.Len(t, toolCall, 1)
		args := toolCall[0].Function.Arguments
		val, exists := args.GetArgument("country")
		assert.True(t, exists)
		assert.Equal(t, "Latvia", val)
	})

	t.Run("int argument is preserved as string", func(t *testing.T) {
		const file = "testdata/args_int.json"
		dat, err := os.ReadFile(file)
		assert.NoError(t, err)

		toolCall := ToolCalls{}
		err = json.Unmarshal(dat, &toolCall)
		assert.NoError(t, err)

		assert.Len(t, toolCall, 1)
		args := toolCall[0].Function.Arguments
		val, exists := args.GetStrArgument("pageIdx")
		assert.True(t, exists)
		assert.Equal(t, "4", val)
	})

	t.Run("bool argument is preserved as bool", func(t *testing.T) {
		const file = "testdata/args_bool.json"
		dat, err := os.ReadFile(file)
		assert.NoError(t, err)

		toolCall := ToolCalls{}
		err = json.Unmarshal(dat, &toolCall)
		assert.NoError(t, err)

		assert.Len(t, toolCall, 1)
		args := toolCall[0].Function.Arguments
		val, exists := args.GetArgument("isTrue")
		assert.True(t, exists)
		assert.Equal(t, true, val)
	})
}
