# Makefile for a Golang project

# Variables
APP_NAME = b.exe
SRC_FILE = b.go
BINARY_NAME = $(APP_NAME)

# Default target
all: build

# Build the binary
build:
	@echo "Building the binary..."
	go build -o $(BINARY_NAME) $(SRC_FILE)

# Run the application
run: build
	@echo "Running the application..."
	./$(BINARY_NAME)

# Test the application
test:
	@echo "Running tests..."
	go test ./...

# Clean up generated files
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

# Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint the code
lint:
	@echo "Linting code..."
	golangci-lint run

# Tidy up dependencies
tidy:
	@echo "Tidying up dependencies..."
	go mod tidy

# Vendor dependencies
vendor:
	@echo "Vendoring dependencies..."
	go mod vendor

.PHONY: all build run test clean fmt lint tidy vendor
