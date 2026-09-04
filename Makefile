ifeq ($(OS),Windows_NT)
  export PATH := $(PATH);C:\Program Files\nodejs;C:\Program Files\Go\bin
endif

.PHONY: start dev server client seed build tidy test help

# Run both backend (:8090) and frontend (:8888) concurrently
start:
	npm run dev

dev:
	npm run dev

# Run Go backend only
server:
	cd server && go run ./cmd/api

# Run React frontend only
client:
	cd client && npm run dev

# Seed initial database records
seed:
	cd server && go run ./cmd/seed

# Build both backend binary and frontend bundle
build:
	npm run build:server
	npm run build:client

tidy:
	cd server && go mod tidy

test:
	cd server && go test ./...

help:
	@echo StarTech Commands:
	@echo   make start   - Run Backend (:8090) and Frontend (:8888) together
	@echo   make dev     - Run Backend and Frontend together
	@echo   make server  - Run Go Echo backend only (:8090)
	@echo   make client  - Run React Vite frontend only (:8888)
	@echo   make seed    - Seed initial data into PostgreSQL
	@echo   make build   - Build backend binary and frontend bundle
