APP_NAME=todolist-api

.PHONY: run migrate-up migrate-down migrate-reset migrate-version swagger

run:
	go run ./cmd/api

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-reset:
	go run ./cmd/migrate reset

migrate-version:
	go run ./cmd/migrate version

swagger:
	swag init -g cmd/api/main.go -o docs
