.PHONY: dev migrate seed gen xls

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
