.PHONY: dev migrate seed gen xls docker-build docker-run

dev:
	go run ./cmd/migrate up && go run ./cmd/api

migrate:
	go run ./cmd/migrate up

seed:
	go run ./cmd/seed

gen:
	sqlc generate

xls:
	go run ./cmd/export-xls $(if $(OUT),-out $(OUT),)

docker-build:
	docker build -t nw-back .

docker-run:
	docker run --rm --env-file .env -p 8080:8080 nw-back
