.PHONY: run test test-unit up down

up:
	docker compose up -d

down:
	docker compose down

run:
	go run ./cmd/api

test-unit:
	go test ./internal/domain ./internal/service ./internal/handler/http -count=1

test:
	go test ./... -count=1 -timeout 5m
