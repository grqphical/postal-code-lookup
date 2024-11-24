build:
	@swag init -g ./cmd/api/main.go
	@go build -o main ./cmd/api

run:
	@swag init -g ./cmd/api/main.go
	@go run ./cmd/api

test:
	@go test ./...