build:
	go build -o bin/urlshort ./cmd/api/main.go

setup-dev-frontend:
	npm install --prefix ./frontend

dev-with-frontend:
	node frontend/index.js

dev-only-frontend:
	npm exec --prefix ./frontend start-hexlet-url-shortener-frontend

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

app_up:
	docker compose --profile full up -d

app_down:
	docker compose --profile full down

db_up:
	docker compose up db -d

db_down:
	docker compose down db

db_connect:
	docker compose exec db psql -U postgres -d appdb

db_migrate:
	go tool \
	goose -dir ./db/migrations postgres "${DATABASE_URL}" \
	up
	
sqlc_generate:
	go tool sqlc generate

.PHONY: build \
	dev-with-frontend dev-only-frontend setup-dev-frontend \
	lint lint-fix format \
	test coverage coverage-html \
	app_up app_down \
	db_up db_down db_connect db_migrate \
	sqlc_generate
