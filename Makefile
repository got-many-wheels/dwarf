DATABASE_NAME ?= dwarf
DATABASE_PORT ?= 5432
DATABASE_USER ?= postgres
DATABASE_PASSWORD ?= root

scdump:
	PGPASSWORD=$(DATABASE_PASSWORD) pg_dump -h localhost -p $(DATABASE_PORT) -s -U $(DATABASE_USER) $(DATABASE_NAME) > migrations/schema.sql

generate:
	cd server && go generate ./...

schema: scdump generate

migrate-up:
	migrate -source file://./migrations -database $(DATABASE_URI) up

migrate-down:
	migrate -source file://./migrations -database $(DATABASE_URI) down

run-server:
	cd server && \
	export DATABASE_URI="postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@localhost:$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable" && \
	export PORT=":8080" && \
	go run ./cmd

.PHONY: scdump generate migrate-up migrate-down
