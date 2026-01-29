.PHONY: build generate validate test lint clean

# Build the CLI tool
build:
	go build -o bin/crx-registry ./cmd/crx-registry

# Generate registry.yaml from pkgs/
generate: build
	./bin/crx-registry generate

# Validate all packages
validate: build
	./bin/crx-registry validate

# Run tests
test:
	go test -race -shuffle=on ./...

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/

# Install the tool locally
install:
	go install ./cmd/crx-registry
