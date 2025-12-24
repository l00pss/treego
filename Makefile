.PHONY: test bench clean fmt vet lint build example

# Test the project
test:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
bench:
	go test -bench=. -benchmem ./...

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Run example
example:
	cd example && go run main.go

# Build example
build-example:
	cd example && go build -o treego-example main.go

# Clean generated files
clean:
	rm -f coverage.out coverage.html
	rm -f example/treego-example

# Run all checks
check: fmt vet test

# Install dependencies and tools
deps:
	go mod tidy
	go mod download

# Generate documentation
docs:
	@which godoc > /dev/null 2>&1 || { echo "Installing godoc..."; go install golang.org/x/tools/cmd/godoc@latest; }
	godoc -http=:6060

# Show package documentation
doc:
	go doc -all .

# Run all quality checks
quality: fmt vet test bench

# Help
help:
	@echo "Available targets:"
	@echo "  test         - Run tests"
	@echo "  test-cover   - Run tests with coverage report"
	@echo "  bench        - Run benchmarks"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  example      - Run example"
	@echo "  build-example- Build example binary"
	@echo "  clean        - Clean generated files"
	@echo "  check        - Run format, vet, and tests"
	@echo "  deps         - Install dependencies"
	@echo "  docs         - Start documentation server"
	@echo "  doc          - Show package documentation"
	@echo "  quality      - Run all quality checks"
	@echo "  help         - Show this help"
