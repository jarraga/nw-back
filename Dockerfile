FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -tags netgo -ldflags "-s -w" -o /out/app ./cmd/api
RUN go build -tags netgo -ldflags "-s -w" -o /out/migrate ./cmd/migrate
RUN go build -tags netgo -ldflags "-s -w" -o /out/seed ./cmd/seed
RUN go build -tags netgo -ldflags "-s -w" -o /out/export-xls ./cmd/export-xls

FROM alpine:3.22

WORKDIR /app

ENV PORT=8080

# Provide PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDATABASE and PGSSLMODE
# at runtime with --env-file, docker compose environment, or the hosting platform.

COPY --from=build /out/app ./app
COPY --from=build /out/migrate ./migrate
COPY --from=build /out/seed ./seed
COPY --from=build /out/export-xls ./export-xls
COPY migrations ./migrations

EXPOSE 8080

CMD ["sh", "-c", "./migrate up && ./app"]
