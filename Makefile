.PHONY: build clean install test

# Binary name
BINARY_NAME=jfrm

# Build the application
build:
	go build -o $(BINARY_NAME) cmd/jfrm/main.go

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Install the application
install: build
	go install ./cmd/jfrm

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Generate documentation
docs:
	godoc -http=:6060

# Run in development mode
dev: build
	./$(BINARY_NAME) --help