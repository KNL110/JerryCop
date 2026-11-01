.PHONY: run backends bench test build

# Run the JerryCop proxy
run:
	go run ./cmd/main.go

# Run the 3 mock backend servers (:8001, :8002, :8003)
backends:
	go run ./test/mock_backends/main.go

# Run latency benchmark (100 requests, 10 concurrency)
bench:
	go run ./test/latency_client/main.go

# Run tests
test:
	go test -v ./...

# Build binary
build:
	go build -o bin/jerrycop.exe ./cmd/main.go
