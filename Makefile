.PHONY: dev seed gen

dev:
	go run cmd/api/main.go

seed:
	go run cmd/seed/main.go

gen:
	sqlc generate
