PACKAGES := $(shell find . -name *.go | grep -v -E "vendor|tools|mocks" | xargs -n1 dirname | sort -u)

.PHONY: test-unit
test-unit:
	go test -short ./...

.PHONY: test
test:
	go test ./...

.PHONY: install
install:
	go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go mod tidy

.PHONY: lint
lint:
	go tool golangci-lint run $(PACKAGES)