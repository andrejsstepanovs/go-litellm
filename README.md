# go-litellm

This package provides client that communicates with LiteLLM server API.

It covers most litellm functionality.

## Functionality
Package implements following LiteLLM endpoints with request and response module structs.

- /models
- /v2/model/info
- /model_group/info
- /utils/token_counter
- /v1/embeddings
- /audio/transcriptions
- /mcp-rest/tools/list
- /mcp-rest/tools/call
- /chat/completions (with image, strict schema and mcp)


## Installation

```bash
go get github.com/andrejsstepanovs/go-litellm@v0.1.0
```

## Examples

Most basic example. (will add more later)
```go
package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/andrejsstepanovs/go-litellm/client"
	"github.com/andrejsstepanovs/go-litellm/conf/connections/litellm"
	"github.com/andrejsstepanovs/go-litellm/models"
	"github.com/andrejsstepanovs/go-litellm/request"
)

func main() {
	baseURL, _ := url.Parse("http://localhost:4000")
	conn := litellm.Connection{
		URL: *baseURL,
		Targets: litellm.Targets{
			System: litellm.Target{Timeout: time.Second}, // litellm internal rest api calls (get models, etc)
			LLM:    litellm.Target{Timeout: time.Minute * 2},
			MCP:    litellm.Target{Timeout: time.Minute * 5},
		},
	}

    // different users can own different api keys, so its separate param.
	cfg := client.Config{
		APIKey:      "sk-1234",
		Temperature: 1,
	}

	ai, err := client.New(cfg, conn)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	model, err := ai.Model(ctx, models.ModelID("claude-4"))
	if err != nil {
		panic(err)
	}

	prompt := "What is the capital of France?"
	messages := request.Messages{request.UserMessageSimple(prompt)}
	req := request.NewCompletionRequest(model, messages, nil, nil, 1)

	resp, err := ai.Completion(ctx, req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.String()) // The capital of France is Paris.
}
```

```bash
go mod tidy
go run main.go
```

## Disclamer
Project is WIP as I'm poring it from another project atm. So if you see something weird or overengineered it is because this was part of different project. In the long run I will iron it out.
Please keep track of versioning.

## Contribution
Consider contributing. Reach out (or create PR) if you have questions or improvements.
