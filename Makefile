APP_NAME=api

.PHONY: run test fmt tidy build

run:
	go run ./cmd/$(APP_NAME)

test:
	go test ./...

fmt:
	gofmt -w $(shell find . -name '*.go' -not -path './vendor/*')

tidy:
	go mod tidy

build:
	mkdir -p bin
	go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)
