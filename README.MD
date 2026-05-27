# nw-back

Backend en Go para Northwind. Usa Chi para las rutas, pgx para conectarse a PostgreSQL, sqlc para generar codigo Go a partir de SQL y gofakeit para generar datos fake.

## Requisitos

- Go instalado
- PostgreSQL corriendo localmente
- sqlc instalado como herramienta de desarrollo
- make instalado para usar los comandos del `Makefile`

En Windows se puede instalar make con:

```powershell
winget install GnuWin32.Make
```

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
El comando `make dev` aplica las migraciones pendientes y, si todo sale bien, levanta la API.

```powershell
make dev
```

Para inicializar o actualizar solo la estructura de la base de datos:

```powershell
make migrate
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

## Generar datos fake

El comando de seed vive en `cmd/seed`. Ahi va todo lo relacionado con generacion de datos fake y carga inicial en la base de datos.
Al correrlo, borra las tablas actuales, vuelve a aplicar migraciones y carga datos demo.

La configuracion del seed esta en:

```text
internal/seed/config.go
```

```powershell
make seed
```

La API tambien expone un endpoint para reiniciar los datos demo desde el front:

```text
POST /admin/reset-demo-data
```

Para borrar todos los datos y dejar la base vacia con estructura creada:

```text
POST /admin/reset-data
```

Tambien se puede reemplazar la base con un Excel del mismo formato que genera la exportacion:

```text
POST /admin/import-xls
```

El archivo se envia como `multipart/form-data` en el campo `file`.

## Exportar clientes a Excel

El comando de exportacion genera un `.xlsx` con una pestaña por año. Cada pestaña incluye los clientes activos y 12 columnas mensuales. En cada mes aparece solo la fecha de pago si existe; si no hay pago registrado, la celda queda vacia.

```powershell
make xls
```

Por defecto guarda el archivo en `exports/`. Tambien se puede indicar una ruta:

```powershell
make xls OUT=exports/customers.xlsx
```

El mismo archivo tambien se puede pedir desde la API:

```text
GET /customers/export/xls
```

## Estructura del proyecto

```text
cmd/
  api/             Entrada principal de la API HTTP
  export-xls/      Entrada para exportar clientes y pagos a Excel
  seed/            Entrada para generar y cargar datos fake
internal/
  handlers/        Handlers HTTP
  routes/          Definicion de rutas
  postgres/        Conexion, schema, queries y codigo generado de Postgres
```

## Estructura de Postgres

La carpeta `internal/postgres` concentra lo relacionado con la base de datos:

```text
internal/postgres/
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

- `internal/postgres/schema` para schemas
- `internal/postgres/queries` para queries
- `internal/postgres/db` como salida del codigo generado

Cuando existan archivos SQL en `schema` y `queries`, generar el codigo con:

```powershell
make gen
```