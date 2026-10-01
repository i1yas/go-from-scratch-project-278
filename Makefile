BINARY := bin/urlshort
GOOSE := go tool goose -dir ./db/migrations postgres "${DATABASE_URL}"

build:
	go build -o $(BINARY) ./cmd/api/main.go

run:
	go run cmd/api/main.go
	
clean:
	rm $(BINARY) coverage.out

##@ Quality check
lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix
	
format:
	golangci-lint fmt
	
test:
	go test -race -coverpkg=./...  -coverprofile=coverage.out ./...

coverage: test
	go tool cover -func=coverage.out

coverage-html: test
	go tool cover -html=coverage.out
	
##@ Frontend
frontend-setup:
	npm install --prefix ./frontend

dev:
	node frontend/index.js

dev-frontend:
	npm exec --prefix ./frontend start-hexlet-url-shortener-frontend

##@ Docker app
up:
	docker compose --profile full up -d

down:
	docker compose --profile full down

##@ Database
db-up:
	docker compose up db -d

db-down:
	docker compose down db

db-remove:
	docker compose down db -v

db-shell:
	docker compose exec db psql -U postgres -d appdb

db-migrate:
	$(GOOSE) up

db-rollback:
	$(GOOSE) down
	
db-seed:
	go run ./cmd/seed/main.go
	
sqlc-generate:
	go tool sqlc generate

.PHONY: build run clean \
	frontend-setup dev dev-frontend \
	lint lint-fix format \
	test coverage coverage-html \
	up down \
	db-up db-down db-remove db-connect db-migrate \
	sqlc-generate
