# Define variables
APP_NAME := wallet-service

# Default target
all: build

# Build the Go application
build:
	@echo "Building the application..."
	go build -o $(APP_NAME) ./cmd/api-service 

# Run the application
run:
	@echo "Running the application..."
	@export $(grep -v '^#' .env | xargs) && go run ./cmd/api/wallet/main.go

# Ensure all imports are satisfied
tidy:
	go mod tidy

# Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...

.PHONY: all build clean run proto test fmt lint


# Run migrate up 
migrate:
	@echo "Running migrate up..."
	dbmate -d "internal/db/migrations" --wait --no-dump-schema up

# swagger commands
check-swagger:
	which swagger || go install github.com/go-swagger/go-swagger/cmd/swagger@latest

swagger: check-swagger
	swagger generate spec -o ./docs/swagger.yaml --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger ./docs/swagger.yaml

# Run the tests
test:
	@echo "Running tests..."
	go test -v ./...