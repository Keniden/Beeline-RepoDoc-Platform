.PHONY: lint test build docker

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

test:
	go test ./...

build:
	go build ./cmd/api
	go build ./cmd/worker

docker:
	docker build -t repodoc/api:latest .
