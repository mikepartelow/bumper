.PHONY: lint fmt test local install image run-image

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

image: test
	@docker build -t bumper .

run-image: image
	@docker run --rm bumper
