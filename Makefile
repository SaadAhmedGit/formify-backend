-include .env
-include test.env

default:
	@go run ./cmd/$(APP_NAME)/main.go

build:
	@go build -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/main.go

test:
	@go test -v ./tests/...