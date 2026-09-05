ifeq ($(OS),Windows_NT)
  export PATH := $(PATH);C:\Program Files\nodejs;C:\Program Files\Go\bin
endif

.PHONY: start dev server backend client frontend seed build tidy test help

# Run both backend (:8090) and frontend (:8888) concurrently
start:
	npm run dev

dev:
	npm run dev

# Run Go backend only
backend:
	cd backend && go run ./cmd/api

server:
	cd backend && go run ./cmd/api

# Run React frontend only
frontend:
	cd frontend && npm run dev

client:
	cd frontend && npm run dev

# Seed initial database records
seed:
	cd backend && go run ./cmd/seed

# Build both backend binary and frontend bundle
build:
	npm run build:backend
	npm run build:frontend

tidy:
	cd backend && go mod tidy

test:
	cd backend && go test ./...

help:
	@echo StarTech Commands:
	@echo   make start    - Run Backend (:8090) and Frontend (:8888) together
	@echo   make dev      - Run Backend and Frontend together
	@echo   make backend  - Run Go Echo backend only (:8090)
	@echo   make frontend - Run React Vite frontend only (:8888)
	@echo   make seed     - Seed initial data into PostgreSQL
	@echo   make build    - Build backend binary and frontend bundle
