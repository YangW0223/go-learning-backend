.PHONY: run test fmt tidy migrate-up migrate-down docker-up docker-down docker-logs

run:
	go run ./cmd/server

test:
	go test ./...

fmt:
	gofmt -w $(shell find . -name '*.go' -not -path './vendor/*')

tidy:
	go mod tidy

migrate-up:
	go run ./scripts/migrate -direction up

migrate-down:
	go run ./scripts/migrate -direction down

docker-up:
	docker compose --env-file docker/.env.example -f docker/docker-compose.yaml up --build

docker-down:
	docker compose --env-file docker/.env.example -f docker/docker-compose.yaml down -v

docker-logs:
	docker compose --env-file docker/.env.example -f docker/docker-compose.yaml logs -f api postgres redis
