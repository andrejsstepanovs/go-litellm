package common_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/common"
)

func Test_ToolCallFunction_Unmarshal_MixedTypes(t *testing.T) {
	const file = "testdata/tool_calls.json"
	dat, err := os.ReadFile(file)
	assert.NoError(t, err)

	var f common.ToolCallFunction
	err = json.Unmarshal(dat, &f)
	assert.NoError(t, err)
	assert.Equal(t, "wait_for", f.Name)

	v, ok := f.Arguments.GetStrArgument("timeout")
	assert.True(t, ok)
	assert.Equal(t, "5000", v)

	v, ok = f.Arguments.GetStrArgument("text")
	assert.True(t, ok)
	assert.Equal(t, "WordPress", v)

	v, ok = f.Arguments.GetStrArgument("timeout")
	assert.True(t, ok)
	assert.Equal(t, "5000", v)

	v, ok = f.Arguments.GetStrArgument("isTrue")
	assert.True(t, ok)
	assert.Equal(t, "true", v)
}
