-include .env

default:
	@go run ./cmd/$(APP_NAME)/main.go

build:
	@go build -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/main.go