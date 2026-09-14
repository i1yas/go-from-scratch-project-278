lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix
	
format:
	golangci-lint fmt
	
test:
	go test -race ./...

.PHONY: lint lint-fix format test
