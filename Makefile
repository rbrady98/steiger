# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	
	go build -o build/main cmd/api/main.go

# Run the application
run:
	@set -a && source .env && set +a && go run cmd/api/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary and the test.db
clean:
	@echo "Cleaning..."
	@rm -f main test.db

# Live Reload
dev:
	@air

lint:
	@golangci-lint run

sqlcgen:
	@sqlc generate

.PHONY: all build run test clean dev sqlcgen
