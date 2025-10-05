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

	v, ok := f.Arguments.GetArgument("timeout")
	assert.True(t, ok)
	assert.Equal(t, int(5000), v)

	v, ok = f.Arguments.GetArgument("text")
	assert.True(t, ok)
	assert.Equal(t, "WordPress", v)

	v, ok = f.Arguments.GetArgument("isTrue")
	assert.True(t, ok)
	assert.Equal(t, true, v)

	timeoutVal, ok := f.Arguments.GetArgument("timeout")
	assert.True(t, ok)
	assert.Equal(t, int(5000), timeoutVal)

	isTrueVal, ok := f.Arguments.GetArgument("isTrue")
	assert.True(t, ok)
	assert.Equal(t, true, isTrueVal)
}

func Test_ToolCallFunction_BooleanArgument(t *testing.T) {
	js := `{
		"name": "click",
		"arguments": "{\"uid\": \"1_146\", \"dblClick\": true}"
	}`

	var f common.ToolCallFunction
	err := json.Unmarshal([]byte(js), &f)
	assert.NoError(t, err)
	assert.Equal(t, "click", f.Name)

	dblClickVal, ok := f.Arguments.GetArgument("dblClick")
	assert.True(t, ok)
	assert.Equal(t, true, dblClickVal)
	assert.IsType(t, true, dblClickVal)

	uidVal, ok := f.Arguments.GetArgument("uid")
	assert.True(t, ok)
	assert.Equal(t, "1_146", uidVal)
	assert.IsType(t, "", uidVal)
}
