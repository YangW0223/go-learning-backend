APP_NAME=api

.PHONY: run test fmt tidy build docker-up docker-down docker-logs

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

docker-up:
	docker compose --env-file docker/.env.example -f docker/docker-compose.redis.yaml up --build

docker-down:
	docker compose --env-file docker/.env.example -f docker/docker-compose.redis.yaml down -v

docker-logs:
	docker compose --env-file docker/.env.example -f docker/docker-compose.redis.yaml logs -f api redis
