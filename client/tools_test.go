package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/client"
	"github.com/andrejsstepanovs/go-litellm/mcp"
)

func Test_Tools_Functional(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	t.Run("success", func(t *testing.T) {
		clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}
		res, err := clientInstance.Tools(context.Background())

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}

func TestTools(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"tools":[{"name":"current_time","description":"Get current date and time.","inputSchema":{"properties":{},"type":"object"},"mcp_info":{"server_name":"datetime-weather"}},{"name":"current_weather","description":"Get current weather. Use it always to get current weather information.","inputSchema":{"properties":{"city":{"description":"Name of the city","type":"string"}},"required":["city"],"type":"object"},"mcp_info":{"server_name":"datetime-weather"}}],"error": null,"message": "Successfully retrieved tools"}`))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}))
		defer server.Close()

		testUrl, err := url.Parse(server.URL)
		assert.NoError(t, err)

		c := getConfig()
		conn := getConn()
		conn.URL = *testUrl
		clientInstance := client.Litellm{Config: c, Connection: conn}

		res, err := clientInstance.Tools(context.Background())

		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Len(t, res, 2)
		assert.Equal(t, mcp.AvailableTools{
			{
				Name:        "current_time",
				Description: "Get current date and time.",
				InputSchema: mcp.AvailableToolInputSchema{Properties: map[string]mcp.AvailableToolProperty{}, Type: "object"},
				McpInfo:     mcp.AvailableToolMcpInfo{ServerName: "datetime-weather"}},
			{
				Name:        "current_weather",
				Description: "Get current weather. Use it always to get current weather information.",
				InputSchema: mcp.AvailableToolInputSchema{Properties: map[string]mcp.AvailableToolProperty{"city": {Description: "Name of the city", Type: "string"}}, Required: []string{"city"}, Type: "object"},
				McpInfo:     mcp.AvailableToolMcpInfo{ServerName: "datetime-weather"},
			},
		}, res)
	})
}
