.PHONY: test

lint:
	@golangci-lint run ./...

fmt: lint
	@go fmt ./...

test: fmt
	@go test -race -cover ./...

local: test
	@go build -o . ./cmd/...

install: test
	@go install ./cmd/...
