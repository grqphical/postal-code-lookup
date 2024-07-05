build:
	@go build -o main ./cmd/api

run:
	@go run ./cmd/api

test:
	@go test ./...