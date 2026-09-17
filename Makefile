APP=users-service

.PHONY: dev test lint build infra-up infra-down

dev:
	go run ./cmd/users

test:
	go test ./...

lint:
	gofmt -w .
	go vet ./...

build:
	go build -o bin/$(APP) ./cmd/users

infra-up:
	docker compose up -d postgres

infra-down:
	docker compose down
