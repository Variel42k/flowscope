SHELL := /bin/sh

.PHONY: up down logs seed test lint fmt

up:
	docker compose up --build -d

down:
	docker compose down -v

logs:
	docker compose logs -f --tail=200

seed:
	docker compose run --rm seed

test:
	go test ./...
	cd web && npm test -- --run

lint:
	golangci-lint run ./...
	cd web && npm run lint

fmt:
	gofmt -w $$(find cmd internal deploy -name '*.go')
	cd web && npm run format
