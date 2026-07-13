BINARY_NAME=palge

.PHONY: run
run:
	go run ./cmd/api

.PHONY: build
build:
	go build -o ./bin/$(BINARY_NAME) ./cmd/api

.PHONY: test
test:
	go test ./...

.PHONY: tidy
tidy:
	go mod tidy
