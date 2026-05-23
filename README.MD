# nw-back

Backend en Go para Northwind. Usa Chi para las rutas, pgx para conectarse a PostgreSQL y sqlc para generar codigo Go a partir de SQL.

## Requisitos

- Go instalado
- PostgreSQL corriendo localmente
- sqlc instalado como herramienta de desarrollo

## Instalar dependencias del proyecto

Desde la raiz del proyecto:

```powershell
go mod tidy
```

## Instalar sqlc

sqlc es una herramienta de linea de comandos. No se usa como dependencia runtime de la aplicacion.

```powershell
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Para verificar la instalacion:

```powershell
sqlc version
```

Si PowerShell no encuentra el comando, revisar que la carpeta `bin` de `GOPATH` este en el `PATH`.

```powershell
go env GOPATH
```

Normalmente el binario queda en:

```text
%USERPROFILE%\go\bin
```

## Configuracion

Crear un archivo `.env` en la raiz del proyecto:

```env
PORT=
PGHOST=
PGPORT=
PGUSER=
PGPASSWORD=
PGDATABASE=
```

El archivo `.env` esta ignorado por git.

## Ejecutar la aplicacion

Antes de iniciar el servidor HTTP, la aplicacion intenta conectarse a PostgreSQL. Si la conexion falla, se loguea el error y el proceso termina.

```powershell
go run .
```

Endpoint disponible:

```text
GET http://localhost:8080/
```

Respuesta esperada:

```json
{
  "message": "Northwind backend",
  "upTime": "10s"
}
```

## Estructura de Postgres

La carpeta `postgres` concentra lo relacionado con la base de datos:

```text
postgres/
  postgres.go      Conexion global con pgx
  schema/          Archivos SQL con estructura de tablas
  queries/         Consultas SQL usadas por sqlc
```

La conexion global queda disponible como:

```go
postgres.DB
```

## Usar sqlc

El archivo `sqlc.yaml` ya apunta a:

- `postgres/schema` para schemas
- `postgres/queries` para queries
- `postgres/db` como salida del codigo generado

Cuando existan archivos SQL en `schema` y `queries`, generar el codigo con:

```powershell
sqlc generate
```

## Verificar el proyecto

```powershell
go test ./...
```
