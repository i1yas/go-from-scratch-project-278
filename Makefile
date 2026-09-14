lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix
	
format:
	golangci-lint fmt

.PHONY: lint lint-fix format
