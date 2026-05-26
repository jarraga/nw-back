.PHONY: dev seed gen xls

dev:
	go run cmd/api/main.go

seed:
	go run cmd/seed/main.go

gen:
	sqlc generate

xls:
	go run cmd/export-xls/main.go $(if $(OUT),-out $(OUT),)
