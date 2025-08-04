package client_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/client"
	"github.com/andrejsstepanovs/go-litellm/request"
)

func TestTokenCounter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping token counter test in short mode")
	}

	t.Run("functional", func(t *testing.T) {
		clientInstance := client.Litellm{Config: getConfig(), Connection: getConn()}
		ctx := context.Background()

		messages := request.Messages{}
		messages.AddMessage(request.UserMessageSimple("Hi"))

		modelMeta, err := clientInstance.Model(ctx, testModel)
		assert.NoError(t, err)

		req := &request.TokenCounterRequest{
			Model:    modelMeta.ModelId,
			Messages: messages,
		}

		resp, err := clientInstance.TokenCounter(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Greater(t, resp.TotalTokens, float64(0))
	})
}
