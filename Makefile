build:
	go build -o bin/urlshort main.go

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix
	
format:
	golangci-lint fmt
	
test:
	go test -race ./...

.PHONY: build lint lint-fix format test
