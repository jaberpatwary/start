ifeq ($(OS),Windows_NT)
  export PATH := $(PATH);C:\Program Files\Go\bin
endif

.PHONY: start dev run build seed test tidy help

# Default start or dev runs Go Echo backend on :8090
start:
	go run ./cmd/api

dev:
	go run ./cmd/api

run:
	go run ./cmd/api

seed:
	go run ./cmd/seed

build:
	go build -o bin/api.exe ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy

help:
	@echo Server Commands:
	@echo   make start - Run Go Echo API server (:8090)
	@echo   make dev   - Run Go Echo API server (:8090)
	@echo   make seed  - Seed initial database data
	@echo   make build - Build Go binary to bin/api.exe
	@echo   make test  - Run tests
	@echo   make tidy  - Run go mod tidy
