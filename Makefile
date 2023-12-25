-include .env

default:
	@go run ./cmd/formify/main.go

build:
	@go build -o ./bin/formify ./cmd/formify/main.go

test:
	@go test -v ./tests/...