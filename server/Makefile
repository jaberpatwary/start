.PHONY: dev build test tidy

dev:
	go run ./src/main.go

build:
	go build -o bin/server ./src/main.go

test:
	go test ./...

tidy:
	go mod tidy
