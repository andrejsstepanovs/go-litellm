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

## Disclamer
Project is WIP as I'm poring it from another project atm. So if you see something weird or overengineered it is because this was part of different project. In the long run I will iron it out.
Please keep track of versioning.

## Contribution
Consider contributing. Reach out (or create PR) if you have questions or improvements.
