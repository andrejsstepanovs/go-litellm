package mcp

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvailableTools(t *testing.T) {
	file := "testdata/response_200.json"
	contents, err := os.ReadFile(file)
	assert.NoError(t, err)
	var response AvailableToolsResponse

	err = json.Unmarshal(contents, &response)
	assert.NoError(t, err)

	assert.NotEmpty(t, response.Tools)
	assert.Equal(t, "Successfully retrieved tools", response.Message)
	assert.Nil(t, response.Error)

	assert.Equal(t, "mcp_bobik-calendar_events", response.Tools[0].Name)
}
