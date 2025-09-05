# go-litellm

A **Go client library** for interacting with the [LiteLLM](https://github.com/BerriAI/litellm) API. This package provides a simple, type-safe, and developer-friendly way to perform completions, audio transcription, image-based queries, structured JSON output, embeddings, token counting, MCP tool operations, model browsing, and more.

---

## Installation

```bash
go get github.com/andrejsstepanovs/go-litellm@latest
```

---

## Quick Start

### Initialize Client

```go
package main

import (
    "context"
    "net/url"
    "time"
    "log"

    "github.com/andrejsstepanovs/go-litellm/client"
    "github.com/andrejsstepanovs/go-litellm/conf/connections/litellm"
)

func main() {
    baseURL, _ := url.Parse("http://localhost:4000")
    conn := litellm.Connection{
        URL: *baseURL,
        Targets: litellm.Targets{
            System: litellm.Target{Timeout: time.Second},
            LLM:    litellm.Target{Timeout: time.Minute * 2},
            MCP:    litellm.Target{Timeout: time.Minute * 5},
        },
    }

    cfg := client.Config{
        APIKey:      "sk-1234",
        Temperature: 0.7,
    }

    ai, err := client.New(cfg, conn)
    if err != nil {
        log.Fatal(err)
    }

    _ = ai // ready to use
}
```

---

## Examples

### 1. Simple Chat Completion

```go
model, _ := ai.Model(ctx, "claude-4")
messages := request.Messages{request.UserMessageSimple("What is the capital of France?")}
req := request.NewCompletionRequest(model, messages, nil, nil, 1)
resp, _ := ai.Completion(ctx, req)
fmt.Println(resp.String()) // The capital of France is Paris.
```

### 2. Audio Transcription (Speech-to-Text)

```go
model, _ := ai.Model(ctx, "whisper-1")
res, _ := ai.SpeechToText(ctx, model, "path/to/audio.oga")
fmt.Println(res.Text)
```

### 3. Image Analysis / Captioning

```go
model, _ := ai.Model(ctx, "gpt-4o-mini")
messages := request.Messages{
    request.UserMessageImage("Describe this image", "https://example.com/image.jpg"),
}
req := request.NewCompletionRequest(model, messages, nil, nil, 1)
resp, _ := ai.Completion(ctx, req)
fmt.Println(resp.String())
```

### 4. Structured JSON Output (Strict Schema)

```go
type City struct {
    CityName        string `json:"city_name"`
    PopulationCount int    `json:"population_count"`
}
type ListOfCities struct {
    Cities []City `json:"cities"`
}

schema := request.JSONSchema{
    Name: "list_of_cities",
    Schema: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "cities": map[string]interface{}{
                "type": "array",
                "items": map[string]interface{}{
                    "type": "object",
                    "properties": map[string]interface{}{
                        "city_name": map[string]interface{}{"type": "string"},
                        "population_count": map[string]interface{}{"type": "integer"},
                    },
                    "required": []string{"city_name", "population_count"},
                },
            },
        },
        "required": []string{"cities"},
    },
    Strict: true,
}

model, _ := ai.Model(ctx, "claude-4")
messages := request.Messages{request.UserMessageSimple("List the 3 largest cities")}
req := request.NewCompletionRequest(model, messages, nil, nil, 0.2)
req.SetJSONSchema(schema)
resp, _ := ai.Completion(ctx, req)

var cities ListOfCities
json.Unmarshal(resp.Bytes(), &cities)
fmt.Printf("%+v\n", cities)
```

### 5. Token Count Calculation

```go
model, _ := ai.Model(ctx, "claude-4")
messages := request.Messages{request.UserMessageSimple("Hello")}
req := &request.TokenCounterRequest{
    Model:    model.ModelId,
    Messages: messages,
}
count, _ := ai.TokenCounter(ctx, req)
fmt.Println("Total tokens:", count.TotalTokens)
```

### 6. List Available MCP Tools

```go
tools, _ := ai.Tools(ctx)
for _, tool := range tools {
    fmt.Printf("%s - %s\n", tool.Name, tool.Description)
}
```

### 7. Call an MCP Tool

```go
tool := common.ToolCallFunction{
    Name: "current_time",
    Arguments: map[string]string{"timezone": "Europe/Riga"},
}
res, _ := ai.ToolCall(ctx, tool)
fmt.Println(res.String())
```

### 8. Browse Available Models

```go
models, _ := ai.Models(ctx)
for _, m := range models {
    fmt.Println(m.ID, m.OwnedBy)
}
```

### 9. Tool-Aware Conversation Example

This example demonstrates how to:

1. Maintain a conversation history.
2. Send the list of available MCP tools so the AI knows what it can call.
3. When AI requests a tool call, execute it, append the result to the history, and send it back.
4. Repeat until AI no longer requests tools.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/url"
    "time"

    "github.com/andrejsstepanovs/go-litellm/client"
    "github.com/andrejsstepanovs/go-litellm/conf/connections/litellm"
    "github.com/andrejsstepanovs/go-litellm/mcp"
    "github.com/andrejsstepanovs/go-litellm/models"
    "github.com/andrejsstepanovs/go-litellm/request"
    "github.com/andrejsstepanovs/go-litellm/response"
)

func main() {
    ctx := context.Background()

    // Create LiteLLM client
    conn := litellm.Connection{
        URL: *mustParseURL("http://localhost:4000"),
        Targets: litellm.Targets{
            System: litellm.Target{Timeout: time.Second},
            LLM:    litellm.Target{Timeout: time.Minute},
            MCP:    litellm.Target{Timeout: time.Minute},
        },
    }
    cfg := client.Config{APIKey: "sk-1234", Temperature: 0.7}
    ai, err := client.New(cfg, conn)
    if err != nil {
        log.Fatal(err)
    }

    // Pick a model and list tools
    model, _ := ai.Model(ctx, "claude-4")
    tools, _ := ai.Tools(ctx)

    // Initial conversation
    messages := request.Messages{
        request.UserMessageSimple("What's the current time in Riga?"),
    }

    finalResp := runToolAwareConversation(ctx, ai, model, tools, messages, 5)
    fmt.Println("Final Answer:", finalResp.String())
}

func runToolAwareConversation(ctx context.Context, ai *client.Litellm, model models.ModelMeta, tools mcp.AvailableTools, messages request.Messages, maxIter int) response.Response {
    if maxIter <= 0 {
        return response.Response{}
    }

    req := request.NewCompletionRequest(model, messages, tools.ToLLMCallTools(), nil, 0.7)
    resp, err := ai.Completion(ctx, req)
    if err != nil {
        log.Fatal("Completion error:", err)
    }

    if resp.Choice().FinishReason == response.FINISH_REASON_TOOL {
        for _, toolCall := range resp.Choice().Message.ToolCalls.SortASC() {
            toolResp, err := ai.ToolCall(ctx, toolCall.Function)
            if err != nil {
                log.Fatal("Tool call error:", err)
            }
            for _, tr := range toolResp {
                messages = append(messages, request.ToolCallMessage(toolCall, tr))
            }
        }
        return runToolAwareConversation(ctx, ai, model, tools, messages, maxIter-1)
    }

    return resp
}

func mustParseURL(s string) *url.URL {
    u, err := url.Parse(s)
    if err != nil {
        log.Fatal(err)
    }
    return u
}
```

---

## Supported Endpoints

* `/models` – list available models
* `/v2/model/info` – detailed model info
* `/model_group/info` – fetch model metadata
* `/utils/token_counter` – count tokens for a given request
* `/v1/embeddings` – generate embeddings
* `/audio/transcriptions` – speech-to-text
* `/mcp-rest/tools/list` – list MCP tools
* `/mcp-rest/tools/call` – invoke MCP tools
* `/chat/completions` – chat completions with support for text, images, strict schemas, and MCP integration

---

## Contributing

Contributions are welcome! Feel free to open an issue or submit a PR.

## License

Apache 2.0
